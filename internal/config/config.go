package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// HandleField will check if the given value is a string in the form of '{{xxx}}' where 'xxx' is the name of an environment variable.
// When that format is found, the environment variable will be fetched and the value converted into the given kind (string, int, float64 or bool)
// If the environment variable doesn't exist or the value cannot be converted, an error is returned.
// Otherwise the original value is returned (if not a string or not in the '{{xxx}}' format)
func HandleField(value interface{}, targetType reflect.Kind) (interface{}, error) {
	if sv, ok := value.(string); ok {
		sv = strings.TrimSpace(sv)
		if strings.HasPrefix(sv, "{{") && strings.HasSuffix(sv, "}}") {
			envVar := sv[2:len(sv)-2]
			envValue, envSet := os.LookupEnv(envVar)
			if !envSet {
				return nil, errors.New("no environment variable with name "+envVar+" found")
			}
			converted, err := convertStringToType(envValue, targetType)
			if err != nil {
				return nil, fmt.Errorf("error when converting environment variable %q: %s", envVar, err.Error())
			}
			return converted, nil
		}
	}

	return value, nil
}

func convertStringToType(v string, k reflect.Kind) (interface{}, error) {
	switch k {
	case reflect.String:
		return v, nil
	case reflect.Int:
		return strconv.Atoi(v)
	case reflect.Float64:
		return strconv.ParseFloat(v, 64)
	case reflect.Bool:
		return strconv.ParseBool(v)
	}
	return nil, fmt.Errorf("Unknown type %q to convert", k.String())
}
