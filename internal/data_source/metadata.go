package data_source

import (
	"context"
	"fmt"
	"time"

	graphql2 "github.com/hasura/go-graphql-client"

	"github.com/raito-io/cli/base/data_source"
	"github.com/raito-io/cli/internal/graphql"
	"github.com/raito-io/cli/internal/target"
)

func SetMetaData(ctx context.Context, config *target.BaseTargetConfig, metadata *data_source.MetaData) error {
	logger := config.TargetLogger.With("datasource", config.DataSourceId)
	start := time.Now()

	var m struct {
		SetDataSourceMetaData struct {
			DataSource struct {
				Id string
			} `graphql:"... on DataSource"`
			BaseError struct {
				Message string
			} `graphql:"... on BaseError"`
		} `graphql:"setDataSourceMetaData(id: $id, input :$input)"`
	}

	type DataSourceMetaDataInput struct {
		DataObjectTypes   []*data_source.DataObjectType `json:"dataObjectTypes,omitempty"`
		SupportedFeatures []string                      `json:"supportedFeatures,omitempty"`
		Type              string                        `json:"type,omitempty"`
		Icon              string                        `json:"icon,omitempty"`
	}

	input := DataSourceMetaDataInput{
		DataObjectTypes:   metadata.DataObjectTypes,
		SupportedFeatures: metadata.SupportedFeatures,
		Type:              metadata.Type,
		Icon:              metadata.Icon,
	}

	err := graphql.NewClient(&config.BaseConfig).Mutate(ctx, &m, map[string]interface{}{"id": graphql2.ID(config.DataSourceId), "input": input})
	if err != nil {
		err = graphql.ParseErrors(err)

		return fmt.Errorf("error while executing SetDataSourceMetaData: %w", err)
	}

	if m.SetDataSourceMetaData.BaseError.Message != "" {
		return fmt.Errorf("error while executing SetDataSourceMetaData: %s", m.SetDataSourceMetaData.BaseError.Message)
	}

	logger.Info(fmt.Sprintf("Done setting DataSource metadata in %s", time.Since(start).Round(time.Millisecond)))

	return nil
}
