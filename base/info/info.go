package info

import (
	"context"

	"github.com/raito-io/cli/base/util/plugin"
)

type InfoImpl struct {
	plugin.UnimplementedInfoServer

	Info *plugin.PluginInfo
}

func (i *InfoImpl) GetInfo(context.Context, *plugin.Empty) (*plugin.PluginInfo, error) {
	return i.Info, nil
}
