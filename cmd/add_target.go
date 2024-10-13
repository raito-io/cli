package cmd

import (
	"context"
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/token"
	"github.com/hashicorp/go-hclog"
	"github.com/pterm/pterm"
	pl "github.com/raito-io/cli/base/util/plugin"
	"github.com/raito-io/cli/internal/auth"
	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/graphql"
	"github.com/raito-io/cli/internal/logging"
	"github.com/raito-io/cli/internal/plugin"
	"github.com/raito-io/cli/internal/target/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/maps"
)

func initAddTargetCommand(rootCmd *cobra.Command) {
	var cmd = &cobra.Command{
		Short:     "Add a new target to the CLI configuration.",
		Long:      "Add a new target to the CLI configuration. If needed, this can also create a new data source in Raito Cloud.",
		Run:       executeAddTargetCmd,
		ValidArgs: []string{},
		Use:       "add-target",
	}

	cmd.PersistentFlags().StringP(constants.NameFlag, "n", "", "The name for the target. If not specified, this will be asked in interactive mode.")
	cmd.PersistentFlags().StringP(constants.ConnectorNameFlag, "c", "", "The name of the connector to use. If not specified, this will be asked in interactive mode.")
	cmd.PersistentFlags().StringP(constants.DataSourceIdFlag, "d", "", "The id of the data source in Raito Cloud. If not specified, this will be asked in interactive mode or you can create a new data source directly in interactive mode.")

	BindFlag(constants.NameFlag, cmd)
	BindFlag(constants.ConnectorNameFlag, cmd)
	BindFlag(constants.DataSourceIdFlag, cmd)

	rootCmd.AddCommand(cmd)
}

func executeAddTargetCmd(cmd *cobra.Command, args []string) {
	logging.SetupLogging(false)

	configFile, configAlreadyExists := getConfigFile()

	pterm.Println()
	pterm.Println("Welcome to the Raito CLI!")

	pterm.Println()

	pterm.Println("This tool will help you to add a new target to your Raito CLI configuration file")
	pterm.Println("The configuration file will be stored at " + pterm.Bold.Sprint(configFile))

	var targetsNode *ast.SequenceNode
	var baseDocument *ast.DocumentNode
	var existingApiUser, existingApiSecret, existingDomain string
	var domainNode, apiUserNode, apiSecretNode *ast.StringNode

	// First we'll search for existing configuration nodes if there are any.
	if configAlreadyExists {
		targetsNode, baseDocument, existingApiUser, existingApiSecret, existingDomain, apiUserNode, apiSecretNode, domainNode = readFromExistingConfigFile(configFile)
	}

	// When not all base configs are created, we need to do an init first
	if domainNode == nil || apiUserNode == nil || apiSecretNode == nil || targetsNode == nil {
		var result bool

		if !configAlreadyExists {
			pterm.Warning.Println("The configuration file does not exist yet.")
			result, _ = pterm.DefaultInteractiveConfirm.Show("Would you like to run the 'init' command first to create the configuration file?")
		} else {
			pterm.Warning.Println("Not all base configuration options are set in the configuration file.")
			result, _ = pterm.DefaultInteractiveConfirm.Show("Would you like to run the 'init' command first to set them?")
		}

		pterm.Println()

		if !result {
			os.Exit(0)
		}

		baseDocument, targetsNode, domainNode, apiUserNode, apiSecretNode = buildBaseConfig(baseDocument, domainNode, apiUserNode, apiSecretNode, targetsNode, existingDomain, existingApiUser, existingApiSecret)
	}

	newTargetNode := buildTargetNode(domainNode.Value, apiUserNode.Value, apiSecretNode.Value)
	if newTargetNode != nil {
		targetsNode.Values = append(targetsNode.Values, newTargetNode)

		storeConfigFile(baseDocument, configFile)
	}
}

func buildTargetNode(domain, apiUser, apiSecret string) *ast.MappingNode {
	targetNode := &ast.MappingNode{
		BaseNode: &ast.BaseNode{},
		Start: &token.Token{
			Position: &token.Position{
				Column: 1,
			},
		},
	}

	// Handling the name of the target
	targetName := viper.GetString(constants.NameFlag)
	if targetName == "" {
		targetName, _ = pterm.DefaultInteractiveTextInput.Show("Enter the name of the target (e.g. 'my-snowflake-target')")
		pterm.Println()
	}

	addStringNodeWithValue("name", targetName, targetNode)

	var pluginInfo *pl.PluginInfo

	connector := viper.GetString(constants.ConnectorNameFlag)

	for {
		if connector == "" {
			connector, _ = pterm.DefaultInteractiveTextInput.Show("Enter the full name of the connector to use (e.g. raito-io/cli-plugin-snowflake")
			pterm.Println()
		}

		pluginInfo = fetchPluginInfo(connector)
		if pluginInfo != nil {
			addStringNodeWithValue("connector-name", connector, targetNode)

			break
		}
	}

	// Fetching the data source and identity store id
	dataSourceId := viper.GetString(constants.DataSourceIdFlag)
	identityStoreId := ""

	if dataSourceId == "" {
		dataSources := fetchDataSources(domain, apiUser, apiSecret)
		if dataSources == nil {
			os.Exit(1)
		}

		if len(dataSources) == 0 {
			fatalError("No data sources found in Raito Cloud. Please create a data source first.")
		}

		dsNames := maps.Keys(dataSources)
		sort.Strings(dsNames)

		selectedOption, _ := pterm.DefaultInteractiveSelect.WithOptions(dsNames).Show()
		selectedDS := dataSources[selectedOption]

		dataSourceId = selectedDS.Id
		identityStoreId = selectedDS.IdentityStore
	} else {
		dsInfo := fetchDataSource(dataSourceId, domain, apiUser, apiSecret)
		if dsInfo == nil {
			os.Exit(1)
		}

		identityStoreId = dsInfo.IdentityStore
	}

	addStringNodeWithValue("data-source-id", dataSourceId, targetNode)
	addStringNodeWithValue("identity-store-id", identityStoreId, targetNode)

	// Now getting the values for all the mandatory parameters for the connector
	optionalParameters := make([]*pl.ParameterInfo, 0, len(pluginInfo.Parameters))

	for _, param := range pluginInfo.Parameters {
		if param.Mandatory {
			value, _ := pterm.DefaultInteractiveTextInput.Show(fmt.Sprintf("Enter the value for parameter %q: %s", param.Name, param.Description))
			pterm.Println()

			addStringNodeWithValue(param.Name, value, targetNode)
		} else {
			optionalParameters = append(optionalParameters, param)
		}
	}

	// Now also handling the option parameters
	if len(optionalParameters) > 0 {
		params := make([]string, 0, len(optionalParameters)+1)

		params = append(params, "No thank you")

		for _, param := range optionalParameters {
			params = append(params, pterm.Bold.Sprint(param.Name)+": "+param.Description)
		}

		for {
			pterm.Println("There are optional parameters for this connector. You can set them now or skip them for now.")
			pterm.Println()

			selectedOption, _ := pterm.DefaultInteractiveSelect.WithOptions(params).Show()
			selectedIndex := slices.Index(params, selectedOption)

			if selectedIndex == 0 {
				break
			}

			param := optionalParameters[selectedIndex-1]
			value, _ := pterm.DefaultInteractiveTextInput.Show(fmt.Sprintf("Enter the value for parameter %q", param.Name))
			pterm.Println()

			addStringNodeWithValue(param.Name, value, targetNode)
		}
	}

	return targetNode
}

func addStringNodeWithValue(key, value string, parent *ast.MappingNode) {
	node := addStringNode(key, parent)
	node.Value = value
}

func fetchPluginInfo(connector string) *pl.PluginInfo {
	spinner, _ := pterm.DefaultSpinner.Start("Looking for connector info...")

	client, err := plugin.NewPluginClient(connector, "", hclog.L())
	if err != nil {
		spinner.Fail(fmt.Sprintf("Error loading plugin %q: %s", connector, err.Error()))
		pterm.Println()

		return nil
	}

	defer client.Close()

	info, err := client.GetInfo()
	if err != nil {
		spinner.Fail(fmt.Sprintf("The plugin (%s) does not implement the Info interface.", connector))
		pterm.Println()

		return nil
	}

	pluginInfo, err := info.GetInfo(context.Background())
	if err != nil {
		spinner.Fail(fmt.Sprintf("Failed to load plugin info: %s", err))
		pterm.Println()

		return nil
	}

	spinner.Success("Loaded plugin info.")
	pterm.Println()

	return pluginInfo
}

func fetchDataSources(domain, apiUser, apiSecret string) map[string]DataSourceInfo {
	spinner, _ := pterm.DefaultSpinner.Start("Loading available data sources from Raito Cloud...")

	tmpConfig := &types.BaseConfig{
		Domain:     domain,
		ApiUser:    apiUser,
		ApiSecret:  apiSecret,
		BaseLogger: hclog.L(),
	}

	gql := `{ "operationName": "dataSources", "variables":{}, "query": "query dataSources {
        dataSources {
          ... on PagedResult {
            edges {
              node {
                ... on DataSource {
                  name
                  id
                  identityStores {
                    id
                    native
                  }
                }
              }
            }
          }
        }
      }"}`

	gql = strings.Replace(gql, "\n", "\\n", -1)

	// Temporarily disable config reload to avoid reloading the config while testing the connection as we need to use this config specifically
	auth.SetNoConfigReload(true)
	defer auth.SetNoConfigReload(false)

	response := DataSourcesResponse{}
	res, err := graphql.ExecuteGraphQL(gql, tmpConfig, &response)

	if err != nil {
		spinner.Fail("An error occurred while fetching data sources: " + err.Error())
		pterm.Println()

		return nil
	}

	if len(res.Errors) > 0 {
		spinner.Fail("An error occurred while fetching data sources: " + res.Errors[0].Message)
		pterm.Println()

		return nil
	}

	spinner.Success("Fetched data sources successfully.")
	pterm.Println()

	return response.GetDataSources()
}

func fetchDataSource(dsId, domain, apiUser, apiSecret string) *DataSourceInfo {
	spinner, _ := pterm.DefaultSpinner.Start("Loading data source from Raito Cloud...")

	tmpConfig := &types.BaseConfig{
		Domain:     domain,
		ApiUser:    apiUser,
		ApiSecret:  apiSecret,
		BaseLogger: hclog.L(),
	}

	gql := fmt.Sprintf(`{ "operationName": "dataSource", "variables":{}, "query": "query dataSource {
        dataSource(id: \"%s\") {
          ... on DataSource {
            name
            id
            identityStores {
              id
              native
            }
          }
        }
      }"}`, dsId)

	gql = strings.Replace(gql, "\n", "\\n", -1)

	// Temporarily disable config reload to avoid reloading the config while testing the connection as we need to use this config specifically
	auth.SetNoConfigReload(true)
	defer auth.SetNoConfigReload(false)

	response := DataSourceResponse{}
	res, err := graphql.ExecuteGraphQL(gql, tmpConfig, &response)

	if err != nil {
		spinner.Fail("An error occurred while fetching data source: " + err.Error())
		pterm.Println()

		return nil
	}

	if len(res.Errors) > 0 {
		spinner.Fail("An error occurred while fetching data source: " + res.Errors[0].Message)
		pterm.Println()

		return nil
	}

	spinner.Success("Fetched data source information.")
	pterm.Println()

	return response.GetDataSourceInfo()
}

type DataSourceInfo struct {
	Name          string
	Id            string
	IdentityStore string
}

type identityStoreNode struct {
	Id     string `json:"id"`
	Native bool   `json:"native"`
}

type dataSourceNode struct {
	Name           string              `json:"name"`
	Id             string              `json:"id"`
	IdentityStores []identityStoreNode `json:"identityStores"`
}

type dataSourceEdge struct {
	Node dataSourceNode `json:"node"`
}

type dataSourcesResponse struct {
	Edges []dataSourceEdge `json:"edges"`
}

type DataSourceResponse struct {
	DataSource dataSourceNode `json:"dataSource"`
}

func (d *DataSourceResponse) GetDataSourceInfo() *DataSourceInfo {
	dsInfo := DataSourceInfo{
		Name: d.DataSource.Name,
		Id:   d.DataSource.Id,
	}

	if len(d.DataSource.IdentityStores) > 0 {
		for _, store := range d.DataSource.IdentityStores {
			if store.Native {
				dsInfo.IdentityStore = store.Id
				break
			}
		}
	}

	return &dsInfo
}

type DataSourcesResponse struct {
	Response dataSourcesResponse `json:"dataSources"`
}

func (d *DataSourcesResponse) GetDataSources() map[string]DataSourceInfo {
	dataSources := make(map[string]DataSourceInfo)

	for _, edge := range d.Response.Edges {
		identityStore := ""

		if len(edge.Node.IdentityStores) > 0 {
			for _, store := range edge.Node.IdentityStores {
				if store.Native {
					identityStore = store.Id
					break
				}
			}
		}

		dataSources[edge.Node.Name] = DataSourceInfo{
			Name:          edge.Node.Name,
			Id:            edge.Node.Id,
			IdentityStore: identityStore,
		}
	}

	return dataSources
}
