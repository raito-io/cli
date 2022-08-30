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

func StartJob(cfg *target.BaseTargetConfig) (string, error) {
	if isRaitoCloudEnabled(cfg) {
		gqlQuery := fmt.Sprintf(`{ "query": "mutation createJob {
        createJob(input: { dataSourceId: \"%s\", eventTime: \"%s\" }) { jobId } }" }"`,
			cfg.DataSourceId, time.Now().Format(time.RFC3339))

		gqlQuery = strings.ReplaceAll(gqlQuery, "\n", "\\n")

		resp := Response{}
		_, err := graphql.ExecuteGraphQL(gqlQuery, cfg, &resp)

		if err != nil {
			return "", fmt.Errorf("error while executing import: %s", err.Error())
		}

		return *resp.Job.JobID, nil
	}

	return "offline-job", nil
}

func AddJobEvent(cfg *target.BaseTargetConfig, jobID, jobType, status string) {
	if isRaitoCloudEnabled(cfg) {
		gqlQuery := fmt.Sprintf(`{ "query":"mutation createJobEvent {
        createJobEvent(input: { jobId: \"%s\", dataSourceId: \"%s\", jobType: \"%s\", status: %s, eventTime: \"%s\" }) { jobId } }" }"`,
			jobID, cfg.DataSourceId, jobType, status, time.Now().Format(time.RFC3339))

		gqlQuery = strings.ReplaceAll(gqlQuery, "\n", "\\n")

		err := graphql.ExecuteGraphQLWithoutResponse(gqlQuery, cfg)
		if err != nil {
			cfg.Logger.Debug("job update failed: %s", err.Error())
		}
	}
}

func GetTaskStatus(cfg *target.BaseTargetConfig, jobID, jobType string, responseResult interface{}) (*JobStatus, error) {
	if isRaitoCloudEnabled(cfg) {
		gqlQuery := fmt.Sprintf(`{ "query": "query getJobTask {
        jobTask(jobId: \"%s\", jobType: \"%s\") {
            jobId
            jobType
            status
            lastUpdate
            result {
            __typename
              ... on DataSourceImportResult {
                  dataObjectsAdded
                  dataObjectsRemoved
                  dataObjectsUpdated
                  errors
              }
              ... on IdentityStoreImportResult {
                  groupsAdded
                  groupsRemoved
                  groupsUpdated
                  usersAdded
                  usersRemoved
                  usersUpdated
                  errors
              }
              ... on AccessProviderImportResult {
                  accessAdded
                  accessRemoved
                  accessUpdated
                  errors
              }
              ... on DataUsageImportResult {
                  edgesCreatedOrUpdated
                  edgesRemoved
                  statementsAdded
                  statementsFailed
                  statementsSkipped
                  errors
              }
            }
        }}"}`, jobID, jobType)

		gqlQuery = strings.ReplaceAll(gqlQuery, "\n", "\\n")

		response := TaskStatusResponse{Task{Result: responseResult}}
		_, err := graphql.ExecuteGraphQL(gqlQuery, cfg, &response)

		if err != nil {
			cfg.Logger.Debug("failed to load Task information: %s", err.Error())
			return nil, err
		}

		return &response.TaskResponse.Status, nil
	}

	return nil, nil
}

type Response struct {
	Job Job `json:"createJob"`
}

type Task struct {
	JobID      string      `json:"jobId"`
	JobType    string      `json:"jobType"`
	Status     JobStatus   `json:"status"`
	LastUpdate time.Time   `json:"lastUpdate"`
	Result     interface{} `json:"result"`
}

type TaskStatusResponse struct {
	TaskResponse Task `json:"jobTask"`
}

type Job struct {
	JobID *string `json:"jobId"`
}

type JobStatus int

const (
	Started JobStatus = iota
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
	DataRetrieve,
	DataUpload,
	Queued,
	DataProcessing,
	Completed,
	Failed,
	Skipped,
}

var jobStatusNames = [...]string{"STARTED", "DATA_RETRIEVE", "DATA_UPLOAD", "QUEUED", "DATA_PROCESSING", "COMPLETED", "FAILED", "SKIPPED"}
var jobStatusNameMap = map[string]JobStatus{
	"STARTED":         Started,
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
	case Started, DataRetrieve, DataUpload, Queued, DataProcessing, Completed, Failed, Skipped:
		return true
	default:
		return false
	}
}

func (e JobStatus) IsRunning() bool {
	switch e {
	case Started, DataRetrieve, DataUpload, Queued, DataProcessing:
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

func isRaitoCloudEnabled(cfg *target.BaseTargetConfig) bool {
	return cfg.ApiUser != "" && cfg.ApiSecret != ""
}
