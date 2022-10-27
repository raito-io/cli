package data_source

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/raito-io/cli/base/data_source"
	"github.com/raito-io/cli/internal/graphql"
	"github.com/raito-io/cli/internal/target"
)

func SetMetaData(config target.BaseTargetConfig, metadata data_source.MetaData) error {
	logger := config.Logger.With("datasource", config.DataSourceId)
	start := time.Now()

	mdm, err := marshalMetaData(metadata)
	if err != nil {
		return err
	}
	md := fixMetaData(mdm)

	// Additional escaping the quotes to embed in graphql string
	md = strings.Replace(md, "\"", "\\\"", -1)

	gqlQuery := fmt.Sprintf(`{ "operationName": "SetMetaData", "variables":{}, "query": "mutation SetMetaData {
        setDataSourceMetaData(id: \"%s\", input: %s) {
          id
        }
    }" }"`, config.DataSourceId, md)

	gqlQuery = strings.Replace(gqlQuery, "\n", "\\n", -1)

	err = graphql.ExecuteGraphQLWithoutResponse(gqlQuery, &config)
	if err != nil {
		return fmt.Errorf("error while executing SetDataSourceMetaData: %s", err.Error())
	}

	logger.Info(fmt.Sprintf("Done setting DataSource metadata in %s", time.Since(start).Round(time.Millisecond)))

	return nil
}

// marshalMetaData marshals the MetaData struct to a string
func marshalMetaData(md data_source.MetaData) (string, error) {
	mdb, err := json.Marshal(md)
	if err != nil {
		return "", fmt.Errorf("error while serializing data source metadata information: %s", err.Error())
	}

	return string(mdb), nil
}

// fixMetaData converts the marshaled JSON into a valid GraphQL input
func fixMetaData(input string) string {
	md := input
	reg := regexp.MustCompile("\"([a-zA-Z]+)\":")
	md = reg.ReplaceAllString(md, "$1:")

	// Replace a 'null' for the 'permissions' field to an empty array
	reg2 := regexp.MustCompile("(,)?permissions:null")
	md = reg2.ReplaceAllString(md, "${1}permissions:[]")

	// Replace a 'null' for the 'children' field to an empty array
	reg3 := regexp.MustCompile("(,)?children:null")
	md = reg3.ReplaceAllString(md, "${1}children:[]")

	// Remove all ',globalPermissions:null'
	reg4 := regexp.MustCompile("(,)?globalPermissions:null")
	md = reg4.ReplaceAllString(md, "")

	return md
}
