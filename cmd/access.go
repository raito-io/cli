package cmd

import (
	"context"
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

	config, err := target.BuildBaseConfigFromFlags(baseLogger, cmd.Flags().Args())
	if err != nil {
		return err
	}

	return target.RunTargets(config, runAccessTarget)
}

func runAccessTarget(targetConfig *target.BaseTargetConfig) error {
	start := time.Now()

	accessFile := viper.GetString(constants.AccessFileFlag)
	if accessFile == "" {
		if fp, ok := targetConfig.Parameters[constants.AccessFileFlag]; ok {
			accessFile = fp
		}
	}

	if accessFile == "" {
		accessFile = defaultAccessFile
	}

	client, err := plugin.NewPluginClient(targetConfig.ConnectorName, targetConfig.ConnectorVersion, targetConfig.TargetLogger)
	if err != nil {
		targetConfig.TargetLogger.Error(fmt.Sprintf("Error initializing connector plugin %q: %s", targetConfig.ConnectorName, err.Error()))
		return err
	}
	defer client.Close()

	as, err := client.GetAccessSyncer()
	if err != nil {
		targetConfig.TargetLogger.Error(fmt.Sprintf("The plugin (%s) does not implement the AccessSyncer interface", targetConfig.ConnectorName))
		return err
	}

	res, err := as.SyncToTarget(context.Background(), &access_provider.AccessSyncToTarget{
		ConfigMap:  &config.ConfigMap{Parameters: targetConfig.Parameters},
		Prefix:     "R",
		SourceFile: accessFile,
	})
	if err != nil {
		return err
	} else if res.Error != nil {
		target.HandleTargetError(res.Error, targetConfig, "synchronizing access information to the data source")
		return err
	}

	sec := time.Since(start).Round(time.Millisecond)
	targetConfig.TargetLogger.Info(fmt.Sprintf("Finished execution in %s", sec), "success")

	return nil
}
