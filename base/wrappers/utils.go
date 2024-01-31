package wrappers

import (
	"time"

	logger_utils "github.com/raito-io/cli/base/util/logger"
)

var logger = logger_utils.InitializeLogger()

func timedExecution(f func() error) (time.Duration, error) {
	start := time.Now()
	err := f()
	sec := time.Since(start).Round(time.Millisecond)

	return sec, err
}
