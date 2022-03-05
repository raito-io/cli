package graphql

import (
	"fmt"
	"github.com/raito-io/cli/internal/target"
	"github.com/raito-io/cli/internal/util/connect"
	"io/ioutil"
)

type GraphqlResponse struct {
	Data interface{} `json:"data"`
	Errors []Error   `json:"errors"`
}

type Error struct {
	Message string `json:"message"`
}

func ExecuteGraphQL(gql string, config *target.BaseTargetConfig) ([]byte, error) {
	resp, err := connect.DoPostToRaito("query", gql, "application/json", config)

	if err != nil {
		return nil, fmt.Errorf("error while executing graphql: %s", err.Error())
	}
	if resp.StatusCode >= 300 {
		buf, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("error (HTTP %d) while executing graphql: %s - %s", resp.StatusCode, resp.Status, string(buf))
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
