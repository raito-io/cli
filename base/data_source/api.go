package data_source

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/raito-io/cli/base/util/error/grpc_error"
	"github.com/raito-io/cli/base/util/version"
	"github.com/raito-io/cli/internal/version_management"
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
	readGlobalPermission   GlobalPermission = "read"
	writeGlobalPermission  GlobalPermission = "write"
	insertGlobalPermission GlobalPermission = "insert"
	updateGlobalPermission GlobalPermission = "update"
	deleteGlobalPermission GlobalPermission = "delete"
)

const (
	Read  string = "read"
	Write string = "write"
	Admin string = "admin"
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

func (s GlobalPermissionSet) StringValues() []string {
	result := make([]string, 0, len(s))

	for permission := range s {
		result = append(result, string(permission))
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

// ReadGlobalPermission Get access to read the data
func ReadGlobalPermission() GlobalPermissionSet {
	set := JoinGlobalPermissionsSets(
		DeleteGlobalPermission(),
		UpdateGlobalPermission(),
		InsertGlobalPermission(),
		WriteGlobalPermission(),
	)
	set.Append(readGlobalPermission)

	return set
}

// DataSourceSyncer interface needs to be implemented by any plugin that wants to import data objects into a Raito data source.
type DataSourceSyncer interface {
	version.CliVersionHandler

	SyncDataSource(ctx context.Context, config *DataSourceSyncConfig) (*DataSourceSyncResult, error)
	GetDataSourceMetaData(ctx context.Context) (*MetaData, error)
}

// DataSourceSyncerPlugin is used on the server (CLI) and client (plugin) side to integrate with the plugin system.
// A plugin should not be using this directly, but instead depend on the cli-plugin-base library to register the plugins.
type DataSourceSyncerPlugin struct {
	plugin.Plugin

	Impl DataSourceSyncer
}

func (p *DataSourceSyncerPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	RegisterDataSourceSyncServiceServer(s, &dataSourceSyncerGRPCServer{Impl: p.Impl})
	return nil
}

func (DataSourceSyncerPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &dataSourceSyncerGRPC{client: NewDataSourceSyncServiceClient(c)}, nil
}

// DataSourceSyncerName constant should not be used directly when implementing plugins.
// It's the registration name for the data source syncer plugin,
// used by the CLI and the cli-plugin-base library (RegisterPlugins function) to register the plugins.
const DataSourceSyncerName = "dataSourceSyncer"

type dataSourceSyncerGRPC struct{ client DataSourceSyncServiceClient }

func (g *dataSourceSyncerGRPC) SyncDataSource(ctx context.Context, config *DataSourceSyncConfig) (*DataSourceSyncResult, error) {
	return grpc_error.ParseErrorResult(g.client.SyncDataSource(ctx, config))
}

func (g *dataSourceSyncerGRPC) GetDataSourceMetaData(ctx context.Context) (*MetaData, error) {
	return grpc_error.ParseErrorResult(g.client.GetDataSourceMetaData(ctx, &emptypb.Empty{}))
}

func (g *dataSourceSyncerGRPC) CliVersionInformation(ctx context.Context) (*version.CliBuildInformation, error) {
	return grpc_error.ParseErrorResult(g.client.CliVersionInformation(ctx, &emptypb.Empty{}))
}

type dataSourceSyncerGRPCServer struct {
	UnimplementedDataSourceSyncServiceServer

	Impl DataSourceSyncer
}

func (s *dataSourceSyncerGRPCServer) SyncDataSource(ctx context.Context, config *DataSourceSyncConfig) (_ *DataSourceSyncResult, err error) {
	defer func() {
		err = grpc_error.GrpcDeferErrorHandling(err)
	}()

	return s.Impl.SyncDataSource(ctx, config)
}

func (s *dataSourceSyncerGRPCServer) GetDataSourceMetaData(ctx context.Context, _ *emptypb.Empty) (_ *MetaData, err error) {
	defer func() {
		err = grpc_error.GrpcDeferErrorHandling(err)
	}()

	return s.Impl.GetDataSourceMetaData(ctx)
}

func (s *dataSourceSyncerGRPCServer) CliVersionInformation(ctx context.Context, _ *emptypb.Empty) (_ *version.CliBuildInformation, err error) {
	defer func() {
		err = grpc_error.GrpcDeferErrorHandling(err)
	}()

	return s.Impl.CliVersionInformation(ctx)
}

type DataSourceSyncerVersionHandler struct {
}

func (h *DataSourceSyncerVersionHandler) CliVersionInformation(ctx context.Context) (*version.CliBuildInformation, error) {
	return version_management.CreateSyncerCliBuildInformation(MinimalCliVersion), nil
}
