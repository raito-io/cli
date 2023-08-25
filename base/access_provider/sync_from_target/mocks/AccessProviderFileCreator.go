// Code generated by mockery v2.33.0. DO NOT EDIT.

package mocks

import (
	sync_from_target "github.com/raito-io/cli/base/access_provider/sync_from_target"
	mock "github.com/stretchr/testify/mock"
)

// AccessProviderFileCreator is an autogenerated mock type for the AccessProviderFileCreator type
type AccessProviderFileCreator struct {
	mock.Mock
}

type AccessProviderFileCreator_Expecter struct {
	mock *mock.Mock
}

func (_m *AccessProviderFileCreator) EXPECT() *AccessProviderFileCreator_Expecter {
	return &AccessProviderFileCreator_Expecter{mock: &_m.Mock}
}

// AddAccessProviders provides a mock function with given fields: dataAccessList
func (_m *AccessProviderFileCreator) AddAccessProviders(dataAccessList ...*sync_from_target.AccessProvider) error {
	_va := make([]interface{}, len(dataAccessList))
	for _i := range dataAccessList {
		_va[_i] = dataAccessList[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(...*sync_from_target.AccessProvider) error); ok {
		r0 = rf(dataAccessList...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AccessProviderFileCreator_AddAccessProviders_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddAccessProviders'
type AccessProviderFileCreator_AddAccessProviders_Call struct {
	*mock.Call
}

// AddAccessProviders is a helper method to define mock.On call
//   - dataAccessList ...*sync_from_target.AccessProvider
func (_e *AccessProviderFileCreator_Expecter) AddAccessProviders(dataAccessList ...interface{}) *AccessProviderFileCreator_AddAccessProviders_Call {
	return &AccessProviderFileCreator_AddAccessProviders_Call{Call: _e.mock.On("AddAccessProviders",
		append([]interface{}{}, dataAccessList...)...)}
}

func (_c *AccessProviderFileCreator_AddAccessProviders_Call) Run(run func(dataAccessList ...*sync_from_target.AccessProvider)) *AccessProviderFileCreator_AddAccessProviders_Call {
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

func (_c *AccessProviderFileCreator_AddAccessProviders_Call) Return(_a0 error) *AccessProviderFileCreator_AddAccessProviders_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AccessProviderFileCreator_AddAccessProviders_Call) RunAndReturn(run func(...*sync_from_target.AccessProvider) error) *AccessProviderFileCreator_AddAccessProviders_Call {
	_c.Call.Return(run)
	return _c
}

// Close provides a mock function with given fields:
func (_m *AccessProviderFileCreator) Close() {
	_m.Called()
}

// AccessProviderFileCreator_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type AccessProviderFileCreator_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *AccessProviderFileCreator_Expecter) Close() *AccessProviderFileCreator_Close_Call {
	return &AccessProviderFileCreator_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *AccessProviderFileCreator_Close_Call) Run(run func()) *AccessProviderFileCreator_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *AccessProviderFileCreator_Close_Call) Return() *AccessProviderFileCreator_Close_Call {
	_c.Call.Return()
	return _c
}

func (_c *AccessProviderFileCreator_Close_Call) RunAndReturn(run func()) *AccessProviderFileCreator_Close_Call {
	_c.Call.Return(run)
	return _c
}

// GetAccessProviderCount provides a mock function with given fields:
func (_m *AccessProviderFileCreator) GetAccessProviderCount() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// AccessProviderFileCreator_GetAccessProviderCount_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAccessProviderCount'
type AccessProviderFileCreator_GetAccessProviderCount_Call struct {
	*mock.Call
}

// GetAccessProviderCount is a helper method to define mock.On call
func (_e *AccessProviderFileCreator_Expecter) GetAccessProviderCount() *AccessProviderFileCreator_GetAccessProviderCount_Call {
	return &AccessProviderFileCreator_GetAccessProviderCount_Call{Call: _e.mock.On("GetAccessProviderCount")}
}

func (_c *AccessProviderFileCreator_GetAccessProviderCount_Call) Run(run func()) *AccessProviderFileCreator_GetAccessProviderCount_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *AccessProviderFileCreator_GetAccessProviderCount_Call) Return(_a0 int) *AccessProviderFileCreator_GetAccessProviderCount_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AccessProviderFileCreator_GetAccessProviderCount_Call) RunAndReturn(run func() int) *AccessProviderFileCreator_GetAccessProviderCount_Call {
	_c.Call.Return(run)
	return _c
}

// NewAccessProviderFileCreator creates a new instance of AccessProviderFileCreator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAccessProviderFileCreator(t interface {
	mock.TestingT
	Cleanup(func())
}) *AccessProviderFileCreator {
	mock := &AccessProviderFileCreator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
