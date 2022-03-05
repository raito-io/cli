package info

import (
	"github.com/raito-io/cli/common/api"
)

type InfoImpl struct {
	Info api.PluginInfo
}

func (i *InfoImpl) PluginInfo() api.PluginInfo {
	return i.Info
}
