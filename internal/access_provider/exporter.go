package access_provider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"

	"github.com/raito-io/cli/base/access_provider"
	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/graphql"
	"github.com/raito-io/cli/internal/job"
	"github.com/raito-io/cli/internal/target/types"
	"github.com/raito-io/cli/internal/util/file"
)

type AccessProviderExportResult struct {
	FileKey      string   `json:"fileKey"`
	FileLocation string   `json:"fileLocation"`
	Warnings     []string `json:"warnings"`
}

type AccessProviderExporterConfig struct {
	types.BaseTargetConfig
}

type AccessProviderExporter interface {
	TriggerExport(ctx context.Context, logger hclog.Logger, jobId string) (job.JobStatus, string, error)
}

type accessProviderExporter struct {
	config        *AccessProviderExporterConfig
	statusUpdater job.TaskEventUpdater
	syncConfig    *access_provider.AccessSyncConfig
}

func NewAccessProviderExporter(config *AccessProviderExporterConfig, statusUpdater job.TaskEventUpdater, syncConfig *access_provider.AccessSyncConfig) AccessProviderExporter {
	dsI := accessProviderExporter{config, statusUpdater, syncConfig}

	return &dsI
}

func (d *accessProviderExporter) TriggerExport(ctx context.Context, logger hclog.Logger, jobId string) (job.JobStatus, string, error) {
	logger = logger.With("AccessProviderExport", d.config.DataSourceId)

	status, subTaskId, err := d.doExport(ctx, logger, jobId)

	if err != nil {
		return job.Failed, "", err
	}

	result := &AccessProviderExportResult{}
	subtask, err := job.WaitForJobToComplete(ctx, logger, jobId, constants.DataAccessSync, subTaskId, result, &d.config.BaseTargetConfig, status)

	if err != nil {
		return job.Failed, "", err
	}

	if subtask.Status == job.Failed {
		return job.Failed, "", fmt.Errorf("export failed on server: [%s]", strings.Join(subtask.Errors, ", "))
	} else if subtask.Status == job.TimeOut {
		return job.TimeOut, "", fmt.Errorf("export timeout on server")
	}

	result.FileLocation, err = d.download(result.FileLocation)

	if err != nil {
		return job.Failed, "", err
	}

	return job.Completed, result.FileLocation, nil
}

func (d *accessProviderExporter) download(url string) (string, error) {
	filePath, err := filepath.Abs(file.CreateUniqueFileNameForTarget(d.config.Name, "toTarget-access", "yaml"))

	if err != nil {
		return "", err
	}

	downloadedFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("error while creating temporary file %q: %s", filePath, err.Error())
	}

	defer downloadedFile.Close()

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

	_, err = downloadedFile.Write(bytes)
	if err != nil {
		return "", fmt.Errorf("error while writing data to file: %s", err.Error())
	}

	return filePath, nil
}

func (d *accessProviderExporter) doExport(ctx context.Context, logger hclog.Logger, jobId string) (job.JobStatus, string, error) {
	start := time.Now()

	var q struct {
		ExportAccessProvidersRequest struct {
			Subtask struct {
				SubtaskId string
				Status    job.JobStatus
			}
		} `graphql:"exportAccessProvidersRequest(input: $input, filter: $filter)"`
	}

	input := AccessProviderExportRequest{
		JobId: jobId,
		ExportSettings: exportSettings{
			DataSource: d.config.DataSourceId,
		},
	}

	filter := AccessProviderExportRequestFilter{}

	if d.config.OnlyOutOfSyncData && d.syncConfig.SupportPartialSync {
		filter.Status = &exportFilterStatus{
			OutOfSync: d.config.OnlyOutOfSyncData && d.syncConfig.SupportPartialSync,
		}
	}

	client := graphql.NewClient(&d.config.BaseConfig)

	err := client.Query(ctx, &q, map[string]interface{}{"input": input, "filter": filter})

	if err != nil {
		return job.Failed, "", fmt.Errorf("error while executing export: %s", err.Error())
	}

	retStatus := q.ExportAccessProvidersRequest.Subtask.Status
	subtaskId := q.ExportAccessProvidersRequest.Subtask.SubtaskId

	logger.Info(fmt.Sprintf("Done submitting export in %s", time.Since(start).Round(time.Millisecond)))

	return retStatus, subtaskId, nil
}

type exportSettings struct {
	DataSource string `json:"dataSource"`
}

type AccessProviderExportRequest struct {
	JobId          string         `json:"jobId"`
	ExportSettings exportSettings `json:"exportSettings"`
}

type exportFilterStatus struct {
	OutOfSync bool `json:"outOfSync"`
}

type exportFilterProperties struct {
	RequiredWhoItemLists []string `json:"requiredWhoItemLists"`
}

type AccessProviderExportRequestFilter struct {
	Status                 *exportFilterStatus     `json:"status,omitempty"`
	ExportFilterProperties *exportFilterProperties `json:"exportProperties,omitempty"`
}
