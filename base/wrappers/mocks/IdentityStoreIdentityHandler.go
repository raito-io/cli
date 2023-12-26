// Code generated by mockery v2.39.1. DO NOT EDIT.

package mocks

import (
	identity_store "github.com/raito-io/cli/base/identity_store"
	mock "github.com/stretchr/testify/mock"
)

// IdentityStoreIdentityHandler is an autogenerated mock type for the IdentityStoreIdentityHandler type
type IdentityStoreIdentityHandler struct {
	mock.Mock
}

type IdentityStoreIdentityHandler_Expecter struct {
	mock *mock.Mock
}

func (_m *IdentityStoreIdentityHandler) EXPECT() *IdentityStoreIdentityHandler_Expecter {
	return &IdentityStoreIdentityHandler_Expecter{mock: &_m.Mock}
}

// AddGroups provides a mock function with given fields: groups
func (_m *IdentityStoreIdentityHandler) AddGroups(groups ...*identity_store.Group) error {
	_va := make([]interface{}, len(groups))
	for _i := range groups {
		_va[_i] = groups[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for AddGroups")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(...*identity_store.Group) error); ok {
		r0 = rf(groups...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IdentityStoreIdentityHandler_AddGroups_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddGroups'
type IdentityStoreIdentityHandler_AddGroups_Call struct {
	*mock.Call
}

// AddGroups is a helper method to define mock.On call
//   - groups ...*identity_store.Group
func (_e *IdentityStoreIdentityHandler_Expecter) AddGroups(groups ...interface{}) *IdentityStoreIdentityHandler_AddGroups_Call {
	return &IdentityStoreIdentityHandler_AddGroups_Call{Call: _e.mock.On("AddGroups",
		append([]interface{}{}, groups...)...)}
}

func (_c *IdentityStoreIdentityHandler_AddGroups_Call) Run(run func(groups ...*identity_store.Group)) *IdentityStoreIdentityHandler_AddGroups_Call {
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

func (_c *IdentityStoreIdentityHandler_AddGroups_Call) Return(_a0 error) *IdentityStoreIdentityHandler_AddGroups_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *IdentityStoreIdentityHandler_AddGroups_Call) RunAndReturn(run func(...*identity_store.Group) error) *IdentityStoreIdentityHandler_AddGroups_Call {
	_c.Call.Return(run)
	return _c
}

// AddUsers provides a mock function with given fields: user
func (_m *IdentityStoreIdentityHandler) AddUsers(user ...*identity_store.User) error {
	_va := make([]interface{}, len(user))
	for _i := range user {
		_va[_i] = user[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for AddUsers")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(...*identity_store.User) error); ok {
		r0 = rf(user...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IdentityStoreIdentityHandler_AddUsers_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddUsers'
type IdentityStoreIdentityHandler_AddUsers_Call struct {
	*mock.Call
}

// AddUsers is a helper method to define mock.On call
//   - user ...*identity_store.User
func (_e *IdentityStoreIdentityHandler_Expecter) AddUsers(user ...interface{}) *IdentityStoreIdentityHandler_AddUsers_Call {
	return &IdentityStoreIdentityHandler_AddUsers_Call{Call: _e.mock.On("AddUsers",
		append([]interface{}{}, user...)...)}
}

func (_c *IdentityStoreIdentityHandler_AddUsers_Call) Run(run func(user ...*identity_store.User)) *IdentityStoreIdentityHandler_AddUsers_Call {
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

func (_c *IdentityStoreIdentityHandler_AddUsers_Call) Return(_a0 error) *IdentityStoreIdentityHandler_AddUsers_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *IdentityStoreIdentityHandler_AddUsers_Call) RunAndReturn(run func(...*identity_store.User) error) *IdentityStoreIdentityHandler_AddUsers_Call {
	_c.Call.Return(run)
	return _c
}

// NewIdentityStoreIdentityHandler creates a new instance of IdentityStoreIdentityHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIdentityStoreIdentityHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *IdentityStoreIdentityHandler {
	mock := &IdentityStoreIdentityHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
