package access_provider

import (
	"github.com/hashicorp/go-plugin"
	"github.com/raito-io/cli/base/util/config"
	error2 "github.com/raito-io/cli/base/util/error"

	//nolint:gosec
	"net/rpc"
)

// AccessSyncConfig represents the configuration that is passed from the CLI to the AccessSyncer plugin interface.
// It contains all the necessary configuration parameters for the plugin to function.
type AccessSyncConfig struct {
	config.ConfigMap
	// SourceFile points to the file containing the access controls that need to be pushed to the data source.
	SourceFile string
	// TargetFile points to the file where the plugin needs to export the access controls to that are read from the data source.
	TargetFile string
	Prefix     string
}

// AccessSyncResult represents the result from the data access sync process.
// A potential error is also modeled in here so specific errors remain intact when passed over RPC.
type AccessSyncResult struct {
	Error *error2.ErrorResult
}

// AccessSyncer interface needs to be implemented by any plugin that wants to sync access controls between Raito and the data source.
// This sync can be in the 2 directions or in just 1 depending on the parameters set in AccessSyncConfig.
type AccessSyncer interface {
	SyncAccess(config *AccessSyncConfig) AccessSyncResult
}

// AccessSyncerPlugin is used on the server (CLI) and client (plugin) side to integrate with the plugin system.
// A plugin should not be using this directly, but instead depend on the cli-plugin-base library to register the plugins.
type AccessSyncerPlugin struct {
	Impl AccessSyncer
}

func (p AccessSyncerPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &accessSyncerRPCServer{Impl: p.Impl}, nil
}

func (AccessSyncerPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &accessSyncerRPC{client: c}, nil
}

// AccessSyncerName constant should not be used directly when implementing plugins.
// It's the registration name for the data access syncer plugin,
// used by the CLI and the cli-plugin-base library (RegisterPlugins function) to register the plugins.
const AccessSyncerName = "accessSyncer"

type accessSyncerRPC struct{ client *rpc.Client }

func (g *accessSyncerRPC) SyncAccess(config *AccessSyncConfig) AccessSyncResult {
	var resp AccessSyncResult

	err := g.client.Call("Plugin.SyncAccess", config, &resp)
	if err != nil && resp.Error == nil {
		resp.Error = error2.ToErrorResult(err)
	}

	return resp
}

type accessSyncerRPCServer struct {
	Impl AccessSyncer
}

func (s *accessSyncerRPCServer) SyncAccess(config *AccessSyncConfig, resp *AccessSyncResult) error {
	*resp = s.Impl.SyncAccess(config)
	return nil
}
