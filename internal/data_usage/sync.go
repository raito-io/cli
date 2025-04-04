package data_usage

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/c2h5oh/datasize"
	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"

	dupc "github.com/raito-io/cli/base/data_usage"
	baseconfig "github.com/raito-io/cli/base/util/config"
	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/job"
	"github.com/raito-io/cli/internal/plugin"
	"github.com/raito-io/cli/internal/target/types"
	"github.com/raito-io/cli/internal/util/file"
	"github.com/raito-io/cli/internal/version_management"
)

type DataUsageSync struct {
	TargetConfig *types.BaseTargetConfig
	JobId        string

	result *job.TaskResult
}

type DataUsageImportResult struct {
	StatementsAdded       int `json:"statementsAdded"`
	StatementsFailed      int `json:"statementsFailed"`
	StatementsSkipped     int `json:"statementsSkipped"`
	EdgesCreatedOrUpdated int `json:"edgesCreatedOrUpdated"`
	EdgesRemoved          int `json:"edgesRemoved"`

	Warnings []string `json:"warnings"`
}

func (s *DataUsageSync) IsClientValid(ctx context.Context, c plugin.PluginClient) (bool, error) {
	dus, err := c.GetDataUsageSyncer()
	if err != nil {
		return false, err
	}

	return version_management.IsValidToSync(ctx, dus, dupc.MinimalCliVersion)
}

func (s *DataUsageSync) GetParts() []job.TaskPart {
	return []job.TaskPart{s}
}

func (s *DataUsageSync) StartSyncAndQueueTaskPart(ctx context.Context, client plugin.PluginClient, statusUpdater job.TaskEventUpdater) (job.JobStatus, string, error) {
	targetFile, err := filepath.Abs(file.CreateUniqueFileNameForTarget(s.TargetConfig.Name, "fromTarget-usage", "json"))

	if err != nil {
		return job.Failed, "", err
	}

	s.TargetConfig.TargetLogger.Debug(fmt.Sprintf("Using %q as data usage target file", targetFile))

	defer s.TargetConfig.HandleTempFile(targetFile, false)

	syncerConfig := dupc.DataUsageSyncConfig{
		ConfigMap:       &baseconfig.ConfigMap{Parameters: s.TargetConfig.Parameters},
		TargetFile:      targetFile,
		MaxBytesPerFile: s.GetMaxBytesPerFile(),
	}

	dus, err := client.GetDataUsageSyncer()
	if err != nil {
		return job.Failed, "", err
	}

	importerConfig := DataUsageImportConfig{
		BaseTargetConfig: *s.TargetConfig,
		TargetFile:       targetFile,
	}
	duImporter := NewDataUsageImporter(&importerConfig, statusUpdater)

	s.TargetConfig.TargetLogger.Info("Fetching last synchronization date")

	firstUsed, lastUsed, err := duImporter.GetLastAndFirstUsage()

	if err != nil {
		hclog.L().Warn(fmt.Sprintf("error retrieving first/last usage for data source %s, last used: %s, first used: %s", importerConfig.DataSourceId, lastUsed, firstUsed))
		return job.Failed, "", err
	}

	if lastUsed != nil {
		lastUsedValue := *lastUsed
		syncerConfig.ConfigMap.Parameters["lastUsed"] = lastUsedValue.Format(time.RFC3339)
	}

	if firstUsed != nil {
		firstUsedValue := *firstUsed
		syncerConfig.ConfigMap.Parameters["firstUsed"] = firstUsedValue.Format(time.RFC3339)
	}

	s.TargetConfig.TargetLogger.Info(fmt.Sprintf("Fetching usage data from the data source, using first used %v and last used %v", syncerConfig.ConfigMap.Parameters["firstUsed"], syncerConfig.ConfigMap.Parameters["lastUsed"]))

	res, err := dus.SyncDataUsage(context.Background(), &syncerConfig)
	if err != nil {
		return job.Failed, "", err
	} else if res.Error != nil { //nolint:staticcheck
		return job.Failed, "", errors.New(res.Error.ErrorMessage) //nolint:staticcheck
	}

	s.TargetConfig.TargetLogger.Info("Importing usage data into Raito")

	var filesCreated []string
	if len(res.TargetFiles) > 0 {
		filesCreated = res.TargetFiles
	} else {
		filesCreated = []string{targetFile}
	}

	status, subtaskId, err := duImporter.TriggerImport(ctx, s.JobId, filesCreated)
	if err != nil {
		return job.Failed, "", err
	}

	if status == job.Queued {
		s.TargetConfig.TargetLogger.Info("Successfully queued import job. Wait until remote processing is done.")
	}

	s.TargetConfig.TargetLogger.Debug(fmt.Sprintf("Current status: %s", status.String()))

	return status, subtaskId, nil
}

func (s *DataUsageSync) ProcessResults(results interface{}) error {
	if duResult, ok := results.(*DataUsageImportResult); ok {
		if duResult != nil && len(duResult.Warnings) > 0 {
			s.TargetConfig.TargetLogger.Info(fmt.Sprintf("Synced data usage with %d warnings (see below). %d statements added, %d failed", len(duResult.Warnings), duResult.StatementsAdded, duResult.StatementsFailed))

			for _, warning := range duResult.Warnings {
				s.TargetConfig.TargetLogger.Warn(warning)
			}
		} else {
			s.TargetConfig.TargetLogger.Info(fmt.Sprintf("Successfully synced data usage. %d statements added, %d failed",
				duResult.StatementsAdded, duResult.StatementsFailed))
		}

		s.result = &job.TaskResult{
			ObjectType: "statements",
			Added:      duResult.StatementsAdded,
			Failed:     duResult.StatementsFailed,
		}

		return nil
	}

	return fmt.Errorf("failed to load results")
}

func (s *DataUsageSync) GetResultObject() interface{} {
	return &DataUsageImportResult{}
}

func (s *DataUsageSync) GetTaskResults() []job.TaskResult {
	if s.result == nil {
		return nil
	}

	return []job.TaskResult{*s.result}
}

func (s *DataUsageSync) GetMaxBytesPerFile() uint64 {
	var v datasize.ByteSize
	if err := v.UnmarshalText([]byte(viper.GetString(constants.MaximumFileSizesFlag))); err != nil {
		s.TargetConfig.TargetLogger.Error(fmt.Sprintf("Error parsing maximum file size: %s. Will use 512MB instead.", err.Error()))
		v = datasize.MB * 512
	}

	return v.Bytes()
}
