package access_provider

// import (
// 	"context"
// 	"testing"
//
//

// 	"github.com/aws/smithy-go/ptr"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"

// 	"github.com/raito-io/cli/base/access_provider/sync_from_target"
// 	"github.com/raito-io/cli/base/access_provider_post_processor"
// 	"github.com/raito-io/cli/base/tag"
// 	"github.com/raito-io/cli/base/wrappers/mocks"
// )

// func TestAccessProviderPostProcessorBase_PostProcess(t *testing.T) {
// 	type fields struct {
// 		setup func(accessProviderPostProcessorHandler *mocks.AccessProviderPostProcessorHandler, wantAp *sync_from_target.AccessProvider)
// 	}
// 	type args struct {
// 		ctx            context.Context
// 		config         *access_provider_post_processor.AccessProviderPostProcessorConfig
// 		accessProvider *sync_from_target.AccessProvider
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    bool
// 		wantAp  *sync_from_target.AccessProvider
// 		wantErr require.ErrorAssertionFunc
// 	}{
// 		{
// 			name: "with owner and name overrides",
// 			fields: fields{
// 				setup: func(accessProviderPostProcessorHandler *mocks.AccessProviderPostProcessorHandler, wantAp *sync_from_target.AccessProvider) {
// 					accessProviderPostProcessorHandler.EXPECT().AddAccessProviders(wantAp).Return(nil).Once()
// 				},
// 			},
// 			args: args{
// 				ctx: context.Background(),
// 				config: &access_provider_post_processor.AccessProviderPostProcessorConfig{
// 					TagOverwriteKeyForName:   "name",
// 					TagOverwriteKeyForOwners: "owner",
// 				},
// 				accessProvider: &sync_from_target.AccessProvider{
// 					Name: "TEST",
// 					Tags: []*tag.Tag{
// 						{Key: "name", Value: "new_value"},
// 						{Key: "owner", Value: "test@test.be,test@test.com"},
// 					},
// 				},
// 			},
// 			want: true,
// 			wantAp: &sync_from_target.AccessProvider{
// 				Name:             "new_value",
// 				NameLocked:       ptr.Bool(true),
// 				NameLockedReason: ptr.String("This Snowflake role cannot be renamed because it has a name tag override attached to it"),
// 				Owner: &sync_from_target.OwnerInput{
// 					Users: []string{"test@test.be", "test@test.com"},
// 				},
// 				Tags: []*tag.Tag{
// 					{Key: "name", Value: "new_value"},
// 					{Key: "owner", Value: "test@test.be,test@test.com"},
// 				},
// 			},
// 			wantErr: require.NoError,
// 		},
// 		{
// 			name: "with owner override",
// 			fields: fields{
// 				setup: func(accessProviderPostProcessorHandler *mocks.AccessProviderPostProcessorHandler, wantAp *sync_from_target.AccessProvider) {
// 					accessProviderPostProcessorHandler.EXPECT().AddAccessProviders(wantAp).Return(nil).Once()
// 				},
// 			},
// 			args: args{
// 				ctx: context.Background(),
// 				config: &access_provider_post_processor.AccessProviderPostProcessorConfig{
// 					TagOverwriteKeyForName:   "",
// 					TagOverwriteKeyForOwners: "owner",
// 				},
// 				accessProvider: &sync_from_target.AccessProvider{
// 					Name: "TEST",
// 					Tags: []*tag.Tag{
// 						{Key: "name", Value: "new_value"},
// 						{Key: "owner", Value: "test@test.be,test@test.com"},
// 					},
// 				},
// 			},
// 			want: true,
// 			wantAp: &sync_from_target.AccessProvider{
// 				Name: "TEST",
// 				Owner: &sync_from_target.OwnerInput{
// 					Users: []string{"test@test.be", "test@test.com"},
// 				},
// 				Tags: []*tag.Tag{
// 					{Key: "name", Value: "new_value"},
// 					{Key: "owner", Value: "test@test.be,test@test.com"},
// 				},
// 			},
// 			wantErr: require.NoError,
// 		},
// 		{
// 			name: "no overrides",
// 			fields: fields{
// 				setup: func(accessProviderPostProcessorHandler *mocks.AccessProviderPostProcessorHandler, wantAp *sync_from_target.AccessProvider) {
// 					accessProviderPostProcessorHandler.EXPECT().AddAccessProviders(wantAp).Return(nil).Once()
// 				},
// 			},
// 			args: args{
// 				ctx: context.Background(),
// 				config: &access_provider_post_processor.AccessProviderPostProcessorConfig{
// 					TagOverwriteKeyForName:   "",
// 					TagOverwriteKeyForOwners: "",
// 				},
// 				accessProvider: &sync_from_target.AccessProvider{
// 					Name: "TEST",
// 					Tags: []*tag.Tag{
// 						{Key: "name", Value: "new_value"},
// 						{Key: "owner", Value: "test@test.be,test@test.com"},
// 					},
// 				},
// 			},
// 			want: false,
// 			wantAp: &sync_from_target.AccessProvider{
// 				Name: "TEST",
// 				Tags: []*tag.Tag{
// 					{Key: "name", Value: "new_value"},
// 					{Key: "owner", Value: "test@test.be,test@test.com"},
// 				},
// 			},
// 			wantErr: require.NoError,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Given
// 			accessProviderPostProcessorHandler := mocks.NewAccessProviderPostProcessorHandler(t)
// 			tt.fields.setup(accessProviderPostProcessorHandler, tt.wantAp)

// 			postProcessor := NewAccessProviderPostProcessorGeneral()
// 			postProcessor.Initialize(tt.args.ctx, accessProviderPostProcessorHandler, tt.args.config)

// 			//When
// 			touched, err := postProcessor.PostProcess(tt.args.ctx, tt.args.accessProvider)

// 			// Then
// 			tt.wantErr(t, err)
// 			assert.Equal(t, tt.want, touched)
// 		})
// 	}
// }
