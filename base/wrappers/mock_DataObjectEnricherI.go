// Code generated by mockery v2.33.0. DO NOT EDIT.

package wrappers

import (
	context "context"

	data_source "github.com/raito-io/cli/base/data_source"
	mock "github.com/stretchr/testify/mock"
)

// MockDataObjectEnricherI is an autogenerated mock type for the DataObjectEnricherI type
type MockDataObjectEnricherI struct {
	mock.Mock
}

type MockDataObjectEnricherI_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDataObjectEnricherI) EXPECT() *MockDataObjectEnricherI_Expecter {
	return &MockDataObjectEnricherI_Expecter{mock: &_m.Mock}
}

// Close provides a mock function with given fields: ctx
func (_m *MockDataObjectEnricherI) Close(ctx context.Context) (int, error) {
	ret := _m.Called(ctx)

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (int, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) int); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDataObjectEnricherI_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type MockDataObjectEnricherI_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockDataObjectEnricherI_Expecter) Close(ctx interface{}) *MockDataObjectEnricherI_Close_Call {
	return &MockDataObjectEnricherI_Close_Call{Call: _e.mock.On("Close", ctx)}
}

func (_c *MockDataObjectEnricherI_Close_Call) Run(run func(ctx context.Context)) *MockDataObjectEnricherI_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockDataObjectEnricherI_Close_Call) Return(_a0 int, _a1 error) *MockDataObjectEnricherI_Close_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDataObjectEnricherI_Close_Call) RunAndReturn(run func(context.Context) (int, error)) *MockDataObjectEnricherI_Close_Call {
	_c.Call.Return(run)
	return _c
}

// Enrich provides a mock function with given fields: ctx, dataObject
func (_m *MockDataObjectEnricherI) Enrich(ctx context.Context, dataObject *data_source.DataObject) error {
	ret := _m.Called(ctx, dataObject)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *data_source.DataObject) error); ok {
		r0 = rf(ctx, dataObject)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDataObjectEnricherI_Enrich_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Enrich'
type MockDataObjectEnricherI_Enrich_Call struct {
	*mock.Call
}

// Enrich is a helper method to define mock.On call
//   - ctx context.Context
//   - dataObject *data_source.DataObject
func (_e *MockDataObjectEnricherI_Expecter) Enrich(ctx interface{}, dataObject interface{}) *MockDataObjectEnricherI_Enrich_Call {
	return &MockDataObjectEnricherI_Enrich_Call{Call: _e.mock.On("Enrich", ctx, dataObject)}
}

func (_c *MockDataObjectEnricherI_Enrich_Call) Run(run func(ctx context.Context, dataObject *data_source.DataObject)) *MockDataObjectEnricherI_Enrich_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*data_source.DataObject))
	})
	return _c
}

func (_c *MockDataObjectEnricherI_Enrich_Call) Return(_a0 error) *MockDataObjectEnricherI_Enrich_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDataObjectEnricherI_Enrich_Call) RunAndReturn(run func(context.Context, *data_source.DataObject) error) *MockDataObjectEnricherI_Enrich_Call {
	_c.Call.Return(run)
	return _c
}

// Initialize provides a mock function with given fields: ctx, dataObjectWriter, config
func (_m *MockDataObjectEnricherI) Initialize(ctx context.Context, dataObjectWriter DataObjectWriter, config map[string]string) error {
	ret := _m.Called(ctx, dataObjectWriter, config)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, DataObjectWriter, map[string]string) error); ok {
		r0 = rf(ctx, dataObjectWriter, config)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDataObjectEnricherI_Initialize_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Initialize'
type MockDataObjectEnricherI_Initialize_Call struct {
	*mock.Call
}

// Initialize is a helper method to define mock.On call
//   - ctx context.Context
//   - dataObjectWriter DataObjectWriter
//   - config map[string]string
func (_e *MockDataObjectEnricherI_Expecter) Initialize(ctx interface{}, dataObjectWriter interface{}, config interface{}) *MockDataObjectEnricherI_Initialize_Call {
	return &MockDataObjectEnricherI_Initialize_Call{Call: _e.mock.On("Initialize", ctx, dataObjectWriter, config)}
}

func (_c *MockDataObjectEnricherI_Initialize_Call) Run(run func(ctx context.Context, dataObjectWriter DataObjectWriter, config map[string]string)) *MockDataObjectEnricherI_Initialize_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(DataObjectWriter), args[2].(map[string]string))
	})
	return _c
}

func (_c *MockDataObjectEnricherI_Initialize_Call) Return(_a0 error) *MockDataObjectEnricherI_Initialize_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDataObjectEnricherI_Initialize_Call) RunAndReturn(run func(context.Context, DataObjectWriter, map[string]string) error) *MockDataObjectEnricherI_Initialize_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockDataObjectEnricherI creates a new instance of MockDataObjectEnricherI. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDataObjectEnricherI(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDataObjectEnricherI {
	mock := &MockDataObjectEnricherI{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
