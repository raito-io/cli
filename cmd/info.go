package cmd

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/hashicorp/go-hclog"
	"github.com/pterm/pterm"
	"github.com/raito-io/cli/internal/logging"
	"github.com/spf13/cobra"

	"github.com/raito-io/cli/internal/plugin"
)

func initInfoCommand(rootCmd *cobra.Command) {
	var cmd = &cobra.Command{
		Short:     "Retrieve information about the given connector.",
		Long:      "Retrieve information about the given connector. You have the option to specify a specific connector version. If not, 'latest' is assumed.",
		Run:       executeInfoCmd,
		ValidArgs: []string{},
		Use:       "info <connector> [<version>]",
	}

	rootCmd.AddCommand(cmd)
}

func executeInfoCmd(cmd *cobra.Command, args []string) {
	logging.SetupLogging(false)

	if len(args) < 1 || len(args) > 2 {
		pterm.Error.Println("Expected 1 or 2 arguments after the info command.")
		return
	}
	connector := args[0]
	version := ""

	if len(args) > 1 {
		version = args[1]
	}

	client, err := plugin.NewPluginClient(connector, version, hclog.L())
	if err != nil {
		return
	}
	defer client.Close()

	info, err := client.GetInfo()
	if err != nil {
		pterm.Warning.Println(fmt.Sprintf("The plugin (%s) does not implement the Info interface. Skipping.", connector))
		return
	}

	pluginInfo, err := info.GetInfo(context.Background())
	if err != nil {
		pterm.Error.Println(fmt.Sprintf("Failed to load plugin info: %s", err))
		return
	}

	pterm.Println("Plugin name: " + pterm.Bold.Sprint(pluginInfo.Name))

	v := pluginInfo.GetVersion()
	sv := semver.New(v.GetMajor(), v.GetMinor(), v.GetPatch(), v.GetPrerelease(), v.GetBuild())
	pterm.Println("Version: " + pterm.Bold.Sprint(sv.String()))

	if pluginInfo.Description != "" {
		pterm.Println()
		pterm.Println(pluginInfo.Description)
	}

	if len(pluginInfo.Parameters) > 0 {
		pterm.Println()
		pterm.Println("Parameters:")

		for _, param := range pluginInfo.Parameters {
			line := "   " + pterm.Bold.Sprintf(param.Name)
			if param.Mandatory {
				line += " (required)"
			}
			line += ": " + param.Description

			pterm.Println(line)
		}
	}

	pterm.Println()
}
