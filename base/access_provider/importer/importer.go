package importer

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/go-hclog"
	"gopkg.in/yaml.v2"

	"github.com/raito-io/cli/base/access_provider"
)

func ParseAccessProviderImportFile(config *access_provider.AccessSyncToTarget) (*AccessProviderImport, error) {
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

type AccessProviderNameTranslationFileCreator interface {
	AddAccessProviderActualName(accessProviderActualName ...AccessProviderActualNameTranslation) error
	Close()
	GetAccessProviderCount() int
}

type accessProviderNameTranslationFileCreator struct {
	config *access_provider.AccessSyncToTarget

	actualNameTargetFile *os.File
	dataAccessCount      int
}

// NewAccessProviderFileCreator creates a new AccessProviderFileCreator based on the configuration coming from
// the Raito CLI.
func NewAccessProviderFileCreator(config *access_provider.AccessSyncToTarget) (AccessProviderNameTranslationFileCreator, error) {
	dsI := accessProviderNameTranslationFileCreator{
		config: config,
	}

	err := dsI.createTargetFile()
	if err != nil {
		return nil, err
	}

	_, err = dsI.actualNameTargetFile.WriteString("[")
	if err != nil {
		return nil, err
	}

	return &dsI, nil
}

func (d *accessProviderNameTranslationFileCreator) AddAccessProviderActualName(accessProviderActualName ...AccessProviderActualNameTranslation) error {
	if len(accessProviderActualName) == 0 {
		return nil
	}

	for _, dant := range accessProviderActualName {
		doBuf, err := json.Marshal(dant)
		if err != nil {
			return fmt.Errorf("error while serializing data object with ID %q and roleName %q", dant.AccessProviderId, dant.AccessProviderActualName)
		}

		if d.dataAccessCount > 0 {
			d.actualNameTargetFile.WriteString(",") //nolint:errcheck
		}

		d.actualNameTargetFile.WriteString("\n") //nolint:errcheck
		_, err = d.actualNameTargetFile.Write(doBuf)

		// Only looking at writing errors at the end, supposing if one fails, all would fail
		if err != nil {
			return fmt.Errorf("error while writing to temp file %q", d.actualNameTargetFile.Name())
		}
		d.dataAccessCount++
	}

	return nil
}

func (d *accessProviderNameTranslationFileCreator) Close() {
	d.actualNameTargetFile.WriteString("\n]") //nolint:errcheck
	d.actualNameTargetFile.Close()
}

func (d *accessProviderNameTranslationFileCreator) GetAccessProviderCount() int {
	return d.dataAccessCount
}

func (d *accessProviderNameTranslationFileCreator) createTargetFile() error {
	f, err := os.Create(d.config.TargetFile)
	if err != nil {
		return fmt.Errorf("error creating temporary file for data source importer: %s", err.Error())
	}
	d.actualNameTargetFile = f

	return nil
}
