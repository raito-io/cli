package identity_store

import (
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
	ReplaceTags     bool
}

type IdentityStoreImporter interface {
	TriggerImport(jobId string) (job.JobStatus, error)
}

type identityStoreImporter struct {
	config        *IdentityStoreImportConfig
	log           hclog.Logger
	statusUpdater func(status job.JobStatus)
}

func NewIdentityStoreImporter(config *IdentityStoreImportConfig, statusUpdater func(status job.JobStatus)) IdentityStoreImporter {
	logger := config.Logger.With("identitystore", config.IdentityStoreId, "userfile", config.UserFile, "groupfile", config.GroupFile)
	isI := identityStoreImporter{config, logger, statusUpdater}

	return &isI
}

func (i *identityStoreImporter) TriggerImport(jobId string) (job.JobStatus, error) {
	env := viper.GetString(constants.EnvironmentFlag)
	if env == constants.EnvironmentDev {
		// In the development environment, we skip the upload and use the local file for the import
		return i.doImport(jobId, i.config.UserFile, i.config.GroupFile)
	} else {
		userKey, groupKey, err := i.upload()
		if err != nil {
			return job.Failed, err
		}

		return i.doImport(jobId, userKey, groupKey)
	}
}

func (i *identityStoreImporter) upload() (string, string, error) {
	i.statusUpdater(job.DataUpload)

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

func (i *identityStoreImporter) doImport(jobId string, userKey string, groupKey string) (job.JobStatus, error) {
	start := time.Now()

	gqlQuery := fmt.Sprintf(`{ "operationName": "ImportIdentityRequest", "variables":{}, "query": "mutation ImportIdentityRequest {
        importIdentityRequest(input: {
          jobId: \"%s\",
          importSettings: {
            identityStore: \"%s\",
            deleteUntouched: %t,
            replaceGroups: %t,
            replaceTags: %t,  
            usersFileKey: \"%s\",
            groupsFileKey: \"%s\"
          }
        }) {
          jobStatus
        }
    }" }"`, jobId, i.config.IdentityStoreId, i.config.DeleteUntouched, i.config.ReplaceGroups, i.config.ReplaceTags, userKey, groupKey)

	gqlQuery = strings.Replace(gqlQuery, "\n", "\\n", -1)

	res := Response{}
	_, err := graphql.ExecuteGraphQL(gqlQuery, &i.config.BaseTargetConfig, &res)

	if err != nil {
		return job.Failed, fmt.Errorf("error while executing identity store import: %s", err.Error())
	}

	ret := res.Respose.Status

	i.log.Info(fmt.Sprintf("Executed import request in %s", time.Since(start).Round(time.Millisecond)))

	return ret, nil
}

type QueryResponse struct {
	Status job.JobStatus `json:"jobStatus"`
}

type Response struct {
	Respose QueryResponse `json:"importIdentityRequest"`
}
