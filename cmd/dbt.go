package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/dbt"
	"github.com/raito-io/cli/internal/target"
	"github.com/raito-io/cli/internal/target/types"
)

func initDbtCommand(rootCmd *cobra.Command) {
	var cmd = &cobra.Command{
		Hidden: true,
		Use:    "dbt",
		Short:  "Run dbt integration",
		Long:   "Run dbt integration",
		Run:    executeDbt,
	}

	cmd.PersistentFlags().StringP(constants.DbtManifestFile, "m", "", "Path to the dbt manifest file")
	BindFlag(constants.DbtManifestFile, cmd)

	rootCmd.AddCommand(cmd)
}

func executeDbt(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	config, err := buildConfig(cmd)
	if err != nil {
		hclog.L().Error(err.Error())
		os.Exit(1)
	}

	dbtService := dbt.NewDbtService(ctx, config)

	updated, failed, err := dbtService.RunDbt(ctx, config.DbtFilePath)
	if err != nil {
		hclog.L().Error(fmt.Sprintf("Error during dbt integration: %s. Updated %d access providers. Failed to update %d access providers", err.Error(), updated, failed))
		os.Exit(1)
	}

	hclog.L().Info(fmt.Sprintf("Successfully finished dbt integration. Updated %d access providers", updated))
}

func buildConfig(cmd *cobra.Command) (*types.DbtConfig, error) {
	otherArgs := cmd.Flags().Args()

	baseConfig, err := target.BuildBaseConfigFromFlags(hclog.L(), otherArgs)
	if err != nil {
		return nil, err
	}

	dbtConfig := types.DbtConfig{
		BaseConfig:   *baseConfig,
		DbtFilePath:  viper.GetString(constants.DbtManifestFile),
		DataSourceId: viper.GetString(constants.DataSourceIdFlag),
	}

	return &dbtConfig, nil
}
