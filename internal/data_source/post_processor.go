package data_source

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/raito-io/cli/base/constants"
	"github.com/raito-io/cli/base/data_source"
	"github.com/raito-io/cli/base/tag"
	"github.com/raito-io/cli/internal/util/jsonstream"
	string2 "github.com/raito-io/cli/internal/util/string"
)

type PostProcessorConfig struct {
	TagOverwriteKeyForOwners string
	DataSourceId             string
	DataObjectParent         string
	DataObjectExcludes       []string
	TargetLogger             hclog.Logger
}
type PostProcessorResult struct {
	DataObjectsTouchedCount int
}

type PostProcessorOutputFileWriter interface {
	AddDataObjects(dataObjects ...*data_source.DataObject) error
}

type PostProcessor struct {
	dataSourceFileCreatorFactory func(config *data_source.DataSourceSyncConfig) (data_source.DataSourceFileCreator, error)

	config *PostProcessorConfig
}

func NewPostProcessor(config *PostProcessorConfig) PostProcessor {
	return PostProcessor{
		dataSourceFileCreatorFactory: data_source.NewDataSourceFileCreator,
		config:                       config,
	}
}

func (p *PostProcessor) NeedsPostProcessing() bool {
	return p.config.TagOverwriteKeyForOwners != ""
}

func (p *PostProcessor) PostProcess(inputFilePath string, outputFile string) (*PostProcessorResult, error) {
	outputWriter, err := p.dataSourceFileCreatorFactory(&data_source.DataSourceSyncConfig{
		TargetFile:         outputFile,
		DataSourceId:       p.config.DataSourceId,
		DataObjectParent:   p.config.DataObjectParent,
		DataObjectExcludes: p.config.DataObjectExcludes,
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	defer outputWriter.Close()

	p.config.TargetLogger.Debug(fmt.Sprintf("Post processor streaming data objects from file %s", inputFilePath))

	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open input file %q: %s", inputFilePath, err.Error())
	}
	defer inputFile.Close()

	dataObjectsRead := 0
	dataObjectsTouched := 0

	decoder := jsonstream.NewJsonArrayStream[data_source.DataObject](inputFile)
	for jsonStreamResult := range decoder.Stream() {
		p.config.TargetLogger.Debug(fmt.Sprintf("Start post processing data object %d", dataObjectsRead))

		if jsonStreamResult.Err != nil {
			return nil, fmt.Errorf("unable to parse data object (%d): %s", dataObjectsRead, jsonStreamResult.Err.Error())
		}

		do := jsonStreamResult.Result
		p.config.TargetLogger.Info(fmt.Sprintf("Start enriching data object %q", do.FullName))

		enriched, err2 := p.postProcessDataObject(do, outputWriter)
		if err2 != nil {
			return nil, fmt.Errorf("unable to enrich data object (%d): %s", dataObjectsRead, err2.Error())
		}

		if enriched {
			dataObjectsTouched++
		}

		dataObjectsRead++
	}

	if outputWriter.GetDataObjectCount() != dataObjectsRead {
		return nil, fmt.Errorf("post processor wrote %d data objects, while expecting %d", outputWriter.GetDataObjectCount(), dataObjectsRead)
	}

	return &PostProcessorResult{
		DataObjectsTouchedCount: dataObjectsTouched,
	}, nil
}

func (p *PostProcessor) postProcessDataObject(do *data_source.DataObject, outputWriter data_source.DataSourceFileCreator) (bool, error) {
	touched := false
	//
	var raitoOwnerTag *tag.Tag

	for _, t := range do.Tags {
		if strings.EqualFold(t.Key, constants.RaitoOwnerTagKey) {
			raitoOwnerTag = t
			raitoOwnerTag.Value = string2.TrimSpaceInCommaSeparatedList(raitoOwnerTag.Value)

			touched = true

			break
		}
	}

	for _, t := range do.Tags {
		if p.matchedWithTagKey(p.config.TagOverwriteKeyForOwners, t) {
			if raitoOwnerTag != nil {
				raitoOwnerTag.Value = raitoOwnerTag.Value + "," + string2.TrimSpaceInCommaSeparatedList(t.Value)
			} else {
				raitoOwnerTag = &tag.Tag{
					Key:    constants.RaitoOwnerTagKey,
					Value:  string2.TrimSpaceInCommaSeparatedList(t.Value),
					Source: t.Source,
				}

				do.Tags = append(do.Tags, raitoOwnerTag)
			}

			touched = true
		}
	}

	if raitoOwnerTag != nil {
		p.config.TargetLogger.Debug(fmt.Sprintf("adjusting owners for DO (fullpath: %v) to %v", do.FullName, raitoOwnerTag.Value))
	}

	err := outputWriter.AddDataObjects(do)
	if err != nil {
		p.config.TargetLogger.Info(fmt.Sprintf("Error while saving data object to writer %q", do.FullName))
		return touched, err
	}

	return touched, nil
}

func (p *PostProcessor) matchedWithTagKey(overwriteKey string, tag *tag.Tag) bool {
	return tag != nil && overwriteKey != "" && strings.EqualFold(tag.Key, overwriteKey) && tag.Value != ""
}

func (p *PostProcessor) Close(_ context.Context) error {
	return nil
}
