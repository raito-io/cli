package wrappers

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/raito-io/cli/base/data_usage"
	du_mocks "github.com/raito-io/cli/base/data_usage/mocks"
	config2 "github.com/raito-io/cli/base/util/config"
)

func TestDataUsageSyncFunction_SyncDataUsage(t *testing.T) {
	//Given
	config := &data_usage.DataUsageSyncConfig{
		TargetFile: "SomeTargetString",
		ConfigMap:  &config2.ConfigMap{Parameters: map[string]string{"key": "value"}},
	}

	fileCreatorMock := du_mocks.NewDataUsageFileCreator(t)
	fileCreatorMock.EXPECT().Close().Return()
	fileCreatorMock.EXPECT().GetStatementCount().Return(0)
	fileCreatorMock.EXPECT().GetImportFileSize().Return(uint64(3))

	syncerMock := NewMockDataUsageSyncer(t)
	syncerMock.EXPECT().SyncDataUsage(mock.Anything, fileCreatorMock, config.ConfigMap).Return(nil)

	syncFunction := dataUsageSyncFunction{
		syncer: NewSyncFactory[config2.ConfigMap, DataUsageSyncer](NewDummySyncFactoryFn[config2.ConfigMap, DataUsageSyncer](syncerMock)),
		fileCreatorFactory: func(config *data_usage.DataUsageSyncConfig) (data_usage.DataUsageFileCreator, error) {
			return fileCreatorMock, nil
		},
	}

	//When
	result, err := syncFunction.SyncDataUsage(context.Background(), config)

	//Then
	assert.NoError(t, err)
	assert.Nil(t, result.Error)
	syncerMock.AssertNumberOfCalls(t, "SyncDataUsage", 1)
	fileCreatorMock.AssertNumberOfCalls(t, "Close", 1)
}

func TestDataUsageSyncFunction_SyncDataUsage_ErrorOnFileCreation(t *testing.T) {
	//Given
	config := &data_usage.DataUsageSyncConfig{
		TargetFile: "SomeTargetString",
		ConfigMap:  &config2.ConfigMap{Parameters: map[string]string{"key": "value"}},
	}

	syncerMock := NewMockDataUsageSyncer(t)

	syncFunction := dataUsageSyncFunction{
		syncer: NewSyncFactory[config2.ConfigMap, DataUsageSyncer](NewDummySyncFactoryFn[config2.ConfigMap, DataUsageSyncer](syncerMock)),
		fileCreatorFactory: func(config *data_usage.DataUsageSyncConfig) (data_usage.DataUsageFileCreator, error) {
			return nil, errors.New("BOOM!")
		},
	}

	//When
	result, err := syncFunction.SyncDataUsage(context.Background(), config)

	//Then
	assert.Error(t, err)
	assert.Nil(t, result)

	syncerMock.AssertNotCalled(t, "SyncDataUsage", mock.Anything, mock.Anything, mock.Anything)
}

func TestDataUsageSyncFunction_SyncDataUsage_ErrorSync(t *testing.T) {
	//Given
	config := &data_usage.DataUsageSyncConfig{
		TargetFile: "SomeTargetString",
		ConfigMap:  &config2.ConfigMap{Parameters: map[string]string{"key": "value"}},
	}

	fileCreatorMock := du_mocks.NewDataUsageFileCreator(t)
	fileCreatorMock.EXPECT().Close().Return()

	syncerMock := NewMockDataUsageSyncer(t)
	syncerMock.EXPECT().SyncDataUsage(mock.Anything, fileCreatorMock, config.ConfigMap).Return(errors.New("BOOM!"))

	syncFunction := dataUsageSyncFunction{
		syncer: NewSyncFactory[config2.ConfigMap, DataUsageSyncer](NewDummySyncFactoryFn[config2.ConfigMap, DataUsageSyncer](syncerMock)),
		fileCreatorFactory: func(config *data_usage.DataUsageSyncConfig) (data_usage.DataUsageFileCreator, error) {
			return fileCreatorMock, nil
		},
	}

	//When
	result, err := syncFunction.SyncDataUsage(context.Background(), config)

	//Then
	assert.Error(t, err)
	assert.Nil(t, result)
	syncerMock.AssertNumberOfCalls(t, "SyncDataUsage", 1)
	fileCreatorMock.AssertNumberOfCalls(t, "Close", 1)
}

func TestDataUsageSyncWrapper(t *testing.T) {
	//Given
	syncerMock := NewMockDataUsageSyncer(t)

	//When
	syncFunction, err := DataUsageSync(syncerMock).syncer.Create(context.Background(), &config2.ConfigMap{})

	//Then
	require.NoError(t, err)
	assert.Equal(t, syncerMock, syncFunction)
}
