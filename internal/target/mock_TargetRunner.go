// Code generated by mockery v2.50.1. DO NOT EDIT.

package target

import (
	context "context"

	types "github.com/raito-io/cli/internal/target/types"
	mock "github.com/stretchr/testify/mock"
)

// MockTargetRunner is an autogenerated mock type for the TargetRunner type
type MockTargetRunner struct {
	mock.Mock
}

type MockTargetRunner_Expecter struct {
	mock *mock.Mock
}

func (_m *MockTargetRunner) EXPECT() *MockTargetRunner_Expecter {
	return &MockTargetRunner_Expecter{mock: &_m.Mock}
}

// Finalize provides a mock function with given fields: ctx, baseConfig, options
func (_m *MockTargetRunner) Finalize(ctx context.Context, baseConfig *types.BaseConfig, options *Options) error {
	ret := _m.Called(ctx, baseConfig, options)

	if len(ret) == 0 {
		panic("no return value specified for Finalize")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.BaseConfig, *Options) error); ok {
		r0 = rf(ctx, baseConfig, options)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockTargetRunner_Finalize_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Finalize'
type MockTargetRunner_Finalize_Call struct {
	*mock.Call
}

// Finalize is a helper method to define mock.On call
//   - ctx context.Context
//   - baseConfig *types.BaseConfig
//   - options *Options
func (_e *MockTargetRunner_Expecter) Finalize(ctx interface{}, baseConfig interface{}, options interface{}) *MockTargetRunner_Finalize_Call {
	return &MockTargetRunner_Finalize_Call{Call: _e.mock.On("Finalize", ctx, baseConfig, options)}
}

func (_c *MockTargetRunner_Finalize_Call) Run(run func(ctx context.Context, baseConfig *types.BaseConfig, options *Options)) *MockTargetRunner_Finalize_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*types.BaseConfig), args[2].(*Options))
	})
	return _c
}

func (_c *MockTargetRunner_Finalize_Call) Return(_a0 error) *MockTargetRunner_Finalize_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTargetRunner_Finalize_Call) RunAndReturn(run func(context.Context, *types.BaseConfig, *Options) error) *MockTargetRunner_Finalize_Call {
	_c.Call.Return(run)
	return _c
}

// RunType provides a mock function with no fields
func (_m *MockTargetRunner) RunType() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for RunType")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockTargetRunner_RunType_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RunType'
type MockTargetRunner_RunType_Call struct {
	*mock.Call
}

// RunType is a helper method to define mock.On call
func (_e *MockTargetRunner_Expecter) RunType() *MockTargetRunner_RunType_Call {
	return &MockTargetRunner_RunType_Call{Call: _e.mock.On("RunType")}
}

func (_c *MockTargetRunner_RunType_Call) Run(run func()) *MockTargetRunner_RunType_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockTargetRunner_RunType_Call) Return(_a0 string) *MockTargetRunner_RunType_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTargetRunner_RunType_Call) RunAndReturn(run func() string) *MockTargetRunner_RunType_Call {
	_c.Call.Return(run)
	return _c
}

// TargetSync provides a mock function with given fields: ctx, targetConfig
func (_m *MockTargetRunner) TargetSync(ctx context.Context, targetConfig *types.BaseTargetConfig) error {
	ret := _m.Called(ctx, targetConfig)

	if len(ret) == 0 {
		panic("no return value specified for TargetSync")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.BaseTargetConfig) error); ok {
		r0 = rf(ctx, targetConfig)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockTargetRunner_TargetSync_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TargetSync'
type MockTargetRunner_TargetSync_Call struct {
	*mock.Call
}

// TargetSync is a helper method to define mock.On call
//   - ctx context.Context
//   - targetConfig *types.BaseTargetConfig
func (_e *MockTargetRunner_Expecter) TargetSync(ctx interface{}, targetConfig interface{}) *MockTargetRunner_TargetSync_Call {
	return &MockTargetRunner_TargetSync_Call{Call: _e.mock.On("TargetSync", ctx, targetConfig)}
}

func (_c *MockTargetRunner_TargetSync_Call) Run(run func(ctx context.Context, targetConfig *types.BaseTargetConfig)) *MockTargetRunner_TargetSync_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*types.BaseTargetConfig))
	})
	return _c
}

func (_c *MockTargetRunner_TargetSync_Call) Return(syncError error) *MockTargetRunner_TargetSync_Call {
	_c.Call.Return(syncError)
	return _c
}

func (_c *MockTargetRunner_TargetSync_Call) RunAndReturn(run func(context.Context, *types.BaseTargetConfig) error) *MockTargetRunner_TargetSync_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockTargetRunner creates a new instance of MockTargetRunner. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockTargetRunner(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockTargetRunner {
	mock := &MockTargetRunner{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
