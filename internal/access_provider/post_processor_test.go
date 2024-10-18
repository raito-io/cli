package access_provider

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/smithy-go/ptr"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	baseAp "github.com/raito-io/cli/base/access_provider"
	"github.com/raito-io/cli/base/access_provider/sync_from_target"
	"github.com/raito-io/cli/base/access_provider/sync_from_target/mocks"
	"github.com/raito-io/cli/base/constants"
	"github.com/raito-io/cli/base/tag"
	"github.com/raito-io/cli/internal/access_provider/post_processing"
	mocks2 "github.com/raito-io/cli/internal/access_provider/post_processing/mocks"
)

var logger = hclog.L()

func TestPostProcessor_PostProcess(t *testing.T) {

	type args struct {
		ctx    context.Context
		config *PostProcessorConfig
	}

	type want struct {
		touchedAps   int
		processedAps []*sync_from_target.AccessProvider
	}
	type fields struct {
		setup func(accessProviderFileCreator *mocks.AccessProviderFileCreator, accessProviderFileReader *mocks2.PostProcessorSourceFileParser, want want) (accessProviderFileCreatorError error, accessProviderFileReaderError error)
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    want
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "rewrite name",
			fields: fields{
				setup: func(accessProviderFileCreator *mocks.AccessProviderFileCreator, accessProviderFileReader *mocks2.PostProcessorSourceFileParser, want want) (accessProviderFileCreatorError error, accessProviderFileReaderError error) {
					accessProviderFileCreator.EXPECT().Close().Return().Once()
					accessProviderFileCreator.EXPECT().GetAccessProviderCount().Return(1).Once()

					if len(want.processedAps) > 0 {
						for _, ap := range want.processedAps {
							accessProviderFileCreator.EXPECT().AddAccessProviders(ap).Return(nil).Once()
						}
					}

					accessProviderFileReader.EXPECT().ParseAccessProviders().Return([]*sync_from_target.AccessProvider{
						{Name: "OLD_NAME", Tags: []*tag.Tag{
							{Key: "RANDOM", Value: "VALUE"},
							{Key: "kEy", Value: "NEW_VALUE"},
						},
						}}, nil).Once()

					accessProviderFileCreatorError = nil
					accessProviderFileReaderError = nil

					return accessProviderFileCreatorError, accessProviderFileReaderError
				},
			},
			args: args{
				ctx: context.Background(),
				config: &PostProcessorConfig{
					TagOverwriteKeyForName: "key",
					TargetLogger:           logger,
				},
			},
			want: want{
				touchedAps: 1,
				processedAps: []*sync_from_target.AccessProvider{{
					Name: "NEW_VALUE",
					Tags: []*tag.Tag{
						{Key: "RANDOM", Value: "VALUE"},
						{Key: "kEy", Value: "NEW_VALUE"},
					},
					NameLocked:       ptr.Bool(true),
					NameLockedReason: ptr.String(nameTagOverrideLockedReason),
				}},
			},
			wantErr: require.NoError,
		},
		{
			name: "no files touched",
			fields: fields{
				setup: func(accessProviderFileCreator *mocks.AccessProviderFileCreator, accessProviderFileReader *mocks2.PostProcessorSourceFileParser, want want) (accessProviderFileCreatorError error, accessProviderFileReaderError error) {
					accessProviderFileCreator.EXPECT().Close().Return().Once()
					accessProviderFileCreator.EXPECT().GetAccessProviderCount().Return(2).Once()

					if len(want.processedAps) > 0 {
						for _, ap := range want.processedAps {
							accessProviderFileCreator.EXPECT().AddAccessProviders(ap).Return(nil).Once()
						}
					}
					accessProviderFileReader.EXPECT().ParseAccessProviders().Return(want.processedAps, nil).Once()

					accessProviderFileCreatorError = nil
					accessProviderFileReaderError = nil

					return accessProviderFileCreatorError, accessProviderFileReaderError
				},
			},
			args: args{
				ctx: context.Background(),
				config: &PostProcessorConfig{
					TagOverwriteKeyForName: "",
					TargetLogger:           logger,
				},
			},
			want: want{
				touchedAps: 0,
				processedAps: []*sync_from_target.AccessProvider{
					{
						Name: "OLD_NAME",
						Tags: []*tag.Tag{
							{Key: "RANDOM", Value: "VALUE"},
							{Key: "kEy", Value: "NEW_VALUE"},
						},
					},
					{
						Name: "OLD_NAME_2",
					},
				},
			},
			wantErr: require.NoError,
		},
		{
			name: "error-on-file-creation",
			fields: fields{
				setup: func(accessProviderFileCreator *mocks.AccessProviderFileCreator, accessProviderFileReader *mocks2.PostProcessorSourceFileParser, want want) (accessProviderFileCreatorError error, accessProviderFileReaderError error) {
					accessProviderFileCreatorError = errors.New("BOOM!")
					accessProviderFileReaderError = nil

					return accessProviderFileCreatorError, accessProviderFileReaderError
				},
			},
			args: args{
				ctx: context.Background(),
				config: &PostProcessorConfig{
					TagOverwriteKeyForName: "key",
				},
			},
			want: want{
				touchedAps: 0,
			},
			wantErr: require.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFileCreator, mockFileReader := createPostProcessor(t)
			accessProviderFileCreatorError, accessProviderFileReaderError := tt.fields.setup(mockFileCreator, mockFileReader, tt.want)

			postProcessorFn := PostProcessor{
				accessFileCreatorFactory: func(config *baseAp.AccessSyncFromTarget) (sync_from_target.AccessProviderFileCreator, error) {
					return mockFileCreator, accessProviderFileCreatorError
				},
				accessProviderParserFactory: func(sourceFile string) (post_processing.PostProcessorSourceFileParser, error) {
					return mockFileReader, accessProviderFileReaderError
				},
				config: tt.args.config,
			}

			result, err := postProcessorFn.PostProcess("", "")
			tt.wantErr(t, err)

			if err != nil {
				return
			}

			assert.Equal(t, tt.want.touchedAps, result.AccessProviderTouchedCount)
		})
	}
}

func TestPostProcessor_matchedWithTagKey(t *testing.T) {

	type args struct {
		ctx               context.Context
		config            *PostProcessorConfig
		tagKeySearchValue string
		tag               *tag.Tag
	}

	type want struct {
		canMerge bool
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "exact match",
			args: args{
				ctx: context.Background(),
				config: &PostProcessorConfig{
					TagOverwriteKeyForName: "key",
					TargetLogger:           logger,
				},

				tagKeySearchValue: "key",
				tag:               &tag.Tag{Key: "key", Value: "NEW_VALUE"},
			},
			want: want{
				canMerge: true,
			},
		},
		{
			name: "non exact match",
			args: args{
				ctx: context.Background(),
				config: &PostProcessorConfig{
					TagOverwriteKeyForName: "key",
					TargetLogger:           logger,
				},

				tagKeySearchValue: "KeY",
				tag:               &tag.Tag{Key: "kEy", Value: "NEW_VALUE"},
			},
			want: want{
				canMerge: true,
			},
		},
		{
			name: "no match",
			args: args{
				ctx: context.Background(),
				config: &PostProcessorConfig{
					TagOverwriteKeyForName: "key",
					TargetLogger:           logger,
				},

				tagKeySearchValue: "KEEEEY",
				tag:               &tag.Tag{Key: "key", Value: "NEW_VALUE"},
			},
			want: want{
				canMerge: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFileCreator, mockFileReader := createPostProcessor(t)

			postProcessorFn := PostProcessor{
				accessFileCreatorFactory: func(config *baseAp.AccessSyncFromTarget) (sync_from_target.AccessProviderFileCreator, error) {
					return mockFileCreator, nil
				},
				accessProviderParserFactory: func(sourceFile string) (post_processing.PostProcessorSourceFileParser, error) {
					return mockFileReader, nil
				},
				config: tt.args.config,
			}

			canMerge := postProcessorFn.matchedWithTagKey(tt.args.tagKeySearchValue, tt.args.tag)

			assert.Equal(t, tt.want.canMerge, canMerge)
		})
	}
}

func TestPostProcessor_overwriteName(t *testing.T) {

	type args struct {
		ctx    context.Context
		config *PostProcessorConfig
		ap     *sync_from_target.AccessProvider
		tag    *tag.Tag
	}

	type want struct {
		touched     bool
		processedAp *sync_from_target.AccessProvider
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "rewrite name",
			args: args{
				ctx: context.Background(),
				config: &PostProcessorConfig{
					TagOverwriteKeyForName: "key",
					TargetLogger:           logger,
				},
				ap: &sync_from_target.AccessProvider{
					Name: "OLD_NAME",
					Tags: []*tag.Tag{
						{Key: "RANDOM", Value: "VALUE"},
						{Key: "kEy", Value: "NEW_VALUE"},
					},
				},
				tag: &tag.Tag{Key: "kEy", Value: "NEW_VALUE"},
			},
			want: want{
				touched: true,
				processedAp: &sync_from_target.AccessProvider{
					Name: "NEW_VALUE",
					Tags: []*tag.Tag{
						{Key: "RANDOM", Value: "VALUE"},
						{Key: "kEy", Value: "NEW_VALUE"},
					},
					NameLocked:       ptr.Bool(true),
					NameLockedReason: ptr.String(nameTagOverrideLockedReason),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFileCreator, mockFileReader := createPostProcessor(t)

			postProcessorFn := PostProcessor{
				accessFileCreatorFactory: func(config *baseAp.AccessSyncFromTarget) (sync_from_target.AccessProviderFileCreator, error) {
					return mockFileCreator, nil
				},
				accessProviderParserFactory: func(sourceFile string) (post_processing.PostProcessorSourceFileParser, error) {
					return mockFileReader, nil
				},
				config: tt.args.config,
			}

			touched := postProcessorFn.overwriteName(tt.args.ap, tt.args.tag)

			assert.Equal(t, tt.want.touched, touched)
			assert.Equal(t, tt.want.processedAp, tt.args.ap)
		})
	}
}

func TestPostProcessor_processOverwriteOwners(t *testing.T) {
	type fields struct {
		accessFileCreatorFactory    func(config *baseAp.AccessSyncFromTarget) (sync_from_target.AccessProviderFileCreator, error)
		accessProviderParserFactory func(sourceFile string) (post_processing.PostProcessorSourceFileParser, error)
		config                      *PostProcessorConfig
	}
	type args struct {
		accessProvider *sync_from_target.AccessProvider
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		expectedAp  *sync_from_target.AccessProvider
		wantTouched bool
	}{
		{
			name: "No owners set",
			fields: fields{
				config: &PostProcessorConfig{
					TagOverwriteKeyForOwners: "owners_overwrite",
					TargetLogger:             hclog.NewNullLogger(),
				},
			},
			args: args{
				accessProvider: &sync_from_target.AccessProvider{
					ExternalId: "apId1",
					Name:       "apName1",
					NamingHint: "apName1",
					Tags: []*tag.Tag{
						{Key: "random_tag", Value: "user1"},
					},
				},
			},
			expectedAp: &sync_from_target.AccessProvider{
				ExternalId: "apId1",
				Name:       "apName1",
				NamingHint: "apName1",
				Tags: []*tag.Tag{
					{Key: "random_tag", Value: "user1"},
				},
			},
			wantTouched: false,
		},
		{
			name: "Owners set by raito_onwers tag",
			fields: fields{
				config: &PostProcessorConfig{
					TagOverwriteKeyForOwners: "owners_overwrite",
					TargetLogger:             hclog.NewNullLogger(),
				},
			},
			args: args{
				accessProvider: &sync_from_target.AccessProvider{
					ExternalId: "apId1",
					Name:       "apName1",
					NamingHint: "apName1",
					Tags: []*tag.Tag{
						{Key: "random_tag", Value: "user1"},
						{
							Key:    constants.RaitoOwnerTagKey,
							Value:  "user1, email:user2@raito.io",
							Source: "source1",
						},
					},
				},
			},
			expectedAp: &sync_from_target.AccessProvider{
				ExternalId: "apId1",
				Name:       "apName1",
				NamingHint: "apName1",
				Tags: []*tag.Tag{
					{Key: "random_tag", Value: "user1"}, {
						Key:    constants.RaitoOwnerTagKey,
						Value:  "user1,email:user2@raito.io",
						Source: "source1",
					},
				},
				OwnersLocked:       ptr.Bool(true),
				OwnersLockedReason: ptr.String(ownersTagOverrideLockedReason),
			},
			wantTouched: true,
		},
		{
			name: "Owners set by specific tag",
			fields: fields{
				config: &PostProcessorConfig{
					TagOverwriteKeyForOwners: "owners_overwrite",
					TargetLogger:             hclog.NewNullLogger(),
				},
			},
			args: args{
				accessProvider: &sync_from_target.AccessProvider{
					ExternalId: "apId1",
					Name:       "apName1",
					NamingHint: "apName1",
					Tags: []*tag.Tag{
						{Key: "random_tag", Value: "user1"},
						{
							Key:    "owners_overwrite",
							Value:  "user1, email:user2@raito.io",
							Source: "source1",
						},
					},
				},
			},
			expectedAp: &sync_from_target.AccessProvider{
				ExternalId: "apId1",
				Name:       "apName1",
				NamingHint: "apName1",
				Tags: []*tag.Tag{
					{Key: "random_tag", Value: "user1"},
					{
						Key:    "owners_overwrite",
						Value:  "user1, email:user2@raito.io",
						Source: "source1",
					},
					{
						Key:    constants.RaitoOwnerTagKey,
						Value:  "user1,email:user2@raito.io",
						Source: "source1",
					},
				},
				OwnersLocked:       ptr.Bool(true),
				OwnersLockedReason: ptr.String(ownersTagOverrideLockedReason),
			},
			wantTouched: true,
		},
		{
			name: "Owners set by specific tag and raito_owners tag",
			fields: fields{
				config: &PostProcessorConfig{
					TagOverwriteKeyForOwners: "owners_overwrite",
					TargetLogger:             hclog.NewNullLogger(),
				},
			},
			args: args{
				accessProvider: &sync_from_target.AccessProvider{
					ExternalId: "apId1",
					Name:       "apName1",
					NamingHint: "apName1",
					Tags: []*tag.Tag{
						{Key: "random_tag", Value: "user1"},
						{
							Key:    "owners_overwrite",
							Value:  ", user1, email:user2@raito.io",
							Source: "source1",
						},
						{
							Key:    constants.RaitoOwnerTagKey,
							Value:  " user3, email:user4@raito.io",
							Source: "source1",
						},
						{
							Key:    "owners_overwrite",
							Value:  ",",
							Source: "source1",
						},
					},
				},
			},
			expectedAp: &sync_from_target.AccessProvider{
				ExternalId: "apId1",
				Name:       "apName1",
				NamingHint: "apName1",
				Tags: []*tag.Tag{
					{Key: "random_tag", Value: "user1"},
					{
						Key:    "owners_overwrite",
						Value:  ", user1, email:user2@raito.io",
						Source: "source1",
					},
					{
						Key:    constants.RaitoOwnerTagKey,
						Value:  "user3,email:user4@raito.io,user1,email:user2@raito.io",
						Source: "source1",
					},
					{
						Key:    "owners_overwrite",
						Value:  ",",
						Source: "source1",
					},
				},
				OwnersLocked:       ptr.Bool(true),
				OwnersLockedReason: ptr.String(ownersTagOverrideLockedReason),
			},
			wantTouched: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PostProcessor{
				accessFileCreatorFactory:    tt.fields.accessFileCreatorFactory,
				accessProviderParserFactory: tt.fields.accessProviderParserFactory,
				config:                      tt.fields.config,
			}
			assert.Equalf(t, tt.wantTouched, p.processOverwriteOwners(tt.args.accessProvider), "processOverwriteOwners(%v)", tt.args.accessProvider)
			assert.Equalf(t, tt.expectedAp, tt.args.accessProvider, "accessProvider after processOverwriteOwners()")
		})
	}
}

func createPostProcessor(t *testing.T) (*mocks.AccessProviderFileCreator, *mocks2.PostProcessorSourceFileParser) {
	accessProviderFileCreator := mocks.NewAccessProviderFileCreator(t)
	accessProviderFileReader := mocks2.NewPostProcessorSourceFileParser(t)

	return accessProviderFileCreator, accessProviderFileReader
}
