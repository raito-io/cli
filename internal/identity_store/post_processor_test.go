package identity_store

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/smithy-go/ptr"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/raito-io/cli/base/identity_store"
	"github.com/raito-io/cli/base/identity_store/mocks"
	"github.com/raito-io/cli/base/tag"
)

var logger = hclog.L()

func TestPostProcessor_PostProcess(t *testing.T) {

	type args struct {
		ctx    context.Context
		config *PostProcessorConfig
	}

	type want struct {
		touchedUsers   int
		processedUsers []*identity_store.User
	}

	type fields struct {
		setup            func(identityStoreFileCreator *mocks.IdentityStoreFileCreator) (identityStoreFileCreatorError error)
		toProcessesUsers []*identity_store.User
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    want
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "no config",
			fields: fields{
				setup: func(accessProviderFileCreator *mocks.IdentityStoreFileCreator) (identityStoreFileCreatorError error) {
					accessProviderFileCreator.EXPECT().GetUserCount().Return(1).Once()

					identityStoreFileCreatorError = nil

					return identityStoreFileCreatorError
				},
				toProcessesUsers: []*identity_store.User{{
					Name: "user_1",
					Tags: []*tag.Tag{
						{Key: "RANDOM", Value: "VALUE"},
					},
				}},
			},
			args: args{
				ctx: context.Background(),
				config: &PostProcessorConfig{
					TagKeyAndValueForUserIsMachine: "",
					TargetLogger:                   logger,
				},
			},
			want: want{
				touchedUsers: 0,
				processedUsers: []*identity_store.User{{
					Name: "user_1",
					Tags: []*tag.Tag{
						{Key: "RANDOM", Value: "VALUE"},
					},
				}},
			},
			wantErr: require.NoError,
		},
		{
			name: "wrong config",
			fields: fields{
				setup: func(accessProviderFileCreator *mocks.IdentityStoreFileCreator) (identityStoreFileCreatorError error) {
					accessProviderFileCreator.EXPECT().GetUserCount().Return(1).Once()

					identityStoreFileCreatorError = nil

					return identityStoreFileCreatorError
				},
				toProcessesUsers: []*identity_store.User{{
					Name: "user_1",
					Tags: []*tag.Tag{
						{Key: "RANDOM", Value: "VALUE"},
					},
				}},
			},
			args: args{
				ctx: context.Background(),
				config: &PostProcessorConfig{
					TagKeyAndValueForUserIsMachine: "wrong,config",
					TargetLogger:                   logger,
				},
			},
			want: want{
				touchedUsers: 0,
				processedUsers: []*identity_store.User{{
					Name: "user_1",
					Tags: []*tag.Tag{
						{Key: "RANDOM", Value: "VALUE"},
					},
				}},
			},
			wantErr: require.NoError,
		}, {
			name: "a machine user",
			fields: fields{
				setup: func(accessProviderFileCreator *mocks.IdentityStoreFileCreator) (identityStoreFileCreatorError error) {
					accessProviderFileCreator.EXPECT().GetUserCount().Return(1).Once()

					identityStoreFileCreatorError = nil

					return identityStoreFileCreatorError
				},
				toProcessesUsers: []*identity_store.User{{
					Name: "user_1",
					Tags: []*tag.Tag{
						{Key: "RANDOM", Value: "VALUE"},
						{Key: "user_type", Value: "machine"},
					},
				}},
			},
			args: args{
				ctx: context.Background(),
				config: &PostProcessorConfig{
					TagKeyAndValueForUserIsMachine: "user_type:machine",
					TargetLogger:                   logger,
				},
			},
			want: want{
				touchedUsers: 1,
				processedUsers: []*identity_store.User{{
					Name: "user_1",
					Tags: []*tag.Tag{
						{Key: "RANDOM", Value: "VALUE"},
						{Key: "user_type", Value: "machine"},
					},
					IsMachine: ptr.Bool(true),
				}},
			},
			wantErr: require.NoError,
		}, {
			name: "two machine users and a normal user",
			fields: fields{
				setup: func(accessProviderFileCreator *mocks.IdentityStoreFileCreator) (identityStoreFileCreatorError error) {
					accessProviderFileCreator.EXPECT().GetUserCount().Return(3).Once()

					identityStoreFileCreatorError = nil

					return identityStoreFileCreatorError
				},
				toProcessesUsers: []*identity_store.User{{
					Name: "user_1",
					Tags: []*tag.Tag{
						{Key: "user_type", Value: "machine"},
					},
				}, {
					Name: "user_2",
					Tags: []*tag.Tag{
						{Key: "RANDOM", Value: "VALUE"},
						{Key: "user_type", Value: "machine"},
					},
				}, {
					Name: "user_3",
					Tags: []*tag.Tag{
						{Key: "user_type", Value: "user"},
					},
				}},
			},
			args: args{
				ctx: context.Background(),
				config: &PostProcessorConfig{
					TagKeyAndValueForUserIsMachine: "user_type:machine",
					TargetLogger:                   logger,
				},
			},
			want: want{
				touchedUsers: 2,
				processedUsers: []*identity_store.User{{
					Name: "user_1",
					Tags: []*tag.Tag{
						{Key: "user_type", Value: "machine"},
					},
					IsMachine: ptr.Bool(true),
				}, {
					Name: "user_2",
					Tags: []*tag.Tag{
						{Key: "RANDOM", Value: "VALUE"},
						{Key: "user_type", Value: "machine"},
					},
					IsMachine: ptr.Bool(true),
				}, {
					Name: "user_3",
					Tags: []*tag.Tag{
						{Key: "user_type", Value: "user"},
					},
				}},
			},
			wantErr: require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "*")
			if err != nil {
				t.Fatal(err)
			}

			defer os.RemoveAll(tmpDir)

			inputFilePath := filepath.Join(tmpDir, "input.json")
			outputFilePath := filepath.Join(tmpDir, "output.json")

			rawJson, err := json.Marshal(tt.fields.toProcessesUsers)
			if err != nil {
				t.Fatal(err)
			}

			err = os.WriteFile(inputFilePath, rawJson, 0644)
			if err != nil {
				t.Fatal(err)
			}

			mockFileCreator := mocks.NewIdentityStoreFileCreator(t)
			identityStoreFileCreatorError := tt.fields.setup(mockFileCreator)

			if len(tt.want.processedUsers) > 0 {
				for _, user := range tt.want.processedUsers {
					mockFileCreator.EXPECT().AddUsers(user).Return(nil).Once()
				}
			}

			postProcessorFn := PostProcessor{
				identityStoreFileCreator: func(config *identity_store.IdentityStoreSyncConfig) (identity_store.IdentityStoreFileCreator, error) {
					return mockFileCreator, identityStoreFileCreatorError
				},
				config: tt.args.config,
			}

			result, err := postProcessorFn.PostProcessUsers(inputFilePath, outputFilePath)
			tt.wantErr(t, err)

			if err != nil {
				return
			}

			assert.Equal(t, tt.want.touchedUsers, result.UsersTouchedCount)
		})
	}
}
