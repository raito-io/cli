package identity_store

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/raito-io/cli/internal/file"
	"github.com/raito-io/cli/internal/graphql"
	"github.com/raito-io/cli/internal/target"
	"strings"
	"time"
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
	UsersAdded    int   `json:"usersAdded"`
	UsersUpdated  int   `json:"usersUpdated"`
	UsersRemoved  int   `json:"usersRemoved"`
	GroupsAdded   int   `json:"groupsAdded"`
	GroupsUpdated int   `json:"groupsUpdated"`
	GroupsRemoved int             `json:"groupsRemoved"`
	Errors        []graphql.Error `json:"_"`
}

type IdentityStoreImporter interface {
	TriggerImport() (*IdentityStoreImportResult, error)
}

type identityStoreImporter struct {
	config *IdentityStoreImportConfig
	log hclog.Logger
}

func NewIdentityStoreImporter(config *IdentityStoreImportConfig) IdentityStoreImporter {
	logger := config.Logger.With("identitystore", config.IdentityStoreId, "userfile", config.UserFile, "groupfile", config.GroupFile)
	isI := identityStoreImporter{config, logger	}
	return &isI
}

func (i *identityStoreImporter) TriggerImport() (*IdentityStoreImportResult, error) {
	userKey, groupKey, err := i.upload()
	if err != nil {
		return nil, err
	}

	return i.doImport(userKey, groupKey)
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

	res, err := graphql.ExecuteGraphQL(gqlQuery, &i.config.BaseTargetConfig)
	if err != nil {
		return nil, fmt.Errorf("error while executing identity store import: %s", err.Error())
	}

	ret, err := i.parseImportResult(res)
	if err != nil {
		return nil, err
	}
	if len(ret.Errors) > 0 {
		return ret, fmt.Errorf("errors while importing: %s", ret.Errors[0].Message)
	}

	i.log.Info(fmt.Sprintf("Executed import in %s", time.Since(start).Round(time.Millisecond)))

	return ret, nil
}

func (i *identityStoreImporter) parseImportResult(res []byte) (*IdentityStoreImportResult, error) {
	resp := Response{}
	gr := graphql.GraphqlResponse{ Data: &resp }
	err := json.Unmarshal(res, &gr)

	if err != nil {
		i.log.Error("error while parsing identity store import result.", "error", err.Error())
		return nil, err
	}

	// Flatten the result
	resp.ImportIdentityStore.Errors = gr.Errors

	return &(resp.ImportIdentityStore), nil
}

type Response struct {
	ImportIdentityStore IdentityStoreImportResult `json:"importIdentityStore"`
}
