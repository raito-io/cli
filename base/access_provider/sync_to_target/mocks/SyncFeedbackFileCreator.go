// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	sync_to_target "github.com/raito-io/cli/base/access_provider/sync_to_target"
	mock "github.com/stretchr/testify/mock"
)

// SyncFeedbackFileCreator is an autogenerated mock type for the SyncFeedbackFileCreator type
type SyncFeedbackFileCreator struct {
	mock.Mock
}

type SyncFeedbackFileCreator_Expecter struct {
	mock *mock.Mock
}

func (_m *SyncFeedbackFileCreator) EXPECT() *SyncFeedbackFileCreator_Expecter {
	return &SyncFeedbackFileCreator_Expecter{mock: &_m.Mock}
}

// AddAccessProviderFeedback provides a mock function with given fields: accessProviderFeedback
func (_m *SyncFeedbackFileCreator) AddAccessProviderFeedback(accessProviderFeedback sync_to_target.AccessProviderSyncFeedback) error {
	ret := _m.Called(accessProviderFeedback)

	if len(ret) == 0 {
		panic("no return value specified for AddAccessProviderFeedback")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(sync_to_target.AccessProviderSyncFeedback) error); ok {
		r0 = rf(accessProviderFeedback)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SyncFeedbackFileCreator_AddAccessProviderFeedback_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddAccessProviderFeedback'
type SyncFeedbackFileCreator_AddAccessProviderFeedback_Call struct {
	*mock.Call
}

// AddAccessProviderFeedback is a helper method to define mock.On call
//   - accessProviderFeedback sync_to_target.AccessProviderSyncFeedback
func (_e *SyncFeedbackFileCreator_Expecter) AddAccessProviderFeedback(accessProviderFeedback interface{}) *SyncFeedbackFileCreator_AddAccessProviderFeedback_Call {
	return &SyncFeedbackFileCreator_AddAccessProviderFeedback_Call{Call: _e.mock.On("AddAccessProviderFeedback", accessProviderFeedback)}
}

func (_c *SyncFeedbackFileCreator_AddAccessProviderFeedback_Call) Run(run func(accessProviderFeedback sync_to_target.AccessProviderSyncFeedback)) *SyncFeedbackFileCreator_AddAccessProviderFeedback_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(sync_to_target.AccessProviderSyncFeedback))
	})
	return _c
}

func (_c *SyncFeedbackFileCreator_AddAccessProviderFeedback_Call) Return(_a0 error) *SyncFeedbackFileCreator_AddAccessProviderFeedback_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *SyncFeedbackFileCreator_AddAccessProviderFeedback_Call) RunAndReturn(run func(sync_to_target.AccessProviderSyncFeedback) error) *SyncFeedbackFileCreator_AddAccessProviderFeedback_Call {
	_c.Call.Return(run)
	return _c
}

// Close provides a mock function with given fields:
func (_m *SyncFeedbackFileCreator) Close() {
	_m.Called()
}

// SyncFeedbackFileCreator_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type SyncFeedbackFileCreator_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *SyncFeedbackFileCreator_Expecter) Close() *SyncFeedbackFileCreator_Close_Call {
	return &SyncFeedbackFileCreator_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *SyncFeedbackFileCreator_Close_Call) Run(run func()) *SyncFeedbackFileCreator_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *SyncFeedbackFileCreator_Close_Call) Return() *SyncFeedbackFileCreator_Close_Call {
	_c.Call.Return()
	return _c
}

func (_c *SyncFeedbackFileCreator_Close_Call) RunAndReturn(run func()) *SyncFeedbackFileCreator_Close_Call {
	_c.Call.Return(run)
	return _c
}

// GetAccessProviderCount provides a mock function with given fields:
func (_m *SyncFeedbackFileCreator) GetAccessProviderCount() int {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetAccessProviderCount")
	}

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// SyncFeedbackFileCreator_GetAccessProviderCount_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAccessProviderCount'
type SyncFeedbackFileCreator_GetAccessProviderCount_Call struct {
	*mock.Call
}

// GetAccessProviderCount is a helper method to define mock.On call
func (_e *SyncFeedbackFileCreator_Expecter) GetAccessProviderCount() *SyncFeedbackFileCreator_GetAccessProviderCount_Call {
	return &SyncFeedbackFileCreator_GetAccessProviderCount_Call{Call: _e.mock.On("GetAccessProviderCount")}
}

func (_c *SyncFeedbackFileCreator_GetAccessProviderCount_Call) Run(run func()) *SyncFeedbackFileCreator_GetAccessProviderCount_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *SyncFeedbackFileCreator_GetAccessProviderCount_Call) Return(_a0 int) *SyncFeedbackFileCreator_GetAccessProviderCount_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *SyncFeedbackFileCreator_GetAccessProviderCount_Call) RunAndReturn(run func() int) *SyncFeedbackFileCreator_GetAccessProviderCount_Call {
	_c.Call.Return(run)
	return _c
}

// NewSyncFeedbackFileCreator creates a new instance of SyncFeedbackFileCreator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSyncFeedbackFileCreator(t interface {
	mock.TestingT
	Cleanup(func())
}) *SyncFeedbackFileCreator {
	mock := &SyncFeedbackFileCreator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
