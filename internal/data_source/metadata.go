package data_source

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	graphql2 "github.com/hasura/go-graphql-client"

	"github.com/raito-io/cli/base/data_source"
	"github.com/raito-io/cli/internal/graphql"
	"github.com/raito-io/cli/internal/target"
	"github.com/raito-io/golang-set/set"
)

func SetMetaData(ctx context.Context, config *target.BaseTargetConfig, metadata *data_source.MetaData) error {
	logger := config.TargetLogger.With("datasource", config.DataSourceId)
	start := time.Now()

	err := metadataConsistencyCheck(metadata)
	if err != nil {
		logger.Error(fmt.Sprintf("error while checking metadata consistency: %s", err.Error()))
		return err
	}

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
		DataObjectTypes     []*data_source.DataObjectType     `json:"dataObjectTypes,omitempty"`
		AccessProviderTypes []*data_source.AccessProviderType `json:"accessProviderTypes,omitempty"`
		SupportedFeatures   []string                          `json:"supportedFeatures,omitempty"`
		Type                string                            `json:"type,omitempty"`
		Icon                string                            `json:"icon,omitempty"`
		UsageMetaInfo       *data_source.UsageMetaInput       `json:"usageMetaInfo,omitempty"`
	}

	input := DataSourceMetaDataInput{
		DataObjectTypes:     metadata.DataObjectTypes,
		AccessProviderTypes: metadata.AccessProviderTypes,
		SupportedFeatures:   metadata.SupportedFeatures,
		Type:                metadata.Type,
		Icon:                metadata.Icon,
		UsageMetaInfo:       metadata.UsageMetaInfo,
	}

	err = graphql.NewClient(&config.BaseConfig).Mutate(ctx, &m, map[string]interface{}{"id": graphql2.ID(config.DataSourceId), "input": input})
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

func metadataConsistencyCheck(metadata *data_source.MetaData) error {
	objectTypeNames := set.Set[string]{}

	for _, objectType := range metadata.DataObjectTypes {
		objectTypeNames.Add(objectType.Name)
	}

	if metadata.UsageMetaInfo == nil || len(metadata.UsageMetaInfo.Levels) == 0 {
		return nil
	}

	if metadata.UsageMetaInfo.DefaultLevel == "" {
		return fmt.Errorf("default level is not provided or empty")
	}

	defaultLevelDefined := false
	errorMessages := []string{}

	for _, level := range metadata.UsageMetaInfo.Levels {
		if strings.EqualFold(metadata.UsageMetaInfo.DefaultLevel, level.Name) {
			defaultLevelDefined = true
		}

		for _, objectType := range level.DataObjectTypes {
			if !objectTypeNames.Contains(objectType) {
				msg := fmt.Sprintf("object type %s/%s not found in metadata", level.Name, objectType)
				errorMessages = append(errorMessages, msg)
			}
		}
	}

	if !defaultLevelDefined {
		errorMessages = append(errorMessages, fmt.Sprintf("default level %s not found in metadata", metadata))
	}

	if len(errorMessages) > 0 {
		return errors.New(strings.Join(errorMessages, ", "))
	}

	return nil
}
