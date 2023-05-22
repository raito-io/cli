// Code generated by mockery v2.27.1. DO NOT EDIT.

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

// SyncAccessProvidersToTarget provides a mock function with given fields: ctx, rolesToRemove, access, feedbackHandler, configMap
func (_m *MockAccessProviderRoleSyncer) SyncAccessProvidersToTarget(ctx context.Context, rolesToRemove []string, access map[string]*sync_to_target.AccessProvider, feedbackHandler wrappers.AccessProviderFeedbackHandler, configMap *config.ConfigMap) error {
	ret := _m.Called(ctx, rolesToRemove, access, feedbackHandler, configMap)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []string, map[string]*sync_to_target.AccessProvider, wrappers.AccessProviderFeedbackHandler, *config.ConfigMap) error); ok {
		r0 = rf(ctx, rolesToRemove, access, feedbackHandler, configMap)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAccessProviderRoleSyncer_SyncAccessProvidersToTarget_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SyncAccessProvidersToTarget'
type MockAccessProviderRoleSyncer_SyncAccessProvidersToTarget_Call struct {
	*mock.Call
}

// SyncAccessProvidersToTarget is a helper method to define mock.On call
//   - ctx context.Context
//   - rolesToRemove []string
//   - access map[string]*sync_to_target.AccessProvider
//   - feedbackHandler wrappers.AccessProviderFeedbackHandler
//   - configMap *config.ConfigMap
func (_e *MockAccessProviderRoleSyncer_Expecter) SyncAccessProvidersToTarget(ctx interface{}, rolesToRemove interface{}, access interface{}, feedbackHandler interface{}, configMap interface{}) *MockAccessProviderRoleSyncer_SyncAccessProvidersToTarget_Call {
	return &MockAccessProviderRoleSyncer_SyncAccessProvidersToTarget_Call{Call: _e.mock.On("SyncAccessProvidersToTarget", ctx, rolesToRemove, access, feedbackHandler, configMap)}
}

func (_c *MockAccessProviderRoleSyncer_SyncAccessProvidersToTarget_Call) Run(run func(ctx context.Context, rolesToRemove []string, access map[string]*sync_to_target.AccessProvider, feedbackHandler wrappers.AccessProviderFeedbackHandler, configMap *config.ConfigMap)) *MockAccessProviderRoleSyncer_SyncAccessProvidersToTarget_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]string), args[2].(map[string]*sync_to_target.AccessProvider), args[3].(wrappers.AccessProviderFeedbackHandler), args[4].(*config.ConfigMap))
	})
	return _c
}

func (_c *MockAccessProviderRoleSyncer_SyncAccessProvidersToTarget_Call) Return(_a0 error) *MockAccessProviderRoleSyncer_SyncAccessProvidersToTarget_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAccessProviderRoleSyncer_SyncAccessProvidersToTarget_Call) RunAndReturn(run func(context.Context, []string, map[string]*sync_to_target.AccessProvider, wrappers.AccessProviderFeedbackHandler, *config.ConfigMap) error) *MockAccessProviderRoleSyncer_SyncAccessProvidersToTarget_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockAccessProviderRoleSyncer interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockAccessProviderRoleSyncer creates a new instance of MockAccessProviderRoleSyncer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockAccessProviderRoleSyncer(t mockConstructorTestingTNewMockAccessProviderRoleSyncer) *MockAccessProviderRoleSyncer {
	mock := &MockAccessProviderRoleSyncer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
