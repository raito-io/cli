package url

import (
	"github.com/raito-io/cli/internal/constants"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRaitoURL(t *testing.T) {
	assert.Equal(t, "https://api.raito.cloud/", GetRaitoURL())
	viper.Set(constants.EnvironmentFlag, "dev")
	assert.Equal(t, "http://localhost:8080/", GetRaitoURL())
	viper.Set(constants.EnvironmentFlag, "test")
	assert.Equal(t, "https://api.raito.dev/", GetRaitoURL())
	viper.Set(constants.EnvironmentFlag, "blah")
	assert.Equal(t, "https://api.raito.cloud/", GetRaitoURL())
}

func TestCreateRaitoURL(t *testing.T) {
	assert.Equal(t, "https://api.raito.cloud/my/path", CreateRaitoURL("https://api.raito.cloud/", "/my/path"))
	assert.Equal(t, "https://api.raito.cloud/my/path", CreateRaitoURL("https://api.raito.cloud", "/my/path"))
	assert.Equal(t, "https://api.raito.cloud/my/path", CreateRaitoURL("https://api.raito.cloud/", "my/path"))
	assert.Equal(t, "https://api.raito.cloud/my/path", CreateRaitoURL("https://api.raito.cloud", "my/path"))
}
