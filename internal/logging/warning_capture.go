package logging

import (
	"sync"

	"github.com/hashicorp/go-hclog"
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

func CreateWarningCapturingLogger(baseLogger hclog.InterceptLogger) (WarningCollector, func()) {
	sink := newWarningCapturingSink()

	baseLogger.RegisterSink(sink)

	return sink, func() {
		baseLogger.DeregisterSink(sink)
	}
}
