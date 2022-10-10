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

func DataAccessSync(syncer AccessProviderSyncer) *DataAccessSyncFunction {
	return &DataAccessSyncFunction{
		Syncer:                           syncer,
		accessFileCreatorFactory:         sync_from_target.NewAccessProviderFileCreator,
		accessFeedbackFileCreatorFactory: sync_to_target.NewFeedbackFileCreator,
		accessProviderParserFactory:      sync_to_target.NewAccessProviderFileParser,
	}
}

type DataAccessSyncFunction struct {
	Syncer                           AccessProviderSyncer
	accessFileCreatorFactory         func(config *access_provider.AccessSyncFromTarget) (sync_from_target.AccessProviderFileCreator, error)
	accessFeedbackFileCreatorFactory func(config *access_provider.AccessSyncToTarget) (sync_to_target.SyncFeedbackFileCreator, error)
	accessProviderParserFactory      func(config *access_provider.AccessSyncToTarget) (sync_to_target.AccessProviderImportFileParser, error)
}

func (s *DataAccessSyncFunction) SyncFromTarget(config *access_provider.AccessSyncFromTarget) access_provider.AccessSyncResult {
	ctx := context.Background()

	logger.Info("Starting data access synchronisation from target")
	logger.Debug("Creating file for storing access providers")

	fileCreator, err := s.accessFileCreatorFactory(config)
	if err != nil {
		logger.Error(err.Error())

		return mapErrorToAccessSyncResult(err)
	}
	defer fileCreator.Close()

	sec, err := timedExecution(func() error {
		return s.Syncer.SyncAccessProvidersFromTarget(ctx, fileCreator, &config.ConfigMap)
	})

	if err != nil {
		logger.Error(err.Error())

		return mapErrorToAccessSyncResult(err)
	}

	logger.Info(fmt.Sprintf("Fetched %d access provider in %s", fileCreator.GetAccessProviderCount(), sec))

	return access_provider.AccessSyncResult{}
}

func (s *DataAccessSyncFunction) SyncToTarget(config *access_provider.AccessSyncToTarget) access_provider.AccessSyncResult {
	ctx := context.Background()

	logger.Info("Starting data access synchronisation to target")

	accessProviderParser, err := s.accessProviderParserFactory(config)
	if err != nil {
		logger.Error(err.Error())

		return mapErrorToAccessSyncResult(err)
	}

	dar, err := accessProviderParser.ParseAccessProviders()
	if err != nil {
		logger.Error(err.Error())

		return mapErrorToAccessSyncResult(err)
	}

	prefix := config.Prefix
	accessAsCode := config.Prefix != ""

	var sec time.Duration

	if accessAsCode {
		sec, err = timedExecution(func() error {
			return s.Syncer.SyncAccessAsCodeToTarget(ctx, dar, prefix, &config.ConfigMap)
		})
	} else {
		feedbackFile, err2 := s.accessFeedbackFileCreatorFactory(config)
		if err2 != nil {
			logger.Error(err2.Error())

			return mapErrorToAccessSyncResult(err2)
		}
		defer feedbackFile.Close()

		sec, err = timedExecution(func() error {
			return s.Syncer.SyncAccessProviderToTarget(ctx, dar, feedbackFile, &config.ConfigMap)
		})
	}

	if err != nil {
		logger.Error(err.Error())

		return mapErrorToAccessSyncResult(err)
	}

	logger.Info(fmt.Sprintf("Successfully synced access providers to target in %s", sec))

	return access_provider.AccessSyncResult{}
}

func mapErrorToAccessSyncResult(err error) access_provider.AccessSyncResult {
	return access_provider.AccessSyncResult{
		Error: e.ToErrorResult(err),
	}
}
