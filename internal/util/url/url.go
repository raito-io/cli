package url

import (
	"github.com/raito-io/cli/common/util/url"
	"github.com/raito-io/cli/internal/constants"
	"github.com/spf13/viper"
	"strings"
)

var TestURL = ""

func GetRaitoURL() string {
	if TestURL != "" {
		return TestURL
	} else if viper.GetBool(constants.DevFlag) {
		return "http://localhost:8080/"
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
