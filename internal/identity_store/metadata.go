package identity_store

import (
	"fmt"
	"github.com/raito-io/cli/base/identity_store"
	"strings"
	"time"

	"github.com/raito-io/cli/internal/graphql"
	"github.com/raito-io/cli/internal/target"
)

func SetMetaData(config target.BaseTargetConfig, metadata identity_store.MetaData) error {
	logger := config.Logger.With("identitystore", config.IdentityStoreId)
	start := time.Now()

	gqlQuery := fmt.Sprintf(`{ "operationName": "SetMetaData", "variables":{}, "query": "mutation SetMetaData {
        setIdentityStoreMetaData(id: \"%s\", input: {
          type: \"%s\",
          icon: \"%s\"
        }) {
          id
        }
    }" }"`, config.IdentityStoreId, metadata.Type, metadata.Icon)

	gqlQuery = strings.Replace(gqlQuery, "\n", "\\n", -1)

	err := graphql.ExecuteGraphQLWithoutResponse(gqlQuery, &config)
	if err != nil {
		return fmt.Errorf("error while executing SetIdentityStoreMetaData: %s", err.Error())
	}

	logger.Info(fmt.Sprintf("Done setting IdentityStore metadata in %s", time.Since(start).Round(time.Millisecond)))

	return nil
}
