package wrappers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/raito-io/cli/base/access_provider"
	"github.com/raito-io/cli/base/access_provider/sync_from_target"
	"github.com/raito-io/cli/base/access_provider/sync_to_target"
	"github.com/raito-io/cli/base/access_provider/sync_to_target/naming_hint"
	"github.com/raito-io/cli/base/util/config"
	e "github.com/raito-io/cli/base/util/error"
)

//go:generate go run github.com/vektra/mockery/v2 --name=AccessProviderHandler --with-expecter
type AccessProviderHandler interface {
	AddAccessProviders(dataAccessList ...*sync_from_target.AccessProvider) error
}

//go:generate go run github.com/vektra/mockery/v2 --name=AccessProviderFeedbackHandler --with-expecter
type AccessProviderFeedbackHandler interface {
	AddAccessProviderFeedback(accessProviderId string, accessFeedback ...sync_to_target.AccessSyncFeedbackInformation) error
}

//go:generate go run github.com/vektra/mockery/v2 --name=AccessProviderSyncer --with-expecter --inpackage
type AccessProviderSyncer interface {
	SyncAccessProviderFromTarget(ctx context.Context, accessProviderHandler AccessProviderHandler, configMap *config.ConfigMap) error
	SyncAccessProviderToTarget(ctx context.Context, rolesToRemove []string, accesses map[string]sync_to_target.EnrichedAccess, feedbackHandler AccessProviderFeedbackHandler, configMap *config.ConfigMap) error
	SyncAccessAsCodeToTarget(ctx context.Context, accesses map[string]sync_to_target.EnrichedAccess, configMap *config.ConfigMap) error
}

func DataAccessSync(syncer AccessProviderSyncer, namingConstraints naming_hint.NamingConstraints) *dataAccessSyncFunction {
	return &dataAccessSyncFunction{
		syncer:                           syncer,
		namingConstraints:                namingConstraints,
		accessFileCreatorFactory:         sync_from_target.NewAccessProviderFileCreator,
		accessFeedbackFileCreatorFactory: sync_to_target.NewFeedbackFileCreator,
		accessProviderParserFactory:      sync_to_target.NewAccessProviderFileParser,
	}
}

type dataAccessSyncFunction struct {
	syncer                           AccessProviderSyncer
	namingConstraints                naming_hint.NamingConstraints
	accessFileCreatorFactory         func(config *access_provider.AccessSyncFromTarget) (sync_from_target.AccessProviderFileCreator, error)
	accessFeedbackFileCreatorFactory func(config *access_provider.AccessSyncToTarget) (sync_to_target.SyncFeedbackFileCreator, error)
	accessProviderParserFactory      func(config *access_provider.AccessSyncToTarget) (sync_to_target.AccessProviderImportFileParser, error)
}

func (s *dataAccessSyncFunction) SyncFromTarget(config *access_provider.AccessSyncFromTarget) access_provider.AccessSyncResult {
	ctx := context.Background()

	logger.Info("Starting data access synchronisation from target")
	logger.Debug("Creating file for storing access providers")

	fileCreator, err := s.accessFileCreatorFactory(config)
	if err != nil {
		logger.Error(err.Error())

		return mapErrorToAccessSyncResult(err)
	}
	defer fileCreator.Close()

	sec, err := timedExecution(func() error {
		return s.syncer.SyncAccessProviderFromTarget(ctx, fileCreator, &config.ConfigMap)
	})

	if err != nil {
		logger.Error(err.Error())

		return mapErrorToAccessSyncResult(err)
	}

	logger.Info(fmt.Sprintf("Fetched %d access provider in %s", fileCreator.GetAccessProviderCount(), sec))

	return access_provider.AccessSyncResult{}
}

func (s *dataAccessSyncFunction) SyncToTarget(config *access_provider.AccessSyncToTarget) access_provider.AccessSyncResult {
	ctx := context.Background()

	logger.Info("Starting data access synchronisation to target")

	accessProviderParser, err := s.accessProviderParserFactory(config)
	if err != nil {
		logger.Error(err.Error())

		return mapErrorToAccessSyncResult(err)
	}

	dar, err := accessProviderParser.ParseAccessProviders()
	if err != nil {
		logger.Error(err.Error())

		return mapErrorToAccessSyncResult(err)
	}

	prefix := config.Prefix
	accessAsCode := config.Prefix != ""

	var sec time.Duration

	if accessAsCode {
		sec, err = s.syncToTargetAccessAsCode(ctx, config, dar, prefix)
	} else {
		sec, err = s.syncToTargetAccessProviders(ctx, config, dar)
	}

	if err != nil {
		logger.Error(err.Error())

		return mapErrorToAccessSyncResult(err)
	}

	logger.Info(fmt.Sprintf("Successfully synced access providers to target in %s", sec))

	return access_provider.AccessSyncResult{}
}

func (s *dataAccessSyncFunction) syncToTargetAccessAsCode(ctx context.Context, config *access_provider.AccessSyncToTarget, dar *sync_to_target.AccessProviderImport, prefix string) (time.Duration, error) {
	roleSeperator := string(s.namingConstraints.SplitCharacter())
	if !strings.HasSuffix(prefix, roleSeperator) {
		prefix += roleSeperator
	}

	uniqueRoleNameGenerator, err := naming_hint.NewUniqueNameGenerator(logger, prefix, &s.namingConstraints)
	if err != nil {
		return 0, err
	}

	logger.Info(fmt.Sprintf("Using prefix %q", prefix))

	apList := dar.AccessProviders
	apMap := make(map[string]sync_to_target.EnrichedAccess)

	for apIndex, ap := range apList {
		roleNames, err := uniqueRoleNameGenerator.GenerateOrdered(&apList[apIndex])
		if err != nil {
			return 0, err
		}

		for i, access := range ap.Access {
			roleName := roleNames[i]

			logger.Info(fmt.Sprintf("Generated rolename %q", roleName))
			apMap[roleName] = sync_to_target.EnrichedAccess{Access: access, AccessProvider: &apList[apIndex]}
		}
	}

	return timedExecution(func() error {
		return s.syncer.SyncAccessAsCodeToTarget(ctx, apMap, &config.ConfigMap)
	})
}

func (s *dataAccessSyncFunction) syncToTargetAccessProviders(ctx context.Context, config *access_provider.AccessSyncToTarget, dar *sync_to_target.AccessProviderImport) (time.Duration, error) {
	feedbackFile, err := s.accessFeedbackFileCreatorFactory(config)
	if err != nil {
		return 0, err
	}
	defer feedbackFile.Close()

	uniqueRoleNameGenerator, err := naming_hint.NewUniqueNameGenerator(logger, "", &s.namingConstraints)
	if err != nil {
		return 0, err
	}

	apList := dar.AccessProviders

	apMap := make(map[string]sync_to_target.EnrichedAccess)
	rolesToRemove := make([]string, 0)

	for apIndex, ap := range apList {
		roleNames, err := uniqueRoleNameGenerator.Generate(&apList[apIndex])
		if err != nil {
			return 0, err
		}

		if ap.Delete {
			for _, access := range ap.Access {
				if access.ActualName == nil {
					logger.Warn(fmt.Sprintf("No actualname defined for deleted access %q. This will be ignored", access.Id))
					continue
				}

				roleName := *access.ActualName

				if !find(rolesToRemove, roleName) {
					rolesToRemove = append(rolesToRemove, roleName)
				}
			}
		} else {
			for _, access := range ap.Access {
				roleName := roleNames[access.Id]
				if _, f := apMap[roleName]; !f {
					apMap[roleName] = sync_to_target.EnrichedAccess{Access: access, AccessProvider: &apList[apIndex]}
				}
			}
		}
	}

	return timedExecution(func() error {
		return s.syncer.SyncAccessProviderToTarget(ctx, rolesToRemove, apMap, feedbackFile, &config.ConfigMap)
	})
}

func mapErrorToAccessSyncResult(err error) access_provider.AccessSyncResult {
	return access_provider.AccessSyncResult{
		Error: e.ToErrorResult(err),
	}
}
