package identity_store

import (
	"context"
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

type IdentityStoreImportConfig struct {
	target.BaseTargetConfig
	UserFile        string
	GroupFile       string
	DeleteUntouched bool
	ReplaceGroups   bool
}

type IdentityStoreImporter interface {
	TriggerImport(ctx context.Context, jobId string) (job.JobStatus, string, error)
}

type identityStoreImporter struct {
	config        *IdentityStoreImportConfig
	log           hclog.Logger
	statusUpdater job.TaskEventUpdater
}

func NewIdentityStoreImporter(config *IdentityStoreImportConfig, statusUpdater job.TaskEventUpdater) IdentityStoreImporter {
	logger := config.TargetLogger.With("identitystore", config.IdentityStoreId, "userfile", config.UserFile, "groupfile", config.GroupFile)
	isI := identityStoreImporter{config, logger, statusUpdater}

	return &isI
}

func (i *identityStoreImporter) TriggerImport(ctx context.Context, jobId string) (job.JobStatus, string, error) {
	if viper.GetBool(constants.SkipFileUpload) {
		// In the development environment, we skip the upload and use the local file for the import
		return i.doImport(jobId, i.config.UserFile, i.config.GroupFile)
	} else {
		userKey, groupKey, err := i.upload(ctx)
		if err != nil {
			return job.Failed, "", err
		}

		return i.doImport(jobId, userKey, groupKey)
	}
}

func (i *identityStoreImporter) upload(ctx context.Context) (string, string, error) {
	i.statusUpdater.SetStatusToDataUpload(ctx)

	userKey, err := file.UploadFile(i.config.UserFile, &i.config.BaseTargetConfig)
	if err != nil {
		return "", "", fmt.Errorf("error while uploading users JSON file to the backend: %s", err.Error())
	}

	groupKey, err := file.UploadFile(i.config.GroupFile, &i.config.BaseTargetConfig)
	if err != nil {
		return "", "", fmt.Errorf("error while uploading groups JSON file to the backend: %s", err.Error())
	}

	return userKey, groupKey, nil
}

func (i *identityStoreImporter) doImport(jobId string, userKey string, groupKey string) (job.JobStatus, string, error) {
	start := time.Now()

	gqlQuery := fmt.Sprintf(`{ "operationName": "ImportIdentityRequest", "variables":{}, "query": "mutation ImportIdentityRequest {
        importIdentityRequest(input: {
          jobId: \"%s\",
          importSettings: {
            identityStore: \"%s\",
            deleteUntouched: %t,
            replaceGroups: %t, 
            usersFileKey: \"%s\",
            groupsFileKey: \"%s\"
          }
        }) {
          subtask {
            status
            subtaskId
          }
        }
    }" }"`, jobId, i.config.IdentityStoreId, i.config.DeleteUntouched, i.config.ReplaceGroups, userKey, groupKey)

	gqlQuery = strings.Replace(gqlQuery, "\n", "\\n", -1)

	res := Response{}
	_, err := graphql.ExecuteGraphQL(gqlQuery, &i.config.BaseConfig, &res)

	if err != nil {
		return job.Failed, "", fmt.Errorf("error while executing identity store import: %s", err.Error())
	}

	subtask := res.Respose.Subtask

	i.log.Info(fmt.Sprintf("Executed import request in %s", time.Since(start).Round(time.Millisecond)))

	return subtask.Status, subtask.SubtaskId, nil
}

type subtaskResponse struct {
	Status    job.JobStatus `json:"status"`
	SubtaskId string        `json:"subtaskId"`
}

type QueryResponse struct {
	Subtask subtaskResponse `json:"subtask"`
}

type Response struct {
	Respose QueryResponse `json:"importIdentityRequest"`
}
