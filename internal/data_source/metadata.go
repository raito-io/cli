package data_source

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/raito-io/cli/common/api/data_source"
	"github.com/raito-io/cli/internal/graphql"
	"github.com/raito-io/cli/internal/target"
)

func SetMetaData(config target.BaseTargetConfig, metadata data_source.MetaData) error {
	logger := config.Logger.With("datasource", config.DataSourceId)
	start := time.Now()

	mdm, err := marshalMetaData(metadata)
	if err != nil {
		return fmt.Errorf("error while serializing metadata information: %s", err.Error())
	}
	md, err := fixMetaData(mdm)
	if err != nil {
		return err
	}
	// Additional escaping the quotes to embed in graphql string
	md = strings.Replace(md, "\"", "\\\"", -1)

	gqlQuery := fmt.Sprintf(`{ "operationName": "SetMetaData", "variables":{}, "query": "mutation SetMetaData {
        setDataSourceMetaData(id: \"%s\", input: %s) {
          id
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

func marshalMetaData(md data_source.MetaData) (string, error) {
	mdb, err := json.Marshal(md)
	if err != nil {
		return "", fmt.Errorf("error while serializing metadata information: %s", err.Error())
	}

	return string(mdb), nil
}

func fixMetaData(input string) (string, error) {
	md := input
	reg, err := regexp.Compile("\"([a-zA-Z]+)\":")
	if err != nil {
		return "", fmt.Errorf("unable to compile regular expression: %s", err.Error())
	}
	md = reg.ReplaceAllString(md, "$1:")

	// Remove all ',permissions:null'
	reg2, err := regexp.Compile("(,)?permissions:null")
	if err != nil {
		return "", fmt.Errorf("unable to compile regular expression: %s", err.Error())
	}
	md = reg2.ReplaceAllString(md, "${1}permissions:[]")

	// Remove all ',children:null'
	reg3, err := regexp.Compile("(,)?children:null")
	if err != nil {
		return "", fmt.Errorf("unable to compile regular expression: %s", err.Error())
	}
	md = reg3.ReplaceAllString(md, "${1}children:[]")

	// Remove all ',globalPermissions:null'
	reg4, err := regexp.Compile("(,)?globalPermissions:null")
	if err != nil {
		return "", fmt.Errorf("unable to compile regular expression: %s", err.Error())
	}
	md = reg4.ReplaceAllString(md, "")

	return md, nil
}
