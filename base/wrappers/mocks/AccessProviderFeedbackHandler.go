// Code generated by mockery v2.21.1. DO NOT EDIT.

package mocks

import (
	sync_to_target "github.com/raito-io/cli/base/access_provider/sync_to_target"
	mock "github.com/stretchr/testify/mock"
)

// AccessProviderFeedbackHandler is an autogenerated mock type for the AccessProviderFeedbackHandler type
type AccessProviderFeedbackHandler struct {
	mock.Mock
}

type AccessProviderFeedbackHandler_Expecter struct {
	mock *mock.Mock
}

func (_m *AccessProviderFeedbackHandler) EXPECT() *AccessProviderFeedbackHandler_Expecter {
	return &AccessProviderFeedbackHandler_Expecter{mock: &_m.Mock}
}

// AddAccessProviderFeedback provides a mock function with given fields: accessProviderId, accessFeedback
func (_m *AccessProviderFeedbackHandler) AddAccessProviderFeedback(accessProviderId string, accessFeedback ...sync_to_target.AccessSyncFeedbackInformation) error {
	_va := make([]interface{}, len(accessFeedback))
	for _i := range accessFeedback {
		_va[_i] = accessFeedback[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, accessProviderId)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, ...sync_to_target.AccessSyncFeedbackInformation) error); ok {
		r0 = rf(accessProviderId, accessFeedback...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AccessProviderFeedbackHandler_AddAccessProviderFeedback_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddAccessProviderFeedback'
type AccessProviderFeedbackHandler_AddAccessProviderFeedback_Call struct {
	*mock.Call
}

// AddAccessProviderFeedback is a helper method to define mock.On call
//   - accessProviderId string
//   - accessFeedback ...sync_to_target.AccessSyncFeedbackInformation
func (_e *AccessProviderFeedbackHandler_Expecter) AddAccessProviderFeedback(accessProviderId interface{}, accessFeedback ...interface{}) *AccessProviderFeedbackHandler_AddAccessProviderFeedback_Call {
	return &AccessProviderFeedbackHandler_AddAccessProviderFeedback_Call{Call: _e.mock.On("AddAccessProviderFeedback",
		append([]interface{}{accessProviderId}, accessFeedback...)...)}
}

func (_c *AccessProviderFeedbackHandler_AddAccessProviderFeedback_Call) Run(run func(accessProviderId string, accessFeedback ...sync_to_target.AccessSyncFeedbackInformation)) *AccessProviderFeedbackHandler_AddAccessProviderFeedback_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]sync_to_target.AccessSyncFeedbackInformation, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(sync_to_target.AccessSyncFeedbackInformation)
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *AccessProviderFeedbackHandler_AddAccessProviderFeedback_Call) Return(_a0 error) *AccessProviderFeedbackHandler_AddAccessProviderFeedback_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AccessProviderFeedbackHandler_AddAccessProviderFeedback_Call) RunAndReturn(run func(string, ...sync_to_target.AccessSyncFeedbackInformation) error) *AccessProviderFeedbackHandler_AddAccessProviderFeedback_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewAccessProviderFeedbackHandler interface {
	mock.TestingT
	Cleanup(func())
}

// NewAccessProviderFeedbackHandler creates a new instance of AccessProviderFeedbackHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAccessProviderFeedbackHandler(t mockConstructorTestingTNewAccessProviderFeedbackHandler) *AccessProviderFeedbackHandler {
	mock := &AccessProviderFeedbackHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
