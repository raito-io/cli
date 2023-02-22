package identity_store

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// IdentityStoreSyncer interface needs to be implemented by any plugin that wants to import users and groups into a Raito identity store.
type IdentityStoreSyncer interface {
	SyncIdentityStore(ctx context.Context, config *IdentityStoreSyncConfig) (*IdentityStoreSyncResult, error)
	GetIdentityStoreMetaData(ctx context.Context) (*MetaData, error)
}

// IdentityStoreSyncerPlugin is used on the server (CLI) and client (plugin) side to integrate with the plugin system.
// A plugin should not be using this directly, but instead depend on the cli-plugin-base library to register the plugins.
type IdentityStoreSyncerPlugin struct {
	plugin.Plugin

	Impl IdentityStoreSyncer
}

func (p *IdentityStoreSyncerPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	RegisterIdentityStoreSyncServiceServer(s, &identityStoreSyncerGRPCServer{Impl: p.Impl})
	return nil
}

func (IdentityStoreSyncerPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &identityStoreSyncerGRPC{client: NewIdentityStoreSyncServiceClient(c)}, nil
}

// IdentityStoreSyncerName constant should not be used directly when implementing plugins.
// It's the registration name for the identity store syncer plugin,
// used by the CLI and the cli-plugin-base library (RegisterPlugins function) to register the plugins.
const IdentityStoreSyncerName = "identityStoreSyncer"

type identityStoreSyncerGRPC struct {
	client IdentityStoreSyncServiceClient
}

func (g *identityStoreSyncerGRPC) SyncIdentityStore(ctx context.Context, config *IdentityStoreSyncConfig) (*IdentityStoreSyncResult, error) {
	return g.client.SyncIdentityStore(ctx, config)
}

func (g *identityStoreSyncerGRPC) GetIdentityStoreMetaData(ctx context.Context) (*MetaData, error) {
	return g.client.GetIdentityStoreMetaData(ctx, &emptypb.Empty{})
}

type identityStoreSyncerGRPCServer struct {
	UnimplementedIdentityStoreSyncServiceServer

	Impl IdentityStoreSyncer
}

func (s *identityStoreSyncerGRPCServer) SyncIdentityStore(ctx context.Context, config *IdentityStoreSyncConfig) (*IdentityStoreSyncResult, error) {
	return s.Impl.SyncIdentityStore(ctx, config)
}

func (s *identityStoreSyncerGRPCServer) GetIdentityStoreMetaData(ctx context.Context, _ *emptypb.Empty) (*MetaData, error) {
	return s.Impl.GetIdentityStoreMetaData(ctx)
}
