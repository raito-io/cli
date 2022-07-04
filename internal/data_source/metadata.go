package data_source

import (
	"encoding/json"
	"fmt"
	"github.com/raito-io/cli/common/api/data_source"
	"github.com/raito-io/cli/internal/graphql"
	"github.com/raito-io/cli/internal/target"
	"strings"
	"time"
)

func SetMetaData(config target.BaseTargetConfig, metadata data_source.MetaData) error {
	logger := config.Logger.With("datasource", config.DataSourceId)
	start := time.Now()

	mdb, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("error while serializing metadata information: %s", err.Error())
	}

	md := string(mdb)

	gqlQuery := fmt.Sprintf(`{ "operationName": "SetMetaData", "variables":{}, "query": "mutation SetMetaData {
        setDataSourceMetaData(id: \"%s\", input: %s) {
          errors
        }
    }" }"`, config.DataSourceId, md)

	gqlQuery = strings.Replace(gqlQuery, "\n", "\\n", -1)

	_, err = graphql.ExecuteGraphQL(gqlQuery, &config)
	if err != nil {
		return fmt.Errorf("error while executing import: %s", err.Error())
	}

	logger.Info(fmt.Sprintf("Done setting DataSource metadata in %s", time.Since(start).Round(time.Millisecond)))

	return nil
}
