// Package identity_store provides the tooling to build the Raito identity store import files.
// Simply use the NewIdentityStoreFileCreator function by passing in the config coming from the CLI
// to create the necessary files.
// The returned IdentityStoreFileCreator can then be used (using the AddUsers and AddGroups functions)
// to write the users and groups to the right file.
// Make sure to call the Close function on the creator at the end (tip: use defer).
package identity_store

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/raito-io/cli/base/tag"
)

//go:generate go run github.com/vektra/mockery/v2 --name=IdentityStoreFileCreator --with-expecter

// Group represents a user group in the format that is suitable to be imported into a Raito identity store.
type Group struct {
	ExternalId             string     `json:"externalId"`
	Name                   string     `json:"name"`
	DisplayName            string     `json:"displayName"`
	Description            string     `json:"description"`
	ParentGroupExternalIds []string   `json:"parentGroupExternalIds"`
	Tags                   []*tag.Tag `json:"tags"`
}

// User represents a user in the format that is suitable to be imported into a Raito identity store.
type User struct {
	ExternalId       string     `json:"externalId"`
	Name             string     `json:"name"`
	UserName         string     `json:"userName"`
	Email            string     `json:"email"`
	GroupExternalIds []string   `json:"groupExternalIds"`
	Tags             []*tag.Tag `json:"tags"`
}

// IdentityStoreFileCreator describes the interface for easily creating the user and group import files
// to be imported by the Raito CLI.
type IdentityStoreFileCreator interface {
	AddGroups(groups ...*Group) error
	AddUsers(users ...*User) error
	Close()
	GetUserCount() int
	GetGroupCount() int
}

type identityStoreImporter struct {
	config *IdentityStoreSyncConfig

	usersFile  *os.File
	groupsFile *os.File
	userCount  int
	groupCount int
}

// NewIdentityStoreFileCreator creates a new IdentityStoreFileCreator based on the configuration coming from
// the Raito CLI.
func NewIdentityStoreFileCreator(config *IdentityStoreSyncConfig) (IdentityStoreFileCreator, error) {
	isI := identityStoreImporter{
		config: config,
	}

	err := isI.createTargetFiles()
	if err != nil {
		return nil, err
	}

	_, err = isI.usersFile.WriteString("[")

	if err != nil {
		return nil, err
	}

	_, err = isI.groupsFile.WriteString("[")

	if err != nil {
		return nil, err
	}

	return &isI, nil
}

// Close finalizes the import files and closes them so they can be correctly read by the Raito CLI.
// This method must be called when all users and groups have been added and before control is given back
// to the CLI. It's advised to call this using 'defer'.
func (i *identityStoreImporter) Close() {
	i.usersFile.WriteString("\n]")  //nolint:errcheck
	i.groupsFile.WriteString("\n]") //nolint:errcheck

	i.usersFile.Close()
	i.groupsFile.Close()
}

// AddGroups adds the slice of groups to the groups import file.
// It returns an error when writing one of the groups fails (it will not process the other groups after that).
// It returns nil if everything went well.
func (i *identityStoreImporter) AddGroups(groups ...*Group) error {
	if len(groups) == 0 {
		return nil
	}

	for _, g := range groups {
		var err error

		if i.groupCount > 0 {
			i.groupsFile.WriteString(",") //nolint:errcheck
		}

		i.groupsFile.WriteString("\n") //nolint:errcheck

		gBuf, _ := json.Marshal(g)

		if err != nil {
			return fmt.Errorf("error while serializing group with externalID %q", g.ExternalId)
		}

		i.groupsFile.WriteString("\n") //nolint:errcheck
		_, err = i.groupsFile.Write(gBuf)

		// Only looking at writing errors at the end, supposing if one fails, all would fail
		if err != nil {
			return fmt.Errorf("error while writing to temp file %q", i.groupsFile.Name())
		}

		i.groupCount++
	}

	return nil
}

// AddUsers adds the slice of users to the users import file.
// It returns an error when writing one of the users fails (it will not process the other users after that).
// It returns nil if everything went well.
func (i *identityStoreImporter) AddUsers(users ...*User) error {
	if len(users) == 0 {
		return nil
	}

	for _, u := range users {
		var err error

		if i.userCount > 0 {
			i.usersFile.WriteString(",") //nolint:errcheck
		}

		i.usersFile.WriteString("\n") //nolint:errcheck

		uBuf, _ := json.Marshal(u)

		if err != nil {
			return fmt.Errorf("error while serializing user with externalID %q", u.ExternalId)
		}

		i.usersFile.WriteString("\n") //nolint:errcheck
		_, err = i.usersFile.Write(uBuf)

		// Only looking at writing errors at the end, supposing if one fails, all would fail
		if err != nil {
			return fmt.Errorf("error while writing to temp file %q", i.usersFile.Name())
		}

		i.userCount++
	}

	return nil
}

// GetUserCount returns the number of users that has been added to the import file.
func (i *identityStoreImporter) GetUserCount() int {
	return i.userCount
}

// GetGroupCount returns the number of groups that has been added to the import file.
func (i *identityStoreImporter) GetGroupCount() int {
	return i.groupCount
}

func (i *identityStoreImporter) createTargetFiles() error {
	f, err := os.Create(i.config.UserFile)
	if err != nil {
		return fmt.Errorf("error creating temporary file for identity store importer (users): %s", err.Error())
	}

	i.usersFile = f

	f2, err := os.Create(i.config.GroupFile)
	if err != nil {
		return fmt.Errorf("error creating temporary file for identity store importer (groups): %s", err.Error())
	}

	i.groupsFile = f2

	return nil
}
