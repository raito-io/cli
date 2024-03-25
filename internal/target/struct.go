package target

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/aws/smithy-go/ptr"

	iconfig "github.com/raito-io/cli/internal/config"
	"github.com/raito-io/cli/internal/constants"
)

// fillStruct will fill the given struct object in the first parameter with the values from the map in the second parameter.
// It will do this automatically by finding the fields in the struct and matching them with the keys in the map.
// This is done by transforming the key from the map to a camel-case to match the field name in the struct (e.g. api-user becomes ApiUser)
func fillStruct(o interface{}, m map[string]interface{}) error {
	for k, v := range m {
		if k != constants.DataObjectEnrichers {
			err := setField(o, k, v)
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

func argumentToString(arg interface{}) (*string, error) {
	if arg == nil {
		return nil, nil
	}

	switch v := arg.(type) {
	case string:
		return &v, nil
	case int:
		return ptr.String(strconv.Itoa(v)), nil
	case bool:
		return ptr.String(strconv.FormatBool(v)), nil
	default:
		jsonb, err := json.Marshal(arg)
		if err != nil {
			return nil, err
		}

		return ptr.String(string(jsonb)), nil
	}
}
