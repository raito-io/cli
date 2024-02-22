package types

import (
	"reflect"
	"strings"

	"github.com/hashicorp/go-hclog"
	iconfig "github.com/raito-io/cli/internal/config"
	"github.com/raito-io/cli/internal/constants"
	"github.com/spf13/viper"

	"github.com/raito-io/cli/base/util/config"
)

type ConfigMap struct {
	Parameters map[string]string
}

func (c *ConfigMap) ToProtobufConfigMap() *config.ConfigMap {
	return &config.ConfigMap{
		Parameters: c.Parameters,
	}
}

type BaseConfig struct {
	ConfigMap
	ApiUser    string
	ApiSecret  string
	Domain     string
	BaseLogger hclog.Logger
	OtherArgs  []string
}

func (c *BaseConfig) ReloadConfig() error {
	apiUser, err := iconfig.HandleField(viper.GetString(constants.ApiUserFlag), reflect.String)
	if err != nil {
		return err
	}

	apiSecret, err := iconfig.HandleField(viper.GetString(constants.ApiSecretFlag), reflect.String)
	if err != nil {
		return err
	}

	domain, err := iconfig.HandleField(viper.GetString(constants.DomainFlag), reflect.String)
	if err != nil {
		return err
	}

	c.ApiUser = apiUser.(string)
	c.ApiSecret = apiSecret.(string)
	c.Domain = domain.(string)

	// Only read the parameters the first time as this is read from the command line + otherwise it would override the parameters as read from the
	if c.Parameters == nil {
		c.Parameters = BuildParameterMapFromArguments(c.OtherArgs)
	}

	return nil
}

func BuildParameterMapFromArguments(args []string) map[string]string {
	params := make(map[string]string)

	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "--") {
			arg := args[i][2:]
			if strings.Contains(arg, "=") {
				// The case where the flag is in the form of "--key=value"
				key := arg[0:strings.Index(arg, "=")]
				value := arg[strings.Index(arg, "=")+1:]
				params[key] = value
			} else if i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
				// The case where the flag is in the form of "--key value"
				params[arg] = args[i+1]
				i++
			} else {
				// Otherwise, we consider this a boolean flag
				params[arg] = "TRUE"
			}
		}
	}

	return params
}

type EnricherConfig struct {
	ConfigMap
	ConnectorName    string
	ConnectorVersion string
	Name             string
}

type BaseTargetConfig struct {
	BaseConfig
	ConnectorName    string
	ConnectorVersion string
	Name             string
	DataSourceId     string
	IdentityStoreId  string

	SkipIdentityStoreSync bool
	SkipDataSourceSync    bool
	SkipDataAccessSync    bool
	SkipDataUsageSync     bool

	LockAllWho            bool
	LockAllInheritance    bool
	LockAllWhat           bool
	LockAllNames          bool
	LockAllDelete         bool
	MakeNotInternalizable string

	TagOverwriteKeyForAccessProviderName   string
	TagOverwriteKeyForAccessProviderOwners string
	TagOverwriteKeyForDataObjectOwners     string

	OnlyOutOfSyncData    bool
	SkipDataAccessImport bool

	DeleteUntouched bool
	DeleteTempFiles bool
	ReplaceGroups   bool

	DataObjectParent   *string
	DataObjectExcludes []string

	DataObjectEnrichers []*EnricherConfig

	TargetLogger hclog.Logger
}
