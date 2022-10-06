package mocks

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/raito-io/cli/base/data_source"
	"github.com/raito-io/cli/base/data_usage"
)

func TestNewSimpleDataUsageStatementHandler(t *testing.T) {
	//Given
	statements1 := []data_usage.Statement{{User: "user1", Credits: 123456}, {User: "user2", Credits: 234567}}
	statements2 := []data_usage.Statement{{User: "user3", Credits: 3141592}}

	mock := NewSimpleDataUsageStatementHandler(t)
	err := mock.AddStatements(statements1)

	assert.NoError(t, err)
	assert.Len(t, mock.Statements, 2)
	assert.Equal(t, statements1, mock.Statements)

	err = mock.AddStatements(statements2)

	assert.NoError(t, err)
	assert.Len(t, mock.Statements, 3)
	assert.Equal(t, statements1, mock.Statements[0:2])
	assert.Equal(t, statements2, mock.Statements[2:])
}

func TestNewSimpleDataSourceObjectHandler_NoCalls(t *testing.T) {
	//Given
	mock := NewSimpleDataSourceObjectHandler(t, 2)

	//Then
	mock.AssertNotCalled(t, "AddDataObjects")
	mock.AssertNotCalled(t, "SetDataSourceName")
	mock.AssertNotCalled(t, "SetDataSourceFullname")
	mock.AssertNotCalled(t, "SetDataSourceDescription")
}

func TestNewSimpleDataSourceObjectHandler(t *testing.T) {
	//Given
	dataObjects := []data_source.DataObject{
		{
			Name:     "ObjectName1",
			Type:     "Table",
			FullName: "ObjectFullName1",
		},
		{
			Name:     "ObjectName2",
			Type:     "Schema",
			FullName: "ObjectFullName2",
		},
		{
			Name:     "ObjectName3",
			Type:     "Table",
			FullName: "ObjectFullName3",
		},
	}

	dataObjectPtrs := make([]*data_source.DataObject, 0, len(dataObjects))

	for i := range dataObjects {
		dataObjectPtrs = append(dataObjectPtrs, &dataObjects[i])
	}

	mock := NewSimpleDataSourceObjectHandler(t, 2)

	//When
	err := mock.AddDataObjects(dataObjectPtrs[0])

	//Then
	assert.NoError(t, err)
	assert.Len(t, mock.DataObjects, 1)
	assert.Equal(t, dataObjects[0], mock.DataObjects[0])

	//When
	err = mock.AddDataObjects(dataObjectPtrs[1:]...)

	//Then
	assert.NoError(t, err)
	assert.Len(t, mock.DataObjects, 3)
	assert.Equal(t, dataObjects, mock.DataObjects)

	//When
	mock.SetDataSourceFullname("DS FullName")
	mock.SetDataSourceName("DS Name")
	mock.SetDataSourceDescription("DS Descr")

	//Then
	assert.Equal(t, "DS FullName", mock.DataSourceFullName)
	assert.Equal(t, "DS Name", mock.DataSourceName)
	assert.Equal(t, "DS Descr", mock.DataSourceDescription)
}
