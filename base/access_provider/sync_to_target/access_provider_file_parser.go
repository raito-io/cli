package sync_to_target

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/go-hclog"
	"gopkg.in/yaml.v2"

	"github.com/raito-io/cli/base/access_provider"
)

//go:generate go run github.com/vektra/mockery/v2 --name=AccessProviderImportFileParser --with-expecter
type AccessProviderImportFileParser interface {
	ParseAccessProviders() (*AccessProviderImport, error)
}

func NewAccessProviderFileParser(config *access_provider.AccessSyncToTarget) (AccessProviderImportFileParser, error) {
	return &accessProviderFileParser{
		SourceFile: config.SourceFile,
	}, nil
}

type accessProviderFileParser struct {
	SourceFile string
}

func (p *accessProviderFileParser) ParseAccessProviders() (*AccessProviderImport, error) {
	var ret AccessProviderImport

	af, err := os.Open(p.SourceFile)
	if err != nil {
		hclog.L().Error(fmt.Sprintf("Error while opening access file %q: %s", p.SourceFile, err.Error()))
		return nil, err
	}

	buf, err := io.ReadAll(af)
	if err != nil {
		hclog.L().Error(fmt.Sprintf("Error while reading access file %q: %s", p.SourceFile, err.Error()))
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
