package sync_from_target

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/smithy-go/ptr"
	"github.com/stretchr/testify/assert"

	"github.com/raito-io/cli/base/access_provider/types"
	"github.com/raito-io/cli/base/tag"

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

	aps := make([]*AccessProvider, 0, 3)
	do := data_source.DataObjectReference{
		FullName: "Data Object 1",
		Type:     "schema",
	}

	aps = append(aps, &AccessProvider{
		ExternalId:        "eid1",
		NotInternalizable: false,
		Name:              "AP1",
		NamingHint:        "AP1Hint",
		Type:              aws.String("role_test"),
		Who: &WhoItem{
			Users:  []string{"uid1"},
			Groups: []string{"gid1"},
		},
		Access: []*Access{
			{
				What: []WhatItem{
					{
						Permissions: []string{"A", "B"},
						DataObject:  &do,
					},
				},
			},
		},
	})

	aps = append(aps, &AccessProvider{
		ExternalId:        "eid2",
		NotInternalizable: true,
		Name:              "AP2",
		NamingHint:        "AP2Hint",
		Action:            types.Filtered,
		Who: &WhoItem{
			Users:  []string{"uid1", "uid2"},
			Groups: []string{"gid1", "gid2"},
		},
		Access: []*Access{
			{
				What: []WhatItem{
					{
						Permissions: []string{"C"},
						DataObject:  &do,
					},
				},
			},
		},
	})

	err = apfc.AddAccessProviders(aps...)
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
	assert.Equal(t, types.Promise, apsr[0].Action)
	assert.Equal(t, 1, len(apsr[0].Who.Users))
	assert.Equal(t, "uid1", apsr[0].Who.Users[0])
	assert.Equal(t, 1, len(apsr[0].Who.Groups))
	assert.Equal(t, "gid1", apsr[0].Who.Groups[0])
	assert.Equal(t, 1, len(apsr[0].What))
	assert.Equal(t, 2, len(apsr[0].What[0].Permissions))
	assert.Equal(t, "A", apsr[0].What[0].Permissions[0])
	assert.Equal(t, "B", apsr[0].What[0].Permissions[1])
	assert.Equal(t, "Data Object 1", apsr[0].What[0].DataObject.FullName)
	assert.Equal(t, "schema", apsr[0].What[0].DataObject.Type)
	assert.NotNil(t, apsr[0].Type)
	assert.Equal(t, "role_test", *apsr[0].Type)

	assert.Equal(t, "eid2", apsr[1].ExternalId)
	assert.True(t, apsr[1].NotInternalizable)
	assert.Equal(t, "AP2", apsr[1].Name)
	assert.Equal(t, "AP2Hint", apsr[1].NamingHint)
	assert.Equal(t, types.Filtered, apsr[1].Action)
	assert.Equal(t, 2, len(apsr[1].Who.Users))
	assert.Equal(t, "uid1", apsr[1].Who.Users[0])
	assert.Equal(t, "uid2", apsr[1].Who.Users[1])
	assert.Equal(t, 2, len(apsr[1].Who.Groups))
	assert.Equal(t, "gid1", apsr[1].Who.Groups[0])
	assert.Equal(t, "gid2", apsr[1].Who.Groups[1])
	assert.Equal(t, 1, len(apsr[1].What))
	assert.Equal(t, 1, len(apsr[1].What[0].Permissions))
	assert.Equal(t, "C", apsr[1].What[0].Permissions[0])
	assert.Equal(t, "Data Object 1", apsr[1].What[0].DataObject.FullName)
	assert.Equal(t, "schema", apsr[1].What[0].DataObject.Type)
	assert.Nil(t, apsr[1].Type)
}

func TestShouldLock(t *testing.T) {
	tests := []struct {
		name         string
		lockAll      bool
		nameLocks    []string
		tagLocks     []string
		onIncomplete bool
		ap           *AccessProvider
		shouldLock   bool
	}{
		{
			name:       "lock all",
			lockAll:    true,
			ap:         &AccessProvider{},
			shouldLock: true,
		},
		{
			name:       "no lock",
			lockAll:    false,
			ap:         &AccessProvider{},
			shouldLock: false,
		},
		{
			name:      "lock by name",
			lockAll:   false,
			nameLocks: []string{"myname1"},
			ap: &AccessProvider{
				Name: "myname1",
			},
			shouldLock: true,
		},
		{
			name:      "lock by name - regex",
			lockAll:   false,
			nameLocks: []string{"my.+", "another.+"},
			ap: &AccessProvider{
				Name: "myname1",
			},
			shouldLock: true,
		},
		{
			name:      "lock by tag",
			lockAll:   false,
			nameLocks: []string{"my.+", "another.+"},
			tagLocks:  []string{"tag1:val1"},
			ap: &AccessProvider{
				Name: "blahname",
				Tags: []*tag.Tag{
					{
						Key:   "tag1",
						Value: "val1",
					},
				},
			},
			shouldLock: true,
		},
		{
			name:      "lock by tag - regex",
			lockAll:   false,
			nameLocks: []string{"my.+", "another.+"},
			tagLocks:  []string{"tag1:.+"},
			ap: &AccessProvider{
				Name: "blahname",
				Tags: []*tag.Tag{
					{
						Key:   "tag1",
						Value: "val1",
					},
				},
			},
			shouldLock: true,
		},
		{
			name:      "lock by tag - regex - no hit",
			lockAll:   false,
			nameLocks: []string{"my.+", "another.+"},
			tagLocks:  []string{"tag1:.+"},
			ap: &AccessProvider{
				Name: "blahname",
				Tags: []*tag.Tag{
					{
						Key:   "tag2",
						Value: "val1",
					},
				},
			},
			shouldLock: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lock, err := shouldLock("lock-name", tt.lockAll, tt.nameLocks, tt.tagLocks, tt.onIncomplete, tt.ap)
			assert.NoError(t, err)
			assert.Equal(t, tt.shouldLock, lock)
		})
	}
}

func TestCheckLocking(t *testing.T) {
	tests := []struct {
		name     string
		config   *access_provider.AccessSyncFromTarget
		ap       *AccessProvider
		resultAp *AccessProvider
	}{
		{
			name:     "no locking",
			config:   &access_provider.AccessSyncFromTarget{},
			ap:       &AccessProvider{},
			resultAp: &AccessProvider{},
		},
		{
			name:     "lock all",
			config:   &access_provider.AccessSyncFromTarget{FullyLockAll: true},
			ap:       &AccessProvider{},
			resultAp: &AccessProvider{NotInternalizable: true},
		},
		{
			name:     "lock all - by name",
			config:   &access_provider.AccessSyncFromTarget{FullyLockByName: []string{"ok", ".+ah"}},
			ap:       &AccessProvider{Name: "blah"},
			resultAp: &AccessProvider{NotInternalizable: true},
		},
		{
			name:     "lock all - by tags",
			config:   &access_provider.AccessSyncFromTarget{FullyLockByTag: []string{"k1:v1"}},
			ap:       &AccessProvider{Tags: []*tag.Tag{{Key: "k1", Value: "v1"}}},
			resultAp: &AccessProvider{NotInternalizable: true},
		},
		{
			name:     "lock all - on incomplete",
			config:   &access_provider.AccessSyncFromTarget{FullyLockWhenIncomplete: true},
			ap:       &AccessProvider{Incomplete: ptr.Bool(true)},
			resultAp: &AccessProvider{NotInternalizable: true},
		},

		{
			name:     "lock who",
			config:   &access_provider.AccessSyncFromTarget{LockAllWho: true},
			ap:       &AccessProvider{},
			resultAp: &AccessProvider{WhoLocked: ptr.Bool(true)},
		},
		{
			name:     "lock who - by name",
			config:   &access_provider.AccessSyncFromTarget{LockWhoByName: []string{"ok", ".+ah"}},
			ap:       &AccessProvider{Name: "blah"},
			resultAp: &AccessProvider{WhoLocked: ptr.Bool(true)},
		},
		{
			name:     "lock who - by tags",
			config:   &access_provider.AccessSyncFromTarget{LockWhoByTag: []string{"k1:.+"}},
			ap:       &AccessProvider{Tags: []*tag.Tag{{Key: "k1", Value: "xxx"}}},
			resultAp: &AccessProvider{WhoLocked: ptr.Bool(true)},
		},
		{
			name:     "lock who - on incomplete",
			config:   &access_provider.AccessSyncFromTarget{LockWhoWhenIncomplete: true},
			ap:       &AccessProvider{Incomplete: ptr.Bool(true)},
			resultAp: &AccessProvider{WhoLocked: ptr.Bool(true)},
		},

		{
			name:     "lock what",
			config:   &access_provider.AccessSyncFromTarget{LockAllWhat: true},
			ap:       &AccessProvider{},
			resultAp: &AccessProvider{WhatLocked: ptr.Bool(true)},
		},
		{
			name:     "lock what - by name",
			config:   &access_provider.AccessSyncFromTarget{LockWhatByName: []string{"ok", ".+ah"}},
			ap:       &AccessProvider{Name: "blah"},
			resultAp: &AccessProvider{WhatLocked: ptr.Bool(true)},
		},
		{
			name:     "lock what - by tags",
			config:   &access_provider.AccessSyncFromTarget{LockWhatByTag: []string{"k1:.+"}},
			ap:       &AccessProvider{Tags: []*tag.Tag{{Key: "k1", Value: "xxx"}}},
			resultAp: &AccessProvider{WhatLocked: ptr.Bool(true)},
		},
		{
			name:     "lock what - on incomplete",
			config:   &access_provider.AccessSyncFromTarget{LockWhatWhenIncomplete: true},
			ap:       &AccessProvider{Incomplete: ptr.Bool(true)},
			resultAp: &AccessProvider{WhatLocked: ptr.Bool(true)},
		},

		{
			name:     "lock delete",
			config:   &access_provider.AccessSyncFromTarget{LockAllDelete: true},
			ap:       &AccessProvider{},
			resultAp: &AccessProvider{DeleteLocked: ptr.Bool(true)},
		},
		{
			name:     "lock delete - by name",
			config:   &access_provider.AccessSyncFromTarget{LockDeleteByName: []string{"ok", ".+ah"}},
			ap:       &AccessProvider{Name: "blah"},
			resultAp: &AccessProvider{DeleteLocked: ptr.Bool(true)},
		},
		{
			name:     "lock delete - by tags",
			config:   &access_provider.AccessSyncFromTarget{LockDeleteByTag: []string{".+:xxx"}},
			ap:       &AccessProvider{Tags: []*tag.Tag{{Key: "yyy", Value: "xxx"}}},
			resultAp: &AccessProvider{DeleteLocked: ptr.Bool(true)},
		},
		{
			name:     "lock delete - on incomplete",
			config:   &access_provider.AccessSyncFromTarget{LockDeleteWhenIncomplete: true},
			ap:       &AccessProvider{Incomplete: ptr.Bool(true)},
			resultAp: &AccessProvider{DeleteLocked: ptr.Bool(true)},
		},

		{
			name:     "lock name",
			config:   &access_provider.AccessSyncFromTarget{LockAllNames: true},
			ap:       &AccessProvider{},
			resultAp: &AccessProvider{NameLocked: ptr.Bool(true)},
		},
		{
			name:     "lock name - by name",
			config:   &access_provider.AccessSyncFromTarget{LockNamesByName: []string{"ok", ".+ah"}},
			ap:       &AccessProvider{Name: "blah"},
			resultAp: &AccessProvider{NameLocked: ptr.Bool(true)},
		},
		{
			name:     "lock name - by tags",
			config:   &access_provider.AccessSyncFromTarget{LockNamesByTag: []string{".+:xxx"}},
			ap:       &AccessProvider{Tags: []*tag.Tag{{Key: "yyy", Value: "xxx"}}},
			resultAp: &AccessProvider{NameLocked: ptr.Bool(true)},
		},
		{
			name:     "lock name - on incomplete",
			config:   &access_provider.AccessSyncFromTarget{LockNamesWhenIncomplete: true},
			ap:       &AccessProvider{Incomplete: ptr.Bool(true)},
			resultAp: &AccessProvider{NameLocked: ptr.Bool(true)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkAp := *tt.ap
			err := checkLocking(&checkAp, tt.config)
			assert.NoError(t, err)
			assert.Equal(t, tt.resultAp.NotInternalizable, checkAp.NotInternalizable)
			assert.Equal(t, tt.resultAp.WhoLocked, checkAp.WhoLocked)
			assert.Equal(t, tt.resultAp.WhatLocked, checkAp.WhatLocked)
			assert.Equal(t, tt.resultAp.DeleteLocked, checkAp.DeleteLocked)
			assert.Equal(t, tt.resultAp.NameLocked, checkAp.NameLocked)
		})
	}
}
