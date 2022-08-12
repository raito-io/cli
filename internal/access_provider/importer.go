package access_provider

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/file"
	"github.com/raito-io/cli/internal/graphql"
	"github.com/raito-io/cli/internal/target"
	"github.com/spf13/viper"
)

type AccessProviderImportConfig struct {
	target.BaseTargetConfig
	TargetFile      string
	DeleteUntouched bool
}

type AccessProviderImportResult struct {
	AccessAdded   int `json:"accessAdded"`
	AccessUpdated int `json:"accessUpdated"`
	AccessRemoved int `json:"accessRemoved"`
}

type AccessProviderImporter interface {
	TriggerImport() (*AccessProviderImportResult, error)
}

type accessProviderImporter struct {
	config *AccessProviderImportConfig
	log    hclog.Logger
}

func NewAccessProviderImporter(config *AccessProviderImportConfig) AccessProviderImporter {
	logger := config.Logger.With("AccessProvider", config.DataSourceId, "file", config.TargetFile)
	dsI := accessProviderImporter{config, logger}

	return &dsI
}

func (d *accessProviderImporter) TriggerImport() (*AccessProviderImportResult, error) {
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

func (d *accessProviderImporter) upload() (string, error) {
	key, err := file.UploadFile(d.config.TargetFile, &d.config.BaseTargetConfig)
	if err != nil {
		return "", fmt.Errorf("error while uploading data source import files to Raito: %s", err.Error())
	}

	return key, nil
}

func (d *accessProviderImporter) doImport(fileKey string) (*AccessProviderImportResult, error) {
	start := time.Now()

	gqlQuery := fmt.Sprintf(`{ "operationName": "ImportAccessProviders", "variables":{}, "query": "mutation ImportAccessProviders {
        importAccessProviders(input: {
          dataSource: \"%s\",
          deleteUntouched: %t,
          fileKey: \"%s\"
        }) {
          accessAdded
          accessUpdated
          accessRemoved
          errors
        }
    }" }"`, d.config.DataSourceId, d.config.DeleteUntouched, fileKey)

	gqlQuery = strings.Replace(gqlQuery, "\n", "\\n", -1)

	res := Response{}
	_, err := graphql.ExecuteGraphQL(gqlQuery, &d.config.BaseTargetConfig, &res)

	if err != nil {
		return nil, fmt.Errorf("error while executing import: %s", err.Error())
	}

	ret := &res.ImportAccessProviders

	d.log.Info(fmt.Sprintf("Done executing import in %s", time.Since(start).Round(time.Millisecond)))

	return ret, nil
}

type Response struct {
	ImportAccessProviders AccessProviderImportResult `json:"importAccessProviders"`
}
