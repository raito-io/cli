package wrappers

import (
	"context"

	"github.com/raito-io/cli/base/util/config"
)

type SyncFactory[S any] struct {
	factory     func(ctx context.Context, configParams *config.ConfigMap) (S, func(), error)
	initialized bool

	s         S
	cleanupFn func()
}

func NewSyncFactory[S any](factory func(ctx context.Context, configParams *config.ConfigMap) (S, func(), error)) SyncFactory[S] {
	return SyncFactory[S]{
		factory: factory,
	}
}

func NewDummySyncFactoryFn[S any](syncer S) func(ctx context.Context, configParams *config.ConfigMap) (S, func(), error) {
	return func(_ context.Context, _ *config.ConfigMap) (S, func(), error) {
		return syncer, func() {}, nil
	}
}

func (s *SyncFactory[S]) Create(ctx context.Context, configParams *config.ConfigMap) (S, error) {
	if !s.initialized {
		var err error

		s.s, s.cleanupFn, err = s.factory(ctx, configParams)
		if err != nil {
			return s.s, err
		}

		s.initialized = true
	}

	return s.s, nil
}

func (s *SyncFactory[S]) Close() {
	if s.initialized {
		s.cleanupFn()
		s.initialized = false
	}
}
