package target_sync

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"

	"github.com/raito-io/cli/base/util/error/grpc_error"
	plugin2 "github.com/raito-io/cli/base/util/plugin"
	"github.com/raito-io/cli/internal/constants"
	gql "github.com/raito-io/cli/internal/graphql"
	"github.com/raito-io/cli/internal/job"
	"github.com/raito-io/cli/internal/logging"
	"github.com/raito-io/cli/internal/plugin"
	"github.com/raito-io/cli/internal/target"
	"github.com/raito-io/cli/internal/target/types"
	"github.com/raito-io/cli/internal/util/array"
	"github.com/raito-io/cli/internal/version_management"
)

type SyncJob struct {
	jobIds      []string
	RunTypeName string
}

func (s *SyncJob) RunType() string {
	return s.RunTypeName
}

func (s *SyncJob) TargetSync(ctx context.Context, targetConfig *types.BaseTargetConfig) (syncError error) {
	targetConfig.TargetLogger.Info("Executing target...")

	var jobId string

	start := time.Now()

	defer func() {
		if syncError != nil {
			targetConfig.TargetLogger.Error(fmt.Sprintf("Failed execution: %s", syncError.Error()), "success")
		} else {
			targetConfig.TargetLogger.Info(fmt.Sprintf("Successfully finished execution in %s", time.Since(start).Round(time.Millisecond)), "success")

			if jobId != "" {
				s.jobIds = append(s.jobIds, jobId)
			}
		}
	}()

	client, err := plugin.NewPluginClient(targetConfig.ConnectorName, targetConfig.ConnectorVersion, targetConfig.TargetLogger)
	if err != nil {
		targetConfig.TargetLogger.Error(fmt.Sprintf("Error initializing connector plugin %q: %s", targetConfig.ConnectorName, err.Error()))
		return err
	}
	defer client.Close()

	jobId, err = job.StartJob(ctx, targetConfig)
	if err != nil {
		return err
	}

	targetConfig.TargetLogger.Info(fmt.Sprintf("Start job with jobID: '%s'", jobId))
	job.UpdateJobEvent(targetConfig, jobId, job.InProgress, nil)

	defer func() {
		if syncError == nil {
			job.UpdateJobEvent(targetConfig, jobId, job.Completed, nil)
		} else {
			job.UpdateJobEvent(targetConfig, jobId, job.Failed, syncError)
		}
	}()

	pluginInfoClient, err := client.GetInfo()
	if err != nil {
		return fmt.Errorf("get plugin info client: %w", err)
	}

	pluginInfo, err := pluginInfoClient.GetInfo(ctx)
	if err != nil {
		return fmt.Errorf("get plugin info: %w", err)
	}

	if len(pluginInfo.Type) == 0 { // Backwards compatibility
		err = dsSyncTargetSync(ctx, targetConfig, client, jobId)
		if err != nil {
			return fmt.Errorf("full data source sync: %w", err)
		}
	}

	for _, pluginType := range pluginInfo.Type {
		switch pluginType {
		case plugin2.PluginType_PLUGIN_TYPE_FULL_DS_SYNC:
			err = dsSyncTargetSync(ctx, targetConfig, client, jobId)
			if err != nil {
				return fmt.Errorf("full data source sync: %w", err)
			}
		case plugin2.PluginType_PLUGIN_TYPE_IS_SYNC:
			err = isSyncTargetSync(ctx, targetConfig, client, jobId)
			if err != nil {
				return fmt.Errorf("identity store sync: %w", err)
			}
		case plugin2.PluginType_PLUGIN_TYPE_TAG_SYNC:
			err = tagSyncTargetSync(ctx, targetConfig, client, jobId)
			if err != nil {
				return fmt.Errorf("tag sync: %w", err)
			}
		case plugin2.PluginType_PLUGIN_TYPE_RESOURCE_PROVIDER:
			err = resourceProviderSync(ctx, targetConfig, client, jobId)
			if err != nil {
				return fmt.Errorf("resource provider sync: %w", err)
			}
		default:
			return fmt.Errorf("unsupported plugin type: %s", pluginType)
		}
	}

	return nil
}

func (s *SyncJob) Finalize(ctx context.Context, baseConfig *types.BaseConfig, options *target.Options) error {
	return sendEndOfTarget(ctx, baseConfig, s.jobIds, options)
}

func execute(ctx context.Context, targetID string, jobID string, syncType string, syncTypeLabel string, skipSync bool, syncTask job.Task, cfg *types.BaseTargetConfig, c plugin.PluginClient) (err error) {
	cfg, warningCollector, loggingCleanUp, err := logging.CreateWarningCapturingLogger(cfg)
	if err != nil {
		return err
	}

	defer loggingCleanUp()

	defer func() {
		if r := recover(); r != nil {
			cfg.TargetLogger.Error(fmt.Sprintf("Panic occurred during %s sync: %v", syncTypeLabel, r))

			err = fmt.Errorf("panic occurred during %s sync", syncTypeLabel)
		}
	}()

	taskEventUpdater := job.NewTaskEventUpdater(cfg, jobID, syncType, warningCollector)

	switch {
	case skipSync:
		taskEventUpdater.SetStatusToSkipped(ctx)
		cfg.TargetLogger.Info("Skipping sync of " + syncTypeLabel)
	case targetID == "":
		taskEventUpdater.SetStatusToSkipped(ctx)

		idField := "data-source-id"
		if syncType == constants.IdentitySync {
			idField = "identity-store-id"
		}

		cfg.TargetLogger.Warn("No " + idField + " argument found. Skipping syncing of " + syncTypeLabel)
	default:
		syncErr := sync(ctx, cfg, syncTypeLabel, taskEventUpdater, syncTask, c, syncType, jobID)
		if syncErr != nil {
			// Sync error is already pushed to task error
			return fmt.Errorf("failed to execute %s sync", syncTypeLabel)
		}
	}

	return nil
}

func logForwardingEnabled(syncType string) bool {
	if viper.GetBool(constants.DisableLogForwarding) {
		return false
	}

	cmdFlag := ""

	switch syncType {
	case constants.DataSourceSync:
		cmdFlag = constants.DisableLogForwardingDataSourceSync
	case constants.IdentitySync:
		cmdFlag = constants.DisableLogForwardingIdentityStoreSync
	case constants.DataAccessSync:
		cmdFlag = constants.DisableLogForwardingDataAccessSync
	case constants.DataUsageSync:
		cmdFlag = constants.DisableLogForwardingDataUsageSync
	case constants.ResourceProviderSync:
		cmdFlag = constants.DisableLogForwardingResourceProviderSync
	case constants.TagSync:
		cmdFlag = constants.DisableLogForwardingTagSync
	}

	return !viper.GetBool(cmdFlag)
}

func sync(ctx context.Context, cfg *types.BaseTargetConfig, syncTypeLabel string, taskEventUpdater job.TaskEventUpdater, syncTask job.Task, c plugin.PluginClient, syncType string, jobID string) (err error) {
	defer func() {
		if err != nil {
			taskEventUpdater.SetStatusToFailed(ctx, err)

			target.HandleTargetError(err, cfg, fmt.Sprintf("Synchronizing %s failed", syncType))
		}
	}()

	if logForwardingEnabled(syncType) {
		targetCfg, cleanup, taskLoggingError := logging.CreateTaskLogger(cfg, jobID, syncType)
		if taskLoggingError != nil {
			return taskLoggingError
		}

		cfg = targetCfg

		defer func() {
			cleanUpErr := cleanup()
			if cleanUpErr != nil {
				cfg.TargetLogger.Warn(fmt.Sprintf("Failed to close logger for task: %s", cleanUpErr.Error()))
			}
		}()
	}

	_, err = syncTask.IsClientValid(ctx, c)
	incompatibleVersionError := version_management.IncompatiblePluginVersionError{}

	var internalPluginStatusError *grpc_error.InternalPluginStatusError

	if errors.As(err, &internalPluginStatusError) {
		if internalPluginStatusError.StatusCode() == codes.Unimplemented {
			cfg.TargetLogger.Info(fmt.Sprintf("Plugin does not implement a syncer for %s. Skipping", syncTypeLabel))
			taskEventUpdater.SetStatusToSkipped(ctx) // Skip should be sent before we send a start status event

			return nil
		}
	}

	cfg.TargetLogger.Info(fmt.Sprintf("Synchronizing %s...", syncTypeLabel))
	taskEventUpdater.SetStatusToStarted(ctx)

	if errors.As(err, &incompatibleVersionError) {
		return fmt.Errorf("unable to execute %s sync: %w", syncTypeLabel, incompatibleVersionError)
	} else if err != nil {
		return err
	}

	syncParts := syncTask.GetParts()

	for i, taskPart := range syncParts {
		err2 := runTaskPartSync(ctx, cfg, syncTypeLabel, taskEventUpdater, jobID, syncType, taskPart, i, syncParts, c)
		if err2 != nil {
			return err2
		}
	}

	taskEventUpdater.SetStatusToCompleted(ctx, syncTask.GetTaskResults())

	return nil
}

func runTaskPartSync(ctx context.Context, cfg *types.BaseTargetConfig, syncTypeLabel string, taskEventUpdater job.TaskEventUpdater, jobID string, syncType string, taskPart job.TaskPart, i int, syncParts []job.TaskPart, c plugin.PluginClient) error {
	cfg.TargetLogger.Debug(fmt.Sprintf("Start sync task part %d out of %d", i+1, len(syncParts)))

	status, subtaskId, err := taskPart.StartSyncAndQueueTaskPart(ctx, c, taskEventUpdater)
	if err != nil {
		err = fmt.Errorf("synchronizing %s : %w", syncType, err)

		return err
	}

	if status == job.Queued {
		taskEventUpdater.SetStatusToQueued(ctx)
		cfg.TargetLogger.Info(fmt.Sprintf("Waiting for server to start processing %s...", syncTypeLabel))
	}

	syncResult := taskPart.GetResultObject()

	if syncResult != nil {
		subtask, err := job.WaitForJobToComplete(ctx, jobID, syncType, subtaskId, syncResult, cfg, status)
		if err != nil {
			return err
		}

		if subtask.Status == job.Failed {
			var subtaskErr error
			subtaskErr = multierror.Append(subtaskErr, array.Map(subtask.Errors, func(err *string) error { return errors.New(*err) })...)

			return subtaskErr
		} else if subtask.Status == job.TimeOut {
			return fmt.Errorf("synchronizing %s timed out", syncType)
		}

		err = taskPart.ProcessResults(syncResult)
		if err != nil {
			return err
		}
	} else if status != job.Completed {
		return fmt.Errorf("unable to load results")
	}

	return nil
}

func sendEndOfTarget(ctx context.Context, baseConfig *types.BaseConfig, jobIds []string, options *target.Options) error {
	if !options.ExternalTrigger {
		var mutation struct {
			EndOfTargetsSyncResult struct {
				EndOfTargetsSync struct {
					Success bool
				} `graphql:"... on EndOfTargetsSync"`
				PermissionDeniedError struct {
					message string
				} `graphql:"... on PermissionDeniedError"`
			} `graphql:"endOfTargetsSync(input: $input)"`
		}

		type EndOfTargetsSyncInput struct {
			JobIds []string `json:"jobIds"`
		}

		input := EndOfTargetsSyncInput{JobIds: jobIds}

		err := gql.NewClient(baseConfig).Mutate(ctx, &mutation, map[string]interface{}{"input": input})
		if err != nil {
			return err
		}

		if mutation.EndOfTargetsSyncResult.PermissionDeniedError.message != "" {
			baseConfig.BaseLogger.Error(fmt.Sprintf("Permission denied to notify end of all targets: %s", mutation.EndOfTargetsSyncResult.PermissionDeniedError.message))
		} else if !mutation.EndOfTargetsSyncResult.EndOfTargetsSync.Success {
			baseConfig.BaseLogger.Warn("Failed to notify end of all targets")
		}
	}

	return nil
}
