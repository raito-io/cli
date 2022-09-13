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

func ParseAccessProviderImportFile(config *access_provider.AccessSyncExportConfig) (*AccessProviderImport, error) {
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
	AddAccessProviderNameTranslation(accesProviderNameTranslations ...AccessProviderNameTranslation) error
	AddAccessProvidersNameTranslations(accesProviderNameTranslations []AccessProviderNameTranslation) error
	Close()
	GetAccessProviderCount() int
}

type accessProviderNameTranslationFileCreator struct {
	config *access_provider.AccessSyncExportConfig

	targetFile      *os.File
	dataAccessCount int
}

// NewAccessProviderFileCreator creates a new AccessProviderFileCreator based on the configuration coming from
// the Raito CLI.
func NewAccessProviderFileCreator(config *access_provider.AccessSyncExportConfig) (AccessProviderNameTranslationFileCreator, error) {
	dsI := accessProviderNameTranslationFileCreator{
		config: config,
	}

	err := dsI.createTargetFile()
	if err != nil {
		return nil, err
	}

	_, err = dsI.targetFile.WriteString("[")
	if err != nil {
		return nil, err
	}

	return &dsI, nil
}

func (d *accessProviderNameTranslationFileCreator) AddAccessProviderNameTranslation(accesProviderNameTranslations ...AccessProviderNameTranslation) error {
	return d.AddAccessProvidersNameTranslations(accesProviderNameTranslations)
}

func (d *accessProviderNameTranslationFileCreator) AddAccessProvidersNameTranslations(accesProviderNameTranslations []AccessProviderNameTranslation) error {
	if len(accesProviderNameTranslations) == 0 {
		return nil
	}

	for _, dant := range accesProviderNameTranslations {
		doBuf, err := json.Marshal(dant)
		if err != nil {
			return fmt.Errorf("error while serializing data object with ID %q and roleName %q", dant.AccessProviderId, dant.AccessProviderActualName)
		}

		if d.dataAccessCount > 0 {
			d.targetFile.WriteString(",") //nolint:errcheck
		}

		d.targetFile.WriteString("\n") //nolint:errcheck
		_, err = d.targetFile.Write(doBuf)

		// Only looking at writing errors at the end, supposing if one fails, all would fail
		if err != nil {
			return fmt.Errorf("error while writing to temp file %q", d.targetFile.Name())
		}
		d.dataAccessCount++
	}

	return nil
}

func (d *accessProviderNameTranslationFileCreator) Close() {
	d.targetFile.WriteString("\n]") //nolint:errcheck
	d.targetFile.Close()
}

func (d *accessProviderNameTranslationFileCreator) GetAccessProviderCount() int {
	return d.dataAccessCount
}

func (d *accessProviderNameTranslationFileCreator) createTargetFile() error {
	f, err := os.Create(d.config.TargetFile)
	if err != nil {
		return fmt.Errorf("error creating temporary file for data source importer: %s", err.Error())
	}
	d.targetFile = f

	return nil
}
