package access_provider

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/raito-io/cli/base/util/version"
	version2 "github.com/raito-io/cli/internal/version"
)

// AccessSyncer interface needs to be implemented by any plugin that wants to sync access controls between Raito and the data source.
// This sync can be in the 2 directions or in just 1 depending on the parameters set in AccessSyncConfig.
type AccessSyncer interface {
	version.CliVersionHandler

	SyncFromTarget(ctx context.Context, config *AccessSyncFromTarget) (*AccessSyncResult, error)
	SyncToTarget(ctx context.Context, config *AccessSyncToTarget) (*AccessSyncResult, error)

	SyncConfig(ctx context.Context) (*AccessSyncConfig, error)
}

// AccessSyncerPlugin is used on the server (CLI) and client (plugin) side to integrate with the plugin system.
// A plugin should not be using this directly, but instead depend on the cli-plugin-base library to register the plugins.
type AccessSyncerPlugin struct {
	plugin.Plugin

	Impl AccessSyncer
}

func (p AccessSyncerPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	RegisterAccessProviderSyncServiceServer(s, &accessSyncerGRPCServer{Impl: p.Impl})
	return nil
}

func (AccessSyncerPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &accessSyncerGRPC{client: NewAccessProviderSyncServiceClient(c)}, nil
}

// AccessSyncerName constant should not be used directly when implementing plugins.
// It's the registration name for the data access syncer plugin,
// used by the CLI and the cli-plugin-base library (RegisterPlugins function) to register the plugins.
const AccessSyncerName = "accessSyncer"

type accessSyncerGRPC struct {
	client AccessProviderSyncServiceClient
}

func (g *accessSyncerGRPC) SyncFromTarget(ctx context.Context, config *AccessSyncFromTarget) (*AccessSyncResult, error) {
	return g.client.SyncFromTarget(ctx, config)
}

func (g *accessSyncerGRPC) SyncToTarget(ctx context.Context, config *AccessSyncToTarget) (*AccessSyncResult, error) {
	return g.client.SyncToTarget(ctx, config)
}

func (g *accessSyncerGRPC) SyncConfig(ctx context.Context) (*AccessSyncConfig, error) {
	return g.client.SyncConfig(ctx, &emptypb.Empty{})
}

func (g *accessSyncerGRPC) CliVersionInformation(ctx context.Context) (*version.CliBuildInformation, error) {
	return g.client.CliVersionInformation(ctx, &emptypb.Empty{})
}

type accessSyncerGRPCServer struct {
	UnimplementedAccessProviderSyncServiceServer

	Impl AccessSyncer
}

func (s *accessSyncerGRPCServer) SyncToTarget(ctx context.Context, config *AccessSyncToTarget) (*AccessSyncResult, error) {
	return s.Impl.SyncToTarget(ctx, config)
}

func (s *accessSyncerGRPCServer) SyncFromTarget(ctx context.Context, config *AccessSyncFromTarget) (*AccessSyncResult, error) {
	return s.Impl.SyncFromTarget(ctx, config)
}

func (s *accessSyncerGRPCServer) SyncConfig(ctx context.Context, _ *emptypb.Empty) (*AccessSyncConfig, error) {
	return s.Impl.SyncConfig(ctx)
}

func (s *accessSyncerGRPCServer) CliVersionInformation(ctx context.Context, _ *emptypb.Empty) (*version.CliBuildInformation, error) {
	return s.Impl.CliVersionInformation(ctx)
}

func WithSupportPartialSync() func(config *AccessSyncConfig) {
	return func(config *AccessSyncConfig) {
		config.SupportPartialSync = true
	}
}

func WithImplicitDeleteInAccessProviderUpdate() func(config *AccessSyncConfig) {
	return func(config *AccessSyncConfig) {
		config.ImplicitDeleteInAccessProviderUpdate = true
	}
}

type AccessSyncerVersionHandler struct {
}

func (h *AccessSyncerVersionHandler) CliVersionInformation(ctx context.Context) (*version.CliBuildInformation, error) {
	return version2.CreateSyncerCliBuildInformation(MinimalCliVersion), nil
}
