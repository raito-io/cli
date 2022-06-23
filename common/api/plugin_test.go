package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPluginInfo(t *testing.T) {
	i := PluginInfo{
		Name:        "TestPlugin",
		Version:     Version{1, 2, 3},
		Description: "Plugin Description!",
		Parameters: []ParameterInfo{
			{"p1", "p1 descr", true},
			{"p2", "p2 descr", false},
		},
	}

	is := i.String()
	assert.Equal(t, "TestPlugin v1.2.3", is)
}
