package config

import (
	"encoding/json"
	"strconv"
)

// GetBool returns the boolean value with the given key in the Parameters map. If not found, or the value cannot be converted to a boolean, false is returned.
func (c *ConfigMap) GetBool(key string) bool {
	return c.GetBoolWithDefault(key, false)
}

// GetBoolWithDefault returns the boolean value with the given key in the Parameters map. If not found, or the value cannot be converted to a boolean, the given default value is returned.
func (c *ConfigMap) GetBoolWithDefault(key string, defaultValue bool) bool {
	if v, ok := c.Parameters[key]; ok {
		ret, err := strconv.ParseBool(v)
		if err != nil {
			return defaultValue
		}

		return ret
	}

	return defaultValue
}

// GetString returns the string value with the given key in the Parameters map. If not found, an empty string is returned.
func (c *ConfigMap) GetString(key string) string {
	return c.GetStringWithDefault(key, "")
}

// GetStringWithDefault returns the string value with the given key in the Parameters map. If not found, the given default value is returned.
func (c *ConfigMap) GetStringWithDefault(key string, defaultValue string) string {
	config := c.Parameters
	if v, ok := config[key]; ok {
		return v
	}

	return defaultValue
}

// GetInt returns the integer value with the given key in the Parameters map. If not found, or it cannot be converted into an integer, 0 is returned.
func (c *ConfigMap) GetInt(key string) int {
	return c.GetIntWithDefault(key, 0)
}

// GetIntWithDefault returns the integer value with the given key in the Parameters map. If not found, or it cannot be converted into an integer, the given default value is returned.
func (c *ConfigMap) GetIntWithDefault(key string, defaultValue int) int {
	config := c.Parameters
	if v, ok := config[key]; ok {
		ret, err := strconv.Atoi(v)
		if err != nil {
			return defaultValue
		}

		return ret
	}

	return defaultValue
}

func (c *ConfigMap) Unmarshal(key string, value interface{}) (bool, error) {
	config := c.Parameters
	if v, ok := config[key]; ok {
		return true, json.Unmarshal([]byte(v), value)
	}

	return false, nil
}
