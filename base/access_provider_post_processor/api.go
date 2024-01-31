package access_provider_post_processor

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/raito-io/cli/base/util/error/grpc_error"
	"github.com/raito-io/cli/base/util/version"
	"github.com/raito-io/cli/internal/version_management"
)

// AccessProviderPostProcessor interface needs to be implemented by any plugin (or CLI generally) that wants to post process access providers, read in sync.
// Post procession can alter the access provider based on for example tags.
type AccessProviderPostProcessor interface {
	version.CliVersionHandler

	PostProcessFromTarget(ctx context.Context, config *AccessProviderPostProcessorConfig) (*AccessProviderPostProcessorResult, error)
}

// AccessProviderPostProcessorPlugin is used on the server (CLI) and client (plugin) side to integrate with the plugin system.
// A plugin should not be using this directly, but instead depend on the cli-plugin-base library to register the plugins.
type AccessProviderPostProcessorPlugin struct {
	plugin.Plugin

	Impl AccessProviderPostProcessor
}

func (p *AccessProviderPostProcessorPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	RegisterAccessProviderPostProcessorServiceServer(s, &accessProviderPostProcessorGRPCServer{Impl: p.Impl})
	return nil
}

func (AccessProviderPostProcessorPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &accessProviderPostProcessorGRPC{client: NewAccessProviderPostProcessorServiceClient(c)}, nil
}

// AccessProviderPostProcessorName constant should not be used directly when implementing plugins.
// It's the registration name for the access provider post processor plugin,
// used by the CLI and the cli-plugin-base library (RegisterPlugins function) to register the plugins.
const AccessProviderPostProcessorName = "accessProviderPostProcessor"

type accessProviderPostProcessorGRPC struct {
	client AccessProviderPostProcessorServiceClient
}

func (g *accessProviderPostProcessorGRPC) PostProcessFromTarget(ctx context.Context, config *AccessProviderPostProcessorConfig) (*AccessProviderPostProcessorResult, error) {
	return grpc_error.ParseErrorResult(g.client.PostProcessFromTarget(ctx, config))
}

func (g *accessProviderPostProcessorGRPC) CliVersionInformation(ctx context.Context) (*version.CliBuildInformation, error) {
	return grpc_error.ParseErrorResult(g.client.CliVersionInformation(ctx, &emptypb.Empty{}))
}

type accessProviderPostProcessorGRPCServer struct {
	UnimplementedAccessProviderPostProcessorServiceServer

	Impl AccessProviderPostProcessor
}

func (s *accessProviderPostProcessorGRPCServer) PostProcessFromTarget(ctx context.Context, config *AccessProviderPostProcessorConfig) (_ *AccessProviderPostProcessorResult, err error) {
	defer func() {
		err = grpc_error.GrpcDeferErrorHandling(err)
	}()

	return s.Impl.PostProcessFromTarget(ctx, config)
}

func (s *accessProviderPostProcessorGRPCServer) CliVersionInformation(ctx context.Context, _ *emptypb.Empty) (_ *version.CliBuildInformation, err error) {
	defer func() {
		err = grpc_error.GrpcDeferErrorHandling(err)
	}()

	return s.Impl.CliVersionInformation(ctx)
}

type AccessProviderPostProcessorVersionHandler struct {
}

func (h *AccessProviderPostProcessorVersionHandler) CliVersionInformation(ctx context.Context) (*version.CliBuildInformation, error) {
	return version_management.CreateSyncerCliBuildInformation(MinimalCliVersion), nil
}
