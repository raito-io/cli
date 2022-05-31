// Package access_provider provides the tooling to build the Raito access provider import file.
// Simply use the NewAccessProviderFileCreator function by passing in the config coming from the CLI
// to create the necessary file(s).
// The returned AccessProviderFileCreator can then be used (using the AddAccessProvider function)
// to write AccessProvider to the file.
// Make sure to call the Close function on the creator at the end (tip: use defer).
package access_provider

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/raito-io/cli/base/data_source"
	"github.com/raito-io/cli/common/api/data_access"
)

// DataObject represents a data object in the format that is suitable to be imported into a Raito data source.
type AccessProvider struct {
	ExternalId    string   `json:"externalId"`
	Name          string   `json:"name"`
	Users         []string `json:"users"`
	AccessObjects []Access `json:"accessObjects"`
}

type Access struct {
	DataObject  *data_source.DataObject `json:"dataobject"`
	Permissions []string                `json:"permissions"`
}

// AccessProviderFileCreator describes the interface for easily creating the access object import files
// to be imported by the Raito CLI.
type AccessProviderFileCreator interface {
	AddAccessProvider(dataAccessList []AccessProvider) error
	Close()
	GetAccessProviderCount() int
}

type accessProviderFileCreator struct {
	config *data_access.DataAccessSyncConfig

	targetFile      *os.File
	dataAccessCount int
}

// NewAccessProviderFileCreator creates a new AccessProviderFileCreator based on the configuration coming from
// the Raito CLI.
func NewAccessProviderFileCreator(config *data_access.DataAccessSyncConfig) (AccessProviderFileCreator, error) {
	dsI := accessProviderFileCreator{
		config: config,
	}

	err := dsI.createTargetFile()
	if err != nil {
		return nil, err
	}

	_, err = dsI.targetFile.Write([]byte("["))
	if err != nil {
		return nil, err
	}

	return &dsI, nil
}

// Close finalizes the import file and close it so it can be correctly read by the Raito CLI.
// This method must be called when all data objects have been added and before control is given back
// to the CLI. It's advised to call this using 'defer'.
func (d *accessProviderFileCreator) Close() {
	d.targetFile.Write([]byte("\n]")) //nolint:errcheck
	d.targetFile.Close()
}

// AddDataAccess adds the slice of data access elements to the import file.
// It returns an error when writing one of the objects fails (it will not process the other data objects after that).
// It returns nil if everything went well.
func (d *accessProviderFileCreator) AddAccessProvider(dataAccessList []AccessProvider) error {
	if len(dataAccessList) == 0 {
		return nil
	}

	for _, da := range dataAccessList {
		var err error

		if d.dataAccessCount > 0 {
			d.targetFile.Write([]byte(",")) //nolint:errcheck
		}
		d.targetFile.Write([]byte("\n")) //nolint:errcheck

		doBuf, err := json.Marshal(da)
		if err != nil {
			return fmt.Errorf("error while serializing data object with externalID %q", da.ExternalId)
		}
		d.targetFile.Write([]byte("\n")) //nolint:errcheck
		_, err = d.targetFile.Write(doBuf)

		// Only looking at writing errors at the end, supposing if one fails, all would fail
		if err != nil {
			return fmt.Errorf("error while writing to temp file %q", d.targetFile.Name())
		}
		d.dataAccessCount++
	}

	return nil
}

// GetAccessProviderCount returns the number of access elements that have been added to the import file.
func (d *accessProviderFileCreator) GetAccessProviderCount() int {
	return d.dataAccessCount
}

func (d *accessProviderFileCreator) createTargetFile() error {
	f, err := os.Create(d.config.TargetFile)
	if err != nil {
		return fmt.Errorf("error creating temporary file for data source importer: %s", err.Error())
	}
	d.targetFile = f
	return nil
}
