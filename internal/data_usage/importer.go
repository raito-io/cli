package data_usage

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/file"
	"github.com/raito-io/cli/internal/graphql"
	"github.com/raito-io/cli/internal/target"
	"github.com/spf13/viper"
	"time"
)

type DataUsageImportConfig struct {
	target.BaseTargetConfig
	TargetFile string
}

type DataUsageImportResult struct {
	TransactionAdded int             `json:"transactionsAdded"`
	Errors           []graphql.Error `json:"_"`
}

type DataUsageImporter interface {
	TriggerImport() (*DataUsageImportResult, error)
}

type dataUsageImporter struct {
	config *DataUsageImportConfig
	log    hclog.Logger
}

func NewDataUsageImporter(config *DataUsageImportConfig) DataUsageImporter {
	logger := config.Logger.With("data-usage", config.DataSourceId, "file", config.TargetFile)
	duI := dataUsageImporter{config, logger}
	return &duI
}

func (d *dataUsageImporter) TriggerImport() (*DataUsageImportResult, error) {
	env := viper.GetString(constants.EnvironmentFlag)
	if env == constants.EnvironmentDev {
		// In the development environment, we skip the upload and use the local file for the import
		return d.doImport(d.config.TargetFile)
	} else {
		key, err := d.upload()
		if err != nil {
			return nil, err
		}

		return d.doImport(key)
	}
}

func (d *dataUsageImporter) upload() (string, error) {
	key, err := file.UploadFile(d.config.TargetFile, &d.config.BaseTargetConfig)
	if err != nil {
		return "", fmt.Errorf("error while uploading data usage import files to Raito: %s", err.Error())
	}
	return key, nil
}

func (d *dataUsageImporter) doImport(fileKey string) (*DataUsageImportResult, error) {
	start := time.Now()

	d.log.Info(fmt.Sprintf("Actual import to appserver still needs to be implemented, reading from file %s", fileKey))
	d.log.Info(fmt.Sprintf("Done executing import in %s", time.Since(start).Round(time.Millisecond)))

	return &DataUsageImportResult{}, nil
}
