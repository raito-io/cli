package data_source

import (
	"testing"

	"github.com/raito-io/cli/base/data_source"
	"github.com/stretchr/testify/assert"
)

func TestMetaData_InputEmpty(t *testing.T) {

	input := data_source.MetaData{}

	err := metadataConsistencyCheck(&input)
	assert.Nil(t, err)
}

func TestMetaData_DefaultLevelEmpty(t *testing.T) {

	input := data_source.MetaData{
		UsageMetaInfo: &data_source.UsageMetaInput{
			DefaultLevel: "",
		},
	}

	err := metadataConsistencyCheck(&input)
	assert.Nil(t, err)
}

func TestMetaData_DefaultLevelNotDefined(t *testing.T) {

	input := data_source.MetaData{
		UsageMetaInfo: &data_source.UsageMetaInput{
			DefaultLevel: "",
			Levels: []*data_source.UsageMetaInputDetail{
				{
					Name:            "not_test",
					DataObjectTypes: []string{"something"},
				},
			},
		},
	}

	err := metadataConsistencyCheck(&input)
	assert.NotNil(t, err)

	input = data_source.MetaData{
		UsageMetaInfo: &data_source.UsageMetaInput{
			DefaultLevel: "test",
			Levels: []*data_source.UsageMetaInputDetail{
				{
					Name:            "not_test",
					DataObjectTypes: []string{"something"},
				},
			},
		},
	}

	err = metadataConsistencyCheck(&input)
	assert.NotNil(t, err)
}

func TestMetaData_LevelsNotDefinedInDataObjectTypes(t *testing.T) {

	input := data_source.MetaData{
		UsageMetaInfo: &data_source.UsageMetaInput{
			DefaultLevel: "table",
			Levels: []*data_source.UsageMetaInputDetail{
				{
					Name:            "table",
					DataObjectTypes: []string{"table"},
				},
			},
		},
	}

	err := metadataConsistencyCheck(&input)
	assert.NotNil(t, err)

	input = data_source.MetaData{
		DataObjectTypes: []*data_source.DataObjectType{
			{
				Name: "table",
			},
		},
		UsageMetaInfo: &data_source.UsageMetaInput{
			DefaultLevel: "table",
			Levels: []*data_source.UsageMetaInputDetail{
				{
					Name:            "table",
					DataObjectTypes: []string{"table"},
				},
			},
		},
	}

	err = metadataConsistencyCheck(&input)
	assert.Nil(t, err)

	input = data_source.MetaData{
		DataObjectTypes: []*data_source.DataObjectType{
			{
				Name: "table",
			},
		},
		UsageMetaInfo: &data_source.UsageMetaInput{
			DefaultLevel: "table",
			Levels: []*data_source.UsageMetaInputDetail{
				{
					Name:            "table",
					DataObjectTypes: []string{"table", "view"},
				},
			},
		},
	}

	err = metadataConsistencyCheck(&input)
	assert.NotNil(t, err)
}
