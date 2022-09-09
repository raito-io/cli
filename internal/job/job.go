package job

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/raito-io/cli/internal/graphql"
	"github.com/raito-io/cli/internal/target"
)

type TaskEventUpdater struct {
	Cfg     *target.BaseTargetConfig
	JobId   string
	JobType string
}

func (u *TaskEventUpdater) AddTaskEvent(status JobStatus) {
	AddTaskEvent(u.Cfg, u.JobId, u.JobType, status)
}

func (u *TaskEventUpdater) GetSubtaskEventUpdater(subtask string) SubtaskEventUpdater {
	return SubtaskEventUpdater{
		Cfg:     u.Cfg,
		JobId:   u.JobId,
		JobType: u.JobType,
		Subtask: subtask,
	}
}

type SubtaskEventUpdater struct {
	Cfg          *target.BaseTargetConfig
	JobId        string
	JobType      string
	Subtask      string
	receivedDate *int64
}

func (u *SubtaskEventUpdater) AddSubtaskEvent(status JobStatus) {
	AddSubtaskEvent(u.Cfg, u.JobId, u.JobType, u.Subtask, status, u.receivedDate)
}

func (u *SubtaskEventUpdater) SetReceivedDate(receivedDate int64) {
	u.receivedDate = &receivedDate
}

func StartJob(cfg *target.BaseTargetConfig) (string, error) {
	gqlQuery := fmt.Sprintf(`{ "query": "mutation createJob {
        createJob(input: { dataSourceId: \"%s\", identityStoreId: \"%s\", eventTime: \"%s\" }) { jobId } }" }"`,
		cfg.DataSourceId, cfg.IdentityStoreId, time.Now().Format(time.RFC3339))

	gqlQuery = strings.ReplaceAll(gqlQuery, "\n", "\\n")

	resp := Response{}
	_, err := graphql.ExecuteGraphQL(gqlQuery, cfg, &resp)

	if err != nil {
		return "", fmt.Errorf("error while executing import: %s", err.Error())
	}

	return *resp.Job.JobID, nil
}

func UpdateJobEvent(cfg *target.BaseTargetConfig, jobID string, status JobStatus) {
	gqlQuery := fmt.Sprintf(`{ "query":"mutation updateJob {
        updateJob(id: \"%s\", input: { dataSourceId: \"%s\", identityStoreId: \"%s\", status: %s, eventTime: \"%s\" }) { jobId } }" }"`,
		jobID, cfg.DataSourceId, cfg.IdentityStoreId, status.String(), time.Now().Format(time.RFC3339))

	gqlQuery = strings.ReplaceAll(gqlQuery, "\n", "\\n")

	err := graphql.ExecuteGraphQLWithoutResponse(gqlQuery, cfg)
	if err != nil {
		cfg.Logger.Debug("job update failed: %s", err.Error())
	}
}

func AddTaskEvent(cfg *target.BaseTargetConfig, jobID, jobType string, status JobStatus) {
	gqlQuery := fmt.Sprintf(`{ "query":"mutation addTaskEvent {
        addTaskEvent(input: { jobId: \"%s\", dataSourceId: \"%s\", identityStoreId: \"%s\", jobType: \"%s\", status: %s, eventTime: \"%s\"}) {jobId } }" }"`,
		jobID, cfg.DataSourceId, cfg.IdentityStoreId, jobType, status.String(), time.Now().Format(time.RFC3339))

	gqlQuery = strings.ReplaceAll(gqlQuery, "\n", "\\n")

	err := graphql.ExecuteGraphQLWithoutResponse(gqlQuery, cfg)
	if err != nil {
		cfg.Logger.Debug("taskEvent update failed: %s", err.Error())
	}
}

func AddSubtaskEvent(cfg *target.BaseTargetConfig, jobID, jobType, subtask string, status JobStatus, receivedDate *int64) {
	gqlQuery := fmt.Sprintf(`{ "query":"mutation addSubtaskEvent {
        addSubtaskEvent(input: { jobId: \"%s\", dataSourceId: \"%s\", identityStoreId: \"%s\", jobType: \"%s\", subtaskId: \"%s\", status: %s, eventTime: \"%s\"`,
		jobID, cfg.DataSourceId, cfg.IdentityStoreId, jobType, subtask, status.String(), time.Now().Format(time.RFC3339))

	if receivedDate != nil {
		gqlQuery += fmt.Sprintf(", receivedDate: %d", *receivedDate)
	}

	gqlQuery += `}) { jobId } }" }"`

	gqlQuery = strings.ReplaceAll(gqlQuery, "\n", "\\n")

	err := graphql.ExecuteGraphQLWithoutResponse(gqlQuery, cfg)
	if err != nil {
		cfg.Logger.Debug("subtask event update failed: %s", err.Error())
	}
}

func GetSubtask(cfg *target.BaseTargetConfig, jobID, jobType, subtaskId string, responseResult interface{}) (*Subtask, error) {
	gqlQuery := fmt.Sprintf(`{ "query": "query getJobSubtask {
        jobSubtask(jobId: \"%s\", jobType: \"%s\", subtaskId: \"%s\") {
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
            }
        }}"}`, jobID, jobType, subtaskId)

	gqlQuery = strings.ReplaceAll(gqlQuery, "\n", "\\n")

	response := SubtaskResponse{Subtask{Result: responseResult}}
	_, err := graphql.ExecuteGraphQL(gqlQuery, cfg, &response)

	if err != nil {
		cfg.Logger.Debug("failed to load Subtask information: %s", err.Error())
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
}

var jobStatusNames = [...]string{"STARTED", "IN_PROGRESS", "DATA_RETRIEVE", "DATA_UPLOAD", "QUEUED", "DATA_PROCESSING", "COMPLETED", "FAILED", "SKIPPED"}
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
}

func (e JobStatus) IsValid() bool {
	switch e {
	case Started, InProgress, DataRetrieve, DataUpload, Queued, DataProcessing, Completed, Failed, Skipped:
		return true
	default:
		return false
	}
}

func (e JobStatus) IsRunning() bool {
	switch e {
	case Started, InProgress, DataRetrieve, DataUpload, Queued, DataProcessing:
		return true
	case Completed, Failed, Skipped:
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
