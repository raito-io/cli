package wrappers

import (
	"context"
	"fmt"
	"time"

	"github.com/raito-io/cli/base/resource_provider"
	"github.com/raito-io/cli/base/util/config"
)

type ResourceProviderSyncer interface {
	UpdateResources(ctx context.Context, config *resource_provider.UpdateResourceInput) (*resource_provider.UpdateResourceResult, error)
}

type ResourceSyncFactoryFn func(ctx context.Context, configParams *config.ConfigMap) (ResourceProviderSyncer, func(), error)

type ResourceProvisionSyncFunction struct {
	resource_provider.ResourceProviderSyncerVersionHandler

	syncer SyncFactory[ResourceProviderSyncer]
}

func ResourceProviderSyncFactory(syncer ResourceSyncFactoryFn) resource_provider.ResourceProviderSyncer {
	return &ResourceProvisionSyncFunction{
		syncer: NewSyncFactory(syncer),
	}
}

func ResourceProviderSync(syncer ResourceProviderSyncer) resource_provider.ResourceProviderSyncer {
	return ResourceProviderSyncFactory(NewDummySyncFactoryFn(syncer))
}

func (r *ResourceProvisionSyncFunction) UpdateResources(ctx context.Context, config *resource_provider.UpdateResourceInput) (*resource_provider.UpdateResourceResult, error) {
	logger.Info("Starting resource provisioning")

	start := time.Now()

	defer func() {
		logger.Info(fmt.Sprintf("Finished resource provisioning in %s", time.Since(start)))
	}()

	syncer, err := r.syncer.Create(ctx, config.ConfigMap)
	if err != nil {
		return nil, fmt.Errorf("create syncer: %w", err)
	}

	result, err := syncer.UpdateResources(ctx, config)
	if err != nil {
		return result, fmt.Errorf("update resources: %w", err)
	}

	return result, nil
}

func (r *ResourceProvisionSyncFunction) Close() {
	r.syncer.Close()
}
