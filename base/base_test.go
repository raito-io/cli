package base

import (
	"github.com/raito-io/cli/common/api"
	"github.com/raito-io/cli/common/api/data_access"
	"github.com/raito-io/cli/common/api/data_source"
	"github.com/raito-io/cli/common/api/identity_store"
	"github.com/stretchr/testify/assert"
	"testing"
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

func TestRegisterDataAccessService(t *testing.T) {
	Logger()
	dasi := dataAccessPlugin{}
	isi := infoPlugin{}
	pluginMap, err := buildPluginMap(&dasi, &isi)

	assert.Nil(t, err)
	assert.NotNil(t, pluginMap)
	assert.Equal(t, 2, len(pluginMap))
	assert.NotNil(t, pluginMap[data_access.DataAccessSyncerName])
}

func TestRegisterComboService(t *testing.T) {
	Logger()
	csi := combo{}
	pluginMap, err := buildPluginMap(&csi)

	assert.Nil(t, err)
	assert.NotNil(t, pluginMap)
	assert.Equal(t, 3, len(pluginMap))
	assert.NotNil(t, pluginMap[data_source.DataSourceSyncerName])
	assert.NotNil(t, pluginMap[api.InfoName])
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
	das1 := dataAccessPlugin{}
	das2 := dataAccessPlugin{}
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
	das1 := dataAccessPlugin{}
	a := another{}
	isi := infoPlugin{}
	pluginMap, err := buildPluginMap(&a, &das1, &isi)

	assert.Nil(t, err)
	assert.NotNil(t, pluginMap)
	assert.Equal(t, 2, len(pluginMap))
}

func TestRegisterNoInfoPlugin(t *testing.T) {
	Logger()
	das1 := dataAccessPlugin{}
	pluginMap, err := buildPluginMap(&das1)

	assert.NotNil(t, err)
	assert.Nil(t, pluginMap)
}

type another struct{}

type combo struct {}

func (s *combo) SyncIdentityStore(config *identity_store.IdentityStoreSyncConfig) identity_store.IdentityStoreSyncResult {
	return identity_store.IdentityStoreSyncResult{}
}

func (s *combo) SyncDataSource(config *data_source.DataSourceSyncConfig) data_source.DataSourceSyncResult {
	return data_source.DataSourceSyncResult{}
}

func (s *combo) PluginInfo() api.PluginInfo {
	return api.PluginInfo{}
}

type identityStoryPlugin struct {}

func (s *identityStoryPlugin) SyncIdentityStore(config *identity_store.IdentityStoreSyncConfig) identity_store.IdentityStoreSyncResult {
	return identity_store.IdentityStoreSyncResult{}
}

type dataSourcePlugin struct {}

func (s *dataSourcePlugin) SyncDataSource(config *data_source.DataSourceSyncConfig) data_source.DataSourceSyncResult {
	return data_source.DataSourceSyncResult{}
}

type dataAccessPlugin struct {}

func (s *dataAccessPlugin) SyncDataAccess(config *data_access.DataAccessSyncConfig) data_access.DataAccessSyncResult {
	return data_access.DataAccessSyncResult{}
}

type infoPlugin struct {}

func (s *infoPlugin) PluginInfo() api.PluginInfo {
	return api.PluginInfo{}
}