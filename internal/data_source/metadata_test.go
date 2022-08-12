package data_source

import (
	"github.com/raito-io/cli/common/api/data_source"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFixMetaData(t *testing.T) {
	input := map[string]data_source.MetaData{
		"{dataObjectTypes:[{name:\"datasource\",label:\"\",icon:\"\",permissions:[{permission:\"APPLY MASKING POLICY\",description:\"\"}],children:[]}],supportedFeatures:[\"columnMasking\"],type:\"\",icon:\"\"}": {
			SupportedFeatures: []string{data_source.ColumnMasking},
			DataObjectTypes: []data_source.DataObjectType{
				{
					Name: data_source.Datasource,
					Permissions: []data_source.DataObjectTypePermission{
						{
							Permission: "APPLY MASKING POLICY",
						},
					},
					Children: []string{},
				},
			},
		},

		"{dataObjectTypes:[{name:\"datasource\",label:\"\",icon:\"\",permissions:[{permission:\"SELECT\",description:\"test\"}],children:[]}],supportedFeatures:[\"columnMasking\"],type:\"\",icon:\"\"}": {
			SupportedFeatures: []string{data_source.ColumnMasking},
			DataObjectTypes: []data_source.DataObjectType{
				{
					Name: data_source.Datasource,
					Permissions: []data_source.DataObjectTypePermission{
						{
							Permission:        "SELECT",
							GlobalPermissions: []string{},
							Description:       "test",
						},
					},
				},
			},
		},

		"{dataObjectTypes:[{name:\"datasource\",label:\"\",icon:\"\",permissions:[],children:[]}],supportedFeatures:[\"columnMasking\"],type:\"snowflake\",icon:\"sf-icon\"}": {
			SupportedFeatures: []string{data_source.ColumnMasking},
			DataObjectTypes: []data_source.DataObjectType{
				{
					Name:        data_source.Datasource,
					Permissions: []data_source.DataObjectTypePermission{},
				},
			},
			Type: "snowflake",
			Icon: "sf-icon",
		},
		"{dataObjectTypes:[{name:\"datasource\",label:\"\",icon:\"\",permissions:[],children:[]}],supportedFeatures:[\"columnFiltering\"],type:\"snowflake\",icon:\"\"}": {
			SupportedFeatures: []string{data_source.ColumnFiltering},
			DataObjectTypes: []data_source.DataObjectType{
				{
					Name: data_source.Datasource,
				},
			},
			Type: "snowflake",
		},
	}

	for expected, i := range input {
		mds, err := marshalMetaData(i)
		assert.NoError(t, err)
		md := fixMetaData(mds)
		assert.Equal(t, expected, md)
	}
}
