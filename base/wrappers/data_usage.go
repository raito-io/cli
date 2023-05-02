package wrappers

import (
	"context"
	"fmt"

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

func DataUsageSync(syncer DataUsageSyncer) *dataUsageSyncFunction {
	return &dataUsageSyncFunction{
		syncer:             syncer,
		fileCreatorFactory: data_usage.NewDataUsageFileCreator,
	}
}

type dataUsageSyncFunction struct {
	data_usage.DataUsageSyncerVersionHandler

	syncer             DataUsageSyncer
	fileCreatorFactory func(config *data_usage.DataUsageSyncConfig) (data_usage.DataUsageFileCreator, error)
}

func (s *dataUsageSyncFunction) SyncDataUsage(ctx context.Context, config *data_usage.DataUsageSyncConfig) (_ *data_usage.DataUsageSyncResult, err error) {
	defer func() {
		if err != nil {
			logger.Error(fmt.Sprintf("Failure during data usage sync: %v", err))
		}
	}()

	logger.Info("Starting data usage synchronisation")
	logger.Debug("Creating file for storing data usage")

	fileCreator, err := s.fileCreatorFactory(config)
	if err != nil {
		return nil, err
	}

	defer fileCreator.Close()

	sec, err := timedExecution(func() error {
		return s.syncer.SyncDataUsage(ctx, fileCreator, config.ConfigMap)
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
