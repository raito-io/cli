package cmd

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterAccessProviders(t *testing.T) {
	outputFile, filtered, total, err := filterAccessProviders("testTarget", "testdata/apply-input1.yaml", "Nicolas.*")
	defer os.Remove(outputFile)

	fmt.Println(fmt.Sprintf("Filtered accessProviders: %s", strings.Join(filtered, ",")))

	assert.NoError(t, err)
	assert.Equal(t, 33, total)
	assert.True(t, len(outputFile) > 0)
	assert.True(t, strings.HasSuffix(outputFile, ".yaml"))
	assert.Equal(t, 4, len(filtered))

	fs, err := readFileStructure(outputFile)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(fs.AccessProviders))

	outputFile2, filtered2, total, err := filterAccessProviders("testTarget", "testdata/apply-input1.yaml", "(?i).*test.*")
	defer os.Remove(outputFile2)

	fmt.Println(fmt.Sprintf("Filtered accessProviders: %s", strings.Join(filtered2, ",")))

	assert.NoError(t, err)
	assert.Equal(t, 33, total)
	assert.True(t, len(outputFile2) > 0)
	assert.True(t, strings.HasSuffix(outputFile2, ".yaml"))
	assert.Equal(t, 12, len(filtered2))

	fs2, err := readFileStructure(outputFile2)
	assert.NoError(t, err)
	assert.Equal(t, 12, len(fs2.AccessProviders))
}
