package data_source

import (
	"encoding/json"
	"github.com/raito-io/cli/common/api/data_source"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestDataSourceFileCreator(t *testing.T) {
	tempFile, _ := os.Create("tempfile-"+strconv.Itoa(rand.Int())+".json")
	defer os.Remove(tempFile.Name())
	config := data_source.DataSourceSyncConfig{
		TargetFile: tempFile.Name(),
	}
	dsfc, err := NewDataSourceFileCreator(&config)
	assert.Nil(t, err)
	assert.NotNil(t, dsfc)

	dos := make([]DataObject, 0, 3)
	dos = append(dos, DataObject{
		ExternalId: "eid1",
		Name: "DO1",
		FullName: "Data Object 1",
		Type: "table",
		Description: "The first data object",
		Tags: map[string]interface{} { "k1": "v1", "k2": "v2" },
	})
	dos = append(dos, DataObject{
		ExternalId: "eid2",
		Name: "DO2",
		FullName: "Data Object 2",
		Type: "schema",
		Description: "The second data object",
		ParentExternalId: "eid1",
		Tags: map[string]interface{} { "k3": "v3" },
	})
	dos = append(dos, DataObject{
		ExternalId: "eid3",
		Name: "DO3",
		FullName: "Data Object 3",
		Type: "database",
		ParentExternalId: "eid2",
	})

	err = dsfc.AddDataObjects(dos)
	assert.Nil(t, err)

	dsfc.Close()

	assert.Equal(t, 3, dsfc.GetDataObjectCount())

	bytes, err := ioutil.ReadAll(tempFile)
	assert.Nil(t, err)

	dosr := make([]DataObject, 0, 3)
	err = json.Unmarshal(bytes, &dosr)
	assert.Nil(t, err)

	assert.Equal(t, 3, len(dosr))
	assert.Equal(t, "eid1", dosr[0].ExternalId)
	assert.Equal(t, "DO1", dosr[0].Name)
	assert.Equal(t, "Data Object 1", dosr[0].FullName)
	assert.Equal(t, "table", dosr[0].Type)
	assert.Equal(t, "The first data object", dosr[0].Description)
	assert.Equal(t, "", dosr[0].ParentExternalId)
	assert.Equal(t, 2, len(dosr[0].Tags))
	assert.Equal(t, "v1", dosr[0].Tags["k1"])
	assert.Equal(t, "v2", dosr[0].Tags["k2"])

	assert.Equal(t, "eid2", dosr[1].ExternalId)
	assert.Equal(t, "DO2", dosr[1].Name)
	assert.Equal(t, "Data Object 2", dosr[1].FullName)
	assert.Equal(t, "schema", dosr[1].Type)
	assert.Equal(t, "The second data object", dosr[1].Description)
	assert.Equal(t, "eid1", dosr[1].ParentExternalId)
	assert.Equal(t, 1, len(dosr[1].Tags))
	assert.Equal(t, "v3", dosr[1].Tags["k3"])
	assert.Nil(t, dosr[1].Tags["k1"])

	assert.Equal(t, "eid3", dosr[2].ExternalId)
	assert.Equal(t, "DO3", dosr[2].Name)
	assert.Equal(t, "Data Object 3", dosr[2].FullName)
	assert.Equal(t, "database", dosr[2].Type)
	assert.Equal(t, "", dosr[2].Description)
	assert.Equal(t, "eid2", dosr[2].ParentExternalId)
	assert.Equal(t, 0, len(dosr[2].Tags))
	assert.Nil(t, dosr[2].Tags["k1"])
}
