package role_based

import (
	"context"
	"fmt"
	"strings"

	"github.com/raito-io/cli/base"
	"github.com/raito-io/cli/base/access_provider"
	"github.com/raito-io/cli/base/access_provider/sync_to_target"
	"github.com/raito-io/cli/base/access_provider/sync_to_target/naming_hint"
	"github.com/raito-io/cli/base/util/config"
	"github.com/raito-io/cli/base/wrappers"
)

var logger = base.Logger()

//go:generate go run github.com/vektra/mockery/v2 --name=AccessProviderRoleSyncer --with-expecter --inpackage
type AccessProviderRoleSyncer interface {
	SyncAccessProvidersFromTarget(ctx context.Context, accessProviderHandler wrappers.AccessProviderHandler, configMap *config.ConfigMap) error

	SyncAccessProviderRolesToTarget(ctx context.Context, apToRemoveMap map[string]*sync_to_target.AccessProvider, apMap map[string]*sync_to_target.AccessProvider, feedbackHandler wrappers.AccessProviderFeedbackHandler, configMap *config.ConfigMap) error
	SyncAccessProviderMasksToTarget(ctx context.Context, apToRemoveMap map[string]*sync_to_target.AccessProvider, apMap map[string]*sync_to_target.AccessProvider, roleNameMap map[string]string, feedbackHandler wrappers.AccessProviderFeedbackHandler, configMap *config.ConfigMap) error

	SyncAccessAsCodeToTarget(ctx context.Context, accesses map[string]*sync_to_target.AccessProvider, prefix string, configMap *config.ConfigMap) error
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
	apMap := make(map[string]*sync_to_target.AccessProvider)

	for _, ap := range apList {
		roleName, err := uniqueRoleNameGenerator.Generate(ap)
		if err != nil {
			return err
		}

		logger.Info(fmt.Sprintf("Generated rolename %q", roleName))
		apMap[roleName] = ap
	}

	return s.syncer.SyncAccessAsCodeToTarget(ctx, apMap, prefix, configMap)
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

	rolesMap := make(map[string]*sync_to_target.AccessProvider)
	rolesToRemove := make(map[string]*sync_to_target.AccessProvider)

	for _, ap := range apList {
		var err2 error
		if ap.Action == sync_to_target.Mask {
			_, masksMap, masksToRemove, err2 = handleAccessProvider(ap, masksMap, masksToRemove, accessProviderFeedbackHandler, uniqueRoleNameGenerator)
		} else {
			var roleName string
			roleName, rolesMap, rolesToRemove, err2 = handleAccessProvider(ap, rolesMap, rolesToRemove, accessProviderFeedbackHandler, uniqueRoleNameGenerator)
			apIdNameMap[ap.Id] = roleName
		}

		if err2 != nil {
			return err2
		}
	}

	// Step 1 first initiate all the masks
	err = s.syncer.SyncAccessProviderMasksToTarget(ctx, masksToRemove, masksMap, apIdNameMap, accessProviderFeedbackHandler, configMap)
	if err != nil {
		return err
	}

	// Step 2 then initiate all the roles
	return s.syncer.SyncAccessProviderRolesToTarget(ctx, rolesToRemove, rolesMap, accessProviderFeedbackHandler, configMap)
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
