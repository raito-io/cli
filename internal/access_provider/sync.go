package access_provider

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-hclog"
	"gopkg.in/yaml.v2"

	dapc "github.com/raito-io/cli/base/access_provider"
	baseconfig "github.com/raito-io/cli/base/util/config"
	error1 "github.com/raito-io/cli/base/util/error"
	"github.com/raito-io/cli/base/util/slice"
	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/job"
	"github.com/raito-io/cli/internal/plugin"
	"github.com/raito-io/cli/internal/target/types"
	"github.com/raito-io/cli/internal/util/file"
	"github.com/raito-io/cli/internal/version_management"
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
	TargetConfig *types.BaseTargetConfig
	JobId        string

	result []job.TaskResult
}

type dataAccessRetrieveInformation struct {
	LastCalculated int64 `yaml:"lastCalculated" json:"lastCalculated"`
	FileBuildTime  int64 `yaml:"fileBuildTime" json:"fileBuildTime"`
}

type dataAccessExportSubtask struct {
	TargetConfig *types.BaseTargetConfig
	JobId        *string

	task *DataAccessSync
}

type dataAccessImportSubtask struct {
	TargetConfig *types.BaseTargetConfig
	JobId        *string

	task *DataAccessSync
}

func (s *DataAccessSync) IsClientValid(ctx context.Context, c plugin.PluginClient) (bool, error) {
	accessSyncer, err := c.GetAccessSyncer()
	if err != nil {
		return false, err
	}

	return version_management.IsValidToSync(ctx, accessSyncer, dapc.MinimalCliVersion)
}

func (s *DataAccessSync) GetParts() []job.TaskPart {
	result := []job.TaskPart{
		&dataAccessExportSubtask{TargetConfig: s.TargetConfig, JobId: &s.JobId, task: s},
	}

	if !s.TargetConfig.SkipDataAccessImport {
		result = append(result, &dataAccessImportSubtask{TargetConfig: s.TargetConfig, JobId: &s.JobId, task: s})
	}

	return result
}

func (s *DataAccessSync) GetTaskResults() []job.TaskResult {
	return s.result
}

func (s *dataAccessImportSubtask) StartSyncAndQueueTaskPart(ctx context.Context, logger hclog.Logger, client plugin.PluginClient, statusUpdater job.TaskEventUpdater, secureImport func(func() error) error) (job.JobStatus, string, error) {
	targetFile, err := filepath.Abs(file.CreateUniqueFileNameForTarget(s.TargetConfig.Name, "fromTarget-access", "json"))
	if err != nil {
		return job.Failed, "", err
	}

	logger.Debug(fmt.Sprintf("Using %q as data access target file", targetFile))

	defer s.TargetConfig.HandleTempFile(logger, targetFile, false)

	err = s.accessSyncImport(logger, client, targetFile)
	if err != nil {
		return job.Failed, "", err
	}

	postProcessor := NewPostProcessor(&PostProcessorConfig{
		TagOverwriteKeyForName:   s.TargetConfig.TagOverwriteKeyForAccessProviderName,
		TagOverwriteKeyForOwners: s.TargetConfig.TagOverwriteKeyForAccessProviderOwners,
	})

	toProcessFile := targetFile
	if postProcessor.NeedsPostProcessing() {
		toProcessFile, _, err = s.postProcessAccessProviders(logger, postProcessor, targetFile)
		if err != nil {
			return job.Failed, "", err
		}
	}

	importerConfig := AccessProviderImportConfig{
		BaseTargetConfig: *s.TargetConfig,
		TargetFile:       toProcessFile,
		DeleteUntouched:  s.TargetConfig.DeleteUntouched,
	}

	daImporter := NewAccessProviderImporter(&importerConfig, statusUpdater)

	var status job.JobStatus
	var subtaskId string

	err = secureImport(func() error {
		status, subtaskId, err = daImporter.TriggerImport(ctx, logger, *s.JobId)
		if err != nil {
			return err
		}

		if status == job.Queued {
			logger.Info("Successfully queued import job. Wait until remote processing is done.")
		}

		logger.Debug(fmt.Sprintf("Current status: %s", status.String()))

		return nil
	})
	if err != nil {
		return job.Failed, "", err
	}

	return status, subtaskId, nil
}

func (s *dataAccessImportSubtask) postProcessAccessProviders(logger hclog.Logger, postProcessor PostProcessor, toProcessFile string) (string, int, error) {
	postProcessedFile := toProcessFile
	fileSuffix := "-post-processed"

	// Generate a unique file name for the post processing
	if strings.Contains(postProcessedFile, fileSuffix) {
		postProcessedFile = postProcessedFile[0:strings.LastIndex(postProcessedFile, fileSuffix)] + fileSuffix + ".json"
	} else {
		postProcessedFile = postProcessedFile[0:strings.LastIndex(postProcessedFile, ".json")] + fileSuffix + ".json"
	}

	res, err := postProcessor.PostProcess(logger, toProcessFile, postProcessedFile)
	if err != nil {
		return toProcessFile, 0, err
	}

	return postProcessedFile, res.AccessProviderTouchedCount, nil
}

func (s *dataAccessImportSubtask) ProcessResults(logger hclog.Logger, results interface{}) error {
	if daResult, ok := results.(*AccessProviderImportResult); ok {
		if len(daResult.Warnings) > 0 {
			logger.Info(fmt.Sprintf("Synced access providers with %d warnings (see below). Added: %d - Removed: %d - Updated: %d", len(daResult.Warnings), daResult.AccessAdded, daResult.AccessRemoved, daResult.AccessUpdated))

			for _, warning := range daResult.Warnings {
				logger.Warn(warning)
			}
		} else {
			logger.Info(fmt.Sprintf("Successfully synced access providers. Added: %d - Removed: %d - Updated: %d", daResult.AccessAdded, daResult.AccessRemoved, daResult.AccessUpdated))
		}

		s.task.result = append(s.task.result, job.TaskResult{
			ObjectType: "imported access providers",
			Added:      daResult.AccessAdded,
			Updated:    daResult.AccessUpdated,
			Removed:    daResult.AccessRemoved,
			Failed:     len(daResult.Warnings),
		})

		return nil
	}

	return fmt.Errorf("failed to load results")
}

func (s *dataAccessImportSubtask) GetResultObject() interface{} {
	return &AccessProviderImportResult{}
}

// Import data from Raito to DS
func (s *dataAccessImportSubtask) accessSyncImport(logger hclog.Logger, client plugin.PluginClient, targetFile string) (returnErr error) {
	syncerConfig := dapc.AccessSyncFromTarget{
		ConfigMap:                     &baseconfig.ConfigMap{Parameters: s.TargetConfig.Parameters},
		Prefix:                        "",
		TargetFile:                    targetFile,
		LockAllWho:                    s.TargetConfig.LockAllWho,
		LockAllInheritance:            s.TargetConfig.LockAllInheritance,
		LockAllWhat:                   s.TargetConfig.LockAllWhat,
		LockAllNames:                  s.TargetConfig.LockAllNames,
		LockAllDelete:                 s.TargetConfig.LockAllDelete,
		LockAllOwners:                 s.TargetConfig.LockAllOwners,
		MakeNotInternalizable:         slice.ParseCommaSeparatedList(s.TargetConfig.MakeNotInternalizable),
		LockWhoByName:                 slice.ParseCommaSeparatedList(s.TargetConfig.LockWhoByName),
		LockWhoByTag:                  slice.ParseCommaSeparatedList(s.TargetConfig.LockWhoByTag),
		LockWhoWhenIncomplete:         s.TargetConfig.LockWhoWhenIncomplete,
		LockWhatByName:                slice.ParseCommaSeparatedList(s.TargetConfig.LockWhatByName),
		LockWhatByTag:                 slice.ParseCommaSeparatedList(s.TargetConfig.LockWhatByTag),
		LockWhatWhenIncomplete:        s.TargetConfig.LockWhatWhenIncomplete,
		LockInheritanceByName:         slice.ParseCommaSeparatedList(s.TargetConfig.LockInheritanceByName),
		LockInheritanceByTag:          slice.ParseCommaSeparatedList(s.TargetConfig.LockInheritanceByTag),
		LockInheritanceWhenIncomplete: s.TargetConfig.LockInheritanceWhenIncomplete,
		FullyLockAll:                  s.TargetConfig.FullyLockAll,
		FullyLockByName:               slice.ParseCommaSeparatedList(s.TargetConfig.FullyLockByName),
		FullyLockByTag:                slice.ParseCommaSeparatedList(s.TargetConfig.FullyLockByTag),
		FullyLockWhenIncomplete:       s.TargetConfig.FullyLockWhenIncomplete,
		LockNamesByName:               slice.ParseCommaSeparatedList(s.TargetConfig.LockNamesByName),
		LockNamesByTag:                slice.ParseCommaSeparatedList(s.TargetConfig.LockNamesByTag),
		LockNamesWhenIncomplete:       s.TargetConfig.LockNamesWhenIncomplete,
		LockDeleteByName:              slice.ParseCommaSeparatedList(s.TargetConfig.LockDeleteByName),
		LockDeleteByTag:               slice.ParseCommaSeparatedList(s.TargetConfig.LockDeleteByTag),
		LockDeleteWhenIncomplete:      s.TargetConfig.LockDeleteWhenIncomplete,
	}

	das, err := client.GetAccessSyncer()
	if err != nil {
		return err
	}

	logger.Info("Synchronizing access providers between data source and Raito")

	res, err := das.SyncFromTarget(context.Background(), &syncerConfig)
	if err != nil {
		return err
	}

	if res.Error != nil { //nolint:staticcheck
		return mapErrorResult(res.Error) //nolint:staticcheck
	}

	return nil
}

func (s *dataAccessExportSubtask) StartSyncAndQueueTaskPart(ctx context.Context, logger hclog.Logger, client plugin.PluginClient, statusUpdater job.TaskEventUpdater, secureImport func(func() error) error) (job.JobStatus, string, error) {
	targetFile, err := filepath.Abs(file.CreateUniqueFileNameForTarget(s.TargetConfig.Name, "toTarget-accessFeedback", "json"))
	if err != nil {
		return job.Failed, "", err
	}

	defer s.TargetConfig.HandleTempFile(logger, targetFile, false)

	logger.Debug(fmt.Sprintf("Using %q as actual access name target file", targetFile))

	statusUpdater.SetStatusToDataRetrieve(ctx)

	return s.accessSyncExport(ctx, logger, client, statusUpdater, targetFile, secureImport)
}

// Export data from Raito to DS
func (s *dataAccessExportSubtask) accessSyncExport(ctx context.Context, logger hclog.Logger, client plugin.PluginClient, statusUpdater job.TaskEventUpdater, targetFile string, secureImport func(func() error) error) (_ job.JobStatus, _ string, returnErr error) {
	subTaskUpdater := statusUpdater.GetSubtaskEventUpdater(constants.SubtaskAccessSync)

	defer func() {
		if returnErr != nil {
			logger.Error(fmt.Sprintf("Access provider sync failed due to error: %s", returnErr.Error()))
			subTaskUpdater.AddSubtaskEvent(ctx, job.Failed)
		} else {
			subTaskUpdater.AddSubtaskEvent(ctx, job.Completed)
		}
	}()

	subTaskUpdater.AddSubtaskEvent(ctx, job.Started)

	logger.Info("Loading plugin")

	das, err := client.GetAccessSyncer()
	if err != nil {
		return job.Failed, "", err
	}

	logger.Info("Fetching access providers for this data source from Raito")

	statusUpdater.SetStatusToDataRetrieve(ctx)

	syncConfig, err := das.SyncConfig(context.Background())
	if err != nil {
		return job.Failed, "", err
	}

	daExporter := NewAccessProviderExporter(&AccessProviderExporterConfig{BaseTargetConfig: *s.TargetConfig}, statusUpdater, syncConfig)

	_, exportedFile, err := daExporter.TriggerExport(ctx, logger, *s.JobId)

	if err != nil {
		return job.Failed, "", err
	}

	defer s.TargetConfig.HandleTempFile(logger, exportedFile, false)

	subTaskUpdater.AddSubtaskEvent(ctx, job.InProgress)

	darInformation, err := s.readDataAccessRetrieveInformation(exportedFile)
	if err != nil {
		return job.Failed, "", err
	}

	subTaskUpdater.SetReceivedDate(darInformation.FileBuildTime)
	s.updateLastCalculated(darInformation)

	syncerConfig := dapc.AccessSyncToTarget{
		ConfigMap:          &baseconfig.ConfigMap{Parameters: s.TargetConfig.Parameters},
		Prefix:             "",
		SourceFile:         exportedFile,
		FeedbackTargetFile: targetFile,
	}

	logger.Info("Synchronizing access providers between Raito and the data source")

	res, err := das.SyncToTarget(context.Background(), &syncerConfig)
	if err != nil {
		return job.Failed, "", err
	}

	if res.Error != nil { //nolint:staticcheck
		return job.Failed, "", mapErrorResult(res.Error) //nolint:staticcheck
	}

	s.task.result = append(s.task.result, job.TaskResult{
		ObjectType: "exported access providers",
		Added:      int(res.AccessProviderCount),
	})

	feedbackImportConfig := AccessProviderExportFeedbackConfig{
		BaseTargetConfig: *s.TargetConfig,
		FeedbackFile:     targetFile,
	}

	importer := NewAccessProviderFeedbackImporter(&feedbackImportConfig, statusUpdater)

	var status job.JobStatus
	var subtaskId string

	err = secureImport(func() error {
		status, subtaskId, err = importer.TriggerFeedbackImport(ctx, logger, *s.JobId)
		if err != nil {
			return err
		}

		if status == job.Queued {
			logger.Info("Successfully queued feedback import job. Wait until remote processing is done.")
		}

		logger.Debug(fmt.Sprintf("Current status: %s", status.String()))

		return nil
	})
	if err != nil {
		return job.Failed, "", err
	}

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

func (s *dataAccessExportSubtask) ProcessResults(logger hclog.Logger, results interface{}) error {
	if daResult, ok := results.(*AccessProviderExportFeedbackResult); ok {
		if len(daResult.Warnings) > 0 {
			logger.Info(fmt.Sprintf("Exported access providers with %d warnings (see below). Added Actual Names: %d", len(daResult.Warnings), daResult.AccessNamesAdded))

			for _, warning := range daResult.Warnings {
				logger.Warn(warning)
			}
		} else {
			logger.Info(fmt.Sprintf("Exported access providers. Added Actual Names: %d", daResult.AccessNamesAdded))
		}

		return nil
	}

	return fmt.Errorf("failed to load results")
}

func (s *dataAccessExportSubtask) GetResultObject() interface{} {
	return &AccessProviderExportFeedbackResult{}
}

func mapErrorResult(result *error1.ErrorResult) error {
	if result == nil {
		return nil
	}

	return errors.New(result.ErrorMessage)
}
