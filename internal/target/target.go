package target

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/jinzhu/copier"
	"github.com/spf13/viper"

	"github.com/raito-io/cli/base/util/config"
	error2 "github.com/raito-io/cli/base/util/error"
	iconfig "github.com/raito-io/cli/internal/config"
	"github.com/raito-io/cli/internal/constants"
)

type BaseConfig struct {
	config.ConfigMap
	ApiUser    string
	ApiSecret  string
	Domain     string
	BaseLogger hclog.Logger
}

type BaseTargetConfig struct {
	BaseConfig
	ConnectorName    string
	ConnectorVersion string
	Name             string
	DataSourceId     string
	IdentityStoreId  string

	SkipIdentityStoreSync bool
	SkipDataSourceSync    bool
	SkipDataAccessSync    bool
	SkipDataUsageSync     bool

	OnlyOutOfSyncData    bool
	SkipDataAccessImport bool

	DeleteUntouched bool
	ReplaceTags     bool
	DeleteTempFiles bool
	ReplaceGroups   bool

	TargetLogger hclog.Logger
}

func RunTargets(baseConfig *BaseConfig, runTarget func(tConfig *BaseTargetConfig) error, opFns ...func(*Options)) error {
	options := createOptions(opFns...)

	if viper.GetString(constants.ConnectorNameFlag) != "" {
		targetConfig := buildTargetConfigFromFlags(baseConfig)

		if !options.SyncDataSourceId(targetConfig.DataSourceId) {
			return nil
		}

		logTargetConfig(targetConfig)

		return runTarget(options.TargetOptions(targetConfig))
	} else {
		return runMultipleTargets(baseConfig, runTarget, &options)
	}
}

func HandleTargetError(err error, config *BaseTargetConfig, during string) {
	if errorResult, ok := err.(error2.ErrorResult); ok {
		if errorResult.ErrorCode == error2.BadInputParameterError || errorResult.ErrorCode == error2.MissingInputParameterError {
			config.TargetLogger.Error(fmt.Sprintf("Error during %s: %s. Execute command 'info <connector>' to print out the expected parameters for the connector.", during, errorResult.ErrorMessage))
			return
		}
	}

	config.TargetLogger.Error(fmt.Sprintf("Error during %s: %s", during, err.Error()))
}

func runMultipleTargets(baseconfig *BaseConfig, runTarget func(tConfig *BaseTargetConfig) error, options *Options) error {
	targets := viper.Get(constants.Targets)
	onlyTargets := make(map[string]struct{})

	onlyTargetsS := viper.GetString(constants.OnlyTargetsFlag)
	if onlyTargetsS != "" {
		for _, ot := range strings.Split(onlyTargetsS, ",") {
			onlyTargets[strings.TrimSpace(ot)] = struct{}{}
		}
	}

	var errorResult error

	if targetList, ok := targets.([]interface{}); ok {
		hclog.L().Debug(fmt.Sprintf("Found %d targets to run.", len(targetList)))

		for _, targetObj := range targetList {
			target, ok := targetObj.(map[string]interface{})
			if !ok {
				errorResult = multierror.Append(errorResult, fmt.Errorf("the target definition could not be parsed correctly (%v)", targetObj))
				hclog.L().Debug(fmt.Sprintf("The target definition could not be parsed correctly (%v)", targetObj))

				break
			}

			tConfig, err := buildTargetConfigFromMap(baseconfig, target)
			if err != nil {
				errorResult = multierror.Append(errorResult, fmt.Errorf("error while parsing target configuration: %s", err.Error()))
				hclog.L().Error(fmt.Sprintf("error while parsing target configuration: %s", err.Error()))

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

			err = runTarget(tConfig)
			if err != nil {
				errorResult = multierror.Append(errorResult, err)

				// In debug as the error should already be outputted, and we are ignoring it here.
				tConfig.TargetLogger.Debug("Error while executing target", "error", err.Error())
			}
		}
	}

	return errorResult
}

func buildTargetConfigFromMap(baseconfig *BaseConfig, target map[string]interface{}) (*BaseTargetConfig, error) {
	tConfig := BaseTargetConfig{
		BaseConfig:      *baseconfig,
		DeleteUntouched: true,
		DeleteTempFiles: true,
		ReplaceTags:     true,
		ReplaceGroups:   true,
	}
	err := fillStruct(&tConfig, target)

	if err != nil {
		return nil, err
	}
	tConfig.Parameters = make(map[string]interface{})

	for k, v := range target {
		if _, f := constants.KnownFlags[k]; !f {
			cv, err := iconfig.HandleField(v, reflect.String)
			if err != nil {
				return nil, err
			}
			tConfig.Parameters[k] = cv
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
		cv, err := iconfig.HandleField(viper.GetString(constants.ApiSecretFlag), reflect.String)
		if err != nil {
			return nil, err
		}
		tConfig.ApiSecret = cv.(string)
	}

	if tConfig.ApiUser == "" {
		cv, err := iconfig.HandleField(viper.GetString(constants.ApiUserFlag), reflect.String)
		if err != nil {
			return nil, err
		}
		tConfig.ApiUser = cv.(string)
	}

	if tConfig.Domain == "" {
		cv, err := iconfig.HandleField(viper.GetString(constants.DomainFlag), reflect.String)
		if err != nil {
			return nil, err
		}
		tConfig.Domain = cv.(string)
	}

	return &tConfig, nil
}

func buildParameterMapFromArguments(args []string) map[string]interface{} {
	params := make(map[string]interface{})

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
				params[arg] = true
			}
		}
	}

	return params
}

func BuildBaseConfigFromFlags(baseLogger hclog.Logger, otherArgs []string) (*BaseConfig, error) {
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

	config := BaseConfig{
		BaseLogger: baseLogger,
		ApiUser:    apiUser.(string),
		ApiSecret:  apiSecret.(string),
		Domain:     domain.(string),
	}

	config.Parameters = buildParameterMapFromArguments(otherArgs)

	return &config, nil
}

func buildTargetConfigFromFlags(baseConfig *BaseConfig) *BaseTargetConfig {
	connector := viper.GetString(constants.ConnectorNameFlag)
	version := viper.GetString(constants.ConnectorVersionFlag)
	name := viper.GetString(constants.NameFlag)

	if name == "" {
		name = connector
	}

	targetConfig := BaseTargetConfig{
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
		TargetLogger:          baseConfig.BaseLogger.With("target", name),
		DeleteUntouched:       true,
		DeleteTempFiles:       true,
		ReplaceTags:           true,
		ReplaceGroups:         true,
	}

	return &targetConfig
}

// logTargetConfig will print out the target configuration in the log (debug level).
// It will censure the sensitive information (secrets and passwords) if it is set.
func logTargetConfig(config *BaseTargetConfig) {
	if !hclog.L().IsDebug() {
		return
	}
	cc := BaseTargetConfig{}

	err := copier.Copy(&cc, config)
	if err != nil {
		hclog.L().Error("Error while copying config")
		return
	}
	cc.Parameters = make(map[string]interface{})

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

func fillStruct(o interface{}, m map[string]interface{}) error {
	for k, v := range m {
		err := setField(o, k, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func setField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		structFieldValue = structValue.FieldByName(toCamelInitCase(name, true))
		if !structFieldValue.IsValid() {
			// Not returning an error but just skipping = ignoring unknown fields.
			return nil
		}
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("cannot set value of field %q field", name)
	}

	structFieldType := structFieldValue.Type()

	value, err := iconfig.HandleField(value, structFieldType.Kind())
	if err != nil {
		return err
	}

	val := reflect.ValueOf(value)

	if structFieldType != val.Type() {
		return fmt.Errorf("provided value type didn't match obj field type for %q", name)
	}

	structFieldValue.Set(val)

	return nil
}

// Converts a string to CamelCase
func toCamelInitCase(s string, initCase bool) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	n := strings.Builder{}
	n.Grow(len(s))
	capNext := initCase

	for i, v := range []byte(s) {
		vIsCap := v >= 'A' && v <= 'Z'
		vIsLow := v >= 'a' && v <= 'z'

		if capNext {
			if vIsLow {
				v += 'A'
				v -= 'a'
			}
		} else if i == 0 {
			if vIsCap {
				v += 'a'
				v -= 'A'
			}
		}

		if vIsCap || vIsLow {
			n.WriteByte(v)
			capNext = false
		} else if vIsNum := v >= '0' && v <= '9'; vIsNum {
			n.WriteByte(v)
			capNext = true
		} else {
			capNext = v == '_' || v == ' ' || v == '-' || v == '.'
		}
	}

	return n.String()
}

type Options struct {
	DataSourceIds map[string]struct{}
	ConfigOption  func(targetConfig *BaseTargetConfig)
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

func (o *Options) TargetOptions(targetConfig *BaseTargetConfig) *BaseTargetConfig {
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

func WithConfigOption(fn func(targetConfig *BaseTargetConfig)) func(o *Options) {
	return func(o *Options) {
		o.ConfigOption = fn
	}
}
