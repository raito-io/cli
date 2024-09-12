package tag

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"

	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/file"
	"github.com/raito-io/cli/internal/graphql"
	"github.com/raito-io/cli/internal/job"
	"github.com/raito-io/cli/internal/target/types"
)

type TagImportConfig struct {
	types.BaseTargetConfig
	TargetFile      string
	TagSourcesScope []string `json:"tagSourcesScope"`
}

type TagImporter interface {
	TriggerImport(ctx context.Context, logger hclog.Logger, jobId string) (job.JobStatus, string, error)
}

type tagImporter struct {
	config        *TagImportConfig
	statusUpdater job.TaskEventUpdater
}

func NewTagImporter(config *TagImportConfig, statusUpdater job.TaskEventUpdater) TagImporter {
	tagI := tagImporter{config, statusUpdater}

	return tagI
}

func (t tagImporter) TriggerImport(ctx context.Context, logger hclog.Logger, jobId string) (job.JobStatus, string, error) {
	logger = logger.With("datasource", t.config.DataSourceId, "file", t.config.TargetFile)

	if viper.GetBool(constants.SkipFileUpload) {
		return t.doImport(ctx, logger, jobId, t.config.TargetFile)
	} else {
		key, err := t.upload(ctx, logger)
		if err != nil {
			return job.Failed, "", err
		}

		return t.doImport(ctx, logger, jobId, key)
	}
}

func (t tagImporter) upload(ctx context.Context, logger hclog.Logger) (string, error) {
	t.statusUpdater.SetStatusToDataUpload(ctx)

	key, err := file.UploadFile(logger, t.config.TargetFile, &t.config.BaseTargetConfig)
	if err != nil {
		return "", fmt.Errorf("uploading tag import files to Raito: %w", err)
	}

	return key, nil
}

func (t tagImporter) doImport(ctx context.Context, logger hclog.Logger, jobId, fileKey string) (status job.JobStatus, subtaskId string, err error) {
	start := time.Now()

	defer func() {
		if err != nil {
			logger.Error(fmt.Sprintf("Error while importing tags: %s", err.Error()))
		} else {
			logger.Info(fmt.Sprintf("Imported tags in %s", time.Since(start).Round(time.Millisecond)))
		}
	}()

	type TagImportRequest struct {
		JobId          string `json:"jobId"`
		ImportSettings struct {
			DataSource      *string  `json:"dataSource,omitempty"`
			IdentitySource  *string  `json:"identitySource,omitempty"`
			FileKey         string   `json:"fileKey"`
			TagSourcesScope []string `json:"tagSourcesScope"`
		} `json:"importSettings"`
	}

	variables := TagImportRequest{
		JobId: jobId,
		ImportSettings: struct {
			DataSource      *string  `json:"dataSource,omitempty"`
			IdentitySource  *string  `json:"identitySource,omitempty"`
			FileKey         string   `json:"fileKey"`
			TagSourcesScope []string `json:"tagSourcesScope"`
		}{
			DataSource:      &t.config.DataSourceId,
			IdentitySource:  &t.config.IdentityStoreId,
			FileKey:         fileKey,
			TagSourcesScope: t.config.TagSourcesScope,
		},
	}

	gqlClient := graphql.NewClient(&t.config.BaseConfig)

	var importTagRequestMutation struct {
		ImportTagsRequest struct {
			Subtask struct {
				Status    job.JobStatus `graphql:"status"`
				SubtaskId string        `graphql:"subtaskId"`
			}
		} `graphql:"importTagsRequest(input: $input)"`
	}

	err = gqlClient.Mutate(ctx, &importTagRequestMutation, map[string]interface{}{"input": variables})
	if err != nil {
		return job.Failed, "", fmt.Errorf("executing import: %w", err)
	}

	subtask := importTagRequestMutation.ImportTagsRequest.Subtask

	return subtask.Status, subtask.SubtaskId, nil
}
