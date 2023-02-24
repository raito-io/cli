package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/raito-io/cli/base/util/version"
)

func TestPluginInfo(t *testing.T) {
	i := PluginInfo{
		Name:        "TestPlugin",
		Version:     &version.SemVer{Major: 1, Minor: 2, Patch: 3},
		Description: "Plugin Description!",
		Parameters: []*ParameterInfo{
			{Name: "p1", Description: "p1 descr", Mandatory: true},
			{Name: "p2", Description: "p2 descr", Mandatory: false},
		},
	}

	is := i.InfoString()
	assert.Equal(t, "TestPlugin v1.2.3", is)
}
