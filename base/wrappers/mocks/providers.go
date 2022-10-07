package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/raito-io/cli/base/data_source"
	"github.com/raito-io/cli/base/data_usage"
	"github.com/raito-io/cli/base/identity_store"
)

type SimpleDataUsageStatementHandler struct {
	*DataUsageStatementHandler
	Statements []data_usage.Statement
}

func NewSimpleDataUsageStatementHandler(t mockConstructorTestingTNewDataUsageStatementHandler) *SimpleDataUsageStatementHandler {
	result := &SimpleDataUsageStatementHandler{
		DataUsageStatementHandler: NewDataUsageStatementHandler(t),
		Statements:                make([]data_usage.Statement, 0),
	}

	result.EXPECT().AddStatements(mock.AnythingOfType("[]data_usage.Statement")).Run(func(statements []data_usage.Statement) {
		result.Statements = append(result.Statements, statements...)
	}).Return(nil)

	return result
}

type SimpleIdentityStoreIdentityHandler struct {
	*IdentityStoreIdentityHandler
	Users  []identity_store.User
	Groups []identity_store.Group
}

func NewSimpleIdentityStoreIdentityHandler(t mockConstructorTestingTNewIdentityStoreIdentityHandler, maxUsersOrGroupsInCall int) *SimpleIdentityStoreIdentityHandler {
	result := &SimpleIdentityStoreIdentityHandler{
		IdentityStoreIdentityHandler: NewIdentityStoreIdentityHandler(t),
		Users:                        make([]identity_store.User, 0),
		Groups:                       make([]identity_store.Group, 0),
	}

	addUsers := func(users ...*identity_store.User) {
		for i := range users {
			result.Users = append(result.Users, *users[i])
		}
	}

	addGroups := func(groups ...*identity_store.Group) {
		for i := range groups {
			result.Groups = append(result.Groups, *groups[i])
		}
	}

	arguments := make([]interface{}, 0)
	for i := 0; i < maxUsersOrGroupsInCall; i++ {
		arguments = append(arguments, mock.Anything)
		result.EXPECT().AddUsers(arguments...).Run(addUsers).Return(nil).Maybe()
		result.EXPECT().AddGroups(arguments...).Run(addGroups).Return(nil).Maybe()
	}

	return result
}

type SimpleDataSourceObjectHandler struct {
	*DataSourceObjectHandler
	DataObjects           []data_source.DataObject
	DataSourceName        string
	DataSourceFullName    string
	DataSourceDescription string
}

func NewSimpleDataSourceObjectHandler(t mockConstructorTestingTNewDataSourceObjectHandler, maxDataObjectsPerCall int) *SimpleDataSourceObjectHandler {
	result := &SimpleDataSourceObjectHandler{
		DataSourceObjectHandler: NewDataSourceObjectHandler(t),
		DataObjects:             make([]data_source.DataObject, 0),
	}

	arguments := make([]interface{}, 0)

	addDataObject := func(dataObjects ...*data_source.DataObject) {
		for _, do := range dataObjects {
			result.DataObjects = append(result.DataObjects, *do)
		}
	}

	for i := 0; i < maxDataObjectsPerCall; i++ {
		arguments = append(arguments, mock.Anything)
		result.EXPECT().AddDataObjects(arguments...).Run(addDataObject).Return(nil).Maybe()
	}

	result.EXPECT().SetDataSourceName(mock.AnythingOfType("string")).Run(func(name string) {
		result.DataSourceName = name
	}).Return().Maybe()

	result.EXPECT().SetDataSourceFullname(mock.AnythingOfType("string")).Run(func(name string) {
		result.DataSourceFullName = name
	}).Return().Maybe()

	result.EXPECT().SetDataSourceDescription(mock.AnythingOfType("string")).Run(func(desc string) {
		result.DataSourceDescription = desc
	}).Return().Maybe()

	return result
}
