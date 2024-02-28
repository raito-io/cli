package tag

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/raito-io/cli/base/util/error/grpc_error"
	"github.com/raito-io/cli/base/util/version"
	"github.com/raito-io/cli/internal/version_management"
)

type TagSyncer interface {
	version.CliVersionHandler

	SyncTags(ctx context.Context, config *TagSyncConfig) (*TagSyncResult, error)
}

type TagSyncerPlugin struct {
	plugin.Plugin

	Impl TagSyncer
}

func (p *TagSyncerPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	RegisterTagSyncServiceServer(s, &tagSyncerGRPCServer{Impl: p.Impl})
	return nil
}

func (p *TagSyncerPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &tagSyncerGRPC{client: NewTagSyncServiceClient(c)}, nil
}

const TagSyncerName = "tagSyncer"

type tagSyncerGRPC struct {
	client TagSyncServiceClient
}

func (g *tagSyncerGRPC) SyncTags(ctx context.Context, config *TagSyncConfig) (*TagSyncResult, error) {
	return grpc_error.ParseErrorResult(g.client.SyncTags(ctx, config))
}

func (g *tagSyncerGRPC) CliVersionInformation(ctx context.Context) (*version.CliBuildInformation, error) {
	return grpc_error.ParseErrorResult(g.client.CliVersionInformation(ctx, &emptypb.Empty{}))
}

type tagSyncerGRPCServer struct {
	UnimplementedTagSyncServiceServer

	Impl TagSyncer
}

func (s *tagSyncerGRPCServer) SyncTags(ctx context.Context, config *TagSyncConfig) (_ *TagSyncResult, err error) {
	defer func() {
		err = grpc_error.GrpcDeferErrorHandling(err)
	}()

	return s.Impl.SyncTags(ctx, config)
}

func (s *tagSyncerGRPCServer) CliVersionInformation(ctx context.Context, _ *emptypb.Empty) (_ *version.CliBuildInformation, err error) {
	defer func() {
		err = grpc_error.GrpcDeferErrorHandling(err)
	}()

	return s.Impl.CliVersionInformation(ctx)
}

type TagSyncerVersionHandler struct{}

func (h *TagSyncerVersionHandler) CliVersionInformation(_ context.Context) (*version.CliBuildInformation, error) {
	return version_management.CreateSyncerCliBuildInformation(MinimalCliVersion), nil
}
