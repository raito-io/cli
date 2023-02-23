package version

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/Masterminds/semver/v3"
	"github.com/hashicorp/go-hclog"

	version2 "github.com/raito-io/cli/base/util/version"
)

var DevVersion = semver.New(0, 0, 0, "error", "noVersionDefined")

var version = *DevVersion
var date = ""

var pluginCliVersion *semver.Version

func SetVersion(setVersion, setDate string) {
	if setVersion != "" {
		semverVersion := semver.MustParse(setVersion)

		version = *semverVersion
	}

	date = setDate
}

func GetVersionString() string {
	return "v" + version.String() + " (" + date + ")"
}

func GetCliVersion() *semver.Version {
	return &version
}

func getCliVersionInPlugin() *semver.Version {
	if pluginCliVersion == nil {
		bi, ok := debug.ReadBuildInfo()
		if !ok {
			return nil
		}

		for _, dep := range bi.Deps {
			if dep.Path == "github.com/raito-io/cli" {
				pluginCliVersion = semver.MustParse(dep.Version)

				break
			}
		}
	}

	return pluginCliVersion
}

func CreateSyncerCliBuildInformation(minimalCliVersion *semver.Version) *version2.CliBuildInformation {
	return &version2.CliBuildInformation{
		CliBuildVersion:   version2.ToSemVer(getCliVersionInPlugin()),
		CliMinimalVersion: version2.ToSemVer(minimalCliVersion),
	}
}

func IsValidToSync(ctx context.Context, plugin version2.CliVersionHandler, syncerMinimalVersion *semver.Version) (bool, error) {
	pluginInformation, err := plugin.CliVersionInformation(ctx)
	if err != nil {
		return false, err
	}

	return isValidToSync(pluginInformation, syncerMinimalVersion, GetCliVersion)
}

func isValidToSync(pluginInformation *version2.CliBuildInformation, syncerMinimalVersion *semver.Version, cliInfo func() *semver.Version) (bool, error) {
	currentCliVersion := cliInfo()
	if currentCliVersion == nil {
		return false, errors.New("could not get current cli version")
	}

	if currentCliVersion.Equal(DevVersion) {
		hclog.L().Warn("Running in dev mode, skipping version check")

		return true, nil
	}

	pluginCliCurrentVersion := pluginInformation.CliBuildVersion.ToVersion()
	pluginMinimalVersion := pluginInformation.CliMinimalVersion.ToVersion()

	if pluginCliCurrentVersion.LessThan(currentCliVersion) {
		//CLI version is newer than plugin version
		if syncerMinimalVersion.GreaterThan(pluginCliCurrentVersion) {
			return false, IncompatibleVersionError{
				pluginVersion: pluginCliCurrentVersion.String(),
				cliVersion:    currentCliVersion.String(),
				updatePlugin:  true,
			}
		} else {
			return true, nil
		}
	} else {
		//CLI version is older than plugin version
		if pluginMinimalVersion.GreaterThan(currentCliVersion) {
			return false, IncompatibleVersionError{
				pluginVersion: pluginCliCurrentVersion.String(),
				cliVersion:    currentCliVersion.String(),
				updatePlugin:  false,
			}
		} else {
			return true, nil
		}
	}
}

type IncompatibleVersionError struct {
	pluginVersion string
	cliVersion    string
	updatePlugin  bool
}

func (e IncompatibleVersionError) Error() string {
	solution := ""
	if e.updatePlugin {
		solution = "Please update the plugin to the latest version."
	} else {
		solution = "Please update the CLI to the latest version."
	}

	return fmt.Sprintf("Plugin is incompatible with CLI version. Plugin is build with CLI version '%s' while current CLI version is '%s', %s", e.pluginVersion, e.cliVersion, solution)
}
