package sync_from_target

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/go-hclog"
	"gopkg.in/yaml.v2"
)

//go:generate go run github.com/vektra/mockery/v2 --name=AccessProviderSyncFromTargetFileParser --with-expecter
type AccessProviderSyncFromTargetFileParser interface {
	ParseAccessProviders() ([]*AccessProvider, error)
}

func NewAccessProviderSyncFromTargetFileParser(sourceFile string) (AccessProviderSyncFromTargetFileParser, error) {
	return &accessProviderSyncFromTargetFileParser{
		SourceFile: sourceFile,
	}, nil
}

type accessProviderSyncFromTargetFileParser struct {
	SourceFile string
}

func (p *accessProviderSyncFromTargetFileParser) ParseAccessProviders() ([]*AccessProvider, error) {
	var ret []*AccessProvider

	af, err := os.Open(p.SourceFile)
	if err != nil {
		hclog.L().Error(fmt.Sprintf("Error while opening access file %q: %s", p.SourceFile, err.Error()))
		return nil, err
	}

	defer af.Close()

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

	return ret, nil
}
