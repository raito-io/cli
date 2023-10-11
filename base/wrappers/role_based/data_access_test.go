package role_based

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/smithy-go/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/raito-io/cli/base/access_provider/sync_from_target/mocks"
	"github.com/raito-io/cli/base/access_provider/sync_to_target"
	mocks2 "github.com/raito-io/cli/base/access_provider/sync_to_target/mocks"
	"github.com/raito-io/cli/base/access_provider/sync_to_target/naming_hint"
	"github.com/raito-io/cli/base/util/config"
)

func TestAccessProviderRoleSyncFunction_SyncAccessProviderFromTarget(t *testing.T) {
	//Given
	configMap := config.ConfigMap{Parameters: map[string]string{"key": "value"}}

	accessProviderFileCreator := mocks.NewAccessProviderFileCreator(t)

	syncerMock := NewMockAccessProviderRoleSyncer(t)
	syncerMock.EXPECT().SyncAccessProvidersFromTarget(mock.Anything, accessProviderFileCreator, &configMap).Return(nil).Once()

	syncFunction := accessProviderRoleSyncFunction{
		syncer: syncerMock,
	}

	//When
	err := syncFunction.SyncAccessProvidersFromTarget(context.Background(), accessProviderFileCreator, &configMap)

	//Then
	assert.NoError(t, err)
}

func TestAccessProviderRoleSyncFunction_SyncAccessProviderFromTarget_WithError(t *testing.T) {
	//Given
	configMap := config.ConfigMap{Parameters: map[string]string{"key": "value"}}

	accessProviderFileCreator := mocks.NewAccessProviderFileCreator(t)

	syncerMock := NewMockAccessProviderRoleSyncer(t)
	syncerMock.EXPECT().SyncAccessProvidersFromTarget(mock.Anything, accessProviderFileCreator, &configMap).Return(fmt.Errorf("boom")).Once()

	syncFunction := accessProviderRoleSyncFunction{
		syncer: syncerMock,
	}

	//When
	err := syncFunction.SyncAccessProvidersFromTarget(context.Background(), accessProviderFileCreator, &configMap)

	//Then
	assert.Error(t, err)
}

func TestAccessProviderRoleSyncFunction_SyncAccessAsCodeToTarget(t *testing.T) {
	//Given
	configMap := config.ConfigMap{Parameters: map[string]string{"key": "value"}}

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
				Id:          "AP3",
				Description: "Descr3",
				Delete:      false,
				Name:        "Ap3",
				NamingHint:  "NameHint1",
				Action:      sync_to_target.Grant,
			},
		},
	}

	syncerMock := NewMockAccessProviderRoleSyncer(t)
	syncerMock.EXPECT().SyncAccessAsCodeToTarget(mock.Anything, mock.Anything, "R_", &configMap).Return(nil).Once()

	syncFunction := accessProviderRoleSyncFunction{
		syncer: syncerMock,
		namingConstraints: naming_hint.NamingConstraints{
			UpperCaseLetters:  true,
			Numbers:           true,
			SpecialCharacters: "_",
			MaxLength:         24,
		},
	}

	//when
	err := syncFunction.SyncAccessAsCodeToTarget(context.Background(), &accessProvidersImport, "R", &configMap)

	//Then
	assert.NoError(t, err)
	syncerMock.AssertCalled(t, "SyncAccessAsCodeToTarget", mock.Anything,
		map[string]*sync_to_target.AccessProvider{
			"R_NAME_HINT1":    accessProvidersImport.AccessProviders[0],
			"R_NAME_HINT2":    accessProvidersImport.AccessProviders[1],
			"R_NAME_HINT1__0": accessProvidersImport.AccessProviders[2],
		}, "R_", &configMap,
	)
}

func TestAccessProviderRoleSyncFunction_SyncAccessAsCodeToTarget_NameGeneratorFactoryError(t *testing.T) {
	//Given
	configMap := config.ConfigMap{Parameters: map[string]string{"key": "value"}}

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

	syncerMock := NewMockAccessProviderRoleSyncer(t)

	syncFunction := accessProviderRoleSyncFunction{
		syncer: syncerMock,
		namingConstraints: naming_hint.NamingConstraints{
			UpperCaseLetters:  true,
			Numbers:           true,
			SpecialCharacters: "",
			MaxLength:         24,
		},
	}

	//when
	err := syncFunction.SyncAccessAsCodeToTarget(context.Background(), &accessProvidersImport, "R", &configMap)

	//Then
	assert.Error(t, err)
	syncerMock.AssertNotCalled(t, "SyncAccessAsCodeToTarget")
}

func TestAccessProviderRoleSyncFunction_SyncAccessProviderToTarget(t *testing.T) {
	//Given
	configMap := config.ConfigMap{Parameters: map[string]string{"key": "value"}}

	accessFeedBackFileCreator := mocks2.NewSyncFeedbackFileCreator(t)

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
			{
				Id:          "Mask1",
				Description: "Mask1Description",
				Delete:      false,
				Name:        "Mask1",
				NamingHint:  "NameHintMask1",
				Action:      sync_to_target.Mask,
			},
			{
				Id:          "Mask2",
				Description: "Mask2Description",
				Delete:      true,
				Name:        "Mask2",
				NamingHint:  "NameHintMask2",
				ActualName:  ptr.String("ActualNameMask2"),
				Action:      sync_to_target.Mask,
			},
		},
	}

	syncerMock := NewMockAccessProviderRoleSyncer(t)
	maskCall := syncerMock.EXPECT().SyncAccessProviderMasksToTarget(mock.Anything, map[string]*sync_to_target.AccessProvider{"ActualNameMask2": {
		Id:          "Mask2",
		Description: "Mask2Description",
		Delete:      true,
		Name:        "Mask2",
		NamingHint:  "NameHintMask2",
		ActualName:  ptr.String("ActualNameMask2"),
		Action:      sync_to_target.Mask,
	}}, map[string]*sync_to_target.AccessProvider{"NAME_HINT_MASK1": {
		Id:          "Mask1",
		Description: "Mask1Description",
		Delete:      false,
		Name:        "Mask1",
		NamingHint:  "NameHintMask1",
		Action:      sync_to_target.Mask,
	},
	}, accessFeedBackFileCreator, &configMap).Return(nil).Once()
	syncerMock.EXPECT().SyncAccessProviderRolesToTarget(mock.Anything, map[string]*sync_to_target.AccessProvider{actualName1: {
		ActualName:  &actualName1,
		Id:          "AP3",
		Description: "Descr3",
		Delete:      true,
		Name:        "Ap3",
		NamingHint:  "NameHint3",
		Action:      sync_to_target.Grant,
	}}, mock.Anything, accessFeedBackFileCreator, &configMap).Return(nil).Once().NotBefore(maskCall)

	syncer := accessProviderRoleSyncFunction{
		syncer: syncerMock,
		namingConstraints: naming_hint.NamingConstraints{
			UpperCaseLetters:  true,
			Numbers:           true,
			SpecialCharacters: "_",
			MaxLength:         24,
		},
	}

	//When
	err := syncer.SyncAccessProviderToTarget(context.Background(), &accessProvidersImport, accessFeedBackFileCreator, &configMap)

	//Then
	assert.NoError(t, err)
	syncerMock.AssertCalled(t, "SyncAccessProviderRolesToTarget", mock.Anything,
		map[string]*sync_to_target.AccessProvider{actualName1: {
			ActualName:  &actualName1,
			Id:          "AP3",
			Description: "Descr3",
			Delete:      true,
			Name:        "Ap3",
			NamingHint:  "NameHint3",
			Action:      sync_to_target.Grant,
		}},
		map[string]*sync_to_target.AccessProvider{
			"NAME_HINT1": accessProvidersImport.AccessProviders[0],
			"NAME_HINT2": accessProvidersImport.AccessProviders[1],
		},
		accessFeedBackFileCreator, &configMap,
	)
}

func TestAccessProviderRoleSyncFunction_SyncAccessProviderToTarget_ErrorOnNameGeneratorFactory(t *testing.T) {
	//Given
	configMap := config.ConfigMap{Parameters: map[string]string{"key": "value"}}

	accessFeedBackFileCreator := mocks2.NewSyncFeedbackFileCreator(t)

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
		},
	}

	syncerMock := NewMockAccessProviderRoleSyncer(t)

	syncer := accessProviderRoleSyncFunction{
		syncer: syncerMock,
		namingConstraints: naming_hint.NamingConstraints{
			UpperCaseLetters:  true,
			Numbers:           true,
			SpecialCharacters: "",
			MaxLength:         24,
		},
	}

	//When
	err := syncer.SyncAccessProviderToTarget(context.Background(), &accessProvidersImport, accessFeedBackFileCreator, &configMap)

	//Then
	assert.Error(t, err)
	syncerMock.AssertNotCalled(t, "SyncAccessProvidersToTarget")
}

func TestAccessProviderRoleSync(t *testing.T) {
	//Given
	syncerMock := NewMockAccessProviderRoleSyncer(t)
	nameConstraints := naming_hint.NamingConstraints{
		UpperCaseLetters: true,
	}

	//When
	syncer := AccessProviderRoleSync(syncerMock, nameConstraints)

	//Then
	actualSyncer := syncer.Syncer.(*accessProviderRoleSyncFunction)
	assert.Equal(t, syncerMock, actualSyncer.syncer)
}
