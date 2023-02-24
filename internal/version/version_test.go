package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	SetVersion("v1.2.3", "yyyy-mm-dd")
	assert.Equal(t, "v1.2.3 (yyyy-mm-dd)", GetVersionString())
}
