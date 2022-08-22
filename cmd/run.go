package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/raito-io/cli/internal/access_provider"
	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/data_source"
	"github.com/raito-io/cli/internal/data_usage"
	"github.com/raito-io/cli/internal/identity_store"
	"github.com/raito-io/cli/internal/job"
	"github.com/raito-io/cli/internal/plugin"
	"github.com/raito-io/cli/internal/target"
)

type task interface {
	StartSyncAndQueueJob(c plugin.PluginClient) (job.JobStatus, error)
	ProcessResults(results interface{}) error
	GetResultObject() interface{}
}

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

	BindFlag(constants.FrequencyFlag, cmd)
	BindFlag(constants.SkipDataSourceSyncFlag, cmd)
	BindFlag(constants.SkipIdentityStoreSyncFlag, cmd)
	BindFlag(constants.SkipDataAccessSyncFlag, cmd)
	BindFlag(constants.SkipDataUsageSyncFlag, cmd)

	cmd.FParseErrWhitelist.UnknownFlags = true

	rootCmd.AddCommand(cmd)
}

func executeRun(cmd *cobra.Command, args []string) {
	otherArgs := cmd.Flags().Args()

	freq := viper.GetInt(constants.FrequencyFlag)
	if freq <= 0 {
		logger.Info("Running synchronization just once.")

		err := executeSingleRun(logger.With("iteration", 0), otherArgs)
		if err != nil {
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	} else {
		logger.Info(fmt.Sprintf("Starting synchronization every %d minutes.", freq))
		logger.Info("Press the letter 'q' (and press return) to stop the program.")

		ticker := time.NewTicker(time.Duration(freq) * time.Minute)
		quit := make(chan struct{})
		finished := make(chan struct{})
		go func() {
			executeSingleRun(logger.With("iteration", 1), otherArgs) //nolint
			it := 2
			for {
				select {
				case <-ticker.C:
					executeSingleRun(logger.With("iteration", it), otherArgs) //nolint
					it++
				case <-quit:
					ticker.Stop()
					finished <- struct{}{}
					return
				}
			}
		}()

		for {
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			if strings.TrimSpace(strings.ToLower(text)) == "q" {
				logger.Info("Waiting for the current synchronization run to end ...")
				quit <- struct{}{}
				break
			} else {
				logger.Info("Press the letter 'q' (and press return) to stop the program.")
			}
		}

		<-finished
		logger.Info("All routines finished. Bye!")
	}
}

func executeSingleRun(baseLogger hclog.Logger, otherArgs []string) error {
	start := time.Now()

	err := runSync(baseLogger, otherArgs)

	sec := time.Since(start).Round(time.Millisecond)
	baseLogger.Info(fmt.Sprintf("Finished execution of all targets in %s", sec))

	return err
}

func runSync(baseLogger hclog.Logger, otherArgs []string) error {
	return target.RunTargets(baseLogger, otherArgs, runTargetSync)
}

func execute(targetID string, jobID string, syncType string, syncTypeLabel string, skipSync bool,
	syncTask task, cfg *target.BaseTargetConfig, c plugin.PluginClient) error {
	switch {
	case skipSync:
		job.AddJobEvent(cfg, jobID, syncType, constants.Skipped)
		cfg.Logger.Info("Skipping sync of " + syncTypeLabel)
	case targetID == "":
		job.AddJobEvent(cfg, jobID, syncType, constants.Skipped)

		idField := "data-source-id"
		if syncType == constants.IdentitySync {
			idField = "identity-store-id"
		}

		cfg.Logger.Info("No " + idField + " argument found. Skipping syncing of " + syncTypeLabel)
	default:
		cfg.Logger.Info(fmt.Sprintf("Synchronizing %s...", syncTypeLabel))
		job.AddJobEvent(cfg, jobID, syncType, constants.Started)

		status, err := syncTask.StartSyncAndQueueJob(c)
		if err != nil {
			target.HandleTargetError(err, cfg, "synchronizing "+syncType)
			job.AddJobEvent(cfg, jobID, syncType, constants.Failed)

			return err
		}

		if status == job.Queued {
			cfg.Logger.Info(fmt.Sprintf("Wait for server to start processing %s...", syncTypeLabel))
		}

		syncResult := syncTask.GetResultObject()
		_, err = waitForJobToComplete(jobID, syncType, syncResult, cfg)

		if err != nil {
			return err
		}

		err = syncTask.ProcessResults(syncResult)
		if err != nil {
			return err
		}
	}

	return nil
}

func waitForJobToComplete(jobID string, syncType string, syncResult interface{}, cfg *target.BaseTargetConfig) (*job.JobStatus, error) {
	status := job.Queued

	i := 0

	for status.IsRunning() {
		time.Sleep(1 * time.Second)
		updatedStatus, err := job.GetTaskStatus(cfg, jobID, syncType, syncResult)

		if err != nil {
			return nil, err
		} else if updatedStatus == nil {
			return nil, fmt.Errorf("received in valid job status")
		}
		status = *updatedStatus
		cfg.Logger.Debug(fmt.Sprintf("Current status on iteration %d: %s", i, status.String()))
		i += 1
	}

	return &status, nil
}

func updateJobStatus(status job.JobStatus, jobId, jobType string, cfg *target.BaseTargetConfig) {
	job.AddJobEvent(cfg, jobId, jobType, status.String())
}

func runTargetSync(targetConfig *target.BaseTargetConfig) error {
	targetConfig.Logger.Info("Executing target...")

	start := time.Now()

	client, err := plugin.NewPluginClient(targetConfig.ConnectorName, targetConfig.ConnectorVersion, targetConfig.Logger)
	if err != nil {
		targetConfig.Logger.Error(fmt.Sprintf("Error initializing connector plugin %q: %s", targetConfig.ConnectorName, err.Error()))
		return err
	}
	defer client.Close()

	jobID, _ := job.StartJob(targetConfig)

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

	targetConfig.Logger.Info(fmt.Sprintf("Successfully finished execution in %s", time.Since(start).Round(time.Millisecond)), "success")

	job.AddJobEvent(targetConfig, jobID, constants.Job, constants.Completed)

	return nil
}

func dataUsageSync(targetConfig *target.BaseTargetConfig, jobID string, client plugin.PluginClient) error {
	dataUsageSyncTask := &data_usage.DataUsageSync{TargetConfig: targetConfig, JobId: jobID, StatusUpdater: func(status job.JobStatus) { updateJobStatus(status, jobID, constants.DataUsageSync, targetConfig) }}

	err := execute(targetConfig.DataSourceId, jobID, constants.DataUsageSync, "data usage", targetConfig.SkipDataUsageSync, dataUsageSyncTask, targetConfig, client)
	if err != nil {
		job.AddJobEvent(targetConfig, jobID, constants.Job, constants.Failed)
		return err
	}

	return nil
}

func dataAccessSync(targetConfig *target.BaseTargetConfig, jobID string, client plugin.PluginClient) error {
	dataAccessSyncTask := &access_provider.DataAccessSync{TargetConfig: targetConfig, JobId: jobID, StatusUpdater: func(status job.JobStatus) { updateJobStatus(status, jobID, constants.DataAccessSync, targetConfig) }}

	err := execute(targetConfig.DataSourceId, jobID, constants.DataAccessSync, "data access", targetConfig.SkipDataAccessSync, dataAccessSyncTask, targetConfig, client)
	if err != nil {
		job.AddJobEvent(targetConfig, jobID, constants.Job, constants.Failed)
		return err
	}

	return nil
}

func identityStoreSync(targetConfig *target.BaseTargetConfig, jobID string, client plugin.PluginClient) error {
	identityStoreSyncTask := &identity_store.IdentityStoreSync{TargetConfig: targetConfig, JobId: jobID, StatusUpdater: func(status job.JobStatus) { updateJobStatus(status, jobID, constants.IdentitySync, targetConfig) }}

	err := execute(targetConfig.IdentityStoreId, jobID, constants.IdentitySync, "identity store", targetConfig.SkipIdentityStoreSync, identityStoreSyncTask, targetConfig, client)
	if err != nil {
		job.AddJobEvent(targetConfig, jobID, constants.Job, constants.Failed)
		return err
	}

	return nil
}

func dataSourceSync(targetConfig *target.BaseTargetConfig, jobID string, client plugin.PluginClient) error {
	dataSourceSyncTask := &data_source.DataSourceSync{TargetConfig: targetConfig, JobId: jobID, StatusUpdater: func(status job.JobStatus) { updateJobStatus(status, jobID, constants.DataSourceSync, targetConfig) }}

	err := execute(targetConfig.DataSourceId, jobID, constants.DataSourceSync, "data source metadata", targetConfig.SkipDataSourceSync, dataSourceSyncTask, targetConfig, client)
	if err != nil {
		job.AddJobEvent(targetConfig, jobID, constants.Job, constants.Failed)
		return err
	}

	return nil
}
