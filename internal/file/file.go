package file

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/raito-io/cli/internal/target"
	"github.com/raito-io/cli/internal/util/connect"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func CreateUniqueFileName(hint, ext string) string {
	r := rand.Intn(10000000)
	t := time.Now().Format("2006-01-02T15-04-05.999999999Z07-00")

	return hint + "-" + t + "-" + strconv.Itoa(r) + "." + ext
}

// UploadFile uploads the file from the given path.
// It returns the key to use to pass to the Raito backend to use the file.
func UploadFile(file string, config *target.BaseTargetConfig) (string, error) {
	start := time.Now()

	url, key, err := getUploadURL(config)
	if err != nil {
		return "", err
	}

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
func getUploadURL(config *target.BaseTargetConfig) (string, string, error) {
	resp, err := connect.DoGetToRaito("file/upload/signed-url", &config.BaseConfig)
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
