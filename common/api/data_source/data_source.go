// Package data_source contains the API for the data source syncer.
package data_source

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/raito-io/cli/common/api"
	"github.com/raito-io/cli/common/util/config"
)

const (
	/*
		Data Source features
	*/
	ColumnMasking   = "columnMasking"
	RowFiltering    = "rowFiltering"
	ColumnFiltering = "columnFiltering"

	/*
		The list of standard Data Object Types
	*/
	Datasource = "datasource"
	Database   = "database"
	Schema     = "schema"
	Table      = "table"
	View       = "view"
	Column     = "column"
	Dataset    = "dataset"

	/*
		The list of global permissions
	*/
	Write = "write"
	Read  = "read"
)

// DataSourceSyncConfig represents the configuration that is passed from the CLI to the DataAccessSyncer plugin interface.
// It contains all the necessary configuration parameters for the plugin to function.
type DataSourceSyncConfig struct {
	config.ConfigMap
	TargetFile string
}

// DataSourceSyncResult represents the result from the data source sync process.
// A potential error is also modeled in here so specific errors remain intact when passed over RPC.
type DataSourceSyncResult struct {
	Error *api.ErrorResult
}

type DataObjectType struct {
	Name        string                     `json:"name"`
	Label       string                     `json:"label"`
	Icon        string                     `json:"icon"`
	Permissions []DataObjectTypePermission `json:"permissions"`
	Children    []string                   `json:"children"`
}

type DataObjectTypePermission struct {
	Permission        string   `json:"permission"`
	GlobalPermissions []string `json:"globalPermission"`
}

type MetaData struct {
	DataObjectTypes   []DataObjectType `json:"dataObjectTypes"`
	SupportedFeatures []string         `json:"supportedFeatures"`
}

// DataSourceSyncer interface needs to be implemented by any plugin that wants to import data objects into a Raito data source.
type DataSourceSyncer interface {
	SyncDataSource(config *DataSourceSyncConfig) DataSourceSyncResult
	GetMetaData() MetaData
}

// DataSourceSyncerPlugin is used on the server (CLI) and client (plugin) side to integrate with the plugin system.
// A plugin should not be using this directly, but instead depend on the cli-plugin-base library to register the plugins.
type DataSourceSyncerPlugin struct {
	Impl DataSourceSyncer
}

func (p *DataSourceSyncerPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &dataSourceSyncerRPCServer{Impl: p.Impl}, nil
}

func (DataSourceSyncerPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &dataSourceSyncerRPC{client: c}, nil
}

// DataSourceSyncerName constant should not be used directly when implementing plugins.
// It's the registration name for the data source syncer plugin,
// used by the CLI and the cli-plugin-base library (RegisterPlugins function) to register the plugins.
const DataSourceSyncerName = "dataSourceSyncer"

type dataSourceSyncerRPC struct{ client *rpc.Client }

func (g *dataSourceSyncerRPC) SyncDataSource(config *DataSourceSyncConfig) DataSourceSyncResult {
	var resp DataSourceSyncResult

	err := g.client.Call("Plugin.SyncDataSource", config, &resp)
	if err != nil && resp.Error == nil {
		resp.Error = api.ToErrorResult(err)
	}

	return resp
}

func (g *dataSourceSyncerRPC) GetMetaData() MetaData {
	var resp MetaData

	err := g.client.Call("Plugin.GetMetaData", nil, &resp)
	if err != nil {
		return MetaData{}
	}

	return resp
}

type dataSourceSyncerRPCServer struct {
	Impl DataSourceSyncer
}

func (s *dataSourceSyncerRPCServer) SyncDataSource(config *DataSourceSyncConfig, resp *DataSourceSyncResult) error {
	*resp = s.Impl.SyncDataSource(config)
	return nil
}

func (s *dataSourceSyncerRPCServer) GetMetaData(args interface{}, resp *MetaData) error {
	*resp = s.Impl.GetMetaData()
	return nil
}
