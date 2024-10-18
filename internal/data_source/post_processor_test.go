package data_source

import (
	"fmt"
	"testing"

	"github.com/aws/smithy-go/ptr"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"

	"github.com/raito-io/cli/base/constants"
	"github.com/raito-io/cli/base/data_source"
	mocks2 "github.com/raito-io/cli/base/data_source/mocks"
	"github.com/raito-io/cli/base/tag"
)

func TestPostProcessor_postProcessDataObject(t *testing.T) {
	type fields struct {
		dataSourceFileCreatorFactory func(config *data_source.DataSourceSyncConfig) (data_source.DataSourceFileCreator, error)
		config                       *PostProcessorConfig
	}
	type args struct {
		do *data_source.DataObject
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		want             bool
		expectedOutputDo *data_source.DataObject
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "no owners",
			fields: fields{
				dataSourceFileCreatorFactory: nil,
				config: &PostProcessorConfig{
					TagOverwriteKeyForOwners: "overwriteKey",
					DataSourceId:             "ds1",
					DataObjectParent:         "",
					DataObjectExcludes:       nil,
					TargetLogger:             hclog.NewNullLogger(),
				},
			},
			args: args{
				do: &data_source.DataObject{
					ExternalId:       "externalId1",
					Name:             "someDataObject",
					FullName:         "ds1.schema1.someDataObject",
					Type:             "table",
					Description:      "A simple data object",
					ParentExternalId: "ds1.schema1",
					Tags:             nil,
					DataType:         ptr.String("table"),
					Owners:           nil,
				},
			},
			want: false,
			expectedOutputDo: &data_source.DataObject{
				ExternalId:       "externalId1",
				Name:             "someDataObject",
				FullName:         "ds1.schema1.someDataObject",
				Type:             "table",
				Description:      "A simple data object",
				ParentExternalId: "ds1.schema1",
				Tags:             nil,
				DataType:         ptr.String("table"),
				Owners:           nil,
			},
			wantErr: assert.NoError,
		},
		{
			name: "Owners by raitoOwnerTag only",
			fields: fields{
				dataSourceFileCreatorFactory: nil,
				config: &PostProcessorConfig{
					TagOverwriteKeyForOwners: "overwriteKey",
					DataSourceId:             "ds1",
					DataObjectParent:         "",
					DataObjectExcludes:       nil,
					TargetLogger:             hclog.NewNullLogger(),
				},
			},
			args: args{
				do: &data_source.DataObject{
					ExternalId:       "externalId1",
					Name:             "someDataObject",
					FullName:         "ds1.schema1.someDataObject",
					Type:             "table",
					Description:      "A simple data object",
					ParentExternalId: "ds1.schema1",
					Tags: []*tag.Tag{
						{
							Key:    "tag1",
							Value:  "value1",
							Source: "someSource",
						},
						{
							Key:    constants.RaitoOwnerTagKey,
							Value:  "email:user@raito.io, otherOwner",
							Source: "someSource",
						},
					},
					DataType: ptr.String("table"),
					Owners:   nil,
				},
			},
			want: true,
			expectedOutputDo: &data_source.DataObject{
				ExternalId:       "externalId1",
				Name:             "someDataObject",
				FullName:         "ds1.schema1.someDataObject",
				Type:             "table",
				Description:      "A simple data object",
				ParentExternalId: "ds1.schema1",
				Tags: []*tag.Tag{
					{
						Key:    "tag1",
						Value:  "value1",
						Source: "someSource",
					},
					{
						Key:    constants.RaitoOwnerTagKey,
						Value:  "email:user@raito.io,otherOwner",
						Source: "someSource",
					},
				},
				DataType: ptr.String("table"),
				Owners:   nil,
			},
			wantErr: assert.NoError,
		},
		{
			name: "Owners by specified tag only",
			fields: fields{
				dataSourceFileCreatorFactory: nil,
				config: &PostProcessorConfig{
					TagOverwriteKeyForOwners: "overwriteKey",
					DataSourceId:             "ds1",
					DataObjectParent:         "",
					DataObjectExcludes:       nil,
					TargetLogger:             hclog.NewNullLogger(),
				},
			},
			args: args{
				do: &data_source.DataObject{
					ExternalId:       "externalId1",
					Name:             "someDataObject",
					FullName:         "ds1.schema1.someDataObject",
					Type:             "table",
					Description:      "A simple data object",
					ParentExternalId: "ds1.schema1",
					Tags: []*tag.Tag{
						{
							Key:    "tag1",
							Value:  "value1",
							Source: "someSource",
						},
						{
							Key:    "overwriteKey",
							Value:  "email:user@raito.io, otherOwner",
							Source: "someSource",
						},
					},
					DataType: ptr.String("table"),
					Owners:   nil,
				},
			},
			want: true,
			expectedOutputDo: &data_source.DataObject{
				ExternalId:       "externalId1",
				Name:             "someDataObject",
				FullName:         "ds1.schema1.someDataObject",
				Type:             "table",
				Description:      "A simple data object",
				ParentExternalId: "ds1.schema1",
				Tags: []*tag.Tag{
					{
						Key:    "tag1",
						Value:  "value1",
						Source: "someSource",
					},
					{
						Key:    "overwriteKey",
						Value:  "email:user@raito.io, otherOwner",
						Source: "someSource",
					},
					{
						Key:    constants.RaitoOwnerTagKey,
						Value:  "email:user@raito.io,otherOwner",
						Source: "someSource",
					},
				},
				DataType: ptr.String("table"),
				Owners:   nil,
			},
			wantErr: assert.NoError,
		},
		{
			name: "Owners by specified tag and raitoOwnerTag",
			fields: fields{
				dataSourceFileCreatorFactory: nil,
				config: &PostProcessorConfig{
					TagOverwriteKeyForOwners: "overwriteKey",
					DataSourceId:             "ds1",
					DataObjectParent:         "",
					DataObjectExcludes:       nil,
					TargetLogger:             hclog.NewNullLogger(),
				},
			},
			args: args{
				do: &data_source.DataObject{
					ExternalId:       "externalId1",
					Name:             "someDataObject",
					FullName:         "ds1.schema1.someDataObject",
					Type:             "table",
					Description:      "A simple data object",
					ParentExternalId: "ds1.schema1",
					Tags: []*tag.Tag{
						{
							Key:    "tag1",
							Value:  "value1",
							Source: "someSource",
						},
						{
							Key:    "overwriteKey",
							Value:  "email:user@raito.io, otherOwner",
							Source: "someSource",
						},
						{
							Key:    constants.RaitoOwnerTagKey,
							Value:  "email:user2@raito.io, andAnotherOwner",
							Source: "someSource",
						},
					},
					DataType: ptr.String("table"),
					Owners:   nil,
				},
			},
			want: true,
			expectedOutputDo: &data_source.DataObject{
				ExternalId:       "externalId1",
				Name:             "someDataObject",
				FullName:         "ds1.schema1.someDataObject",
				Type:             "table",
				Description:      "A simple data object",
				ParentExternalId: "ds1.schema1",
				Tags: []*tag.Tag{
					{
						Key:    "tag1",
						Value:  "value1",
						Source: "someSource",
					},
					{
						Key:    "overwriteKey",
						Value:  "email:user@raito.io, otherOwner",
						Source: "someSource",
					},
					{
						Key:    constants.RaitoOwnerTagKey,
						Value:  "email:user2@raito.io,andAnotherOwner,email:user@raito.io,otherOwner",
						Source: "someSource",
					},
				},
				DataType: ptr.String("table"),
				Owners:   nil,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PostProcessor{
				dataSourceFileCreatorFactory: tt.fields.dataSourceFileCreatorFactory,
				config:                       tt.fields.config,
			}

			outputWriter := mocks2.NewDataSourceFileCreator(t)

			if tt.expectedOutputDo != nil {
				outputWriter.EXPECT().AddDataObjects(tt.expectedOutputDo).Return(nil)
			} else {
				outputWriter.AssertNotCalled(t, "AddDataObjects")
			}

			got, err := p.postProcessDataObject(tt.args.do, outputWriter)
			if !tt.wantErr(t, err, fmt.Sprintf("postProcessDataObject(%+v)", tt.args.do)) {
				return
			}
			assert.Equalf(t, tt.want, got, "postProcessDataObject(%+v)", tt.args.do)
		})
	}
}
