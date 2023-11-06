// Code generated by mockery v2.36.1. DO NOT EDIT.

package mocks

import (
	sync_to_target "github.com/raito-io/cli/base/access_provider/sync_to_target"
	mock "github.com/stretchr/testify/mock"
)

// AccessProviderImportFileParser is an autogenerated mock type for the AccessProviderImportFileParser type
type AccessProviderImportFileParser struct {
	mock.Mock
}

type AccessProviderImportFileParser_Expecter struct {
	mock *mock.Mock
}

func (_m *AccessProviderImportFileParser) EXPECT() *AccessProviderImportFileParser_Expecter {
	return &AccessProviderImportFileParser_Expecter{mock: &_m.Mock}
}

// ParseAccessProviders provides a mock function with given fields:
func (_m *AccessProviderImportFileParser) ParseAccessProviders() (*sync_to_target.AccessProviderImport, error) {
	ret := _m.Called()

	var r0 *sync_to_target.AccessProviderImport
	var r1 error
	if rf, ok := ret.Get(0).(func() (*sync_to_target.AccessProviderImport, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *sync_to_target.AccessProviderImport); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sync_to_target.AccessProviderImport)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AccessProviderImportFileParser_ParseAccessProviders_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ParseAccessProviders'
type AccessProviderImportFileParser_ParseAccessProviders_Call struct {
	*mock.Call
}

// ParseAccessProviders is a helper method to define mock.On call
func (_e *AccessProviderImportFileParser_Expecter) ParseAccessProviders() *AccessProviderImportFileParser_ParseAccessProviders_Call {
	return &AccessProviderImportFileParser_ParseAccessProviders_Call{Call: _e.mock.On("ParseAccessProviders")}
}

func (_c *AccessProviderImportFileParser_ParseAccessProviders_Call) Run(run func()) *AccessProviderImportFileParser_ParseAccessProviders_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *AccessProviderImportFileParser_ParseAccessProviders_Call) Return(_a0 *sync_to_target.AccessProviderImport, _a1 error) *AccessProviderImportFileParser_ParseAccessProviders_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AccessProviderImportFileParser_ParseAccessProviders_Call) RunAndReturn(run func() (*sync_to_target.AccessProviderImport, error)) *AccessProviderImportFileParser_ParseAccessProviders_Call {
	_c.Call.Return(run)
	return _c
}

// NewAccessProviderImportFileParser creates a new instance of AccessProviderImportFileParser. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAccessProviderImportFileParser(t interface {
	mock.TestingT
	Cleanup(func())
}) *AccessProviderImportFileParser {
	mock := &AccessProviderImportFileParser{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
