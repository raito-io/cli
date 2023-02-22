package wrappers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/raito-io/cli/base/data_usage"
	du_mocks "github.com/raito-io/cli/base/data_usage/mocks"
	config2 "github.com/raito-io/cli/base/util/config"
	error2 "github.com/raito-io/cli/base/util/error"
)

func TestDataUsageSyncFunction_SyncDataUsage(t *testing.T) {
	//Given
	config := &data_usage.DataUsageSyncConfig{
		TargetFile: "SomeTargetString",
		ConfigMap:  config2.ConfigMap{Parameters: map[string]string{"key": "value"}},
	}

	fileCreatorMock := du_mocks.NewDataUsageFileCreator(t)
	fileCreatorMock.EXPECT().Close().Return()
	fileCreatorMock.EXPECT().GetStatementCount().Return(0)

	syncerMock := NewMockDataUsageSyncer(t)
	syncerMock.EXPECT().SyncDataUsage(mock.Anything, fileCreatorMock, &config.ConfigMap).Return(nil)

	syncFunction := dataUsageSyncFunction{
		syncer: syncerMock,
		fileCreatorFactory: func(config *data_usage.DataUsageSyncConfig) (data_usage.DataUsageFileCreator, error) {
			return fileCreatorMock, nil
		},
	}

	//When
	result := syncFunction.SyncDataUsage(config)

	//Then
	assert.Nil(t, result.Error)
	syncerMock.AssertNumberOfCalls(t, "SyncDataUsage", 1)
	fileCreatorMock.AssertNumberOfCalls(t, "Close", 1)
}

func TestDataUsageSyncFunction_SyncDataUsage_ErrorOnFileCreation(t *testing.T) {
	//Given
	config := &data_usage.DataUsageSyncConfig{
		TargetFile: "SomeTargetString",
		ConfigMap:  config2.ConfigMap{Parameters: map[string]string{"key": "value"}},
	}

	syncerMock := NewMockDataUsageSyncer(t)

	syncFunction := dataUsageSyncFunction{
		syncer: syncerMock,
		fileCreatorFactory: func(config *data_usage.DataUsageSyncConfig) (data_usage.DataUsageFileCreator, error) {
			return nil, &error2.ErrorResult{
				ErrorMessage: "BOOM!",
				ErrorCode:    error2.ErrorCode_UNKNOWN_ERROR,
			}
		},
	}

	//When
	result := syncFunction.SyncDataUsage(config)

	//Then
	assert.NotNil(t, result.Error)
	assert.Equal(t, "BOOM!", result.Error.ErrorMessage)
	assert.Equal(t, error2.ErrorCode_UNKNOWN_ERROR, result.Error.ErrorCode)

	syncerMock.AssertNotCalled(t, "SyncDataUsage", mock.Anything, mock.Anything, mock.Anything)
}

func TestDataUsageSyncFunction_SyncDataUsage_ErrorSync(t *testing.T) {
	//Given
	config := &data_usage.DataUsageSyncConfig{
		TargetFile: "SomeTargetString",
		ConfigMap:  config2.ConfigMap{Parameters: map[string]string{"key": "value"}},
	}

	fileCreatorMock := du_mocks.NewDataUsageFileCreator(t)
	fileCreatorMock.EXPECT().Close().Return()

	syncerMock := NewMockDataUsageSyncer(t)
	syncerMock.EXPECT().SyncDataUsage(mock.Anything, fileCreatorMock, &config.ConfigMap).Return(&error2.ErrorResult{
		ErrorMessage: "BOOM!",
		ErrorCode:    error2.ErrorCode_BAD_INPUT_PARAMETER_ERROR,
	})

	syncFunction := dataUsageSyncFunction{
		syncer: syncerMock,
		fileCreatorFactory: func(config *data_usage.DataUsageSyncConfig) (data_usage.DataUsageFileCreator, error) {
			return fileCreatorMock, nil
		},
	}

	//When
	result := syncFunction.SyncDataUsage(config)

	//Then
	assert.NotNil(t, result.Error)
	assert.Equal(t, "BOOM!", result.Error.ErrorMessage)
	assert.Equal(t, error2.ErrorCode_BAD_INPUT_PARAMETER_ERROR, result.Error.ErrorCode)
	syncerMock.AssertNumberOfCalls(t, "SyncDataUsage", 1)
	fileCreatorMock.AssertNumberOfCalls(t, "Close", 1)
}

func TestDataUsageSyncWrapper(t *testing.T) {
	//Given
	syncerMock := NewMockDataUsageSyncer(t)

	//When
	syncFunction := DataUsageSync(syncerMock)

	//Then
	assert.Equal(t, syncerMock, syncFunction.syncer)
}
