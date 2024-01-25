package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	idp "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"

	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/target/types"
	"github.com/raito-io/cli/internal/util/url"
)

type userTokens struct {
	userName     string
	idToken      string
	refreshToken string
	expiration   *time.Time
}

var (
	mutex       sync.Mutex
	tokenMap    = make(map[string]*userTokens)
	clientAppId string
)

func AddTokenToHeader(h *http.Header, config *types.BaseConfig) error {
	if viper.GetBool(constants.SkipAuthentication) {
		config.BaseLogger.Debug("Skipping authentication")
		return nil
	}

	tokens, found := tokenMap[config.ApiUser]
	if !found {
		tokens = &userTokens{userName: config.ApiUser}
		tokenMap[config.ApiUser] = tokens
	}

	err := updateTokens(config, tokens)
	if err != nil {
		return err
	}

	h.Add("Authorization", "token "+tokens.idToken)

	return nil
}

func AddToken(r *http.Request, config *types.BaseConfig) error {
	return AddTokenToHeader(&r.Header, config)
}

func updateTokens(config *types.BaseConfig, tokens *userTokens) error {
	if checkTokenValidity(config, tokens) {
		config.BaseLogger.Debug(fmt.Sprintf("Token for user %q is still valid", tokens.userName))
		return nil
	}

	return fetchTokens(config, tokens)
}

func checkTokenValidity(config *types.BaseConfig, tokens *userTokens) bool {
	if tokens.idToken == "" || tokens.refreshToken == "" || tokens.expiration == nil {
		return false
	}

	// Adding a buffer of 10 seconds
	now := time.Now().Add(time.Second * 10)
	if now.After(*tokens.expiration) {
		config.BaseLogger.Debug(fmt.Sprintf("Token for user %q is expired", tokens.userName))
		return false
	}

	return true
}

func fetchTokens(config *types.BaseConfig, tokens *userTokens) error {
	if tokens.refreshToken != "" {
		err := refreshTokens(config, tokens)
		if err != nil {
			config.BaseLogger.Warn(fmt.Sprintf("error while trying to refresh tokens: %s. Trying to fetch tokens from scratch instead", err.Error()))
		} else {
			return nil
		}
	}

	// If no refresh token or refreshing failed
	return fetchNewTokens(config, tokens)
}

func refreshTokens(baseConfig *types.BaseConfig, tokens *userTokens) error {
	baseConfig.BaseLogger.Debug(fmt.Sprintf("Refreshing tokens for user %q", tokens.userName))

	err := fetchClientAppId(baseConfig)
	if err != nil {
		return fmt.Errorf("error while fetching clientAppId: %s", err.Error())
	}

	if isInTest() {
		return setTestTokens(tokens)
	}

	// TODO configurable region
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("eu-central-1"))
	if err != nil {
		return fmt.Errorf("error while configuring AWS SDK: %s", err.Error())
	}
	idpClient := idp.NewFromConfig(cfg)
	output, err := idpClient.InitiateAuth(context.TODO(), &idp.InitiateAuthInput{
		AuthFlow:       "REFRESH_TOKEN_AUTH",
		ClientId:       &clientAppId,
		AuthParameters: map[string]string{"REFRESH_TOKEN": tokens.refreshToken},
	})

	if err != nil {
		return fmt.Errorf("error while refreshing tokens %q: %s", tokens.userName, err.Error())
	}

	return handleAuthOutput(output, tokens)
}

func setTestTokens(tokens *userTokens) error {
	tokens.idToken = "idToken"
	tokens.refreshToken = "refreshToken"
	exp := time.Now().Add(time.Hour)
	tokens.expiration = &exp

	return nil
}

func fetchNewTokens(baseConfig *types.BaseConfig, tokens *userTokens) error {
	baseConfig.BaseLogger.Debug("Fetching new tokens")

	err := fetchClientAppId(baseConfig)
	if err != nil {
		return fmt.Errorf("error while fetching clientAppId: %s", err.Error())
	}

	if isInTest() {
		return setTestTokens(tokens)
	}

	// TODO configurable region
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("eu-central-1"))
	if err != nil {
		return fmt.Errorf("error while configuring AWS SDK: %s", err.Error())
	}
	idpClient := idp.NewFromConfig(cfg)
	output, err := idpClient.InitiateAuth(context.TODO(), &idp.InitiateAuthInput{
		AuthFlow:       "USER_PASSWORD_AUTH",
		ClientId:       &clientAppId,
		AuthParameters: map[string]string{"USERNAME": baseConfig.ApiUser, "PASSWORD": baseConfig.ApiSecret},
	})

	if err != nil {
		return fmt.Errorf("error while initiating authentication flow for user %q: %s", tokens.userName, err.Error())
	}

	return handleAuthOutput(output, tokens)
}

func handleAuthOutput(output *idp.InitiateAuthOutput, tokens *userTokens) error {
	if output.AuthenticationResult != nil {
		if output.AuthenticationResult.IdToken == nil {
			return fmt.Errorf("no id token found in authentication result")
		}

		if output.AuthenticationResult.RefreshToken != nil {
			tokens.refreshToken = *output.AuthenticationResult.RefreshToken
		}

		tokens.idToken = *output.AuthenticationResult.IdToken
		e := time.Now().Add(time.Second * time.Duration(output.AuthenticationResult.ExpiresIn))
		tokens.expiration = &e

		return nil
	} else {
		return fmt.Errorf("invalid authentication result received (challenge %q)", output.ChallengeName)
	}
}

func fetchClientAppId(baseConfig *types.BaseConfig) error {
	mutex.Lock()
	defer mutex.Unlock()

	if clientAppId == "" {
		domain := baseConfig.Domain

		if domain == "" {
			return fmt.Errorf("no domain specified")
		}

		domain = strings.ToLower(domain)
		if !isValidDomain(domain) {
			return fmt.Errorf("invalid domain name %q. A domain should start with a letter and can only contain alphanumeric characters and the dash character. It also should not end with a dash character", domain)
		}

		if isInTest() {
			clientAppId = "testclient"
			return nil
		}

		url := url.CreateRaitoURL(url.GetRaitoURL(), "admin/org/"+domain)

		req, err := http.NewRequest("GET", url, http.NoBody)
		if err != nil {
			return fmt.Errorf("error while creating HTTP GET request to %q: %s", url, err.Error())
		}
		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("error while doing HTTP GET to %q: %s", url, err.Error())
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return fmt.Errorf("unexpected status code %q received when calling URL %q", resp.StatusCode, url)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error while reading body for call to %q: %s", url, err.Error())
		}

		org := orgInfo{}

		err = json.Unmarshal(body, &org)
		if err != nil {
			return fmt.Errorf("error while parsing organization info response from %q: %s", url, err.Error())
		}

		clientAppId = org.ClientAppId
		baseConfig.BaseLogger.Debug(fmt.Sprintf("Received clientAppId %q for domain %q", clientAppId, domain))
	}

	return nil
}

func isValidDomain(domain string) bool {
	matched, err := regexp.Match("^[a-z][a-z0-9-]*[a-z0-9]$", []byte(domain))
	if err != nil {
		hclog.L().Error(fmt.Sprintf("Error while checking domain validity: %s", err.Error()))
		return false
	}

	return matched
}

type orgInfo struct {
	AuthOrgId   string
	ClientAppId string
}

func isInTest() bool {
	return strings.HasSuffix(os.Args[0], ".test")
}
