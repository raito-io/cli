package wrappers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/raito-io/cli/base/data_source"
	ds_mocks "github.com/raito-io/cli/base/data_source/mocks"
	config2 "github.com/raito-io/cli/base/util/config"
	error2 "github.com/raito-io/cli/base/util/error"
)

func TestDataSourceSyncFunction_SyncDataSource(t *testing.T) {
	//Given
	config := &data_source.DataSourceSyncConfig{
		TargetFile:   "targetFile",
		DataSourceId: "DataSourceId",
		ConfigMap:    &config2.ConfigMap{Parameters: map[string]string{"key": "value"}},
	}

	fileCreatorMock := ds_mocks.NewDataSourceFileCreator(t)
	fileCreatorMock.EXPECT().Close().Return().Once()
	fileCreatorMock.EXPECT().GetDataObjectCount().Return(0)

	syncerMock := NewMockDataSourceSyncer(t)
	syncerMock.EXPECT().SyncDataSource(mock.Anything, fileCreatorMock, config.ConfigMap).Return(nil).Once()

	syncFunction := dataSourceSyncFunction{
		syncer: syncerMock,
		fileCreatorFactory: func(config *data_source.DataSourceSyncConfig) (data_source.DataSourceFileCreator, error) {
			return fileCreatorMock, nil
		},
	}

	//When
	result, err := syncFunction.SyncDataSource(context.Background(), config)

	//Then
	assert.NoError(t, err)
	assert.Nil(t, result.Error)
}

func TestDataSourceSyncFunction_SyncDataSource_ErrorOnFile(t *testing.T) {
	//Given
	config := &data_source.DataSourceSyncConfig{
		TargetFile:   "targetFile",
		DataSourceId: "DataSourceId",
		ConfigMap:    &config2.ConfigMap{Parameters: map[string]string{"key": "value"}},
	}

	syncerMock := NewMockDataSourceSyncer(t)

	syncFunction := dataSourceSyncFunction{
		syncer: syncerMock,
		fileCreatorFactory: func(config *data_source.DataSourceSyncConfig) (data_source.DataSourceFileCreator, error) {
			return nil, error2.ErrorResult{
				ErrorCode:    error2.ErrorCode_BAD_INPUT_PARAMETER_ERROR,
				ErrorMessage: "BOOM!",
			}
		},
	}

	//When
	result, err := syncFunction.SyncDataSource(context.Background(), config)

	//Then
	assert.NoError(t, err)
	assert.NotNil(t, result.Error)
	assert.Equal(t, "BOOM!", result.Error.ErrorMessage)
	assert.Equal(t, error2.ErrorCode_BAD_INPUT_PARAMETER_ERROR, result.Error.ErrorCode)

	syncerMock.AssertNotCalled(t, "SyncIdentityStore", mock.Anything, mock.Anything, mock.Anything)
}

func TestDataSourceSyncFunction_SyncDataSource_ErrorSync(t *testing.T) {
	//Given
	config := &data_source.DataSourceSyncConfig{
		TargetFile:   "targetFile",
		DataSourceId: "DataSourceId",
		ConfigMap:    &config2.ConfigMap{Parameters: map[string]string{"key": "value"}},
	}

	fileCreatorMock := ds_mocks.NewDataSourceFileCreator(t)
	fileCreatorMock.EXPECT().Close().Return().Once()

	syncerMock := NewMockDataSourceSyncer(t)
	syncerMock.EXPECT().SyncDataSource(mock.Anything, fileCreatorMock, config.ConfigMap).Return(error2.ErrorResult{
		ErrorCode:    error2.ErrorCode_SOURCE_CONNECTION_ERROR,
		ErrorMessage: "BOOM!",
	}).Once()

	syncFunction := dataSourceSyncFunction{
		syncer: syncerMock,
		fileCreatorFactory: func(config *data_source.DataSourceSyncConfig) (data_source.DataSourceFileCreator, error) {
			return fileCreatorMock, nil
		},
	}

	//When
	result, err := syncFunction.SyncDataSource(context.Background(), config)

	//Then
	assert.NoError(t, err)
	assert.NotNil(t, result.Error)
	assert.Equal(t, "BOOM!", result.Error.ErrorMessage)
	assert.Equal(t, error2.ErrorCode_SOURCE_CONNECTION_ERROR, result.Error.ErrorCode)
}

func TestDataSourceSyncFunctionWrapper(t *testing.T) {
	//Given
	syncerMock := NewMockDataSourceSyncer(t)

	//When
	syncFunction := DataSourceSync(syncerMock)

	//Then
	assert.Equal(t, syncerMock, syncFunction.syncer)
}
