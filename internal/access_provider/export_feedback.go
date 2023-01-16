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

type AccessProviderExportFeedbackConfig struct {
	target.BaseTargetConfig
	FeedbackFile string
}

type AccessProviderExportFeedbackSync interface {
	TriggerFeedbackImport(jobId string) (job.JobStatus, string, error)
}

type accessProviderFeedbackSync struct {
	config        *AccessProviderExportFeedbackConfig
	log           hclog.Logger
	statusUpdater job.TaskEventUpdater
}

func NewAccessProviderFeedbackImporter(config *AccessProviderExportFeedbackConfig, statusUpdater job.TaskEventUpdater) AccessProviderExportFeedbackSync {
	logger := config.TargetLogger.With("AccessProvider", config.DataSourceId, "file", config.FeedbackFile)
	apI := accessProviderFeedbackSync{config, logger, statusUpdater}

	return &apI
}

func (i *accessProviderFeedbackSync) TriggerFeedbackImport(jobId string) (job.JobStatus, string, error) {
	env := viper.GetString(constants.EnvironmentFlag)
	if env == constants.EnvironmentDev {
		// In the development environment, we skip the upload and use the local file for the import
		return i.doImport(jobId, i.config.FeedbackFile)
	} else {
		key, err := i.upload()
		if err != nil {
			return job.Failed, "", err
		}

		return i.doImport(jobId, key)
	}
}

func (i *accessProviderFeedbackSync) upload() (string, error) {
	i.statusUpdater.AddTaskEvent(job.DataUpload)

	key, err := file.UploadFile(i.config.FeedbackFile, &i.config.BaseTargetConfig)
	if err != nil {
		return "", fmt.Errorf("error while uploading access provider feedback import files to Raito: %s", err.Error())
	}

	return key, nil
}

func (i *accessProviderFeedbackSync) doImport(jobId string, fileKey string) (job.JobStatus, string, error) {
	start := time.Now()

	gqlQuery := fmt.Sprintf(`{ "operationName": "ImportAccessProvidersSyncFeedback", "variables":{}, "query": "mutation ImportAccessProvidersSyncFeedback {
        importAccessProvidersSyncFeedback(input: {
          jobId: \"%s\",
          importSettings: {
            dataSource: \"%s\",
            fileKey: \"%s\"
          }
        }) {
          subtask {
            subtaskId
            status            
          }
         }
    }" }"
	`, jobId, i.config.DataSourceId, fileKey)

	gqlQuery = strings.Replace(gqlQuery, "\n", "\\n", -1)

	res := FeedbackResponse{}
	_, err := graphql.ExecuteGraphQL(gqlQuery, &i.config.BaseConfig, &res)

	if err != nil {
		return job.Failed, "", fmt.Errorf("error while executing feedback import: %s", err.Error())
	}

	retStatus := res.Response.Subtask.Status
	subtaskId := res.Response.Subtask.SubtaskId

	i.log.Info(fmt.Sprintf("Done submitting feedback import in %s", time.Since(start).Round(time.Millisecond)))

	return retStatus, subtaskId, nil
}

type FeedbackResponse struct {
	Response QueryResponse `json:"importAccessProvidersSyncFeedback"`
}
