package error

import "fmt"

type ErrorCode int

const (
	NoError ErrorCode = iota
	UnknownError
	BadInputParameterError
	MissingInputParameterError
	SourceConnectionError
)

// ErrorResult represents an error returned by an API call.
type ErrorResult struct {
	ErrorMessage string
	ErrorCode    ErrorCode
}

func (e ErrorResult) Error() string {
	return e.ErrorMessage
}

// CreateMissingInputParameterError is a helper method to create a consistent error result across plugins when a mandatory parameter is missing.
func CreateMissingInputParameterError(parameter string) *ErrorResult {
	return &ErrorResult{
		ErrorMessage: fmt.Sprintf("mandatory parameter %q is missing", parameter),
		ErrorCode:    MissingInputParameterError,
	}
}

// CreateBadInputParameterError is a helper method to create a consistent error result across plugins when an input parameter has an unexpected format.
func CreateBadInputParameterError(parameter, value, explanation string) *ErrorResult {
	var msg string
	if explanation != "" {
		msg = fmt.Sprintf("parameter %q has invalid value %q. %s", parameter, value, explanation)
	} else {
		msg = fmt.Sprintf("parameter %q has invalid value %q", parameter, value)
	}

	return &ErrorResult{
		ErrorMessage: msg,
		ErrorCode:    BadInputParameterError,
	}
}

// CreateSourceConnectionError is a helper method to create a consistent error result across plugins when there is a connection problem to the data source or identity store.
func CreateSourceConnectionError(url, message string) *ErrorResult {
	return &ErrorResult{
		ErrorMessage: fmt.Sprintf("error while connecting to %q: %s", url, message),
		ErrorCode:    SourceConnectionError,
	}
}

// ToErrorResult is a helper method to to create an ErrorResult from an error. If the error already is of type ErrorResult, the original is returned.
func ToErrorResult(err error) *ErrorResult {
	if res, ok := err.(ErrorResult); ok {
		return &res
	}

	if res, ok := err.(*ErrorResult); ok {
		return res
	}

	return &ErrorResult{ErrorMessage: err.Error(), ErrorCode: UnknownError}
}
