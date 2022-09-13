package data_source

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	dspc "github.com/raito-io/cli/base/data_source"
	baseconfig "github.com/raito-io/cli/base/util/config"
	"github.com/raito-io/cli/internal/file"
	"github.com/raito-io/cli/internal/job"
	"github.com/raito-io/cli/internal/plugin"
	"github.com/raito-io/cli/internal/target"
)

type DataSourceImportResult struct {
	DataObjectsAdded   int `json:"dataObjectsAdded"`
	DataObjectsUpdated int `json:"dataObjectsUpdated"`
	DataObjectsRemoved int `json:"dataObjectsRemoved"`

	Warnings []string `json:"warnings"`
}

type DataSourceSync struct {
	TargetConfig *target.BaseTargetConfig
	JobId        string
}

func (s *DataSourceSync) StartSyncAndQueueJob(client plugin.PluginClient, statusUpdater job.TaskEventUpdater) (job.JobStatus, string, error) {
	cn := strings.Replace(s.TargetConfig.ConnectorName, "/", "-", -1)

	targetFile, err := filepath.Abs(file.CreateUniqueFileName(cn+"-ds", "json"))
	if err != nil {
		return job.Failed, "", err
	}

	s.TargetConfig.Logger.Debug(fmt.Sprintf("Using %q as data source target file", targetFile))

	if s.TargetConfig.DeleteTempFiles {
		defer os.RemoveAll(targetFile)
	}

	syncerConfig := dspc.DataSourceSyncConfig{
		ConfigMap:    baseconfig.ConfigMap{Parameters: s.TargetConfig.Parameters},
		TargetFile:   targetFile,
		DataSourceId: s.TargetConfig.DataSourceId,
	}

	dss, err := client.GetDataSourceSyncer()
	if err != nil {
		return job.Failed, "", err
	}

	s.TargetConfig.Logger.Info("Fetching data source metadata configuration")
	md := dss.GetMetaData()

	s.TargetConfig.Logger.Info("Updating data source metadata configuration")
	err = SetMetaData(*s.TargetConfig, md)

	if err != nil {
		return job.Failed, "", err
	}

	s.TargetConfig.Logger.Info("Gathering metadata from the data source")
	res := dss.SyncDataSource(&syncerConfig)

	if res.Error != nil {
		return job.Failed, "", err
	}

	importerConfig := DataSourceImportConfig{
		BaseTargetConfig: *s.TargetConfig,
		TargetFile:       targetFile,
		DeleteUntouched:  s.TargetConfig.DeleteUntouched,
		ReplaceTags:      s.TargetConfig.ReplaceTags,
	}
	dsImporter := NewDataSourceImporter(&importerConfig, statusUpdater)

	s.TargetConfig.Logger.Info("Importing metadata into Raito")
	status, subtaskId, err := dsImporter.TriggerImport(s.JobId)

	if err != nil {
		return job.Failed, "", err
	}

	if status == job.Queued {
		s.TargetConfig.Logger.Info("Successfully queued import job. Wait until remote processing is done.")
	}

	s.TargetConfig.Logger.Debug(fmt.Sprintf("Current status: %s", status))

	return status, subtaskId, nil
}

func (s *DataSourceSync) ProcessResults(results interface{}) error {
	if dsResult, ok := results.(*DataSourceImportResult); ok {
		if dsResult.Warnings != nil && len(dsResult.Warnings) > 0 {
			s.TargetConfig.Logger.Info(fmt.Sprintf("Synced data source with %d warnings (see below). Added: %d - Removed: %d - Updated: %d", len(dsResult.Warnings), dsResult.DataObjectsAdded, dsResult.DataObjectsRemoved, dsResult.DataObjectsUpdated))

			for _, warning := range dsResult.Warnings {
				s.TargetConfig.Logger.Warn(warning)
			}
		} else {
			s.TargetConfig.Logger.Info(fmt.Sprintf("Successfully synced data source. Added: %d - Removed: %d - Updated: %d", dsResult.DataObjectsAdded, dsResult.DataObjectsRemoved, dsResult.DataObjectsUpdated))
		}

		return nil
	}

	return fmt.Errorf("failed to load results")
}

func (s *DataSourceSync) GetResultObject() interface{} {
	return &DataSourceImportResult{}
}
