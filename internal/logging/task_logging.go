package logging

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hasura/go-graphql-client"

	"github.com/raito-io/cli/internal/file"
	gql "github.com/raito-io/cli/internal/graphql"
	"github.com/raito-io/cli/internal/target"
)

func newTaskFileSink(config *target.BaseTargetConfig, jobId, taskId string) (*taskFileSink, error) {
	tmpFile, err := os.CreateTemp("", "*")
	if err != nil {
		return nil, err
	}

	return &taskFileSink{
		jobId:  jobId,
		taskId: taskId,
		config: config,
		writer: tmpFile,
	}, nil
}

type taskFileSink struct {
	jobId  string
	taskId string
	config *target.BaseTargetConfig
	writer *os.File
}

func (s *taskFileSink) Accept(name string, level hclog.Level, msg string, args ...interface{}) {
	var argsBuilder strings.Builder

	for i, arg := range args {
		if i%2 == 0 {
			argsBuilder.WriteString(fmt.Sprintf(" %v=", arg))
		} else {
			argsBuilder.WriteString(fmt.Sprintf("%v", arg))
		}
	}

	s.writer.WriteString(fmt.Sprintf("%s [%s] %s: %s:%s\n", time.Now().Format("2006-01-02T15:04:05.000-0700"), level.String(), name, msg, argsBuilder.String())) //nolint:errcheck
}

func (s *taskFileSink) Close() error {
	err := s.writer.Close()
	if err != nil {
		return err
	}

	key, err := file.UploadLogFile(s.writer.Name(), s.config, s.taskId)
	if err != nil {
		return err
	}

	var mutation struct {
		AddLogFileToTask struct {
			JobId   string
			JobType string
		} `graphql:"addLogFileToTask(input: $input)"`
	}

	variables := map[string]interface{}{
		"input": LogFileToTaskInput{
			JobId:   graphql.ID(s.jobId),
			JobType: s.taskId,
			FileKey: key,
		},
	}

	err = gql.NewClient(&s.config.BaseConfig).Mutate(context.Background(), &mutation, variables)
	if err != nil {
		s.config.BaseLogger.Warn(fmt.Sprintf("Error: %s", err))

		return err
	}

	return os.Remove(s.writer.Name())
}

func CreateTaskLogger(config *target.BaseTargetConfig, jobId string, taskId string) (*target.BaseTargetConfig, func() error, error) {
	if logger, ok := config.BaseLogger.(hclog.InterceptLogger); ok {
		sink, err := newTaskFileSink(config, jobId, taskId)
		if err != nil {
			return nil, nil, err
		}

		logger.RegisterSink(sink)

		return config, func() error {
			logger.DeregisterSink(sink)

			return sink.Close()
		}, nil
	}

	return nil, nil, errors.New("no logger found")
}

type LogFileToTaskInput struct {
	JobId   graphql.ID `json:"jobId"`
	JobType string     `json:"jobType"`
	FileKey string     `json:"fileKey"`
}
