// Package sync_from_target provides the tooling to build the file to export access providers from the data source to be imported into Raito.
// Simply use the NewAccessProviderFileCreator function by passing in the config coming from the CLI
// to create the necessary file(s).
// The returned AccessProviderFileCreator can then be used (using the AddAccessProvider function)
// to write AccessProvider to the file.
// Make sure to call the Close function on the creator at the end (tip: use defer).
package sync_from_target

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/smithy-go/ptr"
	"github.com/raito-io/cli/base/util/match"
	"github.com/raito-io/cli/internal/constants"

	"github.com/raito-io/cli/base/access_provider"
	error2 "github.com/raito-io/cli/base/util/error"
)

//go:generate go run github.com/vektra/mockery/v2 --name=AccessProviderFileCreator --with-expecter

// AccessProviderFileCreator describes the interface for easily creating the access object import files
// to be imported by the Raito CLI.
type AccessProviderFileCreator interface {
	AddAccessProviders(dataAccessList ...*AccessProvider) error
	Close()
	GetAccessProviderCount() int
}

type accessProviderFileCreator struct {
	config *access_provider.AccessSyncFromTarget

	targetFile      *os.File
	dataAccessCount int
}

// NewAccessProviderFileCreator creates a new AccessProviderFileCreator based on the configuration coming from
// the Raito CLI.
func NewAccessProviderFileCreator(config *access_provider.AccessSyncFromTarget) (AccessProviderFileCreator, error) {
	dsI := accessProviderFileCreator{
		config: config,
	}

	err := dsI.createTargetFile()
	if err != nil {
		return nil, err
	}

	_, err = dsI.targetFile.WriteString("[")
	if err != nil {
		return nil, err
	}

	return &dsI, nil
}

// Close finalizes the import file and close it so it can be correctly read by the Raito CLI.
// This method must be called when all data objects have been added and before control is given back
// to the CLI. It's advised to call this using 'defer'.
func (d *accessProviderFileCreator) Close() {
	d.targetFile.WriteString("\n]") //nolint:errcheck
	d.targetFile.Close()
}

func shouldLock(field string, all bool, lockByName []string, lockByTag []string, whenIncomplete bool, ap *AccessProvider) (bool, error) {
	if all {
		return true, nil
	}

	matched, err := match.MatchesAny(ap.Name, lockByName)
	if err != nil {
		return false, fmt.Errorf("parsing %s-by-name: %s", field, err.Error())
	}

	if matched {
		return true, nil
	}

	if len(lockByTag) > 0 && len(ap.Tags) > 0 {
		for _, tag := range ap.Tags {
			fullTag := fmt.Sprintf("%s:%s", tag.Key, tag.Value)

			matched, err = match.MatchesAny(fullTag, lockByTag)
			if err != nil {
				return false, fmt.Errorf("parsing %s-by-tag: %s", field, err.Error())
			}

			if matched {
				return true, nil
			}
		}
	}

	if whenIncomplete && ap.Incomplete != nil && *ap.Incomplete {
		return true, nil
	}

	return false, nil
}

// AddAccessProviders adds the slice of data access elements to the import file.
// It returns an error when writing one of the objects fails (it will not process the other data objects after that).
// It returns nil if everything went well.
func (d *accessProviderFileCreator) AddAccessProviders(accessProviders ...*AccessProvider) error {
	if len(accessProviders) == 0 {
		return nil
	}

	for _, ap := range accessProviders {
		if ap.WhoLocked == nil || !*ap.WhoLocked {
			l, err := shouldLock("lock-who", d.config.LockAllWho, d.config.LockWhoByName, d.config.LockWhoByTag, d.config.LockWhoWhenIncomplete, ap)
			if err != nil {
				return err
			}

			if l {
				ap.WhoLocked = ptr.Bool(true)
			}
		}

		if ap.InheritanceLocked == nil || !*ap.InheritanceLocked {
			l, err := shouldLock("lock-inheritance", d.config.LockAllInheritance, d.config.LockInheritanceByName, d.config.LockInheritanceByTag, d.config.LockInheritanceWhenIncomplete, ap)
			if err != nil {
				return err
			}

			if l {
				ap.InheritanceLocked = ptr.Bool(true)
			}
		}

		if ap.WhatLocked == nil || !*ap.WhatLocked {
			l, err := shouldLock("lock-what", d.config.LockAllWhat, d.config.LockWhatByName, d.config.LockWhatByTag, d.config.LockWhatWhenIncomplete, ap)
			if err != nil {
				return err
			}

			if l {
				ap.WhatLocked = ptr.Bool(true)
			}
		}

		if ap.NameLocked == nil || !*ap.NameLocked {
			l, err := shouldLock("lock-names", d.config.LockAllNames, d.config.LockNamesByName, d.config.LockNamesByTag, d.config.LockNamesWhenIncomplete, ap)
			if err != nil {
				return err
			}

			if l {
				ap.NameLocked = ptr.Bool(true)
			}
		}

		if ap.DeleteLocked == nil || !*ap.DeleteLocked {
			l, err := shouldLock("lock-delete", d.config.LockAllDelete, d.config.LockDeleteByName, d.config.LockDeleteByTag, d.config.LockDeleteWhenIncomplete, ap)
			if err != nil {
				return err
			}

			if l {
				ap.DeleteLocked = ptr.Bool(true)
			}
		}

		if d.config.LockAllOwners {
			ap.OwnersLocked = ptr.Bool(true)
		}

		var err error

		// The legacy field is still supported for now, but will be removed in the future
		shouldMakeNonInternalizable, err := match.MatchesAny(ap.Name, d.config.MakeNotInternalizable)
		if err != nil {
			return fmt.Errorf("parsing parameter %q: %s", constants.MakeNotInternalizableFlag, err.Error())
		}

		if shouldMakeNonInternalizable {
			ap.NotInternalizable = true
		}

		// The new fully lock options to be checked
		l, err := shouldLock("fully-lock", d.config.FullyLockAll, d.config.FullyLockByName, d.config.FullyLockByTag, d.config.FullyLockWhenIncomplete, ap)
		if err != nil {
			return err
		}

		if l {
			ap.NotInternalizable = true
		}

		// TODO REFACTOR to be removed once the old API is removed
		// This now makes sure we send the new model (no more Access layer) to Raito cloud.
		if len(ap.Access) > 0 && ap.What == nil {
			ap.What = ap.Access[0].What
			ap.ActualName = ap.Access[0].ActualName
			ap.Access = nil
		}

		if d.dataAccessCount > 0 {
			d.targetFile.WriteString(",") //nolint:errcheck
		}

		d.targetFile.WriteString("\n") //nolint:errcheck

		doBuf, err := json.Marshal(ap)
		if err != nil {
			return fmt.Errorf("error while serializing data object with externalID %q", ap.ExternalId)
		}

		d.targetFile.WriteString("\n") //nolint:errcheck
		_, err = d.targetFile.Write(doBuf)

		// Only looking at writing errors at the end, supposing if one fails, all would fail
		if err != nil {
			return fmt.Errorf("error while writing to temp file %q", d.targetFile.Name())
		}

		d.dataAccessCount++
	}

	return nil
}

// GetAccessProviderCount returns the number of access elements that have been added to the import file.
func (d *accessProviderFileCreator) GetAccessProviderCount() int {
	return d.dataAccessCount
}

func (d *accessProviderFileCreator) createTargetFile() error {
	f, err := os.Create(d.config.TargetFile)
	if err != nil {
		return error2.CreateErrorFileError(d.config.TargetFile, err)
	}
	d.targetFile = f

	return nil
}

var actionNames = [...]string{"Promise", "Grant", "Deny", "Mask", "Filtered"}
var actionNameMap = map[string]Action{"Promise": Promise, "Grant": Grant, "Deny": Deny, "Mask": Mask, "Filtered": Filtered}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *Action) UnmarshalJSON(b []byte) error {
	var j string

	err := json.Unmarshal(b, &j)
	if err != nil {
		fmt.Println(err.Error()) //nolint:forbidigo
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = actionNameMap[j]

	return nil
}

// MarshalJSON marshals the enum as a quoted json string
func (s Action) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(actionNames[s])
	buffer.WriteString(`"`)

	return buffer.Bytes(), nil
}
