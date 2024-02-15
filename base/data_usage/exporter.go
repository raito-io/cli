package data_usage

import (
	"encoding/json"
	"fmt"
	"os"

	ap "github.com/raito-io/cli/base/access_provider/sync_from_target"
)

//go:generate go run github.com/vektra/mockery/v2 --name=DataUsageFileCreator --with-expecter

type Statement struct {
	ExternalId          string        `json:"externalId"`
	AccessedDataObjects []ap.WhatItem `json:"accessedDataObjects"`
	User                string        `json:"user"`
	Role                string        `json:"role"`
	Success             bool          `json:"success"`
	Status              string        `json:"status"`
	Query               string        `json:"query"`
	StartTime           int64         `json:"startTime"`
	EndTime             int64         `json:"endTime"`
	Bytes               int           `json:"bytes"`
	Rows                int           `json:"rows"`
	Credits             float32       `json:"credits"`
}

// DataUsageFileCreator describes the interface for easily creating the data usage import files
// to be exported from the Raito CLI.
type DataUsageFileCreator interface {
	AddStatements(statements []Statement) error
	Close()
	GetStatementCount() int
	GetImportFileSize() uint64
}

type dataUsageFileCreator struct {
	config         *DataUsageSyncConfig
	targetFile     *os.File
	statementCount int
	fileByteSize   uint64
}

func NewDataUsageFileCreator(config *DataUsageSyncConfig) (DataUsageFileCreator, error) {
	duI := dataUsageFileCreator{
		config:         config,
		statementCount: 0,
		fileByteSize:   2, // 2 bytes for closing the file, '\n]'
	}

	err := duI.createTargetFile()
	if err != nil {
		return nil, err
	}

	_, err = duI.targetFile.WriteString("[")
	if err != nil {
		return nil, err
	}

	duI.fileByteSize += 1

	return &duI, nil
}

// AddStatements adds the slice of data objects to the import file.
// It returns an error when writing one of the data objects fails (it will not process the other data objects after that).
// It returns nil if everything went well.
func (d *dataUsageFileCreator) AddStatements(statements []Statement) error {
	if len(statements) == 0 {
		return nil
	}

	for ind := range statements {
		statement := statements[ind]
		var err error

		if d.statementCount > 0 {
			d.targetFile.WriteString(",") //nolint:errcheck
			d.fileByteSize += 1
		}

		d.targetFile.WriteString("\n") //nolint:errcheck
		d.fileByteSize += 1

		doBuf, err := json.Marshal(statement)
		if err != nil {
			return fmt.Errorf("error while serializing data object with externalID %q", statement.ExternalId)
		}
		_, err = d.targetFile.Write(doBuf)

		// Only looking at writing errors at the end, supposing if one fails, all would fail
		if err != nil {
			return fmt.Errorf("error while writing to temp file %q", d.targetFile.Name())
		}

		d.statementCount++
		d.fileByteSize += uint64(len(doBuf))
	}

	return nil
}

// Close finalizes the import file and close it so it can be correctly read by the Raito CLI.
// This method must be called when all data objects have been added and before control is given back
// to the CLI. It's advised to call this using 'defer'.
func (d *dataUsageFileCreator) Close() {
	d.targetFile.WriteString("\n]") //nolint:errcheck
	d.targetFile.Close()
}

// GetStatementCount returns the number of data objects that has been added to the import file.
func (d *dataUsageFileCreator) GetStatementCount() int {
	return d.statementCount
}

// GetImportFileSize returns the approximate byte size of the data that has been added to the import file.
func (d *dataUsageFileCreator) GetImportFileSize() uint64 {
	return d.fileByteSize
}

func (d *dataUsageFileCreator) createTargetFile() error {
	f, err := os.Create(d.config.TargetFile)
	if err != nil {
		return fmt.Errorf("error creating temporary file for data usage importer: %s", err.Error())
	}
	d.targetFile = f

	return nil
}
