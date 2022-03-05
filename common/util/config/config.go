package config

import "strconv"

type ConfigMap struct {
	Parameters map[string]interface{}
}

// GetBool returns the boolean value with the given key in the Parameters map. If not found, or the value cannot be converted to a boolean, false is returned.
func (c *ConfigMap) GetBool(key string) bool {
	return c.GetBoolWithDefault(key, false)
}

// GetBoolWithDefault returns the boolean value with the given key in the Parameters map. If not found, or the value cannot be converted to a boolean, the given default value is returned.
func (c *ConfigMap) GetBoolWithDefault(key string, defaultValue bool) bool {
	if v, ok := c.Parameters[key]; ok && v != nil {
		if sv, ok := v.(string); ok {
			ret, err := strconv.ParseBool(sv)
			if err != nil {
				return defaultValue
			}
			return ret
		}
		if bv, ok := v.(bool); ok {
			return bv
		}
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
	if v, ok := config[key]; ok && v != nil {
		if sv, ok := v.(string); ok {
			return sv
		}
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
	if v, ok := config[key]; ok && v != nil {
		if sv, ok := v.(string); ok {
			ret, err := strconv.Atoi(sv)
			if err != nil {
				return defaultValue
			}
			return ret
		}
		if bv, ok := v.(int); ok {
			return bv
		}
		if bv, ok := v.(int64); ok {
			return int(bv)
		}
	}
	return defaultValue
}

