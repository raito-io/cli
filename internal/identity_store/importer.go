package identity_store

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"

	"github.com/raito-io/cli/internal/util/tag"

	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/file"
	"github.com/raito-io/cli/internal/graphql"
	"github.com/raito-io/cli/internal/job"
	"github.com/raito-io/cli/internal/target/types"
)

type IdentityStoreImportConfig struct {
	types.BaseTargetConfig
	UserFile        string
	GroupFile       string
	DeleteUntouched bool
	ReplaceGroups   bool

	// TagSourcesScope is the set of sources that will be looked at when merging the tags. Tags with other sources will remain untouched.
	// If not specified, the default is to take all sources for which tags are defined in the import file.
	TagSourcesScope []string `json:"tagSourcesScope"`
}

type IdentityStoreImporter interface {
	TriggerImport(ctx context.Context, logger hclog.Logger, jobId string) (job.JobStatus, string, error)
}

type identityStoreImporter struct {
	config        *IdentityStoreImportConfig
	statusUpdater job.TaskEventUpdater
}

func NewIdentityStoreImporter(config *IdentityStoreImportConfig, statusUpdater job.TaskEventUpdater) IdentityStoreImporter {
	isI := identityStoreImporter{config, statusUpdater}

	return &isI
}

func (i *identityStoreImporter) TriggerImport(ctx context.Context, logger hclog.Logger, jobId string) (job.JobStatus, string, error) {
	logger = logger.With("identitystore", i.config.IdentityStoreId, "userfile", i.config.UserFile, "groupfile", i.config.GroupFile)

	if viper.GetBool(constants.SkipFileUpload) {
		// In the development environment, we skip the upload and use the local file for the import
		return i.doImport(logger, jobId, i.config.UserFile, i.config.GroupFile)
	} else {
		userKey, groupKey, err := i.upload(ctx, logger)
		if err != nil {
			return job.Failed, "", err
		}

		return i.doImport(logger, jobId, userKey, groupKey)
	}
}

func (i *identityStoreImporter) upload(ctx context.Context, logger hclog.Logger) (string, string, error) {
	i.statusUpdater.SetStatusToDataUpload(ctx)

	userKey, err := file.UploadFile(logger, i.config.UserFile, &i.config.BaseTargetConfig)
	if err != nil {
		return "", "", fmt.Errorf("error while uploading users JSON file to the backend: %s", err.Error())
	}

	groupKey, err := file.UploadFile(logger, i.config.GroupFile, &i.config.BaseTargetConfig)
	if err != nil {
		return "", "", fmt.Errorf("error while uploading groups JSON file to the backend: %s", err.Error())
	}

	return userKey, groupKey, nil
}

func (i *identityStoreImporter) doImport(logger hclog.Logger, jobId string, userKey string, groupKey string) (job.JobStatus, string, error) {
	start := time.Now()

	gqlQuery := fmt.Sprintf(`{ "operationName": "ImportIdentityRequest", "variables":{}, "query": "mutation ImportIdentityRequest {
        importIdentityRequest(input: {
          jobId: \"%s\",
          importSettings: {
            identityStore: \"%s\",
            deleteUntouched: %t,
            replaceGroups: %t, 
            usersFileKey: \"%s\",
            groupsFileKey: \"%s\",
            tagSourcesScope: %s
          }
        }) {
          subtask {
            status
            subtaskId
          }
        }
    }" }"`, jobId, i.config.IdentityStoreId, i.config.DeleteUntouched, i.config.ReplaceGroups, userKey, groupKey, strings.Replace(tag.SerializeTagList(i.config.TagSourcesScope), "\"", "\\\"", -1))

	gqlQuery = strings.Replace(gqlQuery, "\n", "\\n", -1)

	res := Response{}
	_, err := graphql.ExecuteGraphQL(gqlQuery, &i.config.BaseConfig, &res)

	if err != nil {
		return job.Failed, "", fmt.Errorf("error while executing identity store import: %s", err.Error())
	}

	subtask := res.Respose.Subtask

	logger.Info(fmt.Sprintf("Executed import request in %s", time.Since(start).Round(time.Millisecond)))

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
