// Code generated by mockery v2.36.1. DO NOT EDIT.

package role_based

import (
	context "context"

	config "github.com/raito-io/cli/base/util/config"

	mock "github.com/stretchr/testify/mock"

	sync_to_target "github.com/raito-io/cli/base/access_provider/sync_to_target"

	wrappers "github.com/raito-io/cli/base/wrappers"
)

// MockAccessProviderRoleSyncer is an autogenerated mock type for the AccessProviderRoleSyncer type
type MockAccessProviderRoleSyncer struct {
	mock.Mock
}

type MockAccessProviderRoleSyncer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAccessProviderRoleSyncer) EXPECT() *MockAccessProviderRoleSyncer_Expecter {
	return &MockAccessProviderRoleSyncer_Expecter{mock: &_m.Mock}
}

// SyncAccessAsCodeToTarget provides a mock function with given fields: ctx, accesses, prefix, configMap
func (_m *MockAccessProviderRoleSyncer) SyncAccessAsCodeToTarget(ctx context.Context, accesses map[string]*sync_to_target.AccessProvider, prefix string, configMap *config.ConfigMap) error {
	ret := _m.Called(ctx, accesses, prefix, configMap)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, map[string]*sync_to_target.AccessProvider, string, *config.ConfigMap) error); ok {
		r0 = rf(ctx, accesses, prefix, configMap)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAccessProviderRoleSyncer_SyncAccessAsCodeToTarget_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SyncAccessAsCodeToTarget'
type MockAccessProviderRoleSyncer_SyncAccessAsCodeToTarget_Call struct {
	*mock.Call
}

// SyncAccessAsCodeToTarget is a helper method to define mock.On call
//   - ctx context.Context
//   - accesses map[string]*sync_to_target.AccessProvider
//   - prefix string
//   - configMap *config.ConfigMap
func (_e *MockAccessProviderRoleSyncer_Expecter) SyncAccessAsCodeToTarget(ctx interface{}, accesses interface{}, prefix interface{}, configMap interface{}) *MockAccessProviderRoleSyncer_SyncAccessAsCodeToTarget_Call {
	return &MockAccessProviderRoleSyncer_SyncAccessAsCodeToTarget_Call{Call: _e.mock.On("SyncAccessAsCodeToTarget", ctx, accesses, prefix, configMap)}
}

func (_c *MockAccessProviderRoleSyncer_SyncAccessAsCodeToTarget_Call) Run(run func(ctx context.Context, accesses map[string]*sync_to_target.AccessProvider, prefix string, configMap *config.ConfigMap)) *MockAccessProviderRoleSyncer_SyncAccessAsCodeToTarget_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(map[string]*sync_to_target.AccessProvider), args[2].(string), args[3].(*config.ConfigMap))
	})
	return _c
}

func (_c *MockAccessProviderRoleSyncer_SyncAccessAsCodeToTarget_Call) Return(_a0 error) *MockAccessProviderRoleSyncer_SyncAccessAsCodeToTarget_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAccessProviderRoleSyncer_SyncAccessAsCodeToTarget_Call) RunAndReturn(run func(context.Context, map[string]*sync_to_target.AccessProvider, string, *config.ConfigMap) error) *MockAccessProviderRoleSyncer_SyncAccessAsCodeToTarget_Call {
	_c.Call.Return(run)
	return _c
}

// SyncAccessProviderMasksToTarget provides a mock function with given fields: ctx, apToRemoveMap, apMap, roleNameMap, feedbackHandler, configMap
func (_m *MockAccessProviderRoleSyncer) SyncAccessProviderMasksToTarget(ctx context.Context, apToRemoveMap map[string]*sync_to_target.AccessProvider, apMap map[string]*sync_to_target.AccessProvider, roleNameMap map[string]string, feedbackHandler wrappers.AccessProviderFeedbackHandler, configMap *config.ConfigMap) error {
	ret := _m.Called(ctx, apToRemoveMap, apMap, roleNameMap, feedbackHandler, configMap)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, map[string]*sync_to_target.AccessProvider, map[string]*sync_to_target.AccessProvider, map[string]string, wrappers.AccessProviderFeedbackHandler, *config.ConfigMap) error); ok {
		r0 = rf(ctx, apToRemoveMap, apMap, roleNameMap, feedbackHandler, configMap)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAccessProviderRoleSyncer_SyncAccessProviderMasksToTarget_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SyncAccessProviderMasksToTarget'
type MockAccessProviderRoleSyncer_SyncAccessProviderMasksToTarget_Call struct {
	*mock.Call
}

// SyncAccessProviderMasksToTarget is a helper method to define mock.On call
//   - ctx context.Context
//   - apToRemoveMap map[string]*sync_to_target.AccessProvider
//   - apMap map[string]*sync_to_target.AccessProvider
//   - roleNameMap map[string]string
//   - feedbackHandler wrappers.AccessProviderFeedbackHandler
//   - configMap *config.ConfigMap
func (_e *MockAccessProviderRoleSyncer_Expecter) SyncAccessProviderMasksToTarget(ctx interface{}, apToRemoveMap interface{}, apMap interface{}, roleNameMap interface{}, feedbackHandler interface{}, configMap interface{}) *MockAccessProviderRoleSyncer_SyncAccessProviderMasksToTarget_Call {
	return &MockAccessProviderRoleSyncer_SyncAccessProviderMasksToTarget_Call{Call: _e.mock.On("SyncAccessProviderMasksToTarget", ctx, apToRemoveMap, apMap, roleNameMap, feedbackHandler, configMap)}
}

func (_c *MockAccessProviderRoleSyncer_SyncAccessProviderMasksToTarget_Call) Run(run func(ctx context.Context, apToRemoveMap map[string]*sync_to_target.AccessProvider, apMap map[string]*sync_to_target.AccessProvider, roleNameMap map[string]string, feedbackHandler wrappers.AccessProviderFeedbackHandler, configMap *config.ConfigMap)) *MockAccessProviderRoleSyncer_SyncAccessProviderMasksToTarget_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(map[string]*sync_to_target.AccessProvider), args[2].(map[string]*sync_to_target.AccessProvider), args[3].(map[string]string), args[4].(wrappers.AccessProviderFeedbackHandler), args[5].(*config.ConfigMap))
	})
	return _c
}

func (_c *MockAccessProviderRoleSyncer_SyncAccessProviderMasksToTarget_Call) Return(_a0 error) *MockAccessProviderRoleSyncer_SyncAccessProviderMasksToTarget_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAccessProviderRoleSyncer_SyncAccessProviderMasksToTarget_Call) RunAndReturn(run func(context.Context, map[string]*sync_to_target.AccessProvider, map[string]*sync_to_target.AccessProvider, map[string]string, wrappers.AccessProviderFeedbackHandler, *config.ConfigMap) error) *MockAccessProviderRoleSyncer_SyncAccessProviderMasksToTarget_Call {
	_c.Call.Return(run)
	return _c
}

// SyncAccessProviderRolesToTarget provides a mock function with given fields: ctx, apToRemoveMap, apMap, feedbackHandler, configMap
func (_m *MockAccessProviderRoleSyncer) SyncAccessProviderRolesToTarget(ctx context.Context, apToRemoveMap map[string]*sync_to_target.AccessProvider, apMap map[string]*sync_to_target.AccessProvider, feedbackHandler wrappers.AccessProviderFeedbackHandler, configMap *config.ConfigMap) error {
	ret := _m.Called(ctx, apToRemoveMap, apMap, feedbackHandler, configMap)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, map[string]*sync_to_target.AccessProvider, map[string]*sync_to_target.AccessProvider, wrappers.AccessProviderFeedbackHandler, *config.ConfigMap) error); ok {
		r0 = rf(ctx, apToRemoveMap, apMap, feedbackHandler, configMap)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAccessProviderRoleSyncer_SyncAccessProviderRolesToTarget_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SyncAccessProviderRolesToTarget'
type MockAccessProviderRoleSyncer_SyncAccessProviderRolesToTarget_Call struct {
	*mock.Call
}

// SyncAccessProviderRolesToTarget is a helper method to define mock.On call
//   - ctx context.Context
//   - apToRemoveMap map[string]*sync_to_target.AccessProvider
//   - apMap map[string]*sync_to_target.AccessProvider
//   - feedbackHandler wrappers.AccessProviderFeedbackHandler
//   - configMap *config.ConfigMap
func (_e *MockAccessProviderRoleSyncer_Expecter) SyncAccessProviderRolesToTarget(ctx interface{}, apToRemoveMap interface{}, apMap interface{}, feedbackHandler interface{}, configMap interface{}) *MockAccessProviderRoleSyncer_SyncAccessProviderRolesToTarget_Call {
	return &MockAccessProviderRoleSyncer_SyncAccessProviderRolesToTarget_Call{Call: _e.mock.On("SyncAccessProviderRolesToTarget", ctx, apToRemoveMap, apMap, feedbackHandler, configMap)}
}

func (_c *MockAccessProviderRoleSyncer_SyncAccessProviderRolesToTarget_Call) Run(run func(ctx context.Context, apToRemoveMap map[string]*sync_to_target.AccessProvider, apMap map[string]*sync_to_target.AccessProvider, feedbackHandler wrappers.AccessProviderFeedbackHandler, configMap *config.ConfigMap)) *MockAccessProviderRoleSyncer_SyncAccessProviderRolesToTarget_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(map[string]*sync_to_target.AccessProvider), args[2].(map[string]*sync_to_target.AccessProvider), args[3].(wrappers.AccessProviderFeedbackHandler), args[4].(*config.ConfigMap))
	})
	return _c
}

func (_c *MockAccessProviderRoleSyncer_SyncAccessProviderRolesToTarget_Call) Return(_a0 error) *MockAccessProviderRoleSyncer_SyncAccessProviderRolesToTarget_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAccessProviderRoleSyncer_SyncAccessProviderRolesToTarget_Call) RunAndReturn(run func(context.Context, map[string]*sync_to_target.AccessProvider, map[string]*sync_to_target.AccessProvider, wrappers.AccessProviderFeedbackHandler, *config.ConfigMap) error) *MockAccessProviderRoleSyncer_SyncAccessProviderRolesToTarget_Call {
	_c.Call.Return(run)
	return _c
}

// SyncAccessProvidersFromTarget provides a mock function with given fields: ctx, accessProviderHandler, configMap
func (_m *MockAccessProviderRoleSyncer) SyncAccessProvidersFromTarget(ctx context.Context, accessProviderHandler wrappers.AccessProviderHandler, configMap *config.ConfigMap) error {
	ret := _m.Called(ctx, accessProviderHandler, configMap)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, wrappers.AccessProviderHandler, *config.ConfigMap) error); ok {
		r0 = rf(ctx, accessProviderHandler, configMap)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAccessProviderRoleSyncer_SyncAccessProvidersFromTarget_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SyncAccessProvidersFromTarget'
type MockAccessProviderRoleSyncer_SyncAccessProvidersFromTarget_Call struct {
	*mock.Call
}

// SyncAccessProvidersFromTarget is a helper method to define mock.On call
//   - ctx context.Context
//   - accessProviderHandler wrappers.AccessProviderHandler
//   - configMap *config.ConfigMap
func (_e *MockAccessProviderRoleSyncer_Expecter) SyncAccessProvidersFromTarget(ctx interface{}, accessProviderHandler interface{}, configMap interface{}) *MockAccessProviderRoleSyncer_SyncAccessProvidersFromTarget_Call {
	return &MockAccessProviderRoleSyncer_SyncAccessProvidersFromTarget_Call{Call: _e.mock.On("SyncAccessProvidersFromTarget", ctx, accessProviderHandler, configMap)}
}

func (_c *MockAccessProviderRoleSyncer_SyncAccessProvidersFromTarget_Call) Run(run func(ctx context.Context, accessProviderHandler wrappers.AccessProviderHandler, configMap *config.ConfigMap)) *MockAccessProviderRoleSyncer_SyncAccessProvidersFromTarget_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(wrappers.AccessProviderHandler), args[2].(*config.ConfigMap))
	})
	return _c
}

func (_c *MockAccessProviderRoleSyncer_SyncAccessProvidersFromTarget_Call) Return(_a0 error) *MockAccessProviderRoleSyncer_SyncAccessProvidersFromTarget_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAccessProviderRoleSyncer_SyncAccessProvidersFromTarget_Call) RunAndReturn(run func(context.Context, wrappers.AccessProviderHandler, *config.ConfigMap) error) *MockAccessProviderRoleSyncer_SyncAccessProvidersFromTarget_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockAccessProviderRoleSyncer creates a new instance of MockAccessProviderRoleSyncer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAccessProviderRoleSyncer(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAccessProviderRoleSyncer {
	mock := &MockAccessProviderRoleSyncer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
