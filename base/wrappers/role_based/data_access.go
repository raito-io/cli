package role_based

import (
	"context"
	"fmt"
	"strings"

	"github.com/raito-io/cli/base/access_provider/sync_to_target"
	"github.com/raito-io/cli/base/access_provider/sync_to_target/naming_hint"
	"github.com/raito-io/cli/base/util/config"
	"github.com/raito-io/cli/base/wrappers"
)

//go:generate go run github.com/vektra/mockery/v2 --name=AccessProviderRoleSyncer --with-expecter --inpackage
type AccessProviderRoleSyncer interface {
	SyncAccessProvidersFromTarget(ctx context.Context, accessProviderHandler wrappers.AccessProviderHandler, configMap *config.ConfigMap) error
	SyncAccessProvidersToTarget(ctx context.Context, rolesToRemove []string, access map[string]sync_to_target.EnrichedAccess, feedbackHandler wrappers.AccessProviderFeedbackHandler, configMap *config.ConfigMap) error
	SyncAccessAsCodeToTarget(ctx context.Context, accesses map[string]sync_to_target.EnrichedAccess, prefix string, configMap *config.ConfigMap) error
}

func AccessProviderRoleSync(syncer AccessProviderRoleSyncer, namingConstraints naming_hint.NamingConstraints) *wrappers.DataAccessSyncFunction {
	roleSync := &accessProviderRoleSyncFunction{
		syncer:            syncer,
		namingConstraints: namingConstraints,
	}

	return wrappers.DataAccessSync(roleSync)
}

type accessProviderRoleSyncFunction struct {
	syncer            AccessProviderRoleSyncer
	namingConstraints naming_hint.NamingConstraints
}

func (s *accessProviderRoleSyncFunction) SyncAccessProvidersFromTarget(ctx context.Context, accessProviderHandler wrappers.AccessProviderHandler, configMap *config.ConfigMap) error {
	return s.syncer.SyncAccessProvidersFromTarget(ctx, accessProviderHandler, configMap)
}

func (s *accessProviderRoleSyncFunction) SyncAccessAsCodeToTarget(ctx context.Context, accessProviders *sync_to_target.AccessProviderImport, prefix string, configMap *config.ConfigMap) error {
	roleSeparator := string(s.namingConstraints.SplitCharacter())
	if !strings.HasSuffix(prefix, roleSeparator) {
		prefix += roleSeparator
	}

	logger.Info(fmt.Sprintf("Using prefix %q", prefix))

	uniqueRoleNameGenerator, err := naming_hint.NewUniqueNameGenerator(logger, prefix, &s.namingConstraints)
	if err != nil {
		return err
	}

	apList := accessProviders.AccessProviders
	apMap := make(map[string]sync_to_target.EnrichedAccess)

	for apIndex, ap := range apList {
		roleNames, err := uniqueRoleNameGenerator.GenerateOrdered(apList[apIndex])
		if err != nil {
			return err
		}

		for i, access := range ap.Access {
			roleName := roleNames[i]

			logger.Info(fmt.Sprintf("Generated rolename %q", roleName))
			apMap[roleName] = sync_to_target.EnrichedAccess{Access: access, AccessProvider: apList[apIndex]}
		}
	}

	return s.syncer.SyncAccessAsCodeToTarget(ctx, apMap, prefix, configMap)
}

func (s *accessProviderRoleSyncFunction) SyncAccessProviderToTarget(ctx context.Context, accessProviders *sync_to_target.AccessProviderImport, accessProviderFeedbackHandler wrappers.AccessProviderFeedbackHandler, configMap *config.ConfigMap) error {
	uniqueRoleNameGenerator, err := naming_hint.NewUniqueNameGenerator(logger, "", &s.namingConstraints)
	if err != nil {
		return err
	}

	apList := accessProviders.AccessProviders

	apMap := make(map[string]sync_to_target.EnrichedAccess)
	rolesToRemove := make([]string, 0)

	for apIndex, ap := range apList {
		roleNames, err := uniqueRoleNameGenerator.Generate(apList[apIndex])
		if err != nil {
			return err
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
					apMap[roleName] = sync_to_target.EnrichedAccess{Access: access, AccessProvider: apList[apIndex]}
				}
			}
		}
	}

	return s.syncer.SyncAccessProvidersToTarget(ctx, rolesToRemove, apMap, accessProviderFeedbackHandler, configMap)
}
