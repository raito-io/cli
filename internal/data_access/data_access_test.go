package data_access

import (
	"github.com/raito-io/cli/common/api/data_access"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestBasicDataAccessParsing(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("testdata/data-access-1.yaml")
	if err != nil {
		t.Fatalf("Error reading test data file: %s", err.Error())
	}
	dar, err := ParseDataAccess(yamlFile)
	if err != nil {
		t.Fatalf("Error while parsing rules: %s", err.Error())
	}

	das := dar.AccessRights
	assert.Equal(t, 2, len(das))

	da0 := das[0]
	assert.Equal(t, "11111111-1111-1111-1111-111111111111", da0.Id)
	assert.Equal(t, "my first rule", da0.Provider.Name)
	assert.Equal(t, "This describes the first rule", da0.Provider.Description)
	assert.Equal(t, "xxxx-xxxx-xxxx", da0.Provider.Id)
	dao0 := *(da0.DataObject)
	assert.Equal(t, "table1", dao0.Name)
	assert.Equal(t, "table", dao0.Type)
	dao0p := *(dao0.Parent)
	assert.Equal(t, "yyy", dao0p.Name)
	assert.Equal(t, "schema", dao0p.Type)
	dao0pp := *(dao0p.Parent)
	assert.Equal(t, "zzz", dao0pp.Name)
	assert.Equal(t, "database", dao0pp.Type)

	assert.Equal(t, 2, len(da0.Permissions))
	assert.Equal(t, "select", da0.Permissions[0])
	assert.Equal(t, "delete", da0.Permissions[1])

	assert.Equal(t, 2, len(da0.Users))
	assert.Equal(t, "bart", da0.Users[0])
	assert.Equal(t, "dieter", da0.Users[1])

	da1 := das[1]
	assert.Equal(t, "11111111-1111-1111-1111-111111111112", da1.Id)
	assert.Equal(t, "my first rule2", da1.Provider.Name)
	assert.Equal(t, "This describes the first rule2", da1.Provider.Description)
	assert.Equal(t, "xxxx-xxxx-xxxx2", da1.Provider.Id)
	dao1 := *(da1.DataObject)
	assert.Equal(t, "table2", dao1.Name)
	assert.Equal(t, "table", dao1.Type)
	dao1p := *(dao1.Parent)
	assert.Equal(t, "yyy", dao1p.Name)
	assert.Equal(t, "schema", dao1p.Type)
	dao1pp := *(dao1p.Parent)
	assert.Equal(t, "zzz", dao1pp.Name)
	assert.Equal(t, "database", dao1pp.Type)

	assert.Equal(t, 1, len(da1.Permissions))
	assert.Equal(t, "select", da1.Permissions[0])

	assert.Equal(t, 2, len(da1.Users))
	assert.Equal(t, "katleen", da1.Users[0])
	assert.Equal(t, "stefanie", da1.Users[1])
}

func TestRulesGeneratedParsing(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("testdata/generated.yaml")
	if err != nil {
		t.Fatalf("Error reading test data file: %s", err.Error())
	}
	dar, err := ParseDataAccess(yamlFile)
	if err != nil {
		t.Fatalf("Error while parsing rules: %s", err.Error())
	}

	das := dar.AccessRights
	assert.Equal(t, 2, len(das))

	da0 := das[0]
	assert.Equal(t, "11111111-1111-1111-1111-111111111111", da0.Id)
	assert.Nil(t, da0.Provider)
	dao0 := *(da0.DataObject)
	assert.Equal(t, "BE_Employees", dao0.Name)
	assert.Equal(t, "Table", dao0.Type)
	assert.Nil(t, dao0.Parent)

	assert.Equal(t, 2, len(da0.Permissions))
	assert.Equal(t, "SELECT", da0.Permissions[0])
	assert.Equal(t, "DELETE", da0.Permissions[1])

	assert.Equal(t, 1, len(da0.Users))
	assert.Equal(t, "hAGzp", da0.Users[0])

	da1 := das[1]
	assert.Equal(t, "a7379157-d4f2-459e-bada-a003f27ad03d", da1.Id)
	assert.Equal(t, "Rule 1", da1.Provider.Name)
	assert.Equal(t, "Some description", da1.Provider.Description)
	assert.Equal(t, "00000000-0000-0000-0009-000000000001", da1.Provider.Id)
	dao1 := *(da1.DataObject)
	assert.Equal(t, "BE_Employees", dao1.Name)
	assert.Equal(t, "Table", dao1.Type)
	assert.Nil(t, dao1.Parent)

	assert.Equal(t, 1, len(da1.Permissions))
	assert.Equal(t, "SELECT", da1.Permissions[0])

	assert.Equal(t, 3964, len(da1.Users))
}

func TestRulesGeneratedMerging(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("testdata/generated.yaml")
	if err != nil {
		t.Fatalf("Error reading test data file: %s", err.Error())
	}
	dar, err := ParseDataAccess(yamlFile)
	if err != nil {
		t.Fatalf("Error while parsing rules: %s", err.Error())
	}

	das := dar.AccessRights
	das = flattenDataAccessList(das)
	assert.Equal(t, 2, len(das))
}

func TestRulesMerging(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("testdata/data-access-merge-1.yaml")
	if err != nil {
		t.Fatalf("Error reading test data file: %s", err.Error())
	}
	dar, err := ParseDataAccess(yamlFile)
	if err != nil {
		t.Fatalf("Error while parsing rules: %s", err.Error())
	}

	das := dar.AccessRights
	das = flattenDataAccessList(das)
	assert.Equal(t, 3, len(das))

	for _, da := range das {
		switch da.DataObject.Path {
		case "zzz.yyy.table1":
			checkTable1(t, da)
		case "zzz.yyy.table2":
			checkTable2(t, da)
		default:
			t.Fatalf("Found unknown dataobject path %q", da.DataObject.Path)
		}
	}
}

func checkTable1(t *testing.T, da *data_access.DataAccess) {
	assert.Equal(t, 3, len(da.Users))
	assert.Contains(t, da.Users, "dieter", "Didn't find user dieter")
	assert.Contains(t, da.Users, "katleen", "Didn't find user katleen")
	assert.Contains(t, da.Users, "bart", "Didn't find user bart")

	assert.Equal(t, 2, len(da.Permissions))
	assert.Contains(t, da.Permissions, "select", "Didn't find permission select")
	assert.Contains(t, da.Permissions, "delete", "Didn't find permission delete")
}

func checkTable2(t *testing.T, da *data_access.DataAccess) {
	if len(da.Permissions) == 1 {
		assert.Equal(t, "select", da.Permissions[0])
		assert.Equal(t, 2, len(da.Users))
		assert.Contains(t, da.Users, "katleen", "Didn't find user katleen")
		assert.Contains(t, da.Users, "stefanie", "Didn't find user stefanie")
	} else if len(da.Permissions) == 2 {
		assert.Equal(t, 2, len(da.Permissions))
		assert.Contains(t, da.Permissions, "select", "Didn't find permission select")
		assert.Contains(t, da.Permissions, "update", "Didn't find permission update")

		assert.Equal(t, 2, len(da.Users))
		assert.Contains(t, da.Users, "jos", "Didn't find user jos")
		assert.Contains(t, da.Users, "stefanie", "Didn't find user stefanie")
	} else {
		t.Fatal("Incorrect number of permissions found in data access element")
	}
}
