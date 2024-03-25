package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/pterm/pterm"
	dapc "github.com/raito-io/cli/base/access_provider"
	baseconfig "github.com/raito-io/cli/base/util/config"
	"github.com/raito-io/cli/base/util/match"
	"github.com/raito-io/cli/base/util/slice"
	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/health_check"
	"github.com/raito-io/cli/internal/logging"
	"github.com/raito-io/cli/internal/plugin"
	"github.com/raito-io/cli/internal/target"
	"github.com/raito-io/cli/internal/target/types"
	"github.com/raito-io/cli/internal/util/file"
	"github.com/raito-io/cli/internal/version_management"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func initApplyAccessCommand(rootCmd *cobra.Command) {
	var cmd = &cobra.Command{
		Short:     "(Re)apply access to the given target from a specific file.",
		Long:      "This can be used to (re)apply a given access file to a specific target. This is a very dangerous operation and should be used with caution.",
		Run:       executeApplyAccessCmd,
		ValidArgs: []string{},
		Use:       "apply-access <target-name> <file-path>",
	}

	cmd.PersistentFlags().String(constants.FilterAccessFlag, "", "To only match a subset of access providers in the file, provide a comma-separated list of access provider names to match. These names can also be specified as regular expressions.")
	BindFlag(constants.FilterAccessFlag, cmd)

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

	filters := viper.GetString(constants.FilterAccessFlag)
	filteredFile, filteredAps, totalAps, err := filterAccessProviders(tConfig.Name, fullPath, filters)

	if err != nil {
		return fmt.Errorf("Error while filtering access providers: %s", err.Error()) //nolint:stylecheck
	}

	isFiltered := filteredFile != fullPath
	fullPath = filteredFile

	if filteredAps == nil {
		pterm.Println(fmt.Sprintf("All %d access providers in the file will be applied. If you would like to filter the access providers, please provide a filter using the --filter-access flag.", totalAps))
	} else if len(filteredAps) == 0 {
		pterm.Println("No access providers in the file match the filter. Nothing to apply.")
		return nil
	} else {
		pterm.Println(fmt.Sprintf("%d of the total %d access providers matches the provided filter. Only these access provides will be applied: %s", len(filteredAps), totalAps, strings.Join(filteredAps, ", ")))
	}

	if !requestConfirmation() {
		return nil
	}

	pterm.Println()
	pterm.Println("Starting to apply access from file ...")
	pterm.Println()

	// Starting the actual run
	err = applyAccess(context.Background(), tConfig, fullPath, isFiltered)
	if err != nil {
		return err
	}

	pterm.Println()
	pterm.Success.Println("Access applied successfully")

	return nil
}

func requestConfirmation() bool {
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
		return false
	}

	return true
}

func applyAccess(ctx context.Context, tConfig *types.BaseTargetConfig, inputFile string, isFiltered bool) error {
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

	defer func() {
		tConfig.HandleTempFile(targetFile, false)

		// If the file is the filtered one, we can remove it after the run like normal, but if it's the original one, we need to keep it.
		tConfig.HandleTempFile(inputFile, !isFiltered)

		tConfig.FinalizeRun()
	}()

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

func filterAccessProviders(targetName, source, filter string) (string, []string, int, error) {
	fileContent, err := readFileStructure(source)
	if err != nil {
		return "", nil, 0, err
	}

	filter = strings.TrimSpace(filter)
	if filter == "" {
		return source, nil, len(fileContent.AccessProviders), nil
	}

	newFileContent, filteredAps, err := filterFileStructure(*fileContent, filter)
	if err != nil {
		return "", nil, 0, err
	}

	filteredFileName, err := filepath.Abs(file.CreateUniqueFileNameForTarget(targetName, "toTarget-filteredAccess", "yaml"))
	if err != nil {
		return "", nil, 0, fmt.Errorf("unable to create temporary file with filtered access: %s", err.Error())
	}

	err = writeFileStructure(filteredFileName, newFileContent)
	if err != nil {
		return "", nil, 0, err
	}

	return filteredFileName, filteredAps, len(fileContent.AccessProviders), nil
}

func writeFileStructure(output string, content FileStructure) error {
	// Writing out the file to a temp file
	filtered, err := yaml.Marshal(content)
	if err != nil {
		return fmt.Errorf("unable to marshal filtered content: %s", err.Error())
	}

	filteredFile, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("unable to create temporary file with filtered access %q: %s", output, err.Error())
	}

	defer filteredFile.Close()

	_, err = filteredFile.Write(filtered)
	if err != nil {
		return fmt.Errorf("unable writing to temporary file with filtered access %q: %s", output, err.Error())
	}

	return nil
}

func readFileStructure(input string) (*FileStructure, error) {
	af, err := os.Open(input)
	if err != nil {
		return nil, fmt.Errorf("unable to open source file %q: %s", input, err.Error())
	}

	buf, err := io.ReadAll(af)
	if err != nil {
		return nil, fmt.Errorf("unable to read source file %q: %s", input, err.Error())
	}

	var fileContent FileStructure

	if json.Valid(buf) {
		err = json.Unmarshal(buf, &fileContent)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal source file %q as JSON: %s", input, err.Error())
		}
	} else {
		err = yaml.Unmarshal(buf, &fileContent)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal source file %q as YAML: %s", input, err.Error())
		}
	}

	return &fileContent, nil
}

func filterFileStructure(input FileStructure, filter string) (FileStructure, []string, error) {
	// Creating the contents for the new file
	newFileContent := FileStructure{
		LastCalculated:  input.LastCalculated,
		AccessProviders: make([]map[string]interface{}, 0, len(input.AccessProviders)),
	}

	// Now filtering the input file and writing the filtered content to a new file structure
	filters := slice.ParseCommaSeparatedList(filter).Slice()

	filteredAps := make([]string, 0, len(input.AccessProviders))

	for _, ap := range input.AccessProviders {
		if name, f := ap["name"]; f {
			if ns, ok := name.(string); ok {
				isMatch, err2 := match.MatchesAny(ns, filters)

				if err2 != nil {
					return newFileContent, filteredAps, fmt.Errorf("invalid filter format: %s", err2.Error())
				}

				if isMatch {
					filteredAps = append(filteredAps, ns)
					newFileContent.AccessProviders = append(newFileContent.AccessProviders, ap)
				}
			}
		}
	}

	return newFileContent, filteredAps, nil
}

type FileStructure struct {
	LastCalculated  int64                    `yaml:"lastCalculated" json:"lastCalculated"`
	AccessProviders []map[string]interface{} `yaml:"accessProviders" json:"accessProviders"`
}
