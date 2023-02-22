package access_provider

import "github.com/Masterminds/semver/v3"

var (
	MinimalCliVersion = semver.MustParse("0.32.0")
	supportedFeatures []string
)
