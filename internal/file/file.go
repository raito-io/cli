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

	"github.com/avast/retry-go/v4"

	"github.com/raito-io/cli/internal/target/types"
	"github.com/raito-io/cli/internal/util/connect"
)

func calculateChecksum(file io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", fmt.Errorf("unable to hash file %q: %w", file, err)
	}

	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

// UploadFile uploads the file from the given path.
// It returns the key to use to pass to the Raito backend to use the file.
func UploadFile(file string, config *types.BaseTargetConfig) (string, error) {
	return uploadHashedFile(file, config, getUploadURL)
}

// UploadLogFile uploads the file from the given path.
// It returns the key to use to pass to the Raito backend to use the file.
func UploadLogFile(file string, config *types.BaseTargetConfig, task string) (string, error) {
	return uploadHashedFile(file, config, func(config *types.BaseTargetConfig, checksum string, fileSize int64) (string, string, map[string][]string, error) {
		return getUploadLogsURL(config, task, checksum, fileSize)
	})
}

func uploadHashedFile(file string, config *types.BaseTargetConfig, uploadURL func(config *types.BaseTargetConfig, checksum string, fileSize int64) (string, string, map[string][]string, error)) (string, error) {
	data, err := os.Open(file)
	if err != nil {
		return "", fmt.Errorf("open file: %w", err)
	}

	defer data.Close()

	checksum, err := calculateChecksum(data)
	if err != nil {
		return "", fmt.Errorf("checksum file: %w", err)
	}

	config.TargetLogger.Debug(fmt.Sprintf("Calculated checksum %q for file %q", checksum, file))

	_, err = data.Seek(0, 0)
	if err != nil {
		return "", fmt.Errorf("seek file: %w", err)
	}

	stats, err := data.Stat()
	if err != nil {
		return "", fmt.Errorf("stat file: %w", err)
	}

	url, key, headers, err := uploadURL(config, checksum, stats.Size())

	if err != nil {
		return "", err
	}

	return uploadFileToBucket(data, url, key, stats.Size(), headers, config)
}

func uploadFileToBucket(data *os.File, url string, key string, contentLength int64, headers map[string][]string, config *types.BaseTargetConfig) (string, error) {
	start := time.Now()

	err := retry.Do(func() error {
		// Ensure to read data from the beginning (in case of retries)
		_, err := data.Seek(0, 0)
		if err != nil {
			return fmt.Errorf("error while seeking file: %s", err.Error())
		}

		req, err := http.NewRequest("PUT", url, data)

		if err != nil {
			return fmt.Errorf("error while executing upload: %s", err.Error())
		}
		req.ContentLength = contentLength

		for headerKey, headerValue := range headers {
			for _, value := range headerValue {
				req.Header.Add(headerKey, value)
			}
		}

		client := &http.Client{}
		res, err := client.Do(req)

		if err != nil {
			return fmt.Errorf("error while executing upload: %s", err.Error())
		}

		defer res.Body.Close()

		if res.StatusCode >= 300 {
			buf, _ := io.ReadAll(res.Body)

			return fmt.Errorf("error (HTTP %d) while executing upload: %s - %s", res.StatusCode, res.Status, string(buf))
		}

		return nil
	}, retry.Attempts(3), retry.DelayType(retry.BackOffDelay), retry.OnRetry(func(attempt uint, err error) {
		config.TargetLogger.Warn(fmt.Sprintf("Failed to upload file with key %q. Will retry (%d/3): %s", key, attempt, err.Error()))
	}))

	if err != nil {
		return "", err
	}

	sec := time.Since(start).Round(time.Millisecond)

	config.TargetLogger.Info(fmt.Sprintf("Successfully uploaded file with key %q (%d bytes) in %s.", key, contentLength, sec))

	return key, nil
}

// GetUploadURL creates an S3 URL to upload a file to.
// It returns the upload URL and the file key to use to pass to the Raito backend to use it.
// Returns two empty strings if something went wrong (error logged)
func getUploadURL(config *types.BaseTargetConfig, checksum string, fileSize int64) (string, string, map[string][]string, error) {
	params := url.Values{}
	params.Add("sha256", checksum)
	params.Add("contentLength", fmt.Sprintf("%d", fileSize))

	return getUploadUrlAndKey(config, "file/upload/signed-url?"+params.Encode())
}

func getUploadLogsURL(config *types.BaseTargetConfig, task string, checksum string, fileSize int64) (string, string, map[string][]string, error) {
	params := url.Values{}
	params.Add("task", task)
	params.Add("sha256", checksum)
	params.Add("contentLength", fmt.Sprintf("%d", fileSize))

	return getUploadUrlAndKey(config, "file/upload/logs/signed-url?"+params.Encode())
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
