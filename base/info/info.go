package info

import (
	"github.com/Masterminds/semver/v3"

	"github.com/raito-io/cli/base/util/plugin"
	"github.com/raito-io/cli/internal/version"
)

type InfoImpl struct {
	Info plugin.PluginInfo
}

func (i *InfoImpl) PluginInfo() plugin.PluginInfo {
	return i.Info
}

func (i *InfoImpl) CliBuildVersion() semver.Version {
	return *version.GetCliVersion()
}

func (i *InfoImpl) PluginCliConstraint() semver.Constraints {
	return *version.CliPluginConstraint()
}
