package wrappers

import (
	"context"
	"fmt"
	"time"

	"github.com/raito-io/cli/base/data_usage"
	"github.com/raito-io/cli/base/util/config"
	e "github.com/raito-io/cli/base/util/error"
)

//go:generate go run github.com/vektra/mockery/v2 --name=DataUsageStatementHandler --with-expecter
type DataUsageStatementHandler interface {
	AddStatements(statements []data_usage.Statement) error
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
	syncer             DataUsageSyncer
	fileCreatorFactory func(config *data_usage.DataUsageSyncConfig) (data_usage.DataUsageFileCreator, error)
}

func (s *dataUsageSyncFunction) SyncDataUsage(config *data_usage.DataUsageSyncConfig) data_usage.DataUsageSyncResult {
	ctx := context.Background()

	logger.Info("Starting data usage synchronisation")
	logger.Debug("Creating file for storing data usage")

	fileCreator, err := s.fileCreatorFactory(config)
	if err != nil {
		logger.Error(err.Error())

		return data_usage.DataUsageSyncResult{
			Error: e.ToErrorResult(err),
		}
	}
	defer fileCreator.Close()

	start := time.Now()

	err = s.syncer.SyncDataUsage(ctx, fileCreator, &config.ConfigMap)
	if err != nil {
		logger.Error(err.Error())

		return data_usage.DataUsageSyncResult{
			Error: e.ToErrorResult(err),
		}
	}

	sec := time.Since(start).Round(time.Millisecond)
	logger.Info(fmt.Sprintf("Retrieved %d rows and written them to file, for a total time of %s",
		fileCreator.GetStatementCount(), sec))

	return data_usage.DataUsageSyncResult{}
}
