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
	gql "github.com/raito-io/cli/internal/graphql"
	"github.com/raito-io/cli/internal/target/types"
)

func RunTargets(ctx context.Context, baseConfig *types.BaseConfig, runTarget func(ctx context.Context, tConfig *types.BaseTargetConfig) (string, error), opFns ...func(*Options)) (err error) {
	options := createOptions(opFns...)

	var jobIds []string

	defer func() {
		baseConfig.BaseLogger.Debug(fmt.Sprintf("Finished targets with jobIds: %v", jobIds))

		notifyErr := sendEndOfTarget(ctx, baseConfig, jobIds, &options)
		if notifyErr != nil {
			err = multierror.Append(err, notifyErr)
		}
	}()

	if viper.GetString(constants.ConnectorNameFlag) != "" {
		targetConfig := buildTargetConfigFromFlags(baseConfig)

		if !options.SyncDataSourceId(targetConfig.DataSourceId) {
			return nil
		}

		logTargetConfig(targetConfig)

		jobId, err := runTarget(ctx, options.TargetOptions(targetConfig))
		if jobId != "" {
			jobIds = append(jobIds, jobId)
		}

		if err != nil {
			return err
		}
	} else {
		jobs, err := runMultipleTargets(ctx, baseConfig, runTarget, &options)
		jobIds = jobs

		if err != nil {
			return err
		}
	}

	return nil
}

func sendEndOfTarget(ctx context.Context, baseConfig *types.BaseConfig, jobIds []string, options *Options) error {
	if !options.ExternalTrigger {
		var mutation struct {
			EndOfTargetsSyncResult struct {
				EndOfTargetsSync struct {
					Success bool
				} `graphql:"... on EndOfTargetsSync"`
				PermissionDeniedError struct {
					message string
				} `graphql:"... on PermissionDeniedError"`
			} `graphql:"endOfTargetsSync(input: $input)"`
		}

		type EndOfTargetsSyncInput struct {
			JobIds []string `json:"jobIds"`
		}

		input := EndOfTargetsSyncInput{JobIds: jobIds}

		err := gql.NewClient(baseConfig).Mutate(ctx, &mutation, map[string]interface{}{"input": input})
		if err != nil {
			return err
		}

		if mutation.EndOfTargetsSyncResult.PermissionDeniedError.message != "" {
			baseConfig.BaseLogger.Error(fmt.Sprintf("Permission denied to notify end of all targets: %s", mutation.EndOfTargetsSyncResult.PermissionDeniedError.message))
		} else if !mutation.EndOfTargetsSyncResult.EndOfTargetsSync.Success {
			baseConfig.BaseLogger.Warn("Failed to notify end of all targets")
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

func runMultipleTargets(ctx context.Context, baseconfig *types.BaseConfig, runTarget func(ctx context.Context, tConfig *types.BaseTargetConfig) (string, error), options *Options) ([]string, error) {
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

	var jobIds []string

	if targetList, ok := targets.([]interface{}); ok {
		hclog.L().Debug(fmt.Sprintf("Found %d targets to run.", len(targetList)))

		for _, targetObj := range targetList {
			target, ok := targetObj.(map[string]interface{})
			if !ok {
				errorResult = multierror.Append(errorResult, fmt.Errorf("the target definition could not be parsed correctly (%v)", targetObj))
				hclog.L().Debug(fmt.Sprintf("The target definition could not be parsed correctly (%v)", targetObj))

				continue
			}

			tConfig, err := buildTargetConfigFromMap(baseconfig, target, dataObjectEnricherMap)
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

			tConfig = options.TargetOptions(tConfig)

			if len(onlyTargets) > 0 {
				if _, found := onlyTargets[tConfig.Name]; !found {
					tConfig.TargetLogger.Info("Skipping target", "success")
					continue
				}
			}

			logTargetConfig(tConfig)

			jobId, err := runTarget(ctx, tConfig)
			if err != nil {
				errorResult = multierror.Append(errorResult, err)

				// In debug as the error should already be outputted, and we are ignoring it here.
				tConfig.TargetLogger.Debug("Error while executing target", "error", err.Error())
			}

			if jobId != "" {
				jobIds = append(jobIds, jobId)
			}
		}
	}

	return jobIds, errorResult
}

func buildTargetConfigFromMap(baseconfig *types.BaseConfig, target map[string]interface{}, dataObjectEnricherMap map[string]*types.EnricherConfig) (*types.BaseTargetConfig, error) {
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

func BuildBaseConfigFromFlags(baseLogger hclog.Logger, otherArgs []string) (*types.BaseConfig, error) {
	apiUser, err := iconfig.HandleField(viper.GetString(constants.ApiUserFlag), reflect.String)
	if err != nil {
		return nil, err
	}

	apiSecret, err := iconfig.HandleField(viper.GetString(constants.ApiSecretFlag), reflect.String)
	if err != nil {
		return nil, err
	}

	domain, err := iconfig.HandleField(viper.GetString(constants.DomainFlag), reflect.String)
	if err != nil {
		return nil, err
	}

	config := types.BaseConfig{
		BaseLogger: baseLogger,
		ApiUser:    apiUser.(string),
		ApiSecret:  apiSecret.(string),
		Domain:     domain.(string),
	}

	config.Parameters = BuildParameterMapFromArguments(otherArgs)

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
		LockAllWho:            viper.GetBool(constants.LockAllWhoFlag),
		LockAllInheritance:    viper.GetBool(constants.LockAllInheritanceFlag),
		LockAllWhat:           viper.GetBool(constants.LockAllWhatFlag),
		LockAllNames:          viper.GetBool(constants.LockAllNamesFlag),
		LockAllDelete:         viper.GetBool(constants.LockAllDeleteFlag),
		TargetLogger:          baseConfig.BaseLogger.With("target", name),
		DeleteUntouched:       true,
		DeleteTempFiles:       true,
		ReplaceGroups:         true,
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
	ExternalTrigger bool
	DataSourceIds   map[string]struct{}
	ConfigOption    func(targetConfig *types.BaseTargetConfig)
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

func (o *Options) TargetOptions(targetConfig *types.BaseTargetConfig) *types.BaseTargetConfig {
	if o.ConfigOption != nil {
		o.ConfigOption(targetConfig)
	}

	return targetConfig
}

func WithDataSourceIds(dataSourceIds ...string) func(o *Options) {
	return func(o *Options) {
		if o.DataSourceIds == nil {
			o.DataSourceIds = map[string]struct{}{}
		}

		for _, dataSourceId := range dataSourceIds {
			o.DataSourceIds[dataSourceId] = struct{}{}
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

func BuildParameterMapFromArguments(args []string) map[string]string {
	params := make(map[string]string)

	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "--") {
			arg := args[i][2:]
			if strings.Contains(arg, "=") {
				// The case where the flag is in the form of "--key=value"
				key := arg[0:strings.Index(arg, "=")]
				value := arg[strings.Index(arg, "=")+1:]
				params[key] = value
			} else if i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
				// The case where the flag is in the form of "--key value"
				params[arg] = args[i+1]
				i++
			} else {
				// Otherwise, we consider this a boolean flag
				params[arg] = "TRUE"
			}
		}
	}

	return params
}
