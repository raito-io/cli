package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/pterm/pterm"
	dapc "github.com/raito-io/cli/base/access_provider"
	baseconfig "github.com/raito-io/cli/base/util/config"
	"github.com/raito-io/cli/internal/health_check"
	"github.com/raito-io/cli/internal/logging"
	"github.com/raito-io/cli/internal/plugin"
	"github.com/raito-io/cli/internal/target"
	"github.com/raito-io/cli/internal/target/types"
	"github.com/raito-io/cli/internal/util/file"
	"github.com/raito-io/cli/internal/version_management"
	"github.com/spf13/cobra"
)

func initApplyAccessCommand(rootCmd *cobra.Command) {
	var cmd = &cobra.Command{
		Short:     "(Re)apply access to the given target from a specific file.",
		Long:      "This can be used to (re)apply a given access file to a specific target. This is a very dangerous operation and should be used with caution.",
		Run:       executeApplyAccessCmd,
		ValidArgs: []string{},
		Use:       "apply-access <target-name> <file-path>",
	}

	rootCmd.AddCommand(cmd)
}

var warningMessage = `This action is highly discouraged and should only be used in emergency situations!
Due to the different way that each data source implements access controls, each connector has different capabilities and behaviour to apply them. 
This means that some connectors depend on explicit delete operations (for access controls, who-items and/or what-items).
This action will simply re-run the connector code, based on the provided file. There will be no communication with Raito Cloud (so no job registered and no feedback provided about the actions taken).

This means that the following things may happen, depending on the connector and situation:
 - explicitly deleted items (access controls, who-items and/or what-items) at that point in time, will be deleted again. During the next full sync, these will potentially be imported again (as ‘Managed on the data source’).
 - access that may have been removed already again in the meanwhile, could be re-applied
So make sure you understand what that means for your case, before you continue with this.`

func executeApplyAccessCmd(cmd *cobra.Command, args []string) {
	logging.SetupLogging(true)

	err := applyAccessCmd(cmd, args)
	if err != nil {
		pterm.Error.Println(err.Error())
		os.Exit(1)
	}
}

func applyAccessCmd(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return errors.New("Expected 2 arguments after the apply-access command.") //nolint:stylecheck
	}

	otherArgs := cmd.Flags().Args()

	baseLogger := hclog.L()

	baseConfig, err := target.BuildBaseConfigFromFlags(baseLogger, health_check.NewDummyHealthChecker(baseLogger), otherArgs)
	if err != nil {
		return fmt.Errorf("An error occurred while parsing the configuration file: %s", err.Error()) //nolint:stylecheck
	}

	tConfig, err := target.GetTargetConfig(args[0], baseConfig)
	if err != nil {
		return fmt.Errorf("An error occurred while locating the requested target in the configuration file: %s", err.Error()) //nolint:stylecheck
	}

	if tConfig == nil {
		return fmt.Errorf("No target %q found in the configuration file.", args[0]) //nolint:stylecheck
	}

	inputFile, err := os.Open(args[1])
	if err != nil {
		return fmt.Errorf("Unable to read file %q: %s", args[1], err.Error()) //nolint:stylecheck
	}

	if inputFile == nil {
		return fmt.Errorf("Unable to read file %q", args[1]) //nolint:stylecheck
	}

	fullPath, err := filepath.Abs(inputFile.Name())
	if err != nil {
		return fmt.Errorf("Unable to determine full path for file %q: %s", args[1], err.Error()) //nolint:stylecheck
	}

	pterm.Println()
	pterm.Println(pterm.Sprintf("Applying access to target %q from file %q", tConfig.Name, fullPath))

	pterm.Println()
	printer := pterm.Error
	printer.Prefix.Text = "ATTENTION"
	printer.Println(warningMessage)

	pterm.Println()

	// Show an interactive confirmation dialog and get the result.
	confirm := pterm.InteractiveTextInputPrinter{
		DefaultText: "Are you sure you want to continue? " + pterm.ThemeDefault.SecondaryStyle.Sprint("[Yes/No]"),
		Delimiter:   ": ",
		TextStyle:   &pterm.ThemeDefault.DefaultText,
		Mask:        "",
	}

	result, _ := confirm.Show()

	if !strings.EqualFold(result, "yes") {
		pterm.Println("Operation cancelled")
		return nil
	}

	pterm.Println()
	pterm.Println("Starting to apply access from file ...")
	pterm.Println()

	// Starting the actual run
	err = applyAccess(context.Background(), tConfig, fullPath)
	if err != nil {
		return err
	}

	pterm.Println()
	pterm.Success.Println("Access applied successfully")

	return nil
}

func applyAccess(ctx context.Context, tConfig *types.BaseTargetConfig, inputFile string) error {
	client, err := plugin.NewPluginClient(tConfig.ConnectorName, tConfig.ConnectorVersion, tConfig.TargetLogger)
	if err != nil {
		return fmt.Errorf("Error while initializing connector plugin %q: %s", tConfig.ConnectorName, err.Error()) //nolint:stylecheck
	}

	defer client.Close()

	accessSyncer, err := client.GetAccessSyncer()
	if err != nil {
		return fmt.Errorf("Error while fetching the access syncer from connector %q: %s", tConfig.ConnectorName, err.Error()) //nolint:stylecheck
	}

	_, err = version_management.IsValidToSync(ctx, accessSyncer, dapc.MinimalCliVersion)
	if err != nil {
		return fmt.Errorf("Error while checking version compatibility of connector %q: %s", tConfig.ConnectorName, err.Error()) //nolint:stylecheck
	}

	err = tConfig.CalculateFileBackupLocationForRun("apply-access")
	if err != nil {
		return fmt.Errorf("Error while setting up backup location: %s", err.Error()) //nolint:stylecheck
	}

	targetFile, err := filepath.Abs(file.CreateUniqueFileNameForTarget(tConfig.Name, "toTarget-accessFeedback", "json"))
	if err != nil {
		return fmt.Errorf("Error while creating access feedback file: %s", err.Error()) //nolint:stylecheck
	}

	defer tConfig.HandleTempFile(targetFile)

	syncerConfig := dapc.AccessSyncToTarget{
		ConfigMap:          &baseconfig.ConfigMap{Parameters: tConfig.Parameters},
		Prefix:             "",
		SourceFile:         inputFile,
		FeedbackTargetFile: targetFile,
	}

	tConfig.TargetLogger.Info("Synchronizing access providers between Raito and the data source")

	_, err = accessSyncer.SyncToTarget(ctx, &syncerConfig)
	if err != nil {
		return fmt.Errorf("Error while syncing access from file to target %q: %s", tConfig.Name, err.Error()) //nolint:stylecheck
	}

	return nil
}
