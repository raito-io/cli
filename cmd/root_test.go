package cmd

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRootCmd(t *testing.T) {
	mem := &exitMemory{}
	Execute("v1.2.3", []string{ "-h" }, mem.Exit)
	assert.Equal(t, 0, mem.code)
}

func TestRootCmdNoArgs(t *testing.T) {
	mem := &exitMemory{}
	Execute("v1.2.3", []string{  }, mem.Exit)
	assert.Equal(t, 0, mem.code)
}

func TestRootCmdConfigFile(t *testing.T) {
	mem := &exitMemory{}
	Execute("v1.2.3", []string{ "--config-file", "testdata/config1.yml" }, mem.Exit)
	assert.Equal(t, 0, mem.code)
}

func TestRootCmdConfigFileNotExists(t *testing.T) {
	mem := &exitMemory{}
	Execute("v1.2.3", []string{ "run --config-file", "testdata/wrong.yml" }, mem.Exit)
	assert.Equal(t, 1, mem.code)
}

func TestRootCmdVersion(t *testing.T) {
	var b bytes.Buffer
	mem := &exitMemory{}
	cmd := newRootCmd("v1.2.3", mem.Exit)
	cmd.cmd.SetOut(&b)
	cmd.cmd.SetArgs([]string{ "-v" })
	err := cmd.cmd.Execute()
	assert.Nil(t, err)
	assert.Equal(t, "raito version v1.2.3\n", b.String())
	assert.Equal(t, 0, mem.code)
}
