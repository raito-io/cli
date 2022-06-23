package url

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRelativePath(t *testing.T) {
	assert.Equal(t, "this/is/it", GetRelativePath("http://blah.io/this/is/it"))
	assert.Equal(t, "this/is/it", GetRelativePath("https://blah.io/this/is/it"))
	assert.Equal(t, "this/is/it", GetRelativePath("/this/is/it"))
	assert.Equal(t, "this/is/it", GetRelativePath("this/is/it"))
}

func TestCutOffSchema(t *testing.T) {
	assert.Equal(t, "blah.io/this/is/it", CutOffSchema("http://blah.io/this/is/it"))
	assert.Equal(t, "blah.io/this/is/it", CutOffSchema("blah.io/this/is/it"))
	assert.Equal(t, "this/is/it", CutOffSchema("ftp://this/is/it"))
}

func TestCutOffPrefix(t *testing.T) {
	assert.Equal(t, "//blah.io/this/is/it", CutOffPrefix("http://blah.io/this/is/it", "http:"))
	assert.Equal(t, "https://blah.io/this/is/it", CutOffPrefix("https://blah.io/this/is/it", "http:"))
	assert.Equal(t, "it", CutOffPrefix("it", "http:"))
	assert.Equal(t, "ithttp:blah", CutOffPrefix("ithttp:blah", "http:"))
}

func TestCutOffSuffix(t *testing.T) {
	assert.Equal(t, "http://blah.io/this/is/", CutOffSuffix("http://blah.io/this/is/it", "it"))
	assert.Equal(t, "https://blah.io/this/is/it", CutOffSuffix("https://blah.io/this/is/it", "ht"))
}
