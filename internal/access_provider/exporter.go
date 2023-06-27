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
	TriggerExport(ctx context.Context, jobId string) (job.JobStatus, string, error)
}

type accessProviderExporter struct {
	config        *AccessProviderExporterConfig
	log           hclog.Logger
	statusUpdater job.TaskEventUpdater
	syncConfig    *access_provider.AccessSyncConfig
}

func NewAccessProviderExporter(config *AccessProviderExporterConfig, statusUpdater job.TaskEventUpdater, syncConfig *access_provider.AccessSyncConfig) AccessProviderExporter {
	logger := config.TargetLogger.With("AccessProviderExport", config.DataSourceId)
	dsI := accessProviderExporter{config, logger, statusUpdater, syncConfig}

	return &dsI
}

func (d *accessProviderExporter) TriggerExport(ctx context.Context, jobId string) (job.JobStatus, string, error) {
	status, subTaskId, err := d.doExport(ctx, jobId)

	if err != nil {
		return job.Failed, "", err
	}

	result := &AccessProviderExportResult{}
	subtask, err := job.WaitForJobToComplete(ctx, jobId, constants.DataAccessSync, subTaskId, result, &d.config.BaseTargetConfig, status)

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

func (d *accessProviderExporter) doExport(ctx context.Context, jobId string) (job.JobStatus, string, error) {
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

	filter := AccessProviderExportRequestFilter{
		Status: exportFilterStatus{
			OutOfSync: d.config.OnlyOutOfSyncData && d.syncConfig.SupportPartialSync,
		},
	}

	if len(d.syncConfig.RequiredExportWhoList) > 0 {
		if filter.ExportFilterProperties == nil {
			filter.ExportFilterProperties = &exportFilterProperties{}
		}

		filter.ExportFilterProperties.RequiredWhoItemLists = make([]string, len(d.syncConfig.RequiredExportWhoList))

		for _, listType := range d.syncConfig.RequiredExportWhoList {
			switch listType { //nolint:exhaustive
			case access_provider.AccessProviderExportWhoList_ACCESSPROVIDER_EXPORT_WHO_LIST_USERS:
				filter.ExportFilterProperties.RequiredWhoItemLists = append(filter.ExportFilterProperties.RequiredWhoItemLists, "users")
			case access_provider.AccessProviderExportWhoList_ACCESSPROVIDER_EXPORT_WHO_LIST_GROUPS:
				filter.ExportFilterProperties.RequiredWhoItemLists = append(filter.ExportFilterProperties.RequiredWhoItemLists, "groups")
			case access_provider.AccessProviderExportWhoList_ACCESSPROVIDER_EXPORT_WHO_LIST_INHERIT_FROM:
				filter.ExportFilterProperties.RequiredWhoItemLists = append(filter.ExportFilterProperties.RequiredWhoItemLists, "inheritFrom")
			case access_provider.AccessProviderExportWhoList_ACCESSPROVIDER_EXPORT_WHO_LIST_USERS_IN_GROUPS:
				filter.ExportFilterProperties.RequiredWhoItemLists = append(filter.ExportFilterProperties.RequiredWhoItemLists, "usersInGroups")
			case access_provider.AccessProviderExportWhoList_ACCESSPROVIDER_EXPORT_WHO_LIST_USERS_INHERITED:
				filter.ExportFilterProperties.RequiredWhoItemLists = append(filter.ExportFilterProperties.RequiredWhoItemLists, "usersInherited")
			case access_provider.AccessProviderExportWhoList_ACCESSPROVIDER_EXPORT_WHO_LIST_NATIVE_GROUPS_INHERITED:
				filter.ExportFilterProperties.RequiredWhoItemLists = append(filter.ExportFilterProperties.RequiredWhoItemLists, "groupsInherited")
			case access_provider.AccessProviderExportWhoList_ACCESSPROVIDER_EXPORT_WHO_LIST_USERS_INHERITED_NATIVE_GROUPS_EXCLUDED:
				filter.ExportFilterProperties.RequiredWhoItemLists = append(filter.ExportFilterProperties.RequiredWhoItemLists, "usersInGroupsExclude")
			default:
				return 0, "", fmt.Errorf("unknown who list type: %s", listType)
			}
		}
	}

	client := graphql.NewClient(&d.config.BaseConfig)

	err := client.Query(ctx, &q, map[string]interface{}{"input": input, "filter": filter})

	if err != nil {
		return job.Failed, "", fmt.Errorf("error while executing export: %s", err.Error())
	}

	retStatus := q.ExportAccessProvidersRequest.Subtask.Status
	subtaskId := q.ExportAccessProvidersRequest.Subtask.SubtaskId

	d.log.Info(fmt.Sprintf("Done submitting export in %s", time.Since(start).Round(time.Millisecond)))

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
	Status                 exportFilterStatus      `json:"status"`
	ExportFilterProperties *exportFilterProperties `json:"exportProperties,omitempty"`
}
