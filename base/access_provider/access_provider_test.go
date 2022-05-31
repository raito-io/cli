package access_provider

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/raito-io/cli/base/data_source"
	"github.com/raito-io/cli/common/api/data_access"
	"github.com/stretchr/testify/assert"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestAccessProviderFileCreator(t *testing.T) {
	tempFile, _ := os.Create("tempfile-" + strconv.Itoa(rand.Int()) + ".json")
	defer os.Remove(tempFile.Name())
	config := data_access.DataAccessSyncConfig{
		TargetFile: tempFile.Name(),
	}
	apfc, err := NewAccessProviderFileCreator(&config)
	assert.Nil(t, err)
	assert.NotNil(t, apfc)

	aps := make([]AccessProvider, 0, 3)
	do := data_source.DataObject{
		ExternalId:       "eid1",
		Name:             "DO1",
		FullName:         "Data Object 1",
		Type:             "schema",
		Description:      "The first data object",
		ParentExternalId: "pid1",
		Tags:             map[string]interface{}{"k1": "v1"},
	}

	aps = append(aps, AccessProvider{
		ExternalId:        "eid1",
		NonInternalizable: false,
		Name:              "AP1",
		Users:             []string{"uid1"},
		Groups:            []string{"gid1"},
		AccessObjects:     []Access{{Permissions: []string{"A", "B"}, DataObject: &do}},
	})

	aps = append(aps, AccessProvider{
		ExternalId:        "eid2",
		NonInternalizable: true,
		Name:              "AP2",
		Users:             []string{"uid1", "uid2"},
		Groups:            []string{"gid1", "gid2"},
		AccessObjects:     []Access{{Permissions: []string{"C"}, DataObject: &do}},
	})

	err = apfc.AddAccessProvider(aps)
	assert.Nil(t, err)

	apfc.Close()

	assert.Equal(t, 2, apfc.GetAccessProviderCount())

	bytes, err := ioutil.ReadAll(tempFile)
	assert.Nil(t, err)

	apsr := make([]AccessProvider, 0, 2)
	err = json.Unmarshal(bytes, &apsr)
	assert.Nil(t, err)

	assert.Equal(t, 2, len(apsr))

	assert.Equal(t, "eid1", apsr[0].ExternalId)
	assert.False(t, apsr[0].NonInternalizable)
	assert.Equal(t, "AP1", apsr[0].Name)
	assert.Equal(t, 1, len(apsr[0].Users))
	assert.Equal(t, "uid1", apsr[0].Users[0])
	assert.Equal(t, 1, len(apsr[0].Groups))
	assert.Equal(t, "gid1", apsr[0].Groups[0])
	assert.Equal(t, 1, len(apsr[0].AccessObjects))
	assert.Equal(t, 2, len(apsr[0].AccessObjects[0].Permissions))
	assert.Equal(t, "A", apsr[0].AccessObjects[0].Permissions[0])
	assert.Equal(t, "B", apsr[0].AccessObjects[0].Permissions[1])
	assert.Equal(t, "eid1", apsr[0].AccessObjects[0].DataObject.ExternalId)
	assert.Equal(t, "DO1", apsr[0].AccessObjects[0].DataObject.Name)
	assert.Equal(t, "Data Object 1", apsr[0].AccessObjects[0].DataObject.FullName)
	assert.Equal(t, "schema", apsr[0].AccessObjects[0].DataObject.Type)
	assert.Equal(t, "The first data object", apsr[0].AccessObjects[0].DataObject.Description)
	assert.Equal(t, "pid1", apsr[0].AccessObjects[0].DataObject.ParentExternalId)
	assert.Equal(t, 1, len(apsr[0].AccessObjects[0].DataObject.Tags))
	assert.Equal(t, "v1", apsr[0].AccessObjects[0].DataObject.Tags["k1"])

	assert.Equal(t, "eid2", apsr[1].ExternalId)
	assert.True(t, apsr[1].NonInternalizable)
	assert.Equal(t, "AP2", apsr[1].Name)
	assert.Equal(t, 2, len(apsr[1].Users))
	assert.Equal(t, "uid1", apsr[1].Users[0])
	assert.Equal(t, "uid2", apsr[1].Users[1])
	assert.Equal(t, 2, len(apsr[1].Groups))
	assert.Equal(t, "gid1", apsr[1].Groups[0])
	assert.Equal(t, "gid2", apsr[1].Groups[1])
	assert.Equal(t, 1, len(apsr[1].AccessObjects))
	assert.Equal(t, 1, len(apsr[1].AccessObjects[0].Permissions))
	assert.Equal(t, "C", apsr[1].AccessObjects[0].Permissions[0])
	assert.Equal(t, "eid1", apsr[1].AccessObjects[0].DataObject.ExternalId)
	assert.Equal(t, "DO1", apsr[1].AccessObjects[0].DataObject.Name)
	assert.Equal(t, "Data Object 1", apsr[1].AccessObjects[0].DataObject.FullName)
	assert.Equal(t, "schema", apsr[1].AccessObjects[0].DataObject.Type)
	assert.Equal(t, "The first data object", apsr[1].AccessObjects[0].DataObject.Description)
	assert.Equal(t, "pid1", apsr[1].AccessObjects[0].DataObject.ParentExternalId)
	assert.Equal(t, 1, len(apsr[1].AccessObjects[0].DataObject.Tags))
	assert.Equal(t, "v1", apsr[1].AccessObjects[0].DataObject.Tags["k1"])
}
