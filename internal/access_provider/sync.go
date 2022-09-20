package access_provider

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hashicorp/go-hclog"

	dapc "github.com/raito-io/cli/base/access_provider"
	baseconfig "github.com/raito-io/cli/base/util/config"
	"github.com/raito-io/cli/internal/file"
	"github.com/raito-io/cli/internal/job"
	"github.com/raito-io/cli/internal/plugin"
	"github.com/raito-io/cli/internal/target"
)

type AccessProviderImportResult struct {
	AccessAdded   int `json:"accessAdded"`
	AccessUpdated int `json:"accessUpdated"`
	AccessRemoved int `json:"accessRemoved"`

	Warnings []string `json:"warnings"`
}

var accessLastCalculated = map[string]int64{}

type DataAccessSync struct {
	TargetConfig  *target.BaseTargetConfig
	JobId         string
	StatusUpdater func(status job.JobStatus)
}

func (s *DataAccessSync) StartSyncAndQueueJob(client plugin.PluginClient) (job.JobStatus, string, error) {
	cn := strings.Replace(s.TargetConfig.ConnectorName, "/", "-", -1)

	targetFile, err := filepath.Abs(file.CreateUniqueFileName(cn+"-da", "json"))
	if err != nil {
		return job.Failed, "", err
	}

	s.TargetConfig.Logger.Debug(fmt.Sprintf("Using %q as data access target file", targetFile))

	if s.TargetConfig.DeleteTempFiles {
		defer os.RemoveAll(targetFile)
	}

	s.TargetConfig.Logger.Info("Fetching access providers for this data source from Raito")
	s.StatusUpdater(job.DataRetrieve)
	daExporter := NewAccessProviderExporter(&AccessProviderExporterConfig{BaseTargetConfig: *s.TargetConfig}, s.StatusUpdater)

	_, dar, err := daExporter.TriggerExport(s.JobId)

	if err != nil {
		return job.Failed, "", err
	}

	if err != nil {
		return job.Failed, "", err
	}

	syncerConfig := dapc.AccessSyncConfig{
		ConfigMap:  baseconfig.ConfigMap{Parameters: s.TargetConfig.Parameters},
		Prefix:     "",
		TargetFile: targetFile,
		SourceFile: dar,
	}

	das, err := client.GetAccessSyncer()
	if err != nil {
		return job.Failed, "", err
	}

	s.TargetConfig.Logger.Info("Synchronizing access providers between Raito and the data source")
	res := das.SyncAccess(&syncerConfig)

	if res.Error != nil {
		return job.Failed, "", err
	}

	importerConfig := AccessProviderImportConfig{
		BaseTargetConfig: *s.TargetConfig,
		TargetFile:       targetFile,
		DeleteUntouched:  s.TargetConfig.DeleteUntouched,
	}

	daImporter := NewAccessProviderImporter(&importerConfig, s.StatusUpdater)

	status, subtaskId, err := daImporter.TriggerImport(s.JobId)
	if err != nil {
		return job.Failed, "", err
	}

	if status == job.Queued {
		s.TargetConfig.Logger.Info("Successfully queued import job. Wait until remote processing is done.")
	}

	s.TargetConfig.Logger.Debug(fmt.Sprintf("Current status: %s", status.String()))

	return status, subtaskId, nil
}

func (s *DataAccessSync) updateLastCalculated(filePath string) error {
	time, err := findLastCalculated(filePath, s.TargetConfig.Logger)
	if err != nil {
		return err
	}
	accessLastCalculated[s.TargetConfig.DataSourceId] = time

	return nil
}

func findLastCalculated(filePath string, logger hclog.Logger) (int64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("error while reading temporary file %q: %s", filePath, err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNr := 0
	searchString := "lastCalculated:"

	for scanner.Scan() {
		line := scanner.Text()
		if index := strings.Index(line, searchString); index >= 0 {
			timeString := strings.TrimSpace(line[index+len(searchString):])
			time, err := strconv.Atoi(timeString)

			if err != nil {
				return 0, fmt.Errorf("unable to parse lastCalculated field in %q: %s", filePath, err.Error())
			}

			return int64(time), nil
		}

		lineNr++

		if lineNr == 10 {
			logger.Info(fmt.Sprintf("Didn't find 'lastCalculated' field in first 10 lines of %q. Giving up.", filePath))
			return 0, nil
		}
	}

	return 0, nil
}

func (s *DataAccessSync) ProcessResults(results interface{}) error {
	if daResult, ok := results.(*AccessProviderImportResult); ok {
		if len(daResult.Warnings) > 0 {
			s.TargetConfig.Logger.Info(fmt.Sprintf("Synced access providers with %d warnings (see below). Added: %d - Removed: %d - Updated: %d", len(daResult.Warnings), daResult.AccessAdded, daResult.AccessRemoved, daResult.AccessUpdated))

			for _, warning := range daResult.Warnings {
				s.TargetConfig.Logger.Warn(warning)
			}
		} else {
			s.TargetConfig.Logger.Info(fmt.Sprintf("Successfully synced access providers. Added: %d - Removed: %d - Updated: %d", daResult.AccessAdded, daResult.AccessRemoved, daResult.AccessUpdated))
		}

		return nil
	}

	return fmt.Errorf("failed to load results")
}

func (s *DataAccessSync) GetResultObject() interface{} {
	return &AccessProviderImportResult{}
}
