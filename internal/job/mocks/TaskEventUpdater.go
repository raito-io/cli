// Code generated by mockery v2.50.1. DO NOT EDIT.

package mocks

import (
	context "context"

	job "github.com/raito-io/cli/internal/job"
	mock "github.com/stretchr/testify/mock"
)

// TaskEventUpdater is an autogenerated mock type for the TaskEventUpdater type
type TaskEventUpdater struct {
	mock.Mock
}

type TaskEventUpdater_Expecter struct {
	mock *mock.Mock
}

func (_m *TaskEventUpdater) EXPECT() *TaskEventUpdater_Expecter {
	return &TaskEventUpdater_Expecter{mock: &_m.Mock}
}

// GetSubtaskEventUpdater provides a mock function with given fields: subtask
func (_m *TaskEventUpdater) GetSubtaskEventUpdater(subtask string) job.SubtaskEventUpdater {
	ret := _m.Called(subtask)

	if len(ret) == 0 {
		panic("no return value specified for GetSubtaskEventUpdater")
	}

	var r0 job.SubtaskEventUpdater
	if rf, ok := ret.Get(0).(func(string) job.SubtaskEventUpdater); ok {
		r0 = rf(subtask)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(job.SubtaskEventUpdater)
		}
	}

	return r0
}

// TaskEventUpdater_GetSubtaskEventUpdater_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetSubtaskEventUpdater'
type TaskEventUpdater_GetSubtaskEventUpdater_Call struct {
	*mock.Call
}

// GetSubtaskEventUpdater is a helper method to define mock.On call
//   - subtask string
func (_e *TaskEventUpdater_Expecter) GetSubtaskEventUpdater(subtask interface{}) *TaskEventUpdater_GetSubtaskEventUpdater_Call {
	return &TaskEventUpdater_GetSubtaskEventUpdater_Call{Call: _e.mock.On("GetSubtaskEventUpdater", subtask)}
}

func (_c *TaskEventUpdater_GetSubtaskEventUpdater_Call) Run(run func(subtask string)) *TaskEventUpdater_GetSubtaskEventUpdater_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *TaskEventUpdater_GetSubtaskEventUpdater_Call) Return(_a0 job.SubtaskEventUpdater) *TaskEventUpdater_GetSubtaskEventUpdater_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TaskEventUpdater_GetSubtaskEventUpdater_Call) RunAndReturn(run func(string) job.SubtaskEventUpdater) *TaskEventUpdater_GetSubtaskEventUpdater_Call {
	_c.Call.Return(run)
	return _c
}

// SetStatusToCompleted provides a mock function with given fields: ctx, results
func (_m *TaskEventUpdater) SetStatusToCompleted(ctx context.Context, results []job.TaskResult) {
	_m.Called(ctx, results)
}

// TaskEventUpdater_SetStatusToCompleted_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetStatusToCompleted'
type TaskEventUpdater_SetStatusToCompleted_Call struct {
	*mock.Call
}

// SetStatusToCompleted is a helper method to define mock.On call
//   - ctx context.Context
//   - results []job.TaskResult
func (_e *TaskEventUpdater_Expecter) SetStatusToCompleted(ctx interface{}, results interface{}) *TaskEventUpdater_SetStatusToCompleted_Call {
	return &TaskEventUpdater_SetStatusToCompleted_Call{Call: _e.mock.On("SetStatusToCompleted", ctx, results)}
}

func (_c *TaskEventUpdater_SetStatusToCompleted_Call) Run(run func(ctx context.Context, results []job.TaskResult)) *TaskEventUpdater_SetStatusToCompleted_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]job.TaskResult))
	})
	return _c
}

func (_c *TaskEventUpdater_SetStatusToCompleted_Call) Return() *TaskEventUpdater_SetStatusToCompleted_Call {
	_c.Call.Return()
	return _c
}

func (_c *TaskEventUpdater_SetStatusToCompleted_Call) RunAndReturn(run func(context.Context, []job.TaskResult)) *TaskEventUpdater_SetStatusToCompleted_Call {
	_c.Run(run)
	return _c
}

// SetStatusToDataProcessing provides a mock function with given fields: ctx
func (_m *TaskEventUpdater) SetStatusToDataProcessing(ctx context.Context) {
	_m.Called(ctx)
}

// TaskEventUpdater_SetStatusToDataProcessing_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetStatusToDataProcessing'
type TaskEventUpdater_SetStatusToDataProcessing_Call struct {
	*mock.Call
}

// SetStatusToDataProcessing is a helper method to define mock.On call
//   - ctx context.Context
func (_e *TaskEventUpdater_Expecter) SetStatusToDataProcessing(ctx interface{}) *TaskEventUpdater_SetStatusToDataProcessing_Call {
	return &TaskEventUpdater_SetStatusToDataProcessing_Call{Call: _e.mock.On("SetStatusToDataProcessing", ctx)}
}

func (_c *TaskEventUpdater_SetStatusToDataProcessing_Call) Run(run func(ctx context.Context)) *TaskEventUpdater_SetStatusToDataProcessing_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *TaskEventUpdater_SetStatusToDataProcessing_Call) Return() *TaskEventUpdater_SetStatusToDataProcessing_Call {
	_c.Call.Return()
	return _c
}

func (_c *TaskEventUpdater_SetStatusToDataProcessing_Call) RunAndReturn(run func(context.Context)) *TaskEventUpdater_SetStatusToDataProcessing_Call {
	_c.Run(run)
	return _c
}

// SetStatusToDataRetrieve provides a mock function with given fields: ctx
func (_m *TaskEventUpdater) SetStatusToDataRetrieve(ctx context.Context) {
	_m.Called(ctx)
}

// TaskEventUpdater_SetStatusToDataRetrieve_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetStatusToDataRetrieve'
type TaskEventUpdater_SetStatusToDataRetrieve_Call struct {
	*mock.Call
}

// SetStatusToDataRetrieve is a helper method to define mock.On call
//   - ctx context.Context
func (_e *TaskEventUpdater_Expecter) SetStatusToDataRetrieve(ctx interface{}) *TaskEventUpdater_SetStatusToDataRetrieve_Call {
	return &TaskEventUpdater_SetStatusToDataRetrieve_Call{Call: _e.mock.On("SetStatusToDataRetrieve", ctx)}
}

func (_c *TaskEventUpdater_SetStatusToDataRetrieve_Call) Run(run func(ctx context.Context)) *TaskEventUpdater_SetStatusToDataRetrieve_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *TaskEventUpdater_SetStatusToDataRetrieve_Call) Return() *TaskEventUpdater_SetStatusToDataRetrieve_Call {
	_c.Call.Return()
	return _c
}

func (_c *TaskEventUpdater_SetStatusToDataRetrieve_Call) RunAndReturn(run func(context.Context)) *TaskEventUpdater_SetStatusToDataRetrieve_Call {
	_c.Run(run)
	return _c
}

// SetStatusToDataUpload provides a mock function with given fields: ctx
func (_m *TaskEventUpdater) SetStatusToDataUpload(ctx context.Context) {
	_m.Called(ctx)
}

// TaskEventUpdater_SetStatusToDataUpload_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetStatusToDataUpload'
type TaskEventUpdater_SetStatusToDataUpload_Call struct {
	*mock.Call
}

// SetStatusToDataUpload is a helper method to define mock.On call
//   - ctx context.Context
func (_e *TaskEventUpdater_Expecter) SetStatusToDataUpload(ctx interface{}) *TaskEventUpdater_SetStatusToDataUpload_Call {
	return &TaskEventUpdater_SetStatusToDataUpload_Call{Call: _e.mock.On("SetStatusToDataUpload", ctx)}
}

func (_c *TaskEventUpdater_SetStatusToDataUpload_Call) Run(run func(ctx context.Context)) *TaskEventUpdater_SetStatusToDataUpload_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *TaskEventUpdater_SetStatusToDataUpload_Call) Return() *TaskEventUpdater_SetStatusToDataUpload_Call {
	_c.Call.Return()
	return _c
}

func (_c *TaskEventUpdater_SetStatusToDataUpload_Call) RunAndReturn(run func(context.Context)) *TaskEventUpdater_SetStatusToDataUpload_Call {
	_c.Run(run)
	return _c
}

// SetStatusToFailed provides a mock function with given fields: ctx, err
func (_m *TaskEventUpdater) SetStatusToFailed(ctx context.Context, err error) {
	_m.Called(ctx, err)
}

// TaskEventUpdater_SetStatusToFailed_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetStatusToFailed'
type TaskEventUpdater_SetStatusToFailed_Call struct {
	*mock.Call
}

// SetStatusToFailed is a helper method to define mock.On call
//   - ctx context.Context
//   - err error
func (_e *TaskEventUpdater_Expecter) SetStatusToFailed(ctx interface{}, err interface{}) *TaskEventUpdater_SetStatusToFailed_Call {
	return &TaskEventUpdater_SetStatusToFailed_Call{Call: _e.mock.On("SetStatusToFailed", ctx, err)}
}

func (_c *TaskEventUpdater_SetStatusToFailed_Call) Run(run func(ctx context.Context, err error)) *TaskEventUpdater_SetStatusToFailed_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(error))
	})
	return _c
}

func (_c *TaskEventUpdater_SetStatusToFailed_Call) Return() *TaskEventUpdater_SetStatusToFailed_Call {
	_c.Call.Return()
	return _c
}

func (_c *TaskEventUpdater_SetStatusToFailed_Call) RunAndReturn(run func(context.Context, error)) *TaskEventUpdater_SetStatusToFailed_Call {
	_c.Run(run)
	return _c
}

// SetStatusToQueued provides a mock function with given fields: ctx
func (_m *TaskEventUpdater) SetStatusToQueued(ctx context.Context) {
	_m.Called(ctx)
}

// TaskEventUpdater_SetStatusToQueued_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetStatusToQueued'
type TaskEventUpdater_SetStatusToQueued_Call struct {
	*mock.Call
}

// SetStatusToQueued is a helper method to define mock.On call
//   - ctx context.Context
func (_e *TaskEventUpdater_Expecter) SetStatusToQueued(ctx interface{}) *TaskEventUpdater_SetStatusToQueued_Call {
	return &TaskEventUpdater_SetStatusToQueued_Call{Call: _e.mock.On("SetStatusToQueued", ctx)}
}

func (_c *TaskEventUpdater_SetStatusToQueued_Call) Run(run func(ctx context.Context)) *TaskEventUpdater_SetStatusToQueued_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *TaskEventUpdater_SetStatusToQueued_Call) Return() *TaskEventUpdater_SetStatusToQueued_Call {
	_c.Call.Return()
	return _c
}

func (_c *TaskEventUpdater_SetStatusToQueued_Call) RunAndReturn(run func(context.Context)) *TaskEventUpdater_SetStatusToQueued_Call {
	_c.Run(run)
	return _c
}

// SetStatusToSkipped provides a mock function with given fields: ctx
func (_m *TaskEventUpdater) SetStatusToSkipped(ctx context.Context) {
	_m.Called(ctx)
}

// TaskEventUpdater_SetStatusToSkipped_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetStatusToSkipped'
type TaskEventUpdater_SetStatusToSkipped_Call struct {
	*mock.Call
}

// SetStatusToSkipped is a helper method to define mock.On call
//   - ctx context.Context
func (_e *TaskEventUpdater_Expecter) SetStatusToSkipped(ctx interface{}) *TaskEventUpdater_SetStatusToSkipped_Call {
	return &TaskEventUpdater_SetStatusToSkipped_Call{Call: _e.mock.On("SetStatusToSkipped", ctx)}
}

func (_c *TaskEventUpdater_SetStatusToSkipped_Call) Run(run func(ctx context.Context)) *TaskEventUpdater_SetStatusToSkipped_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *TaskEventUpdater_SetStatusToSkipped_Call) Return() *TaskEventUpdater_SetStatusToSkipped_Call {
	_c.Call.Return()
	return _c
}

func (_c *TaskEventUpdater_SetStatusToSkipped_Call) RunAndReturn(run func(context.Context)) *TaskEventUpdater_SetStatusToSkipped_Call {
	_c.Run(run)
	return _c
}

// SetStatusToStarted provides a mock function with given fields: ctx
func (_m *TaskEventUpdater) SetStatusToStarted(ctx context.Context) {
	_m.Called(ctx)
}

// TaskEventUpdater_SetStatusToStarted_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetStatusToStarted'
type TaskEventUpdater_SetStatusToStarted_Call struct {
	*mock.Call
}

// SetStatusToStarted is a helper method to define mock.On call
//   - ctx context.Context
func (_e *TaskEventUpdater_Expecter) SetStatusToStarted(ctx interface{}) *TaskEventUpdater_SetStatusToStarted_Call {
	return &TaskEventUpdater_SetStatusToStarted_Call{Call: _e.mock.On("SetStatusToStarted", ctx)}
}

func (_c *TaskEventUpdater_SetStatusToStarted_Call) Run(run func(ctx context.Context)) *TaskEventUpdater_SetStatusToStarted_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *TaskEventUpdater_SetStatusToStarted_Call) Return() *TaskEventUpdater_SetStatusToStarted_Call {
	_c.Call.Return()
	return _c
}

func (_c *TaskEventUpdater_SetStatusToStarted_Call) RunAndReturn(run func(context.Context)) *TaskEventUpdater_SetStatusToStarted_Call {
	_c.Run(run)
	return _c
}

// NewTaskEventUpdater creates a new instance of TaskEventUpdater. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTaskEventUpdater(t interface {
	mock.TestingT
	Cleanup(func())
}) *TaskEventUpdater {
	mock := &TaskEventUpdater{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
