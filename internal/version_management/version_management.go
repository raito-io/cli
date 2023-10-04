package version_management

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/Masterminds/semver/v3"
	"github.com/hashicorp/go-hclog"

	version2 "github.com/raito-io/cli/base/util/version"
	"github.com/raito-io/cli/internal/graphql"
	"github.com/raito-io/cli/internal/target/types"
	"github.com/raito-io/cli/internal/version"
)

var pluginCliVersion *semver.Version

type Compatibility int

const (
	CompatibilityUnknown Compatibility = iota
	NotSupported
	Deprecated
	Supported
)

type CompatibilityInformation struct {
	Compatibility        Compatibility
	DeprecatedWarningMsg *string
	SupportedVersions    string
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

	return isValidToSync(pluginInformation, syncerMinimalVersion, version.GetCliVersion)
}

func isValidToSync(pluginInformation *version2.CliBuildInformation, syncerMinimalVersion *semver.Version, cliInfo func() *semver.Version) (bool, error) {
	currentCliVersion := cliInfo()
	if currentCliVersion == nil {
		return false, errors.New("could not get current cli version")
	}

	if currentCliVersion.Equal(version.DevVersion) {
		hclog.L().Warn("Running in dev mode, skipping version check")

		return true, nil
	}

	pluginCliCurrentVersion := pluginInformation.CliBuildVersion.ToVersion()
	pluginMinimalVersion := pluginInformation.CliMinimalVersion.ToVersion()

	if pluginCliCurrentVersion.LessThan(currentCliVersion) {
		//CLI version is newer than plugin version
		if syncerMinimalVersion.GreaterThan(pluginCliCurrentVersion) {
			return false, IncompatiblePluginVersionError{
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
			return false, IncompatiblePluginVersionError{
				pluginVersion: pluginCliCurrentVersion.String(),
				cliVersion:    currentCliVersion.String(),
				updatePlugin:  false,
			}
		} else {
			return true, nil
		}
	}
}

type IncompatiblePluginVersionError struct {
	pluginVersion string
	cliVersion    string
	updatePlugin  bool
}

func (e IncompatiblePluginVersionError) Error() string {
	solution := ""
	if e.updatePlugin {
		solution = "Please update the plugin to the latest version."
	} else {
		solution = "Please update the CLI to the latest version."
	}

	return fmt.Sprintf("Plugin is incompatible with CLI version. Plugin is build with CLI version '%s' while current CLI version is '%s', %s", e.pluginVersion, e.cliVersion, solution)
}

func IsCompatibleWithRaitoCloud(config *types.BaseConfig) (CompatibilityInformation, error) {
	if version.GetCliVersion().Equal(version.DevVersion) {
		return CompatibilityInformation{
			Compatibility: Supported,
		}, nil
	}

	supportedVersions, err := getCompatibleRaitoCloudVersions(config)
	if err != nil {
		return CompatibilityInformation{
			Compatibility: CompatibilityUnknown,
		}, err
	}

	if supportedVersions.SupportedVersions.Check(version.GetCliVersion()) {
		return CompatibilityInformation{
			Compatibility:     Supported,
			SupportedVersions: supportedVersions.SupportedVersions.String(),
		}, nil
	}

	if supportedVersions.DeprecatedVersions != nil && supportedVersions.DeprecatedVersions.Check(version.GetCliVersion()) {
		return CompatibilityInformation{
			Compatibility:        Deprecated,
			SupportedVersions:    supportedVersions.SupportedVersions.String(),
			DeprecatedWarningMsg: supportedVersions.DeprecatedVersionMsg,
		}, nil
	}

	return CompatibilityInformation{
		Compatibility:     NotSupported,
		SupportedVersions: supportedVersions.SupportedVersions.String(),
	}, nil
}

type SupportedRaitoCloudVersions struct {
	SupportedVersions    *semver.Constraints
	DeprecatedVersions   *semver.Constraints
	DeprecatedVersionMsg *string
}

func getCompatibleRaitoCloudVersions(config *types.BaseConfig) (*SupportedRaitoCloudVersions, error) {
	var query struct {
		SupportedCliVersion struct {
			SupportedVersions  string
			DeprecatedVersions *struct {
				DeprecatedVersions string
				Msg                *string
			}
		} `graphql:"SupportedCLIVersion"`
	}

	err := graphql.NewClient(config).Query(context.Background(), &query, nil)
	if err != nil {
		return nil, fmt.Errorf("compatible Raito Cloud version: %w", err)
	}

	currentVersionConstraint, err := semver.NewConstraint(query.SupportedCliVersion.SupportedVersions)
	if err != nil {
		return nil, fmt.Errorf("compatible Raito Cloud version: %w", err)
	}

	result := &SupportedRaitoCloudVersions{
		SupportedVersions: currentVersionConstraint,
	}

	if query.SupportedCliVersion.DeprecatedVersions != nil {
		result.DeprecatedVersions, err = semver.NewConstraint(query.SupportedCliVersion.DeprecatedVersions.DeprecatedVersions)
		if err != nil {
			return nil, fmt.Errorf("compatible Raito Cloud version: %w", err)
		}

		result.DeprecatedVersionMsg = query.SupportedCliVersion.DeprecatedVersions.Msg
	}

	return result, nil
}
