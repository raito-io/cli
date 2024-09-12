package identity_store

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/smithy-go/ptr"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/raito-io/cli/base/identity_store"
	"github.com/raito-io/cli/base/tag"
	"github.com/raito-io/cli/internal/util/jsonstream"
)

type PostProcessorConfig struct {
	TagKeyAndValueForUserIsMachine string
}

type PostProcessorResult struct {
	UsersTouchedCount int
}

type PostProcessor struct {
	identityStoreFileCreator func(config *identity_store.IdentityStoreSyncConfig) (identity_store.IdentityStoreFileCreator, error)

	config *PostProcessorConfig
}

func NewPostProcessor(config *PostProcessorConfig) PostProcessor {
	return PostProcessor{
		identityStoreFileCreator: identity_store.NewIdentityStoreFileCreator,
		config:                   config,
	}
}

func (p *PostProcessor) NeedsUserPostProcessing() bool {
	tagKeyUserIsMachine, tagValueUserIsMachine := p.splitTagKeyAndValueForUserIsMachine()
	return !strings.EqualFold(tagKeyUserIsMachine, "") && !strings.EqualFold(tagValueUserIsMachine, "")
}

func (p *PostProcessor) splitTagKeyAndValueForUserIsMachine() (string, string) {
	parts := strings.Split(p.config.TagKeyAndValueForUserIsMachine, ":")
	if len(parts) != 2 {
		return "", ""
	}

	return parts[0], parts[1]
}

func (p *PostProcessor) errorWrapper(err error) error {
	return status.Error(codes.Internal, err.Error())
}

func (p *PostProcessor) PostProcessUsers(logger hclog.Logger, usersInputFilePath string, usersOutputFile string) (*PostProcessorResult, error) {
	outputWriter, err := p.identityStoreFileCreator(&identity_store.IdentityStoreSyncConfig{
		UserFile: usersOutputFile,
	})

	if err != nil {
		return nil, p.errorWrapper(err)
	}

	logger.Debug(fmt.Sprintf("Post processor streaming users from file %s", usersInputFilePath))

	usersInputFile, err := os.Open(usersInputFilePath)
	if err != nil {
		return nil, p.errorWrapper(fmt.Errorf("unable to open input file %q: %s", usersInputFilePath, err.Error()))
	}
	defer usersInputFile.Close()

	usersRead := 0
	usersTouched := 0

	tagKeyUserIsMachine, tagValueUserIsMachine := p.splitTagKeyAndValueForUserIsMachine()

	decoder := jsonstream.NewJsonArrayStream[identity_store.User](usersInputFile)
	for jsonStreamResult := range decoder.Stream() {
		logger.Trace(fmt.Sprintf("Start post processing user %d", usersRead))

		if jsonStreamResult.Err != nil {
			return nil, p.errorWrapper(fmt.Errorf("unable to parse user (%d): %s", usersRead, jsonStreamResult.Err.Error()))
		}

		user := jsonStreamResult.Result
		enriched := p.postProcessUser(logger, user, tagKeyUserIsMachine, tagValueUserIsMachine)

		err := outputWriter.AddUsers(user)
		if err != nil {
			return nil, p.errorWrapper(fmt.Errorf("unable to save user (%d) to writer: %s", usersRead, err.Error()))
		}

		if enriched {
			usersTouched++
		}

		usersRead++
	}

	if outputWriter.GetUserCount() != usersRead {
		return nil, p.errorWrapper(fmt.Errorf("post processor wrote %d users, while expecting %d", outputWriter.GetUserCount(), usersRead))
	}

	return &PostProcessorResult{
		UsersTouchedCount: usersTouched,
	}, nil
}

func (p *PostProcessor) postProcessUser(logger hclog.Logger, user *identity_store.User, tagKeyUserIsMachine string, tagValueUserIsMachine string) bool {
	enriched := false

	if len(user.Tags) > 0 {
		for _, tag := range user.Tags {
			if p.matchedWithTagKey(tagKeyUserIsMachine, tag) && strings.EqualFold(tag.Value, tagValueUserIsMachine) {
				user.IsMachine = ptr.Bool(true)
				enriched = true
			}
		}
	}

	if enriched {
		logger.Debug(fmt.Sprintf("Enriched user %q", user.Name))
	}

	return enriched
}

func (p *PostProcessor) matchedWithTagKey(overwriteKey string, tag *tag.Tag) bool {
	return tag != nil && overwriteKey != "" && strings.EqualFold(tag.Key, overwriteKey) && tag.Value != ""
}
