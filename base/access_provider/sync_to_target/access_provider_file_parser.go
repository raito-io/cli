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

	// TODO REFACTOR to be cleaned up when removing the old API.
	// Now this makes sure that the Access stuff is both on the AP layer as on the Access layer
	for _, ap := range ret.AccessProviders {
		if ap.Access == nil && ap.What != nil {
			ap.Access = []*Access{{
				What:       ap.What,
				ActualName: ap.ActualName,
				Id:         ap.Id,
			}}
		} else if len(ap.Access) > 0 && ap.What == nil {
			ap.What = ap.Access[0].What
			ap.ActualName = ap.Access[0].ActualName
		}
	}

	return &ret, nil
}
