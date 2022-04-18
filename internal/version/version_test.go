package version

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVersion(t *testing.T) {
	SetVersion("v1.2.3", "yyyy-mm-dd")
	assert.Equal(t, "v1.2.3 (yyyy-mm-dd)", GetVersion())
}
