package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/pterm/pterm"
	"github.com/spf13/viper"

	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/version"
)

func SetupLogging() {
	var output io.Writer = &nilWriter{}

	if viper.GetBool(constants.LogOutputFlag) {
		output = os.Stdout

		if viper.GetString(constants.LogFileFlag) != "" {
			f, err := os.OpenFile(viper.GetString(constants.LogFileFlag), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
			if err != nil {
				fmt.Printf("error opening file: %v", err) //nolint:forbidigo
			}
			output = io.MultiWriter(f, output)
		}
	} else if viper.GetString(constants.LogFileFlag) != "" {
		f, err := os.OpenFile(viper.GetString(constants.LogFileFlag), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			fmt.Printf("error opening file: %v", err) //nolint:forbidigo
		}
		output = f
	}

	logger := hclog.NewInterceptLogger(&hclog.LoggerOptions{
		Name:   fmt.Sprintf("raito-cli-%s", version.GetCliVersion().String()),
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

			text := fmt.Sprintf("Target %s - %s", tar, msg)

			if level == hclog.Info {
				// During the run of a target, we show info messages in a spinner
				if spinner == nil {
					spinner, _ = pterm.DefaultSpinner.Start(fmt.Sprintf("Running target %s...", tar))
					spinner.RemoveWhenDone = true
					s.progress[spinnerKey] = spinner
				}

				if s.hasSuccess(args) {
					// Target success message, so we stop the spinner and show the success message separately to keep it visible
					//nolint:errcheck
					spinner.Stop()
					s.progress[spinnerKey] = nil

					pterm.Success.Println(text)
				} else {
					// Normal info message during a target run, so just update the text of the spinner
					spinner.Info(text)
				}
			} else if (level == hclog.Error || level == hclog.Warn) && spinner != nil {
				// If we encounter an error or warning, we stop the current spinner if there is any and remove it.
				//nolint:errcheck
				spinner.Stop()
				s.progress[spinnerKey] = nil
			}

			// Show error or warning messages separately
			if level == hclog.Error {
				pterm.Error.Println(text)
			} else if level == hclog.Warn {
				pterm.Warning.Println(text)
			}
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
			if spinner != nil && spinner.IsActive {
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
