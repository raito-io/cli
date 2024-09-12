package data_source

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"

	"github.com/raito-io/cli/internal/util/tag"

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

	// TagSourcesScope is the set of sources that will be looked at when merging the tags. Tags with other sources will remain untouched.
	// If not specified, the default is to take all sources for which tags are defined in the import file.
	TagSourcesScope []string `json:"tagSourcesScope"`
}

type DataSourceImporter interface {
	TriggerImport(ctx context.Context, logger hclog.Logger, jobId string) (job.JobStatus, string, error)
}

type dataSourceImporter struct {
	config        *DataSourceImportConfig
	statusUpdater job.TaskEventUpdater
}

func NewDataSourceImporter(config *DataSourceImportConfig, statusUpdater job.TaskEventUpdater) DataSourceImporter {
	dsI := dataSourceImporter{config, statusUpdater}

	return &dsI
}

func (d *dataSourceImporter) TriggerImport(ctx context.Context, logger hclog.Logger, jobId string) (job.JobStatus, string, error) {
	if viper.GetBool(constants.SkipFileUpload) {
		// In the development environment, we skip the upload and use the local file for the import
		return d.doImport(logger, jobId, d.config.TargetFile)
	} else {
		key, err := d.upload(ctx, logger)
		if err != nil {
			return job.Failed, "", err
		}

		return d.doImport(logger, jobId, key)
	}
}

func (d *dataSourceImporter) upload(ctx context.Context, logger hclog.Logger) (string, error) {
	d.statusUpdater.SetStatusToDataUpload(ctx)

	key, err := file.UploadFile(logger, d.config.TargetFile, &d.config.BaseTargetConfig)
	if err != nil {
		return "", fmt.Errorf("error while uploading data source import files to Raito: %s", err.Error())
	}

	return key, nil
}

func (d *dataSourceImporter) doImport(logger hclog.Logger, jobId, fileKey string) (job.JobStatus, string, error) {
	start := time.Now()

	gqlQuery := fmt.Sprintf(`{ "operationName": "ImportDataSourceRequest", "variables":{}, "query": "mutation ImportDataSourceRequest {
      importDataSourceRequest(input: {
        jobId: \"%s\",
          importSettings: {
            dataSource: \"%s\",
            deleteUntouched: %t,
            fileKey: \"%s\",
            tagSourcesScope: %s
          }
        }) {
          subtask {
            status
            subtaskId
          }
        }
    }" }"`, jobId, d.config.DataSourceId, d.config.DeleteUntouched, fileKey, strings.Replace(tag.SerializeTagList(d.config.TagSourcesScope), "\"", "\\\"", -1))

	gqlQuery = strings.Replace(gqlQuery, "\n", "\\n", -1)

	res := Response{}
	_, err := graphql.ExecuteGraphQL(gqlQuery, &d.config.BaseConfig, &res)

	if err != nil {
		return job.Failed, "", fmt.Errorf("error while executing import: %s", err.Error())
	}

	logger.Info(fmt.Sprintf("Submitted import in %s", time.Since(start).Round(time.Millisecond)))

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
