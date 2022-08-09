package graphql

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/hashicorp/go-multierror"
	"github.com/raito-io/cli/internal/target"
	"github.com/raito-io/cli/internal/util/connect"
)

type GraphqlResponse struct {
	Data   interface{} `json:"data"`
	Errors []Error     `json:"errors"`
}

type Error struct {
	Message string `json:"message"`
}

type dummyResultObject struct{}

func ExecuteGraphQLWithoutResponse(gql string, config *target.BaseTargetConfig) error {
	result := dummyResultObject{}
	_, err := ExecuteGraphQL(gql, config, result)

	return err
}

func ExecuteGraphQL(gql string, config *target.BaseTargetConfig, resultObject interface{}) (*GraphqlResponse, error) {
	rawResponse, err := executeGraphQL(gql, config)

	if err != nil {
		return nil, err
	}

	response := GraphqlResponse{Data: resultObject}

	err = json.Unmarshal(rawResponse, &response)
	if err != nil {
		return nil, err
	}

	if len(response.Errors) > 0 {
		var errors []error
		for _, error := range response.Errors {
			errors = append(errors, fmt.Errorf("graphql server error: %s", error.Message))
		}

		return &response, multierror.Append(nil, errors...)
	}

	return &response, nil
}

func executeGraphQL(gql string, config *target.BaseTargetConfig) ([]byte, error) {
	resp, err := connect.DoPostToRaito("query", gql, "application/json", config)

	if err != nil {
		return nil, fmt.Errorf("error while executing graphql: %s", err.Error())
	}

	if resp.StatusCode >= 300 {
		buf, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error (HTTP %d) while executing graphql: %s - %s", resp.StatusCode, resp.Status, string(buf))
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
