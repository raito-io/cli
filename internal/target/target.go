package target

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/jinzhu/copier"
	"github.com/raito-io/cli/common/api"
	"github.com/raito-io/cli/common/util/config"
	iconfig "github.com/raito-io/cli/internal/config"
	"github.com/raito-io/cli/internal/constants"
	"github.com/spf13/viper"
	"reflect"
	"strings"
)

type BaseTargetConfig struct {
	config.ConfigMap
	ConnectorName         string
	ConnectorVersion      string
	Name                  string
	DataSourceId          string
	IdentityStoreId       string
	ApiUser               string
	ApiSecret             string
	Domain                string
	SkipIdentityStoreSync bool
	SkipDataSourceSync    bool
	SkipDataAccessSync    bool
	SkipDataUsageSync     bool
	Logger                hclog.Logger

	DeleteUntouched bool
	ReplaceTags     bool
	DeleteTempFiles bool
	ReplaceGroups   bool
}

func RunTargets(baseLogger hclog.Logger, otherArgs []string, runTarget func(tConfig *BaseTargetConfig) error) error {
	if viper.GetString(constants.ConnectorNameFlag) != "" {
		targetConfig, _ := buildTargetConfigFromFlags(baseLogger, otherArgs)
		logTargetConfig(targetConfig)
		return runTarget(targetConfig)
	} else {
		return runMultipleTargets(baseLogger, runTarget)
	}
}

func HandleTargetError(err error, config *BaseTargetConfig, during string) {
	if errorResult, ok := err.(api.ErrorResult); ok {
		if errorResult.ErrorCode == api.BadInputParameterError || errorResult.ErrorCode == api.MissingInputParameterError {
			config.Logger.Error(fmt.Sprintf("Error during %s: %s. Execute command 'info <connector>' to print out the expected parameters for the connector.", during, errorResult.ErrorMessage))
			return
		}
	}
	config.Logger.Error(fmt.Sprintf("Error during %s: %s", during, err.Error()))
}

func runMultipleTargets(baseLogger hclog.Logger, runTarget func(tConfig *BaseTargetConfig) error) error {
	targets := viper.Get(constants.Targets)

	onlyTargetsS := viper.GetString(constants.OnlyTargetsFlag)
	onlyTargets := make(map[string]struct{})
	if onlyTargetsS != "" {
		for _, ot := range strings.Split(onlyTargetsS, ",") {
			onlyTargets[strings.TrimSpace(ot)] = struct{}{}
		}
	}

	if targetList, ok := targets.([]interface{}); ok {
		hclog.L().Debug(fmt.Sprintf("Found %d targets to run.", len(targetList)))
		for _, targetObj := range targetList {
			if target, ok := targetObj.(map[interface{}]interface{}); ok {
				tConfig, err := buildTargetConfigFromMap(baseLogger, target)
				if err != nil {
					hclog.L().Error(fmt.Sprintf("error while parsing target configuration: %s", err.Error()))
					continue
				}
				if tConfig == nil {
					continue
				}
				if len(onlyTargets) > 0 {
					if _, found := onlyTargets[tConfig.Name]; !found {
						//time.Sleep(3*time.Second)
						tConfig.Logger.Info("Skipping target", "success")
						continue
					}
				}

				logTargetConfig(tConfig)
				err = runTarget(tConfig)
				if err != nil {
					// In debug as the error should already be outputted, and we are ignoring it here.
					tConfig.Logger.Debug("Error while executing target", "error", err.Error())
				}
			}
		}
	}

	return nil
}

func buildTargetConfigFromMap(baseLogger hclog.Logger, target map[interface{}]interface{}) (*BaseTargetConfig, error) {
	tConfig := BaseTargetConfig{}
	err := fillStruct(&tConfig, target)
	if err != nil {
		return nil, err
	}
	tConfig.Parameters = make(map[string]interface{})
	for k, v := range target {
		if ks, ok := k.(string); ok {
			if _, f := constants.KnownFlags[ks]; !f {
				cv, err := iconfig.HandleField(v, reflect.String)
				if err != nil {
					return nil, err
				}
				tConfig.Parameters[ks] = cv
			}
		}
	}
	if tConfig.Name == "" {
		tConfig.Name = tConfig.ConnectorName
	}

	// Create a logger to add the target log name to each log message.
	tConfig.Logger = baseLogger.With("target", tConfig.Name)

	// Merge with some global parameters
	tConfig.SkipDataAccessSync = tConfig.SkipDataAccessSync || viper.GetBool(constants.SkipDataAccessSyncFlag)
	tConfig.SkipDataSourceSync = tConfig.SkipDataSourceSync || viper.GetBool(constants.SkipDataSourceSyncFlag)
	tConfig.SkipIdentityStoreSync = tConfig.SkipIdentityStoreSync || viper.GetBool(constants.SkipIdentityStoreSyncFlag)
	tConfig.SkipDataUsageSync = tConfig.SkipDataUsageSync || viper.GetBool(constants.SkipDataUsageSyncFlag)

	// Merge with import parameters
	tConfig.DeleteUntouched = tConfig.DeleteUntouched || viper.GetBool(constants.DeleteUntouchedFlag)
	tConfig.ReplaceTags = tConfig.ReplaceTags || viper.GetBool(constants.ReplaceTagsFlag)
	tConfig.DeleteTempFiles = tConfig.DeleteTempFiles || viper.GetBool(constants.DeleteTempFilesFlag)
	tConfig.ReplaceGroups = tConfig.ReplaceGroups || viper.GetBool(constants.ReplaceGroupsFlag)

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

func buildTargetConfigFromFlags(baseLogger hclog.Logger, otherArgs []string) (*BaseTargetConfig, error) {
	connector := viper.GetString(constants.ConnectorNameFlag)
	version := viper.GetString(constants.ConnectorVersionFlag)
	name := viper.GetString(constants.NameFlag)
	if name == "" {
		name = connector
	}

	targetConfig := BaseTargetConfig{
		ConnectorName:         connector,
		ConnectorVersion:      version,
		Name:                  name,
		DataSourceId:          viper.GetString(constants.DataSourceIdFlag),
		IdentityStoreId:       viper.GetString(constants.IdentityStoreIdFlag),
		ApiUser:               viper.GetString(constants.ApiUserFlag),
		ApiSecret:             viper.GetString(constants.ApiSecretFlag),
		Domain:                viper.GetString(constants.DomainFlag),
		SkipIdentityStoreSync: viper.GetBool(constants.SkipIdentityStoreSyncFlag),
		SkipDataSourceSync:    viper.GetBool(constants.SkipDataSourceSyncFlag),
		SkipDataAccessSync:    viper.GetBool(constants.SkipDataAccessSyncFlag),
		SkipDataUsageSync:     viper.GetBool(constants.SkipDataUsageSyncFlag),
		Logger:                baseLogger.With("target", name),
		DeleteUntouched:       viper.GetBool(constants.DeleteUntouchedFlag),
		DeleteTempFiles:       viper.GetBool(constants.DeleteTempFilesFlag),
		ReplaceTags:           viper.GetBool(constants.ReplaceTagsFlag),
		ReplaceGroups:         viper.GetBool(constants.ReplaceGroupsFlag),
	}
	targetConfig.Parameters = buildParameterMapFromArguments(otherArgs)
	return &targetConfig, nil
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

func fillStruct(o interface{}, m map[interface{}]interface{}) error {
	for k, v := range m {
		if ks, ok := k.(string); ok {
			err := setField(o, ks, v)
			if err != nil {
				return err
			}
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
		return fmt.Errorf("Cannot set value of field %q field", name)
	}

	structFieldType := structFieldValue.Type()

	value, err := iconfig.HandleField(value, structFieldType.Kind())
	if err != nil {
		return err
	}

	val := reflect.ValueOf(value)

	if structFieldType != val.Type() {
		return fmt.Errorf("Provided value type didn't match obj field type for %q", name)
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
