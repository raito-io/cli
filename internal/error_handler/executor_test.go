package error_handler

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorExecutor(t *testing.T) {
	// Given
	executed := false

	baseError := NewBaseErrorHandler()
	executor := OnError(baseError, func(e error) {
		executed = true
	})

	assert.False(t, executor.HasError())

	// When
	executor.Error(errors.New("BOOM"))

	// Then
	assert.True(t, executor.HasError())
	assert.True(t, executed)
}
