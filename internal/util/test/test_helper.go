package test

import (
	"github.com/hashicorp/go-hclog"
	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/target/types"
	"github.com/spf13/viper"
)

func CreateBaseConfig(domain, apiUser, apiSecret, urlOverride string) (*types.BaseConfig, func()) {
	if urlOverride != "" {
		viper.Set(constants.URLOverrideFlag, urlOverride)
	}

	viper.Set(constants.DomainFlag, domain)
	viper.Set(constants.ApiUserFlag, apiUser)
	viper.Set(constants.ApiSecretFlag, apiSecret)

	config := types.BaseConfig{
		BaseLogger: hclog.Default(),
	}

	config.ReloadConfig()

	return &config, func() {
		if urlOverride != "" {
			defer viper.Set(constants.URLOverrideFlag, "")
		}

		defer viper.Set(constants.DomainFlag, "")
		defer viper.Set(constants.ApiUserFlag, "")
		defer viper.Set(constants.ApiSecretFlag, "")
	}
}
