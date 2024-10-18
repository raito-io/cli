package data_usage

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

//go:generate go run github.com/vektra/mockery/v2 --name=DataUsageFileCreator --with-expecter

type UsageDataObjectReference struct {
	FullName string `json:"fullName"`
	Type     string `json:"type"`
}

type UsageDataObjectItem struct {
	// Permissions is a list of permissions that were used to access the data object
	Permissions []string `json:"permissions,omitempty"` // Deprecated: please use GlobalPermission instead

	// GlobalPermission is the global permission that was used to access the data object
	GlobalPermission ActionType `json:"globalPermission,omitempty"`

	// DataObject represents the data object that was accessed
	DataObject UsageDataObjectReference `json:"dataObject"`
}

//go:generate go run github.com/raito-io/enumer -json -type=ActionType -transform=lower
type ActionType int

const (
	UnknownAction ActionType = iota
	Read
	Write
	Admin
)

type Statement struct {
	// ExternalId is the unique identifier for the statement
	ExternalId string `json:"externalId"`

	// AccessedDataObjects is a list of data objects that were accessed by the statement
	AccessedDataObjects []UsageDataObjectItem `json:"accessedDataObjects"`

	// User is the user email that executed the statement
	User string `json:"user,omitempty"`

	// Role is the role used to execute the statement (if applicable)
	Role string `json:"role,omitempty"`

	// Success indicates if the statement was successful
	Success bool `json:"success"`

	// Status is the raw status of the statement
	Status string `json:"status,omitempty"`

	// Query is the raw query that was executed
	Query string `json:"query"`

	// StartTime is the time the statement started executing
	StartTime int64 `json:"startTime,omitempty"`

	// EndTime is the time the statement finished executing
	EndTime int64 `json:"endTime,omitempty"`

	// Bytes is the number of bytes returned by the statement
	Bytes int `json:"bytes,omitempty"`

	// Rows is the number of rows returned by the statement
	Rows int `json:"rows,omitempty"`

	// Credits is the number of credits used by the statement
	Credits float32 `json:"credits,omitempty"`

	// Parsing information
	Error string `json:"error,omitempty"`
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
	totalByteSize  uint64

	maxBytesPerFile uint64

	actualFileNames []string
}

func NewDataUsageFileCreator(config *DataUsageSyncConfig) (DataUsageFileCreator, error) {
	duI := dataUsageFileCreator{
		config:          config,
		statementCount:  0,
		fileByteSize:    0,
		totalByteSize:   0,
		maxBytesPerFile: config.MaxBytesPerFile,
	}

	err := duI.openNewFile()
	if err != nil {
		return nil, fmt.Errorf("open new file: %w", err)
	}

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

		size := uint64(len(doBuf))

		if d.maxBytesPerFile != 0 && d.fileByteSize+size > uint64(d.maxBytesPerFile) {
			d.closeFile()

			err = d.openNewFile()
			if err != nil {
				return fmt.Errorf("open new file: %w", err)
			}
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
	d.closeFile()
}

func (d *dataUsageFileCreator) closeFile() {
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

func (d *dataUsageFileCreator) GetActualFileNames() []string {
	return d.actualFileNames
}

func (d *dataUsageFileCreator) createTargetFile() error {
	actualName := d.config.TargetFile

	if strings.HasSuffix(d.config.TargetFile, ".json") || strings.HasSuffix(d.config.TargetFile, ".yaml") {
		split := strings.Split(d.config.TargetFile, ".")

		nr := len(d.actualFileNames)

		actualName = fmt.Sprintf("%s_%d.%s", strings.Join(split[0:len(split)-2], "."), nr, split[1])
	} else {
		actualName = d.config.TargetFile + "_0"
	}

	f, err := os.Create(actualName)
	if err != nil {
		return fmt.Errorf("error creating temporary file for data usage importer: %s", err.Error())
	}
	d.targetFile = f

	d.actualFileNames = append(d.actualFileNames, actualName)

	return nil
}

func (d *dataUsageFileCreator) openNewFile() error {
	err := d.createTargetFile()
	if err != nil {
		return fmt.Errorf("create target file: %w", err)
	}

	_, err = d.targetFile.WriteString("[")
	if err != nil {
		return fmt.Errorf("write open json token: %w", err)
	}

	d.fileByteSize = 3   //Already include '\n]'
	d.totalByteSize += 3 //Already include '\n]'

	return nil
}
