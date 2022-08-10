package job

import (
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
		gqlQuery := fmt.Sprintf(`{ "query": "mutation createJobEvent {
        createJobEvent(input: { jobId: \"%s\", dataSourceId: \"%s\", jobType: \"%s\", status: \"%s\", eventTime: \"%s\" }) { jobId } }" }"`,
			jobID, cfg.DataSourceId, jobType, status, time.Now().Format(time.RFC3339))

		gqlQuery = strings.ReplaceAll(gqlQuery, "\n", "\\n")

		err := graphql.ExecuteGraphQLWithoutResponse(gqlQuery, cfg)
		if err != nil {
			cfg.Logger.Debug("job update failed: %s", err.Error())
		}
	}
}

type Response struct {
	Job Job `json:"createJob"`
}

type Job struct {
	JobID *string `json:"jobId"`
}

func isRaitoCloudEnabled(cfg *target.BaseTargetConfig) bool {
	return cfg.ApiUser != "" && cfg.ApiSecret != ""
}
