package access_provider

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"

	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/file"
	"github.com/raito-io/cli/internal/graphql"
	"github.com/raito-io/cli/internal/job"
	"github.com/raito-io/cli/internal/target"
)

type AccessProviderImportConfig struct {
	target.BaseTargetConfig
	TargetFile      string
	DeleteUntouched bool
}

type AccessProviderImporter interface {
	TriggerImport(jobId string) (job.JobStatus, string, error)
}

type accessProviderImporter struct {
	config        *AccessProviderImportConfig
	log           hclog.Logger
	statusUpdater *job.TaskEventUpdater
}

func NewAccessProviderImporter(config *AccessProviderImportConfig, statusUpdater *job.TaskEventUpdater) AccessProviderImporter {
	logger := config.Logger.With("AccessProvider", config.DataSourceId, "file", config.TargetFile)
	dsI := accessProviderImporter{config, logger, statusUpdater}

	return &dsI
}

func (d *accessProviderImporter) TriggerImport(jobId string) (job.JobStatus, string, error) {
	env := viper.GetString(constants.EnvironmentFlag)
	if env == constants.EnvironmentDev {
		// In the development environment, we skip the upload and use the local file for the import
		return d.doImport(jobId, d.config.TargetFile)
	} else {
		key, err := d.upload()
		if err != nil {
			return job.Failed, "", err
		}

		return d.doImport(jobId, key)
	}
}

func (d *accessProviderImporter) upload() (string, error) {
	d.statusUpdater.AddTaskEvent(job.DataUpload)

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

	gqlQuery = strings.Replace(gqlQuery, "\n", "\\n", -1)

	res := Response{}
	_, err := graphql.ExecuteGraphQL(gqlQuery, &d.config.BaseTargetConfig, &res)

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

type Response struct {
	Response QueryResponse `json:"importAccessProvidersRequest"`
}
