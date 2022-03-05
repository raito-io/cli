package target

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/jinzhu/copier"
	config2 "github.com/raito-io/cli/common/util/config"
	"github.com/raito-io/cli/internal/constants"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestToCamelCase(t *testing.T) {
	assert.Equal(t, "helloString", toCamelInitCase("hello-string", false))
	assert.Equal(t, "HelloString", toCamelInitCase("hello-string", true))
	assert.Equal(t, "", toCamelInitCase("", true))
	assert.Equal(t, "Something", toCamelInitCase("something", true))
	assert.Equal(t, "something", toCamelInitCase("something", false))
	assert.Equal(t, "SomethingCrazy", toCamelInitCase("something crazy", true))
	assert.Equal(t, "upperStart", toCamelInitCase("UpperStart", false))
	assert.Equal(t, "aLLCAPS", toCamelInitCase("ALLCAPS", false))
}

type fillerStruct struct {
	FieldName string
	AnotherOne int
	Another2 float64
	Ok bool
	cannotSet string
}

func TestFillStruct(t *testing.T) {
	fs := fillerStruct{}
	err := fillStruct(&fs, map[interface{}]interface{}{
		"field-name": "blah",
		"another one": 666,
		"another 2": 5.5,
		"Ok": true,
		"unknown": "skip",
	})

	assert.Nil(t, err)
	assert.Equal(t, "blah", fs.FieldName)
	assert.Equal(t, 666, fs.AnotherOne)
	assert.Equal(t, true, fs.Ok)
	assert.Equal(t, 5.5, fs.Another2)
}

func TestFillStructWithEnvironmentVariables(t *testing.T) {
	os.Setenv("RAITO_TEST_FIELDNAME", "eblah")
	os.Setenv("RAITO_TEST_ANOTHERONE", "777")
	os.Setenv("RAITO_TEST_ANOTHERTWO", "6.6")
	os.Setenv("RAITO_TEST_OK", "true")

	fs := fillerStruct{}
	err := fillStruct(&fs, map[interface{}]interface{}{
		"field-name": "{{RAITO_TEST_FIELDNAME}}",
		"another one": "{{RAITO_TEST_ANOTHERONE}}",
		"another 2": "{{RAITO_TEST_ANOTHERTWO}}",
		"Ok": "{{RAITO_TEST_OK}}",
		"unknown": "skip",
	})

	assert.Nil(t, err)
	assert.Equal(t, "eblah", fs.FieldName)
	assert.Equal(t, 777, fs.AnotherOne)
	assert.Equal(t, true, fs.Ok)
	assert.Equal(t, 6.6, fs.Another2)
}

func TestFillStructWithEnvironmentVariablesNotSet(t *testing.T) {
	os.Setenv("RAITO_TEST_FIELDNAME", "eblah")
	os.Setenv("RAITO_TEST_ANOTHERONE", "777")
	os.Unsetenv("RAITO_TEST_ANOTHERTWO")
	os.Setenv("RAITO_TEST_OK", "true")

	fs := fillerStruct{}
	err := fillStruct(&fs, map[interface{}]interface{}{
		"field-name": "{{RAITO_TEST_FIELDNAME}}",
		"another one": "{{RAITO_TEST_ANOTHERONE}}",
		"another 2": "{{RAITO_TEST_ANOTHERTWO}}",
		"Ok": "{{RAITO_TEST_OK}}",
		"unknown": "skip",
	})

	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "RAITO_TEST_ANOTHERTWO"))
}

func TestFillStructWithEnvironmentVariablesWrongType(t *testing.T) {
	os.Setenv("RAITO_TEST_FIELDNAME", "eblah")
	os.Setenv("RAITO_TEST_ANOTHERONE", "xxx")
	os.Unsetenv("RAITO_TEST_ANOTHERTWO")
	os.Setenv("RAITO_TEST_OK", "true")

	fs := fillerStruct{}
	err := fillStruct(&fs, map[interface{}]interface{}{
		"field-name": "{{RAITO_TEST_FIELDNAME}}",
		"another one": "{{RAITO_TEST_ANOTHERONE}}",
		"another 2": "{{RAITO_TEST_ANOTHERTWO}}",
		"Ok": "{{RAITO_TEST_OK}}",
		"unknown": "skip",
	})

	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "RAITO_TEST_ANOTHERONE"))
}

func TestFillStructWrongType(t *testing.T) {
	fs := fillerStruct{}
	err := fillStruct(&fs, map[interface{}]interface{}{
		"field-name": 666,
	})

	assert.NotNil(t, err)
}

func TestFillStructCannotSet(t *testing.T) {
	fs := fillerStruct{}
	err := fillStruct(&fs, map[interface{}]interface{}{
		"cannotSet": "should return error",
	})

	assert.NotNil(t, err)
}

var targets1 = []interface{} {
	map[interface{}]interface{} {
		"name": "snowflake1",
		"connector": "snowflake",
	},
	map[interface{}]interface{} {
		"name": "okta1",
		"connector": "okta",
	},
	map[interface{}]interface{} {
		"name": "snowflake2",
		"connector": "snowflake",
	},
}

func TestBuildTargetConfigFromMapError(t *testing.T) {
	clearViper()
	data := map[interface{}]interface{}{
		"connector-name":                666,
	}
	config, err := buildTargetConfigFromMap(hclog.L(), data)
	assert.NotNil(t, err)
	assert.Nil(t, config)
}

var baseConfigMap = map[interface{}]interface{} {
	"connector-name": "c1",
	"connector-version": "0.1.0",
	"name": "cn1",
	"data-source-id": "xxx",
	"identity-store-id": "yyy",
	"api-user": "c1user",
	"api-secret": "<secret>",
	"domain": "my-raito-domain",
	"skip-identity-store-sync": true,
	"skip-data-source-sync": false,
	"skip-data-access-sync": true,
	"custom1": "v1",
	"custom2": 5,
	"custom3": true,
}

func TestBuildTargetConfigFromMap(t *testing.T) {
	clearViper()
	config, err := buildTargetConfigFromMap(hclog.L(), baseConfigMap)
	assert.Nil(t, err)

	assert.Equal(t, "c1", config.ConnectorName)
	assert.Equal(t, "0.1.0", config.ConnectorVersion)
	assert.Equal(t, "cn1", config.Name)
	assert.Equal(t, "xxx", config.DataSourceId)
	assert.Equal(t, "yyy", config.IdentityStoreId)
	assert.Equal(t, "c1user", config.ApiUser)
	assert.Equal(t, "<secret>", config.ApiSecret)
	assert.Equal(t, "my-raito-domain", config.Domain)
	assert.Equal(t, true, config.SkipIdentityStoreSync)
	assert.Equal(t, false, config.SkipDataSourceSync)
	assert.Equal(t, true, config.SkipDataAccessSync)
	assert.Equal(t, 3, len(config.ConfigMap.Parameters))
	assert.Equal(t, "v1", config.ConfigMap.GetString("custom1"))
	assert.Equal(t, 5, config.ConfigMap.GetInt("custom2"))
	assert.Equal(t, true, config.ConfigMap.GetBoolWithDefault("custom3", false))
}

func TestBuildTargetConfigFromMapNoName(t *testing.T) {
	clearViper()
	var noNameConfigMap = make(map[interface{}]interface{})
	copier.Copy(&noNameConfigMap, &baseConfigMap)
	delete(noNameConfigMap, "name")
	config, err := buildTargetConfigFromMap(hclog.L(), noNameConfigMap)
	assert.Nil(t, err)

	assert.Equal(t, "c1", config.ConnectorName)
	assert.Equal(t, "c1", config.Name)
}

func TestBuildTargetConfigFromMapOverride(t *testing.T) {
	clearViper()
	viper.Set("skip-data-source-sync", true)
	config, err := buildTargetConfigFromMap(hclog.L(), baseConfigMap)
	assert.Nil(t, err)

	assert.Equal(t, true, config.SkipIdentityStoreSync)
	assert.Equal(t, true, config.SkipDataSourceSync)
	assert.Equal(t, true, config.SkipDataAccessSync)
}

func TestBuildTargetConfigFromMapLocalRaitoData(t *testing.T) {
	clearViper()
	viper.Set("api-user", "uuuu")
	viper.Set("domain", "dddd")
	viper.Set("api-secret", "ssss")
	config, err := buildTargetConfigFromMap(hclog.L(), baseConfigMap)
	assert.Nil(t, err)

	assert.Equal(t, "c1user", config.ApiUser)
	assert.Equal(t, "<secret>", config.ApiSecret)
	assert.Equal(t, "my-raito-domain", config.Domain)
}

func TestBuildTargetConfigFromMapGlobalRaitoData(t *testing.T) {
	clearViper()
	// Create the target map
	withoutRaitoStuff := make(map[interface{}]interface{})
	// Copy from the original map to the target map
	for key, value := range baseConfigMap {
		if key != "api-user" && key != "api-secret" && key != "domain" {
			withoutRaitoStuff[key] = value
		}
	}
	viper.Set("api-user", "uuuu")
	viper.Set("api-secret", "ssss")
	viper.Set("domain", "dddd")
	config, err := buildTargetConfigFromMap(hclog.L(), withoutRaitoStuff)
	assert.Nil(t, err)

	assert.Equal(t, "uuuu", config.ApiUser)
	assert.Equal(t, "ssss", config.ApiSecret)
	assert.Equal(t, "dddd", config.Domain)
}

func clearViper() {
	for _, key := range viper.AllKeys() {
		viper.Set(key, nil)
	}
}

func TestBuildTargetConfigFromFlags(t *testing.T) {
	clearViper()

	viper.Set(constants.ConnectorNameFlag, "conn1")
	viper.Set(constants.NameFlag, "name1")

	viper.Set("data-source-id", "aaa")
	viper.Set("identity-store-id", "eee")
	viper.Set("api-user", "conn1user")
	viper.Set("api-secret", "<secret>")
	viper.Set("domain", "my-raito-domain")
	viper.Set("skip-identity-store-sync", false)
	viper.Set("skip-data-source-sync", true)
	viper.Set("skip-data-access-sync", true)

	config, err := buildTargetConfigFromFlags(hclog.L(), []string {"--custom1", "ok"})
	assert.Nil(t, err)
	assert.NotNil(t, config)

	assert.Equal(t, "conn1", config.ConnectorName)
	assert.Equal(t, "name1", config.Name)
	assert.Equal(t, "aaa", config.DataSourceId)
	assert.Equal(t, "eee", config.IdentityStoreId)
	assert.Equal(t, "conn1user", config.ApiUser)
	assert.Equal(t, "<secret>", config.ApiSecret)
	assert.Equal(t, "my-raito-domain", config.Domain)
	assert.Equal(t, false, config.SkipIdentityStoreSync)
	assert.Equal(t, true, config.SkipDataSourceSync)
	assert.Equal(t, true, config.SkipDataAccessSync)
	assert.Equal(t, 1, len(config.ConfigMap.Parameters))
	assert.Equal(t, "ok", config.ConfigMap.GetString("custom1"))
}

func TestBuildTargetConfigFromFlagsNoName(t *testing.T) {
	clearViper()

	viper.Set(constants.ConnectorNameFlag, "conn1")

	config, err := buildTargetConfigFromFlags(hclog.L(), []string { })
	assert.Nil(t, err)
	assert.NotNil(t, config)

	assert.Equal(t, "conn1", config.ConnectorName)
	assert.Equal(t, "conn1", config.Name)
}

func TestBuildParameterMapFromArguments(t *testing.T) {
	params := buildParameterMapFromArguments([]string {"--bool-val", "--string-val=blah", "--another-one", "moremoremore" })
	assert.Equal(t, 3, len(params))
	assert.Equal(t, true, params["bool-val"])
	assert.Equal(t, "blah", params["string-val"])
	assert.Equal(t, "moremoremore", params["another-one"])
}

func TestRunSingleTarget(t *testing.T) {
	clearViper()

	viper.Set(constants.ConnectorNameFlag, "conn1")
	viper.Set(constants.NameFlag, "name1")

	viper.Set("data-source-id", "aaa")
	viper.Set("identity-store-id", "eee")
	viper.Set("api-user", "conn1user")
	viper.Set("api-secret", "<secret>")
	viper.Set("domain", "my-raito-domain")
	viper.Set("skip-identity-store-sync", false)
	viper.Set("skip-data-source-sync", true)
	viper.Set("skip-data-access-sync", true)

	runs := 0
	RunTargets(hclog.L(), []string {}, func(tConfig *BaseTargetConfig) error {
		assert.Equal(t, "name1", tConfig.Name)
		runs++
		return nil
	})
	assert.Equal(t, 1, runs)
}

func TestRunMultipleTargets(t *testing.T) {
	clearViper()

	t1 := map[interface{}]interface{} {
		constants.ConnectorNameFlag: "c1",
		constants.NameFlag:      "cn1",
		"api-secret":            "secret1",
		"other-stuff":           "ok",
	}
	t2 := map[interface{}]interface{} {
		constants.ConnectorNameFlag: "c2",
		"api-secret":            "secret2",
	}

	targets := []interface{} {
		t1, t2,

	}
	viper.Set("targets", targets)

	runs := 0
	RunTargets(hclog.L(), []string {}, func(tConfig *BaseTargetConfig) error {
		if runs == 0 {
			assert.Equal(t, "c1", tConfig.ConnectorName)
			assert.Equal(t, "cn1", tConfig.Name)
			assert.Equal(t, "secret1", tConfig.ApiSecret)
			assert.Equal(t, "ok", tConfig.Parameters["other-stuff"])
		} else if runs == 1 {
			assert.Equal(t, "c2", tConfig.ConnectorName)
			assert.Equal(t, "c2", tConfig.Name)
			assert.Equal(t, "secret2", tConfig.ApiSecret)
		}
		runs++
		return nil
	})
	assert.Equal(t, 2, runs)
}

func TestRunMultipleTargetsWithOnlyTargets(t *testing.T) {
	clearViper()

	t1 := map[interface{}]interface{} {
		constants.ConnectorNameFlag: "c1",
		constants.NameFlag:      "name1",
	}
	t2 := map[interface{}]interface{} {
		constants.ConnectorNameFlag: "c2",
	}

	targets := []interface{} {
		t1, t2,

	}
	viper.Set("targets", targets)

	viper.Set("only-targets", "name1")

	runs := 0
	RunTargets(hclog.L(), []string {}, func(tConfig *BaseTargetConfig) error {
		assert.Equal(t, "c1", tConfig.ConnectorName)
		assert.Equal(t, "name1", tConfig.Name)
		runs++
		return nil
	})
	assert.Equal(t, 1, runs)

	viper.Set("only-targets", "c2")

	runs = 0
	RunTargets(hclog.L(), []string {}, func(tConfig *BaseTargetConfig) error {
		assert.Equal(t, "c2", tConfig.ConnectorName)
		assert.Equal(t, "c2", tConfig.Name)
		runs++
		return nil
	})
	assert.Equal(t, 1, runs)
}

func TestLogTarget(t *testing.T) {
	hclog.L().SetLevel(hclog.Debug)
	var buf bytes.Buffer
	var sbuf bytes.Buffer

	intercept := hclog.NewInterceptLogger(&hclog.LoggerOptions{
		Level:  hclog.Debug,
		Output: &buf,
	})

	sink := hclog.NewSinkAdapter(&hclog.LoggerOptions{
		Level:  hclog.Debug,
		Output: &sbuf,
	})

	intercept.RegisterSink(sink)
	defer intercept.DeregisterSink(sink)
	old := hclog.L()
	hclog.SetDefault(intercept)

	config := BaseTargetConfig{
		ApiSecret: "mylittlesecret",
		ConfigMap: config2.ConfigMap{
			Parameters: map[string]interface{} {
				"password": "anothersecret",
				"normal": "readible",
			},
		},
		ApiUser: "theuser",
	}
	logTargetConfig(&config)

	str := sbuf.String()
	fmt.Println(str)
	assert.NotContains(t, str, "mylittlesecret")
	assert.NotContains(t, str, "anothersecret")
	assert.Contains(t, str, "readible")
	assert.Contains(t, str, "theuser")
	assert.Contains(t, str, "**censured**")

	hclog.SetDefault(old)
}