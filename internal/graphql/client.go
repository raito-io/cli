package graphql

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/hasura/go-graphql-client"

	"github.com/raito-io/cli/internal/auth"
	"github.com/raito-io/cli/internal/target"
	"github.com/raito-io/cli/internal/util/merror"
	"github.com/raito-io/cli/internal/util/url"
	"github.com/raito-io/cli/internal/version"
)

type authedDoer struct {
	config *target.BaseConfig
}

func (d *authedDoer) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", "Raito CLI "+version.GetVersionString())

	err := auth.AddToken(req, d.config)
	if err != nil {
		return nil, fmt.Errorf("error while adding authorization token: %s", err.Error())
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while doing HTTP POST to %q: %s", req.URL.String(), err.Error())
	}

	return resp, nil
}

func NewClient(config *target.BaseConfig) *graphql.Client {
	return graphql.NewClient(url.GetRaitoURL()+"query", &authedDoer{config: config})
}

func ParseErrors(err error) error {
	gqlErrors := graphql.Errors{}
	if errors.As(err, &gqlErrors) {
		err = nil
		for _, gqlErr := range gqlErrors {
			err = merror.Append(err, errors.New(gqlErr.Message))
		}
	}

	return err
}
