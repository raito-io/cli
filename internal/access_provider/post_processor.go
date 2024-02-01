package access_provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/smithy-go/ptr"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	baseAp "github.com/raito-io/cli/base/access_provider"
	"github.com/raito-io/cli/base/access_provider/sync_from_target"
	"github.com/raito-io/cli/base/tag"
	"github.com/raito-io/cli/internal/access_provider/post_processing"
	"github.com/raito-io/cli/internal/util/array"
)

const nameTagOverrideLockedReason = "This Snowflake role cannot be renamed because it has a name tag override attached to it"

type PostProcessorConfig struct {
	TagOverwriteKeyForName   string
	TagOverwriteKeyForOwners string
	TargetLogger             hclog.Logger
}
type PostProcessorResult struct {
	AccessProviderTouchedCount int
}

type PostProcessorOutputFileWriter interface {
	AddAccessProviders(accessProviders ...*sync_from_target.AccessProvider) error
}

type PostProcessor struct {
	accessFileCreatorFactory    func(config *baseAp.AccessSyncFromTarget) (sync_from_target.AccessProviderFileCreator, error)
	accessProviderParserFactory func(sourceFile string) (post_processing.PostProcessorSourceFileParser, error)

	config *PostProcessorConfig
}

func NewAccessProviderPostProcessorGeneral(config *PostProcessorConfig) PostProcessor {
	return PostProcessor{
		accessProviderParserFactory: post_processing.NewPostProcessorSourceFileParser,
		accessFileCreatorFactory:    sync_from_target.NewAccessProviderFileCreator,
		config:                      config,
	}
}

func (p *PostProcessor) NeedsPostProcessing() bool {
	return p.config.TagOverwriteKeyForName != "" || p.config.TagOverwriteKeyForOwners != ""
}

func (p *PostProcessor) PostProcess(inputFilePath string, outputFile string) (*PostProcessorResult, error) {
	accessProviderParser, err := p.accessProviderParserFactory(inputFilePath)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	outputWriter, err := p.accessFileCreatorFactory(&baseAp.AccessSyncFromTarget{
		TargetFile: outputFile,
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	p.config.TargetLogger.Debug("Post processor parsing APs")

	aps, err := accessProviderParser.ParseAccessProviders()
	if err != nil {
		return nil, err
	}

	accessProvidersRead := 0
	accessProviderTouchedCount := 0

	for _, ap := range aps {
		p.config.TargetLogger.Debug(fmt.Sprintf("Start post processing access provider %q", ap.ExternalId))

		touched, err := p.postProcessAp(ap, outputWriter)
		if err != nil {
			return nil, fmt.Errorf("unable to post process access provider (%d): %s", accessProvidersRead, err.Error())
		}

		if touched {
			accessProviderTouchedCount++
		}

		accessProvidersRead++
	}

	if outputWriter.GetAccessProviderCount() != accessProvidersRead {
		return nil, fmt.Errorf("post processor wrote %d access providers, while expecting %d", outputWriter.GetAccessProviderCount(), accessProvidersRead)
	}

	return &PostProcessorResult{
		AccessProviderTouchedCount: accessProviderTouchedCount,
	}, nil
}

func (p *PostProcessor) postProcessAp(accessProvider *sync_from_target.AccessProvider, outputWriter sync_from_target.AccessProviderFileCreator) (bool, error) {
	overwrittenOn := make(map[string]bool)

	if len(accessProvider.Tags) > 0 {
		for _, tag := range accessProvider.Tags {
			if !overwrittenOn["name"] && p.matchedWithTagKey(p.config.TagOverwriteKeyForName, tag) {
				touched := p.overwriteName(accessProvider, tag)
				if touched {
					overwrittenOn["name"] = true
					continue
				}
			}

			if !overwrittenOn["owners"] && p.matchedWithTagKey(p.config.TagOverwriteKeyForOwners, tag) {
				touched := p.overwriteOwners(accessProvider, tag)
				if touched {
					overwrittenOn["owners"] = true
					continue
				}
			}
		}
	}

	touched := len(array.Keys(overwrittenOn)) > 0

	err := outputWriter.AddAccessProviders(accessProvider)
	if err != nil {
		p.config.TargetLogger.Info(fmt.Sprintf("Error while saving AP to writer %q", accessProvider.ExternalId))
		return touched, err
	}

	return touched, nil
}

func (p *PostProcessor) overwriteName(accessProvider *sync_from_target.AccessProvider, tag *tag.Tag) bool {
	if tag.Value != "" {
		p.config.TargetLogger.Debug(fmt.Sprintf("adjusting name for AP (externalId: %v) from %v to %v", accessProvider.ExternalId, accessProvider.Name, tag.Value))

		accessProvider.Name = tag.Value
		accessProvider.NameLocked = ptr.Bool(true)
		accessProvider.NameLockedReason = ptr.String(nameTagOverrideLockedReason)

		return true
	}

	return false
}
func (p *PostProcessor) overwriteOwners(accessProvider *sync_from_target.AccessProvider, tag *tag.Tag) bool {
	if tag.Value != "" {
		overwrittenOwners := strings.Split(tag.Value, ",")

		p.config.TargetLogger.Debug(fmt.Sprintf("adjusting owners for AP (externalId: %v) to %v", accessProvider.ExternalId, overwrittenOwners))

		if accessProvider.Owner == nil {
			accessProvider.Owner = &sync_from_target.OwnerInput{}
		}

		accessProvider.Owner.Users = overwrittenOwners

		return true
	}

	return false
}

func (p *PostProcessor) matchedWithTagKey(overwriteKey string, tag *tag.Tag) bool {
	return tag != nil && overwriteKey != "" && strings.EqualFold(tag.Key, overwriteKey) && tag.Value != ""
}

func (p *PostProcessor) Close(ctx context.Context) error {
	return nil
}
