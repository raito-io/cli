package data_usage

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"

	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/file"
	"github.com/raito-io/cli/internal/graphql"
	"github.com/raito-io/cli/internal/job"
	"github.com/raito-io/cli/internal/target"
)

type DataUsageImportConfig struct {
	target.BaseTargetConfig
	TargetFile string
}

type DataSourceLastUsed struct {
	Id       string `json:"id"`
	LastUsed string `json:"usageLastUsed"`
}

type DataUsageImporter interface {
	TriggerImport(jobId string) (job.JobStatus, string, error)
	GetLastUsage() (*time.Time, error)
}

type dataUsageImporter struct {
	config        *DataUsageImportConfig
	log           hclog.Logger
	statusUpdater func(status job.JobStatus)
}

func NewDataUsageImporter(config *DataUsageImportConfig, statusUpdater func(status job.JobStatus)) DataUsageImporter {
	logger := config.Logger.With("data-usage", config.DataSourceId, "file", config.TargetFile)
	duI := dataUsageImporter{config, logger, statusUpdater}

	return &duI
}

func (d *dataUsageImporter) TriggerImport(jobId string) (job.JobStatus, string, error) {
	env := viper.GetString(constants.EnvironmentFlag)
	if env == constants.EnvironmentDev {
		// In the development environment, we skip the upload and use the local file for the import
		return d.doImport(jobId, d.config.TargetFile)
	} else {
		key, err := d.upload()
		if err != nil {
			return job.Failed, "", err
		}

		return d.doImport(jobId, key)
	}
}

func (d *dataUsageImporter) upload() (string, error) {
	d.statusUpdater(job.DataUpload)
	key, err := file.UploadFile(d.config.TargetFile, &d.config.BaseTargetConfig)

	if err != nil {
		return "", fmt.Errorf("error while uploading data usage import files to Raito: %s", err.Error())
	}

	return key, nil
}

func (d *dataUsageImporter) doImport(jobId string, fileKey string) (job.JobStatus, string, error) {
	gqlQuery := fmt.Sprintf(`{ "operationName": "ImportDataUsageRequest", "variables":{}, "query": "mutation ImportDataUsageRequest {
      importDataUsageRequest(input: {
        jobId: \"%s\",
        importSettings: {
          dataSource: \"%s\",
          fileKey: \"%s\"
        }
      }) {
        subtask {
            subTask
            status            
          }
      }
    }" }"`, jobId, d.config.DataSourceId, fileKey)
	gqlQuery = strings.Replace(gqlQuery, "\n", "\\n", -1)

	res := Response{}
	_, err := graphql.ExecuteGraphQL(gqlQuery, &d.config.BaseTargetConfig, &res)

	if err != nil {
		return job.Failed, "", fmt.Errorf("error while executing data usage import on appserver: %s", err.Error())
	}

	return res.Response.Subtask.Status, res.Response.Subtask.Subtask, nil
}

func (d *dataUsageImporter) GetLastUsage() (*time.Time, error) {
	gqlQuery := fmt.Sprintf(`{"variables":{}, "query": "query {dataSource(id:\"%s\") {id usageLastUsed}}" }`, d.config.DataSourceId)
	gqlQuery = strings.Replace(gqlQuery, "\n", "\\n", -1)
	res := LastUsedResponse{}
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

type subtaskResponse struct {
	Status  job.JobStatus `json:"status"`
	Subtask string        `json:"subTask"`
}

type QueryResponse struct {
	Subtask subtaskResponse `json:"subtask"`
}

type Response struct {
	Response QueryResponse `json:"importDataUsageRequest"`
}

type LastUsedResponse struct {
	DataSourceInfo DataSourceLastUsed `json:"dataSource"`
}
