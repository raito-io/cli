package version

import (
	"github.com/Masterminds/semver/v3"
)

var pluginCliMinimalVersion = semver.MustParse("0.32.0")

var DevVersion = semver.New(0, 0, 0, "error", "noVersionDefined")

var version = *DevVersion
var date = ""

func SetVersion(setVersion, setDate string) {
	semverVersion := semver.MustParse(setVersion)

	version = *semverVersion
	date = setDate
}

func GetVersionString() string {
	return version.String() + " (" + date + ")"
}

func GetCliVersion() *semver.Version {
	return &version
}

func GetMinimalCliVersion() *semver.Version {
	return pluginCliMinimalVersion
}

//go:generate go run github.com/vektra/mockery/v2 --name=CliVersionInformation --with-expecter
type CliVersionInformation interface {
	GetCliVersion() *semver.Version
	GetCliMinimalCompatibleVersion() *semver.Version
}

type CurrentCliVersionInformation struct{}

func (ci *CurrentCliVersionInformation) GetCliVersion() *semver.Version {
	return GetCliVersion()
}

func (ci *CurrentCliVersionInformation) GetCliMinimalCompatibleVersion() *semver.Version {
	return GetMinimalCliVersion()
}
