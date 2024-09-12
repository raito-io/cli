package access_provider

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

type AccessProviderImportConfig struct {
	types.BaseTargetConfig
	TargetFile      string
	DeleteUntouched bool
}

type AccessProviderImporter interface {
	TriggerImport(ctx context.Context, logger hclog.Logger, jobId string) (job.JobStatus, string, error)
}

type accessProviderImporter struct {
	config        *AccessProviderImportConfig
	statusUpdater job.TaskEventUpdater
}

func NewAccessProviderImporter(config *AccessProviderImportConfig, statusUpdater job.TaskEventUpdater) AccessProviderImporter {
	dsI := accessProviderImporter{config, statusUpdater}

	return &dsI
}

func (d *accessProviderImporter) TriggerImport(ctx context.Context, logger hclog.Logger, jobId string) (job.JobStatus, string, error) {
	logger = logger.With("AccessProvider", d.config.DataSourceId, "file", d.config.TargetFile)

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

func (d *accessProviderImporter) upload(ctx context.Context, logger hclog.Logger) (string, error) {
	d.statusUpdater.SetStatusToDataUpload(ctx)

	key, err := file.UploadFile(logger, d.config.TargetFile, &d.config.BaseTargetConfig)
	if err != nil {
		return "", fmt.Errorf("error while uploading data source import files to Raito: %s", err.Error())
	}

	return key, nil
}

func (d *accessProviderImporter) doImport(logger hclog.Logger, jobId string, fileKey string) (job.JobStatus, string, error) {
	start := time.Now()

	gqlQuery := fmt.Sprintf(`{ "operationName": "ImportAccessProvidersRequest", "variables":{}, "query": "mutation ImportAccessProvidersRequest {
        importAccessProvidersRequest(input: {
          jobId: \"%s\",
          importSettings: {
            dataSource: \"%s\",
            deleteUntouched: %t,
            fileKey: \"%s\"
          }
        }) {
          subtask {
            subtaskId
            status            
          }
         }
    }" }"`, jobId, d.config.DataSourceId, d.config.DeleteUntouched, fileKey)

	gqlQuery = strings.Replace(gqlQuery, "\n", "\\n", -1)

	res := ImportResponse{}
	_, err := graphql.ExecuteGraphQL(gqlQuery, &d.config.BaseConfig, &res)

	if err != nil {
		return job.Failed, "", fmt.Errorf("error while executing import: %s", err.Error())
	}

	retStatus := res.Response.Subtask.Status
	subtaskId := res.Response.Subtask.SubtaskId

	logger.Info(fmt.Sprintf("Done submitting import in %s", time.Since(start).Round(time.Millisecond)))

	return retStatus, subtaskId, nil
}

type subtaskResponse struct {
	Status    job.JobStatus `json:"status"`
	SubtaskId string        `json:"subtaskId"`
}

type QueryResponse struct {
	Subtask subtaskResponse `json:"subtask"`
}

type ImportResponse struct {
	Response QueryResponse `json:"importAccessProvidersRequest"`
}
