package identity_store

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ispc "github.com/raito-io/cli/common/api/identity_store"
	baseconfig "github.com/raito-io/cli/common/util/config"
	"github.com/raito-io/cli/internal/file"
	"github.com/raito-io/cli/internal/job"
	"github.com/raito-io/cli/internal/plugin"
	"github.com/raito-io/cli/internal/target"
)

type IdentityStoreSync struct {
	TargetConfig  *target.BaseTargetConfig
	JobId         string
	StatusUpdater func(status job.JobStatus)
}

type IdentityStoreImportResult struct {
	UsersAdded    int `json:"usersAdded"`
	UsersUpdated  int `json:"usersUpdated"`
	UsersRemoved  int `json:"usersRemoved"`
	GroupsAdded   int `json:"groupsAdded"`
	GroupsUpdated int `json:"groupsUpdated"`
	GroupsRemoved int `json:"groupsRemoved"`
}

func (s *IdentityStoreSync) StartSyncAndQueueJob(client plugin.PluginClient) (job.JobStatus, error) {
	s.StatusUpdater(job.Started)

	cn := strings.Replace(s.TargetConfig.ConnectorName, "/", "-", -1)

	userFile, err := filepath.Abs(file.CreateUniqueFileName(cn+"-is-user", "json"))
	if err != nil {
		return job.Failed, err
	}

	groupFile, err := filepath.Abs(file.CreateUniqueFileName(cn+"-is-group", "json"))
	if err != nil {
		return job.Failed, err
	}

	s.TargetConfig.Logger.Debug(fmt.Sprintf("Using %q as user target file", userFile))
	s.TargetConfig.Logger.Debug(fmt.Sprintf("Using %q as groups target file", groupFile))

	if s.TargetConfig.DeleteTempFiles {
		defer os.RemoveAll(userFile)
		defer os.RemoveAll(groupFile)
	}

	syncerConfig := ispc.IdentityStoreSyncConfig{
		ConfigMap: baseconfig.ConfigMap{Parameters: s.TargetConfig.Parameters},
		UserFile:  userFile,
		GroupFile: groupFile,
	}

	iss, err := client.GetIdentityStoreSyncer()
	if err != nil {
		return job.Failed, err
	}

	s.TargetConfig.Logger.Info("Gathering users and groups")
	result := iss.SyncIdentityStore(&syncerConfig)

	if result.Error != nil {
		return job.Failed, *(result.Error)
	}

	importerConfig := IdentityStoreImportConfig{
		BaseTargetConfig: *s.TargetConfig,
		UserFile:         userFile,
		GroupFile:        groupFile,
		DeleteUntouched:  s.TargetConfig.DeleteUntouched,
		ReplaceGroups:    s.TargetConfig.ReplaceGroups,
		ReplaceTags:      s.TargetConfig.ReplaceTags,
	}
	isImporter := NewIdentityStoreImporter(&importerConfig, s.StatusUpdater)

	s.TargetConfig.Logger.Info("Importing users and groups into Raito")
	status, err := isImporter.TriggerImport(s.JobId)

	if err != nil {
		return job.Failed, err
	}

	s.TargetConfig.Logger.Info("Successfully queued import job. Wait until remote processing is done.")
	s.TargetConfig.Logger.Debug(fmt.Sprintf("Current status: %s", status.String()))

	return status, nil
}

func (s *IdentityStoreSync) ProcessResults(results interface{}) error {
	if isResult, ok := results.(*IdentityStoreImportResult); ok {
		s.TargetConfig.Logger.Info(fmt.Sprintf("Successfully synced users and groups. Users: Added: %d - Removed: %d - Updated: %d | Groups: Added: %d - Removed: %d - Updated: %d", isResult.UsersAdded, isResult.UsersRemoved, isResult.UsersUpdated, isResult.GroupsAdded, isResult.GroupsRemoved, isResult.GroupsUpdated))
		return nil
	}

	return fmt.Errorf("failed to load results")
}

func (s *IdentityStoreSync) GetResultObject() interface{} {
	return &IdentityStoreImportResult{}
}