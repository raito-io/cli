package wrappers

import (
	"context"
	"fmt"
	"runtime/debug"
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
	SyncDataSource(ctx context.Context, dataSourceHandler DataSourceObjectHandler, config *data_source.DataSourceSyncConfig) error
	GetDataSourceMetaData(ctx context.Context, configParams *config.ConfigMap) (*data_source.MetaData, error)
}

type DataSourceSyncFactoryFn func(ctx context.Context, configParams *config.ConfigMap) (DataSourceSyncer, func(), error)

func DataSourceSync(syncer DataSourceSyncer) data_source.DataSourceSyncer {
	return DataSourceSyncFactory(NewDummySyncFactoryFn[config.ConfigMap](syncer))
}

func DataSourceSyncFactory(syncer DataSourceSyncFactoryFn) data_source.DataSourceSyncer {
	return &dataSourceSyncFunction{
		syncer:             NewSyncFactory(syncer),
		fileCreatorFactory: data_source.NewDataSourceFileCreator,
	}
}

type dataSourceSyncFunction struct {
	data_source.DataSourceSyncerVersionHandler

	syncer             SyncFactory[config.ConfigMap, DataSourceSyncer]
	fileCreatorFactory func(config *data_source.DataSourceSyncConfig) (data_source.DataSourceFileCreator, error)
}

func (s *dataSourceSyncFunction) SyncDataSource(ctx context.Context, config *data_source.DataSourceSyncConfig) (_ *data_source.DataSourceSyncResult, err error) {
	defer func() {
		if err != nil {
			logger.Error(fmt.Sprintf("Failure during data source sync: %v", err))
		}

		if r := recover(); r != nil {
			err = fmt.Errorf("panic during data source sync: %v", r)

			logger.Error(fmt.Sprintf("Panic during data source sync: %v\n\n%s", r, string(debug.Stack())))
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

	syncer, err := s.syncer.Create(ctx, config.ConfigMap)
	if err != nil {
		return nil, err
	}

	err = syncer.SyncDataSource(ctx, fileCreator, config)
	if err != nil {
		return nil, err
	}

	sec := time.Since(start).Round(time.Millisecond)

	logger.Info(fmt.Sprintf("Fetched %d data objects in %s", fileCreator.GetDataObjectCount(), sec))

	return &data_source.DataSourceSyncResult{
		DataObjects: int32(fileCreator.GetDataObjectCount()), //nolint:gosec
	}, nil
}

func (s *dataSourceSyncFunction) GetDataSourceMetaData(ctx context.Context, configParams *config.ConfigMap) (*data_source.MetaData, error) {
	syncer, err := s.syncer.Create(ctx, configParams)
	if err != nil {
		return nil, err
	}

	return syncer.GetDataSourceMetaData(ctx, configParams)
}

func (s *dataSourceSyncFunction) Close() {
	s.syncer.Close()
}
