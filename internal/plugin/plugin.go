package plugin

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/raito-io/cli/common/api"
	"github.com/raito-io/cli/common/api/data_access"
	"github.com/raito-io/cli/common/api/data_source"
	"github.com/raito-io/cli/common/api/data_usage"
	"github.com/raito-io/cli/common/api/identity_store"
)

// TODO add cancel and async (done context) support

const LATEST = "latest"

var nameRegexp = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9\-]*[a-zA-Z\d]$`)
var versionRegexp = regexp.MustCompile(`^\d+\.\d+\.\d+$`)

var localPluginFolder = "./raito/plugins/"
var globalPluginFolder string

var pluginMap = map[string]plugin.Plugin{
	"identityStoreSyncer": &identity_store.IdentityStoreSyncerPlugin{},
	"dataSourceSyncer":    &data_source.DataSourceSyncerPlugin{},
	"dataAccessSyncer":    &data_access.DataAccessSyncerPlugin{},
	"dataUsageSyncer":     &data_usage.DataUsageSyncerPlugin{},
	"info":                &api.InfoPlugin{},
}

func init() {
	userHome, _ := os.UserHomeDir()
	if !strings.HasSuffix(userHome, "/") {
		userHome += "/"
	}
	globalPluginFolder = userHome + ".raito/plugins/"

	if _, err := os.Stat(globalPluginFolder); err != nil {
		if err := os.MkdirAll(globalPluginFolder, os.ModePerm); err != nil {
			hclog.L().Error("Error creating global Raito folder %q. Make sure permissions are set correctly: %s", globalPluginFolder, err.Error())
		}
	}
}

type PluginClient interface {
	Close()
	GetDataSourceSyncer() (data_source.DataSourceSyncer, error)
	GetIdentityStoreSyncer() (identity_store.IdentityStoreSyncer, error)
	GetDataAccessSyncer() (data_access.DataAccessSyncer, error)
	GetDataUsageSyncer() (data_usage.DataUsageSyncer, error)
	GetInfo() (api.Info, error)
}

func NewPluginClient(connector string, version string, logger hclog.Logger) (PluginClient, error) {
	pluginPath, err := findMatchingPlugin(connector, version, logger)
	if err != nil {
		return nil, fmt.Errorf("error while finding matching plugin for %q (version %q): %s", connector, version, err.Error())
	}

	if pluginPath == "" {
		return nil, fmt.Errorf("unable to find matching plugin for %q (version %q)", connector, version)
	}

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command(pluginPath),
		Logger:          logger,
	})

	// Connecting to see if it works...
	_, err = client.Client()
	if err != nil {
		return nil, fmt.Errorf("error connecting to plugin %q. It may be corrupt or invalid", connector)
	}

	pci := pluginClientImpl{client}

	is, err := pci.GetInfo()
	if err != nil {
		return nil, fmt.Errorf("the plugin (%s) doesn't correctly implement the necessary interfaces", connector)
	}

	info := is.PluginInfo()
	logger.Debug("Using plugin: " + info.String())

	return pci, nil
}

// findMatchingPlugin looks for a plugin with the given connector name and version
// When version is empty, 'latest' is presumed. If versions are found locally, the latest will be used.
// If nothing is found locally, the correct version (or latest) will be downloaded from the first registry that has a hit.
//
// When no matching plugin is found (locally or in registries), an empty string is returned.
// If an error occurred during the search, it will be returned.
func findMatchingPlugin(connector string, version string, logger hclog.Logger) (string, error) {
	pluginRequest, err := parsePluginRequest(connector, version)
	if err != nil {
		return "", err
	}

	logger.Debug(fmt.Sprintf("Using plugin request %+v", pluginRequest))

	// We're looking for a specific plugin version
	if !pluginRequest.IsLatest() {
		// Look locally
		path := localPluginFolder + pluginRequest.Path()
		if _, err := os.Stat(path); err == nil {
			logger.Debug(fmt.Sprintf("Found match for plugin %s version %s locally at path %s", pluginRequest.GroupAndName(), pluginRequest.Version, path))
			return path, nil
		}

		// Look globally
		path = globalPluginFolder + pluginRequest.Path()
		if _, err := os.Stat(path); err == nil {
			logger.Debug(fmt.Sprintf("Found match for plugin %s version %s globally at path %s", pluginRequest.GroupAndName(), pluginRequest.Version, path))
			return path, nil
		}
	} else {
		path := localPluginFolder + pluginRequest.Path()
		matches, err := filepath.Glob(path)
		if len(matches) > 0 && err == nil {
			latest := getLatestVersionFromFiles(matches)
			logger.Debug(fmt.Sprintf("Found latest version for plugin %s locally at path %s", pluginRequest.GroupAndName(), latest))
			return latest, nil
		}

		path = globalPluginFolder + pluginRequest.Path()
		matches, err = filepath.Glob(path)
		if len(matches) > 0 && err == nil {
			latest := getLatestVersionFromFiles(matches)
			logger.Debug(fmt.Sprintf("Found latest version for plugin %s globally at path %s", pluginRequest.GroupAndName(), latest))
			return latest, nil
		}
	}

	logger.Debug(fmt.Sprintf("No matching plugin found for %s version %s on local disk", pluginRequest.GroupAndName(), pluginRequest.Version))

	return downloadAndExtractPluginFromGitHubRepo(pluginRequest, globalPluginFolder, logger)
}

func extractFromDownloadFile(pluginRequest *pluginRequest, downloadedFile, targetPath string) (string, error) {
	tarGzFile, err := os.Open(downloadedFile)
	if err != nil {
		return "", fmt.Errorf("error while reading tar.gz archive %s: %s", downloadedFile, err.Error())
	} else {
		defer tarGzFile.Close()
		extractedFile := targetPath + pluginRequest.Path()
		extractedFile, err := extractTarGz(tarGzFile, extractedFile)

		if err != nil {
			return "", fmt.Errorf("error while reading tar.gz archive %s: %s", downloadedFile, err.Error())
		} else {
			if err := os.Chmod(extractedFile, 0750); err != nil {
				return "", fmt.Errorf("error while setting the right permissions for plugin file %q: %s", extractedFile, err.Error())
			}
			return extractedFile, nil
		}
	}
}

func extractTarGz(gzipStream io.Reader, extractedPath string) (string, error) {
	parentFolder := extractedPath[0 : strings.LastIndex(extractedPath, "/")+1]

	err := os.MkdirAll(parentFolder, fs.ModePerm)
	if err != nil {
		return "", fmt.Errorf("error while creating plugin parent folder %q: %s", parentFolder, err.Error())
	}

	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return "", fmt.Errorf("error while reading gzip stream: %s", err.Error())
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return "", fmt.Errorf("error while extracting file from tar.gz archive: %s", err.Error())
		}

		switch header.Typeflag {
		case tar.TypeDir:
			return "", fmt.Errorf("found directories in the tar.gz archive")
		case tar.TypeReg:
			// goreleaser will also include the LICENSE AND README files, we ignore them.
			// we also ignore other files that are 1MB as they cannot be the binary we are looking for
			if header.Name == "LICENSE" || header.Name == "README" || header.Size < 1024*1024 {
				continue
			}
			outFile, err := os.Create(extractedPath)

			if err != nil {
				return "", fmt.Errorf("error while extracting file from tar.gz archive: %s", err.Error())
			}

			for {
				if _, err := io.CopyN(outFile, tarReader, 1024); err != nil {
					if err != nil {
						if err == io.EOF {
							break
						}

						return "", fmt.Errorf("error while extracting file from tar.gz archive: %s", err.Error())
					}
				}
			}

			outFile.Close()

			return extractedPath, nil
		default:
			return "", errors.New("unknown entry found in tar.gz archive")
		}
	}

	return "", errors.New("no files found to extract from tar.gz archive")
}

func getLatestVersionFromFiles(matches []string) string {
	versionStart := strings.LastIndex(matches[0], "-") + 1
	prefix := matches[0][0:versionStart]
	versions := make([]string, 0, len(matches))

	for _, match := range matches {
		versions = append(versions, match[versionStart:])
	}

	return prefix + getLatestVersion(versions)
}

func getLatestVersion(matches []string) string {
	versions := make([]api.Version, 0, len(matches))

	for _, match := range matches {
		if match == "latest" {
			return match
		}

		versions = append(versions, api.ParseVersion(match))
	}

	sort.SliceStable(versions, func(i, j int) bool {
		v1 := versions[i]
		v2 := versions[j]
		if v1.Major < v2.Major {
			return true
		}
		if v1.Minor < v2.Minor {
			return true
		}
		if v1.Maintenance < v2.Maintenance {
			return true
		}
		return false
	})

	return versions[len(versions)-1].String()
}

func parsePluginRequest(connector string, version string) (*pluginRequest, error) {
	if strings.Count(connector, "/") != 1 {
		return nil, errors.New("the connector name is expected to have exactly 1 slash (/) in it")
	}

	parts := strings.Split(connector, "/")
	if len(parts) != 2 {
		return nil, errors.New("the connector name should be in the format <group>/<name>")
	}

	if parts[0] == "" {
		return nil, errors.New("no connector group specified. The connector name should be in the format <group>/<name>")
	}

	if parts[1] == "" {
		return nil, errors.New("no connector name specified. The connector name should be in the format <group>/<name>")
	}

	group := parts[0]
	if !validateName(group) {
		return nil, errors.New("the connector group should only contain alphanumeric characters and dash (-) characters (not at the start or the end)")
	}

	name := parts[1]
	if !validateName(name) {
		return nil, errors.New("the connector name should only contain alphanumeric characters and dash (-) characters (not at the start or the end)")
	}

	// Handling the version
	version = strings.TrimSpace(version)
	if version == "" {
		version = LATEST
	}

	version = strings.ToLower(version)
	if version != LATEST {
		if !versionRegexp.MatchString(version) {
			return nil, errors.New("the connector version should either be empty, 'latest' or in the format X.Y.Z")
		}
	}

	return &pluginRequest{
		Name:    name,
		Group:   group,
		Version: version,
	}, nil
}

func validateName(name string) bool {
	return nameRegexp.MatchString(name)
}

type pluginClientImpl struct {
	client *plugin.Client
}

func (c pluginClientImpl) Close() {
	c.client.Kill()
}

// TODO once Go Generics are released, these 4 methodes can probably be implemented in 1 helper

func (c pluginClientImpl) GetDataSourceSyncer() (data_source.DataSourceSyncer, error) {
	rpcClient, err := c.client.Client()
	if err != nil {
		return nil, err
	}

	raw, err := rpcClient.Dispense(data_source.DataSourceSyncerName)
	if err != nil {
		return nil, err
	}

	if syncer, ok := raw.(data_source.DataSourceSyncer); ok {
		return syncer, nil
	} else {
		return nil, fmt.Errorf("found plugin doesn't correctly implement the DataSourceSyncer interface")
	}
}

func (c pluginClientImpl) GetIdentityStoreSyncer() (identity_store.IdentityStoreSyncer, error) {
	rpcClient, err := c.client.Client()
	if err != nil {
		return nil, err
	}

	raw, err := rpcClient.Dispense(identity_store.IdentityStoreSyncerName)
	if err != nil {
		return nil, err
	}

	if syncer, ok := raw.(identity_store.IdentityStoreSyncer); ok {
		return syncer, nil
	} else {
		return nil, fmt.Errorf("found plugin doesn't correctly implement the IdentityStoreSyncer interface")
	}
}

func (c pluginClientImpl) GetDataAccessSyncer() (data_access.DataAccessSyncer, error) {
	rpcClient, err := c.client.Client()
	if err != nil {
		return nil, err
	}

	raw, err := rpcClient.Dispense(data_access.DataAccessSyncerName)
	if err != nil {
		return nil, err
	}

	if syncer, ok := raw.(data_access.DataAccessSyncer); ok {
		return syncer, nil
	} else {
		return nil, fmt.Errorf("found plugin doesn't correctly implement the DataAccessSyncer interface")
	}
}

func (c pluginClientImpl) GetDataUsageSyncer() (data_usage.DataUsageSyncer, error) {
	rpcClient, err := c.client.Client()
	if err != nil {
		return nil, err
	}

	raw, err := rpcClient.Dispense(data_usage.DataUsageSyncerName)
	if err != nil {
		return nil, err
	}

	if syncer, ok := raw.(data_usage.DataUsageSyncer); ok {
		return syncer, nil
	} else {
		return nil, fmt.Errorf("found plugin doesn't correctly implement the DataUsageSyncer interface")
	}
}

func (c pluginClientImpl) GetInfo() (api.Info, error) {
	rpcClient, err := c.client.Client()
	if err != nil {
		return nil, err
	}

	raw, err := rpcClient.Dispense(api.InfoName)
	if err != nil {
		return nil, err
	}

	if info, ok := raw.(api.Info); ok {
		return info, nil
	} else {
		return nil, fmt.Errorf("found plugin doesn't correctly implement the Info interface")
	}
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "RAITO_CLI_PLUGIN",
	MagicCookieValue: "Raito Handshake!",
}

type pluginRequest struct {
	Group   string
	Name    string
	Version string
}

func (r *pluginRequest) IsLatest() bool {
	return r.Version == LATEST || r.Version == ""
}

func (r *pluginRequest) GroupAndName() string {
	return r.Group + "/" + r.Name
}

func (r *pluginRequest) Path() string {
	if r.IsLatest() {
		return r.GroupAndName() + "-*"
	}

	return r.GroupAndName() + "-" + r.Version
}

func (r *pluginRequest) RemoteFilePath() string {
	sb := strings.Builder{}
	sb.WriteString(r.GroupAndName())
	sb.WriteString("/")
	sb.WriteString(r.Name)
	sb.WriteString("-")
	sb.WriteString(runtime.GOOS)
	sb.WriteString("_")
	sb.WriteString(runtime.GOARCH)
	sb.WriteString("-")
	sb.WriteString(r.Version)
	sb.WriteString(".tar.gz")

	return sb.String()
}
