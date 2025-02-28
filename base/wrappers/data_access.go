package wrappers

import (
	"context"
	"fmt"
	"runtime/debug"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/raito-io/cli/base/access_provider"
	"github.com/raito-io/cli/base/access_provider/sync_from_target"
	"github.com/raito-io/cli/base/access_provider/sync_to_target"
	"github.com/raito-io/cli/base/util/config"
	error2 "github.com/raito-io/cli/internal/error"
)

//go:generate go run github.com/vektra/mockery/v2 --name=AccessProviderHandler --with-expecter
type AccessProviderHandler interface {
	AddAccessProviders(dataAccessList ...*sync_from_target.AccessProvider) error
}

//go:generate go run github.com/vektra/mockery/v2 --name=AccessProviderFeedbackHandler --with-expecter
type AccessProviderFeedbackHandler interface {
	AddAccessProviderFeedback(accessProviderFeedback sync_to_target.AccessProviderSyncFeedback) error
}

//go:generate go run github.com/vektra/mockery/v2 --name=AccessProviderSyncer --with-expecter --inpackage
type AccessProviderSyncer interface {
	SyncAccessProvidersFromTarget(ctx context.Context, accessProviderHandler AccessProviderHandler, configMap *config.ConfigMap) error
	SyncAccessProviderToTarget(ctx context.Context, accessProviders *sync_to_target.AccessProviderImport, accessProviderFeedbackHandler AccessProviderFeedbackHandler, configMap *config.ConfigMap) error
}

type AccessProviderSyncFactoryFn func(ctx context.Context, configMap *config.ConfigMap) (AccessProviderSyncer, func(), error)

func DataAccessSync(syncer AccessProviderSyncer, configOpt ...func(config *access_provider.AccessSyncConfig)) *DataAccessSyncFunction {
	return DataAccessSyncFactory(NewDummySyncFactoryFn[config.ConfigMap](syncer), configOpt...)
}

func DataAccessSyncFactory(syncer AccessProviderSyncFactoryFn, configOpt ...func(config *access_provider.AccessSyncConfig)) *DataAccessSyncFunction {
	obj := &DataAccessSyncFunction{
		Syncer:                           NewSyncFactory(syncer),
		accessFileCreatorFactory:         sync_from_target.NewAccessProviderFileCreator,
		accessFeedbackFileCreatorFactory: sync_to_target.NewFeedbackFileCreator,
		accessProviderParserFactory:      sync_to_target.NewAccessProviderFileParser,

		config: access_provider.AccessSyncConfig{},
	}

	for _, fn := range configOpt {
		fn(&obj.config)
	}

	return obj
}

type DataAccessSyncFunction struct {
	access_provider.AccessSyncerVersionHandler

	Syncer                           SyncFactory[config.ConfigMap, AccessProviderSyncer]
	accessFileCreatorFactory         func(config *access_provider.AccessSyncFromTarget) (sync_from_target.AccessProviderFileCreator, error)
	accessFeedbackFileCreatorFactory func(config *access_provider.AccessSyncToTarget) (sync_to_target.SyncFeedbackFileCreator, error)
	accessProviderParserFactory      func(config *access_provider.AccessSyncToTarget) (sync_to_target.AccessProviderImportFileParser, error)

	config access_provider.AccessSyncConfig
}

func (s *DataAccessSyncFunction) SyncFromTarget(ctx context.Context, config *access_provider.AccessSyncFromTarget) (_ *access_provider.AccessSyncResult, err error) {
	defer func() {
		if err != nil {
			logger.Error(fmt.Sprintf("Failure during access provider sync from target: %v", err))
		}

		if r := recover(); r != nil {
			err = error2.NewRecoverErrorf("panic during access provider sync from target: %v", r)

			logger.Error(fmt.Sprintf("Panic during access provider sync from target: %v\n\n%s", r, string(debug.Stack())))
		}
	}()

	logger.Info("Starting data access synchronisation from target")
	logger.Debug("Creating file for storing access providers")

	fileCreator, err := s.accessFileCreatorFactory(config)
	if err != nil {
		return nil, err
	}
	defer fileCreator.Close()

	syncer, err := s.Syncer.Create(ctx, config.ConfigMap)
	if err != nil {
		return nil, err
	}

	sec, err := timedExecution(func() error {
		return syncer.SyncAccessProvidersFromTarget(ctx, fileCreator, config.ConfigMap)
	})

	if err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("Fetched %d access providers in %s", fileCreator.GetAccessProviderCount(), sec))

	return &access_provider.AccessSyncResult{
		AccessProviderCount: int32(fileCreator.GetAccessProviderCount()), //nolint:gosec
	}, nil
}

func (s *DataAccessSyncFunction) SyncToTarget(ctx context.Context, config *access_provider.AccessSyncToTarget) (_ *access_provider.AccessSyncResult, err error) {
	defer func() {
		if err != nil {
			logger.Error(fmt.Sprintf("Failure during access provider sync to target: %v", err))
		}

		if r := recover(); r != nil {
			err = error2.NewRecoverErrorf("panic during access provider sync to target: %v", r)

			logger.Error(fmt.Sprintf("Panic during access provider sync to target: %v\n\n%s", r, string(debug.Stack())))
		}
	}()

	logger.Info("Starting data access synchronisation to target")

	syncer, err := s.Syncer.Create(ctx, config.ConfigMap)
	if err != nil {
		return nil, err
	}

	accessProviderParser, err := s.accessProviderParserFactory(config)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	dar, err := accessProviderParser.ParseAccessProviders()
	if err != nil {
		return nil, err
	}

	feedbackFile, err2 := s.accessFeedbackFileCreatorFactory(config)
	if err2 != nil {
		return nil, err2
	}
	defer feedbackFile.Close()

	sec, err := timedExecution(func() error {
		return syncer.SyncAccessProviderToTarget(ctx, dar, feedbackFile, config.ConfigMap)
	})

	if err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("Successfully synced access providers to target in %s", sec))

	return &access_provider.AccessSyncResult{
		AccessProviderCount: int32(len(dar.AccessProviders)), //nolint:gosec
	}, nil
}

func (s *DataAccessSyncFunction) SyncConfig(_ context.Context) (*access_provider.AccessSyncConfig, error) {
	return &s.config, nil
}

func (s *DataAccessSyncFunction) Close() {
	s.Syncer.Close()
}
