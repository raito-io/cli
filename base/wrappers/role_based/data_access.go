package role_based

import (
	"context"
	"fmt"

	"github.com/raito-io/cli/base"
	"github.com/raito-io/cli/base/access_provider"
	"github.com/raito-io/cli/base/access_provider/sync_to_target"
	"github.com/raito-io/cli/base/access_provider/sync_to_target/naming_hint"
	"github.com/raito-io/cli/base/access_provider/types"
	"github.com/raito-io/cli/base/util/config"
	"github.com/raito-io/cli/base/wrappers"
)

var logger = base.Logger()

//go:generate go run github.com/vektra/mockery/v2 --name=AccessProviderRoleSyncer --with-expecter --inpackage
type AccessProviderRoleSyncer interface {
	SyncAccessProvidersFromTarget(ctx context.Context, accessProviderHandler wrappers.AccessProviderHandler, configMap *config.ConfigMap) error

	SyncAccessProviderRolesToTarget(ctx context.Context, apToRemoveMap map[string]*sync_to_target.AccessProvider, apMap map[string]*sync_to_target.AccessProvider, feedbackHandler wrappers.AccessProviderFeedbackHandler, configMap *config.ConfigMap) error
	SyncAccessProviderMasksToTarget(ctx context.Context, apToRemoveMap map[string]*sync_to_target.AccessProvider, apMap map[string]*sync_to_target.AccessProvider, roleNameMap map[string]string, feedbackHandler wrappers.AccessProviderFeedbackHandler, configMap *config.ConfigMap) error
	SyncAccessProviderFiltersToTarget(ctx context.Context, apToRemoveMap map[string]*sync_to_target.AccessProvider, apMap map[string]*sync_to_target.AccessProvider, roleNameMap map[string]string, feedbackHandler wrappers.AccessProviderFeedbackHandler, configMap *config.ConfigMap) error
}

func AccessProviderRoleSync(syncer AccessProviderRoleSyncer, namingConstraints naming_hint.NamingConstraints, configOpt ...func(config *access_provider.AccessSyncConfig)) *wrappers.DataAccessSyncFunction {
	configOpt = append([]func(config *access_provider.AccessSyncConfig){access_provider.WithSupportPartialSync()}, configOpt...)

	roleSync := &accessProviderRoleSyncFunction{
		syncer:            syncer,
		namingConstraints: namingConstraints,
	}

	return wrappers.DataAccessSync(roleSync, configOpt...)
}

type accessProviderRoleSyncFunction struct {
	syncer            AccessProviderRoleSyncer
	namingConstraints naming_hint.NamingConstraints
}

func (s *accessProviderRoleSyncFunction) SyncAccessProvidersFromTarget(ctx context.Context, accessProviderHandler wrappers.AccessProviderHandler, configMap *config.ConfigMap) error {
	return s.syncer.SyncAccessProvidersFromTarget(ctx, accessProviderHandler, configMap)
}

func (s *accessProviderRoleSyncFunction) SyncAccessProviderToTarget(ctx context.Context, accessProviders *sync_to_target.AccessProviderImport, accessProviderFeedbackHandler wrappers.AccessProviderFeedbackHandler, configMap *config.ConfigMap) error {
	uniqueRoleNameGenerator, err := naming_hint.NewUniqueNameGenerator(logger, "", &s.namingConstraints)
	if err != nil {
		return err
	}

	apList := accessProviders.AccessProviders

	apIdNameMap := make(map[string]string)

	masksMap := make(map[string]*sync_to_target.AccessProvider)
	masksToRemove := make(map[string]*sync_to_target.AccessProvider)

	filtersMap := make(map[string]*sync_to_target.AccessProvider)
	filtersToRemove := make(map[string]*sync_to_target.AccessProvider)

	rolesMap := make(map[string]*sync_to_target.AccessProvider)
	rolesToRemove := make(map[string]*sync_to_target.AccessProvider)

	for _, ap := range apList {
		var err2 error

		switch ap.Action {
		case types.Mask:
			_, masksMap, masksToRemove, err2 = handleAccessProvider(ap, masksMap, masksToRemove, accessProviderFeedbackHandler, uniqueRoleNameGenerator)
		case types.Filtered:
			_, filtersMap, filtersToRemove, err2 = handleAccessProvider(ap, filtersMap, filtersToRemove, accessProviderFeedbackHandler, uniqueRoleNameGenerator)
		case types.Grant, types.Purpose:
			var roleName string
			roleName, rolesMap, rolesToRemove, err2 = handleAccessProvider(ap, rolesMap, rolesToRemove, accessProviderFeedbackHandler, uniqueRoleNameGenerator)
			apIdNameMap[ap.Id] = roleName
		default:
			err2 = accessProviderFeedbackHandler.AddAccessProviderFeedback(sync_to_target.AccessProviderSyncFeedback{
				AccessProvider: ap.Id,
				Errors:         []string{fmt.Sprintf("Unsupported action %s", ap.Action.String())},
			})
		}

		if err2 != nil {
			return err2
		}
	}

	// Step 1 first initiate all the masks
	if len(masksMap) > 0 || len(masksToRemove) > 0 {
		err = s.syncer.SyncAccessProviderMasksToTarget(ctx, masksToRemove, masksMap, apIdNameMap, accessProviderFeedbackHandler, configMap)
		if err != nil {
			return fmt.Errorf("sync masks to target: %w", err)
		}
	}

	// Step 2 then initialize all filters
	if len(filtersMap) > 0 || len(filtersToRemove) > 0 {
		err = s.syncer.SyncAccessProviderFiltersToTarget(ctx, filtersToRemove, filtersMap, apIdNameMap, accessProviderFeedbackHandler, configMap)
		if err != nil {
			return fmt.Errorf("sync filters to target: %w", err)
		}
	}

	// Step 3 then initiate all the roles
	err = s.syncer.SyncAccessProviderRolesToTarget(ctx, rolesToRemove, rolesMap, accessProviderFeedbackHandler, configMap)
	if err != nil {
		return fmt.Errorf("sync roles to target: %w", err)
	}

	return nil
}

func handleAccessProvider(ap *sync_to_target.AccessProvider, apMap map[string]*sync_to_target.AccessProvider, apToRemoveMap map[string]*sync_to_target.AccessProvider, accessProviderFeedbackHandler wrappers.AccessProviderFeedbackHandler, roleNameGenerator naming_hint.UniqueGenerator) (string, map[string]*sync_to_target.AccessProvider, map[string]*sync_to_target.AccessProvider, error) {
	var roleName string

	if ap.Delete {
		if ap.ActualName == nil {
			logger.Warn(fmt.Sprintf("No actualname defined for deleted access provider %q. This will be ignored", ap.Id))

			err := accessProviderFeedbackHandler.AddAccessProviderFeedback(sync_to_target.AccessProviderSyncFeedback{
				AccessProvider: ap.Id,
				ActualName:     "",
			})
			if err != nil {
				return "", nil, nil, err
			}

			return "", apMap, apToRemoveMap, nil
		}

		roleName = *ap.ActualName

		apToRemoveMap[roleName] = ap
	} else {
		var err error

		roleName, err = roleNameGenerator.Generate(ap)
		if err != nil {
			return "", apMap, apToRemoveMap, err
		}

		if _, f := apMap[roleName]; !f {
			apMap[roleName] = ap
		}
	}

	return roleName, apMap, apToRemoveMap, nil
}
