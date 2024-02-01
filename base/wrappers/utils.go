package wrappers

import (
	"time"

	"github.com/raito-io/cli/base"
)

var logger = base.Logger()

func timedExecution(f func() error) (time.Duration, error) {
	start := time.Now()
	err := f()
	sec := time.Since(start).Round(time.Millisecond)

	return sec, err
}
