package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"github.com/goccy/go-yaml/printer"
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
	"github.com/raito-io/cli/internal/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.design/x/clipboard"
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
	output := &logging.NilWriter{}
	logger := hclog.NewInterceptLogger(&hclog.LoggerOptions{
		Name:   fmt.Sprintf("raito-cli-%s", version.GetCliVersion().String()),
		Output: output,
	})
	hclog.SetDefault(logger)

	configFile, configAlreadyExists := getConfigFile()

	pterm.Println()
	pterm.Println("#########################")
	pterm.Println("Welcome to the Raito CLI!")
	pterm.Println("#########################")
	pterm.Println()

	pterm.Println("This tool will help you to connect the Raito CLI to your Raito Cloud instance and to add a new target to your Raito CLI configuration file.")
	pterm.Println("The configuration file will be stored at " + pterm.Bold.Sprint(configFile))
	pterm.Println()

	var targetsNode *ast.SequenceNode
	var baseDocument *ast.DocumentNode
	var existingApiUser, existingApiSecret, existingDomain string
	var domainNode, apiUserNode, apiSecretNode *ast.StringNode

	// First we'll search for existing configuration nodes if there are any.
	if configAlreadyExists {
		targetsNode, baseDocument, existingApiUser, existingApiSecret, existingDomain, apiUserNode, apiSecretNode, domainNode = readFromExistingConfigFile(configFile)
	}

	// When not all base configs are created, we need to do an init first
	if domainNode == nil || apiUserNode == nil || apiSecretNode == nil {
		if !configAlreadyExists {
			pterm.Info.Println("The configuration file does not exist yet. We'll need to gather some information first to connect to Raito Cloud.")
			pterm.Println()
		} else {
			pterm.Info.Println("Not all base configuration options are set in the configuration file. We'll need to gather them first to connect to Raito Cloud.")
			pterm.Println()
		}

		baseDocument, domainNode, apiUserNode, apiSecretNode = buildBaseConfig(baseDocument, domainNode, apiUserNode, apiSecretNode, existingDomain, existingApiUser, existingApiSecret)
	}

	if targetsNode == nil {
		targetsNode = addSequenceNode("targets", baseDocument.Body.(*ast.MappingNode))
	}

	newTargetNode, newTargetName := buildTargetNode(domainNode.Value, apiUserNode.Value, apiSecretNode.Value)
	if newTargetNode != nil {
		targetsNode.Values = append(targetsNode.Values, newTargetNode)

		storeConfigFile(baseDocument, configFile)

		runCommand := fmt.Sprintf("raito run --only-targets %q", newTargetName)

		if viper.GetString(constants.ConfigFileFlag) != "" {
			runCommand += fmt.Sprintf(" --config-file %q", viper.GetString(constants.ConfigFileFlag))
		}

		pterm.Println(pterm.Bold.Sprint("Congratulations!") + " You have successfully added a new target to your Raito CLI configuration.")
		pterm.Println(fmt.Sprintf("Next, you can run the following command to start a synchronization with this target: %s", pterm.Bold.Sprint(runCommand)))
		pterm.Println()
		pterm.Println("For your convenience, this command has also been copied to the clipboard.")

		clipboard.Write(clipboard.FmtText, []byte(runCommand))
	}
}

func buildTargetNode(domain, apiUser, apiSecret string) (*ast.MappingNode, string) {
	targetNode := &ast.MappingNode{
		BaseNode: &ast.BaseNode{},
		Start: &token.Token{
			Type:      token.MappingValueType,
			Indicator: token.BlockStructureIndicator,
			Position: &token.Position{
				Column:      3,
				IndentNum:   2,
				IndentLevel: 1,
			},
		},
	}

	// Handling the name of the target
	targetName := viper.GetString(constants.NameFlag)
	if targetName == "" {
		targetName = textInput("Enter the name of the new target to add (e.g. 'my-snowflake-target')", "Target name", "", false)
	}

	addStringNodeWithValue("name", targetName, targetNode)

	var pluginInfo *pl.PluginInfo

	connector := viper.GetString(constants.ConnectorNameFlag)

	for {
		if connector == "" {
			connector = textInput("Enter the full name of the connector (e.g. raito-io/cli-plugin-snowflake)", "Connector name", "", false)
		}

		pluginInfo = fetchPluginInfo(connector)
		if pluginInfo != nil {
			addStringNodeWithValue("connector-name", connector, targetNode)

			break
		} else {
			connector = ""
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

		selectedOptionIndex := 0
		selectedOption := ""

		// If there are data sources already, let the user choose one (or create a new one)
		if len(dataSources) > 0 {
			dsNames := maps.Keys(dataSources)
			sort.Strings(dsNames)

			dsNames = append([]string{"Create a new data source"}, dsNames...)

			selectedOption, _ = pterm.DefaultInteractiveSelect.WithOptions(dsNames).Show()
			selectedOptionIndex = slices.Index(dsNames, selectedOption)

			pterm.Println()
		}

		if selectedOptionIndex == 0 {
			dsName := textInput("Enter the name for the new data source (e.g. 'Snowflake Production')", "Name", "", false) // Chosen to create a new data source

			dsInfo := createDataSource(dsName, domain, apiUser, apiSecret)
			if dsInfo == nil {
				os.Exit(1)
			}

			dataSourceId = dsInfo.Id
			identityStoreId = dsInfo.IdentityStore
		} else {
			selectedDS := dataSources[selectedOption]

			dataSourceId = selectedDS.Id
			identityStoreId = selectedDS.IdentityStore
		}
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
			mask := strings.Contains(param.Name, "passw") || strings.Contains(param.Name, "secret")
			value := textInput("Provide the value for the mandatory parameter "+pterm.Bold.Sprint(param.Name)+": "+param.Description, "Value for "+param.Name, "", mask)

			addStringNodeWithValue(param.Name, value, targetNode)
		} else {
			optionalParameters = append(optionalParameters, param)
		}
	}

	// Now also handling the option parameters
	if len(optionalParameters) > 0 {
		params := make([]string, 0, len(optionalParameters)+1)

		params = append(params, pterm.Bold.Sprint("No thank you"))

		for _, param := range optionalParameters {
			paramText := pterm.Bold.Sprint(param.Name) + ": " + param.Description
			if len(paramText) > pterm.GetTerminalWidth()-3 {
				paramText = paramText[:pterm.GetTerminalWidth()-3] + "..."
			}

			params = append(params, paramText)
		}

		for {
			pterm.Println("There are optional parameters for this connector. You can set them now or skip them for now.")
			pterm.Println()

			selectedOption, _ := pterm.DefaultInteractiveSelect.WithOptions(params).Show()
			selectedIndex := slices.Index(params, selectedOption)

			pterm.Println()

			if selectedIndex == 0 {
				break
			}

			param := optionalParameters[selectedIndex-1]

			mask := strings.Contains(param.Name, "passw") || strings.Contains(param.Name, "secret")
			value := textInput("Provide the value for the optional parameter "+pterm.Bold.Sprint(param.Name)+": "+param.Description, "Value for "+param.Name, "", mask)

			addStringNodeWithValue(param.Name, value, targetNode)
		}
	}

	return targetNode, targetName
}

func addStringNodeWithValue(key, value string, parent *ast.MappingNode) {
	node := addStringNode(key, parent)
	node.Value = value
}

func fetchPluginInfo(connector string) *pl.PluginInfo {
	spinner, _ := pterm.DefaultSpinner.Start("Looking for connector information...")

	client, err := plugin.NewPluginClient(connector, "", hclog.L())
	if err != nil {
		spinner.Fail(fmt.Sprintf("Unable to load connector %q: %s", connector, err.Error()))
		pterm.Println()

		return nil
	}

	defer client.Close()

	info, err := client.GetInfo()
	if err != nil {
		spinner.Fail(fmt.Sprintf("Connector (%s) does not implement the Info interface", connector))
		pterm.Println()

		return nil
	}

	pluginInfo, err := info.GetInfo(context.Background())
	if err != nil {
		spinner.Fail(fmt.Sprintf("Unable to load plugin information: %s", err))
		pterm.Println()

		return nil
	}

	spinner.Success(fmt.Sprintf("Loaded connector information for %s", pluginInfo.Name))
	pterm.Println()

	return pluginInfo
}

func fetchDataSources(domain, apiUser, apiSecret string) map[string]DataSourceInfo {
	spinner, _ := pterm.DefaultSpinner.Start("Loading available data sources from Raito Cloud...")

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

	response := DataSourcesResponse{}
	res, err := executeGraphQLWithCustomConfig(gql, domain, apiUser, apiSecret, &response)

	if err != nil {
		spinner.Fail("Unable to fetch data sources: " + err.Error())
		pterm.Println()

		return nil
	}

	if len(res.Errors) > 0 {
		spinner.Fail("Unable to fetch data sources: " + res.Errors[0].Message)
		pterm.Println()

		return nil
	}

	spinner.Success("Data sources fetched successfully")
	pterm.Println()

	return response.GetDataSources()
}

func createDataSource(dsName, domain, apiUser, apiSecret string) *DataSourceInfo {
	spinner, _ := pterm.DefaultSpinner.Start("Creating new data source in Raito Cloud...")

	gql := fmt.Sprintf(`{ "operationName": "dataSource", "variables":{}, "query": "mutation dataSource {
        createDataSource(input: { name: \"%s\" }) {
          ... on DataSource {
            name
            id
            identityStores {
              id
              native
            }
          }
        }
      }"}`, dsName)

	response := CreateDataSourceResponse{}
	res, err := executeGraphQLWithCustomConfig(gql, domain, apiUser, apiSecret, &response)

	if err != nil {
		spinner.Fail("Unable to create data source: " + err.Error())
		pterm.Println()

		return nil
	}

	if len(res.Errors) > 0 {
		spinner.Fail("Unable to create data source: " + res.Errors[0].Message)
		pterm.Println()

		return nil
	}

	spinner.Success("Data source created successfully")
	pterm.Println()

	return response.GetDataSourceInfo()
}

func fetchDataSource(dsId, domain, apiUser, apiSecret string) *DataSourceInfo {
	spinner, _ := pterm.DefaultSpinner.Start("Loading data source from Raito Cloud...")

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

	response := DataSourceResponse{}
	res, err := executeGraphQLWithCustomConfig(gql, domain, apiUser, apiSecret, &response)

	if err != nil {
		spinner.Fail("Unable to fetch data source: " + err.Error())
		pterm.Println()

		return nil
	}

	if len(res.Errors) > 0 {
		spinner.Fail("Unable to fetch data source: " + res.Errors[0].Message)
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

type CreateDataSourceResponse struct {
	DataSource dataSourceNode `json:"createDataSource"`
}

func (d dataSourceNode) GetDataSourceInfo() *DataSourceInfo {
	dsInfo := DataSourceInfo{
		Name: d.Name,
		Id:   d.Id,
	}

	if len(d.IdentityStores) > 0 {
		for _, store := range d.IdentityStores {
			if store.Native {
				dsInfo.IdentityStore = store.Id
				break
			}
		}
	}

	return &dsInfo
}

func (d *DataSourceResponse) GetDataSourceInfo() *DataSourceInfo {
	return d.DataSource.GetDataSourceInfo()
}

func (d *CreateDataSourceResponse) GetDataSourceInfo() *DataSourceInfo {
	return d.DataSource.GetDataSourceInfo()
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

func executeGraphQLWithCustomConfig(gql, domain, apiUser, apiSecret string, response interface{}) (*graphql.GraphqlResponse, error) {
	tmpConfig := &types.BaseConfig{
		Domain:     domain,
		ApiUser:    apiUser,
		ApiSecret:  apiSecret,
		BaseLogger: hclog.L(),
	}

	gql = strings.Replace(gql, "\n", "\\n", -1)

	// Temporarily disable config reload to avoid reloading the config while testing the connection as we need to use this config specifically
	auth.SetNoConfigReload(true)
	defer auth.SetNoConfigReload(false)

	return graphql.ExecuteGraphQL(gql, tmpConfig, &response)
}

func storeConfigFile(baseDocument *ast.DocumentNode, configFile string) {
	// Marshal the AST back to YAML
	var p printer.Printer
	outputData := p.PrintNode(baseDocument.Body)

	// Write the updated YAML back to the file
	err := os.WriteFile(configFile, outputData, 0600)
	if err != nil {
		fatalError(fmt.Sprintf("Error writing updated YAML to file: %s", err.Error()))
	}

	pterm.Println("Successfully updated the configuration file.")

	pterm.Println()
}

func buildBaseConfig(baseDocument *ast.DocumentNode, domainNode *ast.StringNode, apiUserNode *ast.StringNode, apiSecretNode *ast.StringNode, existingDomain string, existingApiUser string, existingApiSecret string) (*ast.DocumentNode, *ast.StringNode, *ast.StringNode, *ast.StringNode) {
	if baseDocument == nil {
		baseDocument = &ast.DocumentNode{
			BaseNode: &ast.BaseNode{},
			Body: &ast.MappingNode{
				BaseNode: &ast.BaseNode{},
			},
		}
	}

	if domainNode == nil {
		domainNode = addStringNode("domain", baseDocument.Body.(*ast.MappingNode))
	}

	if apiUserNode == nil {
		apiUserNode = addStringNode("api-user", baseDocument.Body.(*ast.MappingNode))
	}

	if apiSecretNode == nil {
		apiSecretNode = addStringNode("api-secret", baseDocument.Body.(*ast.MappingNode))
	}

	domain := existingDomain
	apiUser := existingApiUser
	apiSecret := existingApiSecret

	for {
		domain = readDomain(domain)

		apiUser = textInput("Enter the Raito Cloud user email to use to connect to Raito Cloud. This user needs to have the 'Integrator' role", "Raito user e-mail", apiUser, false)

		apiSecret = textInput("Enter the Raito Cloud user password to use to connect to Raito Cloud", "Raito user password", apiSecret, true)

		res := testConnection(domain, apiUser, apiSecret)

		if res {
			domainNode.Value = domain
			apiUserNode.Value = apiUser
			apiSecretNode.Value = apiSecret

			break
		}
	}

	return baseDocument, domainNode, apiUserNode, apiSecretNode
}

func testConnection(domain, apiUser, apiSecret string) bool {
	spinner, _ := pterm.DefaultSpinner.Start("Checking connection to Raito Cloud...")

	gql := `{ "operationName": "currentUser", "variables":{}, "query": "query currentUser {
        currentUser {
          email
        }
      }"}`

	response := CurrentUserResponse{}
	res, err := executeGraphQLWithCustomConfig(gql, domain, apiUser, apiSecret, &response)

	if err != nil {
		spinner.Fail("An error occurred while checking connection: " + err.Error())
		pterm.Println()

		return false
	}

	if len(res.Errors) > 0 {
		spinner.Fail("An error occurred checking connection: " + res.Errors[0].Message)
		pterm.Println()

		return false
	}

	spinner.Success("Your connection details have been verified successfully!")
	pterm.Println()

	return true
}

type currentUserResponse struct {
	Email string `json:"email"`
}

type CurrentUserQueryResponse struct {
	CurrentUser currentUserResponse `json:"currentUser"`
}

type CurrentUserResponse struct {
	Response CurrentUserQueryResponse `json:"currentUser"`
}

func readDomain(current string) string {
	for {
		domain := textInput("Enter the Raito Cloud sub-domain to connect to (https://<sub-domain>.raito.cloud)", "Raito sub-domain", current, false)

		if regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9-]*[a-zA-Z0-9]$").MatchString(domain) {
			return domain
		}

		pterm.Error.Println("Invalid subdomain format. Make sure to only enter the sub-domain part.")
	}
}

func textInput(title, shortTitle, defaultValue string, masked bool) string {
	pterm.Println(title)

	tip := pterm.DefaultInteractiveTextInput.WithDefaultValue(defaultValue)

	if masked {
		tip = tip.WithMask("*")
	}

	res, _ := tip.Show(shortTitle)

	pterm.Println()

	return res
}

func addSequenceNode(key string, parent *ast.MappingNode) *ast.SequenceNode {
	node := &ast.SequenceNode{
		BaseNode: &ast.BaseNode{},
		Start: &token.Token{
			Type:      token.SequenceEntryType,
			Indicator: token.BlockStructureIndicator,
			Value:     "-",
			Origin:    "\n  -",
			Position: &token.Position{
				Column: 1,
			},
		},
	}

	parent.Values = append(parent.Values, &ast.MappingValueNode{
		BaseNode: &ast.BaseNode{},
		Value:    node,
		Key: &ast.StringNode{
			BaseNode: &ast.BaseNode{},
			Token: &token.Token{
				Position: &token.Position{
					Column: 1,
				},
			},
			Value: key,
		},
	})

	return node
}

func addStringNode(key string, parent *ast.MappingNode) *ast.StringNode {
	node := &ast.StringNode{
		BaseNode: &ast.BaseNode{},
		Token: &token.Token{
			Position: &token.Position{
				Column: 1,
			},
		},
	}

	parent.Values = append(parent.Values, &ast.MappingValueNode{
		BaseNode: &ast.BaseNode{},
		Value:    node,
		Key: &ast.StringNode{
			BaseNode: &ast.BaseNode{},
			Token: &token.Token{
				Position: &token.Position{
					Column: 1,
				},
			},
			Value: key,
		},
	})

	return node
}

func checkFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !errors.Is(err, os.ErrNotExist)
}

func getConfigFile() (string, bool) {
	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		configFile = "./raito.yml"
	}

	configAlreadyExists := checkFileExists(configFile)

	return configFile, configAlreadyExists
}

func readFromExistingConfigFile(configFile string) (targetsNode *ast.SequenceNode, baseDocument *ast.DocumentNode, existingApiUser, existingApiSecret, existingDomain string, apiUserNode, apiSecretNode, domainNode *ast.StringNode) { //nolint:gocritic
	astFile, err := parser.ParseFile(configFile, parser.ParseComments)
	if err != nil {
		fatalError("Error parsing existing configuration file: " + err.Error())
	}

	if len(astFile.Docs) == 1 {
		baseDocument = astFile.Docs[0]
		mappingNode, ok := baseDocument.Body.(*ast.MappingNode)

		if !ok {
			fatalError("Unable to parse existing YAML file")
		}

		for _, baseValue := range mappingNode.Values {
			if stringNode, ok := baseValue.Key.(*ast.StringNode); ok {
				switch strings.ToLower(stringNode.Value) {
				case "api-user":
					apiUserNode = baseValue.Value.(*ast.StringNode)
					existingApiUser = apiUserNode.Value
				case "api-secret":
					apiSecretNode = baseValue.Value.(*ast.StringNode)
					existingApiSecret = apiSecretNode.Value
				case "domain":
					domainNode = baseValue.Value.(*ast.StringNode)
					existingDomain = domainNode.Value
				case "targets":
					targetsNode = baseValue.Value.(*ast.SequenceNode)
				}
			}
		}
	}

	return targetsNode, baseDocument, existingApiUser, existingApiSecret, existingDomain, apiUserNode, apiSecretNode, domainNode
}

func fatalError(err string) {
	pterm.Error.Println(err)
	os.Exit(1)
}
