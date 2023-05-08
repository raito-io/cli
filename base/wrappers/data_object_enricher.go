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
)

type DataObjectWriter interface {
	AddDataObjects(dataObjects ...*data_source.DataObject) error
}

//go:generate go run github.com/vektra/mockery/v2 --name=DataObjectEnricher --with-expecter --inpackage
type DataObjectEnricherI interface {
	// Initialize allows the plugin to do any preparation work (like making a connection to the enrichment source) and store the writer in a variable.
	Initialize(ctx context.Context, dataObjectWriter DataObjectWriter, config map[string]string) error

	// Enrich method receives every data object separately. The plugin can decide to skip or buffer things. All data objects must be written to the DataObjectWriter.
	Enrich(ctx context.Context, dataObject *data_source.DataObject) error

	// Close allows the plugin to close any connections and make sure that all data objects still in the buffer are handled and written to the DataObjectWriter.
	// The first return parameter is the number of data objects that were actually enriched.
	Close(ctx context.Context) (int, error)
}

func DataObjectEnricher(enricher DataObjectEnricherI) *dataObjectEnricherFunction {
	return &dataObjectEnricherFunction{
		enricher: enricher,
	}
}

type dataObjectEnricherFunction struct {
	data_object_enricher.DataObjectEnricherVersionHandler

	enricher DataObjectEnricherI
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

	err = f.enricher.Initialize(ctx, fileCreator, config.ConfigMap.Parameters)
	if err != nil {
		return nil, err
	}

	logger.Info("Enricher initialized")

	dataObjectsRead := 0

	inputFile, err := os.Open(config.InputFile)
	if err != nil {
		return nil, fmt.Errorf("unable to open input file %q: %s", config.InputFile, err.Error())
	}

	decoder := jstream.NewDecoder(inputFile, 1)
	for doRow := range decoder.Stream() {
		logger.Info(fmt.Sprintf("Reading row %d", dataObjectsRead))

		do, err2 := createDataObjectFromRow(doRow)

		logger.Info(fmt.Sprintf("Start enriching data object %q", do.FullName))

		if err2 != nil {
			return nil, fmt.Errorf("unable to parse data object (%d): %s", dataObjectsRead, err2.Error())
		}

		err2 = f.enricher.Enrich(ctx, do)
		if err2 != nil {
			return nil, fmt.Errorf("unable to enrich data object (%d): %s", dataObjectsRead, err2.Error())
		}

		dataObjectsRead++
	}

	enrichmentCount, err := f.enricher.Close(ctx)
	if err != nil {
		return nil, err
	}

	if fileCreator.GetDataObjectCount() != dataObjectsRead {
		return nil, fmt.Errorf("enricher wrote %d data objects, while expecting %d", fileCreator.GetDataObjectCount(), dataObjectsRead)
	}

	return &data_object_enricher.DataObjectEnricherResult{
		Enriched: int32(enrichmentCount),
	}, nil
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
