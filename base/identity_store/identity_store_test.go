package identity_store

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

func TestIdentityStoreFileCreator_Users(t *testing.T) {
	tempFile1, _ := os.Create("tempfile-" + strconv.Itoa(rand.Int()) + ".json")
	defer os.Remove(tempFile1.Name())
	tempFile2, _ := os.Create("tempfile-" + strconv.Itoa(rand.Int()) + ".json")
	defer os.Remove(tempFile2.Name())
	config := IdentityStoreSyncConfig{
		UserFile:  tempFile1.Name(),
		GroupFile: tempFile2.Name(),
	}
	isfc, err := NewIdentityStoreFileCreator(&config)
	assert.Nil(t, err)
	assert.NotNil(t, isfc)

	users := make([]*User, 0, 3)
	users = append(users, &User{
		ExternalId:       "ueid1",
		Name:             "User 1",
		UserName:         "u1",
		Email:            "u1@raito.io",
		GroupExternalIds: []string{"geid1"},
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
	users = append(users, &User{
		ExternalId:       "ueid2",
		Name:             "User 2",
		UserName:         "u2",
		Email:            "u2@raito.io",
		GroupExternalIds: []string{"geid1", "geid2"},
		Tags: []*tag.Tag{
			{
				Key:    "k3",
				Value:  "v3",
				Source: "test",
			},
		},
	})
	users = append(users, &User{
		ExternalId: "ueid3",
		Name:       "User 3",
		UserName:   "u3",
	})

	err = isfc.AddUsers(users...)
	assert.Nil(t, err)

	groups := make([]*Group, 0, 2)
	groups = append(groups, &Group{
		ExternalId:             "geid1",
		Name:                   "g1",
		DisplayName:            "Group1",
		Description:            "A group",
		ParentGroupExternalIds: []string{"geid2"},
		Tags: []*tag.Tag{
			{
				Key:    "gk1",
				Value:  "gv1",
				Source: "test",
			},
			{
				Key:    "gk2",
				Value:  "gv2",
				Source: "test",
			},
		},
	})
	groups = append(groups, &Group{
		ExternalId:  "geid2",
		Name:        "g2",
		DisplayName: "Group2",
		Tags: []*tag.Tag{
			{
				Key:    "gk3",
				Value:  "gv3",
				Source: "test",
			},
		},
	})

	err = isfc.AddGroups(groups...)
	assert.Nil(t, err)

	isfc.Close()

	assert.Equal(t, 3, isfc.GetUserCount())
	assert.Equal(t, 2, isfc.GetGroupCount())

	bytes, err := ioutil.ReadAll(tempFile1)
	assert.Nil(t, err)

	ur := make([]User, 0, 3)
	err = json.Unmarshal(bytes, &ur)
	assert.Nil(t, err)

	assert.Equal(t, 3, len(ur))
	assert.Equal(t, "ueid1", ur[0].ExternalId)
	assert.Equal(t, "User 1", ur[0].Name)
	assert.Equal(t, "u1", ur[0].UserName)
	assert.Equal(t, "u1@raito.io", ur[0].Email)
	assert.Equal(t, 2, len(ur[0].Tags))
	assert.Equal(t, "v1", ur[0].Tags[0].Value)
	assert.Equal(t, "v2", ur[0].Tags[1].Value)
	assert.Equal(t, 1, len(ur[0].GroupExternalIds))
	assert.Equal(t, "geid1", ur[0].GroupExternalIds[0])

	assert.Equal(t, "ueid2", ur[1].ExternalId)
	assert.Equal(t, "User 2", ur[1].Name)
	assert.Equal(t, "u2", ur[1].UserName)
	assert.Equal(t, "u2@raito.io", ur[1].Email)
	assert.Equal(t, 1, len(ur[1].Tags))
	assert.Equal(t, "v3", ur[1].Tags[0].Value)
	assert.Equal(t, 2, len(ur[1].GroupExternalIds))
	assert.Equal(t, "geid1", ur[1].GroupExternalIds[0])
	assert.Equal(t, "geid2", ur[1].GroupExternalIds[1])

	assert.Equal(t, "ueid3", ur[2].ExternalId)
	assert.Equal(t, "User 3", ur[2].Name)
	assert.Equal(t, "u3", ur[2].UserName)
	assert.Equal(t, "", ur[2].Email)
	assert.Equal(t, 0, len(ur[2].Tags))
	assert.Nil(t, ur[2].GroupExternalIds)

	bytes, err = ioutil.ReadAll(tempFile2)
	assert.Nil(t, err)

	gr := make([]Group, 0, 3)
	err = json.Unmarshal(bytes, &gr)
	assert.Nil(t, err)

	assert.Equal(t, 3, len(ur))
	assert.Equal(t, "geid1", gr[0].ExternalId)
	assert.Equal(t, "g1", gr[0].Name)
	assert.Equal(t, "Group1", gr[0].DisplayName)
	assert.Equal(t, "A group", gr[0].Description)
	assert.Equal(t, 2, len(gr[0].Tags))
	assert.Equal(t, "gv1", gr[0].Tags[0].Value)
	assert.Equal(t, "gv2", gr[0].Tags[1].Value)
	assert.Equal(t, 1, len(gr[0].ParentGroupExternalIds))
	assert.Equal(t, "geid2", gr[0].ParentGroupExternalIds[0])

	assert.Equal(t, "geid2", gr[1].ExternalId)
	assert.Equal(t, "g2", gr[1].Name)
	assert.Equal(t, "Group2", gr[1].DisplayName)
	assert.Equal(t, "", gr[1].Description)
	assert.Equal(t, 1, len(gr[1].Tags))
	assert.Equal(t, "gv3", gr[1].Tags[0].Value)
	assert.Nil(t, gr[1].ParentGroupExternalIds)
}
