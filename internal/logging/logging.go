package logging

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hasura/go-graphql-client"
	"github.com/pterm/pterm"
	"github.com/spf13/viper"

	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/file"
	gql "github.com/raito-io/cli/internal/graphql"
	"github.com/raito-io/cli/internal/target"
	"github.com/raito-io/cli/internal/version"
)

func SetupLogging() {
	var output io.Writer = &nilWriter{}

	if viper.GetBool(constants.LogOutputFlag) {
		output = os.Stderr

		if viper.GetString(constants.LogFileFlag) != "" {
			f, err := os.OpenFile(viper.GetString(constants.LogFileFlag), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
			if err != nil {
				fmt.Printf("error opening file: %v", err) //nolint:forbidigo
			}
			output = io.MultiWriter(f, os.Stderr)
		}
	} else if viper.GetString(constants.LogFileFlag) != "" {
		f, err := os.OpenFile(viper.GetString(constants.LogFileFlag), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			fmt.Printf("error opening file: %v", err) //nolint:forbidigo
		}
		output = f
	}

	logger := hclog.NewInterceptLogger(&hclog.LoggerOptions{
		Name:   fmt.Sprintf("raito-cli-%s)", version.GetCliVersion().String()),
		Output: output,
	})

	if !viper.GetBool(constants.LogOutputFlag) {
		logger.RegisterSink(newSinkAdapter())
	}

	if viper.GetBool(constants.DebugFlag) {
		logger.SetLevel(hclog.Debug)
	} else {
		logger.SetLevel(hclog.Info)
	}

	// log the standard logger from 'import "log"'
	log.SetOutput(logger.StandardWriter(&hclog.StandardLoggerOptions{InferLevels: true}))

	hclog.SetDefault(logger)
}

type nilWriter struct {
}

func (w *nilWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

type sinkAdapter struct {
	iteration    int
	progress     map[string]*pterm.SpinnerPrinter
	wasIteration bool
	mu           sync.Mutex
}

func newSinkAdapter() *sinkAdapter {
	sa := &sinkAdapter{}
	sa.iteration = -1

	return sa
}

func (s *sinkAdapter) Accept(name string, level hclog.Level, msg string, args ...interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	it, tar := getIterationAndTarget(args)
	if it >= 0 {
		s.wasIteration = true
		if it != s.iteration {
			if it != 0 {
				pterm.Println(fmt.Sprintf("Starting iteration %d", it))
			}

			s.startNewIteration()
			s.iteration = it
		}

		if tar != "" {
			spinnerKey := strconv.Itoa(it) + "-" + tar

			spinner := s.progress[spinnerKey]
			if spinner == nil {
				spinner, _ = pterm.DefaultSpinner.Start(fmt.Sprintf("Running target %s...", tar))
				s.progress[spinnerKey] = spinner

				// TODO this is to avoid a threading issue with pterm? When not done, fast targets (e.g. when skipped) give weird results in the CLI output (spinner appearing again)
				time.Sleep(500 * time.Millisecond)
			}

			s.handleProgress(spinner, tar, level, msg, args)
		} else {
			s.handleNormalOutput(level, msg)
		}
	} else {
		// Extra line break if we came from an iteration
		if s.wasIteration {
			s.stopIteration()
			pterm.Println()
			s.wasIteration = false
		}
		s.handleNormalOutput(level, msg)
	}
}

func (s *sinkAdapter) handleNormalOutput(level hclog.Level, msg string) {
	if level == hclog.Error {
		pterm.Error.Println(msg)
	} else if level == hclog.Warn {
		pterm.Warning.Println(msg)
	} else if level == hclog.Info {
		pterm.Println(msg)
	}
}

func (s *sinkAdapter) handleProgress(spinner *pterm.SpinnerPrinter, target string, level hclog.Level, msg string, args []interface{}) {
	text := fmt.Sprintf("Target %s - %s", target, msg)

	if level == hclog.Info {
		if s.hasSuccess(args) {
			spinner.Success(text)
		} else {
			spinner.UpdateText(text)
		}
	} else if level == hclog.Error {
		if s.hasSuccess(args) {
			spinner.Fail(text)
		} else {
			pterm.Error.Println(text)
		}
	} else if level == hclog.Warn {
		pterm.Warning.Println(text)
	}
}

func (s *sinkAdapter) hasSuccess(args []interface{}) bool {
	for _, arg := range args {
		if arg == "success" {
			return true
		}
	}

	return false
}

func (s *sinkAdapter) stopIteration() {
	if s.progress != nil {
		for _, spinner := range s.progress {
			if spinner.IsActive {
				spinner.Stop() //nolint
			}
		}
	}
}

func (s *sinkAdapter) startNewIteration() {
	s.stopIteration()
	s.progress = make(map[string]*pterm.SpinnerPrinter)
}

func getIterationAndTarget(args []interface{}) (int, string) {
	iterationFound := false
	targetFound := false

	iteration := -1
	target := ""

	for _, arg := range args {
		if iterationFound {
			iteration = arg.(int)
			iterationFound = false
		}

		if targetFound {
			target = arg.(string)
			targetFound = false
		}

		if arg == "iteration" {
			iterationFound = true
		}

		if arg == "target" {
			targetFound = true
		}
	}

	return iteration, target
}

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

	key, err := file.UploadLogFile(s.writer.Name(), s.config)
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
