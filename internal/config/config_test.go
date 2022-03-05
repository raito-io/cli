package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"testing"
)

func TestHandleFieldIntNoEnv(t *testing.T) {
	v1 := 55
	r1, err := HandleField(v1, reflect.String)
	assert.Nil(t, err)
	assert.Equal(t, 55, r1)
}

func TestHandleFieldStringNoEnv(t *testing.T) {
	v1 := "ok"
	r1, err := HandleField(v1, reflect.String)
	assert.Nil(t, err)
	assert.Equal(t, "ok", r1)
}

func TestHandleFieldStringWrongEnv(t *testing.T) {
	v1 := "{{notexisting}}"
	_, err := HandleField(v1, reflect.String)
	assert.NotNil(t, err)
}

func TestHandleFieldStringEnv(t *testing.T) {
	envVar := "RAITO_TEST_KEY1"
	os.Setenv(envVar, "value")
	v1 := "{{" + envVar + "}}"
	r1, err := HandleField(v1, reflect.String)
	assert.Nil(t, err)
	assert.Equal(t, "value", r1)
}

func TestHandleFieldIntEnv(t *testing.T) {
	envVar := "RAITO_TEST_KEY2"
	os.Setenv(envVar, "5")
	v1 := "{{" + envVar + "}}"
	r1, err := HandleField(v1, reflect.Int)
	assert.Nil(t, err)
	assert.Equal(t, 5, r1)
}

func TestHandleFieldFloatEnv(t *testing.T) {
	envVar := "RAITO_TEST_KEY3"
	os.Setenv(envVar, "5.67")
	v1 := "{{" + envVar + "}}"
	r1, err := HandleField(v1, reflect.Float64)
	assert.Nil(t, err)
	assert.Equal(t, 5.67, r1)
}

func TestHandleFieldBoolEnv(t *testing.T) {
	envVar := "RAITO_TEST_KEY4"
	os.Setenv(envVar, "true")
	v1 := "{{" + envVar + "}}"
	r1, err := HandleField(v1, reflect.Bool)
	assert.Nil(t, err)
	assert.Equal(t, true, r1)
}

func TestHandleFieldEnvBadConversion(t *testing.T) {
	envVar := "RAITO_TEST_KEY5"
	os.Setenv(envVar, "jos")
	v1 := "{{" + envVar + "}}"
	_, err := HandleField(v1, reflect.Float64)
	assert.NotNil(t, err)
}
