package data_source

import (
	"encoding/json"
	"net/rpc"

	"github.com/hashicorp/go-plugin"

	"github.com/raito-io/cli/base/util/config"
	error2 "github.com/raito-io/cli/base/util/error"
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
	Bucket     = "bucket"
	Object     = "object"
	Folder     = "folder"
	File       = "file"
)

type GlobalPermission string

const (
	readGlobalPermission     GlobalPermission = "read"
	writeGlobalPermission    GlobalPermission = "write"
	insertGlobalPermission   GlobalPermission = "insert"
	updateGlobalPermission   GlobalPermission = "update"
	deleteGlobalPermission   GlobalPermission = "delete"
	truncateGlobalPermission GlobalPermission = "truncate"
)

type GlobalPermissionSet map[GlobalPermission]struct{}

func CreateGlobalPermissionSet(permissions ...GlobalPermission) GlobalPermissionSet {
	res := make(GlobalPermissionSet)
	for _, p := range permissions {
		res[p] = struct{}{}
	}

	return res
}

func (s GlobalPermissionSet) Values() []GlobalPermission {
	result := make([]GlobalPermission, 0, len(s))

	for permission := range s {
		result = append(result, permission)
	}

	return result
}

func (s GlobalPermissionSet) Append(permission ...GlobalPermission) {
	for _, p := range permission {
		s[p] = struct{}{}
	}
}

func JoinGlobalPermissionsSets(sets ...GlobalPermissionSet) GlobalPermissionSet {
	res := make(GlobalPermissionSet)
	for _, set := range sets {
		for permission := range set {
			res[permission] = struct{}{}
		}
	}

	return res
}

func (s GlobalPermissionSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values())
}

func (s *GlobalPermissionSet) UnmarshalJSON(data []byte) error {
	var permissions []GlobalPermission

	if err := json.Unmarshal(data, &permissions); err != nil {
		return err
	}

	*s = make(map[GlobalPermission]struct{})

	for _, permission := range permissions {
		(*s)[permission] = struct{}{}
	}

	return nil
}

/*
The list of global permissions
*/

// WriteGlobalPermission Get all rights to (over)write data
func WriteGlobalPermission() GlobalPermissionSet {
	return CreateGlobalPermissionSet(writeGlobalPermission)
}

// InsertGlobalPermission Get rights to add data
func InsertGlobalPermission() GlobalPermissionSet {
	set := WriteGlobalPermission()
	set.Append(insertGlobalPermission)

	return set
}

// UpdateGlobalPermission Get rights to modify data, not to delete a row
func UpdateGlobalPermission() GlobalPermissionSet {
	set := WriteGlobalPermission()
	set.Append(updateGlobalPermission)

	return set
}

// DeleteGlobalPermission Get all rights to delete data and the table
func DeleteGlobalPermission() GlobalPermissionSet {
	set := WriteGlobalPermission()
	set.Append(deleteGlobalPermission)

	return set
}

// TruncateGlobalPermission Get rights to delete data
func TruncateGlobalPermission() GlobalPermissionSet {
	set := DeleteGlobalPermission()
	set.Append(truncateGlobalPermission)

	return set
}

// ReadGlobalPermission Get access to read the data
func ReadGlobalPermission() GlobalPermissionSet {
	set := JoinGlobalPermissionsSets(
		TruncateGlobalPermission(),
		DeleteGlobalPermission(),
		UpdateGlobalPermission(),
		InsertGlobalPermission(),
		WriteGlobalPermission(),
	)
	set.Append(readGlobalPermission)

	return set
}

// DataSourceSyncConfig represents the configuration that is passed from the CLI to the DataAccessSyncer plugin interface.
// It contains all the necessary configuration parameters for the plugin to function.
type DataSourceSyncConfig struct {
	config.ConfigMap
	TargetFile   string
	DataSourceId string
}

// DataSourceSyncResult represents the result from the data source sync process.
// A potential error is also modeled in here so specific errors remain intact when passed over RPC.
type DataSourceSyncResult struct {
	Error *error2.ErrorResult
}

type DataObjectType struct {
	Name        string                     `json:"name"`
	Type        string                     `json:"type"`
	Label       string                     `json:"label"`
	Icon        string                     `json:"icon"`
	Permissions []DataObjectTypePermission `json:"permissions"`
	Children    []string                   `json:"children"`
}

type DataObjectTypePermission struct {
	Permission        string              `json:"permission"`
	GlobalPermissions GlobalPermissionSet `json:"globalPermissions,omitempty"`
	Description       string              `json:"description"`
}

type MetaData struct {
	DataObjectTypes   []DataObjectType `json:"dataObjectTypes"`
	SupportedFeatures []string         `json:"supportedFeatures"`
	Type              string           `json:"type"`
	Icon              string           `json:"icon"`
}

// DataSourceSyncer interface needs to be implemented by any plugin that wants to import data objects into a Raito data source.
type DataSourceSyncer interface {
	SyncDataSource(config *DataSourceSyncConfig) DataSourceSyncResult
	GetDataSourceMetaData() MetaData
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
		resp.Error = error2.ToErrorResult(err)
	}

	return resp
}

func (g *dataSourceSyncerRPC) GetDataSourceMetaData() MetaData {
	var resp MetaData

	err := g.client.Call("Plugin.GetDataSourceMetaData", new(interface{}), &resp)
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

func (s *dataSourceSyncerRPCServer) GetDataSourceMetaData(args interface{}, resp *MetaData) error {
	*resp = s.Impl.GetDataSourceMetaData()
	return nil
}
