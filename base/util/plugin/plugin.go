package plugin

import (
	"context"
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/raito-io/cli/base/util/version"
)

// ParseVersion parses
// the given string version in the form X.Y.Z and returns a Version struct representing it.
// If the input string is invalid, a 0.0.0 version will be returned
func ParseVersion(stringVersion string) *version.SemVer {
	sv := semver.MustParse(stringVersion)

	return &version.SemVer{Major: sv.Major(), Minor: sv.Minor(), Patch: sv.Patch(), Prerelease: sv.Prerelease(), Build: sv.Metadata()}
}

func (i *PluginInfo) InfoString() string {
	var sb strings.Builder

	sb.WriteString(i.Name)
	sb.WriteString(" ")

	if i.GetVersion() != nil {
		sv := semver.New(i.GetVersion().GetMajor(), i.GetVersion().GetMinor(), i.GetVersion().GetPatch(), i.GetVersion().GetPrerelease(), i.GetVersion().GetBuild())

		sb.WriteString("v")
		sb.WriteString(sv.String())
	} else {
		sb.WriteString("unknown version")
	}

	return sb.String()
}

func (i *PluginInfo) FullOverview() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s v%s", i.Name, i.Version))

	if i.Description != "" {
		sb.WriteString(fmt.Sprintf("\n\n%s", i.Description))
	}

	if len(i.Parameters) > 0 {
		sb.WriteString("\n\nParameters:")

		for _, param := range i.Parameters {
			sb.WriteString(fmt.Sprintf("\n   %s", param))
		}
	}

	return sb.String()
}

// Info interface needs to be implemented by all plugins to provide basic plugin information.
type Info interface {
	GetInfo(ctx context.Context) (*PluginInfo, error)
}

// InfoPlugin is used on the server (CLI) and client (plugin) side to integrate with the plugin system.
// A plugin should not be using this directly, but instead depend on the cli-plugin-base library to register the plugins.
type InfoPlugin struct {
	plugin.Plugin

	Impl InfoServiceServer
}

func (p *InfoPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	RegisterInfoServiceServer(s, &infoGRPCServer{Impl: p.Impl})
	return nil
}

func (InfoPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &infoGRPC{client: NewInfoServiceClient(c)}, nil
}

// InfoName constant should not be used directly when implementing plugins.
// It's the registration name for the info plugin,
// used by the CLI and the cli-plugin-base library (RegisterPlugins function) to register the plugins.
const InfoName = "info"

type infoGRPC struct{ client InfoServiceClient }

func (g *infoGRPC) GetInfo(ctx context.Context) (*PluginInfo, error) {
	resp, err := g.client.GetInfo(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type infoGRPCServer struct {
	UnimplementedInfoServiceServer

	// This is the real implementation
	Impl InfoServiceServer
}

func (s *infoGRPCServer) GetInfo(ctx context.Context, in *emptypb.Empty) (*PluginInfo, error) {
	return s.Impl.GetInfo(ctx, in)
}
