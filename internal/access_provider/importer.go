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
	TriggerImport(ctx context.Context, jobId string) (job.JobStatus, string, error)
}

type accessProviderImporter struct {
	config        *AccessProviderImportConfig
	log           hclog.Logger
	statusUpdater job.TaskEventUpdater
}

func NewAccessProviderImporter(config *AccessProviderImportConfig, statusUpdater job.TaskEventUpdater) AccessProviderImporter {
	logger := config.TargetLogger.With("AccessProvider", config.DataSourceId, "file", config.TargetFile)
	dsI := accessProviderImporter{config, logger, statusUpdater}

	return &dsI
}

func (d *accessProviderImporter) TriggerImport(ctx context.Context, jobId string) (job.JobStatus, string, error) {
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

func (d *accessProviderImporter) upload(ctx context.Context) (string, error) {
	d.statusUpdater.SetStatusToDataUpload(ctx)

	key, err := file.UploadFile(d.config.TargetFile, &d.config.BaseTargetConfig)
	if err != nil {
		return "", fmt.Errorf("error while uploading data source import files to Raito: %s", err.Error())
	}

	return key, nil
}

func (d *accessProviderImporter) doImport(jobId string, fileKey string) (job.JobStatus, string, error) {
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

	gqlQuery = strings.ReplaceAll(gqlQuery, "\n", "\\n")

	res := ImportResponse{}
	_, err := graphql.ExecuteGraphQL(gqlQuery, &d.config.BaseConfig, &res)

	if err != nil {
		return job.Failed, "", fmt.Errorf("error while executing import: %s", err.Error())
	}

	retStatus := res.Response.Subtask.Status
	subtaskId := res.Response.Subtask.SubtaskId

	d.log.Info(fmt.Sprintf("Done submitting import in %s", time.Since(start).Round(time.Millisecond)))

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
