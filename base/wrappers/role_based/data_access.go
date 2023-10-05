package role_based

import (
	"context"
	"fmt"
	"strings"

	"github.com/raito-io/golang-set/set"

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

	SyncAccessProviderRolesToTarget(ctx context.Context, rolesToRemove []string, access map[string]*sync_to_target.AccessProvider, feedbackHandler wrappers.AccessProviderFeedbackHandler, configMap *config.ConfigMap) error
	SyncAccessProviderMasksToTarget(ctx context.Context, masksToRemove []string, access []*sync_to_target.AccessProvider, feedbackHandler wrappers.AccessProviderFeedbackHandler, configMap *config.ConfigMap) error

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

	var masksList []*sync_to_target.AccessProvider
	masksToRemove := set.NewSet[string]()

	rolesMap := make(map[string]*sync_to_target.AccessProvider)
	rolesToRemove := set.NewSet[string]()

	for _, ap := range apList {
		if ap.Action == sync_to_target.Mask {
			masksList, masksToRemove = parseMask(ap, masksList, masksToRemove)
		} else {
			rolesMap, rolesToRemove, err = parseRole(ap, rolesMap, rolesToRemove, uniqueRoleNameGenerator)
			if err != nil {
				return err
			}
		}
	}

	// Step 1 first initiate all the masks
	err = s.syncer.SyncAccessProviderMasksToTarget(ctx, masksToRemove.Slice(), masksList, accessProviderFeedbackHandler, configMap)
	if err != nil {
		return err
	}

	// Step 2 then initiate all the roles
	return s.syncer.SyncAccessProviderRolesToTarget(ctx, rolesToRemove.Slice(), rolesMap, accessProviderFeedbackHandler, configMap)
}

func parseMask(mask *sync_to_target.AccessProvider, masksList []*sync_to_target.AccessProvider, masksToRemove set.Set[string]) ([]*sync_to_target.AccessProvider, set.Set[string]) {
	if mask.Delete {
		if mask.ActualName == nil {
			logger.Warn(fmt.Sprintf("No actualname defined for deleted access provider %q. This will be ignored", mask.Id))
			return masksList, masksToRemove
		}

		masksToRemove.Add(*mask.ActualName)
	} else {
		masksList = append(masksList, mask)
	}

	return masksList, masksToRemove
}

func parseRole(ap *sync_to_target.AccessProvider, rolesMap map[string]*sync_to_target.AccessProvider, rolesToRemove set.Set[string], roleNameGenerator naming_hint.UniqueGenerator) (map[string]*sync_to_target.AccessProvider, set.Set[string], error) {
	if ap.Delete {
		if ap.ActualName == nil {
			logger.Warn(fmt.Sprintf("No actualname defined for deleted access provider %q. This will be ignored", ap.Id))
			return rolesMap, rolesToRemove, nil
		}

		rolesToRemove.Add(*ap.ActualName)
	} else {
		roleName, err := roleNameGenerator.Generate(ap)
		if err != nil {
			return rolesMap, rolesToRemove, err
		}

		if _, f := rolesMap[roleName]; !f {
			rolesMap[roleName] = ap
		}
	}

	return rolesMap, rolesToRemove, nil
}
