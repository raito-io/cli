package url

import (
	"testing"

	"github.com/raito-io/cli/internal/constants"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetRaitoURL(t *testing.T) {
	assert.Equal(t, "https://api.raito.cloud/", GetRaitoURL())
	viper.Set(constants.URLOverrideFlag, "https://blah.raito.blah")
	assert.Equal(t, "https://blah.raito.blah", GetRaitoURL())
}

func TestCreateRaitoURL(t *testing.T) {
	assert.Equal(t, "https://api.raito.cloud/my/path", CreateRaitoURL("https://api.raito.cloud/", "/my/path"))
	assert.Equal(t, "https://api.raito.cloud/my/path", CreateRaitoURL("https://api.raito.cloud", "/my/path"))
	assert.Equal(t, "https://api.raito.cloud/my/path", CreateRaitoURL("https://api.raito.cloud/", "my/path"))
	assert.Equal(t, "https://api.raito.cloud/my/path", CreateRaitoURL("https://api.raito.cloud", "my/path"))
}
