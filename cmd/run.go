package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	sync2 "sync"
	"syscall"
	"time"

	"github.com/raito-io/cli/base/util/error/grpc_error"
	"google.golang.org/grpc/codes"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/raito-io/cli/internal/access_provider"
	"github.com/raito-io/cli/internal/clitrigger"
	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/data_source"
	"github.com/raito-io/cli/internal/data_usage"
	"github.com/raito-io/cli/internal/identity_store"
	"github.com/raito-io/cli/internal/job"
	"github.com/raito-io/cli/internal/logging"
	"github.com/raito-io/cli/internal/plugin"
	"github.com/raito-io/cli/internal/target"
	"github.com/raito-io/cli/internal/util/array"
	"github.com/raito-io/cli/internal/version"
	"github.com/raito-io/cli/internal/version_management"
)

func initRunCommand(rootCmd *cobra.Command) {
	var cmd = &cobra.Command{
		Hidden: true,
		Use:    "run",
		Short:  "Run all the configured synchronizations",
		Long:   `Run all the configured synchronizations`,
		Run:    executeRun,
	}
	cmd.PersistentFlags().IntP(constants.FrequencyFlag, "f", 0, "The frequency used to do the sync (in minutes). When not set, the default value '0' is used, which means the sync will run once and quit after.")
	cmd.PersistentFlags().Bool(constants.SkipDataSourceSyncFlag, false, "If set, the data source meta data synchronization step to Raito will be skipped for each of the targets.")
	cmd.PersistentFlags().Bool(constants.SkipIdentityStoreSyncFlag, false, "If set, the identity store synchronization step to Raito will be skipped for each of the targets.")
	cmd.PersistentFlags().Bool(constants.SkipDataAccessSyncFlag, false, "If set, the data access information from Raito will not be synced to the data sources in the target list.")
	cmd.PersistentFlags().Bool(constants.SkipDataUsageSyncFlag, false, "If set, the data usage information synchronization step to Raito will be skipped for each of the targets.")

	cmd.PersistentFlags().Bool(constants.LockAllWhoFlag, false, "If set, the 'who' of all access providers imported into Raito Cloud will be locked.")
	cmd.PersistentFlags().Bool(constants.LockAllWhatFlag, false, "If set, the 'what' of all access providers imported into Raito Cloud will be locked.")
	cmd.PersistentFlags().Bool(constants.LockAllNamesFlag, false, "If set, the names of all access providers imported into Raito Cloud will be locked.")
	cmd.PersistentFlags().Bool(constants.LockAllDeleteFlag, false, "If set, the deletion of all access providers imported into Raito Cloud will be locked.")

	cmd.PersistentFlags().Bool(constants.DisableWebsocketFlag, false, "If set, raito will not setup a websocket to trigger new syncs. This flag has only effect if frequency is set.")
	cmd.PersistentFlags().Bool(constants.DisableLogForwarding, false, "If set, sync logs will not be forwarded to Raito Cloud.")
	cmd.PersistentFlags().Bool(constants.DisableLogForwardingDataSourceSync, false, "If set, data source sync logs will not be forwarded to Raito Cloud.")
	cmd.PersistentFlags().Bool(constants.DisableLogForwardingDataAccessSync, false, "If set, data access sync logs will not be forwarded to Raito Cloud.")
	cmd.PersistentFlags().Bool(constants.DisableLogForwardingIdentityStoreSync, false, "If set, identity store sync logs will not be forwarded to Raito Cloud.")
	cmd.PersistentFlags().Bool(constants.DisableLogForwardingDataUsageSync, false, "If set, data usage sync logs will not be forwarded to Raito Cloud.")

	BindFlag(constants.FrequencyFlag, cmd)
	BindFlag(constants.SkipDataSourceSyncFlag, cmd)
	BindFlag(constants.SkipIdentityStoreSyncFlag, cmd)
	BindFlag(constants.SkipDataAccessSyncFlag, cmd)
	BindFlag(constants.SkipDataUsageSyncFlag, cmd)
	BindFlag(constants.LockAllWhoFlag, cmd)
	BindFlag(constants.LockAllWhatFlag, cmd)
	BindFlag(constants.LockAllNamesFlag, cmd)
	BindFlag(constants.LockAllDeleteFlag, cmd)
	BindFlag(constants.DisableWebsocketFlag, cmd)
	BindFlag(constants.DisableLogForwarding, cmd)
	BindFlag(constants.DisableLogForwardingDataSourceSync, cmd)
	BindFlag(constants.DisableLogForwardingDataAccessSync, cmd)
	BindFlag(constants.DisableLogForwardingIdentityStoreSync, cmd)
	BindFlag(constants.DisableLogForwardingDataUsageSync, cmd)

	cmd.FParseErrWhitelist.UnknownFlags = true

	rootCmd.AddCommand(cmd)
}

func executeRun(cmd *cobra.Command, args []string) {
	otherArgs := cmd.Flags().Args()

	baseConfig, err := target.BuildBaseConfigFromFlags(hclog.L(), otherArgs)
	if err != nil {
		hclog.L().Error(err.Error())
		os.Exit(1)
	}

	freq := viper.GetInt(constants.FrequencyFlag)
	if freq <= 0 {
		hclog.L().Info("Running synchronization just once.")

		baseConfig.BaseLogger = baseConfig.BaseLogger.With("iteration", 0)
		err = executeSingleRun(baseConfig)

		if err != nil {
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	} else {
		hclog.L().Info(fmt.Sprintf("Starting synchronization every %d minutes.", freq))
		hclog.L().Info("Press 'ctrl+c' to stop the program.")

		ticker := time.NewTicker(time.Duration(freq) * time.Minute)

		ctx, cancelFn := context.WithCancel(context.Background())

		waitGroup := sync2.WaitGroup{}

		cliTriggerChannel := make(chan clitrigger.TriggerEvent)
		defer close(cliTriggerChannel)

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGUSR1, os.Interrupt)

		returnSignal := 0

		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()
			defer ticker.Stop()

			cliTriggerCtx, cliTriggerCancel := context.WithCancel(ctx)
			cliTrigger := startListingToCliTriggers(cliTriggerCtx, baseConfig, cliTriggerChannel)

			defer func() {
				cliTriggerCancel()
				cliTrigger.Wait()
			}()

			baseConfig.BaseLogger = baseConfig.BaseLogger.With("iteration", 1)
			if runErr := executeSingleRun(baseConfig); runErr != nil {
				baseConfig.BaseLogger.Error(fmt.Sprintf("Run failed: %s", runErr.Error()))
			}

			it := 2
			for {
				select {
				case <-ticker.C:
					cliTrigger.Reset()

					baseConfig.BaseLogger = baseConfig.BaseLogger.With("iteration", it)
					if runErr := executeSingleRun(baseConfig); runErr != nil {
						baseConfig.BaseLogger.Error(fmt.Sprintf("Run failed: %s", runErr.Error()))
					}

					it++
				case cliTrigger := <-cliTriggerChannel:
					baseConfig.BaseLogger = baseConfig.BaseLogger.With("iteration", it)

					err := handleCliTrigger(baseConfig, &cliTrigger)
					if err != nil {
						baseConfig.BaseLogger.Warn("Cli Trigger failed: %s", err.Error())
					}

					it++
				case <-ctx.Done():
					baseConfig.BaseLogger.Debug("Context done: closing syncing routine.")
					return
				}

				hclog.L().Info("Press 'ctrl+c' to stop the program.")
			}
		}()

		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			defer cancelFn()
			defer hclog.L().Info("Waiting for the current synchronization run to end ...")

			for {
				select {
				case <-ctx.Done():
					hclog.L().Debug("Context done: Will stop all running routines...")
					return
				case s := <-sigs:
					hclog.L().Debug("Received signal: %s. Will stop all running routines...", s.String())

					if sysc, ok := s.(syscall.Signal); ok {
						returnSignal = int(sysc)
					}

					return
				}
			}
		}()

		waitGroup.Wait()
		hclog.L().Info("All routines finished. Bye!")

		if returnSignal != 0 {
			hclog.L().Debug("Exit with code: %d", returnSignal)
			syscall.Exit(returnSignal)
		}
	}
}

func executeSingleRun(baseconfig *target.BaseConfig) error {
	start := time.Now()

	err := runSync(baseconfig)

	sec := time.Since(start).Round(time.Millisecond)
	baseconfig.BaseLogger.Info(fmt.Sprintf("Finished execution of all targets in %s", sec))

	return err
}

func runSync(baseconfig *target.BaseConfig) error {
	compatibilityInformation, err := version_management.IsCompatibleWithRaitoCloud(baseconfig)
	if err != nil {
		baseconfig.BaseLogger.Error(fmt.Sprintf("Failed to check compatibility with Raito Cloud: %s", err.Error()))

		return fmt.Errorf("compatibility check failed: %s", err.Error())
	}

	switch compatibilityInformation.Compatibility {
	case version_management.NotSupported:
		baseconfig.BaseLogger.Error(fmt.Sprintf("CLI version is not compatible with Raito Cloud. Please upgrade to a supported version (%s).", compatibilityInformation.SupportedVersions))

		return errors.New("unsupported CLI version")
	case version_management.Deprecated:
		warning := " "
		if compatibilityInformation.DeprecatedWarningMsg != nil {
			warning += *compatibilityInformation.DeprecatedWarningMsg
		}

		baseconfig.BaseLogger.Warn(fmt.Sprintf("CLI version %s is deprecated.%s Please upgrade to supported version (%s) soon.", version.GetCliVersion().String(), warning, compatibilityInformation.SupportedVersions))

		fallthrough
	case version_management.Supported:
		return target.RunTargets(baseconfig, runTargetSync)
	case version_management.CompatibilityUnknown:
	}

	return errors.New("unknown CLI version")
}

func execute(targetID string, jobID string, syncType string, syncTypeLabel string, skipSync bool,
	syncTask job.Task, cfg *target.BaseTargetConfig, c plugin.PluginClient) error {
	ctx := context.Background()

	cfg, warningCollector, loggingCleanUp, err := logging.CreateWarningCapturingLogger(cfg)
	if err != nil {
		return err
	}

	defer loggingCleanUp()

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
	}

	return !viper.GetBool(cmdFlag)
}

func sync(ctx context.Context, cfg *target.BaseTargetConfig, syncTypeLabel string, taskEventUpdater job.TaskEventUpdater, syncTask job.Task, c plugin.PluginClient, syncType string, jobID string) (err error) {
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

	cfg.TargetLogger.Info(fmt.Sprintf("Synchronizing %s...", syncTypeLabel))

	taskEventUpdater.SetStatusToStarted(ctx)

	_, err = syncTask.IsClientValid(ctx, c)
	incompatibleVersionError := version_management.IncompatiblePluginVersionError{}

	if pluginError, f := err.(*grpc_error.InternalPluginStatusError); f && pluginError.StatusCode() == codes.Unimplemented {
		cfg.TargetLogger.Info(fmt.Sprintf("Plugin does not implement a syncer for %s. Skipping", syncTypeLabel))
		taskEventUpdater.SetStatusToSkipped(ctx)

		return nil
	} else if errors.As(err, &incompatibleVersionError) {
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

func runTaskPartSync(ctx context.Context, cfg *target.BaseTargetConfig, syncTypeLabel string, taskEventUpdater job.TaskEventUpdater, jobID string, syncType string, taskPart job.TaskPart, i int, syncParts []job.TaskPart, c plugin.PluginClient) error {
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

func runTargetSync(targetConfig *target.BaseTargetConfig) (syncError error) {
	targetConfig.TargetLogger.Info("Executing target...")

	start := time.Now()

	defer func() {
		if syncError != nil {
			targetConfig.TargetLogger.Error(fmt.Sprintf("Failed execution: %s", syncError.Error()), "success")
		} else {
			targetConfig.TargetLogger.Info(fmt.Sprintf("Successfully finished execution in %s", time.Since(start).Round(time.Millisecond)), "success")
		}
	}()

	client, err := plugin.NewPluginClient(targetConfig.ConnectorName, targetConfig.ConnectorVersion, targetConfig.TargetLogger)
	if err != nil {
		targetConfig.TargetLogger.Error(fmt.Sprintf("Error initializing connector plugin %q: %s", targetConfig.ConnectorName, err.Error()))
		return err
	}
	defer client.Close()

	jobID, err := job.StartJob(context.Background(), targetConfig)
	if err != nil {
		return err
	}

	targetConfig.TargetLogger.Info(fmt.Sprintf("Start job with jobID: '%s'", jobID))
	job.UpdateJobEvent(targetConfig, jobID, job.InProgress, nil)

	defer func() {
		if syncError == nil {
			job.UpdateJobEvent(targetConfig, jobID, job.Completed, nil)
		} else {
			job.UpdateJobEvent(targetConfig, jobID, job.Failed, syncError)
		}
	}()

	err = dataSourceSync(targetConfig, jobID, client)
	if err != nil {
		return err
	}

	err = identityStoreSync(targetConfig, jobID, client)
	if err != nil {
		return err
	}

	err = dataAccessSync(targetConfig, jobID, client)
	if err != nil {
		return err
	}

	err = dataUsageSync(targetConfig, jobID, client)
	if err != nil {
		return err
	}

	return nil
}

func dataUsageSync(targetConfig *target.BaseTargetConfig, jobID string, client plugin.PluginClient) error {
	dataUsageSyncTask := &data_usage.DataUsageSync{TargetConfig: targetConfig, JobId: jobID}

	err := execute(targetConfig.DataSourceId, jobID, constants.DataUsageSync, "data usage", targetConfig.SkipDataUsageSync, dataUsageSyncTask, targetConfig, client)
	if err != nil {
		return err
	}

	return nil
}

func dataAccessSync(targetConfig *target.BaseTargetConfig, jobID string, client plugin.PluginClient) error {
	dataAccessSyncTask := &access_provider.DataAccessSync{TargetConfig: targetConfig, JobId: jobID}

	err := execute(targetConfig.DataSourceId, jobID, constants.DataAccessSync, "data access", targetConfig.SkipDataAccessSync, dataAccessSyncTask, targetConfig, client)
	if err != nil {
		return err
	}

	return nil
}

func identityStoreSync(targetConfig *target.BaseTargetConfig, jobID string, client plugin.PluginClient) error {
	identityStoreSyncTask := &identity_store.IdentityStoreSync{TargetConfig: targetConfig, JobId: jobID}

	err := execute(targetConfig.IdentityStoreId, jobID, constants.IdentitySync, "identity store", targetConfig.SkipIdentityStoreSync, identityStoreSyncTask, targetConfig, client)
	if err != nil {
		return err
	}

	return nil
}

func dataSourceSync(targetConfig *target.BaseTargetConfig, jobID string, client plugin.PluginClient) error {
	dataSourceSyncTask := &data_source.DataSourceSync{TargetConfig: targetConfig, JobId: jobID}

	err := execute(targetConfig.DataSourceId, jobID, constants.DataSourceSync, "data source metadata", targetConfig.SkipDataSourceSync, dataSourceSyncTask, targetConfig, client)
	if err != nil {
		return err
	}

	return nil
}

func handleCliTrigger(baseConfig *target.BaseConfig, triggerEvent *clitrigger.TriggerEvent) error {
	if triggerEvent.ApUpdate != nil {
		return handleApUpdateTrigger(baseConfig, triggerEvent.ApUpdate)
	}

	return nil
}

func handleApUpdateTrigger(config *target.BaseConfig, apUpdate *clitrigger.ApUpdate) error {
	return target.RunTargets(config, runTargetSync, target.WithDataSourceIds(apUpdate.DataSourceNames...), target.WithConfigOption(func(targetConfig *target.BaseTargetConfig) {
		targetConfig.SkipIdentityStoreSync = true
		targetConfig.SkipDataSourceSync = true
		targetConfig.SkipDataUsageSync = true

		targetConfig.SkipDataAccessImport = true
		targetConfig.OnlyOutOfSyncData = true
	}))
}

func startListingToCliTriggers(ctx context.Context, baseConfig *target.BaseConfig, outputChannel chan clitrigger.TriggerEvent) clitrigger.CliTrigger {
	cliTrigger, err := clitrigger.CreateCliTrigger(baseConfig)
	if err != nil {
		baseConfig.BaseLogger.Warn(fmt.Sprintf("Unable to start asynchronous access provider sync: %s", err.Error()))
		return nil
	}

	ch := cliTrigger.TriggerChannel(ctx)

	go func() {
		for i := range ch {
			outputChannel <- i
		}
	}()

	return cliTrigger
}
