package data_source

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/raito-io/cli/base/data_object_enricher"
	dspc "github.com/raito-io/cli/base/data_source"
	baseconfig "github.com/raito-io/cli/base/util/config"
	"github.com/raito-io/cli/internal/file"
	"github.com/raito-io/cli/internal/job"
	"github.com/raito-io/cli/internal/plugin"
	"github.com/raito-io/cli/internal/target/types"
	"github.com/raito-io/cli/internal/util/tag"
	"github.com/raito-io/cli/internal/version_management"
)

const JSON_EXTENSION = ".json"

type DataSourceImportResult struct {
	DataObjectsAdded   int `json:"dataObjectsAdded"`
	DataObjectsUpdated int `json:"dataObjectsUpdated"`
	DataObjectsRemoved int `json:"dataObjectsRemoved"`

	Warnings []string `json:"warnings"`
}

type DataSourceSync struct {
	TargetConfig *types.BaseTargetConfig
	JobId        string

	result *job.TaskResult
}

func (s *DataSourceSync) IsClientValid(ctx context.Context, c plugin.PluginClient) (bool, error) {
	dss, err := c.GetDataSourceSyncer()
	if err != nil {
		return false, err
	}

	return version_management.IsValidToSync(ctx, dss, dspc.MinimalCliVersion)
}

func (s *DataSourceSync) GetParts() []job.TaskPart {
	return []job.TaskPart{s}
}

func (s *DataSourceSync) StartSyncAndQueueTaskPart(ctx context.Context, client plugin.PluginClient, statusUpdater job.TaskEventUpdater) (job.JobStatus, string, error) {
	cn := strings.Replace(s.TargetConfig.ConnectorName, "/", "-", -1)

	targetFile, err := filepath.Abs(file.CreateUniqueFileName(cn+"-ds", "json"))
	if err != nil {
		return job.Failed, "", err
	}

	s.TargetConfig.TargetLogger.Debug(fmt.Sprintf("Using %q as data source target file", targetFile))

	if s.TargetConfig.DeleteTempFiles {
		defer os.RemoveAll(targetFile)
	}

	doParent := ""
	if s.TargetConfig.DataObjectParent != nil {
		doParent = *s.TargetConfig.DataObjectParent
	}

	syncerConfig := dspc.DataSourceSyncConfig{
		ConfigMap:          &baseconfig.ConfigMap{Parameters: s.TargetConfig.Parameters},
		TargetFile:         targetFile,
		DataSourceId:       s.TargetConfig.DataSourceId,
		DataObjectParent:   doParent,
		DataObjectExcludes: s.TargetConfig.DataObjectExcludes,
	}

	dss, err := client.GetDataSourceSyncer()
	if err != nil {
		return job.Failed, "", err
	}

	s.TargetConfig.TargetLogger.Info("Fetching data source metadata")

	md, err := dss.GetDataSourceMetaData(ctx, syncerConfig.ConfigMap)
	if err != nil {
		return job.Failed, "", err
	}

	s.TargetConfig.TargetLogger.Info("Updating data source metadata")
	err = SetMetaData(ctx, s.TargetConfig, md)

	if err != nil {
		return job.Failed, "", err
	}

	s.TargetConfig.TargetLogger.Info("Gathering data objects from the data source")

	res, err := dss.SyncDataSource(ctx, &syncerConfig)
	if err != nil {
		return job.Failed, "", err
	}

	if res.Error != nil { //nolint:staticcheck
		return job.Failed, "", errors.New(res.Error.ErrorMessage) //nolint:staticcheck
	}

	// Fetching the tagSource from the plugin
	tagSourcesScope, err := tag.FetchTagSourceFromPlugin(ctx, client, nil)

	enrichedTargetFile, createdFiles, tagSourcesScope, err := s.enrichDataObjects(ctx, targetFile, tagSourcesScope)

	if s.TargetConfig.DeleteTempFiles && len(createdFiles) > 0 {
		defer func() {
			for _, f := range createdFiles {
				os.RemoveAll(f)
			}
		}()
	}

	if err != nil {
		return job.Failed, "", err
	}

	postProcessor := NewPostProcessor(&PostProcessorConfig{
		TagOverwriteKeyForOwners: s.TargetConfig.TagOverwriteKeyForDataObjectOwners,
		DataSourceId:             syncerConfig.DataSourceId,
		DataObjectParent:         syncerConfig.DataObjectParent,
		DataObjectExcludes:       syncerConfig.DataObjectExcludes,
		TargetLogger:             s.TargetConfig.TargetLogger,
	})

	toProcessFile := enrichedTargetFile
	if postProcessor.NeedsPostProcessing() {
		toProcessFile, _, err = s.postProcessDataObjects(postProcessor, enrichedTargetFile)
		if err != nil {
			return job.Failed, "", err
		}
	}

	deleteUntouched := s.TargetConfig.DeleteUntouched
	if deleteUntouched && s.TargetConfig.DataObjectParent != nil && *s.TargetConfig.DataObjectParent != "" {
		deleteUntouched = false
	}

	importerConfig := DataSourceImportConfig{
		BaseTargetConfig: *s.TargetConfig,
		TargetFile:       toProcessFile,
		DeleteUntouched:  deleteUntouched,
		TagSourcesScope:  tagSourcesScope,
	}
	dsImporter := NewDataSourceImporter(&importerConfig, statusUpdater)

	s.TargetConfig.TargetLogger.Info("Importing data objects into Raito")
	status, subtaskId, err := dsImporter.TriggerImport(ctx, s.JobId)

	if err != nil {
		return job.Failed, "", err
	}

	if status == job.Queued {
		s.TargetConfig.TargetLogger.Info("Successfully queued import job. Wait until remote processing is done.")
	}

	s.TargetConfig.TargetLogger.Debug(fmt.Sprintf("Current status: %s", status))

	return status, subtaskId, nil
}

func (s *DataSourceSync) enrichDataObjects(ctx context.Context, sourceFile string, tagSourcesScope []string) (string, []string, []string, error) {
	enrichedFile := sourceFile

	var newFiles []string

	var err error

	if len(s.TargetConfig.DataObjectEnrichers) > 0 {
		for i, enricher := range s.TargetConfig.DataObjectEnrichers {
			start := time.Now()

			s.TargetConfig.TargetLogger.Info(fmt.Sprintf("Calling enricher %q", enricher.Name))

			enrichmentCount := 0
			enrichedFile, enrichmentCount, tagSourcesScope, err = s.callEnricher(ctx, enricher, enrichedFile, i, tagSourcesScope)

			if enrichedFile != "" {
				newFiles = append(newFiles, enrichedFile)
			}

			if err != nil {
				return "", newFiles, tagSourcesScope, err
			}

			s.TargetConfig.TargetLogger.Info(fmt.Sprintf("%d data objects enriched (%s) in %s", enrichmentCount, enricher.Name, time.Since(start).Round(time.Millisecond)))
		}
	}

	return enrichedFile, newFiles, tagSourcesScope, nil
}

func (s *DataSourceSync) postProcessDataObjects(postProcessor PostProcessor, toProcessFile string) (string, int, error) {
	postProcessedFile := toProcessFile
	fileSuffix := "-post-processed"

	// Generate a unique file name for the post processing
	if strings.Contains(postProcessedFile, fileSuffix) {
		postProcessedFile = postProcessedFile[0:strings.LastIndex(postProcessedFile, fileSuffix)] + fileSuffix + JSON_EXTENSION
	} else {
		postProcessedFile = postProcessedFile[0:strings.LastIndex(postProcessedFile, JSON_EXTENSION)] + fileSuffix + JSON_EXTENSION
	}

	res, err := postProcessor.PostProcess(toProcessFile, postProcessedFile)
	if err != nil {
		return toProcessFile, 0, err
	}

	return postProcessedFile, res.DataObjectsTouchedCount, nil
}

func (s *DataSourceSync) callEnricher(ctx context.Context, enricher *types.EnricherConfig, sourceFile string, index int, tagSourcesScope []string) (string, int, []string, error) {
	client, err := plugin.NewPluginClient(enricher.ConnectorName, enricher.ConnectorVersion, s.TargetConfig.TargetLogger)
	if err != nil {
		s.TargetConfig.TargetLogger.Error(fmt.Sprintf("Error initializing enricher plugin %q: %s", enricher.ConnectorName, err.Error()))
		return "", 0, tagSourcesScope, fmt.Errorf("creating client for plugin %s: %w", enricher.ConnectorName, err)
	}
	defer client.Close()

	doe, err := client.GetDataObjectEnricher()
	if err != nil {
		return "", 0, tagSourcesScope, fmt.Errorf("fetching enricher interface from plugin %s: %w", enricher.ConnectorName, err)
	}

	targetFile := sourceFile
	// Generate a unique file name for the enrichment
	if strings.Contains(targetFile, "-enriched") {
		targetFile = targetFile[0:strings.LastIndex(targetFile, "-enriched")] + "-enriched" + strconv.Itoa(index) + JSON_EXTENSION
	} else {
		targetFile = targetFile[0:strings.LastIndex(targetFile, JSON_EXTENSION)] + "-enriched" + strconv.Itoa(index) + JSON_EXTENSION
	}

	res, err := doe.Enrich(ctx, &data_object_enricher.DataObjectEnricherConfig{
		ConfigMap: &baseconfig.ConfigMap{
			Parameters: enricher.Parameters,
		},
		InputFile:  sourceFile,
		OutputFile: targetFile,
	})

	if err != nil {
		return targetFile, 0, tagSourcesScope, fmt.Errorf("calling enricher plugin %s: %w", enricher.ConnectorName, err)
	}

	// Fetching the tagSource from the plugin
	tagSourcesScope, err = tag.FetchTagSourceFromPlugin(ctx, client, tagSourcesScope)
	if err != nil {
		return "", 0, tagSourcesScope, fmt.Errorf("fetching tag source from plugin %s: %w", enricher.ConnectorName, err)
	}

	return targetFile, int(res.Enriched), tagSourcesScope, nil
}

func (s *DataSourceSync) ProcessResults(results interface{}) error {
	if dsResult, ok := results.(*DataSourceImportResult); ok {
		if dsResult.Warnings != nil && len(dsResult.Warnings) > 0 {
			s.TargetConfig.TargetLogger.Info(fmt.Sprintf("Synced data source with %d warnings (see below). Added: %d - Removed: %d - Updated: %d", len(dsResult.Warnings), dsResult.DataObjectsAdded, dsResult.DataObjectsRemoved, dsResult.DataObjectsUpdated))

			for _, warning := range dsResult.Warnings {
				s.TargetConfig.TargetLogger.Warn(warning)
			}
		} else {
			s.TargetConfig.TargetLogger.Info(fmt.Sprintf("Successfully synced data source. Added: %d - Removed: %d - Updated: %d", dsResult.DataObjectsAdded, dsResult.DataObjectsRemoved, dsResult.DataObjectsUpdated))
		}

		s.result = &job.TaskResult{
			ObjectType: "data objects",
			Added:      dsResult.DataObjectsAdded,
			Removed:    dsResult.DataObjectsRemoved,
			Updated:    dsResult.DataObjectsUpdated,
			Failed:     len(dsResult.Warnings),
		}

		return nil
	}

	return fmt.Errorf("failed to load results")
}

func (s *DataSourceSync) GetResultObject() interface{} {
	return &DataSourceImportResult{}
}

func (s *DataSourceSync) GetTaskResults() []job.TaskResult {
	if s.result == nil {
		return nil
	}

	return []job.TaskResult{*s.result}
}
