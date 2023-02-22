package data_source

import "github.com/Masterminds/semver/v3"

var (
	MinimalCliVersion = semver.MustParse("0.32.0")
	supportedFeatures []string
)
