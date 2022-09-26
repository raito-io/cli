package cmd

import (
	_ "embed"
	"fmt"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/raito-io/cli/base/access_provider"
	"github.com/raito-io/cli/base/util/config"
	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/plugin"
	"github.com/raito-io/cli/internal/target"
)

const defaultAccessFile = "access.yml"

//go:embed help/access-description.txt
var accessDescription string

func initAccessCommand(rootCmd *cobra.Command) {
	var cmd = &cobra.Command{
		Use:   "access",
		Short: "Update the access permissions of the target with information from a YAML file.",
		Long:  accessDescription,
		RunE:  executeAccessCmd,
	}

	cmd.PersistentFlags().StringP(constants.AccessFileFlag, "a", "", fmt.Sprintf("Use this to specify a custom file path to use for the location of the access definition file. Default is %q. This can also be specified under the target in the configuration file.", defaultAccessFile))

	BindFlag(constants.AccessFileFlag, cmd)

	rootCmd.AddCommand(cmd)
}

func executeAccessCmd(cmd *cobra.Command, args []string) error {
	baseLogger := hclog.L().With("iteration", 0)

	return target.RunTargets(baseLogger, cmd.Flags().Args(), runAccessTarget)
}

func runAccessTarget(targetConfig *target.BaseTargetConfig) error {
	start := time.Now()

	accessFile := viper.GetString(constants.AccessFileFlag)
	if accessFile == "" {
		if fp, ok := targetConfig.Parameters[constants.AccessFileFlag]; ok {
			if af, ok := fp.(string); ok {
				accessFile = af
			}
		}
	}

	if accessFile == "" {
		accessFile = defaultAccessFile
	}

	client, err := plugin.NewPluginClient(targetConfig.ConnectorName, targetConfig.ConnectorVersion, targetConfig.Logger)
	if err != nil {
		targetConfig.Logger.Error(fmt.Sprintf("Error initializing connector plugin %q: %s", targetConfig.ConnectorName, err.Error()))
		return err
	}
	defer client.Close()

	as, err := client.GetAccessSyncer()
	if err != nil {
		targetConfig.Logger.Error(fmt.Sprintf("The plugin (%s) does not implement the AccessSyncer interface", targetConfig.ConnectorName))
		return err
	}

	res := as.SyncToTarget(&access_provider.AccessSyncToTarget{
		ConfigMap:  config.ConfigMap{Parameters: targetConfig.Parameters},
		Prefix:     "R",
		SourceFile: accessFile,
	})
	if res.Error != nil {
		target.HandleTargetError(res.Error, targetConfig, "synchronizing access information to the data source")
		return err
	}

	sec := time.Since(start).Round(time.Millisecond)
	targetConfig.Logger.Info(fmt.Sprintf("Finished execution in %s", sec), "success")

	return nil
}
