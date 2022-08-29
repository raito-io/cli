package data_access

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/raito-io/cli/internal/file"
	"github.com/raito-io/cli/internal/target"
	"github.com/raito-io/cli/internal/util/connect"
)

type AccessSyncConfig struct {
	target.BaseTargetConfig
}

// RetrieveDataAccessListForDataSource fetches the data access right for a specific data source from the Raito backend.
// It will return 'nil, nil' if no changes happened since the date provided in the 'since' parameter.
// Use 0 for the 'since' parameter if you want to force the fetching of the data access rights.
func RetrieveDataAccessListForDataSource(config *AccessSyncConfig, since int64) (string, error) {
	path := "access-provider/data-source/" + config.DataSourceId
	if since > 0 {
		path += "?since=" + strconv.Itoa(int(since))
	}

	resp, err := connect.DoGetToRaito(path, &config.BaseTargetConfig)
	if err != nil {
		return "", fmt.Errorf("error while fetching access controls for datasource %q: %s", config.DataSourceId, err.Error())
	}

	// Nothing changed since the provided date
	if since > 0 && resp.StatusCode == 304 {
		return "", nil
	}

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("error (HTTP %d) while fetching access controls for datasource %q: %s", resp.StatusCode, config.DataSourceId, resp.Status)
	}
	defer resp.Body.Close()

	cn := strings.Replace(config.ConnectorName, "/", "-", -1)
	filePath, err := filepath.Abs(file.CreateUniqueFileName(cn+"-as", "json"))

	if err != nil {
		return "", err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("error while creating temporary file %q: %s", filePath, err.Error())
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", fmt.Errorf("error while storing access controls for data source %q in file %q: %s", config.DataSourceId, filePath, err.Error())
	}

	return filePath, nil
}
