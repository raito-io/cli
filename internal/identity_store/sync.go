package identity_store

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	ispc "github.com/raito-io/cli/base/identity_store"
	baseconfig "github.com/raito-io/cli/base/util/config"
	error1 "github.com/raito-io/cli/base/util/error"
	"github.com/raito-io/cli/internal/job"
	"github.com/raito-io/cli/internal/plugin"
	"github.com/raito-io/cli/internal/target/types"
	"github.com/raito-io/cli/internal/util/file"
	"github.com/raito-io/cli/internal/util/tag"
	"github.com/raito-io/cli/internal/version_management"
)

type IdentityStoreSync struct {
	TargetConfig *types.BaseTargetConfig
	JobId        string

	result []job.TaskResult
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

	return version_management.IsValidToSync(ctx, iss, ispc.MinimalCliVersion)
}

func (s *IdentityStoreSync) GetParts() []job.TaskPart {
	return []job.TaskPart{s}
}

func (s *IdentityStoreSync) StartSyncAndQueueTaskPart(ctx context.Context, client plugin.PluginClient, statusUpdater job.TaskEventUpdater) (job.JobStatus, string, error) {
	userFile, err := filepath.Abs(file.CreateUniqueFileNameForTarget(s.TargetConfig.Name, "fromTarget-users", "json"))
	if err != nil {
		return job.Failed, "", err
	}

	defer s.TargetConfig.HandleTempFile(userFile, false)

	groupFile, err := filepath.Abs(file.CreateUniqueFileNameForTarget(s.TargetConfig.Name, "fromTarget-groups", "json"))
	if err != nil {
		return job.Failed, "", err
	}

	defer s.TargetConfig.HandleTempFile(groupFile, false)

	s.TargetConfig.TargetLogger.Debug(fmt.Sprintf("Using %q as user target file", userFile))
	s.TargetConfig.TargetLogger.Debug(fmt.Sprintf("Using %q as groups target file", groupFile))

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

	md, err := iss.GetIdentityStoreMetaData(ctx, syncerConfig.ConfigMap)
	if err != nil {
		return job.Failed, "", err
	}

	s.TargetConfig.TargetLogger.Info("Updating identity store metadata")
	err = SetMetaData(*s.TargetConfig, md)

	if err != nil {
		return job.Failed, "", err
	}

	s.TargetConfig.TargetLogger.Info("Gathering users and groups")

	result, err := iss.SyncIdentityStore(ctx, &syncerConfig)
	if err != nil {
		return job.Failed, "", err
	}

	if result.Error != nil { //nolint:staticcheck
		return job.Failed, "", s.mapErrorResult(result.Error) //nolint:staticcheck
	}

	postProcessor := NewPostProcessor(&PostProcessorConfig{
		TagKeyForUserIsMachine:   s.TargetConfig.TagKeyForUserIsMachine,
		TagValueForUserIsMachine: s.TargetConfig.TagValueForUserIsMachine,
		TargetLogger:             s.TargetConfig.TargetLogger,
	})

	toProcessUserFile := userFile
	if postProcessor.NeedsUserPostProcessing() {
		toProcessUserFile, err = s.postProcessUsers(postProcessor, userFile)
		if err != nil {
			return job.Failed, "", err
		}
	}

	// Fetching the tagSource from the plugin
	tagSourcesScope, err := tag.FetchTagSourceFromPlugin(ctx, client, nil)
	if err != nil {
		return job.Failed, "", err
	}

	importerConfig := IdentityStoreImportConfig{
		BaseTargetConfig: *s.TargetConfig,
		UserFile:         toProcessUserFile,
		GroupFile:        groupFile,
		DeleteUntouched:  s.TargetConfig.DeleteUntouched,
		ReplaceGroups:    s.TargetConfig.ReplaceGroups,
		TagSourcesScope:  tagSourcesScope,
	}
	isImporter := NewIdentityStoreImporter(&importerConfig, statusUpdater)

	s.TargetConfig.TargetLogger.Info("Importing users and groups into Raito")
	status, subtaskId, err := isImporter.TriggerImport(ctx, s.JobId)

	if err != nil {
		return job.Failed, "", err
	}

	if status == job.Queued {
		s.TargetConfig.TargetLogger.Info("Successfully queued import job. Wait until remote processing is done.")
	}

	s.TargetConfig.TargetLogger.Debug(fmt.Sprintf("Current status: %s", status.String()))

	return status, subtaskId, nil
}

func (s *IdentityStoreSync) postProcessUsers(postProcessor PostProcessor, toProcessFile string) (string, error) {
	postProcessedFile := toProcessFile
	fileSuffix := "-post-processed"

	// Generate a unique file name for the post processing
	if strings.Contains(postProcessedFile, fileSuffix) {
		postProcessedFile = postProcessedFile[0:strings.LastIndex(postProcessedFile, fileSuffix)] + fileSuffix + ".json"
	} else {
		postProcessedFile = postProcessedFile[0:strings.LastIndex(postProcessedFile, ".json")] + fileSuffix + ".json"
	}

	res, err := postProcessor.PostProcessUsers(toProcessFile, postProcessedFile)
	if err != nil {
		return toProcessFile, err
	}

	if res.UsersTouchedCount > 0 {
		s.TargetConfig.TargetLogger.Info(fmt.Sprintf("Successfully updated %d users in post-processing step", res.UsersTouchedCount))
	}

	return postProcessedFile, nil
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

		s.result = append(s.result, job.TaskResult{
			ObjectType: "users",
			Added:      isResult.UsersAdded,
			Removed:    isResult.UsersRemoved,
			Updated:    isResult.UsersUpdated,
		}, job.TaskResult{
			ObjectType: "groups",
			Added:      isResult.GroupsAdded,
			Removed:    isResult.GroupsRemoved,
			Updated:    isResult.GroupsUpdated,
		})

		return nil
	}

	return fmt.Errorf("failed to load results")
}

func (s *IdentityStoreSync) GetResultObject() interface{} {
	return &IdentityStoreImportResult{}
}

func (s *IdentityStoreSync) GetTaskResults() []job.TaskResult {
	return s.result
}

func (s *IdentityStoreSync) mapErrorResult(result *error1.ErrorResult) error {
	if result == nil {
		return nil
	}

	return errors.New(result.ErrorMessage)
}
