package target

import (
	"fmt"
	"github.com/aws/smithy-go/ptr"
	iconfig "github.com/raito-io/cli/internal/config"
	"github.com/raito-io/cli/internal/constants"
	"reflect"
	"strconv"
	"strings"
)

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

func argumentToString(arg interface{}) *string {
	if arg == nil {
		return nil
	}

	switch v := arg.(type) {
	case string:
		return &v
	case int:
		return ptr.String(strconv.Itoa(v))
	case bool:
		return ptr.String(strconv.FormatBool(v))
	}

	return nil
}
