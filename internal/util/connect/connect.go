package connect

import (
	"fmt"
	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/target"
	"github.com/raito-io/cli/internal/util/url"
	"net/http"
	"strings"
)

func AddRaitoAuthorizationHeaders(req *http.Request, config *target.BaseTargetConfig) {
	req.Header.Set(constants.OrgDomainHeader, config.Domain)
	req.Header.Set(constants.ApiUserHeader, config.ApiUser)
	req.Header.Set(constants.ApiSecretHeader, config.ApiSecret)
}

func doPost(host, path, body, contentType string, config *target.BaseTargetConfig) (*http.Response, error) {
	url := url.CreateRaitoURL(host, path)
	config.Logger.Debug("Calling HTTP POST", "URL", url)
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error while creating HTTP GET request to %q: %s", url, err.Error())
	}
	req.Header.Set("Content-Type", contentType)
	AddRaitoAuthorizationHeaders(req, config);
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while doing HTTP GET to %q: %s", url, err.Error())
	}
	return resp, nil
}

func DoPostToRaito(path, body, contentType string, config *target.BaseTargetConfig) (*http.Response, error) {
	return doPost(url.GetRaitoURL(), path, body, contentType, config)
}

func doGet(host, path string, config *target.BaseTargetConfig) (*http.Response, error) {
	url := url.CreateRaitoURL(host, path)
	config.Logger.Debug("Calling HTTP GET", "URL", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error while creating HTTP GET request to %q: %s", url, err.Error())
	}
	AddRaitoAuthorizationHeaders(req, config);
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while doing HTTP GET to %q: %s", url, err.Error())
	}
	return resp, nil
}

func DoGetToRaito(path string, config *target.BaseTargetConfig) (*http.Response, error) {
	return doGet(url.GetRaitoURL(), path, config)
}
