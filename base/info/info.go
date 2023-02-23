package info

import (
	"context"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	"github.com/raito-io/cli/base/util/plugin"
)

type InfoImpl struct {
	plugin.UnimplementedInfoServiceServer

	Info *plugin.PluginInfo
}

func (i *InfoImpl) GetInfo(context.Context, *emptypb.Empty) (*plugin.PluginInfo, error) {
	return i.Info, nil
}
