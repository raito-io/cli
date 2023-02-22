// Code generated by mockery v2.20.0. DO NOT EDIT.

package wrappers

import (
	context "context"

	config "github.com/raito-io/cli/base/util/config"

	identity_store "github.com/raito-io/cli/base/identity_store"

	mock "github.com/stretchr/testify/mock"
)

// MockIdentityStoreSyncer is an autogenerated mock type for the IdentityStoreSyncer type
type MockIdentityStoreSyncer struct {
	mock.Mock
}

type MockIdentityStoreSyncer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockIdentityStoreSyncer) EXPECT() *MockIdentityStoreSyncer_Expecter {
	return &MockIdentityStoreSyncer_Expecter{mock: &_m.Mock}
}

// GetIdentityStoreMetaData provides a mock function with given fields:
func (_m *MockIdentityStoreSyncer) GetIdentityStoreMetaData() identity_store.MetaData {
	ret := _m.Called()

	var r0 identity_store.MetaData
	if rf, ok := ret.Get(0).(func() identity_store.MetaData); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(identity_store.MetaData)
	}

	return r0
}

// MockIdentityStoreSyncer_GetIdentityStoreMetaData_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetIdentityStoreMetaData'
type MockIdentityStoreSyncer_GetIdentityStoreMetaData_Call struct {
	*mock.Call
}

// GetIdentityStoreMetaData is a helper method to define mock.On call
func (_e *MockIdentityStoreSyncer_Expecter) GetIdentityStoreMetaData() *MockIdentityStoreSyncer_GetIdentityStoreMetaData_Call {
	return &MockIdentityStoreSyncer_GetIdentityStoreMetaData_Call{Call: _e.mock.On("GetIdentityStoreMetaData")}
}

func (_c *MockIdentityStoreSyncer_GetIdentityStoreMetaData_Call) Run(run func()) *MockIdentityStoreSyncer_GetIdentityStoreMetaData_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockIdentityStoreSyncer_GetIdentityStoreMetaData_Call) Return(_a0 identity_store.MetaData) *MockIdentityStoreSyncer_GetIdentityStoreMetaData_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIdentityStoreSyncer_GetIdentityStoreMetaData_Call) RunAndReturn(run func() identity_store.MetaData) *MockIdentityStoreSyncer_GetIdentityStoreMetaData_Call {
	_c.Call.Return(run)
	return _c
}

// SyncIdentityStore provides a mock function with given fields: ctx, identityHandler, configMap
func (_m *MockIdentityStoreSyncer) SyncIdentityStore(ctx context.Context, identityHandler IdentityStoreIdentityHandler, configMap *config.ConfigMap) error {
	ret := _m.Called(ctx, identityHandler, configMap)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, IdentityStoreIdentityHandler, *config.ConfigMap) error); ok {
		r0 = rf(ctx, identityHandler, configMap)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockIdentityStoreSyncer_SyncIdentityStore_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SyncIdentityStore'
type MockIdentityStoreSyncer_SyncIdentityStore_Call struct {
	*mock.Call
}

// SyncIdentityStore is a helper method to define mock.On call
//   - ctx context.Context
//   - identityHandler IdentityStoreIdentityHandler
//   - configMap *config.ConfigMap
func (_e *MockIdentityStoreSyncer_Expecter) SyncIdentityStore(ctx interface{}, identityHandler interface{}, configMap interface{}) *MockIdentityStoreSyncer_SyncIdentityStore_Call {
	return &MockIdentityStoreSyncer_SyncIdentityStore_Call{Call: _e.mock.On("SyncIdentityStore", ctx, identityHandler, configMap)}
}

func (_c *MockIdentityStoreSyncer_SyncIdentityStore_Call) Run(run func(ctx context.Context, identityHandler IdentityStoreIdentityHandler, configMap *config.ConfigMap)) *MockIdentityStoreSyncer_SyncIdentityStore_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(IdentityStoreIdentityHandler), args[2].(*config.ConfigMap))
	})
	return _c
}

func (_c *MockIdentityStoreSyncer_SyncIdentityStore_Call) Return(_a0 error) *MockIdentityStoreSyncer_SyncIdentityStore_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIdentityStoreSyncer_SyncIdentityStore_Call) RunAndReturn(run func(context.Context, IdentityStoreIdentityHandler, *config.ConfigMap) error) *MockIdentityStoreSyncer_SyncIdentityStore_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockIdentityStoreSyncer interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockIdentityStoreSyncer creates a new instance of MockIdentityStoreSyncer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockIdentityStoreSyncer(t mockConstructorTestingTNewMockIdentityStoreSyncer) *MockIdentityStoreSyncer {
	mock := &MockIdentityStoreSyncer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
