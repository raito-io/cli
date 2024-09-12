package error_handler

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapError(t *testing.T) {
	// Given
	baseError := NewBaseErrorHandler()
	wrapper := Wrap(baseError, "This is an %s with arguments %s: %w, %d", "ERROR", "testArgs", ErrorPlaceholder, 23)

	// When
	wrapper.Error(errors.New("BOOM"))

	// Then
	assert.True(t, wrapper.HasError())
	assert.True(t, baseError.HasError())
	assert.EqualError(t, wrapper.GetError(), "BOOM")
	assert.EqualError(t, baseError.GetError(), "This is an ERROR with arguments testArgs: BOOM, 23")
}
