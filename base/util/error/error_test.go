package error

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToErrorResult(t *testing.T) {
	res := ToErrorResult(errors.New("An error"))
	assert.Equal(t, ErrorCode_UNKNOWN_ERROR, res.ErrorCode)
	assert.Equal(t, "An error", res.ErrorMessage)
}

func TestToErrorResultAsIs(t *testing.T) {
	res := ToErrorResult(ErrorResult{
		ErrorMessage: "Existing",
		ErrorCode:    ErrorCode_BAD_INPUT_PARAMETER_ERROR,
	})
	assert.Equal(t, ErrorCode_BAD_INPUT_PARAMETER_ERROR, res.ErrorCode)
	assert.Equal(t, "Existing", res.ErrorMessage)
}

func TestToErrorResultAsPointer(t *testing.T) {
	res := ToErrorResult(&ErrorResult{
		ErrorMessage: "Pointer",
		ErrorCode:    ErrorCode_MISSING_INPUT_PARAMETER_ERROR,
	})
	assert.Equal(t, ErrorCode_MISSING_INPUT_PARAMETER_ERROR, res.ErrorCode)
	assert.Equal(t, "Pointer", res.ErrorMessage)
}

func TestCreateSourceConnectionError(t *testing.T) {
	res := CreateSourceConnectionError("http://myurl", "failed")
	assert.Equal(t, ErrorCode_SOURCE_CONNECTION_ERROR, res.ErrorCode)
	assert.Contains(t, res.ErrorMessage, "http://myurl")
	assert.Contains(t, res.ErrorMessage, "connecting")
	assert.Contains(t, res.ErrorMessage, "failed")
}

func TestCreateBadInputParameterError(t *testing.T) {
	res := CreateBadInputParameterError("param1", "v1", "explained")
	assert.Equal(t, ErrorCode_BAD_INPUT_PARAMETER_ERROR, res.ErrorCode)
	assert.Equal(t, "parameter \"param1\" has invalid value \"v1\". explained", res.ErrorMessage)
}

func TestCreateBadInputParameterErrorNoExplanation(t *testing.T) {
	res := CreateBadInputParameterError("param2", "v2", "")
	assert.Equal(t, ErrorCode_BAD_INPUT_PARAMETER_ERROR, res.ErrorCode)
	assert.Contains(t, res.ErrorMessage, "param2")
	assert.Contains(t, res.ErrorMessage, "v2")
}

func TestCreateMissingInputParameterError(t *testing.T) {
	res := CreateMissingInputParameterError("param666")
	assert.Equal(t, ErrorCode_MISSING_INPUT_PARAMETER_ERROR, res.ErrorCode)
	assert.Contains(t, res.ErrorMessage, "missing")
	assert.Contains(t, res.ErrorMessage, "param666")
}
