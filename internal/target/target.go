package target

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/jinzhu/copier"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"

	"github.com/raito-io/cli/base/util/error/grpc_error"
	iconfig "github.com/raito-io/cli/internal/config"
	"github.com/raito-io/cli/internal/constants"
	error2 "github.com/raito-io/cli/internal/error"
	"github.com/raito-io/cli/internal/health_check"
	"github.com/raito-io/cli/internal/target/types"
)

//go:generate go run github.com/vektra/mockery/v2 --name=TargetRunner --with-expecter --inpackage
type TargetRunner interface {
	TargetSync(ctx context.Context, targetConfig *types.BaseTargetConfig) (syncError error)
	Finalize(ctx context.Context, baseConfig *types.BaseConfig, options *Options) error
	RunType() string
}

func GetTargetConfig(targetName string, baseConfig *types.BaseConfig) (*types.BaseTargetConfig, error) {
	targets := viper.Get(constants.Targets)

	if targetList, ok := targets.([]interface{}); ok {
		for _, targetObj := range targetList {
			target, ok := targetObj.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid configuration structure found")
			}

			tConfig, err := buildTargetConfigFromMapForRun(baseConfig, target, nil)
			if err != nil {
				return nil, fmt.Errorf("error while parsing the target configuration: %s", err.Error())
			}

			if tConfig.Name == targetName {
				return tConfig, nil
			}
		}
	}

	return nil, nil
}

func RunTargets(ctx context.Context, baseConfig *types.BaseConfig, runTarget TargetRunner, opFns ...func(*Options)) (err error) {
	options := createOptions(opFns...)

	defer func() {
		notifyErr := runTarget.Finalize(ctx, baseConfig, &options)
		if notifyErr != nil {
			err = multierror.Append(err, notifyErr)
		}
	}()

	defer func() {
		if r := recover(); r != nil {
			err = multierror.Append(err, error2.NewRecoverErrorf("recover from panic in RunTargets: %v", r))
		}
	}()

	if viper.GetString(constants.ConnectorNameFlag) != "" {
		targetConfig := buildTargetConfigFromFlags(baseConfig)

		if !options.SyncDataSourceId(targetConfig.DataSourceId) {
			return nil
		}

		if !options.SyncIdentityStoreId(targetConfig.IdentityStoreId) {
			return nil
		}

		logTargetConfig(targetConfig)

		err2 := runTarget.TargetSync(ctx, options.TargetOptions(targetConfig))
		if err2 != nil {
			return err2
		}
	} else {
		err2 := runMultipleTargets(ctx, baseConfig, runTarget.RunType(), runTarget.TargetSync, &options)
		if err2 != nil {
			return err2
		}
	}

	return nil
}

func HandleTargetError(err error, config *types.BaseTargetConfig, prefix ...string) {
	targetError := &grpc_error.InternalPluginStatusError{}

	prefixString := strings.Join(prefix, " ")

	if prefixString != "" {
		prefixString += ": "
	}

	if errors.As(err, &targetError) && targetError.StatusCode() == codes.InvalidArgument {
		config.TargetLogger.Error(fmt.Sprintf("%s%s. Execute command 'info <connector>' to print out the expected parameters for the connector.", prefixString, targetError.Error()))
		return
	}

	config.TargetLogger.Error(fmt.Sprintf("%s%s", prefixString, err.Error()))
}

func runMultipleTargets(ctx context.Context, baseConfig *types.BaseConfig, runType string, runTarget func(ctx context.Context, tConfig *types.BaseTargetConfig) error, options *Options) error {
	var errorResult error

	dataObjectEnricherMap, err := buildDataObjectEnricherMap()
	if err != nil {
		errorResult = multierror.Append(errorResult, err)
	}

	targets := viper.Get(constants.Targets)
	onlyTargets := make(map[string]struct{})

	onlyTargetsS := viper.GetString(constants.OnlyTargetsFlag)
	if onlyTargetsS != "" {
		for _, ot := range strings.Split(onlyTargetsS, ",") {
			onlyTargets[strings.TrimSpace(ot)] = struct{}{}
		}
	}

	if targetList, ok := targets.([]interface{}); ok {
		hclog.L().Debug(fmt.Sprintf("Found %d targets to run.", len(targetList)))

		for _, targetObj := range targetList {
			target, ok := targetObj.(map[string]interface{})
			if !ok {
				errorResult = multierror.Append(errorResult, fmt.Errorf("the target definition could not be parsed correctly (%v)", targetObj))
				hclog.L().Debug(fmt.Sprintf("The target definition could not be parsed correctly (%v)", targetObj))

				continue
			}

			tConfig, err := buildTargetConfigFromMapForRun(baseConfig, target, dataObjectEnricherMap)
			if err != nil {
				errorResult = multierror.Append(errorResult, fmt.Errorf("error while parsing the target configuration: %s", err.Error()))
				hclog.L().Error(fmt.Sprintf("error while parsing the target configuration: %s", err.Error()))

				continue
			}

			if tConfig == nil {
				continue
			}

			if !options.SyncDataSourceId(tConfig.DataSourceId) {
				continue
			}

			if !options.SyncIdentityStoreId(tConfig.IdentityStoreId) {
				continue
			}

			tConfig = options.TargetOptions(tConfig)

			if len(onlyTargets) > 0 {
				if _, found := onlyTargets[tConfig.Name]; !found {
					tConfig.TargetLogger.Info("Skipping target", "success")
					continue
				}
			}

			logTargetConfig(tConfig)

			err2 := tConfig.CalculateFileBackupLocationForRun(runType)
			if err2 != nil {
				hclog.L().Error(err2.Error())

				continue
			}

			runErr := runTarget(ctx, tConfig)
			if runErr != nil {
				errorResult = multierror.Append(errorResult, runErr)

				// In debug as the error should already be outputted, and we are ignoring it here.
				tConfig.TargetLogger.Debug("Error while executing target", "error", runErr.Error())
			}

			tConfig.FinalizeRun()
		}
	}

	return errorResult
}

func buildTargetConfigFromMapForRun(baseconfig *types.BaseConfig, target map[string]interface{}, dataObjectEnricherMap map[string]*types.EnricherConfig) (*types.BaseTargetConfig, error) { //nolint:cyclop
	tConfig := types.BaseTargetConfig{
		BaseConfig:      *baseconfig,
		DeleteUntouched: true,
		DeleteTempFiles: true,
		ReplaceGroups:   true,
	}

	err := fillStruct(&tConfig, target)
	if err != nil {
		return nil, err
	}

	tConfig.Parameters = make(map[string]string)

	for k, v := range target {
		if _, f := constants.KnownFlags[k]; f {
			continue
		}

		cv, err2 := iconfig.HandleField(v, reflect.String)
		if err2 != nil {
			return nil, err2
		}

		stringValue, err2 := argumentToString(cv)
		if err2 != nil {
			return nil, err2
		}

		if stringValue != nil {
			tConfig.Parameters[k] = *stringValue
		}
	}

	if tConfig.Name == "" {
		tConfig.Name = tConfig.ConnectorName
	}

	// Create a logger to add the target log name to each log message.
	tConfig.TargetLogger = baseconfig.BaseLogger.With("target", tConfig.Name)

	// Merge with some global parameters
	tConfig.SkipDataAccessSync = tConfig.SkipDataAccessSync || viper.GetBool(constants.SkipDataAccessSyncFlag)
	tConfig.SkipDataSourceSync = tConfig.SkipDataSourceSync || viper.GetBool(constants.SkipDataSourceSyncFlag)
	tConfig.SkipIdentityStoreSync = tConfig.SkipIdentityStoreSync || viper.GetBool(constants.SkipIdentityStoreSyncFlag)
	tConfig.SkipDataUsageSync = tConfig.SkipDataUsageSync || viper.GetBool(constants.SkipDataUsageSyncFlag)
	tConfig.SkipResourceProvider = tConfig.SkipResourceProvider || viper.GetBool(constants.SkipResourceProviderFlag)
	tConfig.SkipTagSync = tConfig.SkipTagSync || viper.GetBool(constants.SkipTagFlag)

	// If not set in the target, we take the globally set values.
	if tConfig.ApiSecret == "" {
		cv, err2 := iconfig.HandleField(viper.GetString(constants.ApiSecretFlag), reflect.String)
		if err2 != nil {
			return nil, err2
		}
		tConfig.ApiSecret = cv.(string)
	}

	if tConfig.ApiUser == "" {
		cv, err2 := iconfig.HandleField(viper.GetString(constants.ApiUserFlag), reflect.String)
		if err2 != nil {
			return nil, err2
		}
		tConfig.ApiUser = cv.(string)
	}

	if tConfig.FileBackupLocation == "" {
		cv, err2 := iconfig.HandleField(viper.GetString(constants.FileBackupLocationFlag), reflect.String)
		if err2 != nil {
			return nil, err2
		}
		tConfig.FileBackupLocation = cv.(string)
	}

	if tConfig.MaximumBackupsPerTarget == 0 {
		tConfig.MaximumBackupsPerTarget = viper.GetInt(constants.MaximumBackupsPerTargetFlag)
	}

	if tConfig.Domain == "" {
		cv, err2 := iconfig.HandleField(viper.GetString(constants.DomainFlag), reflect.String)
		if err2 != nil {
			return nil, err2
		}
		tConfig.Domain = cv.(string)
	}

	err = addDataObjectEnrichersToTargetConfig(&tConfig, target, dataObjectEnricherMap)
	if err != nil {
		return nil, err
	}

	return &tConfig, nil
}

func BuildBaseConfigFromFlags(baseLogger hclog.Logger, healthChecker health_check.HealthChecker, otherArgs []string) (*types.BaseConfig, error) {
	config := types.BaseConfig{
		BaseLogger:    baseLogger,
		HealthChecker: healthChecker,
		OtherArgs:     otherArgs,
	}

	err := config.ReloadConfig()
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func buildTargetConfigFromFlags(baseConfig *types.BaseConfig) *types.BaseTargetConfig {
	connector := viper.GetString(constants.ConnectorNameFlag)
	version := viper.GetString(constants.ConnectorVersionFlag)
	name := viper.GetString(constants.NameFlag)

	if name == "" {
		name = connector
	}

	targetConfig := types.BaseTargetConfig{
		BaseConfig:            *baseConfig,
		ConnectorName:         connector,
		ConnectorVersion:      version,
		Name:                  name,
		DataSourceId:          viper.GetString(constants.DataSourceIdFlag),
		IdentityStoreId:       viper.GetString(constants.IdentityStoreIdFlag),
		SkipIdentityStoreSync: viper.GetBool(constants.SkipIdentityStoreSyncFlag),
		SkipDataSourceSync:    viper.GetBool(constants.SkipDataSourceSyncFlag),
		SkipDataAccessSync:    viper.GetBool(constants.SkipDataAccessSyncFlag),
		SkipDataUsageSync:     viper.GetBool(constants.SkipDataUsageSyncFlag),
		SkipResourceProvider:  viper.GetBool(constants.SkipResourceProviderFlag),
		SkipTagSync:           viper.GetBool(constants.SkipTagFlag),
		LockAllWho:            viper.GetBool(constants.LockAllWhoFlag),
		LockWhoByName:         viper.GetString(constants.LockWhoByNameFlag),
		LockWhoByTag:          viper.GetString(constants.LockWhoByTagFlag),
		LockWhoWhenIncomplete: viper.GetBool(constants.LockWhoWhenIncompleteFlag),

		LockAllInheritance:            viper.GetBool(constants.LockAllInheritanceFlag),
		LockInheritanceByName:         viper.GetString(constants.LockInheritanceByNameFlag),
		LockInheritanceByTag:          viper.GetString(constants.LockInheritanceByTagFlag),
		LockInheritanceWhenIncomplete: viper.GetBool(constants.LockInheritanceWhenIncompleteFlag),

		LockAllOwners: viper.GetBool(constants.LockAllOwnersFlag),

		LockAllWhat:            viper.GetBool(constants.LockAllWhatFlag),
		LockWhatByName:         viper.GetString(constants.LockWhatByNameFlag),
		LockWhatByTag:          viper.GetString(constants.LockWhatByTagFlag),
		LockWhatWhenIncomplete: viper.GetBool(constants.LockWhatWhenIncompleteFlag),

		LockAllNames:            viper.GetBool(constants.LockAllNamesFlag),
		LockNamesByName:         viper.GetString(constants.LockNamesByNameFlag),
		LockNamesByTag:          viper.GetString(constants.LockNamesByTagFlag),
		LockNamesWhenIncomplete: viper.GetBool(constants.LockNamesWhenIncompleteFlag),

		LockAllDelete:            viper.GetBool(constants.LockAllDeleteFlag),
		LockDeleteByName:         viper.GetString(constants.LockDeleteByNameFlag),
		LockDeleteByTag:          viper.GetString(constants.LockDeleteByTagFlag),
		LockDeleteWhenIncomplete: viper.GetBool(constants.LockDeleteWhenIncompleteFlag),

		MakeNotInternalizable:   strings.TrimSpace(viper.GetString(constants.MakeNotInternalizableFlag)),
		FullyLockAll:            viper.GetBool(constants.FullyLockAllFlag),
		FullyLockByName:         viper.GetString(constants.FullyLockByNameFlag),
		FullyLockByTag:          viper.GetString(constants.FullyLockByTagFlag),
		FullyLockWhenIncomplete: viper.GetBool(constants.FullyLockWhenIncompleteFlag),

		TargetLogger:    baseConfig.BaseLogger.With("target", name),
		DeleteUntouched: true,
		DeleteTempFiles: true,
		ReplaceGroups:   true,
	}

	return &targetConfig
}

// logTargetConfig will print out the target configuration in the log (debug level).
// It will censure the sensitive information (secrets and passwords) if it is set.
func logTargetConfig(config *types.BaseTargetConfig) {
	if !hclog.L().IsDebug() {
		return
	}
	cc := types.BaseTargetConfig{}

	err := copier.Copy(&cc, config)
	if err != nil {
		hclog.L().Error("Error while copying config")
		return
	}
	cc.Parameters = make(map[string]string)

	err = copier.Copy(&cc.Parameters, config.Parameters)
	if err != nil {
		hclog.L().Error("Error while copying config")
		return
	}

	if cc.ApiSecret != "" {
		cc.ApiSecret = "**censured**"
	}

	for k := range cc.Parameters {
		lk := strings.ToLower(k)
		if strings.Contains(lk, "secret") || strings.Contains(lk, "password") || strings.Contains(lk, "passwd") || strings.Contains(lk, "psswd") {
			cc.Parameters[k] = "**censured**"
		}
	}

	hclog.L().Debug(fmt.Sprintf("Using target config (censured): %+v", cc))
}

type Options struct {
	ExternalTrigger  bool
	DataSourceIds    map[string]struct{}
	IdentityStoreIds map[string]struct{}
	ConfigOption     func(targetConfig *types.BaseTargetConfig)
}

func createOptions(opFns ...func(*Options)) Options {
	result := Options{}
	for _, fn := range opFns {
		fn(&result)
	}

	return result
}

func (o *Options) SyncDataSourceId(dataSourceId string) bool {
	if o.DataSourceIds == nil {
		return true
	}

	_, found := o.DataSourceIds[dataSourceId]

	return found
}

func (o *Options) SyncIdentityStoreId(identityStoreId string) bool {
	if o.IdentityStoreIds == nil {
		return true
	}

	_, found := o.IdentityStoreIds[identityStoreId]

	return found
}

func (o *Options) TargetOptions(targetConfig *types.BaseTargetConfig) *types.BaseTargetConfig {
	if o.ConfigOption != nil {
		o.ConfigOption(targetConfig)
	}

	return targetConfig
}

func WithDataSourceIds(dataSourceIds ...string) func(o *Options) {
	return func(o *Options) {
		if len(dataSourceIds) == 0 {
			return
		}

		if o.DataSourceIds == nil {
			o.DataSourceIds = map[string]struct{}{}
		}

		for _, dataSourceId := range dataSourceIds {
			o.DataSourceIds[dataSourceId] = struct{}{}
		}
	}
}

func WithIdentityStoreIds(identityStoreIds ...string) func(o *Options) {
	return func(o *Options) {
		if len(identityStoreIds) == 0 {
			return
		}

		if o.IdentityStoreIds == nil {
			o.IdentityStoreIds = map[string]struct{}{}
		}

		for _, identityStoreId := range identityStoreIds {
			o.IdentityStoreIds[identityStoreId] = struct{}{}
		}
	}
}

func WithConfigOption(fn func(targetConfig *types.BaseTargetConfig)) func(o *Options) {
	return func(o *Options) {
		o.ConfigOption = fn
	}
}

func WithExternalTrigger() func(o *Options) {
	return func(o *Options) {
		o.ExternalTrigger = true
	}
}
