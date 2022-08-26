package data_usage

import (
	"github.com/raito-io/cli/base/util/config"
	error2 "github.com/raito-io/cli/base/util/error"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// DataUsageSyncConfig represents the configuration that is passed from the CLI to the DataUsageSyncer plugin interface.
// It contains all the necessary configuration parameters for the plugin to function.
type DataUsageSyncConfig struct {
	config.ConfigMap
	TargetFile string
}

// DataUsageSyncResult represents the result from the data usage sync process.
// A potential error is also modeled in here so specific errors remain intact when passed over RPC.
type DataUsageSyncResult struct {
	Error *error2.ErrorResult
}

// DataUsageSyncer interface needs to be implemented by any plugin that wants to import data usage information
// into a Raito data source.
type DataUsageSyncer interface {
	SyncDataUsage(config *DataUsageSyncConfig) DataUsageSyncResult
}

// DataUsageSyncerPlugin is used on the server (CLI) and client (plugin) side to integrate with the plugin system.
// A plugin should not be using this directly, but instead depend on the cli-plugin-base library to register the plugins.
type DataUsageSyncerPlugin struct {
	Impl DataUsageSyncer
}

func (p *DataUsageSyncerPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &dataUsageSyncerRPCServer{Impl: p.Impl}, nil
}

func (DataUsageSyncerPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &dataUsageSyncerRPC{client: c}, nil
}

// DataUsageSyncerName constant should not be used directly when implementing plugins.
// It's the registration name for the data usage syncer plugin,
// used by the CLI and the cli-plugin-base library (RegisterPlugins function) to register the plugins.
const DataUsageSyncerName = "dataUsageSyncer"

type dataUsageSyncerRPC struct{ client *rpc.Client }

func (g *dataUsageSyncerRPC) SyncDataUsage(config *DataUsageSyncConfig) DataUsageSyncResult {
	var resp DataUsageSyncResult

	err := g.client.Call("Plugin.SyncDataUsage", config, &resp)
	if err != nil && resp.Error == nil {
		resp.Error = error2.ToErrorResult(err)
	}

	return resp
}

type dataUsageSyncerRPCServer struct {
	Impl DataUsageSyncer
}

func (s *dataUsageSyncerRPCServer) SyncDataUsage(config *DataUsageSyncConfig, resp *DataUsageSyncResult) error {
	*resp = s.Impl.SyncDataUsage(config)
	return nil
}
