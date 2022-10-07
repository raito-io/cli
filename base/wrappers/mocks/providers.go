package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/raito-io/cli/base/access_provider/sync_from_target"
	"github.com/raito-io/cli/base/access_provider/sync_to_target"
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

type SimpleAccessProviderHandler struct {
	*AccessProviderHandler
	AccessProviders []sync_from_target.AccessProvider
}

func NewSimpleAccessProviderHandler(t mockConstructorTestingTNewAccessProviderHandler, maxAccessProvidersPerCall int) *SimpleAccessProviderHandler {
	result := &SimpleAccessProviderHandler{
		AccessProviderHandler: NewAccessProviderHandler(t),
		AccessProviders:       make([]sync_from_target.AccessProvider, 0),
	}

	arguments := make([]interface{}, 0)

	for i := 0; i < maxAccessProvidersPerCall; i++ {
		arguments = append(arguments, mock.Anything)
		result.EXPECT().AddAccessProviders(arguments...).Run(func(dataAccessList ...*sync_from_target.AccessProvider) {
			for _, da := range dataAccessList {
				result.AccessProviders = append(result.AccessProviders, *da)
			}
		}).Return(nil).Maybe()
	}

	return result
}

type SimpleAccessProviderFeedbackHandler struct {
	*AccessProviderFeedbackHandler
	AccessProviderFeedback map[string][]sync_to_target.AccessSyncFeedbackInformation
}

func NewSimpleAccessProviderFeedbackHandler(t mockConstructorTestingTNewAccessProviderFeedbackHandler, maxAccessFeedbackInformationObjectsPerCall int) *SimpleAccessProviderFeedbackHandler {
	result := &SimpleAccessProviderFeedbackHandler{
		AccessProviderFeedbackHandler: NewAccessProviderFeedbackHandler(t),
		AccessProviderFeedback:        map[string][]sync_to_target.AccessSyncFeedbackInformation{},
	}

	arguments := make([]interface{}, 0)

	for i := 0; i < maxAccessFeedbackInformationObjectsPerCall; i++ {
		arguments = append(arguments, mock.Anything)

		result.EXPECT().AddAccessProviderFeedback(mock.Anything, arguments...).Run(func(accessProviderId string, accessFeedback ...sync_to_target.AccessSyncFeedbackInformation) {
			result.AccessProviderFeedback[accessProviderId] = accessFeedback
		}).Return(nil).Maybe()
	}

	return result
}
