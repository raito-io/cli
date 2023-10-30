package target

import (
	"fmt"
	"reflect"

	"github.com/barkimedes/go-deepcopy"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"

	iconfig "github.com/raito-io/cli/internal/config"
	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/target/types"

	"github.com/spf13/viper"
)

func buildEnricherConfigFromMap(params map[string]interface{}) (*types.EnricherConfig, error) {
	eConfig := types.EnricherConfig{}
	eConfig.Parameters = make(map[string]string)

	err := fillStruct(&eConfig, params)
	if err != nil {
		return nil, err
	}

	for k, v := range params {
		if _, f := constants.KnownFlags[k]; f {
			continue
		}

		cv, err := iconfig.HandleField(v, reflect.String)
		if err != nil {
			return nil, err
		}

		stringValue, err := argumentToString(cv)
		if err != nil {
			return nil, err
		}

		if stringValue != nil {
			eConfig.Parameters[k] = *stringValue
		}
	}

	if eConfig.Name == "" {
		eConfig.Name = eConfig.ConnectorName
	}

	return &eConfig, nil
}

func buildDataObjectEnricherMap() (map[string]*types.EnricherConfig, error) {
	dataObjectEnricherMap := make(map[string]*types.EnricherConfig)

	dataObjectEnricherData := viper.Get(constants.DataObjectEnrichers)

	var errorResult error

	if enricherList, ok := dataObjectEnricherData.([]interface{}); ok {
		hclog.L().Debug(fmt.Sprintf("Found %d enrichers.", len(enricherList)))

		for _, enricherObj := range enricherList {
			enricher, ok := enricherObj.(map[string]interface{})
			if !ok {
				errorResult = multierror.Append(errorResult, fmt.Errorf("the data object enricher definition could not be parsed correctly (%v)", enricherObj))
				hclog.L().Debug(fmt.Sprintf("The data object enricher definition could not be parsed correctly (%v)", enricherObj))

				continue
			}

			eConfig, err := buildEnricherConfigFromMap(enricher)
			if err != nil {
				errorResult = multierror.Append(errorResult, fmt.Errorf("error while parsing the data object enricher configuration: %s", err.Error()))
				hclog.L().Error(fmt.Sprintf("error while parsing the data object enricher configuration: %s", err.Error()))

				continue
			}

			if eConfig == nil {
				continue
			}

			dataObjectEnricherMap[eConfig.Name] = eConfig
		}
	}

	return dataObjectEnricherMap, errorResult
}

func addDataObjectEnrichersToTargetConfig(tConfig *types.BaseTargetConfig, target map[string]interface{}, dataObjectEnricherMap map[string]*types.EnricherConfig) error {
	tConfig.DataObjectEnrichers = make([]*types.EnricherConfig, 0)

	dataObjectEnricherData := target[constants.DataObjectEnrichers]

	var errorResult error

	if enricherList, ok := dataObjectEnricherData.([]interface{}); ok {
		hclog.L().Debug(fmt.Sprintf("Found %d enrichers for target %s", len(enricherList), tConfig.Name))

		for _, enricherObj := range enricherList {
			enricher, ok := enricherObj.(map[string]interface{})
			if !ok {
				errorResult = multierror.Append(errorResult, fmt.Errorf("the data object enricher definition could not be parsed correctly (%v)", enricherObj))
				hclog.L().Debug(fmt.Sprintf("The data object enricher definition could not be parsed correctly (%v)", enricherObj))

				continue
			}

			eConfig, err := buildEnricherConfigFromMap(enricher)
			if err != nil {
				errorResult = multierror.Append(errorResult, fmt.Errorf("error while parsing the data object enricher configuration: %s", err.Error()))
				hclog.L().Error(fmt.Sprintf("error while parsing the data object enricher configuration: %s", err.Error()))

				continue
			}

			if eConfig == nil {
				continue
			}

			if mainEnricher, ok := dataObjectEnricherMap[eConfig.Name]; ok {
				newEnricherObj, err := deepcopy.Anything(mainEnricher)
				if err != nil {
					errorResult = multierror.Append(errorResult, fmt.Errorf("unable to copy data of enricher %q for target %q", eConfig.Name, tConfig.Name))
					hclog.L().Error(fmt.Sprintf("unable to copy data of enricher %q for target %q", eConfig.Name, tConfig.Name))

					continue
				}

				newEnricher := newEnricherObj.(*types.EnricherConfig)
				for k, v := range eConfig.Parameters {
					newEnricher.Parameters[k] = v
				}

				tConfig.DataObjectEnrichers = append(tConfig.DataObjectEnrichers, newEnricher)
			} else {
				tConfig.DataObjectEnrichers = append(tConfig.DataObjectEnrichers, eConfig)
			}
		}
	}

	return errorResult
}
