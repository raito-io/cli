package identity_store

import (
	"fmt"
	"strings"
	"time"

	"github.com/raito-io/cli/base/identity_store"
	"github.com/raito-io/cli/internal/graphql"
	"github.com/raito-io/cli/internal/target/types"
)

func SetMetaData(config types.BaseTargetConfig, metadata *identity_store.MetaData) error {
	logger := config.TargetLogger.With("identitystore", config.IdentityStoreId)
	start := time.Now()

	gqlQuery := fmt.Sprintf(`{ "operationName": "SetMetaData", "variables":{}, "query": "mutation SetMetaData {
        setIdentityStoreMetaData(id: \"%s\", input: {
          type: \"%s\",
          icon: \"%s\",
          canBeLinked: %t,
          canBeMaster: %t,
        }) {    
          ... on IdentityStore { id }
        }
    }" }"`, config.IdentityStoreId, metadata.Type, metadata.Icon, metadata.CanBeLinked, metadata.CanBeMaster)

	gqlQuery = strings.ReplaceAll(gqlQuery, "\n", "\\n")

	err := graphql.ExecuteGraphQLWithoutResponse(gqlQuery, &config.BaseConfig)
	if err != nil {
		return fmt.Errorf("error while executing SetIdentityStoreMetaData: %s", err.Error())
	}

	logger.Info(fmt.Sprintf("Done setting IdentityStore metadata in %s", time.Since(start).Round(time.Millisecond)))

	return nil
}
