package workerpool

import (
	"context"
	"errors"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkerPool_Go(t *testing.T) {
	t.Run("no errors", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		maxParallelExecutions := uint(2)
		workerPool := NewWorkerPool(ctx, maxParallelExecutions)

		var executionsDone atomic.Int32

		for i := 0; i < 10; i++ {
			workerPool.Go(func() error {
				time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

				executionsDone.Add(1)

				return nil
			})
		}

		err := workerPool.Wait()
		require.NoError(t, err)

		assert.Equal(t, int32(10), executionsDone.Load())
	})

	t.Run("execution errors", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		maxParallelExecutions := uint(2)
		workerPool := NewWorkerPool(ctx, maxParallelExecutions)

		var executionsDone atomic.Int32

		for i := 0; i < 10; i++ {
			workerPool.Go(func() error {
				time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

				executionsDone.Add(1)

				if i%3 == 0 {
					return errors.New("some error")
				}

				return nil
			})
		}

		err := workerPool.Wait()
		require.Error(t, err)

		var merr *multierror.Error
		ok := errors.As(err, &merr)
		require.Truef(t, ok, "expected error to be of type *multierror.Error, got: %T", err)

		assert.Len(t, merr.Errors, 4)

		assert.Equal(t, int32(10), executionsDone.Load())
	})

	t.Run("context canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		maxParallelExecutions := uint(2)
		workerPool := NewWorkerPool(ctx, maxParallelExecutions)

		var executionsDone atomic.Int32

		for i := 0; i < 10; i++ {
			workerPool.Go(func() error {
				time.Sleep(time.Duration(rand.Intn(2)+2) * time.Second)

				executionsDone.Add(1)

				return nil
			})
		}

		time.Sleep(1 * time.Second)
		cancel()

		err := workerPool.Wait()
		require.NoError(t, err)

		assert.Equal(t, int32(2), executionsDone.Load())
	})
}
