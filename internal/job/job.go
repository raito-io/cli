package job

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/raito-io/cli/internal/graphql"
	"github.com/raito-io/cli/internal/target"
)

func StartJob(cfg *target.BaseTargetConfig) (string, error) {
	gqlQuery := fmt.Sprintf(`{ "query": "mutation createJob {
        createJob(input: { dataSourceId: \"%s\", eventTime: \"%s\" }) { jobId } }" }"`,
		cfg.DataSourceId, time.Now().Format(time.RFC3339))

	gqlQuery = strings.ReplaceAll(gqlQuery, "\n", "\\n")

	res, err := graphql.ExecuteGraphQL(gqlQuery, cfg)
	if err != nil {
		return "", fmt.Errorf("error while executing import: %s", err.Error())
	}

	resp := Response{}
	gr := graphql.GraphqlResponse{Data: &resp}

	err = json.Unmarshal(res, &gr)
	if err != nil {
		return "", fmt.Errorf("error while parsing job event result: %s", err.Error())
	}

	return *resp.Job.JobID, nil
}

func AddJobEvent(cfg *target.BaseTargetConfig, jobID, jobType, status string) {
	gqlQuery := fmt.Sprintf(`{ "query": "mutation createJobEvent {
        createJobEvent(input: { jobId: \"%s\", dataSourceId: \"%s\", jobType: \"%s\", status: \"%s\", eventTime: \"%s\" }) { jobId } }" }"`,
		jobID, cfg.DataSourceId, jobType, status, time.Now().Format(time.RFC3339))

	gqlQuery = strings.ReplaceAll(gqlQuery, "\n", "\\n")

	_, err := graphql.ExecuteGraphQL(gqlQuery, cfg)
	if err != nil {
		cfg.Logger.Debug("job update failed: %s", err.Error())
	}
}

type Response struct {
	Job Job `json:"createJob"`
}

type Job struct {
	JobID *string `json:"jobId"`
}
