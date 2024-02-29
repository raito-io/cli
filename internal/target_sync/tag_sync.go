package target_sync

import (
	"context"

	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/plugin"
	"github.com/raito-io/cli/internal/tag"
	"github.com/raito-io/cli/internal/target/types"
)

func tagSyncTargetSync(ctx context.Context, targetConfig *types.BaseTargetConfig, client plugin.PluginClient, jobId string) error {
	tagSyncTask := &tag.TagSync{
		TargetConfig: targetConfig,
		JobId:        jobId,
	}

	err := execute(ctx, targetConfig.DataSourceId, jobId, constants.TagSync, "tags", targetConfig.SkipTagsSync, tagSyncTask, targetConfig, client)
	if err != nil {
		return err
	}

	return nil
}
