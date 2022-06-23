package data_usage

import (
	"encoding/json"
	"fmt"
	"os"

	ap "github.com/raito-io/cli/base/access_provider"
	"github.com/raito-io/cli/common/api/data_usage"
)

type Statement struct {
	ExternalId          string      `json:"externalId"`
	AccessedDataObjects []ap.Access `json:"accessedDataObjects"`
	Status              bool        `json:"status"`
	User                string      `json:"user"`
	StartTime           int64       `json:"startTime"`
	EndTime             int64       `json:"endTime"`
	TotalTime           float32     `json:"totalTime"`
	BytesTransferred    int         `json:"bytesTransferred"`
	RowsReturned        int         `json:"rowsReturned"`
}

// DataUsageFileCreator describes the interface for easily creating the data usage import files
// to be exported from the Raito CLI.
type DataUsageFileCreator interface {
	AddStatements(statements []Statement) error
	Close()
	GetStatementCount() int
}

type dataUsageFileCreator struct {
	config         *data_usage.DataUsageSyncConfig
	targetFile     *os.File
	statementCount int
}

func NewDataUsageFileCreator(config *data_usage.DataUsageSyncConfig) (DataUsageFileCreator, error) {
	duI := dataUsageFileCreator{
		config:         config,
		statementCount: 0,
	}

	err := duI.createTargetFile()
	if err != nil {
		return nil, err
	}

	_, err = duI.targetFile.WriteString("[")
	if err != nil {
		return nil, err
	}

	return &duI, nil
}

// AddTransaction adds the slice of data objects to the import file.
// It returns an error when writing one of the data objects fails (it will not process the other data objects after that).
// It returns nil if everything went well.
func (d *dataUsageFileCreator) AddStatements(statements []Statement) error {
	if len(statements) == 0 {
		return nil
	}

	for _, statement := range statements {
		var err error

		if d.statementCount > 0 {
			d.targetFile.WriteString(",") //nolint:errcheck
		}

		d.targetFile.WriteString("\n") //nolint:errcheck

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

// GetTransactionCount returns the number of data objects that has been added to the import file.
func (d *dataUsageFileCreator) GetStatementCount() int {
	return d.statementCount
}

func (d *dataUsageFileCreator) createTargetFile() error {
	f, err := os.Create(d.config.TargetFile)
	if err != nil {
		return fmt.Errorf("error creating temporary file for data usage importer: %s", err.Error())
	}
	d.targetFile = f

	return nil
}
