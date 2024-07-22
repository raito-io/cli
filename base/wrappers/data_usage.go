package wrappers

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/raito-io/cli/base/data_usage"
	"github.com/raito-io/cli/base/util/config"
)

//go:generate go run github.com/vektra/mockery/v2 --name=DataUsageStatementHandler --with-expecter
type DataUsageStatementHandler interface {
	AddStatements(statements []data_usage.Statement) error
	GetImportFileSize() uint64
}

//go:generate go run github.com/vektra/mockery/v2 --name=DataUsageSyncer --with-expecter --inpackage
type DataUsageSyncer interface {
	SyncDataUsage(ctx context.Context, fileCreator DataUsageStatementHandler, configParams *config.ConfigMap) error
}

type DataUsageSyncFactoryFn func(ctx context.Context, configParams *config.ConfigMap) (DataUsageSyncer, func(), error)

func DataUsageSync(syncer DataUsageSyncer) *dataUsageSyncFunction {
	return DataUsageSyncFactory(NewDummySyncFactoryFn[config.ConfigMap](syncer))
}

func DataUsageSyncFactory(syncer DataUsageSyncFactoryFn) *dataUsageSyncFunction {
	return &dataUsageSyncFunction{
		syncer:             NewSyncFactory(syncer),
		fileCreatorFactory: data_usage.NewDataUsageFileCreator,
	}
}

type dataUsageSyncFunction struct {
	data_usage.DataUsageSyncerVersionHandler

	syncer             SyncFactory[config.ConfigMap, DataUsageSyncer]
	fileCreatorFactory func(config *data_usage.DataUsageSyncConfig) (data_usage.DataUsageFileCreator, error)
}

func (s *dataUsageSyncFunction) SyncDataUsage(ctx context.Context, config *data_usage.DataUsageSyncConfig) (_ *data_usage.DataUsageSyncResult, err error) {
	defer func() {
		if err != nil {
			logger.Error(fmt.Sprintf("Failure during data usage sync: %v", err))
		}

		if r := recover(); r != nil {
			err = fmt.Errorf("panic during data usage sync: %v", r)

			logger.Error(fmt.Sprintf("Panic during data usage sync: %v\n\n%s", r, string(debug.Stack())))
		}
	}()

	logger.Info("Starting data usage synchronisation")
	logger.Debug("Creating file for storing data usage")

	fileCreator, err := s.fileCreatorFactory(config)
	if err != nil {
		return nil, err
	}

	defer fileCreator.Close()

	syncer, err := s.syncer.Create(ctx, config.ConfigMap)
	if err != nil {
		return nil, err
	}

	sec, err := timedExecution(func() error {
		return syncer.SyncDataUsage(ctx, fileCreator, config.ConfigMap)
	})

	if err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("Retrieved %d rows and written them to file (total size %d bytes), for a total time of %s",
		fileCreator.GetStatementCount(), fileCreator.GetImportFileSize(), sec))

	return &data_usage.DataUsageSyncResult{
		Statements: int32(fileCreator.GetStatementCount()),
	}, nil
}

func (s *dataUsageSyncFunction) Close() {
	s.syncer.Close()
}
