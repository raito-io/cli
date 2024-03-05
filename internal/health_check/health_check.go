package health_check

import (
	"os"

	"github.com/hashicorp/go-hclog"
)

type HealthChecker struct {
	logger           hclog.Logger
	livenessFilePath string

	livenessFile *os.File
}

func NewHealthChecker(logger hclog.Logger, livenessFilePath string) HealthChecker {
	return HealthChecker{
		logger:           logger,
		livenessFilePath: livenessFilePath,
	}
}

func NewDummyHealthChecker(logger hclog.Logger) HealthChecker {
	return HealthChecker{
		logger:           logger,
		livenessFilePath: "",
	}
}

func (s *HealthChecker) MarkLiveness() error {
	if s.livenessFilePath == "" || s.livenessFile != nil {
		return nil
	}

	livenessFile, _, err := s.createOutputFile(s.livenessFilePath)
	if err != nil {
		return err
	}

	s.logger.Debug("[Healthchecker] Created liveness file: %s", s.livenessFilePath)

	s.livenessFile = livenessFile

	return nil
}

func (s *HealthChecker) RemoveLivenessMark() error {
	if s.livenessFile != nil {
		err := s.livenessFile.Close()

		if err != nil {
			s.logger.Warn("failed to close file: %v", err)
			return err
		}

		err = os.Remove(s.livenessFilePath)
		if err != nil {
			s.logger.Warn("failed to remove file: %v", err)
			return err
		}

		s.livenessFile = nil
	}

	return nil
}

func (s *HealthChecker) Cleanup() {
	_ = s.RemoveLivenessMark()
}

func (s *HealthChecker) createOutputFile(filename string) (*os.File, func() error, error) {
	if filename == os.Stdout.Name() {
		return os.Stdout, func() error { return nil }, nil
	} else if filename == os.Stderr.Name() {
		return os.Stderr, func() error { return nil }, nil
	}

	outputFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return nil, nil, err
	}

	return outputFile, func() error {
		return outputFile.Close()
	}, nil
}
