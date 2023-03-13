package wrappers

import (
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/raito-io/cli/base"
)

var logger = base.Logger()

func timedExecution(f func() error) (time.Duration, error) {
	start := time.Now()
	err := f()
	sec := time.Since(start).Round(time.Millisecond)

	if err != nil {
		if _, found := status.FromError(err); !found {
			err = status.Error(codes.Internal, err.Error())
		}
	}

	return sec, err
}
