package url

import (
	"strings"

	"github.com/spf13/viper"

	"github.com/raito-io/cli/base/util/url"
	"github.com/raito-io/cli/internal/constants"
)

func GetRaitoURL() string {
	override := viper.GetString(constants.URLOverrideFlag)
	if override != "" {
		return override
	}

	return "https://api.raito.cloud/"
}

func CreateRaitoURL(schemaAndHost, path string) string {
	if !strings.HasSuffix(schemaAndHost, "/") {
		schemaAndHost += "/"
	}

	relPath := url.GetRelativePath(path)

	return schemaAndHost + relPath
}
