package error

import (
	"fmt"
	"runtime"
	"strings"
)

type RecoverError struct {
	internalError error
	stackTrace    string
}

func NewRecoverError(err error) RecoverError {
	stackTrace := getStackTrace()

	return RecoverError{
		internalError: err,
		stackTrace:    stackTrace,
	}
}

func NewRecoverErrorf(format string, args ...interface{}) RecoverError {
	return NewRecoverError(fmt.Errorf(format, args...))
}

func (re RecoverError) Error() string {
	return fmt.Sprintf("Recovered: %s\n%s", re.internalError.Error(), re.stackTrace)
}

func (re RecoverError) Unwrap() error {
	return re.internalError
}

func (re RecoverError) StackTrace() string {
	return re.stackTrace
}

func getStackTrace() string {
	var stackTrace strings.Builder

	pc := make([]uintptr, 10)   // Show max 10 stack frames
	n := runtime.Callers(4, pc) // Skip 4 stack frames

	if n < 1 {
		return ""
	}

	frames := runtime.CallersFrames(pc)
	addNewLine := false

	for {
		frame, more := frames.Next()
		if !more {
			break
		}

		if addNewLine {
			stackTrace.WriteString("\n")
		}

		stackTrace.WriteString(fmt.Sprintf("\t[%s:%d] %s", frame.File, frame.Line, frame.Function))
		addNewLine = true
	}

	return stackTrace.String()
}
