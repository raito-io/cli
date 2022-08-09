package identity_store

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/file"
	"github.com/raito-io/cli/internal/graphql"
	"github.com/raito-io/cli/internal/target"
	"github.com/spf13/viper"
)

type IdentityStoreImportConfig struct {
	target.BaseTargetConfig
	UserFile        string
	GroupFile       string
	DeleteUntouched bool
	ReplaceGroups   bool
	ReplaceTags     bool
}

type IdentityStoreImportResult struct {
	UsersAdded    int `json:"usersAdded"`
	UsersUpdated  int `json:"usersUpdated"`
	UsersRemoved  int `json:"usersRemoved"`
	GroupsAdded   int `json:"groupsAdded"`
	GroupsUpdated int `json:"groupsUpdated"`
	GroupsRemoved int `json:"groupsRemoved"`
}

type IdentityStoreImporter interface {
	TriggerImport() (*IdentityStoreImportResult, error)
}

type identityStoreImporter struct {
	config *IdentityStoreImportConfig
	log    hclog.Logger
}

func NewIdentityStoreImporter(config *IdentityStoreImportConfig) IdentityStoreImporter {
	logger := config.Logger.With("identitystore", config.IdentityStoreId, "userfile", config.UserFile, "groupfile", config.GroupFile)
	isI := identityStoreImporter{config, logger}

	return &isI
}

func (i *identityStoreImporter) TriggerImport() (*IdentityStoreImportResult, error) {
	env := viper.GetString(constants.EnvironmentFlag)
	if env == constants.EnvironmentDev {
		// In the development environment, we skip the upload and use the local file for the import
		return i.doImport(i.config.UserFile, i.config.GroupFile)
	} else {
		userKey, groupKey, err := i.upload()
		if err != nil {
			return nil, err
		}

		return i.doImport(userKey, groupKey)
	}
}

func (i *identityStoreImporter) upload() (string, string, error) {
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

func (i *identityStoreImporter) doImport(userKey string, groupKey string) (*IdentityStoreImportResult, error) {
	start := time.Now()

	gqlQuery := fmt.Sprintf(`{ "operationName": "ImportIdentityStore", "variables":{}, "query": "mutation ImportIdentityStore {
        importIdentityStore(input: {
          identityStore: \"%s\",
          deleteUntouched: %t,
          replaceGroups: %t,
          replaceTags: %t,  
          users: \"%s\",
          groups: \"%s\"
        }) {
          usersAdded
          usersUpdated
          usersRemoved
          groupsAdded
          groupsUpdated
          groupsRemoved
          errors
        }
    }" }"`, i.config.IdentityStoreId, i.config.DeleteUntouched, i.config.ReplaceGroups, i.config.ReplaceTags, userKey, groupKey)

	gqlQuery = strings.Replace(gqlQuery, "\n", "\\n", -1)

	res := Response{}
	_, err := graphql.ExecuteGraphQL(gqlQuery, &i.config.BaseTargetConfig, &res)
	if err != nil {
		return nil, fmt.Errorf("error while executing identity store import: %s", err.Error())
	}

	ret := &res.ImportIdentityStore

	i.log.Info(fmt.Sprintf("Executed import in %s", time.Since(start).Round(time.Millisecond)))

	return ret, nil
}

type Response struct {
	ImportIdentityStore IdentityStoreImportResult `json:"importIdentityStore"`
}
