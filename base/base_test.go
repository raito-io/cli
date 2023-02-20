package base

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"

	"github.com/raito-io/cli/base/access_provider"
	"github.com/raito-io/cli/base/data_source"
	"github.com/raito-io/cli/base/identity_store"
	"github.com/raito-io/cli/base/util/plugin"
)

func TestRegisterIdentityStoreService(t *testing.T) {
	Logger()
	issi := identityStoryPlugin{}
	isi := infoPlugin{}
	pluginMap, err := buildPluginMap(&issi, &isi)

	assert.Nil(t, err)
	assert.NotNil(t, pluginMap)
	assert.Equal(t,
		2, len(pluginMap))
	assert.NotNil(t, pluginMap[identity_store.IdentityStoreSyncerName])
}

func TestRegisterDataSourceService(t *testing.T) {
	Logger()
	dssi := dataSourcePlugin{}
	isi := infoPlugin{}
	pluginMap, err := buildPluginMap(&dssi, &isi)

	assert.Nil(t, err)
	assert.NotNil(t, pluginMap)
	assert.Equal(t, 2, len(pluginMap))
	assert.NotNil(t, pluginMap[data_source.DataSourceSyncerName])
}

func TestRegisterAccessSyncService(t *testing.T) {
	Logger()
	dasi := accessSyncPlugin{}
	isi := infoPlugin{}
	pluginMap, err := buildPluginMap(&dasi, &isi)

	assert.Nil(t, err)
	assert.NotNil(t, pluginMap)
	assert.Equal(t, 2, len(pluginMap))
	assert.NotNil(t, pluginMap[access_provider.AccessSyncerName])
}

func TestRegisterComboService(t *testing.T) {
	Logger()
	csi := combo{}
	pluginMap, err := buildPluginMap(&csi)

	assert.Nil(t, err)
	assert.NotNil(t, pluginMap)
	assert.Equal(t, 3, len(pluginMap))
	assert.NotNil(t, pluginMap[data_source.DataSourceSyncerName])
	assert.NotNil(t, pluginMap[plugin.InfoName])
	assert.NotNil(t, pluginMap[identity_store.IdentityStoreSyncerName])
}

func TestRegisterDoubleIdentityStoreService(t *testing.T) {
	Logger()
	csi := combo{}
	issi := identityStoryPlugin{}
	pluginMap, err := buildPluginMap(&issi, &csi)

	assert.NotNil(t, err)
	assert.Nil(t, pluginMap)
}

func TestRegisterDoubleDataSourceService(t *testing.T) {
	Logger()
	csi := combo{}
	dssi := dataSourcePlugin{}
	pluginMap, err := buildPluginMap(&dssi, &csi)

	assert.NotNil(t, err)
	assert.Nil(t, pluginMap)
}

func TestRegisterDoubleDataAccessService(t *testing.T) {
	Logger()
	das1 := accessSyncPlugin{}
	das2 := accessSyncPlugin{}
	isi := infoPlugin{}
	pluginMap, err := buildPluginMap(&das1, &das2, &isi)

	assert.NotNil(t, err)
	assert.Nil(t, pluginMap)
}

func TestRegisterNoServices(t *testing.T) {
	Logger()
	a := another{}
	pluginMap, err := buildPluginMap(&a)

	assert.NotNil(t, err)
	assert.Nil(t, pluginMap)
}

func TestRegisterIgnoreNoService(t *testing.T) {
	Logger()
	das1 := accessSyncPlugin{}
	a := another{}
	isi := infoPlugin{}
	pluginMap, err := buildPluginMap(&a, &das1, &isi)

	assert.Nil(t, err)
	assert.NotNil(t, pluginMap)
	assert.Equal(t, 2, len(pluginMap))
}

func TestRegisterNoInfoPlugin(t *testing.T) {
	Logger()
	das1 := accessSyncPlugin{}
	pluginMap, err := buildPluginMap(&das1)

	assert.NotNil(t, err)
	assert.Nil(t, pluginMap)
}

type another struct{}

type combo struct{}

func (s *combo) SyncIdentityStore(config *identity_store.IdentityStoreSyncConfig) identity_store.IdentityStoreSyncResult {
	return identity_store.IdentityStoreSyncResult{}
}

func (s *combo) SyncDataSource(config *data_source.DataSourceSyncConfig) data_source.DataSourceSyncResult {
	return data_source.DataSourceSyncResult{}
}

func (s *combo) GetDataSourceMetaData() data_source.MetaData {
	return data_source.MetaData{}
}

func (s *combo) GetIdentityStoreMetaData() identity_store.MetaData {
	return identity_store.MetaData{}
}

func (s *combo) PluginInfo() plugin.PluginInfo {
	return plugin.PluginInfo{}
}

func (s *combo) CliBuildVersion() semver.Version {
	return *semver.MustParse("3.0.0")
}

func (s *combo) CliMinimalVersion() semver.Version {
	return *semver.MustParse("0.0.0")
}

type identityStoryPlugin struct{}

func (s *identityStoryPlugin) SyncIdentityStore(config *identity_store.IdentityStoreSyncConfig) identity_store.IdentityStoreSyncResult {
	return identity_store.IdentityStoreSyncResult{}
}

func (s *identityStoryPlugin) GetIdentityStoreMetaData() identity_store.MetaData {
	return identity_store.MetaData{}
}

type dataSourcePlugin struct{}

func (s *dataSourcePlugin) SyncDataSource(config *data_source.DataSourceSyncConfig) data_source.DataSourceSyncResult {
	return data_source.DataSourceSyncResult{}
}

func (s *dataSourcePlugin) GetDataSourceMetaData() data_source.MetaData {
	return data_source.MetaData{}
}

type accessSyncPlugin struct{}

func (s *accessSyncPlugin) SyncFromTarget(config *access_provider.AccessSyncFromTarget) access_provider.AccessSyncResult {
	return access_provider.AccessSyncResult{}
}

func (s *accessSyncPlugin) SyncToTarget(config *access_provider.AccessSyncToTarget) access_provider.AccessSyncResult {
	return access_provider.AccessSyncResult{}
}

func (s *accessSyncPlugin) SyncConfig() access_provider.AccessSyncConfig {
	return access_provider.AccessSyncConfig{}
}

type infoPlugin struct{}

func (s *infoPlugin) PluginInfo() plugin.PluginInfo {
	return plugin.PluginInfo{}
}

func (s *infoPlugin) CliBuildVersion() semver.Version {
	return *semver.MustParse("3.0.0")
}

func (s *infoPlugin) CliMinimalVersion() semver.Version {
	return *semver.MustParse("0.0.0")
}
