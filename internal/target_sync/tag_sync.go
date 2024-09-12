package target_sync

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/error_handler"
	"github.com/raito-io/cli/internal/plugin"
	"github.com/raito-io/cli/internal/tag"
	"github.com/raito-io/cli/internal/target/types"
)

func tagSyncTargetSync(ctx context.Context, logger hclog.Logger, targetConfig *types.BaseTargetConfig, client plugin.PluginClient, jobId string, s *SyncJob, eh error_handler.ErrorHandler) {
	tagSyncTask := &tag.TagSync{
		TargetConfig: targetConfig,
		JobId:        jobId,
	}

	s.execute(ctx, logger, targetConfig.DataSourceId, jobId, constants.TagSync, "tags", targetConfig.SkipTagSync, tagSyncTask, targetConfig, client, eh)
}
