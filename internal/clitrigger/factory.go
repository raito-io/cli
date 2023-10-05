package clitrigger

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"

	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/graphql"
	"github.com/raito-io/cli/internal/target/types"
)

type websocketResult struct {
	CliTriggerUrl struct {
		Url string  `json:"url"`
		Err *string `json:"err,omitempty"`
	} `json:"cliTriggerUrl"`
}

func CreateCliTrigger(config *types.BaseConfig) (CliTrigger, error) {
	if viper.GetBool(constants.DisableWebsocketFlag) {
		config.BaseLogger.Info("Websocket sync is disabled. No CLI triggers will be captured")
		return &DummyCliTrigger{}, nil
	}

	cliTrigger, err := createWebsocketTrigger(config)
	if err != nil || cliTrigger == nil {
		return &DummyCliTrigger{}, err
	}

	return cliTrigger, nil
}

func createWebsocketTrigger(config *types.BaseConfig) (*WebsocketCliTrigger, error) {
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

	return NewWebsocketCliTrigger(config, result.CliTriggerUrl.Url), nil
}

func NewApUpdateTrigger(cliTrigger CliTrigger) *ApUpdateTriggerHandler {
	updateTrigger := NewApUpdateTriggerHandler()

	cliTrigger.Subscribe(updateTrigger)

	return updateTrigger
}
