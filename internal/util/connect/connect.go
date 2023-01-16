package connect

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/raito-io/cli/internal/auth"
	"github.com/raito-io/cli/internal/target"
	"github.com/raito-io/cli/internal/util/url"
	"github.com/raito-io/cli/internal/version"
)

func doPost(host, path, body, contentType string, config *target.BaseConfig) (*http.Response, error) {
	url := url.CreateRaitoURL(host, path)
	config.BaseLogger.Debug("Calling HTTP POST", "URL", url)
	req, err := http.NewRequest("POST", url, strings.NewReader(body))

	if err != nil {
		return nil, fmt.Errorf("error while creating HTTP GET request to %q: %s", url, err.Error())
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", "Raito CLI "+version.GetVersion())

	err = auth.AddToken(req, config)
	if err != nil {
		return nil, fmt.Errorf("error while adding authorization token: %s", err.Error())
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while doing HTTP POST to %q: %s", url, err.Error())
	}

	return resp, nil
}

func DoPostToRaito(path, body, contentType string, config *target.BaseConfig) (*http.Response, error) {
	return doPost(url.GetRaitoURL(), path, body, contentType, config)
}

func doGet(host, path string, config *target.BaseConfig) (*http.Response, error) {
	url := url.CreateRaitoURL(host, path)
	config.BaseLogger.Debug("Calling HTTP GET", "URL", url)
	req, err := http.NewRequest("GET", url, http.NoBody)

	if err != nil {
		return nil, fmt.Errorf("error while creating HTTP GET request to %q: %s", url, err.Error())
	}

	req.Header.Set("User-Agent", "Raito CLI "+version.GetVersion())
	err = auth.AddToken(req, config)

	if err != nil {
		return nil, fmt.Errorf("error while adding authorization token: %s", err.Error())
	}
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error while doing HTTP GET to %q: %s", url, err.Error())
	}

	return resp, nil
}

func DoGetToRaito(path string, config *target.BaseConfig) (*http.Response, error) {
	return doGet(url.GetRaitoURL(), path, config)
}
