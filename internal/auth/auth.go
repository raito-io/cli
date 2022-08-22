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
	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/target"
	"github.com/raito-io/cli/internal/util/url"
	"github.com/spf13/viper"
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

func AddToken(r *http.Request, targetConfig *target.BaseTargetConfig) error {
	env := viper.GetString(constants.EnvironmentFlag)
	if env == constants.EnvironmentDev {
		targetConfig.Logger.Debug("Skipping authentication for development environment.")
		return nil
	}

	tokens, found := tokenMap[targetConfig.ApiUser]
	if !found {
		tokens = &userTokens{userName: targetConfig.ApiUser}
		tokenMap[targetConfig.ApiUser] = tokens
	}

	err := updateTokens(targetConfig, tokens)
	if err != nil {
		return err
	}

	r.Header.Add("Authorization", "token "+tokens.idToken)

	return nil
}

func updateTokens(targetConfig *target.BaseTargetConfig, tokens *userTokens) error {
	if checkTokenValidity(targetConfig, tokens) {
		targetConfig.Logger.Debug(fmt.Sprintf("Token for user %q is still valid", tokens.userName))
		return nil
	}

	return fetchTokens(targetConfig, tokens)
}

func checkTokenValidity(targetConfig *target.BaseTargetConfig, tokens *userTokens) bool {
	if tokens.idToken == "" || tokens.refreshToken == "" || tokens.expiration == nil {
		return false
	}

	// Adding a buffer of 10 seconds
	now := time.Now().Add(time.Second * 10)
	if now.After(*tokens.expiration) {
		targetConfig.Logger.Debug(fmt.Sprintf("Token for user %q is expired", tokens.userName))
		return false
	}

	return true
}

func fetchTokens(targetConfig *target.BaseTargetConfig, tokens *userTokens) error {
	if tokens.refreshToken != "" {
		err := refreshTokens(targetConfig, tokens)
		if err != nil {
			targetConfig.Logger.Error(fmt.Sprintf("error while trying to refresh tokens: %s. Trying to fetch tokens from scratch", err.Error()))
		} else {
			return nil
		}
	}

	// If no refresh token or refreshing failed
	return fetchNewTokens(targetConfig, tokens)
}

func refreshTokens(targetConfig *target.BaseTargetConfig, tokens *userTokens) error {
	targetConfig.Logger.Debug(fmt.Sprintf("Refreshing tokens for user %q", tokens.userName))

	err := fetchClientAppId(targetConfig)
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
		return fmt.Errorf("error while initiating authentication flow for user %q: %s", tokens.userName, err.Error())
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

func fetchNewTokens(targetConfig *target.BaseTargetConfig, tokens *userTokens) error {
	targetConfig.Logger.Debug("Fetching new tokens")

	err := fetchClientAppId(targetConfig)
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
		AuthParameters: map[string]string{"USERNAME": targetConfig.ApiUser, "PASSWORD": targetConfig.ApiSecret},
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

		if output.AuthenticationResult.RefreshToken == nil {
			return fmt.Errorf("no refresh token found in authentication result")
		}

		tokens.idToken = *output.AuthenticationResult.IdToken
		tokens.refreshToken = *output.AuthenticationResult.RefreshToken
		e := time.Now().Add(time.Second * time.Duration(output.AuthenticationResult.ExpiresIn))
		tokens.expiration = &e

		return nil
	} else {
		return fmt.Errorf("invalid authentication result received (challenge %q)", output.ChallengeName)
	}
}

func fetchClientAppId(targetConfig *target.BaseTargetConfig) error {
	mutex.Lock()
	defer mutex.Unlock()

	if clientAppId == "" {
		domain := targetConfig.Domain

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
		targetConfig.Logger.Info(fmt.Sprintf("Received clientAppId %q for domain %q", clientAppId, domain))
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
