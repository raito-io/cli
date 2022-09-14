package exporter

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/raito-io/cli/base/access_provider"
	"github.com/raito-io/cli/base/data_source"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestAccessProviderFileCreator(t *testing.T) {
	tempFile, _ := os.Create("tempfile-" + strconv.Itoa(rand.Int()) + ".json")
	defer os.Remove(tempFile.Name())
	config := access_provider.AccessSyncFromTarget{
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
		NotInternalizable: false,
		Name:              "AP1",
		NamingHint:        "AP1Hint",
		Access: []*Access{
			{
				Who: &WhoItem{
					Users:  []string{"uid1"},
					Groups: []string{"gid1"},
				},
				What: []WhatItem{
					{
						Permissions: []string{"A", "B"},
						DataObject:  &do,
					},
				},
			},
		},
	})

	aps = append(aps, AccessProvider{
		ExternalId:        "eid2",
		NotInternalizable: true,
		Name:              "AP2",
		NamingHint:        "AP2Hint",
		Action:            Filtered,
		Access: []*Access{
			{
				Who: &WhoItem{
					Users:  []string{"uid1", "uid2"},
					Groups: []string{"gid1", "gid2"},
				},
				What: []WhatItem{
					{
						Permissions: []string{"C"},
						DataObject:  &do,
					},
				},
			},
		},
	})

	err = apfc.AddAccessProviders(aps)
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
	assert.False(t, apsr[0].NotInternalizable)
	assert.Equal(t, "AP1", apsr[0].Name)
	assert.Equal(t, "AP1Hint", apsr[0].NamingHint)
	assert.Equal(t, Promise, apsr[0].Action)
	assert.Equal(t, 1, len(apsr[0].Access[0].Who.Users))
	assert.Equal(t, "uid1", apsr[0].Access[0].Who.Users[0])
	assert.Equal(t, 1, len(apsr[0].Access[0].Who.Groups))
	assert.Equal(t, "gid1", apsr[0].Access[0].Who.Groups[0])
	assert.Equal(t, 1, len(apsr[0].Access[0].What))
	assert.Equal(t, 2, len(apsr[0].Access[0].What[0].Permissions))
	assert.Equal(t, "A", apsr[0].Access[0].What[0].Permissions[0])
	assert.Equal(t, "B", apsr[0].Access[0].What[0].Permissions[1])
	assert.Equal(t, "Data Object 1", apsr[0].Access[0].What[0].DataObject.FullName)
	assert.Equal(t, "schema", apsr[0].Access[0].What[0].DataObject.Type)

	assert.Equal(t, "eid2", apsr[1].ExternalId)
	assert.True(t, apsr[1].NotInternalizable)
	assert.Equal(t, "AP2", apsr[1].Name)
	assert.Equal(t, "AP2Hint", apsr[1].NamingHint)
	assert.Equal(t, Filtered, apsr[1].Action)
	assert.Equal(t, 2, len(apsr[1].Access[0].Who.Users))
	assert.Equal(t, "uid1", apsr[1].Access[0].Who.Users[0])
	assert.Equal(t, "uid2", apsr[1].Access[0].Who.Users[1])
	assert.Equal(t, 2, len(apsr[1].Access[0].Who.Groups))
	assert.Equal(t, "gid1", apsr[1].Access[0].Who.Groups[0])
	assert.Equal(t, "gid2", apsr[1].Access[0].Who.Groups[1])
	assert.Equal(t, 1, len(apsr[1].Access[0].What))
	assert.Equal(t, 1, len(apsr[1].Access[0].What[0].Permissions))
	assert.Equal(t, "C", apsr[1].Access[0].What[0].Permissions[0])
	assert.Equal(t, "Data Object 1", apsr[1].Access[0].What[0].DataObject.FullName)
	assert.Equal(t, "schema", apsr[1].Access[0].What[0].DataObject.Type)
}
