// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	data_usage "github.com/raito-io/cli/base/data_usage"
	mock "github.com/stretchr/testify/mock"
)

// DataUsageStatementHandler is an autogenerated mock type for the DataUsageStatementHandler type
type DataUsageStatementHandler struct {
	mock.Mock
}

type DataUsageStatementHandler_Expecter struct {
	mock *mock.Mock
}

func (_m *DataUsageStatementHandler) EXPECT() *DataUsageStatementHandler_Expecter {
	return &DataUsageStatementHandler_Expecter{mock: &_m.Mock}
}

// AddStatements provides a mock function with given fields: statements
func (_m *DataUsageStatementHandler) AddStatements(statements []data_usage.Statement) error {
	ret := _m.Called(statements)

	if len(ret) == 0 {
		panic("no return value specified for AddStatements")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]data_usage.Statement) error); ok {
		r0 = rf(statements)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DataUsageStatementHandler_AddStatements_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddStatements'
type DataUsageStatementHandler_AddStatements_Call struct {
	*mock.Call
}

// AddStatements is a helper method to define mock.On call
//   - statements []data_usage.Statement
func (_e *DataUsageStatementHandler_Expecter) AddStatements(statements interface{}) *DataUsageStatementHandler_AddStatements_Call {
	return &DataUsageStatementHandler_AddStatements_Call{Call: _e.mock.On("AddStatements", statements)}
}

func (_c *DataUsageStatementHandler_AddStatements_Call) Run(run func(statements []data_usage.Statement)) *DataUsageStatementHandler_AddStatements_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]data_usage.Statement))
	})
	return _c
}

func (_c *DataUsageStatementHandler_AddStatements_Call) Return(_a0 error) *DataUsageStatementHandler_AddStatements_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DataUsageStatementHandler_AddStatements_Call) RunAndReturn(run func([]data_usage.Statement) error) *DataUsageStatementHandler_AddStatements_Call {
	_c.Call.Return(run)
	return _c
}

// GetImportFileSize provides a mock function with given fields:
func (_m *DataUsageStatementHandler) GetImportFileSize() uint64 {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetImportFileSize")
	}

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}

// DataUsageStatementHandler_GetImportFileSize_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetImportFileSize'
type DataUsageStatementHandler_GetImportFileSize_Call struct {
	*mock.Call
}

// GetImportFileSize is a helper method to define mock.On call
func (_e *DataUsageStatementHandler_Expecter) GetImportFileSize() *DataUsageStatementHandler_GetImportFileSize_Call {
	return &DataUsageStatementHandler_GetImportFileSize_Call{Call: _e.mock.On("GetImportFileSize")}
}

func (_c *DataUsageStatementHandler_GetImportFileSize_Call) Run(run func()) *DataUsageStatementHandler_GetImportFileSize_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *DataUsageStatementHandler_GetImportFileSize_Call) Return(_a0 uint64) *DataUsageStatementHandler_GetImportFileSize_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DataUsageStatementHandler_GetImportFileSize_Call) RunAndReturn(run func() uint64) *DataUsageStatementHandler_GetImportFileSize_Call {
	_c.Call.Return(run)
	return _c
}

// NewDataUsageStatementHandler creates a new instance of DataUsageStatementHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDataUsageStatementHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *DataUsageStatementHandler {
	mock := &DataUsageStatementHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
