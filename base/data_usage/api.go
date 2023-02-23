package data_usage

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// DataUsageSyncer interface needs to be implemented by any plugin that wants to import data usage information
// into a Raito data source.
type DataUsageSyncer interface {
	SyncDataUsage(ctx context.Context, config *DataUsageSyncConfig) (*DataUsageSyncResult, error)
}

// DataUsageSyncerPlugin is used on the server (CLI) and client (plugin) side to integrate with the plugin system.
// A plugin should not be using this directly, but instead depend on the cli-plugin-base library to register the plugins.
type DataUsageSyncerPlugin struct {
	plugin.Plugin

	Impl DataUsageSyncer
}

func (p *DataUsageSyncerPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	RegisterDataUsageSyncServiceServer(s, &dataUsageSyncerGRPCServer{Impl: p.Impl})
	return nil
}

func (DataUsageSyncerPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &dataUsageSyncerGRPC{client: NewDataUsageSyncServiceClient(c)}, nil
}

// DataUsageSyncerName constant should not be used directly when implementing plugins.
// It's the registration name for the data usage syncer plugin,
// used by the CLI and the cli-plugin-base library (RegisterPlugins function) to register the plugins.
const DataUsageSyncerName = "dataUsageSyncer"

type dataUsageSyncerGRPC struct{ client DataUsageSyncServiceClient }

func (g *dataUsageSyncerGRPC) SyncDataUsage(ctx context.Context, config *DataUsageSyncConfig) (*DataUsageSyncResult, error) {
	return g.client.SyncDataUsage(ctx, config)
}

type dataUsageSyncerGRPCServer struct {
	UnimplementedDataUsageSyncServiceServer

	Impl DataUsageSyncer
}

func (s *dataUsageSyncerGRPCServer) SyncDataUsage(ctx context.Context, config *DataUsageSyncConfig) (*DataUsageSyncResult, error) {
	return s.Impl.SyncDataUsage(ctx, config)
}
