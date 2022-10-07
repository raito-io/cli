package wrappers

import (
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
		ConfigMap:    config2.ConfigMap{Parameters: map[string]interface{}{"key": "value"}},
	}

	fileCreatorMock := ds_mocks.NewDataSourceFileCreator(t)
	fileCreatorMock.EXPECT().Close().Return().Once()
	fileCreatorMock.EXPECT().GetDataObjectCount().Return(0)

	syncerMock := NewMockDataSourceSyncer(t)
	syncerMock.EXPECT().SyncDataSource(mock.Anything, fileCreatorMock, &config.ConfigMap).Return(nil).Once()

	syncFunction := dataSourceSyncFunction{
		syncer: syncerMock,
		fileCreatorFactory: func(config *data_source.DataSourceSyncConfig) (data_source.DataSourceFileCreator, error) {
			return fileCreatorMock, nil
		},
	}

	//When
	result := syncFunction.SyncDataSource(config)

	//Then
	assert.Nil(t, result.Error)
}

func TestDataSourceSyncFunction_SyncDataSource_ErrorOnFile(t *testing.T) {
	//Given
	config := &data_source.DataSourceSyncConfig{
		TargetFile:   "targetFile",
		DataSourceId: "DataSourceId",
		ConfigMap:    config2.ConfigMap{Parameters: map[string]interface{}{"key": "value"}},
	}

	syncerMock := NewMockDataSourceSyncer(t)

	syncFunction := dataSourceSyncFunction{
		syncer: syncerMock,
		fileCreatorFactory: func(config *data_source.DataSourceSyncConfig) (data_source.DataSourceFileCreator, error) {
			return nil, error2.ErrorResult{
				ErrorCode:    error2.BadInputParameterError,
				ErrorMessage: "BOOM!",
			}
		},
	}

	//When
	result := syncFunction.SyncDataSource(config)

	//Then
	assert.NotNil(t, result.Error)
	assert.Equal(t, "BOOM!", result.Error.ErrorMessage)
	assert.Equal(t, error2.BadInputParameterError, result.Error.ErrorCode)

	syncerMock.AssertNotCalled(t, "SyncIdentityStore", mock.Anything, mock.Anything, mock.Anything)
}

func TestDataSourceSyncFunction_SyncDataSource_ErrorSync(t *testing.T) {
	//Given
	config := &data_source.DataSourceSyncConfig{
		TargetFile:   "targetFile",
		DataSourceId: "DataSourceId",
		ConfigMap:    config2.ConfigMap{Parameters: map[string]interface{}{"key": "value"}},
	}

	fileCreatorMock := ds_mocks.NewDataSourceFileCreator(t)
	fileCreatorMock.EXPECT().Close().Return().Once()

	syncerMock := NewMockDataSourceSyncer(t)
	syncerMock.EXPECT().SyncDataSource(mock.Anything, fileCreatorMock, &config.ConfigMap).Return(error2.ErrorResult{
		ErrorCode:    error2.SourceConnectionError,
		ErrorMessage: "BOOM!",
	}).Once()

	syncFunction := dataSourceSyncFunction{
		syncer: syncerMock,
		fileCreatorFactory: func(config *data_source.DataSourceSyncConfig) (data_source.DataSourceFileCreator, error) {
			return fileCreatorMock, nil
		},
	}

	//When
	result := syncFunction.SyncDataSource(config)

	//Then
	assert.NotNil(t, result.Error)
	assert.Equal(t, "BOOM!", result.Error.ErrorMessage)
	assert.Equal(t, error2.SourceConnectionError, result.Error.ErrorCode)
}

func TestDataSourceSyncFunctionWrapper(t *testing.T) {
	//Given
	syncerMock := NewMockDataSourceSyncer(t)

	//When
	syncFunction := DataSourceSync(syncerMock)

	//Then
	assert.Equal(t, syncerMock, syncFunction.syncer)
}
