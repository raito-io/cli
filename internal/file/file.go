package file

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/raito-io/cli/internal/target/types"
	"github.com/raito-io/cli/internal/util/connect"
)

func hashFile(file string) (string, int64, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", 0, fmt.Errorf("unable to open file %q: %w", file, err)
	}

	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return "", 0, fmt.Errorf("unable to get file info of %q: %w", file, err)
	}

	fileSize := fi.Size()

	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil {
		return "", 0, fmt.Errorf("unable to hash file %q: %w", file, err)
	}

	return base64.StdEncoding.EncodeToString(h.Sum(nil)), fileSize, nil
}

// UploadFile uploads the file from the given path.
// It returns the key to use to pass to the Raito backend to use the file.
func UploadFile(file string, config *types.BaseTargetConfig) (string, error) {
	hash, fileSize, err := hashFile(file)
	if err != nil {
		return "", fmt.Errorf("hash file: %w", err)
	}

	url, key, headers, err := getUploadURL(config, hash, fileSize)

	config.TargetLogger.Info(fmt.Sprintf("Uploading file %q to %q", file, url))

	if err != nil {
		return "", err
	}

	return uploadFileToBucket(file, config, url, key, headers)
}

// UploadLogFile uploads the file from the given path.
// It returns the key to use to pass to the Raito backend to use the file.
func UploadLogFile(file string, config *types.BaseTargetConfig, task string) (string, error) {
	url, key, _, err := getUploadLogsURL(config, task)
	if err != nil {
		return "", err
	}

	return uploadFileToBucket(file, config, url, key, map[string][]string{})
}

func uploadFileToBucket(file string, config *types.BaseTargetConfig, url string, key string, headers map[string][]string) (string, error) {
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

	for headerKey, headerValue := range headers {
		for _, value := range headerValue {
			req.Header.Add(headerKey, value)
		}
	}

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
func getUploadURL(config *types.BaseTargetConfig, hash string, fileSize int64) (string, string, map[string][]string, error) {
	params := url.Values{}
	params.Add("sha256", hash)
	params.Add("contentLength", fmt.Sprintf("%d", fileSize))

	return getUploadUrlAndKey(config, "file/upload/signed-url?"+params.Encode())
}

func getUploadLogsURL(config *types.BaseTargetConfig, task string) (string, string, map[string][]string, error) {
	return getUploadUrlAndKey(config, "file/upload/logs/signed-url?task="+task)
}

func getUploadUrlAndKey(config *types.BaseTargetConfig, path string) (string, string, map[string][]string, error) {
	resp, err := connect.DoGetToRaito(path, &config.BaseConfig)
	if err != nil {
		return "", "", nil, fmt.Errorf("error while trying to get a signed upload URL: %s", err.Error())
	}

	if resp.StatusCode >= 300 {
		return "", "", nil, fmt.Errorf("error (HTTP %d) while trying to get a signed upload URL: %s", resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", nil, fmt.Errorf("error while reading result body for getting signed url: %s", err.Error())
	}
	var result signedURL

	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", "", nil, fmt.Errorf("error while parsing result body for getting signed url: %s", err.Error())
	}

	return result.URL, result.Key, result.SignedHeaders, nil
}

type signedURL struct {
	URL           string
	Key           string
	SignedHeaders map[string][]string `json:"signedHeaders,omitempty"`
}
