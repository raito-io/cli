// Code generated by mockery v2.36.1. DO NOT EDIT.

package wrappers

import (
	context "context"

	config "github.com/raito-io/cli/base/util/config"

	mock "github.com/stretchr/testify/mock"

	sync_to_target "github.com/raito-io/cli/base/access_provider/sync_to_target"
)

// MockAccessProviderSyncer is an autogenerated mock type for the AccessProviderSyncer type
type MockAccessProviderSyncer struct {
	mock.Mock
}

type MockAccessProviderSyncer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAccessProviderSyncer) EXPECT() *MockAccessProviderSyncer_Expecter {
	return &MockAccessProviderSyncer_Expecter{mock: &_m.Mock}
}

// SyncAccessAsCodeToTarget provides a mock function with given fields: ctx, accessProviders, prefix, configMap
func (_m *MockAccessProviderSyncer) SyncAccessAsCodeToTarget(ctx context.Context, accessProviders *sync_to_target.AccessProviderImport, prefix string, configMap *config.ConfigMap) error {
	ret := _m.Called(ctx, accessProviders, prefix, configMap)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *sync_to_target.AccessProviderImport, string, *config.ConfigMap) error); ok {
		r0 = rf(ctx, accessProviders, prefix, configMap)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAccessProviderSyncer_SyncAccessAsCodeToTarget_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SyncAccessAsCodeToTarget'
type MockAccessProviderSyncer_SyncAccessAsCodeToTarget_Call struct {
	*mock.Call
}

// SyncAccessAsCodeToTarget is a helper method to define mock.On call
//   - ctx context.Context
//   - accessProviders *sync_to_target.AccessProviderImport
//   - prefix string
//   - configMap *config.ConfigMap
func (_e *MockAccessProviderSyncer_Expecter) SyncAccessAsCodeToTarget(ctx interface{}, accessProviders interface{}, prefix interface{}, configMap interface{}) *MockAccessProviderSyncer_SyncAccessAsCodeToTarget_Call {
	return &MockAccessProviderSyncer_SyncAccessAsCodeToTarget_Call{Call: _e.mock.On("SyncAccessAsCodeToTarget", ctx, accessProviders, prefix, configMap)}
}

func (_c *MockAccessProviderSyncer_SyncAccessAsCodeToTarget_Call) Run(run func(ctx context.Context, accessProviders *sync_to_target.AccessProviderImport, prefix string, configMap *config.ConfigMap)) *MockAccessProviderSyncer_SyncAccessAsCodeToTarget_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*sync_to_target.AccessProviderImport), args[2].(string), args[3].(*config.ConfigMap))
	})
	return _c
}

func (_c *MockAccessProviderSyncer_SyncAccessAsCodeToTarget_Call) Return(_a0 error) *MockAccessProviderSyncer_SyncAccessAsCodeToTarget_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAccessProviderSyncer_SyncAccessAsCodeToTarget_Call) RunAndReturn(run func(context.Context, *sync_to_target.AccessProviderImport, string, *config.ConfigMap) error) *MockAccessProviderSyncer_SyncAccessAsCodeToTarget_Call {
	_c.Call.Return(run)
	return _c
}

// SyncAccessProviderToTarget provides a mock function with given fields: ctx, accessProviders, accessProviderFeedbackHandler, configMap
func (_m *MockAccessProviderSyncer) SyncAccessProviderToTarget(ctx context.Context, accessProviders *sync_to_target.AccessProviderImport, accessProviderFeedbackHandler AccessProviderFeedbackHandler, configMap *config.ConfigMap) error {
	ret := _m.Called(ctx, accessProviders, accessProviderFeedbackHandler, configMap)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *sync_to_target.AccessProviderImport, AccessProviderFeedbackHandler, *config.ConfigMap) error); ok {
		r0 = rf(ctx, accessProviders, accessProviderFeedbackHandler, configMap)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAccessProviderSyncer_SyncAccessProviderToTarget_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SyncAccessProviderToTarget'
type MockAccessProviderSyncer_SyncAccessProviderToTarget_Call struct {
	*mock.Call
}

// SyncAccessProviderToTarget is a helper method to define mock.On call
//   - ctx context.Context
//   - accessProviders *sync_to_target.AccessProviderImport
//   - accessProviderFeedbackHandler AccessProviderFeedbackHandler
//   - configMap *config.ConfigMap
func (_e *MockAccessProviderSyncer_Expecter) SyncAccessProviderToTarget(ctx interface{}, accessProviders interface{}, accessProviderFeedbackHandler interface{}, configMap interface{}) *MockAccessProviderSyncer_SyncAccessProviderToTarget_Call {
	return &MockAccessProviderSyncer_SyncAccessProviderToTarget_Call{Call: _e.mock.On("SyncAccessProviderToTarget", ctx, accessProviders, accessProviderFeedbackHandler, configMap)}
}

func (_c *MockAccessProviderSyncer_SyncAccessProviderToTarget_Call) Run(run func(ctx context.Context, accessProviders *sync_to_target.AccessProviderImport, accessProviderFeedbackHandler AccessProviderFeedbackHandler, configMap *config.ConfigMap)) *MockAccessProviderSyncer_SyncAccessProviderToTarget_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*sync_to_target.AccessProviderImport), args[2].(AccessProviderFeedbackHandler), args[3].(*config.ConfigMap))
	})
	return _c
}

func (_c *MockAccessProviderSyncer_SyncAccessProviderToTarget_Call) Return(_a0 error) *MockAccessProviderSyncer_SyncAccessProviderToTarget_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAccessProviderSyncer_SyncAccessProviderToTarget_Call) RunAndReturn(run func(context.Context, *sync_to_target.AccessProviderImport, AccessProviderFeedbackHandler, *config.ConfigMap) error) *MockAccessProviderSyncer_SyncAccessProviderToTarget_Call {
	_c.Call.Return(run)
	return _c
}

// SyncAccessProvidersFromTarget provides a mock function with given fields: ctx, accessProviderHandler, configMap
func (_m *MockAccessProviderSyncer) SyncAccessProvidersFromTarget(ctx context.Context, accessProviderHandler AccessProviderHandler, configMap *config.ConfigMap) error {
	ret := _m.Called(ctx, accessProviderHandler, configMap)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, AccessProviderHandler, *config.ConfigMap) error); ok {
		r0 = rf(ctx, accessProviderHandler, configMap)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAccessProviderSyncer_SyncAccessProvidersFromTarget_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SyncAccessProvidersFromTarget'
type MockAccessProviderSyncer_SyncAccessProvidersFromTarget_Call struct {
	*mock.Call
}

// SyncAccessProvidersFromTarget is a helper method to define mock.On call
//   - ctx context.Context
//   - accessProviderHandler AccessProviderHandler
//   - configMap *config.ConfigMap
func (_e *MockAccessProviderSyncer_Expecter) SyncAccessProvidersFromTarget(ctx interface{}, accessProviderHandler interface{}, configMap interface{}) *MockAccessProviderSyncer_SyncAccessProvidersFromTarget_Call {
	return &MockAccessProviderSyncer_SyncAccessProvidersFromTarget_Call{Call: _e.mock.On("SyncAccessProvidersFromTarget", ctx, accessProviderHandler, configMap)}
}

func (_c *MockAccessProviderSyncer_SyncAccessProvidersFromTarget_Call) Run(run func(ctx context.Context, accessProviderHandler AccessProviderHandler, configMap *config.ConfigMap)) *MockAccessProviderSyncer_SyncAccessProvidersFromTarget_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(AccessProviderHandler), args[2].(*config.ConfigMap))
	})
	return _c
}

func (_c *MockAccessProviderSyncer_SyncAccessProvidersFromTarget_Call) Return(_a0 error) *MockAccessProviderSyncer_SyncAccessProvidersFromTarget_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAccessProviderSyncer_SyncAccessProvidersFromTarget_Call) RunAndReturn(run func(context.Context, AccessProviderHandler, *config.ConfigMap) error) *MockAccessProviderSyncer_SyncAccessProvidersFromTarget_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockAccessProviderSyncer creates a new instance of MockAccessProviderSyncer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAccessProviderSyncer(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAccessProviderSyncer {
	mock := &MockAccessProviderSyncer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
