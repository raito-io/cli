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
	do := data_source.DataObjectReference{
		FullName: "Data Object 1",
		Type:     "schema",
	}

	aps = append(aps, AccessProvider{
		ExternalId:        "eid1",
		NonInternalizable: false,
		Name:              "AP1",
		Users:             []string{"uid1"},
		Groups:            []string{"gid1"},
		AccessObjects:     []Access{{Permissions: []string{"A", "B"}, DataObjectReference: &do}},
	})

	aps = append(aps, AccessProvider{
		ExternalId:        "eid2",
		NonInternalizable: true,
		Name:              "AP2",
		Users:             []string{"uid1", "uid2"},
		Groups:            []string{"gid1", "gid2"},
		AccessObjects:     []Access{{Permissions: []string{"C"}, DataObjectReference: &do}},
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
	assert.Equal(t, "Data Object 1", apsr[0].AccessObjects[0].DataObjectReference.FullName)
	assert.Equal(t, "schema", apsr[0].AccessObjects[0].DataObjectReference.Type)

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
	assert.Equal(t, "Data Object 1", apsr[1].AccessObjects[0].DataObjectReference.FullName)
	assert.Equal(t, "schema", apsr[1].AccessObjects[0].DataObjectReference.Type)
}
