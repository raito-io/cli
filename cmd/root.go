package cmd

import (
	_ "embed"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/logging"
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
		Use:   "raito",
		Short: "The Raito CLI to take care of all your data access needs.",
		Long:  rootDescription,
		// TODO complete the long description
		Version: version,
	}

	cobra.OnInitialize(root.initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, constants.ConfigFileFlag, "", "The config file (default is ./raito.yml).")
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
		fmt.Printf("error while hiding dev flag.\n") //nolint:forbidigo
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
	if err != nil && !strings.HasSuffix(viper.ConfigFileUsed(), "/raito") {
		// No logger yet
		fmt.Printf("error while reading config file %q: %s\n", viper.ConfigFileUsed(), err.Error()) //nolint:forbidigo
		cmd.exitForError()
	}

	logging.SetupLogging()

	if viper.ConfigFileUsed() != "" {
		hclog.L().Debug(fmt.Sprintf("Using config file: %s", viper.ConfigFileUsed()))
	}
}
