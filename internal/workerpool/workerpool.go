package workerpool

import (
	"context"

	"github.com/hashicorp/go-multierror"
)

type WorkerPool struct {
	ctx     context.Context
	barrier chan struct{}
	wg      multierror.Group
}

func NewWorkerPool(ctx context.Context, maxParallelExecutions uint) WorkerPool {
	barrier := make(chan struct{}, maxParallelExecutions)
	for range maxParallelExecutions {
		barrier <- struct{}{}
	}

	return WorkerPool{
		ctx:     ctx,
		barrier: barrier,
		wg:      multierror.Group{},
	}
}

func (wp *WorkerPool) Go(fn func() error) {
	wp.wg.Go(func() error {
		select {
		case <-wp.ctx.Done():
			return nil
		case <-wp.barrier:
			defer func() {
				wp.barrier <- struct{}{}
			}()
		}

		return fn()
	})
}

func (wp *WorkerPool) Wait() error {
	return wp.wg.Wait().ErrorOrNil()
}
