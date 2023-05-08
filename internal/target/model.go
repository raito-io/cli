package target

import (
	"github.com/hashicorp/go-hclog"
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

	LockAllWho    bool
	LockAllWhat   bool
	LockAllNames  bool
	LockAllDelete bool

	OnlyOutOfSyncData    bool
	SkipDataAccessImport bool

	DeleteUntouched bool
	DeleteTempFiles bool
	ReplaceGroups   bool

	DataObjectEnrichers []*EnricherConfig

	TargetLogger hclog.Logger
}