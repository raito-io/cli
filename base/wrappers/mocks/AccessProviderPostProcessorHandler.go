// Code generated by mockery v2.39.1. DO NOT EDIT.

package mocks

import (
	sync_from_target "github.com/raito-io/cli/base/access_provider/sync_from_target"
	mock "github.com/stretchr/testify/mock"
)

// AccessProviderPostProcessorHandler is an autogenerated mock type for the AccessProviderPostProcessorHandler type
type AccessProviderPostProcessorHandler struct {
	mock.Mock
}

type AccessProviderPostProcessorHandler_Expecter struct {
	mock *mock.Mock
}

func (_m *AccessProviderPostProcessorHandler) EXPECT() *AccessProviderPostProcessorHandler_Expecter {
	return &AccessProviderPostProcessorHandler_Expecter{mock: &_m.Mock}
}

// AddAccessProviders provides a mock function with given fields: accessProviders
func (_m *AccessProviderPostProcessorHandler) AddAccessProviders(accessProviders ...*sync_from_target.AccessProvider) error {
	_va := make([]interface{}, len(accessProviders))
	for _i := range accessProviders {
		_va[_i] = accessProviders[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for AddAccessProviders")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(...*sync_from_target.AccessProvider) error); ok {
		r0 = rf(accessProviders...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AccessProviderPostProcessorHandler_AddAccessProviders_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddAccessProviders'
type AccessProviderPostProcessorHandler_AddAccessProviders_Call struct {
	*mock.Call
}

// AddAccessProviders is a helper method to define mock.On call
//   - accessProviders ...*sync_from_target.AccessProvider
func (_e *AccessProviderPostProcessorHandler_Expecter) AddAccessProviders(accessProviders ...interface{}) *AccessProviderPostProcessorHandler_AddAccessProviders_Call {
	return &AccessProviderPostProcessorHandler_AddAccessProviders_Call{Call: _e.mock.On("AddAccessProviders",
		append([]interface{}{}, accessProviders...)...)}
}

func (_c *AccessProviderPostProcessorHandler_AddAccessProviders_Call) Run(run func(accessProviders ...*sync_from_target.AccessProvider)) *AccessProviderPostProcessorHandler_AddAccessProviders_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]*sync_from_target.AccessProvider, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(*sync_from_target.AccessProvider)
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *AccessProviderPostProcessorHandler_AddAccessProviders_Call) Return(_a0 error) *AccessProviderPostProcessorHandler_AddAccessProviders_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AccessProviderPostProcessorHandler_AddAccessProviders_Call) RunAndReturn(run func(...*sync_from_target.AccessProvider) error) *AccessProviderPostProcessorHandler_AddAccessProviders_Call {
	_c.Call.Return(run)
	return _c
}

// NewAccessProviderPostProcessorHandler creates a new instance of AccessProviderPostProcessorHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAccessProviderPostProcessorHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *AccessProviderPostProcessorHandler {
	mock := &AccessProviderPostProcessorHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}