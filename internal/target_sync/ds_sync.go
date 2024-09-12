package target_sync

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/raito-io/cli/internal/access_provider"
	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/data_source"
	"github.com/raito-io/cli/internal/data_usage"
	"github.com/raito-io/cli/internal/error_handler"
	"github.com/raito-io/cli/internal/identity_store"
	"github.com/raito-io/cli/internal/plugin"
	"github.com/raito-io/cli/internal/target/types"
)

func dsSyncTargetSync(ctx context.Context, logger hclog.Logger, targetConfig *types.BaseTargetConfig, client plugin.PluginClient, jobId string, s *SyncJob, eh error_handler.ErrorHandler) {
	dataSourceSync(ctx, logger, targetConfig, jobId, client, s, error_handler.Wrap(eh, "data source sync: %w", error_handler.ErrorPlaceholder))
	if eh.HasError() {
		return
	}

	identityStoreSync(ctx, logger, targetConfig, jobId, client, s, error_handler.Wrap(eh, "identity store sync: %w", error_handler.ErrorPlaceholder))
	if eh.HasError() {
		return
	}

	dataAccessSync(ctx, logger, targetConfig, jobId, client, s, error_handler.Wrap(eh, "data access sync: %w", error_handler.ErrorPlaceholder))
	if eh.HasError() {
		return
	}

	dataUsageSync(ctx, logger, targetConfig, jobId, client, s, error_handler.Wrap(eh, "data usage sync: %w", error_handler.ErrorPlaceholder))
	if eh.HasError() {
		return
	}

	return
}

func isSyncTargetSync(ctx context.Context, logger hclog.Logger, targetConfig *types.BaseTargetConfig, client plugin.PluginClient, jobId string, s *SyncJob, eh error_handler.ErrorHandler) {
	identityStoreSync(ctx, logger, targetConfig, jobId, client, s, error_handler.Wrap(eh, "identity store sync: %w", error_handler.ErrorPlaceholder))
}

func dataUsageSync(ctx context.Context, logger hclog.Logger, targetConfig *types.BaseTargetConfig, jobID string, client plugin.PluginClient, s *SyncJob, eh error_handler.ErrorHandler) {
	dataUsageSyncTask := &data_usage.DataUsageSync{TargetConfig: targetConfig, JobId: jobID}

	s.execute(ctx, logger, targetConfig.DataSourceId, jobID, constants.DataUsageSync, "data usage", targetConfig.SkipDataUsageSync, dataUsageSyncTask, targetConfig, client, eh)
}

func dataAccessSync(ctx context.Context, logger hclog.Logger, targetConfig *types.BaseTargetConfig, jobID string, client plugin.PluginClient, s *SyncJob, eh error_handler.ErrorHandler) {
	dataAccessSyncTask := &access_provider.DataAccessSync{TargetConfig: targetConfig, JobId: jobID}

	s.execute(ctx, logger, targetConfig.DataSourceId, jobID, constants.DataAccessSync, "data access", targetConfig.SkipDataAccessSync, dataAccessSyncTask, targetConfig, client, eh)
}

func identityStoreSync(ctx context.Context, logger hclog.Logger, targetConfig *types.BaseTargetConfig, jobID string, client plugin.PluginClient, s *SyncJob, eh error_handler.ErrorHandler) {
	identityStoreSyncTask := &identity_store.IdentityStoreSync{TargetConfig: targetConfig, JobId: jobID}

	s.execute(ctx, logger, targetConfig.IdentityStoreId, jobID, constants.IdentitySync, "identity store", targetConfig.SkipIdentityStoreSync, identityStoreSyncTask, targetConfig, client, eh)
}

func dataSourceSync(ctx context.Context, logger hclog.Logger, targetConfig *types.BaseTargetConfig, jobID string, client plugin.PluginClient, s *SyncJob, eh error_handler.ErrorHandler) {
	dataSourceSyncTask := &data_source.DataSourceSync{TargetConfig: targetConfig, JobId: jobID}

	s.execute(ctx, logger, targetConfig.DataSourceId, jobID, constants.DataSourceSync, "data source metadata", targetConfig.SkipDataSourceSync, dataSourceSyncTask, targetConfig, client, eh)
}
