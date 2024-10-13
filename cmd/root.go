package cmd

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/raito-io/cli/internal/constants"
)

//go:embed help/root-description.txt
var rootDescription string

var (
	cfgFile string
)

type rootCmd struct {
	cmd  *cobra.Command
	exit func(int)
}

func (cmd *rootCmd) Execute(args []string) {
	cmd.cmd.SetArgs(args)

	if err := cmd.cmd.Execute(); err != nil {
		cmd.exitForError()
	}
}

func (cmd *rootCmd) exitForError() {
	cmd.exit(1)
}

func newRootCmd(version string, exit func(int)) *rootCmd {
	root := &rootCmd{
		exit: exit,
	}
	rootCmd := &cobra.Command{
		Use:     "raito",
		Short:   "The Raito CLI to take care of all your data access needs.",
		Long:    rootDescription,
		Version: version,
	}

	cobra.OnInitialize(root.initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, constants.ConfigFileFlag, "", "The config file (default is ./raito.yml).")
	rootCmd.PersistentFlags().String(constants.LogFileFlag, "", "The log file to store structured logs in. If not specified, no logging to file is done.")
	rootCmd.PersistentFlags().Bool(constants.LogOutputFlag, false, "When set, logging is sent to the command line (stderr) instead of more human readable output.")
	rootCmd.PersistentFlags().Bool(constants.DebugFlag, false, fmt.Sprintf("If set, extra debug logging is generated. Only useful in combination with %s or %s", constants.LogFileFlag, constants.LogOutputFlag))

	BindFlag(constants.ConfigFileFlag, rootCmd)
	BindFlag(constants.DebugFlag, rootCmd)
	BindFlag(constants.LogFileFlag, rootCmd)
	BindFlag(constants.LogOutputFlag, rootCmd)

	viper.SetDefault(constants.LogFileFlag, "")

	root.cmd = rootCmd

	initRunCommand(rootCmd)
	initInfoCommand(rootCmd)
	initApplyAccessCommand(rootCmd)
	initInitCommand(rootCmd)
	initAddTargetCommand(rootCmd)

	return root
}

func hideConfigOptions(rootCmd *cobra.Command, options ...string) {
	for _, option := range options {
		err := rootCmd.PersistentFlags().MarkHidden(option)
		if err != nil {
			// No logger yet
			fmt.Printf("error while hiding %s flag.\n", option) //nolint:forbidigo
		}
	}
}

func BindFlag(flag string, cmd *cobra.Command) {
	err := viper.BindPFlag(flag, cmd.PersistentFlags().Lookup(flag))
	if err != nil {
		// No logger yet
		fmt.Printf("error while binding flag %q: %s\n", flag, err.Error()) //nolint:forbidigo
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
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.raito")
		viper.SetConfigName("raito")
	}

	viper.SetEnvPrefix("RAITO")
	replacer := strings.NewReplacer(".", "_", "-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	chosenCommand := ""

	if len(os.Args) > 1 {
		chosenCommand = os.Args[1]
	}

	if !strings.EqualFold(chosenCommand, "init") && !strings.EqualFold(chosenCommand, "add-target") {
		if err != nil && !strings.HasSuffix(viper.ConfigFileUsed(), "/raito") && viper.ConfigFileUsed() != "" {
			// No logger yet
			fmt.Printf("error while reading config file %q: %s\n", viper.ConfigFileUsed(), err.Error()) //nolint:forbidigo
			cmd.exitForError()
		}

		viper.WatchConfig()

		if viper.ConfigFileUsed() != "" {
			hclog.L().Debug(fmt.Sprintf("Using config file: %s", viper.ConfigFileUsed()))
		}
	}
}
