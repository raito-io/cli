// Code generated by mockery v2.52.3. DO NOT EDIT.

package mocks

import (
	sync_from_target "github.com/raito-io/cli/base/access_provider/sync_from_target"
	mock "github.com/stretchr/testify/mock"
)

// AccessProviderHandler is an autogenerated mock type for the AccessProviderHandler type
type AccessProviderHandler struct {
	mock.Mock
}

type AccessProviderHandler_Expecter struct {
	mock *mock.Mock
}

func (_m *AccessProviderHandler) EXPECT() *AccessProviderHandler_Expecter {
	return &AccessProviderHandler_Expecter{mock: &_m.Mock}
}

// AddAccessProviders provides a mock function with given fields: dataAccessList
func (_m *AccessProviderHandler) AddAccessProviders(dataAccessList ...*sync_from_target.AccessProvider) error {
	_va := make([]interface{}, len(dataAccessList))
	for _i := range dataAccessList {
		_va[_i] = dataAccessList[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for AddAccessProviders")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(...*sync_from_target.AccessProvider) error); ok {
		r0 = rf(dataAccessList...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AccessProviderHandler_AddAccessProviders_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddAccessProviders'
type AccessProviderHandler_AddAccessProviders_Call struct {
	*mock.Call
}

// AddAccessProviders is a helper method to define mock.On call
//   - dataAccessList ...*sync_from_target.AccessProvider
func (_e *AccessProviderHandler_Expecter) AddAccessProviders(dataAccessList ...interface{}) *AccessProviderHandler_AddAccessProviders_Call {
	return &AccessProviderHandler_AddAccessProviders_Call{Call: _e.mock.On("AddAccessProviders",
		append([]interface{}{}, dataAccessList...)...)}
}

func (_c *AccessProviderHandler_AddAccessProviders_Call) Run(run func(dataAccessList ...*sync_from_target.AccessProvider)) *AccessProviderHandler_AddAccessProviders_Call {
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

func (_c *AccessProviderHandler_AddAccessProviders_Call) Return(_a0 error) *AccessProviderHandler_AddAccessProviders_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AccessProviderHandler_AddAccessProviders_Call) RunAndReturn(run func(...*sync_from_target.AccessProvider) error) *AccessProviderHandler_AddAccessProviders_Call {
	_c.Call.Return(run)
	return _c
}

// NewAccessProviderHandler creates a new instance of AccessProviderHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAccessProviderHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *AccessProviderHandler {
	mock := &AccessProviderHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
