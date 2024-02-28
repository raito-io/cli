package tag

import (
	"encoding/json"
	"fmt"
	"os"
)

type TagImportObject struct {
	DataObjectFullName *string `json:"dataObjectFullName,omitempty"`
	UserId             *string `json:"userId,omitempty"`
	Key                string  `json:"key"`
	StringValue        string  `json:"stringValue"`
	Source             string  `json:"source"`
}

//go:generate go run github.com/vektra/mockery/v2 --name=TagFileCreator --with-expecter
type TagFileCreator interface {
	AddTags(tags ...*TagImportObject) error
	Close()
	GetTagCount() int
}

type tagFileCreator struct {
	config *TagSyncConfig

	targetFile *os.File
	tagCount   int
}

func (t *tagFileCreator) AddTags(tags ...*TagImportObject) error {
	if len(tags) == 0 {
		return nil
	}

	for _, tag := range tags {
		if t.tagCount > 0 {
			t.targetFile.WriteString(",") //nolint: errcheck
		}

		t.targetFile.WriteString("\n") //nolint: errcheck

		tBuf, err := json.Marshal(tag)
		if err != nil {
			return fmt.Errorf("serializing tag: %w", err)
		}

		_, err = t.targetFile.Write(tBuf)
		if err != nil {
			return fmt.Errorf("writing to file %w", err)
		}

		t.tagCount++
	}

	return nil
}

func (t *tagFileCreator) Close() {
	_, _ = t.targetFile.WriteString("\n]")
	_ = t.targetFile.Close()
}

func (t *tagFileCreator) GetTagCount() int {
	return t.tagCount
}

func NewTagFileCreator(config *TagSyncConfig) (TagFileCreator, error) {
	tagFileC := tagFileCreator{
		config: config,
	}

	err := tagFileC.createTargetFile()
	if err != nil {
		return nil, err
	}

	_, err = tagFileC.targetFile.WriteString("[")
	if err != nil {
		return nil, err
	}

	return &tagFileC, nil
}

func (t *tagFileCreator) createTargetFile() error {
	f, err := os.Create(t.config.TargetFile)
	if err != nil {
		return fmt.Errorf("create target file: %w", err)
	}

	t.targetFile = f

	return nil
}
