package tag

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/hashicorp/go-hclog"

	"github.com/raito-io/cli/base/tag"
	baseconfig "github.com/raito-io/cli/base/util/config"
	"github.com/raito-io/cli/internal/job"
	"github.com/raito-io/cli/internal/plugin"
	"github.com/raito-io/cli/internal/target/types"
	"github.com/raito-io/cli/internal/util/file"
	"github.com/raito-io/cli/internal/version_management"
)

type TagImportResult struct {
	TagsAdded   int `json:"tagsAdded"`
	TagsUpdated int `json:"tagsUpdated"`
	TagsRemoved int `json:"tagsRemoved"`

	Warnings []string `json:"warnings"`
}

type TagSync struct {
	TargetConfig *types.BaseTargetConfig
	JobId        string

	result *job.TaskResult
}

func (s *TagSync) IsClientValid(ctx context.Context, c plugin.PluginClient) (bool, error) {
	ts, err := c.GetTagSyncer()
	if err != nil {
		return false, err
	}

	return version_management.IsValidToSync(ctx, ts, tag.MinimalCliVersion)
}

func (s *TagSync) GetParts() []job.TaskPart {
	return []job.TaskPart{s}
}

func (s *TagSync) StartSyncAndQueueTaskPart(ctx context.Context, logger hclog.Logger, client plugin.PluginClient, statusUpdater job.TaskEventUpdater, secureImport func(func() error) error) (job.JobStatus, string, error) {
	targetFile, err := filepath.Abs(file.CreateUniqueFileNameForTarget(s.TargetConfig.Name, "fromTarget-tags", "json"))
	if err != nil {
		return job.Failed, "", err
	}

	logger.Debug(fmt.Sprintf("Using %q as tag taret file", targetFile))

	defer s.TargetConfig.HandleTempFile(logger, targetFile, false)

	syncerConfig := tag.TagSyncConfig{
		ConfigMap:       &baseconfig.ConfigMap{Parameters: s.TargetConfig.Parameters},
		TargetFile:      targetFile,
		DataSourceId:    s.TargetConfig.DataSourceId,
		IdentityStoreId: s.TargetConfig.IdentityStoreId,
	}

	ts, err := client.GetTagSyncer()
	if err != nil {
		return job.Failed, "", fmt.Errorf("fetching tag syncer: %w", err)
	}

	res, err := ts.SyncTags(ctx, &syncerConfig)
	if err != nil {
		return job.Failed, "", fmt.Errorf("syncing tags: %w", err)
	}

	logger.Info(fmt.Sprintf("Successfully synced %d tags.", res.Tags))

	importConfig := TagImportConfig{
		BaseTargetConfig: *s.TargetConfig,
		TargetFile:       targetFile,
		TagSourcesScope:  res.TagSourcesScope,
	}

	tagImporter := NewTagImporter(&importConfig, statusUpdater)

	var status job.JobStatus
	var subtaskId string

	err = secureImport(func() error {
		logger.Info("Importing tags into Raito")

		status, subtaskId, err = tagImporter.TriggerImport(ctx, logger, s.JobId)
		if err != nil {
			err = fmt.Errorf("import tags into Raito: %w", err)

			return err
		}

		if status == job.Queued {
			logger.Info("Successfully queued import job. Wait until remote processing is done.")
		}

		logger.Debug(fmt.Sprintf("Current status: %s", status.String()))

		return nil
	})
	if err != nil {
		return job.Failed, "", err
	}

	return status, subtaskId, nil
}

func (s *TagSync) ProcessResults(logger hclog.Logger, results interface{}) error {
	if tagResult, ok := results.(*TagImportResult); ok {
		numberOfTags := tagResult.TagsAdded + tagResult.TagsUpdated + tagResult.TagsRemoved

		if len(tagResult.Warnings) > 0 {
			logger.Info(fmt.Sprintf("Synced %d tags with %d warnings (see below).", numberOfTags, len(tagResult.Warnings)))

			for _, warning := range tagResult.Warnings {
				logger.Warn(warning)
			}
		} else {
			logger.Info(fmt.Sprintf("Successfully synced %d tags.", numberOfTags))
		}

		s.result = &job.TaskResult{
			ObjectType: "tags",
			Added:      tagResult.TagsAdded,
			Updated:    tagResult.TagsUpdated,
			Removed:    tagResult.TagsRemoved,
			Failed:     len(tagResult.Warnings),
		}

		return nil
	}

	return fmt.Errorf("failed to load results")
}

func (s *TagSync) GetResultObject() interface{} {
	return &TagImportResult{}
}

func (s *TagSync) GetTaskResults() []job.TaskResult {
	if s.result == nil {
		return nil
	}

	return []job.TaskResult{*s.result}
}
