package plugin

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// TODO
//func (i Version) String() string {
//	return fmt.Sprintf("%d.%d.%d", i.Major, i.Minor, i.Maintenance)
//}

// ParseVersion parses
// the given string version in the form X.Y.Z and returns a Version struct representing it.
// If the input string is invalid, a 0.0.0 version will be returned
func ParseVersion(version string) *Version {
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return nil
	}
	major, err := strconv.Atoi(parts[0])

	if err != nil {
		return nil
	}
	minor, err := strconv.Atoi(parts[1])

	if err != nil {
		return nil
	}
	maintenance, err := strconv.Atoi(parts[2])

	if err != nil {
		return nil
	}

	return &Version{Major: int32(major), Minor: int32(minor), Maintenance: int32(maintenance)} //nolint:gosec
}

func (i *PluginInfo) InfoString() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s v%d.%d.%d", i.Name, i.Version.Major, i.Version.Minor, i.Version.Maintenance))

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

	Impl InfoServer
}

func (p *InfoPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	RegisterInfoServer(s, &infoGRPCServer{Impl: p.Impl})
	return nil
}

func (InfoPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &infoGRPC{client: NewInfoClient(c)}, nil
}

// InfoName constant should not be used directly when implementing plugins.
// It's the registration name for the info plugin,
// used by the CLI and the cli-plugin-base library (RegisterPlugins function) to register the plugins.
const InfoName = "info"

type infoGRPC struct{ client InfoClient }

func (g *infoGRPC) GetInfo(ctx context.Context) (*PluginInfo, error) {
	resp, err := g.client.GetInfo(ctx, &Empty{})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type infoGRPCServer struct {
	UnimplementedInfoServer

	// This is the real implementation
	Impl InfoServer
}

func (s *infoGRPCServer) GetInfo(ctx context.Context, in *Empty) (*PluginInfo, error) {
	return s.Impl.GetInfo(ctx, in)
}
