package post_processing

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/go-hclog"
	"gopkg.in/yaml.v2"

	"github.com/raito-io/cli/base/access_provider/sync_from_target"
)

//go:generate go run github.com/vektra/mockery/v2 --name=PostProcessorSourceFileParser --with-expecter
type PostProcessorSourceFileParser interface {
	ParseAccessProviders() ([]*sync_from_target.AccessProvider, error)
}

func NewPostProcessorSourceFileParser(sourceFile string) (PostProcessorSourceFileParser, error) {
	return &postProcessorSourceFileParser{
		SourceFile: sourceFile,
	}, nil
}

type postProcessorSourceFileParser struct {
	SourceFile string
}

func (p *postProcessorSourceFileParser) ParseAccessProviders() ([]*sync_from_target.AccessProvider, error) {
	var ret []*sync_from_target.AccessProvider

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
