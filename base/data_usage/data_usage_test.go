package data_usage

import (
	"encoding/json"
	"github.com/raito-io/cli/base/access_provider"
	"github.com/raito-io/cli/base/data_source"
	"github.com/raito-io/cli/common/api/data_usage"
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

func TestDataUsageFileCreator(t *testing.T) {
	tempFile, _ := os.Create("tempfile-" + strconv.Itoa(rand.Int()) + ".json")
	defer os.Remove(tempFile.Name())
	config := data_usage.DataUsageSyncConfig{
		TargetFile: tempFile.Name(),
	}
	dufc, err := NewDataUsageFileCreator(&config)
	assert.Nil(t, err)
	assert.NotNil(t, dufc)

	dus := make([]Statement, 0, 3)

	dus = append(dus, Statement{
		ExternalId: "transaction1",
		AccessedDataObjects: []access_provider.Access{
			{DataObjectReference: &data_source.DataObjectReference{"schema1.table1.column1", "column"},
				Permissions: []string{"SELECT"}},
		},
		Status:           true,
		User:             "Alice",
		StartTime:        1654073198000,
		EndTime:          1654073198050,
		TotalTime:        0.05,
		BytesTransferred: 52,
		RowsReturned:     3,
	})
	dus = append(dus, Statement{
		ExternalId: "transaction2",
		AccessedDataObjects: []access_provider.Access{
			{DataObjectReference: &data_source.DataObjectReference{"schema1.table2.column3", "column"},
				Permissions: []string{"ALTER"}},
			{DataObjectReference: &data_source.DataObjectReference{"schema1.table2.column5", "column"},
				Permissions: []string{"ALTER"}},
			{DataObjectReference: &data_source.DataObjectReference{"schema1.table2.column7", "column"},
				Permissions: []string{"ALTER"}},
		},
		Status:           false,
		User:             "Alice",
		StartTime:        1654073199000,
		EndTime:          1654073199060,
		TotalTime:        0.06,
		BytesTransferred: 180,
		RowsReturned:     27,
	})
	dus = append(dus, Statement{
		ExternalId: "transaction3",
		AccessedDataObjects: []access_provider.Access{
			{DataObjectReference: &data_source.DataObjectReference{"schema3", "schema"},
				Permissions: []string{"GRANT"}},
		},
		Status:           true,
		User:             "Bob",
		StartTime:        1654073200000,
		EndTime:          1654073200020,
		TotalTime:        0.02,
		BytesTransferred: 0,
		RowsReturned:     0,
	})

	err = dufc.AddStatements(dus)
	assert.Nil(t, err)
	dufc.Close()

	assert.Equal(t, 3, dufc.GetStatementCount())

	bytes, err := ioutil.ReadAll(tempFile)
	assert.Nil(t, err)

	dusr := make([]Statement, 0, 3)
	err = json.Unmarshal(bytes, &dusr)
	assert.Nil(t, err)

	assert.Equal(t, 3, len(dusr))

	assert.Equal(t, "transaction1", dusr[0].ExternalId)
	assert.Equal(t, []string{"SELECT"}, dusr[0].AccessedDataObjects[0].Permissions)
	assert.Equal(t, &data_source.DataObjectReference{FullName: "schema1.table1.column1", Type: "column"}, dusr[0].AccessedDataObjects[0].DataObjectReference)
	assert.Equal(t, true, dusr[0].Status)
	assert.Equal(t, "Alice", dusr[0].User)
	assert.Equal(t, int64(1654073198000), dusr[0].StartTime)
	assert.Equal(t, int64(1654073198050), dusr[0].EndTime)
	assert.Equal(t, float32(0.05), dusr[0].TotalTime)
	assert.Equal(t, 52, dusr[0].BytesTransferred)
	assert.Equal(t, 3, dusr[0].RowsReturned)

	assert.Equal(t, "transaction2", dusr[1].ExternalId)
	assert.Equal(t, []string{"ALTER"}, dusr[1].AccessedDataObjects[0].Permissions)
	assert.Equal(t, []string{"ALTER"}, dusr[1].AccessedDataObjects[1].Permissions)
	assert.Equal(t, []string{"ALTER"}, dusr[1].AccessedDataObjects[2].Permissions)
	assert.Equal(t, &data_source.DataObjectReference{FullName: "schema1.table2.column3", Type: "column"}, dusr[1].AccessedDataObjects[0].DataObjectReference)
	assert.Equal(t, &data_source.DataObjectReference{FullName: "schema1.table2.column5", Type: "column"}, dusr[1].AccessedDataObjects[1].DataObjectReference)
	assert.Equal(t, &data_source.DataObjectReference{FullName: "schema1.table2.column7", Type: "column"}, dusr[1].AccessedDataObjects[2].DataObjectReference)
	assert.Equal(t, false, dusr[1].Status)
	assert.Equal(t, "Alice", dusr[1].User)
	assert.Equal(t, int64(1654073199000), dusr[1].StartTime)
	assert.Equal(t, int64(1654073199060), dusr[1].EndTime)
	assert.Equal(t, float32(0.06), dusr[1].TotalTime)
	assert.Equal(t, 180, dusr[1].BytesTransferred)
	assert.Equal(t, 27, dusr[1].RowsReturned)

	assert.Equal(t, "transaction3", dusr[2].ExternalId)
	assert.Equal(t, []string{"GRANT"}, dusr[2].AccessedDataObjects[0].Permissions)
	assert.Equal(t, &data_source.DataObjectReference{FullName: "schema3", Type: "schema"}, dusr[2].AccessedDataObjects[0].DataObjectReference)
	assert.Equal(t, true, dusr[2].Status)
	assert.Equal(t, "Bob", dusr[2].User)
	assert.Equal(t, int64(1654073200000), dusr[2].StartTime)
	assert.Equal(t, int64(1654073200020), dusr[2].EndTime)
	assert.Equal(t, float32(0.02), dusr[2].TotalTime)
	assert.Equal(t, 0, dusr[2].BytesTransferred)
	assert.Equal(t, 0, dusr[2].RowsReturned)
}
