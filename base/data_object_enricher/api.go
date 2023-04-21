package data_object_enricher

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/raito-io/cli/base/util/error/grpc_error"
	"github.com/raito-io/cli/base/util/version"
	"github.com/raito-io/cli/internal/version_management"
)

// DataObjectEnricher interface needs to be implemented by any plugin that wants to enrich data objects, read by the data source sync.
// Enrichment will add additional metadata to the data objects, coming from different systems.
type DataObjectEnricher interface {
	version.CliVersionHandler

	Enrich(ctx context.Context, config *DataObjectEnricherConfig) (*DataObjectEnricherResult, error)
}

// DataObjectEnricherPlugin is used on the server (CLI) and client (plugin) side to integrate with the plugin system.
// A plugin should not be using this directly, but instead depend on the cli-plugin-base library to register the plugins.
type DataObjectEnricherPlugin struct {
	plugin.Plugin

	Impl DataObjectEnricher
}

func (p *DataObjectEnricherPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	RegisterDataObjectEnricherServiceServer(s, &dataObjectEnricherGRPCServer{Impl: p.Impl})
	return nil
}

func (DataObjectEnricherPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &dataObjectEnricherGRPC{client: NewDataObjectEnricherServiceClient(c)}, nil
}

// DataObjectEnricherName constant should not be used directly when implementing plugins.
// It's the registration name for the data object enricher plugin,
// used by the CLI and the cli-plugin-base library (RegisterPlugins function) to register the plugins.
const DataObjectEnricherName = "dataObjectEnricher"

type dataObjectEnricherGRPC struct {
	client DataObjectEnricherServiceClient
}

func (g *dataObjectEnricherGRPC) Enrich(ctx context.Context, config *DataObjectEnricherConfig) (*DataObjectEnricherResult, error) {
	return grpc_error.ParseErrorResult(g.client.Enrich(ctx, config))
}

func (g *dataObjectEnricherGRPC) CliVersionInformation(ctx context.Context) (*version.CliBuildInformation, error) {
	return grpc_error.ParseErrorResult(g.client.CliVersionInformation(ctx, &emptypb.Empty{}))
}

type dataObjectEnricherGRPCServer struct {
	UnimplementedDataObjectEnricherServiceServer

	Impl DataObjectEnricher
}

func (s *dataObjectEnricherGRPCServer) Enrich(ctx context.Context, config *DataObjectEnricherConfig) (_ *DataObjectEnricherResult, err error) {
	defer func() {
		err = grpc_error.GrpcDeferErrorHandling(err)
	}()

	return s.Impl.Enrich(ctx, config)
}

func (s *dataObjectEnricherGRPCServer) CliVersionInformation(ctx context.Context, _ *emptypb.Empty) (_ *version.CliBuildInformation, err error) {
	defer func() {
		err = grpc_error.GrpcDeferErrorHandling(err)
	}()

	return s.Impl.CliVersionInformation(ctx)
}

type DataObjectEnricherVersionHandler struct {
}

func (h *DataObjectEnricherVersionHandler) CliVersionInformation(ctx context.Context) (*version.CliBuildInformation, error) {
	return version_management.CreateSyncerCliBuildInformation(MinimalCliVersion), nil
}
