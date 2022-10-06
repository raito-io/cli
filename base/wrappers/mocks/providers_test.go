package mocks

import (
	"testing"

	"github.com/stretchr/testify/assert"

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
