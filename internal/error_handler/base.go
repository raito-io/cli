package error_handler

import "fmt"

type Placeholder string

const ErrorPlaceholder Placeholder = "ERROR"

type ErrorHandler interface {
	Error(e error)
	Errorf(format string, args ...any)
	GetError() error
	HasError() bool
}

type BaseErrorHandler struct {
	err error
}

func NewBaseErrorHandler() *BaseErrorHandler {
	return &BaseErrorHandler{}
}

func (eh *BaseErrorHandler) Error(e error) {
	eh.err = e
}

func (eh *BaseErrorHandler) Errorf(format string, args ...any) {
	eh.Error(fmt.Errorf(format, args...))
}

func (eh *BaseErrorHandler) GetError() error {
	return eh.err
}

func (eh *BaseErrorHandler) HasError() bool {
	return eh.err != nil
}
