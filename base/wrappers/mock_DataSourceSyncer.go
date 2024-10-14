// Code generated by mockery v2.46.0. DO NOT EDIT.

package wrappers

import (
	context "context"

	config "github.com/raito-io/cli/base/util/config"

	data_source "github.com/raito-io/cli/base/data_source"

	mock "github.com/stretchr/testify/mock"
)

// MockDataSourceSyncer is an autogenerated mock type for the DataSourceSyncer type
type MockDataSourceSyncer struct {
	mock.Mock
}

type MockDataSourceSyncer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDataSourceSyncer) EXPECT() *MockDataSourceSyncer_Expecter {
	return &MockDataSourceSyncer_Expecter{mock: &_m.Mock}
}

// GetDataSourceMetaData provides a mock function with given fields: ctx, configParams
func (_m *MockDataSourceSyncer) GetDataSourceMetaData(ctx context.Context, configParams *config.ConfigMap) (*data_source.MetaData, error) {
	ret := _m.Called(ctx, configParams)

	if len(ret) == 0 {
		panic("no return value specified for GetDataSourceMetaData")
	}

	var r0 *data_source.MetaData
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *config.ConfigMap) (*data_source.MetaData, error)); ok {
		return rf(ctx, configParams)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *config.ConfigMap) *data_source.MetaData); ok {
		r0 = rf(ctx, configParams)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*data_source.MetaData)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *config.ConfigMap) error); ok {
		r1 = rf(ctx, configParams)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDataSourceSyncer_GetDataSourceMetaData_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetDataSourceMetaData'
type MockDataSourceSyncer_GetDataSourceMetaData_Call struct {
	*mock.Call
}

// GetDataSourceMetaData is a helper method to define mock.On call
//   - ctx context.Context
//   - configParams *config.ConfigMap
func (_e *MockDataSourceSyncer_Expecter) GetDataSourceMetaData(ctx interface{}, configParams interface{}) *MockDataSourceSyncer_GetDataSourceMetaData_Call {
	return &MockDataSourceSyncer_GetDataSourceMetaData_Call{Call: _e.mock.On("GetDataSourceMetaData", ctx, configParams)}
}

func (_c *MockDataSourceSyncer_GetDataSourceMetaData_Call) Run(run func(ctx context.Context, configParams *config.ConfigMap)) *MockDataSourceSyncer_GetDataSourceMetaData_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*config.ConfigMap))
	})
	return _c
}

func (_c *MockDataSourceSyncer_GetDataSourceMetaData_Call) Return(_a0 *data_source.MetaData, _a1 error) *MockDataSourceSyncer_GetDataSourceMetaData_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDataSourceSyncer_GetDataSourceMetaData_Call) RunAndReturn(run func(context.Context, *config.ConfigMap) (*data_source.MetaData, error)) *MockDataSourceSyncer_GetDataSourceMetaData_Call {
	_c.Call.Return(run)
	return _c
}

// SyncDataSource provides a mock function with given fields: ctx, dataSourceHandler, _a2
func (_m *MockDataSourceSyncer) SyncDataSource(ctx context.Context, dataSourceHandler DataSourceObjectHandler, _a2 *data_source.DataSourceSyncConfig) error {
	ret := _m.Called(ctx, dataSourceHandler, _a2)

	if len(ret) == 0 {
		panic("no return value specified for SyncDataSource")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, DataSourceObjectHandler, *data_source.DataSourceSyncConfig) error); ok {
		r0 = rf(ctx, dataSourceHandler, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDataSourceSyncer_SyncDataSource_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SyncDataSource'
type MockDataSourceSyncer_SyncDataSource_Call struct {
	*mock.Call
}

// SyncDataSource is a helper method to define mock.On call
//   - ctx context.Context
//   - dataSourceHandler DataSourceObjectHandler
//   - _a2 *data_source.DataSourceSyncConfig
func (_e *MockDataSourceSyncer_Expecter) SyncDataSource(ctx interface{}, dataSourceHandler interface{}, _a2 interface{}) *MockDataSourceSyncer_SyncDataSource_Call {
	return &MockDataSourceSyncer_SyncDataSource_Call{Call: _e.mock.On("SyncDataSource", ctx, dataSourceHandler, _a2)}
}

func (_c *MockDataSourceSyncer_SyncDataSource_Call) Run(run func(ctx context.Context, dataSourceHandler DataSourceObjectHandler, _a2 *data_source.DataSourceSyncConfig)) *MockDataSourceSyncer_SyncDataSource_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(DataSourceObjectHandler), args[2].(*data_source.DataSourceSyncConfig))
	})
	return _c
}

func (_c *MockDataSourceSyncer_SyncDataSource_Call) Return(_a0 error) *MockDataSourceSyncer_SyncDataSource_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDataSourceSyncer_SyncDataSource_Call) RunAndReturn(run func(context.Context, DataSourceObjectHandler, *data_source.DataSourceSyncConfig) error) *MockDataSourceSyncer_SyncDataSource_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockDataSourceSyncer creates a new instance of MockDataSourceSyncer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDataSourceSyncer(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDataSourceSyncer {
	mock := &MockDataSourceSyncer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
