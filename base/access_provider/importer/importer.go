package importer

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/raito-io/cli/base/access_provider"
	"gopkg.in/yaml.v2"
)

func ParseAccessProviderImportFile(config *access_provider.AccessSyncConfig) (*AccessProviderImport, error) {
	var ret AccessProviderImport

	af, err := os.Open(config.SourceFile)
	if err != nil {
		hclog.L().Error(fmt.Sprintf("Error while opening access file %q: %s", config.SourceFile, err.Error()))
		return nil, err
	}

	buf, err := io.ReadAll(af)
	if err != nil {
		hclog.L().Error(fmt.Sprintf("Error while reading access file %q: %s", config.SourceFile, err.Error()))
		return nil, err
	}

	if json.Valid(buf) {
		err = json.Unmarshal(buf, &ret)
		if err != nil {
			return nil, err
		}
	} else {
		err = yaml.Unmarshal(buf, &ret)
		if err != nil {
			return nil, err
		}
	}

	return &ret, nil
}
