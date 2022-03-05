package cmd

import (
	"fmt"
	"github.com/raito-io/cli/internal/plugin"
	"github.com/spf13/cobra"
)

func initInfoCommand(rootCmd *cobra.Command) {
	var cmd = &cobra.Command{
		Short: "Retrieve information about the given connector.",
		Long: "Retrieve information about the given connector. You have the option to specify a specific connector version. If not, 'latest' is assumed.",
		Run: executeInfoCmd,
		ValidArgs: []string {  },
		Use: "info <connector> [<version>]",
	}

	rootCmd.AddCommand(cmd)
}

func executeInfoCmd(cmd *cobra.Command, args []string) {
	if len(args) < 1 || len(args) > 2 {
		logger.Error("Expected 1 or 2 arguments after the info command.")
		return
	}
	connector := args[0]
	version := ""
	if len(args) > 1 {
		version = args[1]
	}

	client, err := plugin.NewPluginClient(connector, version, logger)
	if err != nil {
		return
	}
	defer client.Close()

	info, err := client.GetInfo()
	if err != nil {
		logger.Warn(fmt.Sprintf("The plugin (%s) does not implement the Info interface. Skipping.", connector))
		return
	}

	logger.Info(info.PluginInfo().FullOverview())
}
