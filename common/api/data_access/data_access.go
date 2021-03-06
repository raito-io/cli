package data_access

import (
	"crypto/sha1" //nolint:gosec
	"encoding/base64"
	"net/rpc"
	"sort"
	"strings"

	"github.com/hashicorp/go-plugin"
	"github.com/raito-io/cli/common/api"
	"github.com/raito-io/cli/common/util/config"
)

// DataAccessSyncConfig represents the configuration that is passed from the CLI to the DataAccessSyncer plugin interface.
// It contains all the necessary configuration parameters for the plugin to function.
type DataAccessSyncConfig struct {
	config.ConfigMap
	DataAccess *DataAccessResult
	Prefix     string
	TargetFile string
	RunImport  bool
}

// DataAccessResult is the main structure containing the information coming from Raito describing the data access rules for this data source.
type DataAccessResult struct {
	LastCalculated int64         `yaml:"lastCalculated"`
	AccessRights   []*DataAccess `yaml:"accessRights"`
}

// DataAccess is the structure for one data access element. It has:
//  - Id: the UUID of the data access element. Typically, this is not needed.
//  - DataObject: the data object (e.g. schema, table, column) this data access is applicable to.
//  - Permissions: the list of (Raito) permissions that are granted to the users on the data object.
//  - Users: the list of users the permissions are granted to.
//  - Provider (optional): the Raito Access Provider this data access is generated from. Can be nil.
type DataAccess struct {
	Id          string
	Delete      bool
	NamingHint  string      `yaml:"namingHint"`
	DataObject  *DataObject `yaml:"dataObject"`
	Permissions []string
	Users       []string
	Provider    *Provider
}

// Provider represents the Access Provider that generates the data access.
type Provider struct {
	Name        string
	Description string
	Id          string
}

// DataObject represents the information about a data object. It will refer to a parent data object.
// Parent will be nil if this is a top-level data-object.
type DataObject struct {
	Type   string
	Name   string
	Parent *DataObject
	Path   string `yaml:"-"`
	Source string
}

// CalculateHash calculates a hash for this data access element.
// It's used in the CLI code to flatten a list of data access elements for a data source.
func (d *DataAccess) CalculateHash() string {
	sort.Strings(d.Permissions)
	permissions := strings.Join(d.Permissions, ",")

	path := d.DataObject.Path
	if path == "" {
		path = d.DataObject.BuildPath(".")
	}

	hasher := sha1.New() //nolint:gosec

	hasher.Write([]byte(path))
	hasher.Write([]byte("|"))
	hasher.Write([]byte(permissions))

	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

// Merge merges multiple data access elements together in one.
// It's used in the CLI code to flatten a list of data access elements for a data source.
func (d *DataAccess) Merge(input []*DataAccess) *DataAccess {
	ret := DataAccess{
		DataObject:  d.DataObject,
		Permissions: d.Permissions,
		Users:       d.Users,
	}
	for _, da := range input {
		ret.Users = append(ret.Users, da.Users...)
	}

	// Now remove the duplicates
	check := make(map[string]int)

	for _, val := range ret.Users {
		check[val] = 1
	}
	mergedUsers := make([]string, 0)

	for u := range check {
		mergedUsers = append(mergedUsers, u)
	}
	ret.Users = mergedUsers

	return &ret
}

// BuildPath builds the full path of a data object, using the given separator.
// For example: table 'Employees' in schema 'CompanyX' in database 'Internal' will result in 'Internal.CompanyX.Employees'
// when using a dot (.) as separator.
func (d *DataObject) BuildPath(sep string) string {
	if d.Parent != nil {
		d.Path = d.Parent.BuildPath(sep) + sep + d.Name
	} else {
		d.Path = d.Name
	}

	return d.Path
}

// DataAccessSyncResult represents the result from the data access sync process.
// A potential error is also modeled in here so specific errors remain intact when passed over RPC.
type DataAccessSyncResult struct {
	Error *api.ErrorResult
}

// DataAccessSyncer interface needs to be implemented by any plugin that wants to push data access rules from Raito to its underlying data source.
type DataAccessSyncer interface {
	SyncDataAccess(config *DataAccessSyncConfig) DataAccessSyncResult
}

// DataAccessSyncerPlugin is used on the server (CLI) and client (plugin) side to integrate with the plugin system.
// A plugin should not be using this directly, but instead depend on the cli-plugin-base library to register the plugins.
type DataAccessSyncerPlugin struct {
	Impl DataAccessSyncer
}

func (p *DataAccessSyncerPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &dataAccessSyncerRPCServer{Impl: p.Impl}, nil
}

func (DataAccessSyncerPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &dataAccessSyncerRPC{client: c}, nil
}

// DataAccessSyncerName constant should not be used directly when implementing plugins.
// It's the registration name for the data access syncer plugin,
// used by the CLI and the cli-plugin-base library (RegisterPlugins function) to register the plugins.
const DataAccessSyncerName = "dataAccessSyncer"

type dataAccessSyncerRPC struct{ client *rpc.Client }

func (g *dataAccessSyncerRPC) SyncDataAccess(config *DataAccessSyncConfig) DataAccessSyncResult {
	var resp DataAccessSyncResult

	err := g.client.Call("Plugin.SyncDataAccess", config, &resp)
	if err != nil && resp.Error == nil {
		resp.Error = api.ToErrorResult(err)
	}

	return resp
}

type dataAccessSyncerRPCServer struct {
	Impl DataAccessSyncer
}

func (s *dataAccessSyncerRPCServer) SyncDataAccess(config *DataAccessSyncConfig, resp *DataAccessSyncResult) error {
	*resp = s.Impl.SyncDataAccess(config)
	return nil
}
