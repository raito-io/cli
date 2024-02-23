package resource_provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/viper"

	"github.com/raito-io/cli/base/resource_provider"
	baseconfig "github.com/raito-io/cli/base/util/config"
	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/job"
	"github.com/raito-io/cli/internal/plugin"
	"github.com/raito-io/cli/internal/target/types"
	"github.com/raito-io/cli/internal/version_management"
)

type ResourceSync struct {
	TargetConfig *types.BaseTargetConfig
	JobId        string

	result *job.TaskResult
}

func (s *ResourceSync) IsClientValid(ctx context.Context, c plugin.PluginClient) (bool, error) {
	rsc, err := c.GetResourceProvider()
	if err != nil {
		return false, err
	}

	return version_management.IsValidToSync(ctx, rsc, resource_provider.MinimalCliVersion)
}

func (s *ResourceSync) GetParts() []job.TaskPart {
	return []job.TaskPart{s}
}

func (s *ResourceSync) StartSyncAndQueueTaskPart(ctx context.Context, client plugin.PluginClient, statusUpdater job.TaskEventUpdater) (job.JobStatus, string, error) {
	var urlOverridePtr *string

	urlOverride := viper.GetString(constants.URLOverrideFlag)
	if urlOverride != "" {
		urlOverridePtr = &urlOverride
	}

	syncInput := resource_provider.UpdateResourceInput{
		ConfigMap:       &baseconfig.ConfigMap{Parameters: s.TargetConfig.Parameters},
		Domain:          s.TargetConfig.Domain,
		DataSourceId:    s.TargetConfig.DataSourceId,
		IdentityStoreId: s.TargetConfig.IdentityStoreId,
		UrlOverride:     urlOverridePtr,
		Credentials: &resource_provider.ApiCredentials{
			Username: s.TargetConfig.ApiUser,
			Password: s.TargetConfig.ApiSecret,
		},
	}

	rsc, err := client.GetResourceProvider()
	if err != nil {
		return job.Failed, "", fmt.Errorf("get resource provider client: %w", err)
	}

	statusUpdater.SetStatusToDataProcessing(ctx)

	result, err := rsc.UpdateResources(ctx, &syncInput)
	if err != nil {
		return job.Failed, "", fmt.Errorf("update resources: %w", err)
	}

	s.result = &job.TaskResult{
		ObjectType: "Resource objects",
		Added:      int(result.AddedObjects),
		Updated:    int(result.UpdatedObjects),
		Removed:    int(result.DeletedObjects),
		Failed:     int(result.Failures),
	}

	s.TargetConfig.TargetLogger.Info(fmt.Sprintf("Successfully synced resource objects. Added: %d - Updated: %d - Removed: %d - Failures: %d", result.AddedObjects, result.UpdatedObjects, result.DeletedObjects, result.Failures))

	return job.Completed, "", nil
}

func (s *ResourceSync) ProcessResults(_ interface{}) error {
	return errors.New("unexpected result processing")
}

func (s *ResourceSync) GetResultObject() interface{} {
	return nil
}

func (s *ResourceSync) GetTaskResults() []job.TaskResult {
	if s.result == nil {
		return nil
	}

	return []job.TaskResult{*s.result}
}
