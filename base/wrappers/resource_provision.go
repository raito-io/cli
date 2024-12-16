package wrappers

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/raito-io/cli/base/resource_provider"
	error2 "github.com/raito-io/cli/internal/error"
)

type ResourceProviderSyncer interface {
	UpdateResources(ctx context.Context, config *resource_provider.UpdateResourceInput) (*resource_provider.UpdateResourceResult, error)
}

type ResourceSyncFactoryFn func(ctx context.Context, configParams *resource_provider.UpdateResourceInput) (ResourceProviderSyncer, func(), error)

type ResourceProvisionSyncFunction struct {
	resource_provider.ResourceProviderSyncerVersionHandler

	syncer SyncFactory[resource_provider.UpdateResourceInput, ResourceProviderSyncer]
}

func ResourceProviderSyncFactory(syncer ResourceSyncFactoryFn) resource_provider.ResourceProviderSyncer {
	return &ResourceProvisionSyncFunction{
		syncer: NewSyncFactory(syncer),
	}
}

func ResourceProviderSync(syncer ResourceProviderSyncer) resource_provider.ResourceProviderSyncer {
	return ResourceProviderSyncFactory(NewDummySyncFactoryFn[resource_provider.UpdateResourceInput](syncer))
}

func (r *ResourceProvisionSyncFunction) UpdateResources(ctx context.Context, config *resource_provider.UpdateResourceInput) (_ *resource_provider.UpdateResourceResult, err error) {
	logger.Info("Starting resource provisioning")

	start := time.Now()

	defer func() {
		if err != nil {
			logger.Error(fmt.Sprintf("Failure during resource provisioning: %v", err))
		} else if r := recover(); r != nil {
			err = error2.NewRecoverErrorf("panic during resource provisioning: %v", r)

			logger.Error(fmt.Sprintf("Panic during resource provisioning: %v\n\n%s", r, string(debug.Stack())))
		} else {
			logger.Info(fmt.Sprintf("Finished resource provisioning in %s", time.Since(start)))
		}
	}()

	syncer, err := r.syncer.Create(ctx, config)
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
