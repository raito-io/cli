package access_provider

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	dapc "github.com/raito-io/cli/common/api/data_access"
	baseconfig "github.com/raito-io/cli/common/util/config"
	"github.com/raito-io/cli/internal/data_access"
	"github.com/raito-io/cli/internal/file"
	"github.com/raito-io/cli/internal/job"
	"github.com/raito-io/cli/internal/plugin"
	"github.com/raito-io/cli/internal/target"
)

type AccessProviderImportResult struct {
	AccessAdded   int `json:"accessAdded"`
	AccessUpdated int `json:"accessUpdated"`
	AccessRemoved int `json:"accessRemoved"`
}

var accessRightsLastUpdated int64 = 0

type DataAccessSync struct {
	TargetConfig  *target.BaseTargetConfig
	JobId         string
	StatusUpdater func(status job.JobStatus)
}

func (s *DataAccessSync) StartSyncAndQueueJob(client plugin.PluginClient) (job.JobStatus, error) {
	s.StatusUpdater(job.Started)

	cn := strings.Replace(s.TargetConfig.ConnectorName, "/", "-", -1)

	targetFile, err := filepath.Abs(file.CreateUniqueFileName(cn+"-da", "json"))
	if err != nil {
		return job.Failed, err
	}

	s.TargetConfig.Logger.Debug(fmt.Sprintf("Using %q as data access target file", targetFile))

	if s.TargetConfig.DeleteTempFiles {
		defer os.RemoveAll(targetFile)
	}

	config := data_access.DataAccessConfig{
		BaseTargetConfig: *s.TargetConfig,
	}

	s.StatusUpdater(job.DataRetrieve)
	s.TargetConfig.Logger.Info("Fetching access providers for this data source from Raito")
	dar, err := data_access.RetrieveDataAccessListForDataSource(&config, accessRightsLastUpdated, true)

	if err != nil {
		return job.Failed, err
	}

	if dar == nil {
		s.TargetConfig.Logger.Info("No changes in the access providers recorded since previous sync. Skipping.", "datasource", config.DataSourceId)
		return job.Failed, nil
	}

	accessRightsLastUpdated = dar.LastCalculated

	syncerConfig := dapc.DataAccessSyncConfig{
		ConfigMap:  baseconfig.ConfigMap{Parameters: s.TargetConfig.Parameters},
		Prefix:     "",
		TargetFile: targetFile,
		RunImport:  true, // signal syncer to also run raito import
	}
	syncerConfig.DataAccess = dar

	das, err := client.GetDataAccessSyncer()
	if err != nil {
		return job.Failed, err
	}

	s.TargetConfig.Logger.Info("Synchronizing access providers between Raito and the data source")
	res := das.SyncDataAccess(&syncerConfig)

	if res.Error != nil {
		return job.Failed, err
	}

	importerConfig := AccessProviderImportConfig{
		BaseTargetConfig: *s.TargetConfig,
		TargetFile:       targetFile,
		DeleteUntouched:  s.TargetConfig.DeleteUntouched,
	}
	daImporter := NewAccessProviderImporter(&importerConfig, s.StatusUpdater)

	status, err := daImporter.TriggerImport(s.JobId)
	if err != nil {
		return job.Failed, err
	}

	s.TargetConfig.Logger.Info("Successfully queued import job. Wait until remote processing is done.")
	s.TargetConfig.Logger.Debug(fmt.Sprintf("Current status: %s", status.String()))

	return status, nil
}

func (s *DataAccessSync) ProcessResults(results interface{}) error {
	if daResult, ok := results.(*AccessProviderImportResult); ok {
		s.TargetConfig.Logger.Info(fmt.Sprintf("Successfully synced access providers. Added: %d - Removed: %d - Updated: %d", daResult.AccessAdded, daResult.AccessRemoved, daResult.AccessUpdated))
		return nil
	}

	return fmt.Errorf("failed to load results")
}

func (s *DataAccessSync) GetResultObject() interface{} {
	return &AccessProviderImportResult{}
}
