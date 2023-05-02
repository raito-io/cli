package data_usage

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/raito-io/cli/base/access_provider/sync_from_target"

	"github.com/stretchr/testify/assert"

	"github.com/raito-io/cli/base/data_source"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestDataUsageFileCreator(t *testing.T) {
	tempFile, _ := os.Create("tempfile-" + strconv.Itoa(rand.Int()) + ".json")
	defer os.Remove(tempFile.Name())
	config := DataUsageSyncConfig{
		TargetFile: tempFile.Name(),
	}
	dufc, err := NewDataUsageFileCreator(&config)
	assert.Nil(t, err)
	assert.NotNil(t, dufc)

	dus := make([]Statement, 0, 3)

	dus = append(dus, Statement{
		ExternalId: "transaction1",
		AccessedDataObjects: []sync_from_target.WhatItem{
			{DataObject: &data_source.DataObjectReference{"schema1.table1.column1", "column"},
				Permissions: []string{"SELECT"}},
		},
		Success:   true,
		Status:    "",
		User:      "Alice",
		StartTime: 1654073198000,
		EndTime:   1654073198050,
		Bytes:     52,
		Rows:      3,
		Credits:   0,
	})
	dus = append(dus, Statement{
		ExternalId: "transaction2",
		AccessedDataObjects: []sync_from_target.WhatItem{
			{DataObject: &data_source.DataObjectReference{"schema1.table2.column3", "column"},
				Permissions: []string{"ALTER"}},
			{DataObject: &data_source.DataObjectReference{"schema1.table2.column5", "column"},
				Permissions: []string{"ALTER"}},
			{DataObject: &data_source.DataObjectReference{"schema1.table2.column7", "column"},
				Permissions: []string{"ALTER"}},
		},
		Success:   false,
		Status:    "",
		User:      "Alice",
		StartTime: 1654073199000,
		EndTime:   1654073199060,
		Bytes:     180,
		Rows:      27,
	})
	dus = append(dus, Statement{
		ExternalId: "transaction3",
		AccessedDataObjects: []sync_from_target.WhatItem{
			{DataObject: &data_source.DataObjectReference{"schema3", "schema"},
				Permissions: []string{"GRANT"}},
		},
		Success:   true,
		Status:    "",
		User:      "Bob",
		StartTime: 1654073200000,
		EndTime:   1654073200020,
		Bytes:     0,
		Rows:      0,
		Credits:   0,
	})

	err = dufc.AddStatements(dus)
	assert.Nil(t, err)
	dufc.Close()

	assert.Equal(t, 3, dufc.GetStatementCount())
	assert.Equal(t, int64(1045), dufc.GetImportFileSize())

	bytes, err := ioutil.ReadAll(tempFile)
	assert.Nil(t, err)

	dusr := make([]Statement, 0, 3)
	err = json.Unmarshal(bytes, &dusr)
	assert.Nil(t, err)

	assert.Equal(t, 3, len(dusr))

	assert.Equal(t, "transaction1", dusr[0].ExternalId)
	assert.Equal(t, []string{"SELECT"}, dusr[0].AccessedDataObjects[0].Permissions)
	assert.Equal(t, &data_source.DataObjectReference{FullName: "schema1.table1.column1", Type: "column"}, dusr[0].AccessedDataObjects[0].DataObject)
	assert.Equal(t, true, dusr[0].Success)
	assert.Equal(t, "Alice", dusr[0].User)
	assert.Equal(t, int64(1654073198000), dusr[0].StartTime)
	assert.Equal(t, int64(1654073198050), dusr[0].EndTime)
	assert.Equal(t, 52, dusr[0].Bytes)
	assert.Equal(t, 3, dusr[0].Rows)

	assert.Equal(t, "transaction2", dusr[1].ExternalId)
	assert.Equal(t, []string{"ALTER"}, dusr[1].AccessedDataObjects[0].Permissions)
	assert.Equal(t, []string{"ALTER"}, dusr[1].AccessedDataObjects[1].Permissions)
	assert.Equal(t, []string{"ALTER"}, dusr[1].AccessedDataObjects[2].Permissions)
	assert.Equal(t, &data_source.DataObjectReference{FullName: "schema1.table2.column3", Type: "column"}, dusr[1].AccessedDataObjects[0].DataObject)
	assert.Equal(t, &data_source.DataObjectReference{FullName: "schema1.table2.column5", Type: "column"}, dusr[1].AccessedDataObjects[1].DataObject)
	assert.Equal(t, &data_source.DataObjectReference{FullName: "schema1.table2.column7", Type: "column"}, dusr[1].AccessedDataObjects[2].DataObject)
	assert.Equal(t, false, dusr[1].Success)
	assert.Equal(t, "Alice", dusr[1].User)
	assert.Equal(t, int64(1654073199000), dusr[1].StartTime)
	assert.Equal(t, int64(1654073199060), dusr[1].EndTime)
	assert.Equal(t, 180, dusr[1].Bytes)
	assert.Equal(t, 27, dusr[1].Rows)

	assert.Equal(t, "transaction3", dusr[2].ExternalId)
	assert.Equal(t, []string{"GRANT"}, dusr[2].AccessedDataObjects[0].Permissions)
	assert.Equal(t, &data_source.DataObjectReference{FullName: "schema3", Type: "schema"}, dusr[2].AccessedDataObjects[0].DataObject)
	assert.Equal(t, true, dusr[2].Success)
	assert.Equal(t, "Bob", dusr[2].User)
	assert.Equal(t, int64(1654073200000), dusr[2].StartTime)
	assert.Equal(t, int64(1654073200020), dusr[2].EndTime)
	assert.Equal(t, 0, dusr[2].Bytes)
	assert.Equal(t, 0, dusr[2].Rows)
}
