package wrappers

import (
	"context"
	"fmt"
	"time"

	"github.com/raito-io/cli/base/data_source"
	"github.com/raito-io/cli/base/util/config"
	e "github.com/raito-io/cli/base/util/error"
)

//go:generate go run github.com/vektra/mockery/v2 --name=DataSourceObjectHandler --with-expecter
type DataSourceObjectHandler interface {
	AddDataObjects(dataObjects ...*data_source.DataObject) error
	SetDataSourceName(name string)
	SetDataSourceFullname(name string)
	SetDataSourceDescription(desc string)
}

//go:generate go run github.com/vektra/mockery/v2 --name=DataSourceSyncer --with-expecter --inpackage
type DataSourceSyncer interface {
	SyncDataSource(ctx context.Context, dataSourceHandler DataSourceObjectHandler, configParams *config.ConfigMap) error
}

func DataSourceSync(syncer DataSourceSyncer) *dataSourceSyncFunction {
	return &dataSourceSyncFunction{
		syncer:             syncer,
		fileCreatorFactory: data_source.NewDataSourceFileCreator,
	}
}

type dataSourceSyncFunction struct {
	syncer             DataSourceSyncer
	fileCreatorFactory func(config *data_source.DataSourceSyncConfig) (data_source.DataSourceFileCreator, error)
}

func (s *dataSourceSyncFunction) SyncDataSource(config *data_source.DataSourceSyncConfig) data_source.DataSourceSyncResult {
	ctx := context.Background()

	logger.Info("Starting data source synchronisation")
	logger.Debug("Creating file for storing data usage")

	fileCreator, err := s.fileCreatorFactory(config)
	if err != nil {
		logger.Error(err.Error())

		return mapErrorToDataSourceSyncResult(err)
	}
	defer fileCreator.Close()

	start := time.Now()

	err = s.syncer.SyncDataSource(ctx, fileCreator, &config.ConfigMap)
	if err != nil {
		logger.Error(err.Error())

		return mapErrorToDataSourceSyncResult(err)
	}

	sec := time.Since(start).Round(time.Millisecond)

	logger.Info(fmt.Sprintf("Fetched %d data objects in %s", fileCreator.GetDataObjectCount(), sec))

	return data_source.DataSourceSyncResult{}
}

func mapErrorToDataSourceSyncResult(err error) data_source.DataSourceSyncResult {
	return data_source.DataSourceSyncResult{
		Error: e.ToErrorResult(err),
	}
}
