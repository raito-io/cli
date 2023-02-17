package version

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
)

const pluginCliMinimalVersion = "0.32.0"

var DevVersion = semver.New(0, 0, 0, "error", "noVersionDefined")

var version = *DevVersion
var date = ""

func SetVersion(setVersion, setDate string) {
	semverVersion, err := semver.NewVersion(setVersion)
	if err != nil {
		panic(err)
	}

	version = *semverVersion
	date = setDate
}

func GetVersionString() string {
	return version.String() + " (" + date + ")"
}

func GetCliVersion() *semver.Version {
	return &version
}

func CliPluginConstraint() *semver.Constraints {
	constraint, err := semver.NewConstraint(fmt.Sprintf("%s - %s", pluginCliMinimalVersion, GetCliVersion().String()))
	if err != nil {
		panic(err)
	}

	return constraint
}

//go:generate go run github.com/vektra/mockery/v2 --name=CliVersionInformation --with-expecter
type CliVersionInformation interface {
	GetCliVersion() *semver.Version
	CliPluginConstraint() *semver.Constraints
}

type CurrentCliVersionInformation struct{}

func (ci *CurrentCliVersionInformation) GetCliVersion() *semver.Version {
	return GetCliVersion()
}

func (ci *CurrentCliVersionInformation) CliPluginConstraint() *semver.Constraints {
	return CliPluginConstraint()
}
