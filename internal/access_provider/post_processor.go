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
	"github.com/raito-io/cli/base/constants"
	"github.com/raito-io/cli/base/tag"
	"github.com/raito-io/cli/internal/access_provider/post_processing"
	"github.com/raito-io/cli/internal/util/stringops"
)

const nameTagOverrideLockedReason = "This Snowflake role can't be renamed because it has a name tag override attached to it"
const ownersTagOverrideLockedReason = "This Snowflake role can't update its owners as it has an owners tag override attached to it"

type PostProcessorConfig struct {
	TagOverwriteKeyForName   string
	TagOverwriteKeyForOwners string
	TargetLogger             hclog.Logger
}
type PostProcessorResult struct {
	AccessProviderTouchedCount int
}

type PostProcessor struct {
	accessFileCreatorFactory    func(config *baseAp.AccessSyncFromTarget) (sync_from_target.AccessProviderFileCreator, error)
	accessProviderParserFactory func(sourceFile string) (post_processing.PostProcessorSourceFileParser, error)

	config *PostProcessorConfig
}

func NewPostProcessor(config *PostProcessorConfig) PostProcessor {
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

	defer outputWriter.Close()

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
	touched := p.processOverwriteName(accessProvider)
	touched = p.processOverwriteOwners(accessProvider) || touched

	err := outputWriter.AddAccessProviders(accessProvider)
	if err != nil {
		p.config.TargetLogger.Info(fmt.Sprintf("Error while saving AP to writer %q", accessProvider.ExternalId))
		return touched, err
	}

	return touched, nil
}

func (p *PostProcessor) processOverwriteName(accessProvider *sync_from_target.AccessProvider) bool {
	if p.config.TagOverwriteKeyForName == "" {
		return false
	}

	for _, t := range accessProvider.Tags {
		if p.matchedWithTagKey(p.config.TagOverwriteKeyForName, t) {
			return p.overwriteName(accessProvider, t)
		}
	}

	return false
}

func (p *PostProcessor) processOverwriteOwners(accessProvider *sync_from_target.AccessProvider) (touched bool) {
	touched = false

	if p.config.TagOverwriteKeyForOwners == "" {
		return touched
	}

	var raitoOwnerTag *tag.Tag

	for _, t := range accessProvider.Tags {
		if !strings.EqualFold(t.Key, constants.RaitoOwnerTagKey) {
			continue
		}

		value := stringops.TrimSpaceInCommaSeparatedList(t.Value)

		if value == "" {
			continue
		}

		raitoOwnerTag = t
		raitoOwnerTag.Value = value

		touched = true

		break
	}

	for _, t := range accessProvider.Tags {
		if p.matchedWithTagKey(p.config.TagOverwriteKeyForOwners, t) {
			value := stringops.TrimSpaceInCommaSeparatedList(t.Value)
			if value == "" {
				continue
			}

			if raitoOwnerTag != nil {
				raitoOwnerTag.Value = raitoOwnerTag.Value + "," + value
			} else {
				raitoOwnerTag = &tag.Tag{
					Key:    constants.RaitoOwnerTagKey,
					Value:  value,
					Source: t.Source,
				}

				accessProvider.Tags = append(accessProvider.Tags, raitoOwnerTag)
			}

			touched = true
		}
	}

	if raitoOwnerTag != nil {
		accessProvider.OwnersLocked = ptr.Bool(true)
		accessProvider.OwnersLockedReason = ptr.String(ownersTagOverrideLockedReason)
	}

	return touched
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

func (p *PostProcessor) matchedWithTagKey(overwriteKey string, tag *tag.Tag) bool {
	return tag != nil && overwriteKey != "" && strings.EqualFold(tag.Key, overwriteKey) && tag.Value != ""
}

func (p *PostProcessor) Close(_ context.Context) error {
	return nil
}
