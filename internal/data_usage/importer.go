package data_usage

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

type DataUsageImportConfig struct {
	target.BaseTargetConfig
	TargetFile string
}

type DataUsageImportResult struct {
	StatementsAdded       int `json:"statementsAdded"`
	StatementsFailed      int `json:"statementsFailed"`
	StatementsSkipped     int `json:"statementsSkipped"`
	EdgesCreatedOrUpdated int `json:"edgesCreatedOrUpdated"`
	EdgesRemoved          int `json:"edgesRemoved"`
}

type DataSourceLastUsed struct {
	Id       string `json:"id"`
	LastUsed string `json:"usageLastUsed"`
}

type DataUsageImporter interface {
	TriggerImport() (*DataUsageImportResult, error)
	GetLastUsage() (*time.Time, error)
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

	gqlQuery := fmt.Sprintf(`{ "operationName": "ImportDataUsage", "variables":{}, "query": "mutation ImportDataUsage {
        importDataUsage(input: {
          dataSource: \"%s\",
          fileKey: \"%s\"
        }) {
          statementsAdded
          statementsFailed
          errors
        }
    }" }"`, d.config.DataSourceId, fileKey)
	gqlQuery = strings.Replace(gqlQuery, "\n", "\\n", -1)

	res := Response{}
	_, err := graphql.ExecuteGraphQL(gqlQuery, &d.config.BaseTargetConfig, &res)

	if err != nil {
		return nil, fmt.Errorf("error while executing data usage import on appserver: %s", err.Error())
	}

	ret := &res.ImportDataUsage

	d.log.Info(fmt.Sprintf("Successfully imported %d data usage statements, %d failures, %d skipped; %d relationships created/updated, %d relationships deleted; in %s",
		ret.StatementsAdded, ret.StatementsFailed, ret.StatementsSkipped, ret.EdgesCreatedOrUpdated, ret.EdgesRemoved, time.Since(start).Round(time.Millisecond)))

	return &DataUsageImportResult{StatementsAdded: ret.StatementsAdded, StatementsFailed: ret.StatementsFailed}, nil
}

func (d *dataUsageImporter) GetLastUsage() (*time.Time, error) {
	gqlQuery := fmt.Sprintf(`{"variables":{}, "query": "query {dataSource(id:\"%s\") {id usageLastUsed}}" }`, d.config.DataSourceId)
	gqlQuery = strings.Replace(gqlQuery, "\n", "\\n", -1)
	res := LastUsedReponse{}
	_, err := graphql.ExecuteGraphQL(gqlQuery, &d.config.BaseTargetConfig, &res)

	if err != nil {
		return nil, fmt.Errorf("error while executing data usage import on appserver: %s", err.Error())
	}

	finalResult := time.Unix(int64(0), 0)
	if res.DataSourceInfo.LastUsed != "" {
		finalResultRaw, err := time.Parse(time.RFC3339, res.DataSourceInfo.LastUsed)
		if err == nil {
			finalResult = finalResultRaw
		}
	}

	return &finalResult, nil
}

type Response struct {
	ImportDataUsage DataUsageImportResult `json:"importDataUsage"`
}

type LastUsedReponse struct {
	DataSourceInfo DataSourceLastUsed `json:"dataSource"`
}
