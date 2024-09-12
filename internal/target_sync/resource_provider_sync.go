package target_sync

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/error_handler"
	"github.com/raito-io/cli/internal/plugin"
	"github.com/raito-io/cli/internal/resource_provider"
	"github.com/raito-io/cli/internal/target/types"
)

func resourceProviderSync(ctx context.Context, logger hclog.Logger, targetConfig *types.BaseTargetConfig, client plugin.PluginClient, jobId string, s *SyncJob, eh error_handler.ErrorHandler) {
	resourceSourceSyncTask := &resource_provider.ResourceSync{
		TargetConfig: targetConfig,
		JobId:        jobId,
	}

	s.execute(ctx, logger, targetConfig.DataSourceId, jobId, constants.ResourceProviderSync, "resource provisioning", targetConfig.SkipResourceProvider, resourceSourceSyncTask, targetConfig, client, eh)
}
