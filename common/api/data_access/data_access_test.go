package data_access

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"sort"
	"testing"
)

func parseDataAccess(input []byte) (*DataAccessResult, error) {
	var ret DataAccessResult
	err := yaml.Unmarshal(input, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func TestBasicDataAccessParsing(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("testdata/data-access-1.yaml")
	if err != nil {
		t.Fatalf("Error reading test data file: %s", err.Error())
	}
	dar, err := parseDataAccess(yamlFile)
	if err != nil {
		t.Fatalf("Error while parsing rules: %s", err.Error())
	}

	das := dar.AccessRights
	assert.Equal(t, 2, len(das))

	da0 := das[0]
	assert.Equal(t, "11111111-1111-1111-1111-111111111111", da0.Id)
	assert.Equal(t, "my first rule", da0.Rule.Name)
	assert.Equal(t, "This describes the first rule", da0.Rule.Description)
	assert.Equal(t, "xxxx-xxxx-xxxx", da0.Rule.Id)
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
	assert.Equal(t, "my first rule2", da1.Rule.Name)
	assert.Equal(t, "This describes the first rule2", da1.Rule.Description)
	assert.Equal(t, "xxxx-xxxx-xxxx2", da1.Rule.Id)
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
	dar, err := parseDataAccess(yamlFile)
	if err != nil {
		t.Fatalf("Error while parsing rules: %s", err.Error())
	}

	das := dar.AccessRights
	assert.Equal(t, 2, len(das))

	da0 := das[0]
	assert.Equal(t, "11111111-1111-1111-1111-111111111111", da0.Id)
	assert.Nil(t, da0.Rule)
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
	assert.Equal(t, "Rule 1", da1.Rule.Name)
	assert.Equal(t, "Some description", da1.Rule.Description)
	assert.Equal(t, "00000000-0000-0000-0009-000000000001", da1.Rule.Id)
	dao1 := *(da1.DataObject)
	assert.Equal(t, "BE_Employees", dao1.Name)
	assert.Equal(t, "Table", dao1.Type)
	assert.Nil(t, dao1.Parent)

	assert.Equal(t, 1, len(da1.Permissions))
	assert.Equal(t, "SELECT", da1.Permissions[0])

	assert.Equal(t, 3964, len(da1.Users))
}

func TestDataAccess_CalculateHash(t *testing.T) {
	yamlFile, _ := ioutil.ReadFile("testdata/generated.yaml")
	dar, _ := parseDataAccess(yamlFile)
	das := dar.AccessRights
	assert.Equal(t, 2, len(das))

	da0 := das[0]
	hash := da0.CalculateHash()
	hash2 := da0.CalculateHash()

	assert.NotNil(t, hash)
	assert.NotEmpty(t, hash)
	assert.Equal(t, hash2, hash)
}

func TestDataAccess_Merge(t *testing.T) {
	do := DataObject{
		Type: "table",
		Name: "good_table",
	}
	da1 := DataAccess{
		DataObject: &do,
		Permissions: []string { "p1", "p2" },
		Users: []string { "u1", "u2" },
	}
	da2 := DataAccess{
		DataObject: &do,
		Permissions: []string { "p1", "p2" },
		Users: []string { "u1", "u3" },
	}
	da3 := DataAccess{
		DataObject: &do,
		Permissions: []string { "p1", "p2" },
		Users: []string { "u1", "u2" },
	}

	merged := da1.Merge([]*DataAccess { &da2, &da3 })

	assert.Equal(t, 3, len(merged.Users))
	sort.Strings(merged.Users)
	assert.Equal(t, "u1", merged.Users[0])
	assert.Equal(t, "u2", merged.Users[1])
	assert.Equal(t, "u3", merged.Users[2])
}