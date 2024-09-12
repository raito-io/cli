package error_handler

type ErrorExecutorHandler struct {
	eh ErrorHandler
	f  func(e error)
}

func OnError(eh ErrorHandler, f func(e error)) *ErrorExecutorHandler {
	return &ErrorExecutorHandler{
		eh: eh,
		f:  f,
	}
}

func (eeh *ErrorExecutorHandler) Error(e error) {
	eeh.f(e)
	eeh.eh.Error(e)
}

func (eeh *ErrorExecutorHandler) Errorf(format string, args ...any) {
	eeh.eh.Errorf(format, args...)
}

func (eeh *ErrorExecutorHandler) GetError() error {
	return eeh.eh.GetError()
}

func (eeh *ErrorExecutorHandler) HasError() bool {
	return eeh.eh.HasError()
}
