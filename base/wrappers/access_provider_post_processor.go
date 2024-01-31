package wrappers

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/raito-io/cli/base/access_provider"
	"github.com/raito-io/cli/base/access_provider/sync_from_target"
	"github.com/raito-io/cli/base/access_provider_post_processor"
	"github.com/raito-io/cli/base/util/config"
)

//go:generate go run github.com/vektra/mockery/v2 --name=AccessProviderPostProcessorHandler --with-expecter
type AccessProviderPostProcessorHandler interface {
	AddAccessProviders(accessProviders ...*sync_from_target.AccessProvider) error
}

//go:generate go run github.com/vektra/mockery/v2 --name=AccessProviderPostProcessorI --with-expecter --inpackage
type AccessProviderPostProcessorI interface {
	// Initialize allows the plugin to do any preparation work
	Initialize(ctx context.Context, accessProviderWriter AccessProviderPostProcessorHandler, config *access_provider_post_processor.AccessProviderPostProcessorConfig) error

	// PostProcess method receives every access provider separately. The plugin can decide to skip or buffer things. All access providers must be written to the AccessProviderWriter.
	// First argument indicates if the access provider was touched during post processing,
	PostProcess(ctx context.Context, accessProvider *sync_from_target.AccessProvider) (bool, error)
}

type AccessProviderPostProcessorFactoryFn func(ctx context.Context, config *config.ConfigMap) (AccessProviderPostProcessorI, func(), error)

func AccessProviderPostProcessor(postProcessor AccessProviderPostProcessorI) *accessProviderPostProcessorFunction {
	return AccessProviderPostProcessorFactory(NewDummySyncFactoryFn(postProcessor))
}

func AccessProviderPostProcessorFactory(postProcessor AccessProviderPostProcessorFactoryFn) *accessProviderPostProcessorFunction {
	return &accessProviderPostProcessorFunction{
		postProcessor: NewSyncFactory(postProcessor),

		accessFileCreatorFactory:    sync_from_target.NewAccessProviderFileCreator,
		accessProviderParserFactory: sync_from_target.NewAccessProviderSyncFromTargetFileParser,
	}
}

type accessProviderPostProcessorFunction struct {
	access_provider_post_processor.AccessProviderPostProcessorVersionHandler

	postProcessor SyncFactory[AccessProviderPostProcessorI]

	accessFileCreatorFactory    func(config *access_provider.AccessSyncFromTarget) (sync_from_target.AccessProviderFileCreator, error)
	accessProviderParserFactory func(sourceFile string) (sync_from_target.AccessProviderSyncFromTargetFileParser, error)
}

func (p *accessProviderPostProcessorFunction) PostProcessFromTarget(ctx context.Context, config *access_provider_post_processor.AccessProviderPostProcessorConfig) (_ *access_provider_post_processor.AccessProviderPostProcessorResult, err error) {
	logger.Debug("Post processing access providers...")

	fileCreator, err := p.accessFileCreatorFactory(&access_provider.AccessSyncFromTarget{
		TargetFile: config.OutputFile,
	})
	if err != nil {
		return nil, err
	}
	defer fileCreator.Close()

	logger.Info(fmt.Sprintf("Post processor creating instance of plugin with params - %v", config.ConfigMap))

	postProcessor, err := p.postProcessor.Create(ctx, config.ConfigMap)
	if err != nil {
		return nil, fmt.Errorf("create post processor: %w", err)
	}

	err = postProcessor.Initialize(ctx, fileCreator, config)
	if err != nil {
		return nil, err
	}

	accessProvidersRead := 0

	logger.Debug(fmt.Sprintf("Post processor reading file - %s", config.InputFile))

	accessProviderParser, err := p.accessProviderParserFactory(config.InputFile)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	logger.Debug("Post processor parsing APs")

	aps, err := accessProviderParser.ParseAccessProviders()
	if err != nil {
		return nil, err
	}

	accessProviderTouchedCount := 0

	for _, ap := range aps {
		logger.Debug(fmt.Sprintf("Start post processing access provider %q", ap.ExternalId))

		postProcessed, err2 := postProcessor.PostProcess(ctx, ap)
		if err2 != nil {
			return nil, fmt.Errorf("unable to post process access provider (%d): %s", accessProvidersRead, err2.Error())
		}

		if postProcessed {
			accessProviderTouchedCount++
		}

		accessProvidersRead++
	}

	if fileCreator.GetAccessProviderCount() != accessProvidersRead {
		return nil, fmt.Errorf("post processor wrote %d access providers, while expecting %d", fileCreator.GetAccessProviderCount(), accessProvidersRead)
	}

	return &access_provider_post_processor.AccessProviderPostProcessorResult{
		AccessProviderTouchedCount: int32(accessProviderTouchedCount),
	}, nil
}

func (p *accessProviderPostProcessorFunction) Close() {
	p.postProcessor.Close()
}
