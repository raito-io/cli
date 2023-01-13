package access_provider

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"

	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/file"
	"github.com/raito-io/cli/internal/graphql"
	"github.com/raito-io/cli/internal/job"
	"github.com/raito-io/cli/internal/target"
)

type AccessProviderExportResult struct {
	FileKey      string   `json:"fileKey"`
	FileLocation string   `json:"fileLocation"`
	Warnings     []string `json:"warnings"`
}

type AccessProviderExporterConfig struct {
	target.BaseTargetConfig
}

type AccessProviderExporter interface {
	TriggerExport(jobId string) (job.JobStatus, string, error)
}

type accessProviderExporter struct {
	config        *AccessProviderExporterConfig
	log           hclog.Logger
	statusUpdater job.TaskEventUpdater
}

func NewAccessProviderExporter(config *AccessProviderExporterConfig, statusUpdater job.TaskEventUpdater) AccessProviderExporter {
	logger := config.TargetLogger.With("AccessProviderExport", config.DataSourceId)
	dsI := accessProviderExporter{config, logger, statusUpdater}

	return &dsI
}

func (d *accessProviderExporter) TriggerExport(jobId string) (job.JobStatus, string, error) {
	status, subTaskId, err := d.doExport(jobId)

	if err != nil {
		return job.Failed, "", err
	}

	result := &AccessProviderExportResult{}
	_, err = job.WaitForJobToComplete(jobId, constants.DataAccessSync, subTaskId, result, &d.config.BaseTargetConfig, status)

	if err != nil {
		return job.Failed, "", err
	}

	result.FileLocation, err = d.download(result.FileLocation)

	if err != nil {
		return job.Failed, "", err
	}

	return job.Completed, result.FileLocation, nil
}

func (d *accessProviderExporter) download(url string) (string, error) {
	cn := strings.Replace(d.config.ConnectorName, "/", "-", -1)
	filePath, err := filepath.Abs(file.CreateUniqueFileName(cn+"-as", "json"))

	if err != nil {
		return "", err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("error while creating temporary file %q: %s", filePath, err.Error())
	}
	defer file.Close()

	resp, err := http.Get(url) //nolint

	if err != nil {
		return "", fmt.Errorf("error while fetching access controls for datasource %q: %s", d.config.DataSourceId, err.Error())
	}

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("error (HTTP %d) while fetching access controls for datasource %q: %s", resp.StatusCode, d.config.DataSourceId, resp.Status)
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error while reading bytes from export file: %s", err.Error())
	}

	_, err = file.Write(bytes)
	if err != nil {
		return "", fmt.Errorf("error while writing data to file: %s", err.Error())
	}

	return filePath, nil
}

func (d *accessProviderExporter) doExport(jobId string) (job.JobStatus, string, error) {
	start := time.Now()

	filter := ""

	if d.config.ModifiedAfter != nil {
		filter = fmt.Sprintf(`, filter : {
					modifiedAfter: \"%s\"
			    }`, d.config.ModifiedAfter.Format(time.RFC3339))
	}

	gqlQuery := fmt.Sprintf(`{ "operationName": "ExportAccessProvidersRequest", "variables":{}, "query": "query ExportAccessProvidersRequest {
        exportAccessProvidersRequest(input: {
          jobId: \"%s\",
          exportSettings: {
            dataSource: \"%s\"
          }
        }%s) {
          subtask {
            subtaskId
            status            
          }
         }
    }" }"`, jobId, d.config.DataSourceId, filter)

	gqlQuery = strings.Replace(gqlQuery, "\n", "\\n", -1)

	res := exportResponse{}
	_, err := graphql.ExecuteGraphQL(gqlQuery, &d.config.BaseConfig, &res)

	if err != nil {
		return job.Failed, "", fmt.Errorf("error while executing export: %s", err.Error())
	}

	retStatus := res.Response.Subtask.Status
	subtaskId := res.Response.Subtask.SubtaskId

	d.log.Info(fmt.Sprintf("Done submitting export in %s", time.Since(start).Round(time.Millisecond)))

	return retStatus, subtaskId, nil
}

type exportResponse struct {
	Response QueryResponse `json:"exportAccessProvidersRequest"`
}
