package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	dapc "github.com/raito-io/cli/common/api/data_access"
	dspc "github.com/raito-io/cli/common/api/data_source"
	dupc "github.com/raito-io/cli/common/api/data_usage"
	ispc "github.com/raito-io/cli/common/api/identity_store"
	baseconfig "github.com/raito-io/cli/common/util/config"
	"github.com/raito-io/cli/internal/access_provider"
	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/data_access"
	"github.com/raito-io/cli/internal/data_source"
	"github.com/raito-io/cli/internal/data_usage"
	"github.com/raito-io/cli/internal/file"
	"github.com/raito-io/cli/internal/identity_store"
	"github.com/raito-io/cli/internal/job"
	"github.com/raito-io/cli/internal/plugin"
	"github.com/raito-io/cli/internal/target"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var accessRightsLastUpdated int64 = 0

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
	syncFunc func(c plugin.PluginClient, cfg target.BaseTargetConfig) error,
	cfg *target.BaseTargetConfig, c plugin.PluginClient) error {
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

		err := syncFunc(c, *cfg)
		if err != nil {
			target.HandleTargetError(err, cfg, "synchronizing "+syncType)
			job.AddJobEvent(cfg, jobID, syncType, constants.Failed)

			return err
		}

		job.AddJobEvent(cfg, jobID, syncType, constants.Completed)
	}

	return nil
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

	err = execute(targetConfig.DataSourceId, jobID, constants.DataSourceSync, "data source metadata", targetConfig.SkipDataSourceSync, syncDataSource, targetConfig, client)
	if err != nil {
		job.AddJobEvent(targetConfig, jobID, constants.Job, constants.Failed)
		return err
	}

	err = execute(targetConfig.IdentityStoreId, jobID, constants.IdentitySync, "identity store", targetConfig.SkipIdentityStoreSync, syncIdentityStore, targetConfig, client)
	if err != nil {
		job.AddJobEvent(targetConfig, jobID, constants.Job, constants.Failed)
		return err
	}

	err = execute(targetConfig.DataSourceId, jobID, constants.DataAccessSync, "data access", targetConfig.SkipDataAccessSync, syncDataAccess, targetConfig, client)
	if err != nil {
		job.AddJobEvent(targetConfig, jobID, constants.Job, constants.Failed)
		return err
	}

	err = execute(targetConfig.DataSourceId, jobID, constants.DataUsageSync, "data usage", targetConfig.SkipDataUsageSync, syncDataUsage, targetConfig, client)
	if err != nil {
		job.AddJobEvent(targetConfig, jobID, constants.Job, constants.Failed)
		return err
	}

	targetConfig.Logger.Info(fmt.Sprintf("Successfully finished execution in %s", time.Since(start).Round(time.Millisecond)), "success")

	job.AddJobEvent(targetConfig, jobID, constants.Job, constants.Completed)

	return nil
}

func syncDataSource(client plugin.PluginClient, targetConfig target.BaseTargetConfig) error {
	cn := strings.Replace(targetConfig.ConnectorName, "/", "-", -1)

	targetFile, err := filepath.Abs(file.CreateUniqueFileName(cn+"-ds", "json"))
	if err != nil {
		return err
	}

	targetConfig.Logger.Debug(fmt.Sprintf("Using %q as data source target file", targetFile))

	if targetConfig.DeleteTempFiles {
		defer os.Remove(targetFile)
	}

	syncerConfig := dspc.DataSourceSyncConfig{
		ConfigMap:  baseconfig.ConfigMap{Parameters: targetConfig.Parameters},
		TargetFile: targetFile,
	}

	dss, err := client.GetDataSourceSyncer()
	if err != nil {
		return err
	}

	md := dss.GetMetaData()
	err = data_source.SetMetaData(targetConfig, md)

	if err != nil {
		return err
	}

	res := dss.SyncDataSource(&syncerConfig)
	if res.Error != nil {
		return err
	}

	importerConfig := data_source.DataSourceImportConfig{
		BaseTargetConfig: targetConfig,
		TargetFile:       targetFile,
		DeleteUntouched:  targetConfig.DeleteUntouched,
		ReplaceTags:      targetConfig.ReplaceTags,
	}
	dsImporter := data_source.NewDataSourceImporter(&importerConfig)

	dsResult, err := dsImporter.TriggerImport()
	if err != nil {
		return err
	}

	targetConfig.Logger.Info(fmt.Sprintf("Successfully synced data source. Added: %d - Removed: %d - Updated: %d", dsResult.DataObjectsAdded, dsResult.DataObjectsRemoved, dsResult.DataObjectsUpdated))

	return nil
}

func syncIdentityStore(client plugin.PluginClient, targetConfig target.BaseTargetConfig) error {
	cn := strings.Replace(targetConfig.ConnectorName, "/", "-", -1)

	userFile, err := filepath.Abs(file.CreateUniqueFileName(cn+"-is-user", "json"))
	if err != nil {
		return err
	}

	groupFile, err := filepath.Abs(file.CreateUniqueFileName(cn+"-is-group", "json"))
	if err != nil {
		return err
	}

	targetConfig.Logger.Debug(fmt.Sprintf("Using %q as user target file", userFile))
	targetConfig.Logger.Debug(fmt.Sprintf("Using %q as groups target file", groupFile))

	if targetConfig.DeleteTempFiles {
		defer os.Remove(userFile)
		defer os.Remove(groupFile)
	}

	syncerConfig := ispc.IdentityStoreSyncConfig{
		ConfigMap: baseconfig.ConfigMap{Parameters: targetConfig.Parameters},
		UserFile:  userFile,
		GroupFile: groupFile,
	}

	iss, err := client.GetIdentityStoreSyncer()
	if err != nil {
		return err
	}

	result := iss.SyncIdentityStore(&syncerConfig)
	if result.Error != nil {
		return *(result.Error)
	}

	importerConfig := identity_store.IdentityStoreImportConfig{
		BaseTargetConfig: targetConfig,
		UserFile:         userFile,
		GroupFile:        groupFile,
		DeleteUntouched:  targetConfig.DeleteUntouched,
		ReplaceGroups:    targetConfig.ReplaceGroups,
		ReplaceTags:      targetConfig.ReplaceTags,
	}
	isImporter := identity_store.NewIdentityStoreImporter(&importerConfig)

	isResult, err := isImporter.TriggerImport()
	if err != nil {
		return err
	}

	targetConfig.Logger.Info(fmt.Sprintf("Successfully synced users and groups. Users: Added: %d - Removed: %d - Updated: %d | Groups: Added: %d - Removed: %d - Updated: %d", isResult.UsersAdded, isResult.UsersRemoved, isResult.UsersUpdated, isResult.GroupsAdded, isResult.GroupsRemoved, isResult.GroupsUpdated))

	return nil
}

func syncDataAccess(client plugin.PluginClient, targetConfig target.BaseTargetConfig) error {
	cn := strings.Replace(targetConfig.ConnectorName, "/", "-", -1)

	targetFile, err := filepath.Abs(file.CreateUniqueFileName(cn+"-da", "json"))
	if err != nil {
		return err
	}

	targetConfig.Logger.Debug(fmt.Sprintf("Using %q as data access target file", targetFile))

	if targetConfig.DeleteTempFiles {
		defer os.Remove(targetFile)
	}

	config := data_access.DataAccessConfig{
		BaseTargetConfig: targetConfig,
	}

	dar, err := data_access.RetrieveDataAccessListForDataSource(&config, accessRightsLastUpdated, true)
	if err != nil {
		return err
	}

	if dar == nil {
		targetConfig.Logger.Info("No changes in the data access rights recorded since previous sync. Skipping.", "datasource", config.DataSourceId)
		return nil
	}

	accessRightsLastUpdated = dar.LastCalculated

	syncerConfig := dapc.DataAccessSyncConfig{
		ConfigMap:  baseconfig.ConfigMap{Parameters: targetConfig.Parameters},
		Prefix:     "",
		TargetFile: targetFile,
		RunImport:  true, // signal syncer to also run raito import
	}
	syncerConfig.DataAccess = dar

	das, err := client.GetDataAccessSyncer()
	if err != nil {
		return err
	}

	res := das.SyncDataAccess(&syncerConfig)
	if res.Error != nil {
		return err
	}

	importerConfig := access_provider.AccessProviderImportConfig{
		BaseTargetConfig: targetConfig,
		TargetFile:       targetFile,
		DeleteUntouched:  targetConfig.DeleteUntouched,
	}
	daImporter := access_provider.NewAccessProviderImporter(&importerConfig)

	daResult, err := daImporter.TriggerImport()
	if err != nil {
		return err
	}

	targetConfig.Logger.Info(fmt.Sprintf("Successfully synced access providers. Added: %d - Removed: %d - Updated: %d", daResult.AccessAdded, daResult.AccessRemoved, daResult.AccessUpdated))

	return nil
}

func syncDataUsage(client plugin.PluginClient, targetConfig target.BaseTargetConfig) error {
	cn := strings.Replace(targetConfig.ConnectorName, "/", "-", -1)
	targetFile, err := filepath.Abs(file.CreateUniqueFileName(cn+"-du", "json"))

	if err != nil {
		return err
	}

	targetConfig.Logger.Debug(fmt.Sprintf("Using %q as data usage target file", targetFile))

	if targetConfig.DeleteTempFiles {
		defer os.Remove(targetFile)
	}

	syncerConfig := dupc.DataUsageSyncConfig{
		ConfigMap:  baseconfig.ConfigMap{Parameters: targetConfig.Parameters},
		TargetFile: targetFile,
	}

	dus, err := client.GetDataUsageSyncer()
	if err != nil {
		return err
	}

	res := dus.SyncDataUsage(&syncerConfig)
	if res.Error != nil {
		return err
	}

	importerConfig := data_usage.DataUsageImportConfig{
		BaseTargetConfig: targetConfig,
		TargetFile:       targetFile,
	}
	duImporter := data_usage.NewDataUsageImporter(&importerConfig)

	duResult, err := duImporter.TriggerImport()
	if err != nil {
		return err
	}

	targetConfig.Logger.Info(fmt.Sprintf("Successfully synced data usage. %d statements added, %d failed",
		duResult.StatementsAdded, duResult.StatementsFailed))

	return nil
}
