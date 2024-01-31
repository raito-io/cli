package implementations

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/smithy-go/ptr"

	"github.com/raito-io/cli/base/access_provider/sync_from_target"
	"github.com/raito-io/cli/base/access_provider_post_processor"
	logger_utils "github.com/raito-io/cli/base/util/logger"
	"github.com/raito-io/cli/base/wrappers"
)

const nameTagOverrideLockedReason = "This Snowflake role cannot be renamed because it has a name tag override attached to it"

var _ wrappers.AccessProviderPostProcessorI = (*AccessProviderPostProcessorBase)(nil)
var logger = logger_utils.InitializeLogger()

type AccessProviderPostProcessorBase struct {
	accessProviderWriter wrappers.AccessProviderPostProcessorHandler
	config               *access_provider_post_processor.AccessProviderPostProcessorConfig
}

func NewAccessProviderPostProcessorGeneral() *AccessProviderPostProcessorBase {
	return &AccessProviderPostProcessorBase{}
}

func (e *AccessProviderPostProcessorBase) Initialize(ctx context.Context, accessProviderWriter wrappers.AccessProviderPostProcessorHandler, config *access_provider_post_processor.AccessProviderPostProcessorConfig) error {
	e.accessProviderWriter = accessProviderWriter
	e.config = config

	logger.Info(fmt.Sprintf("Generic post processor initialized - %v", config))

	return nil
}

func (e *AccessProviderPostProcessorBase) PostProcess(ctx context.Context, accessProvider *sync_from_target.AccessProvider) (bool, error) {
	touched := false

	var overwrittenName *string = nil
	var overwrittenOwners []string = nil

	if len(accessProvider.Tags) > 0 {
		for _, tag := range accessProvider.Tags {
			if overwrittenName == nil && e.config.TagOverwriteKeyForName != "" && strings.EqualFold(tag.Key, e.config.TagOverwriteKeyForName) {
				overwrittenName = ptr.String(tag.Value)
				continue
			}

			if overwrittenOwners == nil && e.config.TagOverwriteKeyForOwners != "" && strings.EqualFold(tag.Key, e.config.TagOverwriteKeyForOwners) {
				overwrittenOwners = strings.Split(tag.Value, ",")
				continue
			}
		}
	}

	if overwrittenName != nil {
		logger.Debug(fmt.Sprintf("adjusting name for AP (externalId: %v) from %v to %v", accessProvider.ExternalId, accessProvider.Name, *overwrittenName))

		accessProvider.Name = *overwrittenName
		accessProvider.NameLocked = ptr.Bool(true)
		accessProvider.NameLockedReason = ptr.String(nameTagOverrideLockedReason)
		touched = true
	}

	if overwrittenOwners != nil {
		logger.Debug(fmt.Sprintf("adjusting owners for AP (externalId: %v) to %v", accessProvider.ExternalId, overwrittenOwners))

		if accessProvider.Owner == nil {
			accessProvider.Owner = &sync_from_target.OwnerInput{}
		}

		accessProvider.Owner.Users = overwrittenOwners
		touched = true
	}

	err := e.accessProviderWriter.AddAccessProviders(accessProvider)
	if err != nil {
		logger.Info(fmt.Sprintf("Error while saving AP to writer %q", accessProvider.ExternalId))
		return touched, err
	}

	return touched, nil
}

func (e *AccessProviderPostProcessorBase) Close(ctx context.Context) error {
	return nil
}
