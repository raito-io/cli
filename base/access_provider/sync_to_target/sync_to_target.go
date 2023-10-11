package sync_to_target

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/go-hclog"
	"gopkg.in/yaml.v2"

	"github.com/raito-io/cli/base/access_provider"
)

//go:generate go run github.com/vektra/mockery/v2 --name=SyncFeedbackFileCreator --with-expecter

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

type SyncFeedbackFileCreator interface {
	AddAccessProviderFeedback(accessProviderFeedback AccessProviderSyncFeedback) error
	Close()
	GetAccessProviderCount() int
}

type syncFeedbackFileCreator struct {
	config *access_provider.AccessSyncToTarget

	feedbackFile           *os.File
	dataAccessCount        int
	definedAccessProviders map[string]struct{}
}

// NewFeedbackFileCreator creates a new SyncFeedbackFileCreator based on the configuration coming from the Raito CLI.
func NewFeedbackFileCreator(config *access_provider.AccessSyncToTarget) (SyncFeedbackFileCreator, error) {
	dsI := syncFeedbackFileCreator{
		config:                 config,
		feedbackFile:           nil,
		dataAccessCount:        0,
		definedAccessProviders: map[string]struct{}{},
	}

	err := dsI.createTargetFile()
	if err != nil {
		return nil, err
	}

	_, err = dsI.feedbackFile.WriteString("[")
	if err != nil {
		return nil, err
	}

	return &dsI, nil
}

func (d *syncFeedbackFileCreator) AddAccessProviderFeedback(accessProviderFeedback AccessProviderSyncFeedback) error {
	if _, found := d.definedAccessProviders[accessProviderFeedback.AccessProvider]; found {
		return errors.New("access provider is already defined in feedback file")
	}

	apFeedback := accessProviderFeedbackInformation{
		ExternalId: accessProviderFeedback.AccessProvider,
		AccessFeedbackObjects: []accessSyncFeedbackInformation{
			{
				AccessId:   accessProviderFeedback.AccessProvider,
				ActualName: accessProviderFeedback.ActualName,
				ExternalId: accessProviderFeedback.ExternalId,
				Type:       accessProviderFeedback.Type,
				Errors:     accessProviderFeedback.Errors,
				Warnings:   accessProviderFeedback.Warnings,
			},
		},
	}

	doBuf, err := json.Marshal(apFeedback)
	if err != nil {
		return fmt.Errorf("error while serialzing access provider feedback data object with ID %q: %s", accessProviderFeedback.AccessProvider, err.Error())
	}

	newLine := bytes.NewBufferString("")
	if d.dataAccessCount > 0 {
		_, err = newLine.WriteString(",")
		if err != nil {
			return err
		}
	}

	_, err = newLine.WriteString("\n")
	if err != nil {
		return err
	}

	_, err = newLine.Write(doBuf)
	if err != nil {
		return err
	}

	_, err = d.feedbackFile.Write(newLine.Bytes())
	if err != nil {
		return fmt.Errorf("error while writing to temp file %q: %s", d.feedbackFile.Name(), err.Error())
	}

	d.dataAccessCount++
	d.definedAccessProviders[accessProviderFeedback.AccessProvider] = struct{}{}

	return nil
}

func (d *syncFeedbackFileCreator) Close() {
	d.feedbackFile.WriteString("\n]") //nolint:errcheck
	d.feedbackFile.Close()
}

func (d *syncFeedbackFileCreator) GetAccessProviderCount() int {
	return d.dataAccessCount
}

func (d *syncFeedbackFileCreator) createTargetFile() error {
	f, err := os.Create(d.config.FeedbackTargetFile)
	if err != nil {
		return fmt.Errorf("error creating temporary file for data source importer: %s", err.Error())
	}
	d.feedbackFile = f

	return nil
}
