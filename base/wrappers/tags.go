package wrappers

import (
	"context"
	"fmt"
	"time"

	"github.com/raito-io/cli/base/tag"
)

//go:generate go run github.com/vektra/mockery/v2 --name=TagHandler --with-expecter
type TagHandler interface {
	AddTags(tags ...*tag.TagImportObject) error
}

type TagSyncer interface {
	SyncTags(ctx context.Context, tagsHandler TagHandler, config *tag.TagSyncConfig) ([]string, error)
}

type TagSyncFactoryFn func(ctx context.Context, configParams *tag.TagSyncConfig) (TagSyncer, func(), error)

func TagSyncFactory(syncer TagSyncFactoryFn) tag.TagSyncer {
	return &tagSyncFunction{
		syncer:             NewSyncFactory[tag.TagSyncConfig](syncer),
		fileCreatorFactory: tag.NewTagFileCreator,
	}
}

type tagSyncFunction struct {
	tag.TagSyncerVersionHandler

	syncer             SyncFactory[tag.TagSyncConfig, TagSyncer]
	fileCreatorFactory func(config *tag.TagSyncConfig) (tag.TagFileCreator, error)
}

func (t *tagSyncFunction) SyncTags(ctx context.Context, config *tag.TagSyncConfig) (_ *tag.TagSyncResult, err error) {
	defer func() {
		if err != nil {
			logger.Error(fmt.Sprintf("Failure during tag sync: %v", err))
		}
	}()

	logger.Info("Starting tag synchronisation")
	logger.Debug("Creating file for storing tags")

	fileCreator, err := t.fileCreatorFactory(config)
	if err != nil {
		return nil, fmt.Errorf("create file: %w", err)
	}

	defer fileCreator.Close()

	start := time.Now()

	syncer, err := t.syncer.Create(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create syncer: %w", err)
	}

	tagSources, err := syncer.SyncTags(ctx, fileCreator, config)
	if err != nil {
		return nil, fmt.Errorf("sync tags: %w", err)
	}

	sec := time.Since(start).Round(time.Millisecond)

	logger.Info(fmt.Sprintf("Finished synchronising %d tags in %s", fileCreator.GetTagCount(), sec))

	return &tag.TagSyncResult{
		Tags:            int32(fileCreator.GetTagCount()),
		TagSourcesScope: tagSources,
	}, nil
}

func (t *tagSyncFunction) Close() {
	t.syncer.Close()
}
