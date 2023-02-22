package wrappers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/raito-io/cli/base/identity_store"
	is_mocks "github.com/raito-io/cli/base/identity_store/mocks"
	config2 "github.com/raito-io/cli/base/util/config"
	error2 "github.com/raito-io/cli/base/util/error"
)

func TestIdentityStoreSyncFunction_SyncIdentityStore(t *testing.T) {
	//Given
	config := &identity_store.IdentityStoreSyncConfig{
		GroupFile: "GroupFile",
		UserFile:  "UserFile",
		ConfigMap: &config2.ConfigMap{Parameters: map[string]string{"key": "value"}},
	}

	fileCreatorMock := is_mocks.NewIdentityStoreFileCreator(t)
	fileCreatorMock.EXPECT().Close().Return().Once()
	fileCreatorMock.EXPECT().GetUserCount().Return(0)
	fileCreatorMock.EXPECT().GetGroupCount().Return(0)

	syncerMock := NewMockIdentityStoreSyncer(t)
	syncerMock.EXPECT().SyncIdentityStore(mock.Anything, fileCreatorMock, config.ConfigMap).Return(nil).Once()

	syncFunction := identityStoreSyncFunction{
		syncer: syncerMock,
		identityHandlerFactory: func(config *identity_store.IdentityStoreSyncConfig) (identity_store.IdentityStoreFileCreator, error) {
			return fileCreatorMock, nil
		},
	}

	//When
	result, err := syncFunction.SyncIdentityStore(context.Background(), config)

	//Then
	assert.NoError(t, err)
	assert.Nil(t, result.Error)
}

func TestDataUsageSyncFunction_SyncDataUsage_ErrorOfFileCreation(t *testing.T) {
	//Given
	config := &identity_store.IdentityStoreSyncConfig{
		GroupFile: "GroupFile",
		UserFile:  "UserFile",
		ConfigMap: &config2.ConfigMap{Parameters: map[string]string{"key": "value"}},
	}

	syncerMock := NewMockIdentityStoreSyncer(t)

	syncFunction := identityStoreSyncFunction{
		syncer: syncerMock,
		identityHandlerFactory: func(config *identity_store.IdentityStoreSyncConfig) (identity_store.IdentityStoreFileCreator, error) {
			return nil, &error2.ErrorResult{
				ErrorMessage: "BOOM!",
				ErrorCode:    error2.ErrorCode_BAD_INPUT_PARAMETER_ERROR,
			}
		},
	}

	//When
	result, err := syncFunction.SyncIdentityStore(context.Background(), config)

	//Then
	assert.NoError(t, err)
	assert.NotNil(t, result.Error)
	assert.Equal(t, "BOOM!", result.Error.ErrorMessage)
	assert.Equal(t, error2.ErrorCode_BAD_INPUT_PARAMETER_ERROR, result.Error.ErrorCode)

	syncerMock.AssertNotCalled(t, "SyncIdentityStore", mock.Anything, mock.Anything, mock.Anything)
}

func TestMockDataUsageSyncer_SyncDataUsage_ErrorSync(t *testing.T) {
	//Given
	config := &identity_store.IdentityStoreSyncConfig{
		GroupFile: "GroupFile",
		UserFile:  "UserFile",
		ConfigMap: &config2.ConfigMap{Parameters: map[string]string{"key": "value"}},
	}

	fileCreatorMock := is_mocks.NewIdentityStoreFileCreator(t)
	fileCreatorMock.EXPECT().Close().Return().Once()

	syncerMock := NewMockIdentityStoreSyncer(t)
	syncerMock.EXPECT().SyncIdentityStore(mock.Anything, fileCreatorMock, config.ConfigMap).Return(&error2.ErrorResult{
		ErrorMessage: "BOOM!",
		ErrorCode:    error2.ErrorCode_SOURCE_CONNECTION_ERROR,
	}).Once()

	syncFunction := identityStoreSyncFunction{
		syncer: syncerMock,
		identityHandlerFactory: func(config *identity_store.IdentityStoreSyncConfig) (identity_store.IdentityStoreFileCreator, error) {
			return fileCreatorMock, nil
		},
	}

	//When
	result, err := syncFunction.SyncIdentityStore(context.Background(), config)

	//Then
	assert.NoError(t, err)
	assert.NotNil(t, result.Error)
	assert.Equal(t, "BOOM!", result.Error.ErrorMessage)
	assert.Equal(t, error2.ErrorCode_SOURCE_CONNECTION_ERROR, result.Error.ErrorCode)
}

func TestIdentityStoreSyncWrapper(t *testing.T) {
	//Given
	syncerMock := NewMockIdentityStoreSyncer(t)

	//When
	syncFunction := IdentityStoreSync(syncerMock)

	//Then
	assert.Equal(t, syncerMock, syncFunction.syncer)
}
