package target_sync

import (
	"context"
	"fmt"

	"github.com/raito-io/cli/internal/access_provider"
	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/data_source"
	"github.com/raito-io/cli/internal/data_usage"
	"github.com/raito-io/cli/internal/identity_store"
	"github.com/raito-io/cli/internal/plugin"
	"github.com/raito-io/cli/internal/target/types"
)

func dsSyncTargetSync(ctx context.Context, targetConfig *types.BaseTargetConfig, client plugin.PluginClient, jobId string) error {
	err := dataSourceSync(ctx, targetConfig, jobId, client)
	if err != nil {
		return fmt.Errorf("data source sync: %w", err)
	}

	err = identityStoreSync(ctx, targetConfig, jobId, client)
	if err != nil {
		return fmt.Errorf("identity store sync: %w", err)
	}

	err = dataAccessSync(ctx, targetConfig, jobId, client)
	if err != nil {
		return fmt.Errorf("data access sync: %w", err)
	}

	err = dataUsageSync(ctx, targetConfig, jobId, client)
	if err != nil {
		return fmt.Errorf("data usage sync: %w", err)
	}

	return nil
}

func isSyncTargetSync(ctx context.Context, targetConfig *types.BaseTargetConfig, client plugin.PluginClient, jobId string) error {
	err := identityStoreSync(ctx, targetConfig, jobId, client)
	if err != nil {
		return fmt.Errorf("identity store sync: %w", err)
	}

	return nil
}

func dataUsageSync(ctx context.Context, targetConfig *types.BaseTargetConfig, jobID string, client plugin.PluginClient) error {
	dataUsageSyncTask := &data_usage.DataUsageSync{TargetConfig: targetConfig, JobId: jobID}

	err := execute(ctx, targetConfig.DataSourceId, jobID, constants.DataUsageSync, "data usage", targetConfig.SkipDataUsageSync, dataUsageSyncTask, targetConfig, client)
	if err != nil {
		return err
	}

	return nil
}

func dataAccessSync(ctx context.Context, targetConfig *types.BaseTargetConfig, jobID string, client plugin.PluginClient) error {
	dataAccessSyncTask := &access_provider.DataAccessSync{TargetConfig: targetConfig, JobId: jobID}

	err := execute(ctx, targetConfig.DataSourceId, jobID, constants.DataAccessSync, "data access", targetConfig.SkipDataAccessSync, dataAccessSyncTask, targetConfig, client)
	if err != nil {
		return err
	}

	return nil
}

func identityStoreSync(ctx context.Context, targetConfig *types.BaseTargetConfig, jobID string, client plugin.PluginClient) error {
	identityStoreSyncTask := &identity_store.IdentityStoreSync{TargetConfig: targetConfig, JobId: jobID}

	err := execute(ctx, targetConfig.IdentityStoreId, jobID, constants.IdentitySync, "identity store", targetConfig.SkipIdentityStoreSync, identityStoreSyncTask, targetConfig, client)
	if err != nil {
		return err
	}

	return nil
}

func dataSourceSync(ctx context.Context, targetConfig *types.BaseTargetConfig, jobID string, client plugin.PluginClient) error {
	dataSourceSyncTask := &data_source.DataSourceSync{TargetConfig: targetConfig, JobId: jobID}

	err := execute(ctx, targetConfig.DataSourceId, jobID, constants.DataSourceSync, "data source metadata", targetConfig.SkipDataSourceSync, dataSourceSyncTask, targetConfig, client)
	if err != nil {
		return err
	}

	return nil
}
