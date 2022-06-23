package api

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToErrorResult(t *testing.T) {
	res := ToErrorResult(errors.New("An error"))
	assert.Equal(t, UnknownError, res.ErrorCode)
	assert.Equal(t, "An error", res.ErrorMessage)
}

func TestToErrorResultAsIs(t *testing.T) {
	res := ToErrorResult(ErrorResult{
		ErrorMessage: "Existing",
		ErrorCode:    BadInputParameterError,
	})
	assert.Equal(t, BadInputParameterError, res.ErrorCode)
	assert.Equal(t, "Existing", res.ErrorMessage)
}

func TestToErrorResultAsPointer(t *testing.T) {
	res := ToErrorResult(&ErrorResult{
		ErrorMessage: "Pointer",
		ErrorCode:    MissingInputParameterError,
	})
	assert.Equal(t, MissingInputParameterError, res.ErrorCode)
	assert.Equal(t, "Pointer", res.ErrorMessage)
}

func TestCreateSourceConnectionError(t *testing.T) {
	res := CreateSourceConnectionError("http://myurl", "failed")
	assert.Equal(t, SourceConnectionError, res.ErrorCode)
	assert.Contains(t, res.ErrorMessage, "http://myurl")
	assert.Contains(t, res.ErrorMessage, "connecting")
	assert.Contains(t, res.ErrorMessage, "failed")
}

func TestCreateBadInputParameterError(t *testing.T) {
	res := CreateBadInputParameterError("param1", "v1", "explained")
	assert.Equal(t, BadInputParameterError, res.ErrorCode)
	assert.Equal(t, "parameter \"param1\" has invalid value \"v1\". explained", res.ErrorMessage)
}

func TestCreateBadInputParameterErrorNoExplanation(t *testing.T) {
	res := CreateBadInputParameterError("param2", "v2", "")
	assert.Equal(t, BadInputParameterError, res.ErrorCode)
	assert.Contains(t, res.ErrorMessage, "param2")
	assert.Contains(t, res.ErrorMessage, "v2")
}

func TestCreateMissingInputParameterError(t *testing.T) {
	res := CreateMissingInputParameterError("param666")
	assert.Equal(t, MissingInputParameterError, res.ErrorCode)
	assert.Contains(t, res.ErrorMessage, "missing")
	assert.Contains(t, res.ErrorMessage, "param666")
}
