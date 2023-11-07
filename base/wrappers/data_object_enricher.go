package wrappers

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/bcicen/jstream"

	"github.com/raito-io/cli/base/data_object_enricher"
	"github.com/raito-io/cli/base/data_source"
	"github.com/raito-io/cli/base/tag"
	"github.com/raito-io/cli/base/util/config"
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
	return DataObjectEnricherFactory(NewDummySyncFactoryFn(enricher))
}

func DataObjectEnricherFactory(enricher DataObjectEnricherFactoryFn) *dataObjectEnricherFunction {
	return &dataObjectEnricherFunction{
		enricher: NewSyncFactory(enricher),
	}
}

type dataObjectEnricherFunction struct {
	data_object_enricher.DataObjectEnricherVersionHandler

	enricher SyncFactory[DataObjectEnricherI]
}

func (f *dataObjectEnricherFunction) Enrich(ctx context.Context, config *data_object_enricher.DataObjectEnricherConfig) (*data_object_enricher.DataObjectEnricherResult, error) {
	logger.Info("Enriching data objects...")

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

	enrichmentCount := 0

	decoder := jstream.NewDecoder(inputFile, 1)
	for doRow := range decoder.Stream() {
		logger.Info(fmt.Sprintf("Reading row %d", dataObjectsRead))

		do, err2 := createDataObjectFromRow(doRow)

		logger.Info(fmt.Sprintf("Start enriching data object %q", do.FullName))

		if err2 != nil {
			return nil, fmt.Errorf("unable to parse data object (%d): %s", dataObjectsRead, err2.Error())
		}

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
		Enriched: int32(enrichmentCount),
	}, nil
}

func (f *dataObjectEnricherFunction) Close() {
	f.enricher.Close()
}

func createDataObjectFromRow(row *jstream.MetaValue) (*data_source.DataObject, error) {
	if row.ValueType != jstream.Object {
		return nil, errors.New("illegal format for data object definition in source file")
	}

	var values = row.Value.(map[string]interface{})

	do := data_source.DataObject{
		ExternalId:       getStringValue(values, "externalId"),
		Name:             getStringValue(values, "name"),
		FullName:         getStringValue(values, "fullName"),
		Type:             getStringValue(values, "type"),
		Description:      getStringValue(values, "description"),
		ParentExternalId: getStringValue(values, "parentExternalId"),
	}

	if t, found := values["tags"]; found && t != nil {
		if tags, ok := t.([]interface{}); ok {
			do.Tags = make([]*tag.Tag, 0, len(tags))

			for _, tagInput := range tags {
				if tagObj, tok := tagInput.(map[string]interface{}); tok {
					do.Tags = append(do.Tags, &tag.Tag{
						Key:    getStringValue(tagObj, "key"),
						Value:  getStringValue(tagObj, "value"),
						Source: getStringValue(tagObj, "source"),
					})
				}
			}
		}
	}

	return &do, nil
}

func getStringValue(row map[string]interface{}, key string) string {
	if v, found := row[key]; found {
		if vs, ok := v.(string); ok {
			return vs
		}
	}

	return ""
}
