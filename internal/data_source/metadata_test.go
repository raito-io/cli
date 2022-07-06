package data_source

import (
	"github.com/raito-io/cli/common/api/data_source"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFixMetaData(t *testing.T) {
	input := map[string]data_source.MetaData{
		"{dataObjectTypes:[{name:\"datasource\",label:\"\",icon:\"\",permissions:[{permission:\"APPLY MASKING POLICY\"}],children:[]}],supportedFeatures:[\"columnMasking\"]}": {
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

		"{dataObjectTypes:[{name:\"datasource\",label:\"\",icon:\"\",permissions:[{permission:\"SELECT\"}],children:[]}],supportedFeatures:[\"columnMasking\"]}": {
			SupportedFeatures: []string{data_source.ColumnMasking},
			DataObjectTypes: []data_source.DataObjectType{
				{
					Name: data_source.Datasource,
					Permissions: []data_source.DataObjectTypePermission{
						{
							Permission:        "SELECT",
							GlobalPermissions: []string{},
						},
					},
				},
			},
		},

		"{dataObjectTypes:[{name:\"datasource\",label:\"\",icon:\"\",permissions:[],children:[]}],supportedFeatures:[\"columnMasking\"]}": {
			SupportedFeatures: []string{data_source.ColumnMasking},
			DataObjectTypes: []data_source.DataObjectType{
				{
					Name:        data_source.Datasource,
					Permissions: []data_source.DataObjectTypePermission{},
				},
			},
		},
		"{dataObjectTypes:[{name:\"datasource\",label:\"\",icon:\"\",permissions:[],children:[]}],supportedFeatures:[\"columnFiltering\"]}": {
			SupportedFeatures: []string{data_source.ColumnFiltering},
			DataObjectTypes: []data_source.DataObjectType{
				{
					Name: data_source.Datasource,
				},
			},
		},
	}

	for expected, i := range input {
		mds, err := marshalMetaData(i)
		assert.NoError(t, err)
		md, err := fixMetaData(mds)
		assert.NoError(t, err)
		assert.Equal(t, expected, md)
	}
}
