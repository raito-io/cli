package tag

import (
	"context"
	"fmt"
	"path/filepath"

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

func (s *TagSync) StartSyncAndQueueTaskPart(ctx context.Context, client plugin.PluginClient, statusUpdater job.TaskEventUpdater) (job.JobStatus, string, error) {
	targetFile, err := filepath.Abs(file.CreateUniqueFileNameForTarget(s.TargetConfig.Name, "fromTarget-tags", "json"))
	if err != nil {
		return job.Failed, "", err
	}

	s.TargetConfig.TargetLogger.Debug(fmt.Sprintf("Using %q as tag taret file", targetFile))

	defer s.TargetConfig.HandleTempFile(targetFile, false)

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

	s.TargetConfig.TargetLogger.Info(fmt.Sprintf("Successfully synced %d tags.", res.Tags))

	importConfig := TagImportConfig{
		BaseTargetConfig: *s.TargetConfig,
		TargetFile:       targetFile,
		TagSourcesScope:  res.TagSourcesScope,
	}

	tagImporter := NewTagImporter(&importConfig, statusUpdater)

	s.TargetConfig.TargetLogger.Info("Importing tags into Raito")

	status, subtaskId, err := tagImporter.TriggerImport(ctx, s.JobId)
	if err != nil {
		return job.Failed, "", fmt.Errorf("import tags into Raito: %w", err)
	}

	if status == job.Queued {
		s.TargetConfig.TargetLogger.Info("Successfully queued import job. Wait until remote processing is done.")
	}

	s.TargetConfig.TargetLogger.Debug(fmt.Sprintf("Current status: %s", status.String()))

	return status, subtaskId, nil
}

func (s *TagSync) ProcessResults(results interface{}) error {
	if tagResult, ok := results.(*TagImportResult); ok {
		numberOfTags := tagResult.TagsAdded + tagResult.TagsUpdated + tagResult.TagsRemoved

		if len(tagResult.Warnings) > 0 {
			s.TargetConfig.TargetLogger.Info(fmt.Sprintf("Synced %d tags with %d warnings (see below).", numberOfTags, len(tagResult.Warnings)))

			for _, warning := range tagResult.Warnings {
				s.TargetConfig.TargetLogger.Warn(warning)
			}
		} else {
			s.TargetConfig.TargetLogger.Info(fmt.Sprintf("Successfully synced %d tags.", numberOfTags))
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
