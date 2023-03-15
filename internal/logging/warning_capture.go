package logging

import (
	"errors"
	"sync"

	"github.com/hashicorp/go-hclog"

	"github.com/raito-io/cli/internal/target"
)

const bufferSize = 1024

type WarningCollector interface {
	GetWarnings() []string
}

func newWarningCapturingSink() *warningCapturingSink {
	return &warningCapturingSink{
		warnings: make([]string, 0, bufferSize),
	}
}

type warningCapturingSink struct {
	mutex    sync.Mutex
	warnings []string
}

func (s *warningCapturingSink) Accept(name string, level hclog.Level, msg string, args ...interface{}) {
	if level == hclog.Warn {
		s.mutex.Lock()
		defer s.mutex.Unlock()

		s.warnings = append(s.warnings, msg)
	}
}

func (s *warningCapturingSink) GetWarnings() []string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var warnings []string
	if len(s.warnings) > 0 {
		warnings = append(warnings, s.warnings...)
		s.warnings = make([]string, 0, bufferSize)
	}

	return warnings
}

func CreateWarningCapturingLogger(config *target.BaseTargetConfig) (*target.BaseTargetConfig, WarningCollector, func(), error) {
	if logger, ok := config.BaseLogger.(hclog.InterceptLogger); ok {
		sink := newWarningCapturingSink()

		logger.RegisterSink(sink)

		return config, sink, func() {
			logger.DeregisterSink(sink)
		}, nil
	}

	return nil, nil, nil, errors.New("no logger found")
}
