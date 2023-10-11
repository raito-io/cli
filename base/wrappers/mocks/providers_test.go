package mocks

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/raito-io/cli/base/access_provider/sync_from_target"
	"github.com/raito-io/cli/base/access_provider/sync_to_target"
	"github.com/raito-io/cli/base/data_source"
	"github.com/raito-io/cli/base/data_usage"
	"github.com/raito-io/cli/base/identity_store"
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

func TestNewSimpleIdentityStoreIdentityHandler_NoCalls(t *testing.T) {
	//Given
	mock := NewSimpleIdentityStoreIdentityHandler(t, 2)

	//Then
	mock.AssertNotCalled(t, "AddUsers")
	mock.AssertNotCalled(t, "AddGroups")
}

func TestNewSimpleIdentityStoreIdentityHandler(t *testing.T) {
	//Given
	mock := NewSimpleIdentityStoreIdentityHandler(t, 2)

	groups := []identity_store.Group{{Name: "GroupName1", DisplayName: "Group1"}, {Name: "GroupName2", DisplayName: "Group2"}}
	users := []identity_store.User{{Name: "User1", UserName: "user1"}, {UserName: "user2", Name: "User2"}}

	groupsPtr := make([]*identity_store.Group, len(groups))
	for i := range groups {
		groupsPtr[i] = &groups[i]
	}

	usersPtr := make([]*identity_store.User, len(users))
	for i := range users {
		usersPtr[i] = &users[i]
	}

	//When
	err := mock.AddUsers(usersPtr[0])

	//Then
	assert.NoError(t, err)
	assert.Equal(t, users[0], mock.Users[0])

	//When
	err = mock.AddGroups(groupsPtr...)

	//Then
	assert.NoError(t, err)
	assert.Equal(t, groups, mock.Groups)

	//When
	err = mock.AddUsers(usersPtr[1:]...)

	//Then
	assert.NoError(t, err)
	assert.Equal(t, users, mock.Users)

	mock.AssertNumberOfCalls(t, "AddGroups", 1)
	mock.AssertNumberOfCalls(t, "AddUsers", 2)
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

func TestNewSimpleAccessProviderHandler(t *testing.T) {
	//Given
	mock := NewSimpleAccessProviderHandler(t, 2)

	accessProviders := []sync_from_target.AccessProvider{
		{Name: "AP1"}, {Name: "AP2"}, {Name: "AP3"},
	}

	accessProviderPtrs := make([]*sync_from_target.AccessProvider, len(accessProviders))

	for i := range accessProviders {
		accessProviderPtrs[i] = &accessProviders[i]
	}

	//When
	err := mock.AddAccessProviders(accessProviderPtrs[0])

	//Then
	assert.NoError(t, err)
	assert.Len(t, mock.AccessProviders, 1)
	assert.Equal(t, accessProviders[0], mock.AccessProviders[0])

	//When
	err = mock.AddAccessProviders(accessProviderPtrs[1:]...)

	//Then
	assert.NoError(t, err)
	assert.Len(t, mock.AccessProviders, 3)
	assert.Equal(t, accessProviders, mock.AccessProviders)
}

func TestNewSimpleAccessProviderFeedbackHandler(t *testing.T) {
	//Given
	mock := NewSimpleAccessProviderFeedbackHandler(t)

	accessProviderFeedbackMap := []sync_to_target.AccessProviderSyncFeedback{
		{
			AccessProvider: "AccessId1",
			ActualName:     "ActualName1",
		},
		{
			AccessProvider: "AccessId2",
			ActualName:     "ActualName2",
		},
	}

	//When
	err := mock.AddAccessProviderFeedback(accessProviderFeedbackMap[0])

	//Then
	assert.NoError(t, err)

	//When
	err = mock.AddAccessProviderFeedback(accessProviderFeedbackMap[1])

	//Then
	assert.NoError(t, err)
	assert.Equal(t, accessProviderFeedbackMap, mock.AccessProviderFeedback)

}
