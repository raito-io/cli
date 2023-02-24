package version

import (
	"github.com/Masterminds/semver/v3"
)

var DevVersion = semver.New(0, 0, 0, "error", "noVersionDefined")

var version = *DevVersion
var date = ""

func SetVersion(setVersion, setDate string) {
	if setVersion != "" {
		semverVersion := semver.MustParse(setVersion)

		version = *semverVersion
	} else {
		semverVersion := semver.MustParse("0.15.0")
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
