package target_sync

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"

	"github.com/raito-io/cli/base/util/error/grpc_error"
	plugin2 "github.com/raito-io/cli/base/util/plugin"
	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/error_handler"
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
	runTypeName string

	syncStartMutex sync.Mutex
	syncQueueMutex sync.Mutex
	wg             sync.WaitGroup
}

func NewSyncJob(runTypeName string) *SyncJob {
	return &SyncJob{
		runTypeName: runTypeName,
	}
}

func (s *SyncJob) RunType() string {
	return s.runTypeName
}

func (s *SyncJob) TargetSync(ctx context.Context, logger hclog.Logger, targetConfig *types.BaseTargetConfig) (syncError error) {
	logger.Info("Executing target...")

	var jobId string

	start := time.Now()

	defer func() {
		if syncError != nil {
			logger.Error(fmt.Sprintf("Failed execution: %s", syncError.Error()), "success")
		} else {
			logger.Info(fmt.Sprintf("Successfully finished execution in %s", time.Since(start).Round(time.Millisecond)), "success")

			if jobId != "" {
				s.jobIds = append(s.jobIds, jobId)
			}
		}
	}()

	client, err := plugin.NewPluginClient(targetConfig.ConnectorName, targetConfig.ConnectorVersion, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("Error initializing connector plugin %q: %s", targetConfig.ConnectorName, err.Error()))
		return err
	}
	defer client.Close()

	jobId, err = job.StartJob(ctx, targetConfig)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("Start job with jobID: '%s'", jobId))
	job.UpdateJobEvent(targetConfig, logger, jobId, job.InProgress, nil)

	defer func() {
		if syncError == nil {
			job.UpdateJobEvent(targetConfig, logger, jobId, job.Completed, nil)
		} else {
			job.UpdateJobEvent(targetConfig, logger, jobId, job.Failed, syncError)
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

	eh := error_handler.NewBaseErrorHandler()

	if len(pluginInfo.Type) == 0 { // Backwards compatibility
		dsSyncTargetSync(ctx, logger, targetConfig, client, jobId, s, error_handler.Wrap(eh, "full data source sync: %w", error_handler.ErrorPlaceholder))
	}

	for _, pluginType := range pluginInfo.Type {
		switch pluginType {
		case plugin2.PluginType_PLUGIN_TYPE_FULL_DS_SYNC:
			dsSyncTargetSync(ctx, logger, targetConfig, client, jobId, s, error_handler.Wrap(eh, "full data source sync: %w", error_handler.ErrorPlaceholder))
		case plugin2.PluginType_PLUGIN_TYPE_IS_SYNC:
			isSyncTargetSync(ctx, logger, targetConfig, client, jobId, s, error_handler.Wrap(eh, "identity store sync: %w", error_handler.ErrorPlaceholder))
		case plugin2.PluginType_PLUGIN_TYPE_TAG_SYNC:
			tagSyncTargetSync(ctx, logger, targetConfig, client, jobId, s, error_handler.Wrap(eh, "tag sync: %w", error_handler.ErrorPlaceholder))
		case plugin2.PluginType_PLUGIN_TYPE_RESOURCE_PROVIDER:
			resourceProviderSync(ctx, logger, targetConfig, client, jobId, s, error_handler.Wrap(eh, "resource provider sync: %w", error_handler.ErrorPlaceholder))
		default:
			return fmt.Errorf("unsupported plugin type: %s", pluginType)
		}
	}

	s.wg.Wait()

	return eh.GetError()
}

func (s *SyncJob) Finalize(ctx context.Context, baseConfig *types.BaseConfig, options *target.Options) error {
	return sendEndOfTarget(ctx, baseConfig, s.jobIds, options)
}

func (s *SyncJob) execute(ctx context.Context, logger hclog.Logger, targetID string, jobID string, syncType string, syncTypeLabel string, skipSync bool, syncTask job.Task, cfg *types.BaseTargetConfig, c plugin.PluginClient, eh error_handler.ErrorHandler) {
	logger = logger.With("jobID", jobID, "syncType", syncTypeLabel)

	warningCollector, loggingCleanUp := logging.CreateWarningCapturingLogger(logger.(hclog.InterceptLogger))
	defer loggingCleanUp()

	taskEventUpdater := job.NewTaskEventUpdater(cfg, logger, jobID, syncType, warningCollector)

	defer func() {
		if r := recover(); r != nil {
			logger.Error(fmt.Sprintf("Panic occurred during %s sync: %v", syncTypeLabel, r))

			err := fmt.Errorf("panic occurred during %s sync", syncTypeLabel)
			eh.Error(err)

			taskEventUpdater.SetStatusToFailed(ctx, err)
		}
	}()

	switch {
	case skipSync:
		taskEventUpdater.SetStatusToSkipped(ctx)
		logger.Info("Skipping sync of " + syncTypeLabel)
	case targetID == "":
		taskEventUpdater.SetStatusToSkipped(ctx)

		idField := "data-source-id"
		if syncType == constants.IdentitySync {
			idField = "identity-store-id"
		}

		logger.Warn("No " + idField + " argument found. Skipping syncing of " + syncTypeLabel)
	default:
		syncErrHandler := error_handler.Wrap(eh, "failed to execute %s sync: %w", syncTypeLabel, error_handler.ErrorPlaceholder)
		s.sync(ctx, logger, cfg, syncTypeLabel, taskEventUpdater, syncTask, c, syncType, jobID, syncErrHandler)

		if syncErrHandler.HasError() {
			// Sync error is already pushed to task error
			return
		}
	}
}

func (s *SyncJob) logForwardingEnabled(syncType string) bool {
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

func (s *SyncJob) sync(ctx context.Context, logger hclog.Logger, cfg *types.BaseTargetConfig, syncTypeLabel string, taskEventUpdater job.TaskEventUpdater, syncTask job.Task, c plugin.PluginClient, syncType string, jobID string, eh error_handler.ErrorHandler) {
	eh = error_handler.OnError(eh, func(err error) {
		taskEventUpdater.SetStatusToFailed(ctx, err)

		target.HandleTargetError(err, logger, fmt.Sprintf("Synchronizing %s failed", syncType))
	})

	var taskMutex sync.Mutex

	if s.logForwardingEnabled(syncType) {
		targetCfg, cleanup, taskLoggingError := logging.CreateTaskLogger(cfg, jobID, syncType)
		if taskLoggingError != nil {
			eh.Error(taskLoggingError)

			return
		}

		cfg = targetCfg

		defer func() {
			s.wg.Add(1)

			go func() {
				defer s.wg.Done()

				taskMutex.Lock()

				cleanUpErr := cleanup()
				if cleanUpErr != nil {
					logger.Warn(fmt.Sprintf("Failed to close logger for task: %s", cleanUpErr.Error()))
				}
			}()
		}()
	}

	_, err := syncTask.IsClientValid(ctx, c)
	incompatibleVersionError := version_management.IncompatiblePluginVersionError{}

	var internalPluginStatusError *grpc_error.InternalPluginStatusError

	if errors.As(err, &internalPluginStatusError) {
		if internalPluginStatusError.StatusCode() == codes.Unimplemented {
			logger.Info(fmt.Sprintf("Plugin does not implement a syncer for %s. Skipping", syncTypeLabel))
			taskEventUpdater.SetStatusToSkipped(ctx) // Skip should be sent before we send a start status event

			return
		}
	}

	s.syncStartMutex.Lock()
	logger.Info(fmt.Sprintf("Synchronizing %s...", syncTypeLabel))
	taskEventUpdater.SetStatusToStarted(ctx)
	s.syncStartMutex.Unlock()

	if errors.As(err, &incompatibleVersionError) {
		eh.Errorf("unable to execute %s sync: %w", syncTypeLabel, incompatibleVersionError)

		return
	} else if err != nil {
		eh.Error(err)

		return
	}

	syncParts := syncTask.GetParts()

	taskWg := sync.WaitGroup{}

	for i, taskPart := range syncParts {
		s.runTaskPartSync(ctx, logger.With("subtask", i), cfg, syncTypeLabel, taskEventUpdater, jobID, syncType, taskPart, i, syncParts, c, &taskWg, eh)

		if eh.HasError() {
			return
		}
	}

	taskMutex.Lock()
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()
		defer taskMutex.Unlock()

		taskWg.Wait()
		taskEventUpdater.SetStatusToCompleted(ctx, syncTask.GetTaskResults())
	}()
}

func (s *SyncJob) runTaskPartSync(ctx context.Context, logger hclog.Logger, cfg *types.BaseTargetConfig, syncTypeLabel string, taskEventUpdater job.TaskEventUpdater, jobID string, syncType string, taskPart job.TaskPart, i int,
	syncParts []job.TaskPart, c plugin.PluginClient, wg *sync.WaitGroup, eh error_handler.ErrorHandler) {
	s.syncStartMutex.Lock()
	hasImportPart := false

	wg.Add(1)

	go func() {
		defer wg.Done()

		logger.Debug(fmt.Sprintf("Start sync task part %d out of %d", i+1, len(syncParts)))

		status, subtaskId, err := taskPart.StartSyncAndQueueTaskPart(ctx, logger, c, taskEventUpdater, func(f func() error) error {
			hasImportPart = true

			s.syncQueueMutex.Lock()
			defer s.syncQueueMutex.Unlock()

			s.syncStartMutex.Unlock()

			return f()
		})

		if !hasImportPart {
			s.syncStartMutex.Unlock()
		}

		if err != nil {
			eh.Errorf("synchronizing %s : %w", syncTypeLabel, err)

			return
		}

		if status == job.Queued {
			taskEventUpdater.SetStatusToQueued(ctx)
			logger.Info(fmt.Sprintf("Waiting for server to start processing %s...", syncTypeLabel))
		}

		syncResult := taskPart.GetResultObject()

		if syncResult != nil {
			subtask, err := job.WaitForJobToComplete(ctx, logger, jobID, syncType, subtaskId, syncResult, cfg, status)
			if err != nil {
				eh.Error(err)

				return
			}

			if subtask.Status == job.Failed {
				var subtaskErr error
				subtaskErr = multierror.Append(subtaskErr, array.Map(subtask.Errors, func(err *string) error { return errors.New(*err) })...)

				eh.Error(subtaskErr)

				return
			} else if subtask.Status == job.TimeOut {
				eh.Errorf("synchronizing %s timed out", syncTypeLabel)

				return
			}

			err = taskPart.ProcessResults(logger, syncResult)
			if err != nil {
				eh.Error(err)

				return
			}
		} else if status != job.Completed {
			eh.Errorf("unable to load results")

			return
		}
	}()
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
