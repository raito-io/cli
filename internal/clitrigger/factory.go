package clitrigger

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"

	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/graphql"
	"github.com/raito-io/cli/internal/target"
)

type websocketResult struct {
	CliTriggerUrl struct {
		Url string  `json:"url"`
		Err *string `json:"err,omitempty"`
	} `json:"cliTriggerUrl"`
}

func CreateCliTrigger(config *target.BaseConfig) (CliTrigger, error) {
	if viper.GetBool(constants.DisableWebsocketFlag) {
		config.BaseLogger.Info("Websocket sync is disabled. No CLI triggers will be captured")
		return nil, nil
	}

	return createWebsocketTrigger(config)
}

func createWebsocketTrigger(config *target.BaseConfig) (*WebsocketCliTrigger, error) {
	query := "{ \"query\": \"query CliTriggerWebSocket {\n    cliTriggerUrl {\n        ... on CliTriggerUrl {\n            url\n        }\n        ... on PermissionDeniedError {\n            err: message\n        }\n    }\n}\"}"
	query = strings.ReplaceAll(query, "\n", "\\n")

	result := websocketResult{}

	_, err := graphql.ExecuteGraphQL(query, config, &result)
	if err != nil {
		return nil, fmt.Errorf("create websocket trigger: %w", err)
	}

	if result.CliTriggerUrl.Err != nil {
		return nil, fmt.Errorf("create websocket trigger: %s", *result.CliTriggerUrl.Err)
	}

	if result.CliTriggerUrl.Url == "" {
		return nil, nil
	}

	return NewWebsocketCliTrigger(result.CliTriggerUrl.Url), nil
}
