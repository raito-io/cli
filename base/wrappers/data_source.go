package wrappers

import (
	"context"
	"fmt"
	"time"

	"github.com/raito-io/cli/base/data_source"
	"github.com/raito-io/cli/base/util/config"
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
	GetDataSourceMetaData(ctx context.Context) (*data_source.MetaData, error)
}

func DataSourceSync(syncer DataSourceSyncer) *dataSourceSyncFunction {
	return &dataSourceSyncFunction{
		syncer:             syncer,
		fileCreatorFactory: data_source.NewDataSourceFileCreator,
	}
}

type dataSourceSyncFunction struct {
	data_source.DataSourceSyncerVersionHandler

	syncer             DataSourceSyncer
	fileCreatorFactory func(config *data_source.DataSourceSyncConfig) (data_source.DataSourceFileCreator, error)
}

func (s *dataSourceSyncFunction) SyncDataSource(ctx context.Context, config *data_source.DataSourceSyncConfig) (_ *data_source.DataSourceSyncResult, err error) {
	defer func() {
		if err != nil {
			logger.Error(fmt.Sprintf("Failure during data source sync: %v", err))
		}
	}()

	logger.Info("Starting data source synchronisation")
	logger.Debug("Creating file for storing data source")

	fileCreator, err := s.fileCreatorFactory(config)
	if err != nil {
		return nil, err
	}
	defer fileCreator.Close()

	start := time.Now()

	err = s.syncer.SyncDataSource(ctx, fileCreator, config.ConfigMap)
	if err != nil {
		return nil, err
	}

	sec := time.Since(start).Round(time.Millisecond)

	logger.Info(fmt.Sprintf("Fetched %d data objects in %s", fileCreator.GetDataObjectCount(), sec))

	return &data_source.DataSourceSyncResult{
		DataObjects: int32(fileCreator.GetDataObjectCount()),
	}, nil
}

func (s *dataSourceSyncFunction) GetDataSourceMetaData(ctx context.Context) (*data_source.MetaData, error) {
	return s.syncer.GetDataSourceMetaData(ctx)
}
