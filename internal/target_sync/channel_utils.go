package target_sync

import (
	"fmt"
	"sync"
)

func ChanMap[I any, O any](c <-chan I, f func(I) O) <-chan O {
	out := make(chan O)
	go func() {
		defer close(out)

		for i := range c {
			out <- f(i)
		}
	}()

	return out
}

func ErrorWrap(c <-chan error, format string) <-chan error {
	return ChanMap(c, func(err error) error {
		return fmt.Errorf(format, err)
	})
}

func Merge[T any](cs ...<-chan T) <-chan T {
	out := make(chan T)

	var wg sync.WaitGroup
	wg.Add(len(cs))

	for _, c := range cs {
		go func(c <-chan T) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func Single[T any](c <-chan T) T {
	for v := range c {
		return v
	}

	var defaultVal T
	return defaultVal
}
