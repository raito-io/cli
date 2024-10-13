package cmd

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"github.com/goccy/go-yaml/printer"
	"github.com/goccy/go-yaml/token"
	"github.com/hashicorp/go-hclog"
	"github.com/pterm/pterm"
	"github.com/raito-io/cli/internal/auth"
	"github.com/raito-io/cli/internal/graphql"
	"github.com/raito-io/cli/internal/logging"
	"github.com/raito-io/cli/internal/target/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func initInitCommand(rootCmd *cobra.Command) {
	var cmd = &cobra.Command{
		Short:     "Initialize the CLI configuration.",
		Long:      "Helper tool to initialize the CLI by setting up the configuration file.",
		Run:       executeInitCmd,
		ValidArgs: []string{},
		Use:       "init",
	}

	rootCmd.AddCommand(cmd)
}

func executeInitCmd(cmd *cobra.Command, args []string) {
	logging.SetupLogging(false)

	configFile, configAlreadyExists := getConfigFile()

	pterm.Println()
	pterm.Println("Welcome to the Raito CLI!")

	pterm.Println()

	pterm.Println("This tool will help you set up the configuration file for the Raito CLI.")
	pterm.Println("The configuration file will be stored at " + pterm.Bold.Sprint(configFile))

	var targetsNode *ast.SequenceNode
	var baseDocument *ast.DocumentNode
	var existingApiUser, existingApiSecret, existingDomain string
	var domainNode, apiUserNode, apiSecretNode *ast.StringNode

	// First we'll search for existing configuration nodes if there are any.
	if configAlreadyExists {
		targetsNode, baseDocument, existingApiUser, existingApiSecret, existingDomain, apiUserNode, apiSecretNode, domainNode = readFromExistingConfigFile(configFile)
	}

	if domainNode != nil || apiUserNode != nil || apiSecretNode != nil || targetsNode != nil {
		pterm.Println()
		pterm.Warning.Println("Warning: the configuration file already exists with some of the configuration options set. These will be overwritten.")

		result, _ := pterm.DefaultInteractiveConfirm.Show("Are you sure you want to continue?")

		pterm.Println()

		if !result {
			os.Exit(0)
		}
	}

	baseDocument, _, _, _, _ = buildBaseConfig(baseDocument, domainNode, apiUserNode, apiSecretNode, targetsNode, existingDomain, existingApiUser, existingApiSecret) //nolint:dogsled

	storeConfigFile(baseDocument, configFile)
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

func buildBaseConfig(baseDocument *ast.DocumentNode, domainNode *ast.StringNode, apiUserNode *ast.StringNode, apiSecretNode *ast.StringNode, targetsNode *ast.SequenceNode, existingDomain string, existingApiUser string, existingApiSecret string) (*ast.DocumentNode, *ast.SequenceNode, *ast.StringNode, *ast.StringNode, *ast.StringNode) {
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

	if targetsNode == nil {
		targetsNode = addSequenceNode("targets", baseDocument.Body.(*ast.MappingNode))
	}

	domain := existingDomain
	apiUser := existingApiUser
	apiSecret := existingApiSecret

	for {
		domain = readDomain(domain)

		apiUser, _ = pterm.DefaultInteractiveTextInput.WithDefaultValue(apiUser).Show("Enter the e-mail address of your Raito user (needs the Integrator role)")
		pterm.Println()

		apiSecret, _ = pterm.DefaultInteractiveTextInput.WithDefaultValue(apiSecret).WithMask("*").Show("Enter the password of your Raito user")
		pterm.Println()

		res := testConnection(domain, apiUser, apiSecret)

		if res {
			domainNode.Value = domain
			apiUserNode.Value = apiUser
			apiSecretNode.Value = apiSecret

			break
		}
	}

	return baseDocument, targetsNode, domainNode, apiUserNode, apiSecretNode
}

func testConnection(domain, apiUser, apiSecret string) bool {
	spinner, _ := pterm.DefaultSpinner.Start("Checking connection to Raito Cloud...")

	tmpConfig := &types.BaseConfig{
		Domain:     domain,
		ApiUser:    apiUser,
		ApiSecret:  apiSecret,
		BaseLogger: hclog.L(),
	}

	gql := `{ "operationName": "currentUser", "variables":{}, "query": "query currentUser {
        currentUser {
          email
        }
      }"}`

	gql = strings.Replace(gql, "\n", "\\n", -1)

	// Temporarily disable config reload to avoid reloading the config while testing the connection as we need to use this config specifically
	auth.SetNoConfigReload(true)
	defer auth.SetNoConfigReload(false)

	response := CurrentUserResponse{}
	res, err := graphql.ExecuteGraphQL(gql, tmpConfig, &response)

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
		domain, _ := pterm.DefaultInteractiveTextInput.WithDefaultValue(current).Show("Enter the domain of your Raito instance (https://<mysubdomain>.raito.cloud)")
		pterm.Println()

		if regexp.MustCompile("[a-zA-Z][a-zA-Z0-9-]*[a-zA-Z0-9]").MatchString(domain) {
			return domain
		}

		pterm.Error.Println("Invalid subdomain format. Make sure to only enter the subdomain part.")
	}
}

func addSequenceNode(key string, parent *ast.MappingNode) *ast.SequenceNode {
	node := &ast.SequenceNode{
		BaseNode: &ast.BaseNode{},
		Start: &token.Token{
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
