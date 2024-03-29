package data_source

import (
	"encoding/json"
	"github.com/raito-io/cli/base/tag"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestDataSourceFileCreator(t *testing.T) {
	tempFile, _ := os.Create("tempfile-" + strconv.Itoa(rand.Int()) + ".json")
	defer os.Remove(tempFile.Name())
	config := DataSourceSyncConfig{
		TargetFile:   tempFile.Name(),
		DataSourceId: "myDataSource",
	}
	dsfc, err := NewDataSourceFileCreator(&config)
	assert.Nil(t, err)
	assert.NotNil(t, dsfc)

	dos := make([]*DataObject, 0, 3)
	dos = append(dos, &DataObject{
		ExternalId:  "eid1",
		Name:        "DO1",
		FullName:    "Data Object 1",
		Type:        "table",
		Description: "The first data object",
		Tags: []*tag.Tag{
			{
				Key:    "k1",
				Value:  "v1",
				Source: "test",
			},
			{
				Key:    "k2",
				Value:  "v2",
				Source: "test",
			},
		},
	})
	dos = append(dos, &DataObject{
		ExternalId:       "eid2",
		Name:             "DO2",
		FullName:         "Data Object 2",
		Type:             "schema",
		Description:      "The second data object",
		ParentExternalId: "eid1",
		Tags: []*tag.Tag{
			{
				Key:    "k3",
				Value:  "v3",
				Source: "test",
			},
		},
	})
	dos = append(dos, &DataObject{
		ExternalId:       "eid3",
		Name:             "DO3",
		FullName:         "Data Object 3",
		Type:             "database",
		ParentExternalId: "eid2",
	})

	err = dsfc.AddDataObjects(dos...)
	assert.Nil(t, err)

	dsfc.Close()

	assert.Equal(t, 4, dsfc.GetDataObjectCount())

	bytes, err := ioutil.ReadAll(tempFile)
	assert.Nil(t, err)

	dosr := make([]DataObject, 0, 4)
	err = json.Unmarshal(bytes, &dosr)
	assert.Nil(t, err)

	assert.Equal(t, 4, len(dosr))
	assert.Equal(t, config.DataSourceId, dosr[0].ExternalId)
	assert.Equal(t, config.DataSourceId, dosr[0].Name)
	assert.Equal(t, config.DataSourceId, dosr[0].FullName)
	assert.Equal(t, "datasource", dosr[0].Type)
	assert.Equal(t, "", dosr[0].Description)
	assert.Equal(t, "", dosr[0].ParentExternalId)
	assert.Equal(t, 0, len(dosr[0].Tags))

	assert.Equal(t, "eid1", dosr[1].ExternalId)
	assert.Equal(t, "DO1", dosr[1].Name)
	assert.Equal(t, "Data Object 1", dosr[1].FullName)
	assert.Equal(t, "table", dosr[1].Type)
	assert.Equal(t, "The first data object", dosr[1].Description)
	assert.Equal(t, config.DataSourceId, dosr[1].ParentExternalId)
	assert.Equal(t, 2, len(dosr[1].Tags))
	assert.Equal(t, "v1", dosr[1].Tags[0].Value)
	assert.Equal(t, "v2", dosr[1].Tags[1].Value)

	assert.Equal(t, "eid2", dosr[2].ExternalId)
	assert.Equal(t, "DO2", dosr[2].Name)
	assert.Equal(t, "Data Object 2", dosr[2].FullName)
	assert.Equal(t, "schema", dosr[2].Type)
	assert.Equal(t, "The second data object", dosr[2].Description)
	assert.Equal(t, "eid1", dosr[2].ParentExternalId)
	assert.Equal(t, 1, len(dosr[2].Tags))
	assert.Equal(t, "v3", dosr[2].Tags[0].Value)

	assert.Equal(t, "eid3", dosr[3].ExternalId)
	assert.Equal(t, "DO3", dosr[3].Name)
	assert.Equal(t, "Data Object 3", dosr[3].FullName)
	assert.Equal(t, "database", dosr[3].Type)
	assert.Equal(t, "", dosr[3].Description)
	assert.Equal(t, "eid2", dosr[3].ParentExternalId)
	assert.Equal(t, 0, len(dosr[3].Tags))
}

func TestDataSourceFileCreatorPartial(t *testing.T) {
	tempFile, _ := os.Create("tempfile-" + strconv.Itoa(rand.Int()) + ".json")
	defer os.Remove(tempFile.Name())
	config := DataSourceSyncConfig{
		TargetFile:       tempFile.Name(),
		DataSourceId:     "myDataSource",
		DataObjectParent: "myParent",
	}
	dsfc, err := NewDataSourceFileCreator(&config)
	assert.Nil(t, err)
	assert.NotNil(t, dsfc)

	dos := make([]*DataObject, 0, 3)
	dos = append(dos, &DataObject{
		ExternalId:  "eid1",
		Name:        "DO1",
		FullName:    "Data Object 1",
		Type:        "table",
		Description: "The first data object",
	})
	dos = append(dos, &DataObject{
		ExternalId:       "eid2",
		Name:             "DO2",
		FullName:         "Data Object 2",
		Type:             "schema",
		Description:      "The second data object",
		ParentExternalId: "eid1",
	})
	dos = append(dos, &DataObject{
		ExternalId:       "eid3",
		Name:             "DO3",
		FullName:         "Data Object 3",
		Type:             "database",
		ParentExternalId: "eid2",
	})

	err = dsfc.AddDataObjects(dos...)
	assert.Nil(t, err)

	dsfc.Close()

	// The data source node should not have been added as this is a partial sync
	assert.Equal(t, 3, dsfc.GetDataObjectCount())
}

func TestDataSourceDetails(t *testing.T) {
	tempFile, _ := os.Create("tempfile-" + strconv.Itoa(rand.Int()) + ".json")
	defer os.Remove(tempFile.Name())
	config := DataSourceSyncConfig{
		TargetFile:   tempFile.Name(),
		DataSourceId: "myDataSource",
	}
	dsfc, err := NewDataSourceFileCreator(&config)
	assert.Nil(t, err)
	assert.NotNil(t, dsfc)

	dos := make([]*DataObject, 0, 1)
	dos = append(dos, &DataObject{
		ExternalId:  "eid1",
		Name:        "DO1",
		FullName:    "Data Object 1",
		Type:        "table",
		Description: "The first data object",
		Tags: []*tag.Tag{
			{
				Key:    "k1",
				Value:  "v1",
				Source: "test",
			},
			{
				Key:    "k2",
				Value:  "v2",
				Source: "test",
			},
		},
	})

	// set data source specs
	dsfc.GetDataSourceDetails().SetName("dsName")
	dsfc.GetDataSourceDetails().SetFullname("dsFullName")
	dsfc.GetDataSourceDetails().SetDescription("dsDesc")

	err = dsfc.AddDataObjects(dos...)
	assert.Nil(t, err)

	dsfc.Close()

	assert.Equal(t, 2, dsfc.GetDataObjectCount())

	bytes, err := ioutil.ReadAll(tempFile)
	assert.Nil(t, err)

	dosr := make([]DataObject, 0, 2)
	err = json.Unmarshal(bytes, &dosr)
	assert.Nil(t, err)

	assert.Equal(t, 2, len(dosr))
	assert.Equal(t, config.DataSourceId, dosr[0].ExternalId)
	assert.Equal(t, "dsName", dosr[0].Name)
	assert.Equal(t, "dsFullName", dosr[0].FullName)
	assert.Equal(t, "datasource", dosr[0].Type)
	assert.Equal(t, "dsDesc", dosr[0].Description)
	assert.Equal(t, "", dosr[0].ParentExternalId)
	assert.Equal(t, 0, len(dosr[0].Tags))
}
