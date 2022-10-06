package wrappers

import (
	"context"
	"fmt"
	"time"

	"github.com/raito-io/cli/base/identity_store"
	"github.com/raito-io/cli/base/util/config"
	e "github.com/raito-io/cli/base/util/error"
)

//go:generate go run github.com/vektra/mockery/v2 --name=IdentityStoreIdentityHandler --with-expecter
type IdentityStoreIdentityHandler interface {
	AddGroups(groups ...*identity_store.Group) error
	AddUsers(user ...*identity_store.User) error
}

//go:generate go run github.com/vektra/mockery/v2 --name=IdentityStoreSyncer --with-expecter --inpackage
type IdentityStoreSyncer interface {
	SyncIdentityStore(ctx context.Context, identityHandler IdentityStoreIdentityHandler, configMap *config.ConfigMap) error
}

func IdentityStoreSync(syncer IdentityStoreSyncer) *identityStoreSyncFunction {
	return &identityStoreSyncFunction{
		syncer:                 syncer,
		identityHandlerFactory: identity_store.NewIdentityStoreFileCreator,
	}
}

type identityStoreSyncFunction struct {
	syncer                 IdentityStoreSyncer
	identityHandlerFactory func(config *identity_store.IdentityStoreSyncConfig) (identity_store.IdentityStoreFileCreator, error)
}

func (s *identityStoreSyncFunction) SyncIdentityStore(config *identity_store.IdentityStoreSyncConfig) identity_store.IdentityStoreSyncResult {
	ctx := context.Background()

	logger.Info("Starting identity store synchronisation")
	logger.Debug("Creating file for storing identity information")

	fileCreator, err := s.identityHandlerFactory(config)
	if err != nil {
		logger.Error(err.Error())

		return mapError(err)
	}
	defer fileCreator.Close()

	start := time.Now()

	err = s.syncer.SyncIdentityStore(ctx, fileCreator, &config.ConfigMap)
	if err != nil {
		logger.Error(err.Error())

		return mapError(err)
	}

	sec := time.Since(start).Round(time.Millisecond)
	logger.Info(fmt.Sprintf("Fetched %d users and %d groups in %s", fileCreator.GetUserCount(), fileCreator.GetGroupCount(), sec))

	return identity_store.IdentityStoreSyncResult{}
}

func mapError(err error) identity_store.IdentityStoreSyncResult {
	return identity_store.IdentityStoreSyncResult{
		Error: e.ToErrorResult(err),
	}
}
