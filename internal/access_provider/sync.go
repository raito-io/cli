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

type AccessProviderExportFeedbackResult struct {
	AccessNamesAdded int `json:"accessNamesAdded"`

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
	result := []job.TaskPart{
		&dataAccessExportSubtask{TargetConfig: s.TargetConfig, JobId: &s.JobId},
	}

	if !s.TargetConfig.SkipDataAccessImport {
		result = append(result, &dataAccessImportSubtask{TargetConfig: s.TargetConfig, JobId: &s.JobId})
	}

	return result
}

func (s *dataAccessImportSubtask) StartSyncAndQueueTaskPart(client plugin.PluginClient, statusUpdater job.TaskEventUpdater) (job.JobStatus, string, error) {
	cn := strings.Replace(s.TargetConfig.ConnectorName, "/", "-", -1)

	targetFile, err := filepath.Abs(file.CreateUniqueFileName(cn+"-da", "json"))
	if err != nil {
		return job.Failed, "", err
	}

	s.TargetConfig.TargetLogger.Debug(fmt.Sprintf("Using %q as data access target file", targetFile))

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
		s.TargetConfig.TargetLogger.Info("Successfully queued import job. Wait until remote processing is done.")
	}

	s.TargetConfig.TargetLogger.Debug(fmt.Sprintf("Current status: %s", status.String()))

	return status, subtaskId, nil
}

func (s *dataAccessImportSubtask) ProcessResults(results interface{}) error {
	if daResult, ok := results.(*AccessProviderImportResult); ok {
		if len(daResult.Warnings) > 0 {
			s.TargetConfig.TargetLogger.Info(fmt.Sprintf("Synced access providers with %d warnings (see below). Added: %d - Removed: %d - Updated: %d", len(daResult.Warnings), daResult.AccessAdded, daResult.AccessRemoved, daResult.AccessUpdated))

			for _, warning := range daResult.Warnings {
				s.TargetConfig.TargetLogger.Warn(warning)
			}
		} else {
			s.TargetConfig.TargetLogger.Info(fmt.Sprintf("Successfully synced access providers. Added: %d - Removed: %d - Updated: %d", daResult.AccessAdded, daResult.AccessRemoved, daResult.AccessUpdated))
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
	syncerConfig := dapc.AccessSyncFromTarget{
		ConfigMap:  baseconfig.ConfigMap{Parameters: s.TargetConfig.Parameters},
		Prefix:     "",
		TargetFile: targetFile,
	}

	das, err := client.GetAccessSyncer()
	if err != nil {
		return err
	}

	s.TargetConfig.TargetLogger.Info("Synchronizing access providers between data source and the Raito")
	res := das.SyncFromTarget(&syncerConfig)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (s *dataAccessExportSubtask) StartSyncAndQueueTaskPart(client plugin.PluginClient, statusUpdater job.TaskEventUpdater) (job.JobStatus, string, error) {
	cn := strings.Replace(s.TargetConfig.ConnectorName, "/", "-", -1)

	targetFile, err := filepath.Abs(file.CreateUniqueFileName(cn+"-da-feedback", "json"))
	if err != nil {
		return job.Failed, "", err
	}

	if s.TargetConfig.DeleteTempFiles {
		defer os.RemoveAll(targetFile)
	}

	s.TargetConfig.TargetLogger.Debug(fmt.Sprintf("Using %q as actual access name target file", targetFile))

	statusUpdater.AddTaskEvent(job.DataRetrieve)

	return s.accessSyncExport(client, statusUpdater, targetFile)
}

// Export data from Raito to DS
func (s *dataAccessExportSubtask) accessSyncExport(client plugin.PluginClient, statusUpdater job.TaskEventUpdater, targetFile string) (jobStatus job.JobStatus, subTaskId string, returnErr error) {
	subTaskUpdater := statusUpdater.GetSubtaskEventUpdater(constants.SubtaskAccessSync)

	defer func() {
		if returnErr != nil {
			subTaskUpdater.AddSubtaskEvent(job.Failed)
		} else {
			subTaskUpdater.AddSubtaskEvent(job.Completed)
		}
	}()

	subTaskUpdater.AddSubtaskEvent(job.Started)

	s.TargetConfig.TargetLogger.Info("Fetching access providers for this data source from Raito")

	statusUpdater.AddTaskEvent(job.DataRetrieve)

	daExporter := NewAccessProviderExporter(&AccessProviderExporterConfig{BaseTargetConfig: *s.TargetConfig}, statusUpdater)

	_, dar, err := daExporter.TriggerExport(*s.JobId)

	if err != nil {
		return job.Failed, "", err
	}

	subTaskUpdater.AddSubtaskEvent(job.InProgress)

	darInformation, err := s.readDataAccessRetrieveInformation(dar)
	if err != nil {
		return job.Failed, "", err
	}

	subTaskUpdater.SetReceivedDate(darInformation.FileBuildTime)
	s.updateLastCalculated(darInformation)

	syncerConfig := dapc.AccessSyncToTarget{
		ConfigMap:          baseconfig.ConfigMap{Parameters: s.TargetConfig.Parameters},
		Prefix:             "",
		SourceFile:         dar,
		FeedbackTargetFile: targetFile,
	}

	das, err := client.GetAccessSyncer()
	if err != nil {
		return job.Failed, "", err
	}

	s.TargetConfig.TargetLogger.Info("Synchronizing access providers between Raito and the data source")
	res := das.SyncToTarget(&syncerConfig)

	if res.Error != nil {
		return job.Failed, "", res.Error
	}

	feedbackImportConfig := AccessProviderExportFeedbackConfig{
		BaseTargetConfig: *s.TargetConfig,
		FeedbackFile:     targetFile,
	}

	importer := NewAccessProviderFeedbackImporter(&feedbackImportConfig, statusUpdater)

	status, subtaskId, err := importer.TriggerFeedbackImport(*s.JobId)
	if err != nil {
		return job.Failed, "", err
	}

	if status == job.Queued {
		s.TargetConfig.TargetLogger.Info("Successfully queued feedback import job. Wait until remote processing is done.")
	}

	s.TargetConfig.TargetLogger.Debug(fmt.Sprintf("Current status: %s", status.String()))

	return status, subtaskId, nil
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
	if daResult, ok := results.(*AccessProviderExportFeedbackResult); ok {
		if len(daResult.Warnings) > 0 {
			s.TargetConfig.TargetLogger.Info(fmt.Sprintf("Exported access providers with %d warnings (see below). Added Actual Names: %d", len(daResult.Warnings), daResult.AccessNamesAdded))

			for _, warning := range daResult.Warnings {
				s.TargetConfig.TargetLogger.Warn(warning)
			}
		} else {
			s.TargetConfig.TargetLogger.Info(fmt.Sprintf("Exported access providers. Added Actual Names: %d", daResult.AccessNamesAdded))
		}

		return nil
	}

	return fmt.Errorf("failed to load results")
}

func (s *dataAccessExportSubtask) GetResultObject() interface{} {
	return &AccessProviderExportFeedbackResult{}
}
