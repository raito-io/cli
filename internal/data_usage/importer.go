package data_usage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"

	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/file"
	"github.com/raito-io/cli/internal/graphql"
	"github.com/raito-io/cli/internal/job"
	"github.com/raito-io/cli/internal/target/types"
)

type DataUsageImportConfig struct {
	types.BaseTargetConfig
	TargetFile string
}

type DataSourceUsageInfo struct {
	Id        string `json:"id"`
	LastUsed  string `json:"usageLastUsed"`
	FirstUsed string `json:"usageFirstUsed"`
}

type DataUsageImporter interface {
	TriggerImport(ctx context.Context, jobId string, files []string) (job.JobStatus, string, error)
	GetLastAndFirstUsage() (*time.Time, *time.Time, error)
}

type dataUsageImporter struct {
	config        *DataUsageImportConfig
	log           hclog.Logger
	statusUpdater job.TaskEventUpdater
}

func NewDataUsageImporter(config *DataUsageImportConfig, statusUpdater job.TaskEventUpdater) DataUsageImporter {
	logger := config.TargetLogger.With("data-usage", config.DataSourceId, "file", config.TargetFile)
	duI := dataUsageImporter{config, logger, statusUpdater}

	return &duI
}

func (d *dataUsageImporter) TriggerImport(ctx context.Context, jobId string, files []string) (job.JobStatus, string, error) {
	if viper.GetBool(constants.SkipFileUpload) {
		// In the development environment, we skip the upload and use the local file for the import
		return d.doImport(ctx, jobId, files)
	} else {
		keys := make([]string, 0, len(files))

		for _, file := range files {
			key, err := d.upload(ctx, file)
			if err != nil {
				return job.Failed, "", err
			}

			keys = append(keys, key)
		}

		return d.doImport(ctx, jobId, keys)
	}
}

func (d *dataUsageImporter) upload(ctx context.Context, filePath string) (string, error) {
	d.statusUpdater.SetStatusToDataUpload(ctx)
	key, err := file.UploadFile(filePath, &d.config.BaseTargetConfig)

	if err != nil {
		return "", fmt.Errorf("error while uploading data usage import files to Raito: %s", err.Error())
	}

	return key, nil
}

func (d *dataUsageImporter) doImport(ctx context.Context, jobId string, fileKeys []string) (job.JobStatus, string, error) {
	var mutation struct {
		ImportDataUsageRequest struct {
			Subtask struct {
				SubtaskId string
				Status    job.JobStatus
			}
		} `graphql:"importDataUsageRequest(input: $request)"`
	}

	variables := map[string]interface{}{
		"request": DataUsageImportRequest{
			JobId: ID(jobId),
			ImportSettings: DataUsageImportSettings{
				DataSource: ID(d.config.DataSourceId),
				FileKeys:   fileKeys,
			},
		},
	}

	err := graphql.NewClient(&d.config.BaseConfig).Mutate(ctx, &mutation, variables)

	if err != nil {
		return job.Failed, "", fmt.Errorf("error while executing data usage import on appserver: %s", err.Error())
	}

	return mutation.ImportDataUsageRequest.Subtask.Status, mutation.ImportDataUsageRequest.Subtask.SubtaskId, nil
}

func (d *dataUsageImporter) GetLastAndFirstUsage() (*time.Time, *time.Time, error) {
	gqlQuery := fmt.Sprintf(`{"variables":{}, "query": "query {dataSource(id:\"%s\") { ... on DataSource {id usageLastUsed usageFirstUsed }}}" }`, d.config.DataSourceId)
	gqlQuery = strings.ReplaceAll(gqlQuery, "\n", "\\n")
	res := LastUsedResponse{}
	_, err := graphql.ExecuteGraphQL(gqlQuery, &d.config.BaseConfig, &res)

	if err != nil {
		return nil, nil, fmt.Errorf("error while executing data usage import on appserver: %s", err.Error())
	}

	var finalResultFirstUsage, finalResultLastUsage *time.Time

	if res.DataSourceInfo.LastUsed != "" {
		finalResultRaw, err := time.Parse(time.RFC3339, res.DataSourceInfo.LastUsed)
		if err == nil {
			finalResultLastUsage = &finalResultRaw
		}
	}

	if res.DataSourceInfo.FirstUsed != "" {
		finalResultRaw, err := time.Parse(time.RFC3339, res.DataSourceInfo.FirstUsed)
		if err == nil {
			finalResultFirstUsage = &finalResultRaw
		}
	}

	return finalResultFirstUsage, finalResultLastUsage, nil
}

type ID string

type DataUsageImportSettings struct {
	DataSource ID       `json:"dataSource"`
	FileKeys   []string `json:"fileKeys"`
}

type DataUsageImportRequest struct {
	JobId          ID                      `json:"jobId"`
	ImportSettings DataUsageImportSettings `json:"importSettings"`
}

type LastUsedResponse struct {
	DataSourceInfo DataSourceUsageInfo `json:"dataSource"`
}
