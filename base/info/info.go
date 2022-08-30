package info

import (
	"github.com/raito-io/cli/base/util/plugin"
)

type InfoImpl struct {
	Info plugin.PluginInfo
}

func (i *InfoImpl) PluginInfo() plugin.PluginInfo {
	return i.Info
}
