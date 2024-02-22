package target_sync

import (
	"context"

	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/plugin"
	"github.com/raito-io/cli/internal/resource_provider"
	"github.com/raito-io/cli/internal/target/types"
)

func resourceProviderSync(ctx context.Context, targetConfig *types.BaseTargetConfig, client plugin.PluginClient, jobId string) error {
	resourceSourceSyncTask := &resource_provider.ResourceSync{
		TargetConfig: targetConfig,
		JobId:        jobId,
	}

	err := execute(ctx, targetConfig.DataSourceId, jobId, constants.ResourceProviderSync, "resource provisioning", targetConfig.SkipResourceProvider, resourceSourceSyncTask, targetConfig, client)
	if err != nil {
		return err
	}

	return nil
}
