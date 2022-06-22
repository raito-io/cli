package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/raito-io/cli/internal/graphql"
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

func startJob(cfg *target.BaseTargetConfig) (string, error) {
	gqlQuery := fmt.Sprintf(`{ "query": "mutation createJob {
        createJob(input: { dataSourceId: \"%s\", eventTime: \"%s\" }) { jobId } }" }"`,
		cfg.DataSourceId, time.Now().Format(time.RFC3339))

	gqlQuery = strings.ReplaceAll(gqlQuery, "\n", "\\n")

	res, err := graphql.ExecuteGraphQL(gqlQuery, cfg)
	if err != nil {
		return "", fmt.Errorf("error while executing import: %s", err.Error())
	}

	resp := Response{}
	gr := graphql.GraphqlResponse{Data: &resp}

	err = json.Unmarshal(res, &gr)
	if err != nil {
		return "", fmt.Errorf("error while parsing job event result: %s", err.Error())
	}

	return *resp.Job.JobID, nil
}

type Response struct {
	Job Job `json:"createJob"`
}

type Job struct {
	JobID *string `json:"jobId"`
}

func addJobEvent(cfg *target.BaseTargetConfig, jobID, jobType, status string) {
	gqlQuery := fmt.Sprintf(`{ "query": "mutation createJobEvent {
        createJobEvent(input: { jobId: \"%s\", dataSourceId: \"%s\", jobType: \"%s\", status: \"%s\", eventTime: \"%s\" }) { jobId } }" }"`,
		jobID, cfg.DataSourceId, jobType, status, time.Now().Format(time.RFC3339))

	gqlQuery = strings.ReplaceAll(gqlQuery, "\n", "\\n")

	id, err := graphql.ExecuteGraphQL(gqlQuery, cfg)
	if err != nil {
		cfg.Logger.Info("job update failed: %s", err.Error())
	}

	cfg.Logger.Info("Job ID: " + string(id))
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

	jobID, jobErr := startJob(targetConfig)
	// TODO: Don't send jobEvents if job creation failed
	if jobErr != nil {
		targetConfig.Logger.Warn(fmt.Sprintf("Error creating status job: %s", err.Error()))
	}

	if targetConfig.DataSourceId != "" && !targetConfig.SkipDataSourceSync {
		targetConfig.Logger.Info("Synchronizing data source meta data...")
		addJobEvent(targetConfig, jobID, constants.DataSourceSync, constants.Started)

		err := syncDataSource(&client, *targetConfig)
		if err != nil {
			target.HandleTargetError(err, targetConfig, "synchronizing data source meta data")
			addJobEvent(targetConfig, jobID, constants.DataSourceSync, constants.Failed)

			return err
		} else {
			addJobEvent(targetConfig, jobID, constants.DataSourceSync, constants.Completed)
		}
	} else {
		addJobEvent(targetConfig, jobID, constants.DataSourceSync, constants.Skipped)
		if targetConfig.DataSourceId == "" {
			targetConfig.Logger.Info("No data-source-id argument found. Skipping syncing of data source meta data")
		} else {
			targetConfig.Logger.Info("Skipping syncing of data source meta data")
		}
	}

	if targetConfig.IdentityStoreId != "" && !targetConfig.SkipIdentityStoreSync {
		targetConfig.Logger.Info("Synchronizing identity store data...")
		addJobEvent(targetConfig, jobID, constants.IdentitySync, constants.Started)

		err := syncIdentityStore(&client, *targetConfig)
		if err != nil {
			target.HandleTargetError(err, targetConfig, "sychronizing identity store data")
			addJobEvent(targetConfig, jobID, constants.IdentitySync, constants.Failed)

			return err
		} else {
			addJobEvent(targetConfig, jobID, constants.IdentitySync, constants.Completed)
		}
	} else {
		addJobEvent(targetConfig, jobID, constants.IdentitySync, constants.Skipped)

		if targetConfig.DataSourceId == "" {
			targetConfig.Logger.Info("No identity-store-id argument found. Skipping identity store syncing")
		} else {
			targetConfig.Logger.Info("Skipping identity store syncing")
		}
	}

	if targetConfig.DataSourceId != "" && !targetConfig.SkipDataAccessSync {
		targetConfig.Logger.Info("Synchronizing data access...")
		addJobEvent(targetConfig, jobID, constants.DataAccessSync, constants.Started)

		err := syncDataAccess(&client, *targetConfig)
		if err != nil {
			target.HandleTargetError(err, targetConfig, "sychronizing data access information to the data source")
			addJobEvent(targetConfig, jobID, constants.DataAccessSync, constants.Failed)

			return err
		} else {
			addJobEvent(targetConfig, jobID, constants.DataAccessSync, constants.Completed)
		}
	} else {
		addJobEvent(targetConfig, jobID, constants.DataAccessSync, constants.Skipped)

		if targetConfig.DataSourceId == "" {
			targetConfig.Logger.Info("No data-source-id argument found. Skipping data access syncing")
		} else {
			targetConfig.Logger.Info("Skipping data access syncing")
		}
	}

	if targetConfig.DataSourceId != "" && !targetConfig.SkipDataUsageSync {
		targetConfig.Logger.Info("Synchronizing data usage...")
		addJobEvent(targetConfig, jobID, constants.DataUsageSync, constants.Started)

		err := syncDataUsage(&client, *targetConfig)
		if err != nil {
			target.HandleTargetError(err, targetConfig, "sychronizing data usage information to the data source")
			addJobEvent(targetConfig, jobID, constants.DataUsageSync, constants.Failed)

			return err
		} else {
			addJobEvent(targetConfig, jobID, constants.DataUsageSync, constants.Completed)
		}
	} else {
		addJobEvent(targetConfig, jobID, constants.DataUsageSync, constants.Skipped)

		if targetConfig.DataSourceId == "" {
			targetConfig.Logger.Info("No data-source-id argument found. Skipping data usage syncing")
		} else {
			targetConfig.Logger.Info("Skipping data usage syncing")
		}
	}

	targetConfig.Logger.Info(fmt.Sprintf("Successfully finished execution in %s", time.Since(start).Round(time.Millisecond)), "success")

	// TODO: If one fails, fail the whole job?
	addJobEvent(targetConfig, jobID, constants.Job, constants.Completed)

	return nil
}

func syncDataSource(client *plugin.PluginClient, targetConfig target.BaseTargetConfig) error {
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
	dss, err := (*client).GetDataSourceSyncer()
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

func syncIdentityStore(client *plugin.PluginClient, targetConfig target.BaseTargetConfig) error {
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

	iss, err := (*client).GetIdentityStoreSyncer()
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

func syncDataAccess(client *plugin.PluginClient, targetConfig target.BaseTargetConfig) error {
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

	das, err := (*client).GetDataAccessSyncer()
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

func syncDataUsage(client *plugin.PluginClient, targetConfig target.BaseTargetConfig) error {
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
	dus, err := (*client).GetDataUsageSyncer()
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
