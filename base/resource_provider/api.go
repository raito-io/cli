package resource_provider

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/raito-io/cli/base/util/error/grpc_error"
	"github.com/raito-io/cli/base/util/version"
	"github.com/raito-io/cli/internal/version_management"
)

// ResourceProviderSyncer interface needs to be implemented by any plugin that wants to initialize resources on raito Cloud
type ResourceProviderSyncer interface {
	version.CliVersionHandler

	UpdateResources(ctx context.Context, config *UpdateResourceInput) (*UpdateResourceResult, error)
}

// ResourceProviderSyncerPlugin is used on the server (CLI) and client (plugin) side to integrate with the plugin system.
// A plugin should not be using this directly, but instead depend on the cli-plugin-base library to register the plugins.
type ResourceProviderSyncerPlugin struct {
	plugin.Plugin

	Impl ResourceProviderSyncer
}

func (p *ResourceProviderSyncerPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	RegisterResourceProviderServiceServer(s, &resourceProviderSyncerGRPCServer{Impl: p.Impl})
	return nil
}

func (p *ResourceProviderSyncerPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &resourceProviderSyncerGRPC{client: NewResourceProviderServiceClient(c)}, nil
}

const ResourceProviderSyncerName = "resourceProviderSyncer"

type resourceProviderSyncerGRPCServer struct {
	UnimplementedResourceProviderServiceServer

	Impl ResourceProviderSyncer
}

func (s *resourceProviderSyncerGRPCServer) UpdateResources(ctx context.Context, config *UpdateResourceInput) (_ *UpdateResourceResult, err error) {
	defer func() {
		err = grpc_error.GrpcDeferErrorHandling(err)
	}()

	return s.Impl.UpdateResources(ctx, config)
}

func (s *resourceProviderSyncerGRPCServer) CliVersionInformation(ctx context.Context, _ *emptypb.Empty) (_ *version.CliBuildInformation, err error) {
	defer func() {
		err = grpc_error.GrpcDeferErrorHandling(err)
	}()

	return s.Impl.CliVersionInformation(ctx)
}

type resourceProviderSyncerGRPC struct {
	client ResourceProviderServiceClient
}

func (g *resourceProviderSyncerGRPC) UpdateResources(ctx context.Context, config *UpdateResourceInput) (*UpdateResourceResult, error) {
	return grpc_error.ParseErrorResult(g.client.UpdateResources(ctx, config))
}

func (g *resourceProviderSyncerGRPC) CliVersionInformation(ctx context.Context) (*version.CliBuildInformation, error) {
	return grpc_error.ParseErrorResult(g.client.CliVersionInformation(ctx, &emptypb.Empty{}))
}

type ResourceProviderSyncerVersionHandler struct {}

func (h *ResourceProviderSyncerVersionHandler) CliVersionInformation(ctx context.Context) (*version.CliBuildInformation, error) {
	return version_management.CreateSyncerCliBuildInformation(MinimalCliVersion), nil
}