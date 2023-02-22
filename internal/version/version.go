package version

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/hashicorp/go-hclog"
	"github.com/raito-io/golang-set/set"

	version2 "github.com/raito-io/cli/base/util/version"
)

var v = ""

var DevVersion = semver.New(0, 0, 0, "error", "noVersionDefined")

var version = *DevVersion
var date = ""

func SetVersion(setVersion, setDate string) {
	semverVersion := semver.MustParse(setVersion)

	version = *semverVersion
	date = setDate
}

func GetVersionString() string {
	return "v" + version.String() + " (" + date + ")"
}

func GetCliVersion() *semver.Version {
	return &version
}

func CreateSyncerCliBuildInformation(minimalCliVersion *semver.Version, supportedFeatures ...string) *version2.CliBuildInformation {
	return &version2.CliBuildInformation{
		CliBuildVersion:   version2.ToSemVer(GetCliVersion()),
		CliMinimalVersion: version2.ToSemVer(minimalCliVersion),
		SupportedFeatures: supportedFeatures,
	}
}

func IsValidToSync(ctx context.Context, plugin version2.CliVersionHandler, syncerMinimalVersion *semver.Version) (bool, set.Set[string], error) {
	pluginInformation, err := plugin.CliVersionInformation(ctx)
	if err != nil {
		return false, nil, err
	}

	return isValidToSync(pluginInformation, syncerMinimalVersion, GetCliVersion)
}

func isValidToSync(pluginInformation *version2.CliBuildInformation, syncerMinimalVersion *semver.Version, cliInfo func() *semver.Version) (bool, set.Set[string], error) {
	currentCliVersion := cliInfo()
	supportedFeaturesFn := func() set.Set[string] { return set.NewSet(pluginInformation.SupportedFeatures...) }

	if currentCliVersion.Equal(DevVersion) {
		hclog.L().Warn("Running in dev mode, skipping version check")

		return true, supportedFeaturesFn(), nil
	}

	pluginCliCurrentVersion := pluginInformation.CliBuildVersion.ToVersion()
	pluginMinimalVersion := pluginInformation.CliMinimalVersion.ToVersion()

	if pluginCliCurrentVersion.LessThan(currentCliVersion) {
		//CLI version is newer than plugin version
		if syncerMinimalVersion.GreaterThan(pluginCliCurrentVersion) {
			return false, nil, IncompatibleVersionError{
				pluginVersion: pluginCliCurrentVersion.String(),
				cliVersion:    currentCliVersion.String(),
				updatePlugin:  true,
			}
		} else {
			return true, supportedFeaturesFn(), nil
		}
	} else {
		//CLI version is older than plugin version
		if pluginMinimalVersion.GreaterThan(currentCliVersion) {
			return false, nil, IncompatibleVersionError{
				pluginVersion: pluginCliCurrentVersion.String(),
				cliVersion:    currentCliVersion.String(),
				updatePlugin:  false,
			}
		} else {
			return true, supportedFeaturesFn(), nil
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
