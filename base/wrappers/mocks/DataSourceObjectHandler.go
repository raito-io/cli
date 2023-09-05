// Code generated by mockery v2.33.1. DO NOT EDIT.

package mocks

import (
	data_source "github.com/raito-io/cli/base/data_source"
	mock "github.com/stretchr/testify/mock"
)

// DataSourceObjectHandler is an autogenerated mock type for the DataSourceObjectHandler type
type DataSourceObjectHandler struct {
	mock.Mock
}

type DataSourceObjectHandler_Expecter struct {
	mock *mock.Mock
}

func (_m *DataSourceObjectHandler) EXPECT() *DataSourceObjectHandler_Expecter {
	return &DataSourceObjectHandler_Expecter{mock: &_m.Mock}
}

// AddDataObjects provides a mock function with given fields: dataObjects
func (_m *DataSourceObjectHandler) AddDataObjects(dataObjects ...*data_source.DataObject) error {
	_va := make([]interface{}, len(dataObjects))
	for _i := range dataObjects {
		_va[_i] = dataObjects[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(...*data_source.DataObject) error); ok {
		r0 = rf(dataObjects...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DataSourceObjectHandler_AddDataObjects_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddDataObjects'
type DataSourceObjectHandler_AddDataObjects_Call struct {
	*mock.Call
}

// AddDataObjects is a helper method to define mock.On call
//   - dataObjects ...*data_source.DataObject
func (_e *DataSourceObjectHandler_Expecter) AddDataObjects(dataObjects ...interface{}) *DataSourceObjectHandler_AddDataObjects_Call {
	return &DataSourceObjectHandler_AddDataObjects_Call{Call: _e.mock.On("AddDataObjects",
		append([]interface{}{}, dataObjects...)...)}
}

func (_c *DataSourceObjectHandler_AddDataObjects_Call) Run(run func(dataObjects ...*data_source.DataObject)) *DataSourceObjectHandler_AddDataObjects_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]*data_source.DataObject, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(*data_source.DataObject)
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *DataSourceObjectHandler_AddDataObjects_Call) Return(_a0 error) *DataSourceObjectHandler_AddDataObjects_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DataSourceObjectHandler_AddDataObjects_Call) RunAndReturn(run func(...*data_source.DataObject) error) *DataSourceObjectHandler_AddDataObjects_Call {
	_c.Call.Return(run)
	return _c
}

// SetDataSourceDescription provides a mock function with given fields: desc
func (_m *DataSourceObjectHandler) SetDataSourceDescription(desc string) {
	_m.Called(desc)
}

// DataSourceObjectHandler_SetDataSourceDescription_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetDataSourceDescription'
type DataSourceObjectHandler_SetDataSourceDescription_Call struct {
	*mock.Call
}

// SetDataSourceDescription is a helper method to define mock.On call
//   - desc string
func (_e *DataSourceObjectHandler_Expecter) SetDataSourceDescription(desc interface{}) *DataSourceObjectHandler_SetDataSourceDescription_Call {
	return &DataSourceObjectHandler_SetDataSourceDescription_Call{Call: _e.mock.On("SetDataSourceDescription", desc)}
}

func (_c *DataSourceObjectHandler_SetDataSourceDescription_Call) Run(run func(desc string)) *DataSourceObjectHandler_SetDataSourceDescription_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *DataSourceObjectHandler_SetDataSourceDescription_Call) Return() *DataSourceObjectHandler_SetDataSourceDescription_Call {
	_c.Call.Return()
	return _c
}

func (_c *DataSourceObjectHandler_SetDataSourceDescription_Call) RunAndReturn(run func(string)) *DataSourceObjectHandler_SetDataSourceDescription_Call {
	_c.Call.Return(run)
	return _c
}

// SetDataSourceFullname provides a mock function with given fields: name
func (_m *DataSourceObjectHandler) SetDataSourceFullname(name string) {
	_m.Called(name)
}

// DataSourceObjectHandler_SetDataSourceFullname_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetDataSourceFullname'
type DataSourceObjectHandler_SetDataSourceFullname_Call struct {
	*mock.Call
}

// SetDataSourceFullname is a helper method to define mock.On call
//   - name string
func (_e *DataSourceObjectHandler_Expecter) SetDataSourceFullname(name interface{}) *DataSourceObjectHandler_SetDataSourceFullname_Call {
	return &DataSourceObjectHandler_SetDataSourceFullname_Call{Call: _e.mock.On("SetDataSourceFullname", name)}
}

func (_c *DataSourceObjectHandler_SetDataSourceFullname_Call) Run(run func(name string)) *DataSourceObjectHandler_SetDataSourceFullname_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *DataSourceObjectHandler_SetDataSourceFullname_Call) Return() *DataSourceObjectHandler_SetDataSourceFullname_Call {
	_c.Call.Return()
	return _c
}

func (_c *DataSourceObjectHandler_SetDataSourceFullname_Call) RunAndReturn(run func(string)) *DataSourceObjectHandler_SetDataSourceFullname_Call {
	_c.Call.Return(run)
	return _c
}

// SetDataSourceName provides a mock function with given fields: name
func (_m *DataSourceObjectHandler) SetDataSourceName(name string) {
	_m.Called(name)
}

// DataSourceObjectHandler_SetDataSourceName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetDataSourceName'
type DataSourceObjectHandler_SetDataSourceName_Call struct {
	*mock.Call
}

// SetDataSourceName is a helper method to define mock.On call
//   - name string
func (_e *DataSourceObjectHandler_Expecter) SetDataSourceName(name interface{}) *DataSourceObjectHandler_SetDataSourceName_Call {
	return &DataSourceObjectHandler_SetDataSourceName_Call{Call: _e.mock.On("SetDataSourceName", name)}
}

func (_c *DataSourceObjectHandler_SetDataSourceName_Call) Run(run func(name string)) *DataSourceObjectHandler_SetDataSourceName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *DataSourceObjectHandler_SetDataSourceName_Call) Return() *DataSourceObjectHandler_SetDataSourceName_Call {
	_c.Call.Return()
	return _c
}

func (_c *DataSourceObjectHandler_SetDataSourceName_Call) RunAndReturn(run func(string)) *DataSourceObjectHandler_SetDataSourceName_Call {
	_c.Call.Return(run)
	return _c
}

// NewDataSourceObjectHandler creates a new instance of DataSourceObjectHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDataSourceObjectHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *DataSourceObjectHandler {
	mock := &DataSourceObjectHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
