package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/raito-io/cli/internal/access_provider"
	"github.com/raito-io/cli/internal/clitrigger"
	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/data_source"
	"github.com/raito-io/cli/internal/data_usage"
	"github.com/raito-io/cli/internal/identity_store"
	"github.com/raito-io/cli/internal/job"
	"github.com/raito-io/cli/internal/plugin"
	"github.com/raito-io/cli/internal/target"
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
	cmd.PersistentFlags().Bool(constants.DisableWebsocketFlag, false, "If set, raito will not setup a websocket to trigger new syncs. This flag has only effect if frequency is set.")

	BindFlag(constants.FrequencyFlag, cmd)
	BindFlag(constants.SkipDataSourceSyncFlag, cmd)
	BindFlag(constants.SkipIdentityStoreSyncFlag, cmd)
	BindFlag(constants.SkipDataAccessSyncFlag, cmd)
	BindFlag(constants.SkipDataUsageSyncFlag, cmd)
	BindFlag(constants.DisableWebsocketFlag, cmd)

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
		hclog.L().Info("Press the letter 'q' (and press return) to stop the program.")

		ticker := time.NewTicker(time.Duration(freq) * time.Minute)

		ctx, cancelFn := context.WithCancel(context.Background())

		finished := make(chan struct{})
		defer close(finished)

		cliTriggerChannel := make(chan clitrigger.TriggerEvent)
		defer close(cliTriggerChannel)

		go func() {
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
					baseConfig.BaseLogger = baseConfig.BaseLogger.With("iteration", 1)
					if runErr := executeSingleRun(baseConfig); runErr != nil {
						baseConfig.BaseLogger.Error(fmt.Sprintf("Run failed: %s", runErr.Error()))
					}

					it++
				case cliTrigger := <-cliTriggerChannel:
					err := handleCliTrigger(baseConfig, &cliTrigger)
					if err != nil {
						baseConfig.BaseLogger.Warn("Cli Trigger failed: %s", err.Error())
					}
				case <-ctx.Done():
					finished <- struct{}{}
					return
				}
			}
		}()

		for {
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			if strings.TrimSpace(strings.ToLower(text)) == "q" {
				hclog.L().Info("Waiting for the current synchronization run to end ...")
				cancelFn()
				break
			} else {
				hclog.L().Info("Press the letter 'q' (and press return) to stop the program.")
			}
		}

		<-finished
		hclog.L().Info("All routines finished. Bye!")
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
	compatibilityInformation, err := version_management.IsCompatibleWithAppServer(baseconfig)
	if err != nil {
		baseconfig.BaseLogger.Error(fmt.Sprintf("Failed to check compatibility with app server: %s", err.Error()))

		return fmt.Errorf("compatibility check failed: %s", err.Error())
	}

	switch compatibilityInformation.Compatibility {
	case version_management.NotSupported:
		baseconfig.BaseLogger.Error(fmt.Sprintf("CLI version is not compatible with app server. Please upgrade to a supported version (%s).", compatibilityInformation.SupportedVersions))

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
	}

	return errors.New("unknown CLI version")
}

func execute(targetID string, jobID string, syncType string, syncTypeLabel string, skipSync bool,
	syncTask job.Task, cfg *target.BaseTargetConfig, c plugin.PluginClient) error {
	taskEventUpdater := job.NewTaskEventUpdater(cfg, jobID, syncType)

	switch {
	case skipSync:
		taskEventUpdater.AddTaskEvent(job.Skipped)
		cfg.TargetLogger.Info("Skipping sync of " + syncTypeLabel)
	case targetID == "":
		taskEventUpdater.AddTaskEvent(job.Skipped)

		idField := "data-source-id"
		if syncType == constants.IdentitySync {
			idField = "identity-store-id"
		}

		cfg.TargetLogger.Info("No " + idField + " argument found. Skipping syncing of " + syncTypeLabel)
	default:
		err := sync(cfg, syncTypeLabel, taskEventUpdater, syncTask, c, syncType, jobID)
		if err != nil {
			return err
		}
	}

	return nil
}

func sync(cfg *target.BaseTargetConfig, syncTypeLabel string, taskEventUpdater job.TaskEventUpdater, syncTask job.Task, c plugin.PluginClient, syncType string, jobID string) (err error) {
	defer func() {
		if err != nil {
			cfg.TargetLogger.Error(fmt.Sprintf("Synchronizing %s failed: %s", syncTypeLabel, err.Error()))
		}
	}()

	cfg.TargetLogger.Info(fmt.Sprintf("Synchronizing %s...", syncTypeLabel))

	taskEventUpdater.AddTaskEvent(job.Started)

	_, err = syncTask.IsClientValid(context.Background(), c)
	incompatibleVersionError := version_management.IncompatiblePluginVersionError{}

	if errors.As(err, &incompatibleVersionError) {
		cfg.TargetLogger.Error(fmt.Sprintf("Unable to execute %s sync: %s", syncTypeLabel, incompatibleVersionError.Error()))
		taskEventUpdater.AddTaskEvent(job.Failed)

		return nil
	} else if err != nil {
		return err
	}

	syncParts := syncTask.GetParts()

	for i, taskPart := range syncParts {
		cfg.TargetLogger.Debug(fmt.Sprintf("Start sync task part %d out of %d", i+1, len(syncParts)))

		status, subtaskId, err := taskPart.StartSyncAndQueueTaskPart(c, taskEventUpdater)
		if err != nil {
			target.HandleTargetError(err, cfg, "synchronizing "+syncType)
			taskEventUpdater.AddTaskEvent(job.Failed)

			return err
		}

		if status != job.Completed {
			taskEventUpdater.AddTaskEvent(status)
		}

		if status == job.Queued {
			cfg.TargetLogger.Info(fmt.Sprintf("Waiting for server to start processing %s...", syncTypeLabel))
		}

		syncResult := taskPart.GetResultObject()

		if syncResult != nil {
			subtask, err := job.WaitForJobToComplete(jobID, syncType, subtaskId, syncResult, cfg, status)
			if err != nil {
				taskEventUpdater.AddTaskEvent(job.Failed)
				return err
			}

			if subtask.Status == job.Failed {
				taskEventUpdater.AddTaskEvent(job.Failed)
				return fmt.Errorf("%s", strings.Join(subtask.Errors, ", "))
			}

			err = taskPart.ProcessResults(syncResult)
			if err != nil {
				taskEventUpdater.AddTaskEvent(job.Failed)
				return err
			}
		} else if status != job.Completed {
			taskEventUpdater.AddTaskEvent(job.Failed)
			return fmt.Errorf("unable to load results")
		}
	}

	taskEventUpdater.AddTaskEvent(job.Completed)

	return nil
}

func runTargetSync(targetConfig *target.BaseTargetConfig) (syncError error) {
	targetConfig.TargetLogger.Info("Executing target...")

	start := time.Now()

	client, err := plugin.NewPluginClient(targetConfig.ConnectorName, targetConfig.ConnectorVersion, targetConfig.TargetLogger)
	if err != nil {
		targetConfig.TargetLogger.Error(fmt.Sprintf("Error initializing connector plugin %q: %s", targetConfig.ConnectorName, err.Error()))
		return err
	}
	defer client.Close()

	jobID, _ := job.StartJob(targetConfig)
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

	targetConfig.TargetLogger.Info(fmt.Sprintf("Successfully finished execution in %s", time.Since(start).Round(time.Millisecond)), "success")

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
