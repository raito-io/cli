// Code generated by mockery v2.52.3. DO NOT EDIT.

package mocks

import (
	tag "github.com/raito-io/cli/base/tag"
	mock "github.com/stretchr/testify/mock"
)

// TagFileCreator is an autogenerated mock type for the TagFileCreator type
type TagFileCreator struct {
	mock.Mock
}

type TagFileCreator_Expecter struct {
	mock *mock.Mock
}

func (_m *TagFileCreator) EXPECT() *TagFileCreator_Expecter {
	return &TagFileCreator_Expecter{mock: &_m.Mock}
}

// AddTags provides a mock function with given fields: tags
func (_m *TagFileCreator) AddTags(tags ...*tag.TagImportObject) error {
	_va := make([]interface{}, len(tags))
	for _i := range tags {
		_va[_i] = tags[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for AddTags")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(...*tag.TagImportObject) error); ok {
		r0 = rf(tags...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TagFileCreator_AddTags_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddTags'
type TagFileCreator_AddTags_Call struct {
	*mock.Call
}

// AddTags is a helper method to define mock.On call
//   - tags ...*tag.TagImportObject
func (_e *TagFileCreator_Expecter) AddTags(tags ...interface{}) *TagFileCreator_AddTags_Call {
	return &TagFileCreator_AddTags_Call{Call: _e.mock.On("AddTags",
		append([]interface{}{}, tags...)...)}
}

func (_c *TagFileCreator_AddTags_Call) Run(run func(tags ...*tag.TagImportObject)) *TagFileCreator_AddTags_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]*tag.TagImportObject, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(*tag.TagImportObject)
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *TagFileCreator_AddTags_Call) Return(_a0 error) *TagFileCreator_AddTags_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TagFileCreator_AddTags_Call) RunAndReturn(run func(...*tag.TagImportObject) error) *TagFileCreator_AddTags_Call {
	_c.Call.Return(run)
	return _c
}

// Close provides a mock function with no fields
func (_m *TagFileCreator) Close() {
	_m.Called()
}

// TagFileCreator_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type TagFileCreator_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *TagFileCreator_Expecter) Close() *TagFileCreator_Close_Call {
	return &TagFileCreator_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *TagFileCreator_Close_Call) Run(run func()) *TagFileCreator_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *TagFileCreator_Close_Call) Return() *TagFileCreator_Close_Call {
	_c.Call.Return()
	return _c
}

func (_c *TagFileCreator_Close_Call) RunAndReturn(run func()) *TagFileCreator_Close_Call {
	_c.Run(run)
	return _c
}

// GetTagCount provides a mock function with no fields
func (_m *TagFileCreator) GetTagCount() int {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetTagCount")
	}

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// TagFileCreator_GetTagCount_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTagCount'
type TagFileCreator_GetTagCount_Call struct {
	*mock.Call
}

// GetTagCount is a helper method to define mock.On call
func (_e *TagFileCreator_Expecter) GetTagCount() *TagFileCreator_GetTagCount_Call {
	return &TagFileCreator_GetTagCount_Call{Call: _e.mock.On("GetTagCount")}
}

func (_c *TagFileCreator_GetTagCount_Call) Run(run func()) *TagFileCreator_GetTagCount_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *TagFileCreator_GetTagCount_Call) Return(_a0 int) *TagFileCreator_GetTagCount_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TagFileCreator_GetTagCount_Call) RunAndReturn(run func() int) *TagFileCreator_GetTagCount_Call {
	_c.Call.Return(run)
	return _c
}

// NewTagFileCreator creates a new instance of TagFileCreator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTagFileCreator(t interface {
	mock.TestingT
	Cleanup(func())
}) *TagFileCreator {
	mock := &TagFileCreator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
