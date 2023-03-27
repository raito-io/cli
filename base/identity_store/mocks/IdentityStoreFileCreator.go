// Code generated by mockery v2.23.1. DO NOT EDIT.

package mocks

import (
	identity_store "github.com/raito-io/cli/base/identity_store"
	mock "github.com/stretchr/testify/mock"
)

// IdentityStoreFileCreator is an autogenerated mock type for the IdentityStoreFileCreator type
type IdentityStoreFileCreator struct {
	mock.Mock
}

type IdentityStoreFileCreator_Expecter struct {
	mock *mock.Mock
}

func (_m *IdentityStoreFileCreator) EXPECT() *IdentityStoreFileCreator_Expecter {
	return &IdentityStoreFileCreator_Expecter{mock: &_m.Mock}
}

// AddGroups provides a mock function with given fields: groups
func (_m *IdentityStoreFileCreator) AddGroups(groups ...*identity_store.Group) error {
	_va := make([]interface{}, len(groups))
	for _i := range groups {
		_va[_i] = groups[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(...*identity_store.Group) error); ok {
		r0 = rf(groups...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IdentityStoreFileCreator_AddGroups_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddGroups'
type IdentityStoreFileCreator_AddGroups_Call struct {
	*mock.Call
}

// AddGroups is a helper method to define mock.On call
//   - groups ...*identity_store.Group
func (_e *IdentityStoreFileCreator_Expecter) AddGroups(groups ...interface{}) *IdentityStoreFileCreator_AddGroups_Call {
	return &IdentityStoreFileCreator_AddGroups_Call{Call: _e.mock.On("AddGroups",
		append([]interface{}{}, groups...)...)}
}

func (_c *IdentityStoreFileCreator_AddGroups_Call) Run(run func(groups ...*identity_store.Group)) *IdentityStoreFileCreator_AddGroups_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]*identity_store.Group, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(*identity_store.Group)
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *IdentityStoreFileCreator_AddGroups_Call) Return(_a0 error) *IdentityStoreFileCreator_AddGroups_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *IdentityStoreFileCreator_AddGroups_Call) RunAndReturn(run func(...*identity_store.Group) error) *IdentityStoreFileCreator_AddGroups_Call {
	_c.Call.Return(run)
	return _c
}

// AddUsers provides a mock function with given fields: users
func (_m *IdentityStoreFileCreator) AddUsers(users ...*identity_store.User) error {
	_va := make([]interface{}, len(users))
	for _i := range users {
		_va[_i] = users[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(...*identity_store.User) error); ok {
		r0 = rf(users...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IdentityStoreFileCreator_AddUsers_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddUsers'
type IdentityStoreFileCreator_AddUsers_Call struct {
	*mock.Call
}

// AddUsers is a helper method to define mock.On call
//   - users ...*identity_store.User
func (_e *IdentityStoreFileCreator_Expecter) AddUsers(users ...interface{}) *IdentityStoreFileCreator_AddUsers_Call {
	return &IdentityStoreFileCreator_AddUsers_Call{Call: _e.mock.On("AddUsers",
		append([]interface{}{}, users...)...)}
}

func (_c *IdentityStoreFileCreator_AddUsers_Call) Run(run func(users ...*identity_store.User)) *IdentityStoreFileCreator_AddUsers_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]*identity_store.User, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(*identity_store.User)
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *IdentityStoreFileCreator_AddUsers_Call) Return(_a0 error) *IdentityStoreFileCreator_AddUsers_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *IdentityStoreFileCreator_AddUsers_Call) RunAndReturn(run func(...*identity_store.User) error) *IdentityStoreFileCreator_AddUsers_Call {
	_c.Call.Return(run)
	return _c
}

// Close provides a mock function with given fields:
func (_m *IdentityStoreFileCreator) Close() {
	_m.Called()
}

// IdentityStoreFileCreator_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type IdentityStoreFileCreator_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *IdentityStoreFileCreator_Expecter) Close() *IdentityStoreFileCreator_Close_Call {
	return &IdentityStoreFileCreator_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *IdentityStoreFileCreator_Close_Call) Run(run func()) *IdentityStoreFileCreator_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *IdentityStoreFileCreator_Close_Call) Return() *IdentityStoreFileCreator_Close_Call {
	_c.Call.Return()
	return _c
}

func (_c *IdentityStoreFileCreator_Close_Call) RunAndReturn(run func()) *IdentityStoreFileCreator_Close_Call {
	_c.Call.Return(run)
	return _c
}

// GetGroupCount provides a mock function with given fields:
func (_m *IdentityStoreFileCreator) GetGroupCount() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// IdentityStoreFileCreator_GetGroupCount_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetGroupCount'
type IdentityStoreFileCreator_GetGroupCount_Call struct {
	*mock.Call
}

// GetGroupCount is a helper method to define mock.On call
func (_e *IdentityStoreFileCreator_Expecter) GetGroupCount() *IdentityStoreFileCreator_GetGroupCount_Call {
	return &IdentityStoreFileCreator_GetGroupCount_Call{Call: _e.mock.On("GetGroupCount")}
}

func (_c *IdentityStoreFileCreator_GetGroupCount_Call) Run(run func()) *IdentityStoreFileCreator_GetGroupCount_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *IdentityStoreFileCreator_GetGroupCount_Call) Return(_a0 int) *IdentityStoreFileCreator_GetGroupCount_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *IdentityStoreFileCreator_GetGroupCount_Call) RunAndReturn(run func() int) *IdentityStoreFileCreator_GetGroupCount_Call {
	_c.Call.Return(run)
	return _c
}

// GetUserCount provides a mock function with given fields:
func (_m *IdentityStoreFileCreator) GetUserCount() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// IdentityStoreFileCreator_GetUserCount_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserCount'
type IdentityStoreFileCreator_GetUserCount_Call struct {
	*mock.Call
}

// GetUserCount is a helper method to define mock.On call
func (_e *IdentityStoreFileCreator_Expecter) GetUserCount() *IdentityStoreFileCreator_GetUserCount_Call {
	return &IdentityStoreFileCreator_GetUserCount_Call{Call: _e.mock.On("GetUserCount")}
}

func (_c *IdentityStoreFileCreator_GetUserCount_Call) Run(run func()) *IdentityStoreFileCreator_GetUserCount_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *IdentityStoreFileCreator_GetUserCount_Call) Return(_a0 int) *IdentityStoreFileCreator_GetUserCount_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *IdentityStoreFileCreator_GetUserCount_Call) RunAndReturn(run func() int) *IdentityStoreFileCreator_GetUserCount_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewIdentityStoreFileCreator interface {
	mock.TestingT
	Cleanup(func())
}

// NewIdentityStoreFileCreator creates a new instance of IdentityStoreFileCreator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIdentityStoreFileCreator(t mockConstructorTestingTNewIdentityStoreFileCreator) *IdentityStoreFileCreator {
	mock := &IdentityStoreFileCreator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
