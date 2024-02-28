package wrappers

import (
	"context"
)

type SyncFactory[C any, S any] struct {
	factory     func(ctx context.Context, configParams *C) (S, func(), error)
	initialized bool

	s         S
	cleanupFn func()
}

func NewSyncFactory[C any, S any](factory func(ctx context.Context, configParams *C) (S, func(), error)) SyncFactory[C, S] {
	return SyncFactory[C, S]{
		factory: factory,
	}
}

func NewDummySyncFactoryFn[C any, S any](syncer S) func(ctx context.Context, configParams *C) (S, func(), error) {
	return func(_ context.Context, _ *C) (S, func(), error) {
		return syncer, func() {}, nil
	}
}

func (s *SyncFactory[C, S]) Create(ctx context.Context, configParams *C) (S, error) {
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

func (s *SyncFactory[C, S]) Close() {
	if s.initialized {
		s.cleanupFn()
		s.initialized = false
	}
}
