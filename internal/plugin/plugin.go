package plugin

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"github.com/raito-io/cli/base/data_object_enricher"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/raito-io/cli/base/access_provider"
	"github.com/raito-io/cli/base/data_source"
	"github.com/raito-io/cli/base/data_usage"
	"github.com/raito-io/cli/base/identity_store"
	plugin2 "github.com/raito-io/cli/base/util/plugin"
)

const LATEST = "latest"

var nameRegexp = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9\-]*[a-zA-Z\d]$`)
var versionRegexp = regexp.MustCompile(`^v?([0-9]+)(\.[0-9]+)(\.[0-9]+)` +
	`(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?` +
	`(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?$`)

var localPluginFolder = "./raito/plugins/"
var globalPluginFolder string

var pluginMap = map[string]plugin.Plugin{
	"identityStoreSyncer": &identity_store.IdentityStoreSyncerPlugin{},
	"dataSourceSyncer":    &data_source.DataSourceSyncerPlugin{},
	"accessSyncer":        &access_provider.AccessSyncerPlugin{},
	"dataUsageSyncer":     &data_usage.DataUsageSyncerPlugin{},
	"info":                &plugin2.InfoPlugin{},
	"dataObjectEnricher":  &data_object_enricher.DataObjectEnricherPlugin{},
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
	GetDataObjectEnricher() (data_object_enricher.DataObjectEnricher, error)
	GetIdentityStoreSyncer() (identity_store.IdentityStoreSyncer, error)
	GetAccessSyncer() (access_provider.AccessSyncer, error)
	GetDataUsageSyncer() (data_usage.DataUsageSyncer, error)
	GetInfo() (plugin2.Info, error)
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
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
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

	info, err := is.GetInfo(context.Background())
	if err != nil {
		return nil, fmt.Errorf("plugininfo: %w", err)
	}

	logger.Debug("Using plugin: " + info.InfoString())

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

	latestVersion := ""
	latestPath := ""

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
			latestPath, latestVersion = getLatestVersionFromFiles(matches)
			logger.Debug(fmt.Sprintf("Found version for plugin %s locally at path %s", pluginRequest.GroupAndName(), latestPath))
		}

		if latestVersion == "" {
			path = globalPluginFolder + pluginRequest.Path()
			matches, err = filepath.Glob(path)
			if len(matches) > 0 && err == nil {
				latestPath, latestVersion = getLatestVersionFromFiles(matches)
				logger.Debug(fmt.Sprintf("Found version for plugin %s globally at path %s", pluginRequest.GroupAndName(), latestPath))
			}
		}
	}

	if latestVersion != "" && pluginRequest.IsLatest() {
		logger.Debug(fmt.Sprintf("A matching plugin found for %s version %s on local disk. Will check if there is a newer version available online.", pluginRequest.GroupAndName(), pluginRequest.Version))

		if latestVersion == LATEST {
			logger.Warn("Using special development version of the plugin. Remove the '-latest' plugin if you want to go back to using the released plugins.")
			return latestPath, nil
		}
	} else {
		logger.Debug(fmt.Sprintf("No matching plugin found for %s version %s on local disk", pluginRequest.GroupAndName(), pluginRequest.Version))
	}

	return downloadAndExtractPluginFromGitHubRepo(pluginRequest, globalPluginFolder, latestVersion, latestPath, logger)
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

func getLatestVersionFromFiles(matches []string) (string, string) {
	versionStart := strings.LastIndex(matches[0], "-") + 1
	prefix := matches[0][0:versionStart]
	versions := make([]string, 0, len(matches))

	for _, match := range matches {
		versions = append(versions, match[versionStart:])
	}

	latestVersion := getLatestVersion(versions)

	return prefix + latestVersion, latestVersion
}

func getLatestVersion(matches []string) string {
	versions := make([]*semver.Version, 0, len(matches))

	for _, match := range matches {
		if match == "latest" {
			return match
		}

		version, err := semver.StrictNewVersion(match)
		if err != nil {
			continue
		}

		versions = append(versions, version)
	}

	sort.SliceStable(versions, func(i, j int) bool {
		v1 := versions[i]
		v2 := versions[j]
		return v1.LessThan(v2)
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
		if !validateVersion(version) {
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

func validateVersion(version string) bool {
	return versionRegexp.MatchString(version)
}

type pluginClientImpl struct {
	client *plugin.Client
}

func (c pluginClientImpl) Close() {
	c.client.Kill()
}

func (c pluginClientImpl) GetDataSourceSyncer() (data_source.DataSourceSyncer, error) {
	raw, err := c.getPlugin(data_source.DataSourceSyncerName)
	if err != nil {
		return nil, err
	}

	if syncer, ok := raw.(data_source.DataSourceSyncer); ok {
		return syncer, nil
	} else {
		return nil, fmt.Errorf("found plugin doesn't correctly implement the DataSourceSyncer interface")
	}
}

func (c pluginClientImpl) GetDataObjectEnricher() (data_object_enricher.DataObjectEnricher, error) {
	raw, err := c.getPlugin(data_object_enricher.DataObjectEnricherName)
	if err != nil {
		return nil, err
	}

	if syncer, ok := raw.(data_object_enricher.DataObjectEnricher); ok {
		return syncer, nil
	} else {
		return nil, fmt.Errorf("found plugin doesn't correctly implement the DataObjectEnricher interface")
	}
}

func (c pluginClientImpl) GetIdentityStoreSyncer() (identity_store.IdentityStoreSyncer, error) {
	raw, err := c.getPlugin(identity_store.IdentityStoreSyncerName)
	if err != nil {
		return nil, err
	}

	if syncer, ok := raw.(identity_store.IdentityStoreSyncer); ok {
		return syncer, nil
	} else {
		return nil, fmt.Errorf("found plugin doesn't correctly implement the IdentityStoreSyncer interface")
	}
}

func (c pluginClientImpl) GetAccessSyncer() (access_provider.AccessSyncer, error) {
	raw, err := c.getPlugin(access_provider.AccessSyncerName)
	if err != nil {
		return nil, err
	}

	if syncer, ok := raw.(access_provider.AccessSyncer); ok {
		return syncer, nil
	} else {
		return nil, fmt.Errorf("found plugin doesn't correctly implement the AccessSyncer interface")
	}
}

func (c pluginClientImpl) GetDataUsageSyncer() (data_usage.DataUsageSyncer, error) {
	raw, err := c.getPlugin(data_usage.DataUsageSyncerName)
	if err != nil {
		return nil, err
	}

	if syncer, ok := raw.(data_usage.DataUsageSyncer); ok {
		return syncer, nil
	} else {
		return nil, fmt.Errorf("found plugin doesn't correctly implement the DataUsageSyncer interface")
	}
}

func (c pluginClientImpl) GetInfo() (plugin2.Info, error) {
	raw, err := c.getPlugin(plugin2.InfoName)
	if err != nil {
		return nil, err
	}

	if info, ok := raw.(plugin2.Info); ok {
		return info, nil
	} else {
		return nil, fmt.Errorf("found plugin doesn't correctly implement the Info interface")
	}
}

func (c pluginClientImpl) getPlugin(plugin string) (interface{}, error) {
	rpcClient, err := c.client.Client()
	if err != nil {
		return nil, err
	}

	return rpcClient.Dispense(plugin)
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
