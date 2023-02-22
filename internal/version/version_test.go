package version

import (
	"fmt"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/raito-io/golang-set/set"
	"github.com/stretchr/testify/assert"

	version2 "github.com/raito-io/cli/base/util/version"
)

func TestVersion(t *testing.T) {
	SetVersion("v1.2.3", "yyyy-mm-dd")
	assert.Equal(t, "v1.2.3 (yyyy-mm-dd)", GetVersionString())
}

func Test_isValidToSync(t *testing.T) {
	type args struct {
		pluginInformation    *version2.CliBuildInformation
		syncerMinimalVersion *semver.Version
		cliInfo              func() *semver.Version
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		want1   set.Set[string]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Cli > Plugin: Compatible",
			args: args{
				pluginInformation: &version2.CliBuildInformation{
					CliBuildVersion:   &version2.SemVer{Major: 1, Minor: 4, Patch: 2},
					CliMinimalVersion: &version2.SemVer{Major: 1, Minor: 0, Patch: 0},
					SupportedFeatures: []string{"feature1", "feature2"},
				},
				syncerMinimalVersion: semver.New(2, 0, 4, "", ""),
				cliInfo: func() *semver.Version {
					return semver.New(1, 0, 0, "", "")
				},
			},
			want:    true,
			want1:   set.NewSet("feature1", "feature2"),
			wantErr: assert.NoError,
		},
		{
			name: "Cli = Plugin: Compatible",
			args: args{
				pluginInformation: &version2.CliBuildInformation{
					CliBuildVersion:   &version2.SemVer{Major: 2, Minor: 0, Patch: 4},
					CliMinimalVersion: &version2.SemVer{Major: 1, Minor: 0, Patch: 0},
					SupportedFeatures: []string{"feature1", "feature3"},
				},
				syncerMinimalVersion: semver.New(2, 0, 4, "", ""),
				cliInfo: func() *semver.Version {
					return semver.New(1, 0, 0, "", "")
				},
			},
			want:    true,
			want1:   set.NewSet("feature1", "feature3"),
			wantErr: assert.NoError,
		},
		{
			name: "Cli < Plugin: Compatible",
			args: args{
				pluginInformation: &version2.CliBuildInformation{
					CliBuildVersion:   &version2.SemVer{Major: 2, Minor: 0, Patch: 4},
					CliMinimalVersion: &version2.SemVer{Major: 1, Minor: 0, Patch: 0},
					SupportedFeatures: []string{"feature1", "feature3"},
				},
				syncerMinimalVersion: semver.New(1, 0, 0, "", ""),
				cliInfo: func() *semver.Version {
					return semver.New(1, 4, 2, "", "")
				},
			},
			want:    true,
			want1:   set.NewSet("feature1", "feature3"),
			wantErr: assert.NoError,
		},
		{
			name: "Cli > Plugin: Not Compatible",
			args: args{
				pluginInformation: &version2.CliBuildInformation{
					CliBuildVersion:   &version2.SemVer{Major: 1, Minor: 4, Patch: 2},
					CliMinimalVersion: &version2.SemVer{Major: 1, Minor: 0, Patch: 0},
					SupportedFeatures: []string{"feature1", "feature3"},
				},
				syncerMinimalVersion: semver.New(2, 0, 4, "", ""),
				cliInfo: func() *semver.Version {
					return semver.New(2, 0, 0, "", "")
				},
			},
			want:  false,
			want1: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorAs(t, err, &IncompatibleVersionError{})
			},
		},
		{
			name: "Cli < Plugin: Not Compatible",
			args: args{
				pluginInformation: &version2.CliBuildInformation{
					CliBuildVersion:   &version2.SemVer{Major: 2, Minor: 0, Patch: 4},
					CliMinimalVersion: &version2.SemVer{Major: 2, Minor: 0, Patch: 0},
					SupportedFeatures: []string{"feature1", "feature3"},
				},
				syncerMinimalVersion: semver.New(1, 4, 2, "", ""),
				cliInfo: func() *semver.Version {
					return semver.New(1, 0, 0, "", "")
				},
			},
			want:  false,
			want1: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorAs(t, err, &IncompatibleVersionError{})
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := isValidToSync(tt.args.pluginInformation, tt.args.syncerMinimalVersion, tt.args.cliInfo)
			if !tt.wantErr(t, err, fmt.Sprintf("isValidToSync(%v, %v, %v)", tt.args.pluginInformation, tt.args.syncerMinimalVersion, tt.args.cliInfo())) {
				return
			}
			assert.Equalf(t, tt.want, got, "isValidToSync(%v, %v, %v)", tt.args.pluginInformation, tt.args.syncerMinimalVersion, tt.args.cliInfo())
			assert.Equalf(t, tt.want1, got1, "isValidToSync(%v, %v, %v)", tt.args.pluginInformation, tt.args.syncerMinimalVersion, tt.args.cliInfo())
		})
	}
}
