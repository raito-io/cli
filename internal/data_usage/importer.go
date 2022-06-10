package data_usage

import (
	"encoding/json"
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
	StatementsAdded  int             `json:"statementsAdded"`
	StatementsFailed int             `json:"statementsFailed"`
	Errors           []graphql.Error `json:"errors"`
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

	gqlQuery := fmt.Sprintf(`{ "operationName": "ImportDataUsage", "variables":{}, "query": "mutation ImportDataUsage {
        importDataUsage(input: {
          dataSource: \"%s\",
          dataObjects: \"%s\"
        }) {
          statementsAdded
          statementsFailed
          errors
        }
    }" }"`, d.config.DataSourceId, fileKey)
	gqlQuery = strings.Replace(gqlQuery, "\n", "\\n", -1)

	res, err := graphql.ExecuteGraphQL(gqlQuery, &d.config.BaseTargetConfig)
	if err != nil {
		return nil, fmt.Errorf("error while executing data usage import on appserver: %s", err.Error())
	}

	ret, err := d.parseImportResult(res)
	if err != nil {
		return nil, err
	}
	if len(ret.Errors) > 0 {
		return ret, fmt.Errorf("errors while importing/processing data usage: %s", ret.Errors[0].Message)
	}

	d.log.Info(fmt.Sprintf("Successfully imported %d data usage statements, %d failures, in %s", ret.StatementsAdded, ret.StatementsFailed, time.Since(start).Round(time.Millisecond)))

	return &DataUsageImportResult{StatementsAdded: ret.StatementsAdded, StatementsFailed: ret.StatementsFailed}, nil
}

func (d *dataUsageImporter) parseImportResult(res []byte) (*DataUsageImportResult, error) {
	resp := Response{}
	gr := graphql.GraphqlResponse{Data: &resp}
	err := json.Unmarshal(res, &gr)
	if err != nil {
		return nil, fmt.Errorf("error while parsing data usage import result: %s", err.Error())
	}

	// Flatten the result
	resp.ImportDataUsage.Errors = gr.Errors

	return &(resp.ImportDataUsage), nil
}

type Response struct {
	ImportDataUsage DataUsageImportResult `json:"importDataUsage"`
}
