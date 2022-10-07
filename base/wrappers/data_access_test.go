package wrappers

import (
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
	"github.com/raito-io/cli/base/access_provider/sync_to_target/naming_hint"
	config2 "github.com/raito-io/cli/base/util/config"
	error2 "github.com/raito-io/cli/base/util/error"
)

func TestDataAccessSyncFunction_SyncFromTarget(t *testing.T) {
	//Given
	config := &access_provider.AccessSyncFromTarget{
		Prefix:     "prefix",
		TargetFile: "targetFile",
		ConfigMap:  config2.ConfigMap{Parameters: map[string]interface{}{"key": "value"}},
	}

	accessProviderFilCreator := mocks.NewAccessProviderFileCreator(t)
	accessProviderFilCreator.EXPECT().Close().Return().Once()
	accessProviderFilCreator.EXPECT().GetAccessProviderCount().Return(0).Once()

	syncerMock := NewMockAccessProviderSyncer(t)
	syncerMock.EXPECT().SyncAccessProviderFromTarget(mock.Anything, accessProviderFilCreator, &config.ConfigMap).Return(nil).Once()

	syncFunction := dataAccessSyncFunction{
		syncer: syncerMock,
		accessFileCreatorFactory: func(config *access_provider.AccessSyncFromTarget) (sync_from_target.AccessProviderFileCreator, error) {
			return accessProviderFilCreator, nil
		},
	}

	//When
	result := syncFunction.SyncFromTarget(config)

	//Then
	assert.Nil(t, result.Error)
}

func TestDataAccessSyncFunction_SyncFromTarget_ErrorOnFileCreation(t *testing.T) {
	//Given
	config := &access_provider.AccessSyncFromTarget{
		Prefix:     "prefix",
		TargetFile: "targetFile",
		ConfigMap:  config2.ConfigMap{Parameters: map[string]interface{}{"key": "value"}},
	}

	syncerMock := NewMockAccessProviderSyncer(t)

	syncFunction := dataAccessSyncFunction{
		syncer: syncerMock,
		accessFileCreatorFactory: func(config *access_provider.AccessSyncFromTarget) (sync_from_target.AccessProviderFileCreator, error) {
			return nil, &error2.ErrorResult{
				ErrorMessage: "BOOM!",
				ErrorCode:    error2.UnknownError,
			}
		},
	}

	//When
	result := syncFunction.SyncFromTarget(config)

	//Then
	assert.NotNil(t, result.Error)
	assert.Equal(t, "BOOM!", result.Error.ErrorMessage)
	assert.Equal(t, error2.UnknownError, result.Error.ErrorCode)

	syncerMock.AssertNotCalled(t, "SyncAccessProviderFromTarget", mock.Anything, mock.Anything, mock.Anything)
}

func TestDataAccessSyncFunction_SyncFromTarget_ErrorOnSync(t *testing.T) {
	//Given
	config := &access_provider.AccessSyncFromTarget{
		Prefix:     "prefix",
		TargetFile: "targetFile",
		ConfigMap:  config2.ConfigMap{Parameters: map[string]interface{}{"key": "value"}},
	}

	accessProviderFilCreator := mocks.NewAccessProviderFileCreator(t)
	accessProviderFilCreator.EXPECT().Close().Return().Once()

	syncerMock := NewMockAccessProviderSyncer(t)
	syncerMock.EXPECT().SyncAccessProviderFromTarget(mock.Anything, accessProviderFilCreator, &config.ConfigMap).Return(
		&error2.ErrorResult{
			ErrorMessage: "BOOM!",
			ErrorCode:    error2.UnknownError,
		}).Once()

	syncFunction := dataAccessSyncFunction{
		syncer: syncerMock,
		accessFileCreatorFactory: func(config *access_provider.AccessSyncFromTarget) (sync_from_target.AccessProviderFileCreator, error) {
			return accessProviderFilCreator, nil
		},
	}

	//When
	result := syncFunction.SyncFromTarget(config)

	//Then
	assert.NotNil(t, result.Error)
	assert.Equal(t, "BOOM!", result.Error.ErrorMessage)
	assert.Equal(t, error2.UnknownError, result.Error.ErrorCode)
}

func TestDataAccessSyncFunction_SyncToTarget_AccessProviders(t *testing.T) {
	//Given
	config := &access_provider.AccessSyncToTarget{
		SourceFile:         "SourceFile",
		FeedbackTargetFile: "FeedbackTargetFile",
		ConfigMap:          config2.ConfigMap{Parameters: map[string]interface{}{"key": "value"}},
	}

	accessFeedBackFileCreator := mocks2.NewSyncFeedbackFileCreator(t)
	accessFeedBackFileCreator.EXPECT().Close().Once()

	actualName1 := "ActualName1"

	accessProviders := []sync_to_target.AccessProvider{
		{
			Access: []*sync_to_target.Access{
				{
					Id: "AccessId1",
				},
				{
					Id: "AccessId2",
				},
			},
			Id:          "AP1",
			Description: "Descr",
			Delete:      false,
			Name:        "Ap1",
			NamingHint:  "NameHint1",
			Action:      sync_to_target.Grant,
		},
		{
			Access: []*sync_to_target.Access{
				{
					Id: "AccessId3",
				},
			},
			Id:          "AP2",
			Description: "Descr2",
			Delete:      false,
			Name:        "Ap2",
			NamingHint:  "NameHint2",
			Action:      sync_to_target.Grant,
		},
		{
			Access: []*sync_to_target.Access{
				{
					Id:         "SomeAccessToDelete",
					ActualName: &actualName1,
				},
			},
			Id:          "AP3",
			Description: "Descr3",
			Delete:      true,
			Name:        "Ap3",
			NamingHint:  "NameHint3",
			Action:      sync_to_target.Grant,
		},
		{
			Access: []*sync_to_target.Access{
				{
					Id: "AnotherAccessToDelete",
				},
			},
			Id:          "AP4",
			Description: "Descr4",
			Delete:      true,
			Name:        "Ap4",
			NamingHint:  "NameHint4",
			Action:      sync_to_target.Grant,
		},
	}

	accessProviderParser := mocks2.NewAccessProviderImportFileParser(t)
	accessProviderParser.EXPECT().ParseAccessProviders().Return(&sync_to_target.AccessProviderImport{
		LastCalculated:  time.Now().Unix(),
		AccessProviders: accessProviders,
	}, nil).Once()

	syncerMock := NewMockAccessProviderSyncer(t)
	syncerMock.EXPECT().SyncAccessProviderToTarget(mock.Anything, mock.Anything, mock.Anything, accessFeedBackFileCreator, &config.ConfigMap).Return(nil).Once()

	syncFunction := dataAccessSyncFunction{
		syncer: syncerMock,
		accessFeedbackFileCreatorFactory: func(config *access_provider.AccessSyncToTarget) (sync_to_target.SyncFeedbackFileCreator, error) {
			return accessFeedBackFileCreator, nil
		},
		accessProviderParserFactory: func(config *access_provider.AccessSyncToTarget) (sync_to_target.AccessProviderImportFileParser, error) {
			return accessProviderParser, nil
		},
		namingConstraints: naming_hint.NamingConstraints{
			UpperCaseLetters:  true,
			Numbers:           true,
			SpecialCharacters: "_",
			MaxLength:         24,
		},
	}

	//When
	result := syncFunction.SyncToTarget(config)

	//Then
	assert.Nil(t, result.Error)

	syncerMock.AssertCalled(t, "SyncAccessProviderToTarget", mock.Anything,
		[]string{actualName1},
		map[string]sync_to_target.EnrichedAccess{
			"NAME_HINT1": {
				AccessProvider: &accessProviders[0],
				Access:         accessProviders[0].Access[0],
			},
			"NAME_HINT1__0": {
				AccessProvider: &accessProviders[0],
				Access:         accessProviders[0].Access[1],
			},
			"NAME_HINT2": {
				AccessProvider: &accessProviders[1],
				Access:         accessProviders[1].Access[0],
			},
		},
		accessFeedBackFileCreator, &config.ConfigMap,
	)
}

func TestDataAccessSyncFunction_SyncToTarget_AccessAsCode(t *testing.T) {
	//Given
	config := &access_provider.AccessSyncToTarget{
		Prefix:             "R",
		SourceFile:         "SourceFile",
		FeedbackTargetFile: "FeedbackTargetFile",
		ConfigMap:          config2.ConfigMap{Parameters: map[string]interface{}{"key": "value"}},
	}

	accessProviders := []sync_to_target.AccessProvider{
		{
			Access: []*sync_to_target.Access{
				{
					Id: "AccessId1",
				},
				{
					Id: "AccessId2",
				},
			},
			Id:          "AP1",
			Description: "Descr",
			Delete:      false,
			Name:        "Ap1",
			NamingHint:  "NameHint1",
			Action:      sync_to_target.Grant,
		},
		{
			Access: []*sync_to_target.Access{
				{
					Id: "AccessId3",
				},
			},
			Id:          "AP2",
			Description: "Descr2",
			Delete:      false,
			Name:        "Ap2",
			NamingHint:  "NameHint2",
			Action:      sync_to_target.Grant,
		},
	}

	accessProviderParser := mocks2.NewAccessProviderImportFileParser(t)
	accessProviderParser.EXPECT().ParseAccessProviders().Return(&sync_to_target.AccessProviderImport{
		LastCalculated:  time.Now().Unix(),
		AccessProviders: accessProviders,
	}, nil).Once()

	syncerMock := NewMockAccessProviderSyncer(t)
	syncerMock.EXPECT().SyncAccessAsCodeToTarget(mock.Anything, mock.Anything, &config.ConfigMap).Return(nil).Once()

	syncFunction := dataAccessSyncFunction{
		syncer: syncerMock,
		accessProviderParserFactory: func(config *access_provider.AccessSyncToTarget) (sync_to_target.AccessProviderImportFileParser, error) {
			return accessProviderParser, nil
		},
		namingConstraints: naming_hint.NamingConstraints{
			UpperCaseLetters:  true,
			Numbers:           true,
			SpecialCharacters: "_",
			MaxLength:         24,
		},
	}

	//When
	result := syncFunction.SyncToTarget(config)

	//Then
	assert.Nil(t, result.Error)

	syncerMock.AssertCalled(t, "SyncAccessAsCodeToTarget", mock.Anything,
		map[string]sync_to_target.EnrichedAccess{
			"R_NAME_HINT1": {
				AccessProvider: &accessProviders[0],
				Access:         accessProviders[0].Access[0],
			},
			"R_NAME_HINT1__0": {
				AccessProvider: &accessProviders[0],
				Access:         accessProviders[0].Access[1],
			},
			"R_NAME_HINT2": {
				AccessProvider: &accessProviders[1],
				Access:         accessProviders[1].Access[0],
			},
		}, &config.ConfigMap,
	)
}

func TestDataAccessSyncFunction_SyncToTarget_ErrorOnFileParsingFactory(t *testing.T) {
	//Given
	config := &access_provider.AccessSyncToTarget{
		SourceFile:         "SourceFile",
		FeedbackTargetFile: "FeedbackTargetFile",
		ConfigMap:          config2.ConfigMap{Parameters: map[string]interface{}{"key": "value"}},
	}

	syncerMock := NewMockAccessProviderSyncer(t)

	syncFunction := dataAccessSyncFunction{
		syncer: syncerMock,
		accessProviderParserFactory: func(config *access_provider.AccessSyncToTarget) (sync_to_target.AccessProviderImportFileParser, error) {
			return nil, error2.ToErrorResult(fmt.Errorf("boom"))
		},
		namingConstraints: naming_hint.NamingConstraints{
			UpperCaseLetters:  true,
			Numbers:           true,
			SpecialCharacters: "_",
			MaxLength:         24,
		},
	}

	//When
	result := syncFunction.SyncToTarget(config)

	//Then
	assert.NotNil(t, result.Error)
	assert.Equal(t, "boom", result.Error.ErrorMessage)
	assert.Equal(t, error2.UnknownError, result.Error.ErrorCode)

	syncerMock.AssertNotCalled(t, "SyncAccessProviderToTarget", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	syncerMock.AssertNotCalled(t, "SyncAccessAsCodeToTarget", mock.Anything, mock.Anything, mock.Anything)
}

func TestDataAccessSyncFunction_SyncToTarget_ErrorOnFileParsing(t *testing.T) {
	//Given
	config := &access_provider.AccessSyncToTarget{
		SourceFile:         "SourceFile",
		FeedbackTargetFile: "FeedbackTargetFile",
		ConfigMap:          config2.ConfigMap{Parameters: map[string]interface{}{"key": "value"}},
	}

	accessProviderParser := mocks2.NewAccessProviderImportFileParser(t)
	accessProviderParser.EXPECT().ParseAccessProviders().Return(nil, &error2.ErrorResult{
		ErrorMessage: "BOOM!",
		ErrorCode:    error2.SourceConnectionError,
	}).Once()

	syncerMock := NewMockAccessProviderSyncer(t)

	syncFunction := dataAccessSyncFunction{
		syncer: syncerMock,
		accessProviderParserFactory: func(config *access_provider.AccessSyncToTarget) (sync_to_target.AccessProviderImportFileParser, error) {
			return accessProviderParser, nil
		},
		namingConstraints: naming_hint.NamingConstraints{
			UpperCaseLetters:  true,
			Numbers:           true,
			SpecialCharacters: "_",
			MaxLength:         24,
		},
	}

	//When
	result := syncFunction.SyncToTarget(config)

	//Then
	assert.NotNil(t, result.Error)
	assert.Equal(t, "BOOM!", result.Error.ErrorMessage)
	assert.Equal(t, error2.SourceConnectionError, result.Error.ErrorCode)

	syncerMock.AssertNotCalled(t, "SyncAccessProviderToTarget", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	syncerMock.AssertNotCalled(t, "SyncAccessAsCodeToTarget", mock.Anything, mock.Anything, mock.Anything)
}

func TestDataAccessSyncFunction_SyncToTarget_AccessProviders_ErrorOnFeedbackFileCreation(t *testing.T) {
	//Given
	config := &access_provider.AccessSyncToTarget{
		SourceFile:         "SourceFile",
		FeedbackTargetFile: "FeedbackTargetFile",
		ConfigMap:          config2.ConfigMap{Parameters: map[string]interface{}{"key": "value"}},
	}

	accessProviders := []sync_to_target.AccessProvider{
		{
			Access: []*sync_to_target.Access{
				{
					Id: "AccessId1",
				},
				{
					Id: "AccessId2",
				},
			},
			Id:          "AP1",
			Description: "Descr",
			Delete:      false,
			Name:        "Ap1",
			NamingHint:  "NameHint1",
			Action:      sync_to_target.Grant,
		},
		{
			Access: []*sync_to_target.Access{
				{
					Id: "AccessId3",
				},
			},
			Id:          "AP2",
			Description: "Descr2",
			Delete:      false,
			Name:        "Ap2",
			NamingHint:  "NameHint2",
			Action:      sync_to_target.Grant,
		},
	}

	accessProviderParser := mocks2.NewAccessProviderImportFileParser(t)
	accessProviderParser.EXPECT().ParseAccessProviders().Return(&sync_to_target.AccessProviderImport{
		LastCalculated:  time.Now().Unix(),
		AccessProviders: accessProviders,
	}, nil).Once()

	syncerMock := NewMockAccessProviderSyncer(t)

	syncFunction := dataAccessSyncFunction{
		syncer: syncerMock,
		accessFeedbackFileCreatorFactory: func(config *access_provider.AccessSyncToTarget) (sync_to_target.SyncFeedbackFileCreator, error) {
			return nil, fmt.Errorf("boom")
		},
		accessProviderParserFactory: func(config *access_provider.AccessSyncToTarget) (sync_to_target.AccessProviderImportFileParser, error) {
			return accessProviderParser, nil
		},
		namingConstraints: naming_hint.NamingConstraints{
			UpperCaseLetters:  true,
			Numbers:           true,
			SpecialCharacters: "_",
			MaxLength:         24,
		},
	}

	//When
	result := syncFunction.SyncToTarget(config)

	//Then
	assert.NotNil(t, result.Error)
	assert.Equal(t, "boom", result.Error.ErrorMessage)
	assert.Equal(t, error2.UnknownError, result.Error.ErrorCode)

	syncerMock.AssertNotCalled(t, "SyncAccessProviderToTarget", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestDataAccessSyncFunction_SyncToTarget_AccessProviders_ErrorOnNameGeneratorFactory(t *testing.T) {
	//Given
	config := &access_provider.AccessSyncToTarget{
		SourceFile:         "SourceFile",
		FeedbackTargetFile: "FeedbackTargetFile",
		ConfigMap:          config2.ConfigMap{Parameters: map[string]interface{}{"key": "value"}},
	}

	accessFeedBackFileCreator := mocks2.NewSyncFeedbackFileCreator(t)
	accessFeedBackFileCreator.EXPECT().Close().Once()

	accessProviders := []sync_to_target.AccessProvider{
		{
			Access: []*sync_to_target.Access{
				{
					Id: "AccessId1",
				},
				{
					Id: "AccessId2",
				},
			},
			Id:          "AP1",
			Description: "Descr",
			Delete:      false,
			Name:        "Ap1",
			NamingHint:  "NameHint1",
			Action:      sync_to_target.Grant,
		},
		{
			Access: []*sync_to_target.Access{
				{
					Id: "AccessId3",
				},
			},
			Id:          "AP2",
			Description: "Descr2",
			Delete:      false,
			Name:        "Ap2",
			NamingHint:  "NameHint2",
			Action:      sync_to_target.Grant,
		},
	}

	accessProviderParser := mocks2.NewAccessProviderImportFileParser(t)
	accessProviderParser.EXPECT().ParseAccessProviders().Return(&sync_to_target.AccessProviderImport{
		LastCalculated:  time.Now().Unix(),
		AccessProviders: accessProviders,
	}, nil).Once()

	syncerMock := NewMockAccessProviderSyncer(t)
	syncFunction := dataAccessSyncFunction{
		syncer: syncerMock,
		accessFeedbackFileCreatorFactory: func(config *access_provider.AccessSyncToTarget) (sync_to_target.SyncFeedbackFileCreator, error) {
			return accessFeedBackFileCreator, nil
		},
		accessProviderParserFactory: func(config *access_provider.AccessSyncToTarget) (sync_to_target.AccessProviderImportFileParser, error) {
			return accessProviderParser, nil
		},
		namingConstraints: naming_hint.NamingConstraints{
			UpperCaseLetters:  true,
			Numbers:           true,
			SpecialCharacters: "",
			MaxLength:         24,
		},
	}

	//When
	result := syncFunction.SyncToTarget(config)

	//Then
	assert.NotNil(t, result.Error)
}

func TestDataAccessSyncFunction_SyncToTarget_AccessAsCode_ErrorOnNameGeneratorFactory(t *testing.T) {
	//Given
	config := &access_provider.AccessSyncToTarget{
		Prefix:             "R",
		SourceFile:         "SourceFile",
		FeedbackTargetFile: "FeedbackTargetFile",
		ConfigMap:          config2.ConfigMap{Parameters: map[string]interface{}{"key": "value"}},
	}

	accessProviders := []sync_to_target.AccessProvider{
		{
			Access: []*sync_to_target.Access{
				{
					Id: "AccessId1",
				},
				{
					Id: "AccessId2",
				},
			},
			Id:          "AP1",
			Description: "Descr",
			Delete:      false,
			Name:        "Ap1",
			NamingHint:  "NameHint1",
			Action:      sync_to_target.Grant,
		},
		{
			Access: []*sync_to_target.Access{
				{
					Id: "AccessId3",
				},
			},
			Id:          "AP2",
			Description: "Descr2",
			Delete:      false,
			Name:        "Ap2",
			NamingHint:  "NameHint2",
			Action:      sync_to_target.Grant,
		},
	}

	accessProviderParser := mocks2.NewAccessProviderImportFileParser(t)
	accessProviderParser.EXPECT().ParseAccessProviders().Return(&sync_to_target.AccessProviderImport{
		LastCalculated:  time.Now().Unix(),
		AccessProviders: accessProviders,
	}, nil).Once()

	syncerMock := NewMockAccessProviderSyncer(t)

	syncFunction := dataAccessSyncFunction{
		syncer: syncerMock,
		accessProviderParserFactory: func(config *access_provider.AccessSyncToTarget) (sync_to_target.AccessProviderImportFileParser, error) {
			return accessProviderParser, nil
		},
		namingConstraints: naming_hint.NamingConstraints{
			UpperCaseLetters:  true,
			Numbers:           true,
			SpecialCharacters: "",
			MaxLength:         24,
		},
	}

	//When
	result := syncFunction.SyncToTarget(config)

	//Then
	assert.NotNil(t, result.Error)
}

func TestDataAccessSync(t *testing.T) {
	//Given
	syncerMock := NewMockAccessProviderSyncer(t)
	namingConstraints := naming_hint.NamingConstraints{
		UpperCaseLetters:  true,
		Numbers:           true,
		SpecialCharacters: "_",
	}

	//When
	syncFunction := DataAccessSync(syncerMock, namingConstraints)

	//Then
	assert.Equal(t, syncerMock, syncFunction.syncer)
	assert.Equal(t, namingConstraints, syncFunction.namingConstraints)
}
