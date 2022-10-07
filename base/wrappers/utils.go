package wrappers

import (
	"strings"
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

func find(s []string, q string) bool {
	for _, r := range s {
		if strings.EqualFold(r, q) {
			return true
		}
	}

	return false
}
