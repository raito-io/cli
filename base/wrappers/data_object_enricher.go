package wrappers

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/raito-io/cli/base/data_object_enricher"
	"github.com/raito-io/cli/base/data_source"
	"github.com/raito-io/cli/base/util/config"
	error2 "github.com/raito-io/cli/internal/error"
	"github.com/raito-io/cli/internal/util/jsonstream"
)

type DataObjectWriter interface {
	AddDataObjects(dataObjects ...*data_source.DataObject) error
}

//go:generate go run github.com/vektra/mockery/v2 --name=DataObjectEnricherI --with-expecter --inpackage
type DataObjectEnricherI interface {
	// Initialize allows the plugin to do any preparation work (like making a connection to the enrichment source) and store the writer in a variable.
	Initialize(ctx context.Context, dataObjectWriter DataObjectWriter, config map[string]string) error

	// Enrich method receives every data object separately. The plugin can decide to skip or buffer things. All data objects must be written to the DataObjectWriter.
	// First argument indicates if the data object was enriched
	Enrich(ctx context.Context, dataObject *data_source.DataObject) (bool, error)
}

type DataObjectEnricherFactoryFn func(ctx context.Context, config *config.ConfigMap) (DataObjectEnricherI, func(), error)

func DataObjectEnricher(enricher DataObjectEnricherI) *dataObjectEnricherFunction {
	return DataObjectEnricherFactory(NewDummySyncFactoryFn[config.ConfigMap](enricher))
}

func DataObjectEnricherFactory(enricher DataObjectEnricherFactoryFn) *dataObjectEnricherFunction {
	return &dataObjectEnricherFunction{
		enricher: NewSyncFactory(enricher),
	}
}

type dataObjectEnricherFunction struct {
	data_object_enricher.DataObjectEnricherVersionHandler

	enricher SyncFactory[config.ConfigMap, DataObjectEnricherI]
}

func (f *dataObjectEnricherFunction) Enrich(ctx context.Context, config *data_object_enricher.DataObjectEnricherConfig) (_ *data_object_enricher.DataObjectEnricherResult, err error) {
	logger.Info("Enriching data objects...")

	defer func() {
		if err != nil {
			logger.Error(fmt.Sprintf("Failure during data object enrich: %v", err))
		}

		if r := recover(); r != nil {
			err = error2.NewRecoverErrorf("panic during data object enrich: %v", r)

			logger.Error(fmt.Sprintf("Panic during data object enrich: %v\n\n%s", r, string(debug.Stack())))
		}
	}()

	fileCreator, err := data_source.NewDataSourceFileCreator(&data_source.DataSourceSyncConfig{
		TargetFile:   config.OutputFile,
		DataSourceId: "",
	})
	if err != nil {
		return nil, err
	}

	logger.Info("File creator initialized")

	enricher, err := f.enricher.Create(ctx, config.ConfigMap)
	if err != nil {
		return nil, fmt.Errorf("create enricher: %w", err)
	}

	err = enricher.Initialize(ctx, fileCreator, config.ConfigMap.Parameters)
	if err != nil {
		return nil, err
	}

	logger.Info("Enricher initialized")

	dataObjectsRead := 0

	inputFile, err := os.Open(config.InputFile)
	if err != nil {
		return nil, fmt.Errorf("unable to open input file %q: %s", config.InputFile, err.Error())
	}
	defer inputFile.Close()

	enrichmentCount := 0

	decoder := jsonstream.NewJsonArrayStream[data_source.DataObject](inputFile)
	for jsonStreamResult := range decoder.Stream() {
		logger.Info(fmt.Sprintf("Reading row %d", dataObjectsRead))

		if jsonStreamResult.Err != nil {
			return nil, fmt.Errorf("unable to parse data object (%d): %s", dataObjectsRead, jsonStreamResult.Err.Error())
		}

		do := jsonStreamResult.Result

		enriched, err2 := enricher.Enrich(ctx, do)
		if err2 != nil {
			return nil, fmt.Errorf("unable to enrich data object (%d): %s", dataObjectsRead, err2.Error())
		}

		if enriched {
			enrichmentCount++
		}

		dataObjectsRead++
	}

	if fileCreator.GetDataObjectCount() != dataObjectsRead {
		return nil, fmt.Errorf("enricher wrote %d data objects, while expecting %d", fileCreator.GetDataObjectCount(), dataObjectsRead)
	}

	return &data_object_enricher.DataObjectEnricherResult{
		Enriched: int32(enrichmentCount), //nolint:gosec
	}, nil
}

func (f *dataObjectEnricherFunction) Close() {
	f.enricher.Close()
}
