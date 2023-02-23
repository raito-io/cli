package data_source

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/raito-io/cli/base/data_source"
)

func TestFixMetaData(t *testing.T) {
	input := map[string]data_source.MetaData{
		"{dataObjectTypes:[{name:\"datasource\",type:\"datasource\",permissions:[{permission:\"APPLY MASKING POLICY\"}]}],supportedFeatures:[\"columnMasking\"]}": {
			SupportedFeatures: []string{data_source.ColumnMasking},
			DataObjectTypes: []*data_source.DataObjectType{
				{
					Name: data_source.Datasource,
					Type: data_source.Datasource,
					Permissions: []*data_source.DataObjectTypePermission{
						{
							Permission: "APPLY MASKING POLICY",
						},
					},
					Children: []string{},
				},
			},
		},

		"{dataObjectTypes:[{name:\"datasource\",type:\"datasource\",permissions:[{permission:\"SELECT\",description:\"test\"}]}],supportedFeatures:[\"columnMasking\"]}": {
			SupportedFeatures: []string{data_source.ColumnMasking},
			DataObjectTypes: []*data_source.DataObjectType{
				{
					Name: data_source.Datasource,
					Type: data_source.Datasource,
					Permissions: []*data_source.DataObjectTypePermission{
						{
							Permission:        "SELECT",
							GlobalPermissions: []string{},
							Description:       "test",
						},
					},
				},
			},
		},

		"{dataObjectTypes:[{name:\"datasource\",type:\"datasource\"}],supportedFeatures:[\"columnMasking\"],type:\"snowflake\",icon:\"sf-icon\"}": {
			SupportedFeatures: []string{data_source.ColumnMasking},
			DataObjectTypes: []*data_source.DataObjectType{
				{
					Name:        data_source.Datasource,
					Type:        data_source.Datasource,
					Permissions: []*data_source.DataObjectTypePermission{},
				},
			},
			Type: "snowflake",
			Icon: "sf-icon",
		},
		"{dataObjectTypes:[{name:\"datasource\",type:\"datasource\"}],supportedFeatures:[\"columnFiltering\"],type:\"snowflake\"}": {
			SupportedFeatures: []string{data_source.ColumnFiltering},
			DataObjectTypes: []*data_source.DataObjectType{
				{
					Name: data_source.Datasource,
					Type: data_source.Datasource,
				},
			},
			Type: "snowflake",
		},
	}

	for expected, i := range input {
		mds, err := marshalMetaData(&i)
		assert.NoError(t, err)
		md := fixMetaData(mds)
		assert.Equal(t, expected, md)
	}
}
