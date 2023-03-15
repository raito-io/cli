package grpc_error

import (
	"errors"
	"fmt"
	"reflect"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	error2 "github.com/raito-io/cli/base/util/error"
)

type InternalPluginStatusError struct {
	statusCode codes.Code
	errorMsg   string
}

func (e *InternalPluginStatusError) Error() string {
	return fmt.Sprintf("%s (%s)", e.errorMsg, e.statusCode.String())
}

func (e *InternalPluginStatusError) StatusCode() codes.Code {
	return e.statusCode
}

var errorStatusMap = map[reflect.Type]codes.Code{
	reflect.TypeOf((*error2.MissingInputParameterError)(nil)): codes.InvalidArgument,
	reflect.TypeOf((*error2.BadInputParameterError)(nil)):     codes.InvalidArgument,
	reflect.TypeOf((*error2.SourceConnectionError)(nil)):      codes.Unavailable,
	reflect.TypeOf((*error2.CreateFileError)(nil)):            codes.Internal,
}

func ToStatusError(err error) error {
	msg := err.Error()

	statusCode := codes.Internal

	unwrappedErr := err

	for {
		if unwrappedErr == nil {
			break
		}

		if code, ok := errorStatusMap[reflect.TypeOf(unwrappedErr)]; ok {
			statusCode = code
			break
		}

		unwrappedErr = errors.Unwrap(unwrappedErr)
	}

	return status.Error(statusCode, msg)
}

func FromStatusError(err error) error {
	if err == nil {
		return nil
	}

	if s, ok := status.FromError(err); ok {
		return &InternalPluginStatusError{
			statusCode: s.Code(),
			errorMsg:   s.Message(),
		}
	}

	return err
}

func ParseErrorResult[T any](result T, err error) (T, error) {
	err = FromStatusError(err)

	return result, err
}

func GrpcDeferErrorHandling(err error) error {
	if r := recover(); r != nil {
		err = status.Errorf(codes.Unknown, "recover after panic: %v", r)
	}

	if err != nil {
		err = ToStatusError(err)
	}

	return err
}
