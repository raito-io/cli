package error

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func _createRecoverError(format string, args ...any) RecoverError {
	return NewRecoverErrorf(format, args...)
}

func TestRecoverError_Error(t *testing.T) {
	re := _createRecoverError("Recovered: test error")

	assert.Error(t, re)

	assert.Equal(t, 3, strings.Count(re.Error(), "\n"))

	t.Log(re.Error())
}
