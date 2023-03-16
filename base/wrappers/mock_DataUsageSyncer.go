// Code generated by mockery v2.21.1. DO NOT EDIT.

package wrappers

import (
	context "context"

	config "github.com/raito-io/cli/base/util/config"

	mock "github.com/stretchr/testify/mock"
)

// MockDataUsageSyncer is an autogenerated mock type for the DataUsageSyncer type
type MockDataUsageSyncer struct {
	mock.Mock
}

type MockDataUsageSyncer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDataUsageSyncer) EXPECT() *MockDataUsageSyncer_Expecter {
	return &MockDataUsageSyncer_Expecter{mock: &_m.Mock}
}

// SyncDataUsage provides a mock function with given fields: ctx, fileCreator, configParams
func (_m *MockDataUsageSyncer) SyncDataUsage(ctx context.Context, fileCreator DataUsageStatementHandler, configParams *config.ConfigMap) error {
	ret := _m.Called(ctx, fileCreator, configParams)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, DataUsageStatementHandler, *config.ConfigMap) error); ok {
		r0 = rf(ctx, fileCreator, configParams)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDataUsageSyncer_SyncDataUsage_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SyncDataUsage'
type MockDataUsageSyncer_SyncDataUsage_Call struct {
	*mock.Call
}

// SyncDataUsage is a helper method to define mock.On call
//   - ctx context.Context
//   - fileCreator DataUsageStatementHandler
//   - configParams *config.ConfigMap
func (_e *MockDataUsageSyncer_Expecter) SyncDataUsage(ctx interface{}, fileCreator interface{}, configParams interface{}) *MockDataUsageSyncer_SyncDataUsage_Call {
	return &MockDataUsageSyncer_SyncDataUsage_Call{Call: _e.mock.On("SyncDataUsage", ctx, fileCreator, configParams)}
}

func (_c *MockDataUsageSyncer_SyncDataUsage_Call) Run(run func(ctx context.Context, fileCreator DataUsageStatementHandler, configParams *config.ConfigMap)) *MockDataUsageSyncer_SyncDataUsage_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(DataUsageStatementHandler), args[2].(*config.ConfigMap))
	})
	return _c
}

func (_c *MockDataUsageSyncer_SyncDataUsage_Call) Return(_a0 error) *MockDataUsageSyncer_SyncDataUsage_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDataUsageSyncer_SyncDataUsage_Call) RunAndReturn(run func(context.Context, DataUsageStatementHandler, *config.ConfigMap) error) *MockDataUsageSyncer_SyncDataUsage_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockDataUsageSyncer interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockDataUsageSyncer creates a new instance of MockDataUsageSyncer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockDataUsageSyncer(t mockConstructorTestingTNewMockDataUsageSyncer) *MockDataUsageSyncer {
	mock := &MockDataUsageSyncer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
