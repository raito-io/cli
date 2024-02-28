package tag

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"testing"

	"github.com/aws/smithy-go/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTagFileCreator(t *testing.T) {
	tempFile, _ := os.Create("tempfile-tags- " + strconv.Itoa(rand.Int()) + ".json")
	defer os.Remove(tempFile.Name())

	config := TagSyncConfig{
		TargetFile:   tempFile.Name(),
		DataSourceId: "myDataSource",
	}

	tfc, err := NewTagFileCreator(&config)
	require.NoError(t, err)
	require.NotNil(t, tfc)

	tags := []*TagImportObject{
		{
			DataObjectFullName: ptr.String("do1"),
			Key:                "key1",
			StringValue:        "value1",
			Source:             "source1",
		},
		{
			DataObjectFullName: ptr.String("do2"),
			Key:                "key2",
			StringValue:        "value2",
			Source:             "source1",
		},
		{
			DataObjectFullName: ptr.String("do3"),
			Key:                "key3",
			StringValue:        "value3",
			Source:             "source2",
		},
	}

	err = tfc.AddTags(tags...)
	require.NoError(t, err)

	tfc.Close()

	assert.Equal(t, 3, tfc.GetTagCount())

	bytes, err := ioutil.ReadAll(tempFile)
	require.Nil(t, err)

	tagsr := make([]TagImportObject, 0, 4)

	err = json.Unmarshal(bytes, &tagsr)
	require.NoError(t, err)

	require.Len(t, tagsr, 3)
	assert.Equal(t, *tags[0], tagsr[0])
	assert.Equal(t, *tags[1], tagsr[1])
	assert.Equal(t, *tags[2], tagsr[2])

}
