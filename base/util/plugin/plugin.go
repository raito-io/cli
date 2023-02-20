package plugin

import (
	"fmt"
	"net/rpc"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/hashicorp/go-plugin"
)

// Version contains semantic versioning information of the plugin
type Version struct {
	Major       int
	Minor       int
	Maintenance int
}

func (i Version) String() string {
	return fmt.Sprintf("%d.%d.%d", i.Major, i.Minor, i.Maintenance)
}

// ParseVersion parses
// the given string version in the form X.Y.Z and returns a Version struct representing it.
// If the input string is invalid, a 0.0.0 version will be returned
func ParseVersion(version string) Version {
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return Version{}
	}
	major, err := strconv.Atoi(parts[0])

	if err != nil {
		return Version{}
	}
	minor, err := strconv.Atoi(parts[1])

	if err != nil {
		return Version{}
	}
	maintenance, err := strconv.Atoi(parts[2])

	if err != nil {
		return Version{}
	}

	return Version{Major: major, Minor: minor, Maintenance: maintenance}
}

// ParameterInfo contains the information about a parameter.
// This is used to inform the CLI user what command-line parameters are expected explicitly for this target (plugin).
type ParameterInfo struct {
	Name        string
	Description string
	Mandatory   bool
}

func (i ParameterInfo) String() string {
	if i.Mandatory {
		return fmt.Sprintf("%s (mandatory): %s", i.Name, i.Description)
	}

	return fmt.Sprintf("%s (optional): %s", i.Name, i.Description)
}

// PluginInfo represents the information about a plugin.
type PluginInfo struct {
	Name        string
	Description string
	Version     Version
	Parameters  []ParameterInfo
}

func (i PluginInfo) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s v%s", i.Name, i.Version))

	return sb.String()
}

func (i PluginInfo) FullOverview() string {
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
//
//go:generate go run github.com/vektra/mockery/v2 --name=Info --with-expecter
type Info interface {
	PluginInfo() PluginInfo
	CliBuildVersion() semver.Version
	CliMinimalVersion() semver.Version
}

// InfoPlugin is used on the server (CLI) and client (plugin) side to integrate with the plugin system.
// A plugin should not be using this directly, but instead depend on the cli-plugin-base library to register the plugins.
type InfoPlugin struct {
	Impl Info
}

func (p *InfoPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &infoRPCServer{Impl: p.Impl}, nil
}

func (InfoPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &infoRPC{client: c}, nil
}

// InfoName constant should not be used directly when implementing plugins.
// It's the registration name for the info plugin,
// used by the CLI and the cli-plugin-base library (RegisterPlugins function) to register the plugins.
const InfoName = "info"

type infoRPC struct{ client *rpc.Client }

func (g *infoRPC) PluginInfo() PluginInfo {
	var resp PluginInfo

	err := g.client.Call("Plugin.PluginInfo", new(interface{}), &resp)
	if err != nil {
		return PluginInfo{}
	}

	return resp
}

func (g *infoRPC) CliBuildVersion() semver.Version {
	var resp semver.Version

	err := g.client.Call("Plugin.CliBuildVersion", new(interface{}), &resp)
	if err != nil {
		return semver.Version{}
	}

	return resp
}

func (g *infoRPC) CliMinimalVersion() semver.Version {
	var resp semver.Version

	err := g.client.Call("Plugin.PluginCliConstraint", new(interface{}), &resp)
	if err != nil {
		return semver.Version{}
	}

	return resp
}

type infoRPCServer struct {
	Impl Info
}

func (s *infoRPCServer) PluginInfo(args interface{}, resp *PluginInfo) error {
	*resp = s.Impl.PluginInfo()
	return nil
}

func (s *infoRPCServer) CliBuildVersion(args interface{}, resp *semver.Version) error {
	*resp = s.Impl.CliBuildVersion()
	return nil
}

func (s *infoRPCServer) CliMinimalVersion(args interface{}, resp *semver.Version) error {
	*resp = s.Impl.CliMinimalVersion()
	return nil
}
