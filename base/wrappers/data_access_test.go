package wrappers

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/raito-io/cli/base/access_provider"
	"github.com/raito-io/cli/base/access_provider/sync_from_target"
	"github.com/raito-io/cli/base/access_provider/sync_from_target/mocks"
	"github.com/raito-io/cli/base/access_provider/sync_to_target"
	mocks2 "github.com/raito-io/cli/base/access_provider/sync_to_target/mocks"
	config2 "github.com/raito-io/cli/base/util/config"
	error2 "github.com/raito-io/cli/base/util/error"
)

func TestDataAccessSyncFunction_SyncFromTarget(t *testing.T) {
	//Given
	config := &access_provider.AccessSyncFromTarget{
		Prefix:     "prefix",
		TargetFile: "targetFile",
		ConfigMap:  &config2.ConfigMap{Parameters: map[string]string{"key": "value"}},
	}

	accessProviderFileCreator := mocks.NewAccessProviderFileCreator(t)
	accessProviderFileCreator.EXPECT().Close().Return().Once()
	accessProviderFileCreator.EXPECT().GetAccessProviderCount().Return(1).Twice()

	syncerMock := NewMockAccessProviderSyncer(t)
	syncerMock.EXPECT().SyncAccessProvidersFromTarget(mock.Anything, accessProviderFileCreator, config.ConfigMap).Return(nil).Once()

	syncFunction := DataAccessSyncFunction{
		Syncer: syncerMock,
		accessFileCreatorFactory: func(config *access_provider.AccessSyncFromTarget) (sync_from_target.AccessProviderFileCreator, error) {
			return accessProviderFileCreator, nil
		},
	}

	//When
	result, err := syncFunction.SyncFromTarget(context.Background(), config)

	//Then
	assert.NoError(t, err)
	assert.Equal(t, int32(1), result.AccessProviderCount)
	assert.Nil(t, result.Error)
}

func TestDataAccessSyncFunction_SyncFromTarget_ErrorOnFileCreation(t *testing.T) {
	//Given
	config := &access_provider.AccessSyncFromTarget{
		Prefix:     "prefix",
		TargetFile: "targetFile",
		ConfigMap:  &config2.ConfigMap{Parameters: map[string]string{"key": "value"}},
	}

	syncerMock := NewMockAccessProviderSyncer(t)

	syncFunction := DataAccessSyncFunction{
		Syncer: syncerMock,
		accessFileCreatorFactory: func(config *access_provider.AccessSyncFromTarget) (sync_from_target.AccessProviderFileCreator, error) {
			return nil, errors.New("BOOM!")
		},
	}

	//When
	result, err := syncFunction.SyncFromTarget(context.Background(), config)

	//Then
	assert.Error(t, err)
	assert.Nil(t, result)

	syncerMock.AssertNotCalled(t, "SyncAccessProvidersFromTarget", mock.Anything, mock.Anything, mock.Anything)
}

func TestDataAccessSyncFunction_SyncFromTarget_ErrorOnSync(t *testing.T) {
	//Given
	config := &access_provider.AccessSyncFromTarget{
		Prefix:     "prefix",
		TargetFile: "targetFile",
		ConfigMap:  &config2.ConfigMap{Parameters: map[string]string{"key": "value"}},
	}

	accessProviderFilCreator := mocks.NewAccessProviderFileCreator(t)
	accessProviderFilCreator.EXPECT().Close().Return().Once()

	syncerMock := NewMockAccessProviderSyncer(t)
	syncerMock.EXPECT().SyncAccessProvidersFromTarget(mock.Anything, accessProviderFilCreator, config.ConfigMap).Return(
		errors.New("BOOM!")).Once()

	syncFunction := DataAccessSyncFunction{
		Syncer: syncerMock,
		accessFileCreatorFactory: func(config *access_provider.AccessSyncFromTarget) (sync_from_target.AccessProviderFileCreator, error) {
			return accessProviderFilCreator, nil
		},
	}

	//When
	result, err := syncFunction.SyncFromTarget(context.Background(), config)

	//Then
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestDataAccessSyncFunction_SyncToTarget_AccessProviders(t *testing.T) {
	//Given
	config := &access_provider.AccessSyncToTarget{
		SourceFile:         "SourceFile",
		FeedbackTargetFile: "FeedbackTargetFile",
		ConfigMap:          &config2.ConfigMap{Parameters: map[string]string{"key": "value"}},
	}

	accessFeedBackFileCreator := mocks2.NewSyncFeedbackFileCreator(t)
	accessFeedBackFileCreator.EXPECT().Close().Once()

	actualName1 := "ActualName1"

	accessProvidersImport := sync_to_target.AccessProviderImport{
		LastCalculated: time.Now().Unix(),
		AccessProviders: []*sync_to_target.AccessProvider{
			{
				Id:          "AP1",
				Description: "Descr",
				Delete:      false,
				Name:        "Ap1",
				NamingHint:  "NameHint1",
				Action:      sync_to_target.Grant,
			},
			{
				Id:          "AP2",
				Description: "Descr2",
				Delete:      false,
				Name:        "Ap2",
				NamingHint:  "NameHint2",
				Action:      sync_to_target.Grant,
			},
			{
				ActualName:  &actualName1,
				Id:          "AP3",
				Description: "Descr3",
				Delete:      true,
				Name:        "Ap3",
				NamingHint:  "NameHint3",
				Action:      sync_to_target.Grant,
			},
			{
				Id:          "AP4",
				Description: "Descr4",
				Delete:      true,
				Name:        "Ap4",
				NamingHint:  "NameHint4",
				Action:      sync_to_target.Grant,
			},
		},
	}

	accessProviderParser := mocks2.NewAccessProviderImportFileParser(t)
	accessProviderParser.EXPECT().ParseAccessProviders().Return(&accessProvidersImport, nil).Once()

	syncerMock := NewMockAccessProviderSyncer(t)
	syncerMock.EXPECT().SyncAccessProviderToTarget(mock.Anything, &accessProvidersImport, accessFeedBackFileCreator, config.ConfigMap).Return(nil).Once()

	syncFunction := DataAccessSyncFunction{
		Syncer: syncerMock,
		accessFeedbackFileCreatorFactory: func(config *access_provider.AccessSyncToTarget) (sync_to_target.SyncFeedbackFileCreator, error) {
			return accessFeedBackFileCreator, nil
		},
		accessProviderParserFactory: func(config *access_provider.AccessSyncToTarget) (sync_to_target.AccessProviderImportFileParser, error) {
			return accessProviderParser, nil
		},
	}

	//When
	result, err := syncFunction.SyncToTarget(context.Background(), config)

	//Then
	assert.NoError(t, err)
	assert.Nil(t, result.Error)
}

func TestDataAccessSyncFunction_SyncToTarget_AccessAsCode(t *testing.T) {
	//Given
	config := &access_provider.AccessSyncToTarget{
		Prefix:             "R",
		SourceFile:         "SourceFile",
		FeedbackTargetFile: "FeedbackTargetFile",
		ConfigMap:          &config2.ConfigMap{Parameters: map[string]string{"key": "value"}},
	}

	accessProvidersImport := sync_to_target.AccessProviderImport{
		LastCalculated: time.Now().Unix(),
		AccessProviders: []*sync_to_target.AccessProvider{
			{
				Id:          "AP1",
				Description: "Descr",
				Delete:      false,
				Name:        "Ap1",
				NamingHint:  "NameHint1",
				Action:      sync_to_target.Grant,
			},
			{
				Id:          "AP2",
				Description: "Descr2",
				Delete:      false,
				Name:        "Ap2",
				NamingHint:  "NameHint2",
				Action:      sync_to_target.Grant,
			},
		},
	}

	accessProviderParser := mocks2.NewAccessProviderImportFileParser(t)
	accessProviderParser.EXPECT().ParseAccessProviders().Return(&accessProvidersImport, nil).Once()

	syncerMock := NewMockAccessProviderSyncer(t)
	syncerMock.EXPECT().SyncAccessAsCodeToTarget(mock.Anything, &accessProvidersImport, "R", config.ConfigMap).Return(nil).Once()

	syncFunction := DataAccessSyncFunction{
		Syncer: syncerMock,
		accessProviderParserFactory: func(config *access_provider.AccessSyncToTarget) (sync_to_target.AccessProviderImportFileParser, error) {
			return accessProviderParser, nil
		},
	}

	//When
	result, err := syncFunction.SyncToTarget(context.Background(), config)

	//Then
	assert.NoError(t, err)
	assert.Nil(t, result.Error)
}

func TestDataAccessSyncFunction_SyncToTarget_ErrorOnFileParsingFactory(t *testing.T) {
	//Given
	config := &access_provider.AccessSyncToTarget{
		SourceFile:         "SourceFile",
		FeedbackTargetFile: "FeedbackTargetFile",
		ConfigMap:          &config2.ConfigMap{Parameters: map[string]string{"key": "value"}},
	}

	syncerMock := NewMockAccessProviderSyncer(t)

	syncFunction := DataAccessSyncFunction{
		Syncer: syncerMock,
		accessProviderParserFactory: func(config *access_provider.AccessSyncToTarget) (sync_to_target.AccessProviderImportFileParser, error) {
			return nil, error2.ToErrorResult(fmt.Errorf("boom"))
		},
	}

	//When
	result, err := syncFunction.SyncToTarget(context.Background(), config)

	//Then
	assert.Error(t, err)
	assert.Nil(t, result)

	syncerMock.AssertNotCalled(t, "SyncAccessAsCodeToTarget", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	syncerMock.AssertNotCalled(t, "SyncAccessProvidersToTarget", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestDataAccessSyncFunction_SyncToTarget_ErrorOnFileParsing(t *testing.T) {
	//Given
	config := &access_provider.AccessSyncToTarget{
		SourceFile:         "SourceFile",
		FeedbackTargetFile: "FeedbackTargetFile",
		ConfigMap:          &config2.ConfigMap{Parameters: map[string]string{"key": "value"}},
	}

	accessProviderParser := mocks2.NewAccessProviderImportFileParser(t)
	accessProviderParser.EXPECT().ParseAccessProviders().Return(nil, errors.New("BOOM!")).Once()

	syncerMock := NewMockAccessProviderSyncer(t)

	syncFunction := DataAccessSyncFunction{
		Syncer: syncerMock,
		accessProviderParserFactory: func(config *access_provider.AccessSyncToTarget) (sync_to_target.AccessProviderImportFileParser, error) {
			return accessProviderParser, nil
		},
	}

	//When
	result, err := syncFunction.SyncToTarget(context.Background(), config)

	//Then
	assert.Error(t, err)
	assert.Nil(t, result)

	syncerMock.AssertNotCalled(t, "SyncAccessProvidersToTarget", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	syncerMock.AssertNotCalled(t, "SyncAccessAsCodeToTarget", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestDataAccessSyncFunction_SyncToTarget_AccessProviders_ErrorOnFeedbackFileCreation(t *testing.T) {
	//Given
	config := &access_provider.AccessSyncToTarget{
		SourceFile:         "SourceFile",
		FeedbackTargetFile: "FeedbackTargetFile",
		ConfigMap:          &config2.ConfigMap{Parameters: map[string]string{"key": "value"}},
	}

	accessProvidersImport := sync_to_target.AccessProviderImport{
		LastCalculated: time.Now().Unix(),
		AccessProviders: []*sync_to_target.AccessProvider{
			{
				Id:          "AP1",
				Description: "Descr",
				Delete:      false,
				Name:        "Ap1",
				NamingHint:  "NameHint1",
				Action:      sync_to_target.Grant,
			},
			{
				Id:          "AP2",
				Description: "Descr2",
				Delete:      false,
				Name:        "Ap2",
				NamingHint:  "NameHint2",
				Action:      sync_to_target.Grant,
			},
		},
	}

	accessProviderParser := mocks2.NewAccessProviderImportFileParser(t)
	accessProviderParser.EXPECT().ParseAccessProviders().Return(&accessProvidersImport, nil).Once()

	syncerMock := NewMockAccessProviderSyncer(t)

	syncFunction := DataAccessSyncFunction{
		Syncer: syncerMock,
		accessFeedbackFileCreatorFactory: func(config *access_provider.AccessSyncToTarget) (sync_to_target.SyncFeedbackFileCreator, error) {
			return nil, fmt.Errorf("boom")
		},
		accessProviderParserFactory: func(config *access_provider.AccessSyncToTarget) (sync_to_target.AccessProviderImportFileParser, error) {
			return accessProviderParser, nil
		},
	}

	//When
	result, err := syncFunction.SyncToTarget(context.Background(), config)

	//Then
	assert.Error(t, err)
	assert.Nil(t, result)

	syncerMock.AssertNotCalled(t, "SyncAccessProvidersToTarget", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	syncerMock.AssertNotCalled(t, "SyncAccessAsCodeToTarget", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	syncerMock.AssertNotCalled(t, "SyncAccessProvidersToTarget", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestDataAccessSyncFunction_SyncToTarget_AccessProviders_ErrorOnSync(t *testing.T) {
	//Given
	config := &access_provider.AccessSyncToTarget{
		SourceFile:         "SourceFile",
		FeedbackTargetFile: "FeedbackTargetFile",
		ConfigMap:          &config2.ConfigMap{Parameters: map[string]string{"key": "value"}},
	}

	accessFeedBackFileCreator := mocks2.NewSyncFeedbackFileCreator(t)
	accessFeedBackFileCreator.EXPECT().Close().Once()

	accessProvidersImport := sync_to_target.AccessProviderImport{
		LastCalculated: time.Now().Unix(),
		AccessProviders: []*sync_to_target.AccessProvider{
			{
				Id:          "AP1",
				Description: "Descr",
				Delete:      false,
				Name:        "Ap1",
				NamingHint:  "NameHint1",
				Action:      sync_to_target.Grant,
			},
			{
				Id:          "AP4",
				Description: "Descr4",
				Delete:      true,
				Name:        "Ap4",
				NamingHint:  "NameHint4",
				Action:      sync_to_target.Grant,
			},
		},
	}

	accessProviderParser := mocks2.NewAccessProviderImportFileParser(t)
	accessProviderParser.EXPECT().ParseAccessProviders().Return(&accessProvidersImport, nil).Once()

	syncerMock := NewMockAccessProviderSyncer(t)
	syncerMock.EXPECT().SyncAccessProviderToTarget(mock.Anything, &accessProvidersImport, accessFeedBackFileCreator, config.ConfigMap).Return(fmt.Errorf("boom")).Once()

	syncFunction := DataAccessSyncFunction{
		Syncer: syncerMock,
		accessFeedbackFileCreatorFactory: func(config *access_provider.AccessSyncToTarget) (sync_to_target.SyncFeedbackFileCreator, error) {
			return accessFeedBackFileCreator, nil
		},
		accessProviderParserFactory: func(config *access_provider.AccessSyncToTarget) (sync_to_target.AccessProviderImportFileParser, error) {
			return accessProviderParser, nil
		},
	}

	//When
	result, err := syncFunction.SyncToTarget(context.Background(), config)

	//Then
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestDataAccessSyncFunction_SyncToTarget_AccessAsCode_ErrorOnSync(t *testing.T) {
	//Given
	config := &access_provider.AccessSyncToTarget{
		Prefix:             "R",
		SourceFile:         "SourceFile",
		FeedbackTargetFile: "FeedbackTargetFile",
		ConfigMap:          &config2.ConfigMap{Parameters: map[string]string{"key": "value"}},
	}

	accessProvidersImport := sync_to_target.AccessProviderImport{
		LastCalculated: time.Now().Unix(),
		AccessProviders: []*sync_to_target.AccessProvider{
			{
				Id:          "AP1",
				Description: "Descr",
				Delete:      false,
				Name:        "Ap1",
				NamingHint:  "NameHint1",
				Action:      sync_to_target.Grant,
			},
			{
				Id:          "AP2",
				Description: "Descr2",
				Delete:      false,
				Name:        "Ap2",
				NamingHint:  "NameHint2",
				Action:      sync_to_target.Grant,
			},
		},
	}

	accessProviderParser := mocks2.NewAccessProviderImportFileParser(t)
	accessProviderParser.EXPECT().ParseAccessProviders().Return(&accessProvidersImport, nil).Once()

	syncerMock := NewMockAccessProviderSyncer(t)
	syncerMock.EXPECT().SyncAccessAsCodeToTarget(mock.Anything, &accessProvidersImport, "R", config.ConfigMap).Return(fmt.Errorf("boom")).Once()

	syncFunction := DataAccessSyncFunction{
		Syncer: syncerMock,
		accessProviderParserFactory: func(config *access_provider.AccessSyncToTarget) (sync_to_target.AccessProviderImportFileParser, error) {
			return accessProviderParser, nil
		},
	}

	//When
	result, err := syncFunction.SyncToTarget(context.Background(), config)

	//Then
	assert.Error(t, err)
	assert.Nil(t, result)
}

//	func TestDataAccessSyncFunction_SyncToTarget_AccessProviders_ErrorOnNameGeneratorFactory(t *testing.T) {
//		//Given
//		config := &access_provider.AccessSyncToTarget{
//			SourceFile:         "SourceFile",
//			FeedbackTargetFile: "FeedbackTargetFile",
//			ConfigMap:          config2.ConfigMap{Parameters: map[string]interface{}{"key": "value"}},
//		}
//
//		accessFeedBackFileCreator := mocks2.NewSyncFeedbackFileCreator(t)
//		accessFeedBackFileCreator.EXPECT().Close().Once()
//
//		accessProviders := []sync_to_target.AccessProvider{
//			{
//				Access: []*sync_to_target.Access{
//					{
//						Id: "AccessId1",
//					},
//					{
//						Id: "AccessId2",
//					},
//				},
//				Id:          "AP1",
//				Description: "Descr",
//				Delete:      false,
//				Name:        "Ap1",
//				NamingHint:  "NameHint1",
//				Action:      sync_to_target.Grant,
//			},
//			{
//				Access: []*sync_to_target.Access{
//					{
//						Id: "AccessId3",
//					},
//				},
//				Id:          "AP2",
//				Description: "Descr2",
//				Delete:      false,
//				Name:        "Ap2",
//				NamingHint:  "NameHint2",
//				Action:      sync_to_target.Grant,
//			},
//		}
//
//		accessProviderParser := mocks2.NewAccessProviderImportFileParser(t)
//		accessProviderParser.EXPECT().ParseAccessProviders().Return(&sync_to_target.AccessProviderImport{
//			LastCalculated:  time.Now().Unix(),
//			AccessProviders: accessProviders,
//		}, nil).Once()
//
//		syncerMock := NewMockAccessProviderSyncer(t)
//		syncFunction := dataAccessSyncFunction{
//			syncer: syncerMock,
//			accessFeedbackFileCreatorFactory: func(config *access_provider.AccessSyncToTarget) (sync_to_target.SyncFeedbackFileCreator, error) {
//				return accessFeedBackFileCreator, nil
//			},
//			accessProviderParserFactory: func(config *access_provider.AccessSyncToTarget) (sync_to_target.AccessProviderImportFileParser, error) {
//				return accessProviderParser, nil
//			},
//			namingConstraints: naming_hint.NamingConstraints{
//				UpperCaseLetters:  true,
//				Numbers:           true,
//				SpecialCharacters: "",
//				MaxLength:         24,
//			},
//		}
//
//		//When
//		result := syncFunction.SyncToTarget(config)
//
//		//Then
//		assert.NotNil(t, result.Error)
//	}
//
//	func TestDataAccessSyncFunction_SyncToTarget_AccessAsCode_ErrorOnNameGeneratorFactory(t *testing.T) {
//		//Given
//		config := &access_provider.AccessSyncToTarget{
//			Prefix:             "R",
//			SourceFile:         "SourceFile",
//			FeedbackTargetFile: "FeedbackTargetFile",
//			ConfigMap:          config2.ConfigMap{Parameters: map[string]interface{}{"key": "value"}},
//		}
//
//		accessProviders := []sync_to_target.AccessProvider{
//			{
//				Access: []*sync_to_target.Access{
//					{
//						Id: "AccessId1",
//					},
//					{
//						Id: "AccessId2",
//					},
//				},
//				Id:          "AP1",
//				Description: "Descr",
//				Delete:      false,
//				Name:        "Ap1",
//				NamingHint:  "NameHint1",
//				Action:      sync_to_target.Grant,
//			},
//			{
//				Access: []*sync_to_target.Access{
//					{
//						Id: "AccessId3",
//					},
//				},
//				Id:          "AP2",
//				Description: "Descr2",
//				Delete:      false,
//				Name:        "Ap2",
//				NamingHint:  "NameHint2",
//				Action:      sync_to_target.Grant,
//			},
//		}
//
//		accessProviderParser := mocks2.NewAccessProviderImportFileParser(t)
//		accessProviderParser.EXPECT().ParseAccessProviders().Return(&sync_to_target.AccessProviderImport{
//			LastCalculated:  time.Now().Unix(),
//			AccessProviders: accessProviders,
//		}, nil).Once()
//
//		syncerMock := NewMockAccessProviderSyncer(t)
//
//		syncFunction := dataAccessSyncFunction{
//			syncer: syncerMock,
//			accessProviderParserFactory: func(config *access_provider.AccessSyncToTarget) (sync_to_target.AccessProviderImportFileParser, error) {
//				return accessProviderParser, nil
//			},
//			namingConstraints: naming_hint.NamingConstraints{
//				UpperCaseLetters:  true,
//				Numbers:           true,
//				SpecialCharacters: "",
//				MaxLength:         24,
//			},
//		}
//
//		//When
//		result := syncFunction.SyncToTarget(config)
//
//		//Then
//		assert.NotNil(t, result.Error)
//	}
func TestDataAccessSync(t *testing.T) {
	//Given
	syncerMock := NewMockAccessProviderSyncer(t)

	//When
	syncFunction := DataAccessSync(syncerMock)

	//Then
	assert.Equal(t, syncerMock, syncFunction.Syncer)
}
