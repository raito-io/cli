package data_usage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"

	dupc "github.com/raito-io/cli/base/data_usage"
	baseconfig "github.com/raito-io/cli/base/util/config"
	"github.com/raito-io/cli/internal/file"
	"github.com/raito-io/cli/internal/job"
	"github.com/raito-io/cli/internal/plugin"
	"github.com/raito-io/cli/internal/target"
)

type DataUsageSync struct {
	TargetConfig  *target.BaseTargetConfig
	JobId         string
	StatusUpdater func(status job.JobStatus)
}

type DataUsageImportResult struct {
	StatementsAdded       int `json:"statementsAdded"`
	StatementsFailed      int `json:"statementsFailed"`
	StatementsSkipped     int `json:"statementsSkipped"`
	EdgesCreatedOrUpdated int `json:"edgesCreatedOrUpdated"`
	EdgesRemoved          int `json:"edgesRemoved"`

	Warnings []string `json:"warnings"`
}

func (s *DataUsageSync) StartSyncAndQueueJob(client plugin.PluginClient) (job.JobStatus, string, error) {
	cn := strings.Replace(s.TargetConfig.ConnectorName, "/", "-", -1)
	targetFile, err := filepath.Abs(file.CreateUniqueFileName(cn+"-du", "json"))

	if err != nil {
		return job.Failed, "", err
	}

	s.TargetConfig.Logger.Debug(fmt.Sprintf("Using %q as data usage target file", targetFile))

	if s.TargetConfig.DeleteTempFiles {
		defer os.RemoveAll(targetFile)
	}

	syncerConfig := dupc.DataUsageSyncConfig{
		ConfigMap:  baseconfig.ConfigMap{Parameters: s.TargetConfig.Parameters},
		TargetFile: targetFile,
	}

	dus, err := client.GetDataUsageSyncer()
	if err != nil {
		return job.Failed, "", err
	}

	importerConfig := DataUsageImportConfig{
		BaseTargetConfig: *s.TargetConfig,
		TargetFile:       targetFile,
	}
	duImporter := NewDataUsageImporter(&importerConfig, s.StatusUpdater)

	s.TargetConfig.Logger.Info("Fetching last synchronization date")

	lastUsed, err := duImporter.GetLastUsage()

	if err != nil || lastUsed == nil {
		hclog.L().Warn(fmt.Sprintf("error retrieving last usage for data source %s, last used: %s", importerConfig.DataSourceId, lastUsed))
		timeValue := time.Unix(int64(0), 0)
		lastUsed = &timeValue
	}

	lastUsedValue := *lastUsed
	syncerConfig.ConfigMap.Parameters["lastUsed"] = lastUsedValue.Format(time.RFC3339)

	s.TargetConfig.Logger.Info("Fetching usage data from the data source")

	res := dus.SyncDataUsage(&syncerConfig)
	if res.Error != nil {
		return job.Failed, "", err
	}

	s.TargetConfig.Logger.Info("Importing usage data into Raito")

	status, subtaskId, err := duImporter.TriggerImport(s.JobId)
	if err != nil {
		return job.Failed, "", err
	}

	if status == job.Queued {
		s.TargetConfig.Logger.Info("Successfully queued import job. Wait until remote processing is done.")
	}

	s.TargetConfig.Logger.Debug("Current status: %s", status.String())

	return status, subtaskId, nil
}

func (s *DataUsageSync) ProcessResults(results interface{}) error {
	if duResult, ok := results.(*DataUsageImportResult); ok {
		if duResult != nil && len(duResult.Warnings) > 0 {
			s.TargetConfig.Logger.Info(fmt.Sprintf("Synced data usage with %d warnings (see below). %d statements added, %d failed", len(duResult.Warnings), duResult.StatementsAdded, duResult.StatementsFailed))

			for _, warning := range duResult.Warnings {
				s.TargetConfig.Logger.Warn(warning)
			}
		} else {
			s.TargetConfig.Logger.Info(fmt.Sprintf("Successfully synced data usage. %d statements added, %d failed",
				duResult.StatementsAdded, duResult.StatementsFailed))
		}

		return nil
	}

	return fmt.Errorf("failed to load results")
}

func (s *DataUsageSync) GetResultObject() interface{} {
	return &DataUsageImportResult{}
}
