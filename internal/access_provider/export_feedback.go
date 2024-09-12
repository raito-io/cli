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

type AccessProviderExportFeedbackConfig struct {
	types.BaseTargetConfig
	FeedbackFile string
}

type AccessProviderExportFeedbackSync interface {
	TriggerFeedbackImport(ctx context.Context, logger hclog.Logger, jobId string) (job.JobStatus, string, error)
}

type accessProviderFeedbackSync struct {
	config        *AccessProviderExportFeedbackConfig
	statusUpdater job.TaskEventUpdater
}

func NewAccessProviderFeedbackImporter(config *AccessProviderExportFeedbackConfig, statusUpdater job.TaskEventUpdater) AccessProviderExportFeedbackSync {
	apI := accessProviderFeedbackSync{config, statusUpdater}

	return &apI
}

func (i *accessProviderFeedbackSync) TriggerFeedbackImport(ctx context.Context, logger hclog.Logger, jobId string) (job.JobStatus, string, error) {
	logger = logger.With("AccessProvider", i.config.DataSourceId, "file", i.config.FeedbackFile)

	if viper.GetBool(constants.SkipFileUpload) {
		// In the development environment, we skip the upload and use the local file for the import
		return i.doImport(logger, jobId, i.config.FeedbackFile)
	} else {
		key, err := i.upload(ctx, logger)
		if err != nil {
			return job.Failed, "", err
		}

		return i.doImport(logger, jobId, key)
	}
}

func (i *accessProviderFeedbackSync) upload(ctx context.Context, logger hclog.Logger) (string, error) {
	i.statusUpdater.SetStatusToDataUpload(ctx)

	key, err := file.UploadFile(logger, i.config.FeedbackFile, &i.config.BaseTargetConfig)
	if err != nil {
		return "", fmt.Errorf("error while uploading access provider feedback import files to Raito: %s", err.Error())
	}

	return key, nil
}

func (i *accessProviderFeedbackSync) doImport(logger hclog.Logger, jobId string, fileKey string) (job.JobStatus, string, error) {
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

	logger.Info(fmt.Sprintf("Done submitting feedback import in %s", time.Since(start).Round(time.Millisecond)))

	return retStatus, subtaskId, nil
}

type FeedbackResponse struct {
	Response QueryResponse `json:"importAccessProvidersSyncFeedback"`
}
