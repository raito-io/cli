package data_source

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

type DataSourceImportConfig struct {
	types.BaseTargetConfig
	TargetFile      string
	DeleteUntouched bool
}

type DataSourceImporter interface {
	TriggerImport(ctx context.Context, jobId string) (job.JobStatus, string, error)
}

type dataSourceImporter struct {
	config        *DataSourceImportConfig
	log           hclog.Logger
	statusUpdater job.TaskEventUpdater
}

func NewDataSourceImporter(config *DataSourceImportConfig, statusUpdater job.TaskEventUpdater) DataSourceImporter {
	logger := config.TargetLogger.With("datasource", config.DataSourceId, "file", config.TargetFile)
	dsI := dataSourceImporter{config, logger, statusUpdater}

	return &dsI
}

func (d *dataSourceImporter) TriggerImport(ctx context.Context, jobId string) (job.JobStatus, string, error) {
	if viper.GetBool(constants.SkipFileUpload) {
		// In the development environment, we skip the upload and use the local file for the import
		return d.doImport(jobId, d.config.TargetFile)
	} else {
		key, err := d.upload(ctx)
		if err != nil {
			return job.Failed, "", err
		}

		return d.doImport(jobId, key)
	}
}

func (d *dataSourceImporter) upload(ctx context.Context) (string, error) {
	d.statusUpdater.SetStatusToDataUpload(ctx)

	key, err := file.UploadFile(d.config.TargetFile, &d.config.BaseTargetConfig)
	if err != nil {
		return "", fmt.Errorf("error while uploading data source import files to Raito: %s", err.Error())
	}

	return key, nil
}

func (d *dataSourceImporter) doImport(jobId, fileKey string) (job.JobStatus, string, error) {
	start := time.Now()

	gqlQuery := fmt.Sprintf(`{ "operationName": "ImportDataSourceRequest", "variables":{}, "query": "mutation ImportDataSourceRequest {
      importDataSourceRequest(input: {
        jobId: \"%s\",
          importSettings: {
            dataSource: \"%s\",
            deleteUntouched: %t, 
            fileKey: \"%s\"
          }
        }) {
          subtask {
            status
            subtaskId
          }
        }
    }" }"`, jobId, d.config.DataSourceId, d.config.DeleteUntouched, fileKey)

	gqlQuery = strings.Replace(gqlQuery, "\n", "\\n", -1)

	res := Response{}
	_, err := graphql.ExecuteGraphQL(gqlQuery, &d.config.BaseConfig, &res)

	if err != nil {
		return job.Failed, "", fmt.Errorf("error while executing import: %s", err.Error())
	}

	d.log.Info(fmt.Sprintf("Submitted import in %s", time.Since(start).Round(time.Millisecond)))

	subtask := res.Response.Subtask

	return subtask.Status, subtask.SubtaskId, nil
}

type subtaskResponse struct {
	Status    job.JobStatus `json:"status"`
	SubtaskId string        `json:"subtaskId"`
}

type QueryResponse struct {
	Subtask subtaskResponse `json:"subtask"`
}

type Response struct {
	Response QueryResponse `json:"importDataSourceRequest"`
}
