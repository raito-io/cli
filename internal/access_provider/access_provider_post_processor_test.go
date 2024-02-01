package access_provider

// import (
// 	"context"
// 	"errors"
// 	"testing"
//
//

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/require"

// 	"github.com/raito-io/cli/base/access_provider"
// 	"github.com/raito-io/cli/base/access_provider/sync_from_target"
// 	"github.com/raito-io/cli/base/access_provider/sync_from_target/mocks"
// 	mocks3 "github.com/raito-io/cli/base/access_provider/sync_from_target/mocks"
// 	"github.com/raito-io/cli/base/access_provider_post_processor"
// 	"github.com/raito-io/cli/base/tag"
// 	config2 "github.com/raito-io/cli/base/util/config"
// )

// func TestAccessProviderPostProcessor(t *testing.T) {

// 	type args struct {
// 		ctx    context.Context
// 		config *access_provider_post_processor.AccessProviderPostProcessorConfig
// 	}
// 	type fields struct {
// 		setup func(accessProviderFileCreator *mocks3.AccessProviderFileCreator, accessProviderFileReader *mocks3.AccessProviderSyncFromTargetFileParser, postProcessor *MockAccessProviderPostProcessorI, args args) (accessProviderFileCreatorError error, accessProviderFileReaderError error)
// 	}

// 	type want struct {
// 		touchedAps int32
// 	}

// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    want
// 		wantErr require.ErrorAssertionFunc
// 	}{
// 		{
// 			name: "basic",
// 			fields: fields{
// 				setup: func(accessProviderFileCreator *mocks3.AccessProviderFileCreator, accessProviderFileReader *mocks3.AccessProviderSyncFromTargetFileParser, postProcessor *MockAccessProviderPostProcessorI, args args) (accessProviderFileCreatorError error, accessProviderFileReaderError error) {
// 					accessProviderFileCreator.EXPECT().Close().Return().Once()
// 					accessProviderFileCreator.EXPECT().GetAccessProviderCount().Return(1).Once()

// 					accessProviderFileReader.EXPECT().ParseAccessProviders().Return([]*sync_from_target.AccessProvider{
// 						{Name: "OLD_NAME", Tags: []*tag.Tag{
// 							{Key: "RANDOM", Value: "VALUE"},
// 							{Key: "kEy", Value: "NEW_VALUE"},
// 						},
// 						}}, nil).Once()

// 					postProcessor.EXPECT().Initialize(mock.Anything, accessProviderFileCreator, args.config).Return(nil).Once()
// 					postProcessor.EXPECT().PostProcess(mock.Anything, mock.Anything).Return(true, nil).Once()

// 					accessProviderFileCreatorError = nil
// 					accessProviderFileReaderError = nil

// 					return accessProviderFileCreatorError, accessProviderFileReaderError
// 				},
// 			},
// 			args: args{
// 				ctx: context.Background(),
// 				config: &access_provider_post_processor.AccessProviderPostProcessorConfig{
// 					TagOverwriteKeyForName: "key",
// 					ConfigMap:              &config2.ConfigMap{Parameters: map[string]string{}},
// 				},
// 			},
// 			want: want{
// 				touchedAps: int32(1),
// 			},
// 			wantErr: require.NoError,
// 		},
// 		{
// 			name: "error-post-processor",
// 			fields: fields{
// 				setup: func(accessProviderFileCreator *mocks3.AccessProviderFileCreator, accessProviderFileReader *mocks3.AccessProviderSyncFromTargetFileParser, postProcessor *MockAccessProviderPostProcessorI, args args) (accessProviderFileCreatorError error, accessProviderFileReaderError error) {
// 					accessProviderFileCreator.EXPECT().Close().Return().Once()

// 					postProcessor.EXPECT().Initialize(mock.Anything, accessProviderFileCreator, args.config).Return(errors.New("BOOM!")).Once()

// 					accessProviderFileCreatorError = nil
// 					accessProviderFileReaderError = nil

// 					return accessProviderFileCreatorError, accessProviderFileReaderError
// 				},
// 			},
// 			args: args{
// 				ctx: context.Background(),
// 				config: &access_provider_post_processor.AccessProviderPostProcessorConfig{
// 					TagOverwriteKeyForName: "key",
// 					ConfigMap:              &config2.ConfigMap{Parameters: map[string]string{}},
// 				},
// 			},
// 			want: want{
// 				touchedAps: int32(0),
// 			},
// 			wantErr: require.Error,
// 		},
// 		{
// 			name: "error-on-file-creation",
// 			fields: fields{
// 				setup: func(accessProviderFileCreator *mocks3.AccessProviderFileCreator, accessProviderFileReader *mocks3.AccessProviderSyncFromTargetFileParser, postProcessor *MockAccessProviderPostProcessorI, args args) (accessProviderFileCreatorError error, accessProviderFileReaderError error) {
// 					accessProviderFileCreatorError = errors.New("BOOM!")
// 					accessProviderFileReaderError = nil

// 					return accessProviderFileCreatorError, accessProviderFileReaderError
// 				},
// 			},
// 			args: args{
// 				ctx: context.Background(),
// 				config: &access_provider_post_processor.AccessProviderPostProcessorConfig{
// 					TagOverwriteKeyForName: "key",
// 					ConfigMap:              &config2.ConfigMap{Parameters: map[string]string{}},
// 				},
// 			},
// 			want: want{
// 				touchedAps: int32(0),
// 			},
// 			wantErr: require.Error,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockFileCreator, mockFileReader, mockPostProcessor := createPostProcessor(t)
// 			accessProviderFileCreatorError, accessProviderFileReaderError := tt.fields.setup(mockFileCreator, mockFileReader, mockPostProcessor, tt.args)

// 			postProcessorFn := accessProviderPostProcessorFunction{
// 				postProcessor: NewSyncFactory(NewDummySyncFactoryFn[AccessProviderPostProcessorI](mockPostProcessor)),
// 				accessFileCreatorFactory: func(config *access_provider.AccessSyncFromTarget) (sync_from_target.AccessProviderFileCreator, error) {
// 					return mockFileCreator, accessProviderFileCreatorError
// 				},
// 				accessProviderParserFactory: func(sourceFile string) (sync_from_target.AccessProviderSyncFromTargetFileParser, error) {
// 					return mockFileReader, accessProviderFileReaderError
// 				},
// 			}

// 			result, err := postProcessorFn.PostProcessFromTarget(tt.args.ctx, tt.args.config)
// 			tt.wantErr(t, err)

// 			if err != nil {
// 				return
// 			}

// 			assert.Equal(t, tt.want.touchedAps, result.AccessProviderTouchedCount)
// 		})
// 	}
// }

// func createPostProcessor(t *testing.T) (*mocks3.AccessProviderFileCreator, *mocks3.AccessProviderSyncFromTargetFileParser, *MockAccessProviderPostProcessorI) {
// 	accessProviderFileCreator := mocks.NewAccessProviderFileCreator(t)

// 	accessProviderFileReader := mocks3.NewAccessProviderSyncFromTargetFileParser(t)

// 	postProcessor := NewMockAccessProviderPostProcessorI(t)

// 	return accessProviderFileCreator, accessProviderFileReader, postProcessor
// }
