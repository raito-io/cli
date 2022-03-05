package data_source

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/raito-io/cli/internal/file"
	"github.com/raito-io/cli/internal/graphql"
	"github.com/raito-io/cli/internal/target"
	"strings"
	"time"
)

type DataSourceImportConfig struct {
	target.BaseTargetConfig
	TargetFile string
	DeleteUntouched bool
	ReplaceTags bool
}

type DataSourceImportResult struct {
	DataObjectsAdded    int     `json:"dataObjectsAdded"`
	DataObjectsUpdated  int     `json:"dataObjectsUpdated"`
	DataObjectsRemoved  int             `json:"dataObjectsRemoved"`
	Errors              []graphql.Error `json:"_"`
}

type DataSourceImporter interface {
	TriggerImport() (*DataSourceImportResult, error)
}

type dataSourceImporter struct {
	config *DataSourceImportConfig
	log hclog.Logger
}

func NewDataSourceImporter(config *DataSourceImportConfig) DataSourceImporter {
	logger := config.Logger.With("datasource", config.DataSourceId, "file", config.TargetFile)
	dsI := dataSourceImporter{config, logger	}
	return &dsI
}

func (d *dataSourceImporter) TriggerImport() (*DataSourceImportResult, error) {
	key, err := d.upload()
	if err != nil {
		return nil, err
	}

	return d.doImport(key)
}

func (d *dataSourceImporter) upload() (string, error) {
	key, err := file.UploadFile(d.config.TargetFile, &d.config.BaseTargetConfig)
	if err != nil {
		return "", fmt.Errorf("error while uploading data source import files to Raito: %s", err.Error())
	}
	return key, nil
}

func (d *dataSourceImporter) doImport(fileKey string) (*DataSourceImportResult, error) {
	start := time.Now()

	gqlQuery := fmt.Sprintf(`{ "operationName": "ImportDataSource", "variables":{}, "query": "mutation ImportDataSource {
        importDataSource(input: {
          dataSource: \"%s\",
          deleteUntouched: %t,
          replaceTags: %t,  
          dataObjects: \"%s\"
        }) {
          dataObjectsAdded
          dataObjectsUpdated
          dataObjectsRemoved
          errors
        }
    }" }"`, d.config.DataSourceId, d.config.DeleteUntouched, d.config.ReplaceTags, fileKey)

	gqlQuery = strings.Replace(gqlQuery, "\n", "\\n", -1)

	res, err := graphql.ExecuteGraphQL(gqlQuery, &d.config.BaseTargetConfig)
	if err != nil {
		return nil, fmt.Errorf("error while executing import: %s", err.Error())
	}

	ret, err := d.parseImportResult(res)
	if err != nil {
		return nil, err
	}
	if len(ret.Errors) > 0 {
		return ret, fmt.Errorf("errors while importing into data source: %s", ret.Errors[0].Message)
	}

	d.log.Info(fmt.Sprintf("Done executing import in %s", time.Since(start).Round(time.Millisecond)))

	return ret, nil
}

func (d *dataSourceImporter) parseImportResult(res []byte) (*DataSourceImportResult, error) {
	resp := Response{}
	gr := graphql.GraphqlResponse{ Data: &resp }
	err := json.Unmarshal(res, &gr)

	if err != nil {
		return nil, fmt.Errorf("error while parsing data source import result: %s", err.Error())
	}

	// Flatten the result
	resp.ImportDataSource.Errors = gr.Errors

	return &(resp.ImportDataSource), nil
}

type Response struct {
	ImportDataSource DataSourceImportResult `json:"importDataSource"`
}
