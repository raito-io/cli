package error_handler

import "fmt"

type ErrorHandlerWrapper struct {
	eh     ErrorHandler
	err    error
	format string
	args   []any
}

func Wrap(eh ErrorHandler, format string, args ...any) *ErrorHandlerWrapper {
	return &ErrorHandlerWrapper{
		eh:     eh,
		format: format,
		args:   args,
	}
}

func (ehw *ErrorHandlerWrapper) Error(e error) {
	ehw.err = e
	ehw.eh.Error(ehw.wrapError(e))
}

func (ehw *ErrorHandlerWrapper) Errorf(format string, args ...any) {
	ehw.eh.Errorf(ehw.format, ehw.args...)
}

func (ehw *ErrorHandlerWrapper) GetError() error {
	return ehw.err
}

func (ehw *ErrorHandlerWrapper) HasError() bool {
	return ehw.err != nil
}

func (ehw *ErrorHandlerWrapper) wrapError(e error) error {
	args := make([]any, 0, len(ehw.args))

	for _, arg := range ehw.args {
		if arg == ErrorPlaceholder {
			args = append(args, e)
		} else {
			args = append(args, arg)
		}
	}

	return fmt.Errorf(ehw.format, args...)
}
