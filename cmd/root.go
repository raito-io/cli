package cmd

import (
	_ "embed"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/pterm/pterm"
	"github.com/raito-io/cli/internal/constants"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

//go:embed help/root-description.txt
var rootDescription string

var (
	logger  hclog.InterceptLogger
	cfgFile string
)

type rootCmd struct {
	cmd  *cobra.Command
	exit func(int)
}

func (cmd *rootCmd) Execute(args []string) {
	cmd.cmd.SetArgs(args)
	if err := cmd.cmd.Execute(); err != nil {
		cmd.exitForError(err)
	}
}

func (cmd *rootCmd) exitForError(err error) {
	cmd.exit(1)
}

func newRootCmd(version string, exit func(int)) *rootCmd {
	root := &rootCmd{
		exit: exit,
	}
	rootCmd := &cobra.Command{
		Use:   "raito",
		Short: "The Raito CLI to take care of all your data access needs.",
		Long:  rootDescription,
		// TODO complete the long description
		Version: version,
	}
	cobra.OnInitialize((*root).initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, constants.ConfigFileFlag, "", "The config file (default is ./raito.yml)")
	rootCmd.PersistentFlags().String(constants.IdentityStoreIdFlag, "", "The ID of the identity store in Raito to import the user and group information to. This is only applicable if specifying the (single) target information in the commandline.")
	rootCmd.PersistentFlags().String(constants.DataSourceIdFlag, "", "The ID of the data source in Raito to import the meta data information to or get the access permissions from. This is only applicable if specifying the (single) target information in the commandline.")
	rootCmd.PersistentFlags().StringP(constants.DomainFlag, "d", "", "The subdomain to your Raito instance (https://<subdomain>.raito.io).")
	rootCmd.PersistentFlags().StringP(constants.ApiUserFlag, "u", "", "The username of the API user to authenticate against Raito.")
	rootCmd.PersistentFlags().StringP(constants.ApiSecretFlag, "s", "", "The API key secret to authenticate against Raito.")
	rootCmd.PersistentFlags().String(constants.LogFileFlag, "", "The log file to store structured logs in. If not specified, no logging to file is done.")
	rootCmd.PersistentFlags().Bool(constants.LogOutputFlag, false, "When set, logging is sent to the command line (stderr) instead of more human readable output.")
	rootCmd.PersistentFlags().String(constants.EnvironmentFlag, constants.EnvironmentProd, "")
	rootCmd.PersistentFlags().Bool(constants.DebugFlag, false, fmt.Sprintf("If set, extra debug logging is generated. Only useful in combination with %s or %s", constants.LogFileFlag, constants.LogOutputFlag))
	rootCmd.PersistentFlags().StringP(constants.OnlyTargetsFlag, "t", "", "Can be used to only execute a subset of the defined targets in the configuration file. To specify multiple, use a comma-separated list.")
	rootCmd.PersistentFlags().String(constants.ConnectorNameFlag, "", "The name of the connector to use. If not set, the CLI will use a configuration file to define the targets.")
	rootCmd.PersistentFlags().String(constants.ConnectorVersionFlag, "", "The version of the connector to use. This is only relevant if the 'connector' flag is set as well. If not set (but the 'connector' flag is), then 'latest' is used.")
	rootCmd.PersistentFlags().StringP(constants.NameFlag, "n", "", "The name for the target. This is only relevant if the 'connector' flag is set as well. If not set, the name of the connector will be used.")
	err := rootCmd.PersistentFlags().MarkHidden(constants.EnvironmentFlag)
	if err != nil {
		// No logger yet
		fmt.Printf("error while hiding dev flag.\n")
	}

	BindFlag(constants.IdentityStoreIdFlag, rootCmd)
	BindFlag(constants.DataSourceIdFlag, rootCmd)
	BindFlag(constants.OnlyTargetsFlag, rootCmd)
	BindFlag(constants.ConnectorNameFlag, rootCmd)
	BindFlag(constants.ConnectorVersionFlag, rootCmd)
	BindFlag(constants.NameFlag, rootCmd)
	BindFlag(constants.ConfigFileFlag, rootCmd)
	BindFlag(constants.DomainFlag, rootCmd)
	BindFlag(constants.ApiUserFlag, rootCmd)
	BindFlag(constants.ApiSecretFlag, rootCmd)
	BindFlag(constants.EnvironmentFlag, rootCmd)
	BindFlag(constants.DebugFlag, rootCmd)
	BindFlag(constants.LogFileFlag, rootCmd)
	BindFlag(constants.LogOutputFlag, rootCmd)

	viper.SetDefault(constants.LogFileFlag, "")

	root.cmd = rootCmd

	initRunCommand(rootCmd)
	initAccessCommand(rootCmd)
	initInfoCommand(rootCmd)

	return root
}

func BindFlag(flag string, cmd *cobra.Command) {
	err := viper.BindPFlag(flag, cmd.PersistentFlags().Lookup(flag))
	if err != nil {
		// No logger yet
		fmt.Printf("error while binding flag %q: %s\n", flag, err.Error())
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string, args []string, exit func(int)) {
	newRootCmd(version, exit).Execute(args)
}

func (cmd *rootCmd) initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.raito")
		viper.SetConfigType("yaml")
		viper.SetConfigName("raito")
	}
	viper.SetEnvPrefix("RAITO")
	replacer := strings.NewReplacer(".", "_", "-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if cfgFile != "" && err != nil {
		// No logger yet
		fmt.Printf("error while reading config file: %s\n", err.Error())
		cmd.exitForError(err)
	}

	setupLogging()
}

// TODO move all logging stuff to a separate package (under pkg)
func setupLogging() {
	var output io.Writer = &nilWriter{}

	if viper.GetBool(constants.LogOutputFlag) {
		output = os.Stderr
		if viper.GetString(constants.LogFileFlag) != "" {
			f, err := os.OpenFile(viper.GetString(constants.LogFileFlag), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
			if err != nil {
				fmt.Printf("error opening file: %v", err)
			}
			output = io.MultiWriter(f, os.Stderr)
		}
	} else {
		if viper.GetString(constants.LogFileFlag) != "" {
			f, err := os.OpenFile(viper.GetString(constants.LogFileFlag), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
			if err != nil {
				fmt.Printf("error opening file: %v", err)
			}
			output = f
		}
	}

	logger = hclog.NewInterceptLogger(&hclog.LoggerOptions{
		Name:   "raito-cli",
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
	it, tar := s.getIterationAndTarget(args)
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
			s.handleProgress(spinner, it, tar, level, msg, args)
		} else {
			s.handleNormalOutput(name, level, msg, args)
		}
	} else {
		// Extra line break if we came from an iteration
		if s.wasIteration {
			s.stopIteration()
			pterm.Println()
			s.wasIteration = false
		}
		s.handleNormalOutput(name, level, msg, args)
	}
}

func (s *sinkAdapter) handleNormalOutput(name string, level hclog.Level, msg string, args []interface{}) {
	if level == hclog.Error {
		pterm.Error.Println(msg)
	} else if level == hclog.Info {
		pterm.Println(msg)
	}
}

func (s *sinkAdapter) handleProgress(spinner *pterm.SpinnerPrinter, iteration int, target string, level hclog.Level, msg string, args []interface{}) {
	text := fmt.Sprintf("Target %s - %s", target, msg)
	if level == hclog.Info {
		if s.hasSuccess(args) {
			spinner.Success(text)
		} else {
			spinner.UpdateText(text)
		}
	} else if level == hclog.Error {
		spinner.Fail(text)
	} else if level == hclog.Warn {
		spinner.UpdateText("Warning: " + text)
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

func (s *sinkAdapter) getIterationAndTarget(args []interface{}) (int, string) {
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
