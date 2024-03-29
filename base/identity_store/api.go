package identity_store

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/raito-io/cli/base/util/config"
	"github.com/raito-io/cli/base/util/error/grpc_error"
	"github.com/raito-io/cli/base/util/version"
	"github.com/raito-io/cli/internal/version_management"
)

// IdentityStoreSyncer interface needs to be implemented by any plugin that wants to import users and groups into a Raito identity store.
type IdentityStoreSyncer interface {
	version.CliVersionHandler

	SyncIdentityStore(ctx context.Context, config *IdentityStoreSyncConfig) (*IdentityStoreSyncResult, error)
	GetIdentityStoreMetaData(ctx context.Context, config *config.ConfigMap) (*MetaData, error)
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
	return grpc_error.ParseErrorResult(g.client.SyncIdentityStore(ctx, config))
}

func (g *identityStoreSyncerGRPC) GetIdentityStoreMetaData(ctx context.Context, configMap *config.ConfigMap) (*MetaData, error) {
	return grpc_error.ParseErrorResult(g.client.GetIdentityStoreMetaData(ctx, configMap))
}

func (g *identityStoreSyncerGRPC) CliVersionInformation(ctx context.Context) (*version.CliBuildInformation, error) {
	return grpc_error.ParseErrorResult(g.client.CliVersionInformation(ctx, &emptypb.Empty{}))
}

type identityStoreSyncerGRPCServer struct {
	UnimplementedIdentityStoreSyncServiceServer

	Impl IdentityStoreSyncer
}

func (s *identityStoreSyncerGRPCServer) SyncIdentityStore(ctx context.Context, config *IdentityStoreSyncConfig) (_ *IdentityStoreSyncResult, err error) {
	defer func() {
		err = grpc_error.GrpcDeferErrorHandling(err)
	}()

	return s.Impl.SyncIdentityStore(ctx, config)
}

func (s *identityStoreSyncerGRPCServer) GetIdentityStoreMetaData(ctx context.Context, configMap *config.ConfigMap) (_ *MetaData, err error) {
	defer func() {
		err = grpc_error.GrpcDeferErrorHandling(err)
	}()

	return s.Impl.GetIdentityStoreMetaData(ctx, configMap)
}

func (s *identityStoreSyncerGRPCServer) CliVersionInformation(ctx context.Context, _ *emptypb.Empty) (_ *version.CliBuildInformation, err error) {
	defer func() {
		err = grpc_error.GrpcDeferErrorHandling(err)
	}()

	return s.Impl.CliVersionInformation(ctx)
}

type IdentityStoreSyncerVersionHandler struct {
}

func (h *IdentityStoreSyncerVersionHandler) CliVersionInformation(ctx context.Context) (*version.CliBuildInformation, error) {
	return version_management.CreateSyncerCliBuildInformation(MinimalCliVersion), nil
}
