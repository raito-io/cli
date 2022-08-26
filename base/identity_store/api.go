package identity_store

import (
	"github.com/raito-io/cli/base/util/config"
	error2 "github.com/raito-io/cli/base/util/error"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// IdentityStoreSyncConfig represents the configuration that is passed from the CLI to the IdentityStoreSyncer plugin interface.
// It contains all the necessary configuration parameters for the plugin to function.
type IdentityStoreSyncConfig struct {
	config.ConfigMap
	UserFile  string
	GroupFile string
}

// IdentityStoreSyncResult represents the result from the identity store sync process.
// A potential error is also modeled in here so specific errors remain intact when passed over RPC.
type IdentityStoreSyncResult struct {
	Error *error2.ErrorResult
}

// IdentityStoreSyncer interface needs to be implemented by any plugin that wants to import users and groups into a Raito identity store.
type IdentityStoreSyncer interface {
	SyncIdentityStore(config *IdentityStoreSyncConfig) IdentityStoreSyncResult
}

// IdentityStoreSyncerPlugin is used on the server (CLI) and client (plugin) side to integrate with the plugin system.
// A plugin should not be using this directly, but instead depend on the cli-plugin-base library to register the plugins.
type IdentityStoreSyncerPlugin struct {
	Impl IdentityStoreSyncer
}

func (p *IdentityStoreSyncerPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &identityStoreSyncerRPCServer{Impl: p.Impl}, nil
}

func (IdentityStoreSyncerPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &identityStoreSyncerRPC{client: c}, nil
}

// IdentityStoreSyncerName constant should not be used directly when implementing plugins.
// It's the registration name for the identity store syncer plugin,
// used by the CLI and the cli-plugin-base library (RegisterPlugins function) to register the plugins.
const IdentityStoreSyncerName = "identityStoreSyncer"

type identityStoreSyncerRPC struct{ client *rpc.Client }

func (g *identityStoreSyncerRPC) SyncIdentityStore(config *IdentityStoreSyncConfig) IdentityStoreSyncResult {
	var resp IdentityStoreSyncResult

	err := g.client.Call("Plugin.SyncIdentityStore", config, &resp)
	if err != nil && resp.Error == nil {
		resp.Error = error2.ToErrorResult(err)
	}

	return resp
}

type identityStoreSyncerRPCServer struct {
	Impl IdentityStoreSyncer
}

func (s *identityStoreSyncerRPCServer) SyncIdentityStore(config *IdentityStoreSyncConfig, resp *IdentityStoreSyncResult) error {
	*resp = s.Impl.SyncIdentityStore(config)

	return nil
}
