// Code generated by mockery v2.52.3. DO NOT EDIT.

package mocks

import (
	data_usage "github.com/raito-io/cli/base/data_usage"
	mock "github.com/stretchr/testify/mock"
)

// DataUsageFileCreator is an autogenerated mock type for the DataUsageFileCreator type
type DataUsageFileCreator struct {
	mock.Mock
}

type DataUsageFileCreator_Expecter struct {
	mock *mock.Mock
}

func (_m *DataUsageFileCreator) EXPECT() *DataUsageFileCreator_Expecter {
	return &DataUsageFileCreator_Expecter{mock: &_m.Mock}
}

// AddStatements provides a mock function with given fields: statements
func (_m *DataUsageFileCreator) AddStatements(statements []data_usage.Statement) error {
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

// DataUsageFileCreator_AddStatements_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddStatements'
type DataUsageFileCreator_AddStatements_Call struct {
	*mock.Call
}

// AddStatements is a helper method to define mock.On call
//   - statements []data_usage.Statement
func (_e *DataUsageFileCreator_Expecter) AddStatements(statements interface{}) *DataUsageFileCreator_AddStatements_Call {
	return &DataUsageFileCreator_AddStatements_Call{Call: _e.mock.On("AddStatements", statements)}
}

func (_c *DataUsageFileCreator_AddStatements_Call) Run(run func(statements []data_usage.Statement)) *DataUsageFileCreator_AddStatements_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]data_usage.Statement))
	})
	return _c
}

func (_c *DataUsageFileCreator_AddStatements_Call) Return(_a0 error) *DataUsageFileCreator_AddStatements_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DataUsageFileCreator_AddStatements_Call) RunAndReturn(run func([]data_usage.Statement) error) *DataUsageFileCreator_AddStatements_Call {
	_c.Call.Return(run)
	return _c
}

// Close provides a mock function with no fields
func (_m *DataUsageFileCreator) Close() {
	_m.Called()
}

// DataUsageFileCreator_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type DataUsageFileCreator_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *DataUsageFileCreator_Expecter) Close() *DataUsageFileCreator_Close_Call {
	return &DataUsageFileCreator_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *DataUsageFileCreator_Close_Call) Run(run func()) *DataUsageFileCreator_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *DataUsageFileCreator_Close_Call) Return() *DataUsageFileCreator_Close_Call {
	_c.Call.Return()
	return _c
}

func (_c *DataUsageFileCreator_Close_Call) RunAndReturn(run func()) *DataUsageFileCreator_Close_Call {
	_c.Run(run)
	return _c
}

// GetActualFileNames provides a mock function with no fields
func (_m *DataUsageFileCreator) GetActualFileNames() []string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetActualFileNames")
	}

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// DataUsageFileCreator_GetActualFileNames_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetActualFileNames'
type DataUsageFileCreator_GetActualFileNames_Call struct {
	*mock.Call
}

// GetActualFileNames is a helper method to define mock.On call
func (_e *DataUsageFileCreator_Expecter) GetActualFileNames() *DataUsageFileCreator_GetActualFileNames_Call {
	return &DataUsageFileCreator_GetActualFileNames_Call{Call: _e.mock.On("GetActualFileNames")}
}

func (_c *DataUsageFileCreator_GetActualFileNames_Call) Run(run func()) *DataUsageFileCreator_GetActualFileNames_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *DataUsageFileCreator_GetActualFileNames_Call) Return(_a0 []string) *DataUsageFileCreator_GetActualFileNames_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DataUsageFileCreator_GetActualFileNames_Call) RunAndReturn(run func() []string) *DataUsageFileCreator_GetActualFileNames_Call {
	_c.Call.Return(run)
	return _c
}

// GetImportFileSize provides a mock function with no fields
func (_m *DataUsageFileCreator) GetImportFileSize() uint64 {
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

// DataUsageFileCreator_GetImportFileSize_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetImportFileSize'
type DataUsageFileCreator_GetImportFileSize_Call struct {
	*mock.Call
}

// GetImportFileSize is a helper method to define mock.On call
func (_e *DataUsageFileCreator_Expecter) GetImportFileSize() *DataUsageFileCreator_GetImportFileSize_Call {
	return &DataUsageFileCreator_GetImportFileSize_Call{Call: _e.mock.On("GetImportFileSize")}
}

func (_c *DataUsageFileCreator_GetImportFileSize_Call) Run(run func()) *DataUsageFileCreator_GetImportFileSize_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *DataUsageFileCreator_GetImportFileSize_Call) Return(_a0 uint64) *DataUsageFileCreator_GetImportFileSize_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DataUsageFileCreator_GetImportFileSize_Call) RunAndReturn(run func() uint64) *DataUsageFileCreator_GetImportFileSize_Call {
	_c.Call.Return(run)
	return _c
}

// GetStatementCount provides a mock function with no fields
func (_m *DataUsageFileCreator) GetStatementCount() int {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetStatementCount")
	}

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// DataUsageFileCreator_GetStatementCount_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetStatementCount'
type DataUsageFileCreator_GetStatementCount_Call struct {
	*mock.Call
}

// GetStatementCount is a helper method to define mock.On call
func (_e *DataUsageFileCreator_Expecter) GetStatementCount() *DataUsageFileCreator_GetStatementCount_Call {
	return &DataUsageFileCreator_GetStatementCount_Call{Call: _e.mock.On("GetStatementCount")}
}

func (_c *DataUsageFileCreator_GetStatementCount_Call) Run(run func()) *DataUsageFileCreator_GetStatementCount_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *DataUsageFileCreator_GetStatementCount_Call) Return(_a0 int) *DataUsageFileCreator_GetStatementCount_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DataUsageFileCreator_GetStatementCount_Call) RunAndReturn(run func() int) *DataUsageFileCreator_GetStatementCount_Call {
	_c.Call.Return(run)
	return _c
}

// NewDataUsageFileCreator creates a new instance of DataUsageFileCreator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDataUsageFileCreator(t interface {
	mock.TestingT
	Cleanup(func())
}) *DataUsageFileCreator {
	mock := &DataUsageFileCreator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
