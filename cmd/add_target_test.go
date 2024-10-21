package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToStringValue(t *testing.T) {
	assert.False(t, mustBeQuoted("123"))
	assert.False(t, mustBeQuoted("a123b"))
	assert.False(t, mustBeQuoted("A123"))
	assert.False(t, mustBeQuoted("aBc"))

	assert.True(t, mustBeQuoted(" "))
	assert.True(t, mustBeQuoted("123 ok"))
	assert.True(t, mustBeQuoted("special#char"))
}
