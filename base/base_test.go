package base

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/raito-io/cli/base/access_provider"
	"github.com/raito-io/cli/base/data_source"
	"github.com/raito-io/cli/base/identity_store"
	"github.com/raito-io/cli/base/util/plugin"
)

func TestRegisterIdentityStoreService(t *testing.T) {
	Logger()
	issi := identityStoryPlugin{}
	isi := plugin.UnimplementedInfoServiceServer{}
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
	isi := plugin.UnimplementedInfoServiceServer{}
	pluginMap, err := buildPluginMap(&dssi, &isi)

	assert.Nil(t, err)
	assert.NotNil(t, pluginMap)
	assert.Equal(t, 2, len(pluginMap))
	assert.NotNil(t, pluginMap[data_source.DataSourceSyncerName])
}

func TestRegisterAccessSyncService(t *testing.T) {
	Logger()
	dasi := accessSyncPlugin{}
	isi := plugin.UnimplementedInfoServiceServer{}
	pluginMap, err := buildPluginMap(&dasi, &isi)

	assert.Nil(t, err)
	assert.NotNil(t, pluginMap)
	assert.Equal(t, 2, len(pluginMap))
	assert.NotNil(t, pluginMap[access_provider.AccessSyncerName])
}

func TestRegisterDoubleDataAccessService(t *testing.T) {
	Logger()
	das1 := accessSyncPlugin{}
	das2 := accessSyncPlugin{}
	isi := plugin.UnimplementedInfoServiceServer{}
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
	isi := plugin.UnimplementedInfoServiceServer{}
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

type identityStoryPlugin struct {
	identity_store.IdentityStoreSyncerVersionHandler
}

func (s *identityStoryPlugin) SyncIdentityStore(_ context.Context, _ *identity_store.IdentityStoreSyncConfig) (*identity_store.IdentityStoreSyncResult, error) {
	return &identity_store.IdentityStoreSyncResult{}, nil
}

func (s *identityStoryPlugin) GetIdentityStoreMetaData(_ context.Context) (*identity_store.MetaData, error) {
	return &identity_store.MetaData{}, nil
}

type dataSourcePlugin struct {
	data_source.DataSourceSyncerVersionHandler
}

func (s *dataSourcePlugin) SyncDataSource(_ context.Context, _ *data_source.DataSourceSyncConfig) (*data_source.DataSourceSyncResult, error) {
	return &data_source.DataSourceSyncResult{}, nil
}

func (s *dataSourcePlugin) GetDataSourceMetaData(_ context.Context) (*data_source.MetaData, error) {
	return &data_source.MetaData{}, nil
}

type accessSyncPlugin struct {
	access_provider.AccessSyncerVersionHandler
}

func (s *accessSyncPlugin) SyncFromTarget(_ context.Context, _ *access_provider.AccessSyncFromTarget) (*access_provider.AccessSyncResult, error) {
	return &access_provider.AccessSyncResult{}, nil
}

func (s *accessSyncPlugin) SyncToTarget(_ context.Context, _ *access_provider.AccessSyncToTarget) (*access_provider.AccessSyncResult, error) {
	return &access_provider.AccessSyncResult{}, nil
}

func (s *accessSyncPlugin) SyncConfig(_ context.Context) (*access_provider.AccessSyncConfig, error) {
	return &access_provider.AccessSyncConfig{}, nil
}
