package job

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/raito-io/cli/internal/graphql"
	"github.com/raito-io/cli/internal/logging"
	"github.com/raito-io/cli/internal/plugin"
	"github.com/raito-io/cli/internal/target/types"
)

const NrOfRetriesOnError = 5

//go:generate go run github.com/vektra/mockery/v2 --name=TaskEventUpdater --with-expecter
type TaskEventUpdater interface {
	SetStatusToStarted(ctx context.Context)
	SetStatusToDataRetrieve(ctx context.Context)
	SetStatusToDataUpload(ctx context.Context)
	SetStatusToQueued(ctx context.Context)
	SetStatusToDataProcessing(ctx context.Context)
	SetStatusToCompleted(ctx context.Context, results []TaskResult)
	SetStatusToFailed(ctx context.Context, err error)
	SetStatusToSkipped(ctx context.Context)

	GetSubtaskEventUpdater(subtask string) SubtaskEventUpdater
}

type SubtaskEventUpdater interface {
	AddSubtaskEvent(ctx context.Context, status JobStatus)
	SetReceivedDate(receivedDate int64)
}

type Task interface {
	IsClientValid(ctx context.Context, c plugin.PluginClient) (bool, error)
	GetParts() []TaskPart
	GetTaskResults() []TaskResult
}

type TaskPart interface {
	StartSyncAndQueueTaskPart(ctx context.Context, c plugin.PluginClient, statusUpdater TaskEventUpdater) (JobStatus, string, error)
	ProcessResults(results interface{}) error
	GetResultObject() interface{}
}

type TaskResult struct {
	ObjectType string `json:"objectType"`
	Added      int    `json:"added"`
	Updated    int    `json:"updated"`
	Removed    int    `json:"removed"`
	Failed     int    `json:"failed"`
}

type taskEventUpdater struct {
	Cfg              *types.BaseTargetConfig
	JobId            string
	JobType          string
	warningCollector logging.WarningCollector
}

func NewTaskEventUpdater(cfg *types.BaseTargetConfig, jobId, jobType string, warningCollector logging.WarningCollector) TaskEventUpdater {
	return &taskEventUpdater{cfg, jobId, jobType, warningCollector}
}

func (u *taskEventUpdater) setStatus(ctx context.Context, status JobStatus, results []TaskResult, err error) {
	var errors []error
	if err != nil {
		errors = append(errors, err)
	}

	var warnings []string
	if u.warningCollector != nil {
		warnings = u.warningCollector.GetWarnings()
	}

	AddTaskEvent(ctx, u.Cfg, u.JobId, u.JobType, status, results, warnings, errors)
}

func (u *taskEventUpdater) SetStatusToStarted(ctx context.Context) {
	u.setStatus(ctx, Started, nil, nil)
}

func (u *taskEventUpdater) SetStatusToDataRetrieve(ctx context.Context) {
	u.setStatus(ctx, DataRetrieve, nil, nil)
}

func (u *taskEventUpdater) SetStatusToDataUpload(ctx context.Context) {
	u.setStatus(ctx, DataUpload, nil, nil)
}

func (u *taskEventUpdater) SetStatusToQueued(ctx context.Context) {
	u.setStatus(ctx, Queued, nil, nil)
}

func (u *taskEventUpdater) SetStatusToDataProcessing(ctx context.Context) {
	u.setStatus(ctx, DataProcessing, nil, nil)
}

func (u *taskEventUpdater) SetStatusToCompleted(ctx context.Context, results []TaskResult) {
	u.setStatus(ctx, Completed, results, nil)
}

func (u *taskEventUpdater) SetStatusToFailed(ctx context.Context, err error) {
	u.setStatus(ctx, Failed, nil, err)
}

func (u *taskEventUpdater) SetStatusToSkipped(ctx context.Context) {
	u.setStatus(ctx, Skipped, nil, nil)
}

func (u *taskEventUpdater) GetSubtaskEventUpdater(subtask string) SubtaskEventUpdater {
	return &subtaskEventUpdater{
		Cfg:     u.Cfg,
		JobId:   u.JobId,
		JobType: u.JobType,
		Subtask: subtask,
	}
}

type subtaskEventUpdater struct {
	Cfg          *types.BaseTargetConfig
	JobId        string
	JobType      string
	Subtask      string
	receivedDate *int64
}

func (u *subtaskEventUpdater) AddSubtaskEvent(ctx context.Context, status JobStatus) {
	AddSubtaskEvent(ctx, u.Cfg, u.JobId, u.JobType, u.Subtask, status, u.receivedDate)
}

func (u *subtaskEventUpdater) SetReceivedDate(receivedDate int64) {
	u.receivedDate = &receivedDate
}

func StartJob(ctx context.Context, cfg *types.BaseTargetConfig) (string, error) {
	var mutation struct {
		CreateJob struct {
			JobId string
		} `graphql:"createJob(input: $input)"`
	}

	type JobInput struct {
		DataSourceId    *string   `json:"dataSourceId"`
		IdentityStoreId *string   `json:"identityStoreId"`
		EventTime       time.Time `json:"eventTime"`
	}

	input := JobInput{
		EventTime: time.Now(),
	}

	if cfg.DataSourceId != "" {
		input.DataSourceId = &cfg.DataSourceId
	}

	if cfg.IdentityStoreId != "" {
		input.IdentityStoreId = &cfg.IdentityStoreId
	}

	err := graphql.NewClient(&cfg.BaseConfig).Mutate(ctx, &mutation, map[string]interface{}{"input": input})
	if err != nil {
		return "", fmt.Errorf("error while executing import: %s", err.Error())
	}

	return mutation.CreateJob.JobId, nil
}

func UpdateJobEvent(cfg *types.BaseTargetConfig, jobID string, status JobStatus, inputErr error) {
	var mutation struct {
		UpdateJob struct {
			JobId string
		} `graphql:"updateJob(id: $id, input: $input)"`
	}

	type JobInput struct {
		DataSourceId    *string   `json:"dataSourceId"`
		IdentityStoreId *string   `json:"identityStoreId"`
		Status          JobStatus `json:"status"`
		EventTime       time.Time `json:"eventTime"`
		Errors          []string  `json:"errors"`
	}

	input := JobInput{
		Status:    status,
		EventTime: time.Now(),
	}

	if inputErr != nil {
		input.Errors = append(input.Errors, inputErr.Error())
	}

	if cfg.DataSourceId != "" {
		input.DataSourceId = &cfg.DataSourceId
	}

	if cfg.IdentityStoreId != "" {
		input.IdentityStoreId = &cfg.IdentityStoreId
	}

	err := graphql.NewClient(&cfg.BaseConfig).Mutate(context.Background(), &mutation, map[string]interface{}{"id": jobID, "input": input})
	if err != nil {
		cfg.TargetLogger.Debug(fmt.Sprintf("job update failed: %s", err.Error()))
	}
}

func AddTaskEvent(ctx context.Context, cfg *types.BaseTargetConfig, jobID, jobType string, status JobStatus, taskResults []TaskResult, warnings []string, errors []error) {
	var mutation struct {
		AddTaskEvent struct {
			JobId string
		} `graphql:"addTaskEvent(input: $input)"`
	}

	type TaskEventInput struct {
		JobId           string       `json:"jobId"`
		JobType         string       `json:"jobType"`
		DataSourceId    *string      `json:"dataSourceId"`
		IdentityStoreId *string      `json:"identityStoreId"`
		Status          JobStatus    `json:"status"`
		EventTime       time.Time    `json:"eventTime"`
		Errors          []string     `json:"errors"`
		Warnings        []string     `json:"warnings"`
		Result          []TaskResult `json:"result"`
	}

	var errorMsgs []string
	if len(errors) > 0 {
		errorMsgs = make([]string, len(errors))
		for i, err := range errors {
			errorMsgs[i] = err.Error()
		}
	}

	input := TaskEventInput{
		JobId:     jobID,
		JobType:   jobType,
		EventTime: time.Now(),
		Status:    status,
		Warnings:  warnings,
		Errors:    errorMsgs,
		Result:    taskResults,
	}

	if cfg.DataSourceId != "" {
		input.DataSourceId = &cfg.DataSourceId
	}

	if cfg.IdentityStoreId != "" {
		input.IdentityStoreId = &cfg.IdentityStoreId
	}

	err := graphql.NewClient(&cfg.BaseConfig).Mutate(ctx, &mutation, map[string]interface{}{"input": input})
	if err != nil {
		cfg.TargetLogger.Debug(fmt.Sprintf("taskEvent update failed: %s", err.Error()))
	}
}

func AddSubtaskEvent(ctx context.Context, cfg *types.BaseTargetConfig, jobID, jobType, subtask string, status JobStatus, receivedDate *int64) {
	var mutation struct {
		AddSubtaskEvent struct {
			JobId string
		} `graphql:"addSubtaskEvent(input: $input)"`
	}

	type SubtaskInput struct {
		JobId           string    `json:"jobId"`
		JobType         string    `json:"jobType"`
		SubtaskId       string    `json:"subtaskId"`
		DataSourceId    *string   `json:"dataSourceId"`
		IdentityStoreId *string   `json:"identityStoreId"`
		Status          JobStatus `json:"status"`
		EventTime       time.Time `json:"eventTime"`
		ReceivedDate    *int64    `json:"receivedDate"`
	}

	input := SubtaskInput{
		JobId:        jobID,
		JobType:      jobType,
		SubtaskId:    subtask,
		Status:       status,
		EventTime:    time.Now(),
		ReceivedDate: receivedDate,
	}

	if cfg.DataSourceId != "" {
		input.DataSourceId = &cfg.DataSourceId
	}

	if cfg.IdentityStoreId != "" {
		input.IdentityStoreId = &cfg.IdentityStoreId
	}

	err := graphql.NewClient(&cfg.BaseConfig).Mutate(ctx, &mutation, map[string]interface{}{"input": input})
	if err != nil {
		cfg.TargetLogger.Debug(fmt.Sprintf("subtask event update failed: %s", err.Error()))
	}
}

func GetSubtask(ctx context.Context, cfg *types.BaseTargetConfig, jobID, jobType, subtaskId string, responseResult interface{}) (*Subtask, error) {
	gqlQuery := fmt.Sprintf(`query jobSubtask{
		jobSubtask(jobId: "%s", jobType: "%s", subtaskId: "%s") {
            jobId
            jobType
            subtaskId
            status
            lastUpdate
            errors
            result {
            __typename
              ... on DataSourceImportResult {
                  dataObjectsAdded
                  dataObjectsRemoved
                  dataObjectsUpdated
                  warnings
              }
              ... on IdentityStoreImportResult {
                  groupsAdded
                  groupsRemoved
                  groupsUpdated
                  usersAdded
                  usersRemoved
                  usersUpdated
                  warnings
              }
              ... on AccessProviderImportResult {
                  accessAdded
                  accessRemoved
                  accessUpdated
                  warnings
              }
              ... on DataUsageImportResult {
                  edgesCreatedOrUpdated
                  edgesRemoved
                  statementsAdded
                  statementsFailed
                  statementsSkipped
                  warnings
              }
              ... on AccessProviderExportResult {
                  fileKey
                  fileLocation
                  warnings
			  }
              ... on AccessProviderSyncFeedbackResult {
                  accessNamesAdded
                  warnings
              }
            }
        }}`, jobID, jobType, subtaskId)

	gqlQuery = strings.ReplaceAll(gqlQuery, "\t", "")

	rawResponse, err := graphql.NewClient(&cfg.BaseConfig).ExecRaw(ctx, gqlQuery, nil)
	if err != nil {
		cfg.TargetLogger.Debug(fmt.Sprintf("failed loading subtask: %s", err.Error()))
		return nil, err
	}

	response := SubtaskResponse{Subtask{Result: responseResult}}

	err = json.Unmarshal(rawResponse, &response)
	if err != nil {
		cfg.TargetLogger.Debug(fmt.Sprintf("failed parsing subtask: %s", err.Error()))
		return nil, err
	}

	return &response.SubtaskResponse, nil
}

type Response struct {
	Job Job `json:"createJob"`
}

type Subtask struct {
	JobID      string      `json:"jobId"`
	JobType    string      `json:"jobType"`
	SubtaskId  string      `json:"subtaskId"`
	Status     JobStatus   `json:"status"`
	LastUpdate time.Time   `json:"lastUpdate"`
	Result     interface{} `json:"result"`
	Errors     []string    `json:"errors"`
}

type SubtaskResponse struct {
	SubtaskResponse Subtask `json:"jobSubtask"`
}

type Job struct {
	JobID *string `json:"jobId"`
}

type JobStatus int

const (
	Started JobStatus = iota
	InProgress
	DataRetrieve
	DataUpload
	Queued
	DataProcessing
	Completed
	Failed
	Skipped
	TimeOut
)

var AllJobStatus = []JobStatus{
	Started,
	InProgress,
	DataRetrieve,
	DataUpload,
	Queued,
	DataProcessing,
	Completed,
	Failed,
	Skipped,
	TimeOut,
}

var jobStatusNames = [...]string{"STARTED", "IN_PROGRESS", "DATA_RETRIEVE", "DATA_UPLOAD", "QUEUED", "DATA_PROCESSING", "COMPLETED", "FAILED", "SKIPPED", "TIMED_OUT"}
var jobStatusNameMap = map[string]JobStatus{
	"STARTED":         Started,
	"IN_PROGRESS":     InProgress,
	"DATA_RETRIEVE":   DataRetrieve,
	"DATA_UPLOAD":     DataUpload,
	"QUEUED":          Queued,
	"DATA_PROCESSING": DataProcessing,
	"COMPLETED":       Completed,
	"FAILED":          Failed,
	"SKIPPED":         Skipped,
	"TIMED_OUT":       TimeOut,
}

func (e JobStatus) IsValid() bool {
	switch e {
	case Started, InProgress, DataRetrieve, DataUpload, Queued, DataProcessing, Completed, Failed, Skipped, TimeOut:
		return true
	default:
		return false
	}
}

func (e JobStatus) IsRunning() bool {
	switch e {
	case Started, InProgress, DataRetrieve, DataUpload, Queued, DataProcessing:
		return true
	case Completed, Failed, Skipped, TimeOut:
		return false
	default:
		return false
	}
}

func (e JobStatus) String() string {
	return jobStatusNames[e]
}

func (e *JobStatus) UnmarshalJSON(b []byte) error {
	var j string

	err := json.Unmarshal(b, &j)

	if err != nil {
		return err
	}

	*e = jobStatusNameMap[j]

	return nil
}

func (e JobStatus) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(e.String())
	buffer.WriteString(`"`)

	return buffer.Bytes(), nil
}

func WaitForJobToComplete(ctx context.Context, jobID string, syncType string, subtaskId string, syncResult interface{}, cfg *types.BaseTargetConfig, currentStatus JobStatus, waitInterval int) (*Subtask, error) {
	i := 0
	errorCount := 0

	var subtask *Subtask
	var err error

	for currentStatus.IsRunning() || i == 0 {
		if currentStatus.IsRunning() {
			time.Sleep(time.Duration(waitInterval) * time.Second)
		}

		subtask, err = GetSubtask(ctx, cfg, jobID, syncType, subtaskId, syncResult)

		if err != nil || subtask == nil {
			errorCount++
			if errorCount > NrOfRetriesOnError {
				if err == nil {
					err = fmt.Errorf("received invalid job status")
				}

				cfg.TargetLogger.Error(fmt.Sprintf("exceeded threshold for number of retries for fetching subtask progress: %w", err))

				return nil, err
			} else {
				continue
			}
		}

		if currentStatus != subtask.Status {
			cfg.TargetLogger.Info(fmt.Sprintf("Update task status to %s", subtask.Status.String()))
		}

		currentStatus = subtask.Status
		cfg.TargetLogger.Debug(fmt.Sprintf("Current status on iteration %d: %s", i, currentStatus.String()))
		i += 1
	}

	return subtask, nil
}
