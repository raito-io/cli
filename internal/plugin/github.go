package plugin

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"runtime"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/raito-io/cli/internal/config"
	"github.com/raito-io/cli/internal/constants"
	"github.com/spf13/viper"
)

func downloadAndExtractPluginFromGitHubRepo(pluginRequest *pluginRequest, targetPath string, logger hclog.Logger) (string, error) {
	asset, err := getGitHubAsset(pluginRequest, logger)
	if err != nil {
		return "", fmt.Errorf("error looking for plugin to download from Github for %q (version %q): %s", pluginRequest.GroupAndName(), pluginRequest.Version, err.Error())
	}

	if asset == nil {
		return "", nil
	}

	downloadedFile, err := downloadGitHubAsset(pluginRequest, asset.URL, logger)
	if downloadedFile != "" {
		defer os.Remove(downloadedFile)
	}

	if err != nil {
		return "", fmt.Errorf("error downloading plugin from Github for %q (version %q): %s", pluginRequest.GroupAndName(), pluginRequest.Version, err.Error())
	}

	extractedFile, err := extractFromDownloadFile(pluginRequest, downloadedFile, targetPath)
	if err != nil {
		return extractedFile, fmt.Errorf("error extracting plugin binary from release asset from %q: %s", asset.URL, err.Error())
	}

	return extractedFile, nil
}

// getGitHubAsset returns the release asset on github that corresponds with the incoming plugin request and OS+Arch.
// If an error occurs during the search, the error is returned.
// If no asset is found, nil is returned.
func getGitHubAsset(pluginRequest *pluginRequest, logger hclog.Logger) (*gitHubReleaseAsset, error) {
	url := getGitHubReleaseURL(pluginRequest)

	client := retryablehttp.NewClient()
	client.Logger = logger

	request, err := retryablehttp.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Accept", "application/vnd.github.v3+json")

	token, err := findGitHubToken(pluginRequest)
	if err != nil {
		return nil, err
	}

	if token != "" {
		logger.Debug(fmt.Sprintf("found token for repository %q", pluginRequest.Group))
		request.Header.Set("Authorization", "token "+token)
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error while fetching release assets from %q: %s", url, err.Error())
	}

	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading response body from releases request to %q: %s", url, err.Error())
	}

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unable to fetch releases from %q: %s", url, string(respBytes))
	}

	releaseInfo := gitHubReleaseInfo{}

	err = json.Unmarshal(respBytes, &releaseInfo)
	if err != nil {
		return nil, fmt.Errorf("error while parsing response body from releases request to %q: %s", url, err.Error())
	}

	if pluginRequest.IsLatest() {
		// Fill in with the actual version now that we resolved what 'latest' is.
		pluginRequest.Version = strings.TrimPrefix(releaseInfo.TagName, "v")
	}

	return findMatchingGitHubAsset(pluginRequest, &releaseInfo), nil
}

// downloadGitHubAsset downloads the given asset file from github.
// Returns the filename of the file created (if there is one, otherwise, empty string).
// Returns an error if an error occurred.
func downloadGitHubAsset(pluginRequest *pluginRequest, url string, logger hclog.Logger) (string, error) {
	client := retryablehttp.NewClient()
	client.Logger = logger

	request, err := retryablehttp.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	request.Header.Set("Accept", "application/octet-stream")

	token, err := findGitHubToken(pluginRequest)
	if err != nil {
		return "", err
	}

	if token != "" {
		logger.Debug(fmt.Sprintf("found token for repository %q", pluginRequest.Group))
		request.Header.Set("Authorization", "token "+token)
	}

	resp, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("error while fetching release asset from %q: %s", url, err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("error while fetching release asset from %q", url)
	}

	defer resp.Body.Close()

	// Create the file
	out, err := ioutil.TempFile("", "plugin-download-")
	if err != nil {
		return "", fmt.Errorf("error while creating temporary file for asset download: %s", err.Error())
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return out.Name(), fmt.Errorf("error while storing release asset from %q to temporary file: %s", url, err.Error())
	}

	return out.Name(), nil
}

// getGitHubReleaseURL builds the github URL to fetch the assets of a specific release (either latest or a given version).
func getGitHubReleaseURL(pluginRequest *pluginRequest) string {
	url := "https://api.github.com/repos/" + pluginRequest.Group + "/" + pluginRequest.Name + "/releases/"
	if pluginRequest.IsLatest() {
		url += "latest"
	} else {
		url += "tags/"
		if !strings.HasPrefix(pluginRequest.Version, "v") {
			url += "v"
		}
		url += pluginRequest.Version
	}

	return url
}

// findMatchingAsset looks for the asset that matches our OS and architecture.
// Matching happens by looking for a file in the form <name>-<version>-<OS>_<Arch>.tar.gz
// The <version> part in this is ignored (as that matching is already done with the release.
// If nothing is found, nil is returned
func findMatchingGitHubAsset(pluginRequest *pluginRequest, releaseInfo *gitHubReleaseInfo) *gitHubReleaseAsset {
	if releaseInfo.Assets != nil {
		suffix := "-" + runtime.GOOS + "_" + runtime.GOARCH + ".tar.gz"
		prefix := pluginRequest.Name + "-"

		for _, asset := range releaseInfo.Assets {
			if strings.HasPrefix(asset.Name, prefix) && strings.HasSuffix(asset.Name, suffix) {
				return &asset
			}
		}
	}

	// Not found
	return nil
}

func findGitHubToken(pluginRequest *pluginRequest) (string, error) {
	repos := viper.Get(constants.Repositories)
	if repoList, ok := repos.([]interface{}); ok {
		for _, repoObj := range repoList {
			if repo, ok := repoObj.(map[interface{}]interface{}); ok {
				repoName := repo[constants.NameFlag]
				if repoName == pluginRequest.Group {
					if v, f := repo[constants.GitHubToken]; f {
						un, err := config.HandleField(v, reflect.String)

						if err != nil {
							return "", fmt.Errorf("error while handling username field for repository %q: %s", repoName, err.Error())
						}

						if sv, f := un.(string); f {
							return sv, nil
						}
					}
				}
			}
		}
	}

	return "", nil
}

type gitHubReleaseInfo struct {
	Name       string
	TagName    string `json:"tag_name"`
	Prerelease bool
	Message    string
	Assets     []gitHubReleaseAsset
}

type gitHubReleaseAsset struct {
	URL  string
	Name string
}
