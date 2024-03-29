package file

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/raito-io/cli/internal/target/types"
	"github.com/raito-io/cli/internal/util/connect"
)

// UploadFile uploads the file from the given path.
// It returns the key to use to pass to the Raito backend to use the file.
func UploadFile(file string, config *types.BaseTargetConfig) (string, error) {
	url, key, err := getUploadURL(config)
	if err != nil {
		return "", err
	}

	return uploadFileToBucket(file, config, url, key)
}

// UploadLogFile uploads the file from the given path.
// It returns the key to use to pass to the Raito backend to use the file.
func UploadLogFile(file string, config *types.BaseTargetConfig, task string) (string, error) {
	url, key, err := getUploadLogsURL(config, task)
	if err != nil {
		return "", err
	}

	return uploadFileToBucket(file, config, url, key)
}

func uploadFileToBucket(file string, config *types.BaseTargetConfig, url string, key string) (string, error) {
	start := time.Now()

	data, err := os.Open(file)
	if err != nil {
		return "", fmt.Errorf("unable to open file %q: %s", file, err.Error())
	}

	defer data.Close()
	stat, err := data.Stat()

	if err != nil {
		return "", fmt.Errorf("error while getting file size of %q: %s", file, err.Error())
	}
	req, err := http.NewRequest("PUT", url, data)

	if err != nil {
		return "", fmt.Errorf("error while executing upload: %s", err.Error())
	}
	req.ContentLength = stat.Size()

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return "", fmt.Errorf("error while executing upload: %s", err.Error())
	}

	defer res.Body.Close()

	if res.StatusCode >= 300 {
		buf, _ := io.ReadAll(res.Body)

		return "", fmt.Errorf("error (HTTP %d) while executing upload: %s - %s", res.StatusCode, res.Status, string(buf))
	}

	sec := time.Since(start).Round(time.Millisecond)

	config.TargetLogger.Info(fmt.Sprintf("Successfully uploaded file with key %q (%d bytes) in %s.", key, stat.Size(), sec))

	return key, nil
}

// GetUploadURL creates an S3 URL to upload a file to.
// It returns the upload URL and the file key to use to pass to the Raito backend to use it.
// Returns two empty strings if something went wrong (error logged)
func getUploadURL(config *types.BaseTargetConfig) (string, string, error) {
	return getUploadUrlAndKey(config, "file/upload/signed-url")
}

func getUploadLogsURL(config *types.BaseTargetConfig, task string) (string, string, error) {
	return getUploadUrlAndKey(config, "file/upload/logs/signed-url?task="+task)
}

func getUploadUrlAndKey(config *types.BaseTargetConfig, path string) (string, string, error) {
	resp, err := connect.DoGetToRaito(path, &config.BaseConfig)
	if err != nil {
		return "", "", fmt.Errorf("error while trying to get a signed upload URL: %s", err.Error())
	}

	if resp.StatusCode >= 300 {
		return "", "", fmt.Errorf("error (HTTP %d) while trying to get a signed upload URL: %s", resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("error while reading result body for getting signed url: %s", err.Error())
	}
	var result signedURL

	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", "", fmt.Errorf("error while parsing result body for getting signed url: %s", err.Error())
	}

	return result.URL, result.Key, nil
}

type signedURL struct {
	URL string
	Key string
}
