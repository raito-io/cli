// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	plugin "github.com/raito-io/cli/base/util/plugin"
	mock "github.com/stretchr/testify/mock"

	semver "github.com/Masterminds/semver/v3"
)

// Info is an autogenerated mock type for the Info type
type Info struct {
	mock.Mock
}

type Info_Expecter struct {
	mock *mock.Mock
}

func (_m *Info) EXPECT() *Info_Expecter {
	return &Info_Expecter{mock: &_m.Mock}
}

// CliBuildVersion provides a mock function with given fields:
func (_m *Info) CliBuildVersion() semver.Version {
	ret := _m.Called()

	var r0 semver.Version
	if rf, ok := ret.Get(0).(func() semver.Version); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(semver.Version)
	}

	return r0
}

// Info_CliBuildVersion_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CliBuildVersion'
type Info_CliBuildVersion_Call struct {
	*mock.Call
}

// CliBuildVersion is a helper method to define mock.On call
func (_e *Info_Expecter) CliBuildVersion() *Info_CliBuildVersion_Call {
	return &Info_CliBuildVersion_Call{Call: _e.mock.On("CliBuildVersion")}
}

func (_c *Info_CliBuildVersion_Call) Run(run func()) *Info_CliBuildVersion_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Info_CliBuildVersion_Call) Return(_a0 semver.Version) *Info_CliBuildVersion_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Info_CliBuildVersion_Call) RunAndReturn(run func() semver.Version) *Info_CliBuildVersion_Call {
	_c.Call.Return(run)
	return _c
}

// CliMinimalVersion provides a mock function with given fields:
func (_m *Info) CliMinimalVersion() semver.Version {
	ret := _m.Called()

	var r0 semver.Version
	if rf, ok := ret.Get(0).(func() semver.Version); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(semver.Version)
	}

	return r0
}

// Info_CliMinimalVersion_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CliMinimalVersion'
type Info_CliMinimalVersion_Call struct {
	*mock.Call
}

// CliMinimalVersion is a helper method to define mock.On call
func (_e *Info_Expecter) CliMinimalVersion() *Info_CliMinimalVersion_Call {
	return &Info_CliMinimalVersion_Call{Call: _e.mock.On("CliMinimalVersion")}
}

func (_c *Info_CliMinimalVersion_Call) Run(run func()) *Info_CliMinimalVersion_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Info_CliMinimalVersion_Call) Return(_a0 semver.Version) *Info_CliMinimalVersion_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Info_CliMinimalVersion_Call) RunAndReturn(run func() semver.Version) *Info_CliMinimalVersion_Call {
	_c.Call.Return(run)
	return _c
}

// PluginInfo provides a mock function with given fields:
func (_m *Info) PluginInfo() plugin.PluginInfo {
	ret := _m.Called()

	var r0 plugin.PluginInfo
	if rf, ok := ret.Get(0).(func() plugin.PluginInfo); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(plugin.PluginInfo)
	}

	return r0
}

// Info_PluginInfo_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PluginInfo'
type Info_PluginInfo_Call struct {
	*mock.Call
}

// PluginInfo is a helper method to define mock.On call
func (_e *Info_Expecter) PluginInfo() *Info_PluginInfo_Call {
	return &Info_PluginInfo_Call{Call: _e.mock.On("PluginInfo")}
}

func (_c *Info_PluginInfo_Call) Run(run func()) *Info_PluginInfo_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Info_PluginInfo_Call) Return(_a0 plugin.PluginInfo) *Info_PluginInfo_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Info_PluginInfo_Call) RunAndReturn(run func() plugin.PluginInfo) *Info_PluginInfo_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewInfo interface {
	mock.TestingT
	Cleanup(func())
}

// NewInfo creates a new instance of Info. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewInfo(t mockConstructorTestingTNewInfo) *Info {
	mock := &Info{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
