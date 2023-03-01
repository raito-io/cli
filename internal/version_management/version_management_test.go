package version_management

import (
	"fmt"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"

	version2 "github.com/raito-io/cli/base/util/version"
)

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
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Cli > Plugin: Compatible",
			args: args{
				pluginInformation: &version2.CliBuildInformation{
					CliBuildVersion:   &version2.SemVer{Major: 1, Minor: 4, Patch: 2},
					CliMinimalVersion: &version2.SemVer{Major: 1, Minor: 0, Patch: 0},
				},
				syncerMinimalVersion: semver.New(2, 0, 4, "", ""),
				cliInfo: func() *semver.Version {
					return semver.New(1, 0, 0, "", "")
				},
			},
			want:    true,
			wantErr: assert.NoError,
		},
		{
			name: "Cli = Plugin: Compatible",
			args: args{
				pluginInformation: &version2.CliBuildInformation{
					CliBuildVersion:   &version2.SemVer{Major: 2, Minor: 0, Patch: 4},
					CliMinimalVersion: &version2.SemVer{Major: 1, Minor: 0, Patch: 0},
				},
				syncerMinimalVersion: semver.New(2, 0, 4, "", ""),
				cliInfo: func() *semver.Version {
					return semver.New(1, 0, 0, "", "")
				},
			},
			want:    true,
			wantErr: assert.NoError,
		},
		{
			name: "Cli < Plugin: Compatible",
			args: args{
				pluginInformation: &version2.CliBuildInformation{
					CliBuildVersion:   &version2.SemVer{Major: 2, Minor: 0, Patch: 4},
					CliMinimalVersion: &version2.SemVer{Major: 1, Minor: 0, Patch: 0},
				},
				syncerMinimalVersion: semver.New(1, 0, 0, "", ""),
				cliInfo: func() *semver.Version {
					return semver.New(1, 4, 2, "", "")
				},
			},
			want:    true,
			wantErr: assert.NoError,
		},
		{
			name: "Cli > Plugin: Not Compatible",
			args: args{
				pluginInformation: &version2.CliBuildInformation{
					CliBuildVersion:   &version2.SemVer{Major: 1, Minor: 4, Patch: 2},
					CliMinimalVersion: &version2.SemVer{Major: 1, Minor: 0, Patch: 0},
				},
				syncerMinimalVersion: semver.New(2, 0, 4, "", ""),
				cliInfo: func() *semver.Version {
					return semver.New(2, 0, 0, "", "")
				},
			},
			want: false,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorAs(t, err, &IncompatiblePluginVersionError{})
			},
		},
		{
			name: "Cli < Plugin: Not Compatible",
			args: args{
				pluginInformation: &version2.CliBuildInformation{
					CliBuildVersion:   &version2.SemVer{Major: 2, Minor: 0, Patch: 4},
					CliMinimalVersion: &version2.SemVer{Major: 2, Minor: 0, Patch: 0},
				},
				syncerMinimalVersion: semver.New(1, 4, 2, "", ""),
				cliInfo: func() *semver.Version {
					return semver.New(1, 0, 0, "", "")
				},
			},
			want: false,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorAs(t, err, &IncompatiblePluginVersionError{})
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isValidToSync(tt.args.pluginInformation, tt.args.syncerMinimalVersion, tt.args.cliInfo)
			if !tt.wantErr(t, err, fmt.Sprintf("isValidToSync(%v, %v, %v)", tt.args.pluginInformation, tt.args.syncerMinimalVersion, tt.args.cliInfo())) {
				return
			}
			assert.Equalf(t, tt.want, got, "isValidToSync(%v, %v, %v)", tt.args.pluginInformation, tt.args.syncerMinimalVersion, tt.args.cliInfo())
		})
	}
}
