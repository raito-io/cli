package access_provider

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/raito-io/cli/base/util/config"
	error2 "github.com/raito-io/cli/base/util/error"

	"net/rpc"
)

// AccessSyncToTarget contains all necessary configuration parameters to export Data from Raito into DS
type AccessSyncToTarget struct {
	config.ConfigMap
	// SourceFile points to the file containing the access controls that need to be pushed to the data source.
	SourceFile string
	// FeedbackTargetFile points to the file where the plugin needs to export the access controls feedback to.
	FeedbackTargetFile string
	Prefix             string
}

// AccessSyncFromTarget contains all necessary configuration parameters to import Data from Raito into DS
type AccessSyncFromTarget struct {
	config.ConfigMap
	// TargetFile points to the file where the plugin needs to export the access control naming.
	TargetFile string
	Prefix     string
}

// AccessSyncResult represents the result from the data access sync process.
// A potential error is also modeled in here so specific errors remain intact when passed over RPC.
type AccessSyncResult struct {
	Error *error2.ErrorResult
}

// AccessSyncConfig gives us information on how the CLI can sync access providers
type AccessSyncConfig struct {
	// SupportPartialSync if true, syncing only out of sync access providers is allowed
	SupportPartialSync bool

	// ImplicitDeleteInAccessProviderUpdate if true, access providers can be deleted by name only
	ImplicitDeleteInAccessProviderUpdate bool
}

// AccessSyncer interface needs to be implemented by any plugin that wants to sync access controls between Raito and the data source.
// This sync can be in the 2 directions or in just 1 depending on the parameters set in AccessSyncConfig.
type AccessSyncer interface {
	SyncFromTarget(config *AccessSyncFromTarget) AccessSyncResult
	SyncToTarget(config *AccessSyncToTarget) AccessSyncResult

	SyncConfig() AccessSyncConfig
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

func (g *accessSyncerRPC) SyncFromTarget(config *AccessSyncFromTarget) AccessSyncResult {
	var resp AccessSyncResult

	err := g.client.Call("Plugin.SyncFromTarget", config, &resp)
	if err != nil && resp.Error == nil {
		resp.Error = error2.ToErrorResult(err)
	}

	return resp
}

func (g *accessSyncerRPC) SyncToTarget(config *AccessSyncToTarget) AccessSyncResult {
	var resp AccessSyncResult

	err := g.client.Call("Plugin.SyncToTarget", config, &resp)
	if err != nil && resp.Error == nil {
		resp.Error = error2.ToErrorResult(err)
	}

	return resp
}

func (g *accessSyncerRPC) SyncConfig() AccessSyncConfig {
	var resp AccessSyncConfig

	err := g.client.Call("Plugin.SyncConfig", new(interface{}), &resp)
	if err != nil {
		hclog.L().Warn(fmt.Sprintf("Failed to load sync config from plugin. Will use default settings. %s", err.Error()))
		return AccessSyncConfig{}
	}

	return resp
}

type accessSyncerRPCServer struct {
	Impl AccessSyncer
}

func (s *accessSyncerRPCServer) SyncToTarget(config *AccessSyncToTarget, resp *AccessSyncResult) error {
	*resp = s.Impl.SyncToTarget(config)
	return nil
}

func (s *accessSyncerRPCServer) SyncFromTarget(config *AccessSyncFromTarget, resp *AccessSyncResult) error {
	*resp = s.Impl.SyncFromTarget(config)
	return nil
}

func (s *accessSyncerRPCServer) SyncConfig(args interface{}, resp *AccessSyncConfig) error {
	*resp = s.Impl.SyncConfig()
	return nil
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
