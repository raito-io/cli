package url

import (
	"github.com/raito-io/cli/internal/constants"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRaitoURL(t *testing.T) {
	assert.Equal(t, "https://api.raito.io/", GetRaitoURL())
	viper.Set(constants.DevFlag, true)
	assert.Equal(t, "http://localhost:8080/", GetRaitoURL())
}

func TestCreateRaitoURL(t *testing.T) {
	assert.Equal(t, "https://api.raito.io/my/path", CreateRaitoURL("https://api.raito.io/", "/my/path"))
	assert.Equal(t, "https://api.raito.io/my/path", CreateRaitoURL("https://api.raito.io", "/my/path"))
	assert.Equal(t, "https://api.raito.io/my/path", CreateRaitoURL("https://api.raito.io/", "my/path"))
	assert.Equal(t, "https://api.raito.io/my/path", CreateRaitoURL("https://api.raito.io", "my/path"))
}
