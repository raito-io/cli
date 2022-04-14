package url

import (
	"github.com/raito-io/cli/common/util/url"
	"github.com/raito-io/cli/internal/constants"
	"github.com/spf13/viper"
	"strings"
)

var TestURL = ""

func GetRaitoURL() string {
	env := viper.GetString(constants.EnvironmentFlag)
	if TestURL != "" {
		return TestURL
	} else if env == constants.EnvironmentDev {
		return "http://localhost:8080/"
	} else if env == constants.EnvironmentTest {
		return "https://api.raito.dev/"
	} else {
		return "https://api.raito.io/"
	}
}

func CreateRaitoURL(schemaAndHost, path string) string {
	if !strings.HasSuffix(schemaAndHost, "/") {
		schemaAndHost = schemaAndHost + "/"
	}
	relPath := url.GetRelativePath(path)
	return schemaAndHost + relPath
}
