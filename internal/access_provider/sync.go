package access_provider

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	dapc "github.com/raito-io/cli/base/access_provider"
	baseconfig "github.com/raito-io/cli/base/util/config"
	"github.com/raito-io/cli/internal/constants"
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

	Warnings []string `json:"warnings"`
}

var accessLastCalculated = map[string]int64{}

type DataAccessSync struct {
	TargetConfig *target.BaseTargetConfig
	JobId        string
}

type dataAccessRetrieveInformation struct {
	LastCalculated int64 `yaml:"lastCalculated" json:"lastCalculated"`
	FileBuildTime  int64 `yaml:"fileBuildTime" json:"fileBuildTime"`
}

type dataAccessExportSubtask struct {
	TargetConfig *target.BaseTargetConfig
	JobId        *string
}

type dataAccessImportSubtask struct {
	TargetConfig *target.BaseTargetConfig
	JobId        *string
}

func (s *DataAccessSync) GetParts() []job.TaskPart {
	return []job.TaskPart{
		&dataAccessExportSubtask{TargetConfig: s.TargetConfig, JobId: &s.JobId},
		&dataAccessImportSubtask{TargetConfig: s.TargetConfig, JobId: &s.JobId},
	}
}

func (s *dataAccessImportSubtask) StartSyncAndQueueTaskPart(client plugin.PluginClient, statusUpdater job.TaskEventUpdater) (job.JobStatus, string, error) {
	cn := strings.Replace(s.TargetConfig.ConnectorName, "/", "-", -1)

	targetFile, err := filepath.Abs(file.CreateUniqueFileName(cn+"-da", "json"))
	if err != nil {
		return job.Failed, "", err
	}

	s.TargetConfig.Logger.Debug(fmt.Sprintf("Using %q as data access target file", targetFile))

	if s.TargetConfig.DeleteTempFiles {
		defer os.RemoveAll(targetFile)
	}

	err = s.accessSyncImport(client, targetFile)
	if err != nil {
		return job.Failed, "", err
	}

	importerConfig := AccessProviderImportConfig{
		BaseTargetConfig: *s.TargetConfig,
		TargetFile:       targetFile,
		DeleteUntouched:  s.TargetConfig.DeleteUntouched,
	}

	daImporter := NewAccessProviderImporter(&importerConfig, statusUpdater)

	status, subtaskId, err := daImporter.TriggerImport(*s.JobId)
	if err != nil {
		return job.Failed, "", err
	}

	if status == job.Queued {
		s.TargetConfig.Logger.Info("Successfully queued import job. Wait until remote processing is done.")
	}

	s.TargetConfig.Logger.Debug(fmt.Sprintf("Current status: %s", status.String()))

	return status, subtaskId, nil
}

func (s *dataAccessImportSubtask) ProcessResults(results interface{}) error {
	if daResult, ok := results.(*AccessProviderImportResult); ok {
		if len(daResult.Warnings) > 0 {
			s.TargetConfig.Logger.Info(fmt.Sprintf("Synced access providers with %d warnings (see below). Added: %d - Removed: %d - Updated: %d", len(daResult.Warnings), daResult.AccessAdded, daResult.AccessRemoved, daResult.AccessUpdated))

			for _, warning := range daResult.Warnings {
				s.TargetConfig.Logger.Warn(warning)
			}
		} else {
			s.TargetConfig.Logger.Info(fmt.Sprintf("Successfully synced access providers. Added: %d - Removed: %d - Updated: %d", daResult.AccessAdded, daResult.AccessRemoved, daResult.AccessUpdated))
		}

		return nil
	}

	return fmt.Errorf("failed to load results")
}

func (s *dataAccessImportSubtask) GetResultObject() interface{} {
	return &AccessProviderImportResult{}
}

// Import data from Raito to DS
func (s *dataAccessImportSubtask) accessSyncImport(client plugin.PluginClient, targetFile string) (returnErr error) {
	syncerConfig := dapc.AccessSyncToTarget{
		ConfigMap:  baseconfig.ConfigMap{Parameters: s.TargetConfig.Parameters},
		Prefix:     "",
		TargetFile: targetFile,
	}

	das, err := client.GetAccessSyncer()
	if err != nil {
		return err
	}

	s.TargetConfig.Logger.Info("Synchronizing access providers between data source and the Raito")
	res := das.SyncToTarget(&syncerConfig)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (s *dataAccessExportSubtask) StartSyncAndQueueTaskPart(client plugin.PluginClient, statusUpdater job.TaskEventUpdater) (job.JobStatus, string, error) {
	cn := strings.Replace(s.TargetConfig.ConnectorName, "/", "-", -1)

	targetFile, err := filepath.Abs(file.CreateUniqueFileName(cn+"-da-naming", "json"))
	if err != nil {
		return job.Failed, "", err
	}

	if s.TargetConfig.DeleteTempFiles {
		defer os.RemoveAll(targetFile)
	}

	s.TargetConfig.Logger.Debug(fmt.Sprintf("Using %q as actual access name target file", targetFile))

	statusUpdater.AddTaskEvent(job.DataRetrieve)

	err = s.accessSyncExport(client, statusUpdater, targetFile)
	if err != nil {
		return job.Failed, "", err
	}

	return job.Completed, "", err
}

// Export data from Raito to DS
func (s *dataAccessExportSubtask) accessSyncExport(client plugin.PluginClient, statusUpdater job.TaskEventUpdater, targetFile string) (returnErr error) {
	subTaskUpdater := statusUpdater.GetSubtaskEventUpdater(constants.SubtaskAccessSync)

	defer func() {
		if returnErr != nil {
			subTaskUpdater.AddSubtaskEvent(job.Failed)
		} else {
			subTaskUpdater.AddSubtaskEvent(job.Completed)
		}
	}()

	subTaskUpdater.AddSubtaskEvent(job.Started)

	config := data_access.AccessSyncConfig{
		BaseTargetConfig: *s.TargetConfig,
	}

	lastUpdated := accessLastCalculated[s.TargetConfig.DataSourceId]

	s.TargetConfig.Logger.Info("Fetching access providers for this data source from Raito")
	subTaskUpdater.AddSubtaskEvent(job.DataRetrieve)
	dar, err := data_access.RetrieveDataAccessListForDataSource(&config, lastUpdated)

	if err != nil {
		return err
	}

	subTaskUpdater.AddSubtaskEvent(job.InProgress)

	darInformation, err := s.readDataAccessRetrieveInformation(dar)
	if err != nil {
		return err
	}

	subTaskUpdater.SetReceivedDate(darInformation.FileBuildTime)
	s.updateLastCalculated(darInformation)

	syncerConfig := dapc.AccessSyncFromTarget{
		ConfigMap:  baseconfig.ConfigMap{Parameters: s.TargetConfig.Parameters},
		Prefix:     "",
		SourceFile: dar,
		TargetFile: targetFile,
	}

	das, err := client.GetAccessSyncer()
	if err != nil {
		return err
	}

	s.TargetConfig.Logger.Info("Synchronizing access providers between Raito and the data source")
	res := das.SyncFromTarget(&syncerConfig)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (s *dataAccessExportSubtask) readDataAccessRetrieveInformation(filePath string) (*dataAccessRetrieveInformation, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	darInf := &dataAccessRetrieveInformation{}
	err = yaml.NewDecoder(file).Decode(darInf)

	return darInf, err
}

func (s *dataAccessExportSubtask) updateLastCalculated(information *dataAccessRetrieveInformation) {
	accessLastCalculated[s.TargetConfig.DataSourceId] = information.LastCalculated
}

func (s *dataAccessExportSubtask) ProcessResults(results interface{}) error {
	return nil
}

func (s *dataAccessExportSubtask) GetResultObject() interface{} {
	return nil
}
