package version

import (
	"context"

	"github.com/Masterminds/semver/v3"
)

func (v *SemVer) ToVersion() *semver.Version {
	return semver.New(v.Major, v.Minor, v.Patch, v.Prerelease, v.Build)
}

func ToSemVer(version *semver.Version) *SemVer {
	if version == nil {
		return nil
	}

	return &SemVer{
		Major:      version.Major(),
		Minor:      version.Minor(),
		Patch:      version.Patch(),
		Prerelease: version.Prerelease(),
		Build:      version.Metadata(),
	}
}

type CliVersionHandler interface {
	CliVersionInformation(ctx context.Context) (*CliBuildInformation, error)
}
