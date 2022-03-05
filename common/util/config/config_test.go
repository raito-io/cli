package config


import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetBool(t *testing.T) {
	c := ConfigMap{
		Parameters: map[string]interface{} {
			"bool-ok": "true",
			"bool-ok-1": "1",
			"bool-nok": "nok",
			"bool-false": "false",
			"bool-real": true,
		},
	}
	assert.Equal(t, true, c.GetBoolWithDefault("bool-real", false))
	assert.Equal(t, true, c.GetBool("bool-real"))
	assert.Equal(t, true, c.GetBoolWithDefault("bool-ok", false))
	assert.Equal(t, true, c.GetBool("bool-ok"))
	assert.Equal(t, true, c.GetBoolWithDefault("bool-ok-1", false))
	assert.Equal(t, true, c.GetBool("bool-ok-1"))
	assert.Equal(t, false, c.GetBoolWithDefault("bool-nok", false))
	assert.Equal(t, true, c.GetBoolWithDefault("bool-nok", true))
	assert.Equal(t, false, c.GetBool("bool-nok"))
	assert.Equal(t, false, c.GetBoolWithDefault("bool-false", false))
	assert.Equal(t, false, c.GetBoolWithDefault("bool-false", true))
	assert.Equal(t, false, c.GetBool("bool-false"))
	assert.Equal(t, true, c.GetBoolWithDefault("not-exists", true))
	assert.Equal(t, false, c.GetBool("not-exists"))
}

func TestGetString(t *testing.T) {
	c := ConfigMap{
		Parameters: map[string]interface{} {
			"string-ok": "some string",
			"string-empty": "",
		},
	}
	assert.Equal(t, "some string", c.GetStringWithDefault("string-ok", "default"))
	assert.Equal(t, "some string", c.GetString("string-ok"))
	assert.Equal(t, "", c.GetStringWithDefault("string-empty", "default"))
	assert.Equal(t, "", c.GetString("string-empty"))
	assert.Equal(t, "default", c.GetStringWithDefault("string-notexists", "default"))
	assert.Equal(t, "", c.GetString("string-notexists"))
}

func TestGetInt(t *testing.T) {
	c := ConfigMap{
		Parameters: map[string]interface{} {
			"int-ok": "123",
			"int-empty": "",
			"int-wrong": "wrong",
			"int-real": 555,
			"int64": int64(777),
		},
	}
	assert.Equal(t, 123, c.GetIntWithDefault("int-ok", 666))
	assert.Equal(t, 123, c.GetInt("int-ok"))
	assert.Equal(t, 666, c.GetIntWithDefault("int-empty", 666))
	assert.Equal(t, 0, c.GetInt("int-empty"))
	assert.Equal(t, 666, c.GetIntWithDefault("int-wrong", 666))
	assert.Equal(t, 0, c.GetInt("int-wrong"))

	assert.Equal(t, 555, c.GetIntWithDefault("int-real", 666))
	assert.Equal(t, 555, c.GetInt("int-real"))
	assert.Equal(t, 777, c.GetIntWithDefault("int64", 666))
	assert.Equal(t, 777, c.GetInt("int64"))
}

