package url

import (
	"strings"

	"github.com/spf13/viper"

	"github.com/raito-io/cli/base/util/url"
	"github.com/raito-io/cli/internal/constants"
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
	} else if env == constants.EnvironmentStaging {
		return "https://api.staging.raito.dev/"
	} else {
		return "https://api.raito.cloud/"
	}
}

func CreateRaitoURL(schemaAndHost, path string) string {
	if !strings.HasSuffix(schemaAndHost, "/") {
		schemaAndHost += "/"
	}

	relPath := url.GetRelativePath(path)

	return schemaAndHost + relPath
}
