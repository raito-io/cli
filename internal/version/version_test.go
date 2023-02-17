package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	SetVersion("v1.2.3", "yyyy-mm-dd")
	assert.Equal(t, "1.2.3 (yyyy-mm-dd)", GetVersionString())
}

func TestSemanticVersion(t *testing.T) {
	SetVersion("1.2.3", "yyyy-mm-dd")
	assert.Equal(t, "1.2.3", GetCliVersion().String())
}

func TestCliPluginConstraint(t *testing.T) {
	constraint := CliPluginConstraint()
	assert.NotNil(t, constraint)
}
