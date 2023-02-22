package wrappers

import (
	"context"
	"fmt"
	"time"

	"github.com/raito-io/cli/base/access_provider"
	"github.com/raito-io/cli/base/access_provider/sync_from_target"
	"github.com/raito-io/cli/base/access_provider/sync_to_target"
	"github.com/raito-io/cli/base/util/config"
	e "github.com/raito-io/cli/base/util/error"
)

//go:generate go run github.com/vektra/mockery/v2 --name=AccessProviderHandler --with-expecter
type AccessProviderHandler interface {
	AddAccessProviders(dataAccessList ...*sync_from_target.AccessProvider) error
}

//go:generate go run github.com/vektra/mockery/v2 --name=AccessProviderFeedbackHandler --with-expecter
type AccessProviderFeedbackHandler interface {
	AddAccessProviderFeedback(accessProviderId string, accessFeedback ...sync_to_target.AccessSyncFeedbackInformation) error
}

//go:generate go run github.com/vektra/mockery/v2 --name=AccessProviderSyncer --with-expecter --inpackage
type AccessProviderSyncer interface {
	SyncAccessProvidersFromTarget(ctx context.Context, accessProviderHandler AccessProviderHandler, configMap *config.ConfigMap) error
	SyncAccessAsCodeToTarget(ctx context.Context, accessProviders *sync_to_target.AccessProviderImport, prefix string, configMap *config.ConfigMap) error
	SyncAccessProviderToTarget(ctx context.Context, accessProviders *sync_to_target.AccessProviderImport, accessProviderFeedbackHandler AccessProviderFeedbackHandler, configMap *config.ConfigMap) error
}

func DataAccessSync(syncer AccessProviderSyncer, configOpt ...func(config *access_provider.AccessSyncConfig)) *DataAccessSyncFunction {
	config := access_provider.AccessSyncConfig{}

	for _, fn := range configOpt {
		fn(&config)
	}

	return &DataAccessSyncFunction{
		Syncer:                           syncer,
		accessFileCreatorFactory:         sync_from_target.NewAccessProviderFileCreator,
		accessFeedbackFileCreatorFactory: sync_to_target.NewFeedbackFileCreator,
		accessProviderParserFactory:      sync_to_target.NewAccessProviderFileParser,

		config: config,
	}
}

type DataAccessSyncFunction struct {
	Syncer                           AccessProviderSyncer
	accessFileCreatorFactory         func(config *access_provider.AccessSyncFromTarget) (sync_from_target.AccessProviderFileCreator, error)
	accessFeedbackFileCreatorFactory func(config *access_provider.AccessSyncToTarget) (sync_to_target.SyncFeedbackFileCreator, error)
	accessProviderParserFactory      func(config *access_provider.AccessSyncToTarget) (sync_to_target.AccessProviderImportFileParser, error)

	config access_provider.AccessSyncConfig
}

func (s *DataAccessSyncFunction) SyncFromTarget(ctx context.Context, config *access_provider.AccessSyncFromTarget) (*access_provider.AccessSyncResult, error) {
	logger.Info("Starting data access synchronisation from target")
	logger.Debug("Creating file for storing access providers")

	fileCreator, err := s.accessFileCreatorFactory(config)
	if err != nil {
		logger.Error(err.Error())

		return mapErrorToAccessSyncResult(err), nil
	}
	defer fileCreator.Close()

	sec, err := timedExecution(func() error {
		return s.Syncer.SyncAccessProvidersFromTarget(ctx, fileCreator, config.ConfigMap)
	})

	if err != nil {
		logger.Error(err.Error())

		return mapErrorToAccessSyncResult(err), nil
	}

	logger.Info(fmt.Sprintf("Fetched %d access provider in %s", fileCreator.GetAccessProviderCount(), sec))

	return &access_provider.AccessSyncResult{}, nil
}

func (s *DataAccessSyncFunction) SyncToTarget(ctx context.Context, config *access_provider.AccessSyncToTarget) (*access_provider.AccessSyncResult, error) {
	logger.Info("Starting data access synchronisation to target")

	accessProviderParser, err := s.accessProviderParserFactory(config)
	if err != nil {
		logger.Error(err.Error())

		return mapErrorToAccessSyncResult(err), nil
	}

	dar, err := accessProviderParser.ParseAccessProviders()
	if err != nil {
		logger.Error(err.Error())

		return mapErrorToAccessSyncResult(err), nil
	}

	prefix := config.Prefix
	accessAsCode := config.Prefix != ""

	var sec time.Duration

	if accessAsCode {
		sec, err = timedExecution(func() error {
			return s.Syncer.SyncAccessAsCodeToTarget(ctx, dar, prefix, config.ConfigMap)
		})
	} else {
		feedbackFile, err2 := s.accessFeedbackFileCreatorFactory(config)
		if err2 != nil {
			logger.Error(err2.Error())

			return mapErrorToAccessSyncResult(err2), nil
		}
		defer feedbackFile.Close()

		sec, err = timedExecution(func() error {
			return s.Syncer.SyncAccessProviderToTarget(ctx, dar, feedbackFile, config.ConfigMap)
		})
	}

	if err != nil {
		logger.Error(err.Error())

		return mapErrorToAccessSyncResult(err), nil
	}

	logger.Info(fmt.Sprintf("Successfully synced access providers to target in %s", sec))

	return &access_provider.AccessSyncResult{}, nil
}

func (s *DataAccessSyncFunction) SyncConfig(_ context.Context) (*access_provider.AccessSyncConfig, error) {
	return &s.config, nil
}

func mapErrorToAccessSyncResult(err error) *access_provider.AccessSyncResult {
	return &access_provider.AccessSyncResult{
		Error: e.ToErrorResult(err),
	}
}
