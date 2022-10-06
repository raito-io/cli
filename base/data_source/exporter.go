// Package data_source provides the tooling to build the Raito data source import file.
// Simply use the NewDataSourceFileCreator function by passing in the config coming from the CLI
// to create the necessary file(s).
// The returned DataSourceFileCreator can then be used (using the AddDataObjects function)
// to write DataObjects to the file.
// Make sure to call the Close function on the creator at the end (tip: use defer).
package data_source

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

//go:generate go run github.com/vektra/mockery/v2 --name=DataSourceFileCreator --with-expecter

// DataObject represents a data object in the format that is suitable to be imported into a Raito data source.
type DataObject struct {
	ExternalId       string                 `json:"externalId"`
	Name             string                 `json:"name"`
	FullName         string                 `json:"fullName"`
	Type             string                 `json:"type"`
	Description      string                 `json:"description"`
	ParentExternalId string                 `json:"parentExternalId"`
	Tags             map[string]interface{} `json:"tags"`
}

// DataObjectReference represents the reference to a DataObject suitable for e.g. defining the What in Access Provider import
type DataObjectReference struct {
	FullName string `json:"fullName" yaml:"fullName"`
	Type     string `json:"type" yaml:"type"`
}

// DataSourceFileCreator describes the interface for easily creating the data object import files
// to be imported by the Raito CLI.
// Use GetDataSourceDetails() to access DataSourceDetails setters like SetAvailablePermission()
type DataSourceFileCreator interface {
	AddDataObjects(dataObjects ...*DataObject) error
	SetDataSourceName(name string)
	SetDataSourceFullname(name string)
	SetDataSourceDescription(desc string)
	Close()
	GetDataObjectCount() int
	GetDataSourceDetails() *DataSourceDetails
}

type dataSourceFileCreator struct {
	config *DataSourceSyncConfig

	targetFile        *os.File
	dataObjectCount   int
	dataSourceDetails DataSourceDetails
}

// NewDataSourceFileCreator creates a new DataSourceFileCreator based on the configuration coming from
// the Raito CLI.
func NewDataSourceFileCreator(config *DataSourceSyncConfig) (DataSourceFileCreator, error) {
	dsI := dataSourceFileCreator{
		config:            config,
		dataSourceDetails: newDataSourceDetails(config.DataSourceId),
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

// Close finalizes the import file and close it so it can be correctly read by the Raito CLI.
// This method must be called when all data objects have been added and before control is given back
// to the CLI. It's advised to call this using 'defer'.
func (d *dataSourceFileCreator) Close() {
	d.targetFile.WriteString("\n]") //nolint:errcheck
	d.targetFile.Close()
}

func (d *dataSourceFileCreator) GetDataSourceDetails() *DataSourceDetails {
	return &d.dataSourceDetails
}

func (d *dataSourceFileCreator) prependDataSourceDataObject(dataObjects []*DataObject) []*DataObject {
	finalDataObjects := []*DataObject{&d.dataSourceDetails.dataSource}
	finalDataObjects = append(finalDataObjects, dataObjects...)

	return finalDataObjects
}

// AddDataObjects adds the slice of data objects to the import file.
// It returns an error when writing one of the data objects fails (it will not process the other data objects after that).
// It returns nil if everything went well.
func (d *dataSourceFileCreator) AddDataObjects(dataObjects ...*DataObject) error {
	if len(dataObjects) == 0 {
		return nil
	}

	// validate wether the first DataObject represents the Data Source, if not generate one
	if d.dataObjectCount == 0 {
		if !strings.EqualFold(dataObjects[0].Type, "datasource") {
			dataObjects = d.prependDataSourceDataObject(dataObjects)
		}
	}

	for _, do := range dataObjects { //nolint

		// set datasource DO as the parent for all direct children
		if do.ParentExternalId == "" && do.Type != "datasource" {
			do.ParentExternalId = d.dataSourceDetails.dataSource.ExternalId
		}

		var err error

		if d.dataObjectCount > 0 {
			d.targetFile.WriteString(",") //nolint:errcheck
		}

		d.targetFile.WriteString("\n") //nolint:errcheck

		doBuf, err := json.Marshal(do)
		if err != nil {
			return fmt.Errorf("error while serializing data object with externalID %q", do.ExternalId)
		}

		d.targetFile.WriteString("\n") //nolint:errcheck
		_, err = d.targetFile.Write(doBuf)

		// Only looking at writing errors at the end, supposing if one fails, all would fail
		if err != nil {
			return fmt.Errorf("error while writing to temp file %q", d.targetFile.Name())
		}
		d.dataObjectCount++
	}

	return nil
}

func (d *dataSourceFileCreator) SetDataSourceName(name string) {
	d.dataSourceDetails.SetName(name)
}

func (d *dataSourceFileCreator) SetDataSourceFullname(name string) {
	d.dataSourceDetails.SetFullname(name)
}

func (d *dataSourceFileCreator) SetDataSourceDescription(desc string) {
	d.dataSourceDetails.SetDescription(desc)
}

// GetDataObjectCount returns the number of data objects that has been added to the import file.
func (d *dataSourceFileCreator) GetDataObjectCount() int {
	return d.dataObjectCount
}

func (d *dataSourceFileCreator) createTargetFile() error {
	f, err := os.Create(d.config.TargetFile)
	if err != nil {
		return fmt.Errorf("error creating temporary file for data source importer: %s", err.Error())
	}
	d.targetFile = f

	return nil
}
