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
	TagKeyForUserIsMachine   string
	TagValueForUserIsMachine string
	TargetLogger             hclog.Logger
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
	return p.config.TagKeyForUserIsMachine != "" && p.config.TagValueForUserIsMachine != ""
}

func (p *PostProcessor) errorWrapper(err error) error {
	return status.Error(codes.Internal, err.Error())
}

func (p *PostProcessor) PostProcessUsers(usersInputFilePath string, usersOutputFile string) (*PostProcessorResult, error) {
	outputWriter, err := p.identityStoreFileCreator(&identity_store.IdentityStoreSyncConfig{
		UserFile: usersOutputFile,
	})

	if err != nil {
		return nil, p.errorWrapper(err)
	}

	p.config.TargetLogger.Debug(fmt.Sprintf("Post processor streaming users from file %s", usersInputFilePath))

	usersInputFile, err := os.Open(usersInputFilePath)
	if err != nil {
		return nil, p.errorWrapper(fmt.Errorf("unable to open input file %q: %s", usersInputFilePath, err.Error()))
	}
	defer usersInputFile.Close()

	usersRead := 0
	usersTouched := 0

	decoder := jsonstream.NewJsonArrayStream[identity_store.User](usersInputFile)
	for jsonStreamResult := range decoder.Stream() {
		p.config.TargetLogger.Trace(fmt.Sprintf("Start post processing user %d", usersRead))

		if jsonStreamResult.Err != nil {
			return nil, p.errorWrapper(fmt.Errorf("unable to parse user (%d): %s", usersRead, jsonStreamResult.Err.Error()))
		}

		user := jsonStreamResult.Result
		enriched := p.postProcessUser(user)

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

func (p *PostProcessor) postProcessUser(user *identity_store.User) bool {
	enriched := false

	if len(user.Tags) > 0 {
		for _, tag := range user.Tags {
			if p.matchedWithTagKey(p.config.TagKeyForUserIsMachine, tag) && strings.EqualFold(tag.Value, p.config.TagValueForUserIsMachine) {
				user.IsMachine = ptr.Bool(true)
				enriched = true
			}
		}
	}

	if enriched {
		p.config.TargetLogger.Debug(fmt.Sprintf("Enriched user %q", user.Name))
	}

	return enriched
}

func (p *PostProcessor) matchedWithTagKey(overwriteKey string, tag *tag.Tag) bool {
	return tag != nil && overwriteKey != "" && strings.EqualFold(tag.Key, overwriteKey) && tag.Value != ""
}
