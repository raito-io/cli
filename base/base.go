// Package base provides some helper functionalities that should be used by every plugin.
package base

import (
	"errors"
	"sync"

	"github.com/raito-io/cli/base/data_object_enricher"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/raito-io/cli/base/access_provider"
	"github.com/raito-io/cli/base/data_source"
	"github.com/raito-io/cli/base/data_usage"
	"github.com/raito-io/cli/base/identity_store"
	plugin2 "github.com/raito-io/cli/base/util/plugin"
)

var logger hclog.Logger
var onlyOnce sync.Once

// Logger creates a new logger that should be used as a basis for all logging in the plugin.
// So it's advised to call this method first and store the logger in a (global) variable.
func Logger() hclog.Logger {
	onlyOnce.Do(func() {
		logger = hclog.New(&hclog.LoggerOptions{
			JSONFormat: true,
		})
	})

	return logger
}

func buildPluginMap(pluginImpls ...interface{}) (plugin.PluginSet, error) {
	var pluginMap = plugin.PluginSet{}

	infoFound := false

	for _, plugin := range pluginImpls {
		switch p := plugin.(type) {
		case identity_store.IdentityStoreSyncer:
			if _, f := pluginMap[identity_store.IdentityStoreSyncerName]; f {
				return nil, errors.New("multiple implementations for IdentityStoreSyncer Plugin found. There should be only one")
			}
			pluginMap[identity_store.IdentityStoreSyncerName] = &identity_store.IdentityStoreSyncerPlugin{Impl: p}

			logger.Debug("Registered IdentityStoreSyncer Plugin")
		case data_source.DataSourceSyncer:
			if _, f := pluginMap[data_source.DataSourceSyncerName]; f {
				return nil, errors.New("multiple implementations for DataSourceSyncer Plugin found. There should be only one")
			}
			pluginMap[data_source.DataSourceSyncerName] = &data_source.DataSourceSyncerPlugin{Impl: p}

			logger.Debug("Registered DataSourceSyncer Plugin")
		case access_provider.AccessSyncer:
			if _, f := pluginMap[access_provider.AccessSyncerName]; f {
				return nil, errors.New("multiple implementations for AccessSyncer Plugin found. There should be only one")
			}
			pluginMap[access_provider.AccessSyncerName] = &access_provider.AccessSyncerPlugin{Impl: p}

			logger.Debug("Registered AccessSyncer Plugin")
		case data_usage.DataUsageSyncer:
			if _, f := pluginMap[data_usage.DataUsageSyncerName]; f {
				return nil, errors.New("multiple implementations for DataUsageSyncer Plugin found. There should be only one")
			}
			pluginMap[data_usage.DataUsageSyncerName] = &data_usage.DataUsageSyncerPlugin{Impl: p}

			logger.Debug("Registered DataUsageSyncer Plugin")
		case plugin2.InfoServiceServer:
			if _, f := pluginMap[plugin2.InfoName]; f {
				return nil, errors.New("multiple implementation for Info Plugin found. There should be only one")
			}
			pluginMap[plugin2.InfoName] = &plugin2.InfoPlugin{Impl: p}

			logger.Debug("Registered Info Plugin")

			infoFound = true
		case data_object_enricher.DataObjectEnricher:
			if _, f := pluginMap[data_object_enricher.DataObjectEnricherName]; f {
				return nil, errors.New("multiple implementations for DataObjectEnricher Plugin found. There should be only one")
			}
			pluginMap[data_object_enricher.DataObjectEnricherName] = &data_object_enricher.DataObjectEnricherPlugin{Impl: p}

			logger.Debug("Registered DataObjectEnricher Plugin")
		}
	}

	if len(pluginMap) == 0 {
		return nil, errors.New("no plugin implementations found")
	}

	if !infoFound {
		return nil, errors.New("no info plugin implementation found. This infoPlugin mandatory")
	}

	return pluginMap, nil
}

// RegisterPlugins takes a list of objects that implement the different plugin API interfaces.
// It will automatically detect which of the interfaces are implemented and will register them as plugins.
// This way, the underlying plugin system infoPlugin abstracted away for everybody implementing plugins for the Raito CLI.
func RegisterPlugins(pluginImpls ...interface{}) error {
	Logger()

	pluginMap, err := buildPluginMap(pluginImpls...)
	if err != nil {
		return err
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		GRPCServer:      plugin.DefaultGRPCServer,
	})

	return nil
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "RAITO_CLI_PLUGIN",
	MagicCookieValue: "Raito Handshake!",
}
