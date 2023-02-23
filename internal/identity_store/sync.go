package identity_store

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ispc "github.com/raito-io/cli/base/identity_store"
	baseconfig "github.com/raito-io/cli/base/util/config"
	"github.com/raito-io/cli/internal/file"
	"github.com/raito-io/cli/internal/job"
	"github.com/raito-io/cli/internal/plugin"
	"github.com/raito-io/cli/internal/target"
	"github.com/raito-io/cli/internal/version"
)

type IdentityStoreSync struct {
	TargetConfig *target.BaseTargetConfig
	JobId        string
}

type IdentityStoreImportResult struct {
	UsersAdded    int `json:"usersAdded"`
	UsersUpdated  int `json:"usersUpdated"`
	UsersRemoved  int `json:"usersRemoved"`
	GroupsAdded   int `json:"groupsAdded"`
	GroupsUpdated int `json:"groupsUpdated"`
	GroupsRemoved int `json:"groupsRemoved"`

	Warnings []string `json:"warnings"`
}

func (s *IdentityStoreSync) IsClientValid(ctx context.Context, c plugin.PluginClient) (bool, error) {
	iss, err := c.GetIdentityStoreSyncer()
	if err != nil {
		return false, err
	}

	return version.IsValidToSync(ctx, iss, ispc.MinimalCliVersion)
}

func (s *IdentityStoreSync) GetParts() []job.TaskPart {
	return []job.TaskPart{s}
}

func (s *IdentityStoreSync) StartSyncAndQueueTaskPart(client plugin.PluginClient, statusUpdater job.TaskEventUpdater) (job.JobStatus, string, error) {
	cn := strings.Replace(s.TargetConfig.ConnectorName, "/", "-", -1)

	userFile, err := filepath.Abs(file.CreateUniqueFileName(cn+"-is-user", "json"))
	if err != nil {
		return job.Failed, "", err
	}

	groupFile, err := filepath.Abs(file.CreateUniqueFileName(cn+"-is-group", "json"))
	if err != nil {
		return job.Failed, "", err
	}

	s.TargetConfig.TargetLogger.Debug(fmt.Sprintf("Using %q as user target file", userFile))
	s.TargetConfig.TargetLogger.Debug(fmt.Sprintf("Using %q as groups target file", groupFile))

	if s.TargetConfig.DeleteTempFiles {
		defer os.RemoveAll(userFile)
		defer os.RemoveAll(groupFile)
	}

	syncerConfig := ispc.IdentityStoreSyncConfig{
		ConfigMap: &baseconfig.ConfigMap{Parameters: s.TargetConfig.Parameters},
		UserFile:  userFile,
		GroupFile: groupFile,
	}

	iss, err := client.GetIdentityStoreSyncer()
	if err != nil {
		return job.Failed, "", err
	}

	s.TargetConfig.TargetLogger.Info("Fetching identity store metadata")

	md, err := iss.GetIdentityStoreMetaData(context.Background())
	if err != nil {
		return job.Failed, "", err
	}

	s.TargetConfig.TargetLogger.Info("Updating identity store metadata")
	err = SetMetaData(*s.TargetConfig, md)

	if err != nil {
		return job.Failed, "", err
	}

	s.TargetConfig.TargetLogger.Info("Gathering users and groups")

	result, err := iss.SyncIdentityStore(context.Background(), &syncerConfig)
	if err != nil {
		return job.Failed, "", err
	}

	if result.Error != nil {
		return job.Failed, "", result.Error
	}

	importerConfig := IdentityStoreImportConfig{
		BaseTargetConfig: *s.TargetConfig,
		UserFile:         userFile,
		GroupFile:        groupFile,
		DeleteUntouched:  s.TargetConfig.DeleteUntouched,
		ReplaceGroups:    s.TargetConfig.ReplaceGroups,
		ReplaceTags:      s.TargetConfig.ReplaceTags,
	}
	isImporter := NewIdentityStoreImporter(&importerConfig, statusUpdater)

	s.TargetConfig.TargetLogger.Info("Importing users and groups into Raito")
	status, subtaskId, err := isImporter.TriggerImport(s.JobId)

	if err != nil {
		return job.Failed, "", err
	}

	if status == job.Queued {
		s.TargetConfig.TargetLogger.Info("Successfully queued import job. Wait until remote processing is done.")
	}

	s.TargetConfig.TargetLogger.Debug(fmt.Sprintf("Current status: %s", status.String()))

	return status, subtaskId, nil
}

func (s *IdentityStoreSync) ProcessResults(results interface{}) error {
	if isResult, ok := results.(*IdentityStoreImportResult); ok {
		if isResult != nil && len(isResult.Warnings) > 0 {
			s.TargetConfig.TargetLogger.Info(fmt.Sprintf("Synced users and groups with %d warnings (see below). Users: Added: %d - Removed: %d - Updated: %d | Groups: Added: %d - Removed: %d - Updated: %d",
				len(isResult.Warnings), isResult.UsersAdded, isResult.UsersRemoved, isResult.UsersUpdated, isResult.GroupsAdded, isResult.GroupsRemoved, isResult.GroupsUpdated))

			for _, warning := range isResult.Warnings {
				s.TargetConfig.TargetLogger.Warn(warning)
			}
		} else {
			s.TargetConfig.TargetLogger.Info(fmt.Sprintf("Successfully synced users and groups. Users: Added: %d - Removed: %d - Updated: %d | Groups: Added: %d - Removed: %d - Updated: %d", isResult.UsersAdded, isResult.UsersRemoved, isResult.UsersUpdated, isResult.GroupsAdded, isResult.GroupsRemoved, isResult.GroupsUpdated))
		}

		return nil
	}

	return fmt.Errorf("failed to load results")
}

func (s *IdentityStoreSync) GetResultObject() interface{} {
	return &IdentityStoreImportResult{}
}
