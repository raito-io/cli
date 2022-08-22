package cmd

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/hashicorp/go-hclog"
	dapc "github.com/raito-io/cli/common/api/data_access"
	"github.com/raito-io/cli/common/util/config"
	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/data_access"
	"github.com/raito-io/cli/internal/plugin"
	"github.com/raito-io/cli/internal/target"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	hclog.L().Info("")
	defer hclog.L().Info("")

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

	config := dapc.DataAccessSyncConfig{
		ConfigMap: config.ConfigMap{Parameters: targetConfig.Parameters},
		Prefix:    "R",
		RunImport: false,
	}

	af, err := os.Open(accessFile)
	if err != nil {
		targetConfig.Logger.Error(fmt.Sprintf("Error while opening data access file %q: %s", accessFile, err.Error()))
		return err
	}

	buf, err := io.ReadAll(af)
	if err != nil {
		targetConfig.Logger.Error(fmt.Sprintf("Error while reading data access file %q: %s", accessFile, err.Error()))
		return err
	}

	dar, err := data_access.ParseDataAccess(buf)
	if err != nil {
		targetConfig.Logger.Error(fmt.Sprintf("Error while parsing data access file %q: %s", accessFile, err.Error()))
		return err
	}
	config.DataAccess = dar

	das, err := client.GetDataAccessSyncer()
	if err != nil {
		targetConfig.Logger.Error(fmt.Sprintf("The plugin (%s) does not implement the DataAccessSyncer interface", targetConfig.ConnectorName))
		return err
	}

	res := das.SyncDataAccess(&config)
	if res.Error != nil {
		target.HandleTargetError(res.Error, targetConfig, "sychronizing data access information to the data source")
		return err
	}

	sec := time.Since(start).Round(time.Millisecond)
	targetConfig.Logger.Info(fmt.Sprintf("Finished execution in %s", sec), "success")

	return nil
}
