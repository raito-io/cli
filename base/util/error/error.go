package error

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MissingInputParameterError struct {
	parameter string
}

func (e *MissingInputParameterError) Error() string {
	return fmt.Sprintf("mandatory parameter %q is missing", e.parameter)
}

// CreateMissingInputParameterError is a helper method to create a consistent error result across plugins when a mandatory parameter is missing.
func CreateMissingInputParameterError(parameter string) error {
	return &MissingInputParameterError{parameter: parameter}
}

type BadInputParameterError struct {
	parameter   string
	value       string
	explanation string
}

func (e *BadInputParameterError) Error() string {
	if e.explanation != "" {
		return fmt.Sprintf("parameter %q has invalid value %q. %s", e.parameter, e.value, e.explanation)
	}

	return fmt.Sprintf("parameter %q has invalid value %q", e.parameter, e.value)
}

// CreateBadInputParameterError is a helper method to create a consistent error result across plugins when an input parameter has an unexpected format.
func CreateBadInputParameterError(parameter, value, explanation string) error {
	return &BadInputParameterError{parameter: parameter, value: value, explanation: explanation}
}

type SourceConnectionError struct {
	url     string
	message string
}

func (e *SourceConnectionError) Error() string {
	return fmt.Sprintf("error while connecting to %q: %s", e.url, e.message)
}

// CreateSourceConnectionError is a helper method to create a consistent error result across plugins when there is a connection problem to the data source or identity store.
func CreateSourceConnectionError(url, message string) error {
	return &SourceConnectionError{url: url, message: message}
}

type CreateFileError struct {
	filename     string
	wrappedError error
}

func (e *CreateFileError) Error() string {
	return fmt.Sprintf("error creating temporary file %q for data source importer: %s", e.filename, e.wrappedError.Error())
}

func (e *CreateFileError) Unwrap() error {
	return e.wrappedError
}

func CreateErrorFileError(filename string, err error) error {
	return &CreateFileError{filename: filename, wrappedError: err}
}

// ToErrorResult is a helper method to create an ErrorResult from an error. If the error already is of type ErrorResult, the original is returned.
func ToErrorResult(err error) error {
	return status.Error(codes.Internal, err.Error())
}
