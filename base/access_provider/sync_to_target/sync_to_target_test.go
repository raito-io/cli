package sync_to_target

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/raito-io/cli/base/access_provider"
	"github.com/raito-io/cli/base/access_provider/types"
)

func TestParseImportFileYaml(t *testing.T) {
	config := &access_provider.AccessSyncToTarget{
		SourceFile: "./testdata/data-access.yaml",
	}
	parsed, err := ParseAccessProviderImportFile(config)
	validateParsedAccessFile(t, parsed, err)
}

func TestParseImportFileJSON(t *testing.T) {
	config := &access_provider.AccessSyncToTarget{
		SourceFile: "./testdata/data-access.json",
	}
	parsed, err := ParseAccessProviderImportFile(config)
	validateParsedAccessFile(t, parsed, err)
}

func validateParsedAccessFile(t *testing.T, parsed *AccessProviderImport, err error) {
	assert.NotNil(t, parsed)
	assert.Nil(t, err)

	fmt.Println(parsed.LastCalculated)
	assert.Equal(t, int64(100), parsed.LastCalculated)
	assert.Equal(t, 1, len(parsed.AccessProviders))

	ap := parsed.AccessProviders[0]
	assert.Equal(t, "11111111", ap.Id)
	assert.Equal(t, "blah", ap.Name)
	assert.Equal(t, "Lots of blah", ap.Description)
	assert.Equal(t, "Blah_", ap.NamingHint)
	assert.Equal(t, types.Mask, ap.Action)
	require.NotNil(t, ap.Type)
	assert.Equal(t, "role_test", *ap.Type)

	assert.Equal(t, "Blahkes", *ap.ActualName)
	assert.Equal(t, 2, len(ap.Who.Users))
	assert.Equal(t, 2, len(ap.What))
	assert.Equal(t, "zzz.yyy.table1", ap.What[0].DataObject.FullName)
	assert.Equal(t, 2, len(ap.What[0].Permissions))
	assert.Equal(t, "zzz.yyy.table2", ap.What[1].DataObject.FullName)
	assert.Equal(t, 3, len(ap.What[1].Permissions))

	assert.Len(t, ap.Owners, 2)
	assert.Equal(t, "owner@raito.io", *ap.Owners[0].Email)
	assert.Equal(t, "ownerAccount", *ap.Owners[0].AccountName)
	assert.Equal(t, "ownerGroup", *ap.Owners[1].GroupName)
}
