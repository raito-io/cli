package jsonstream

import (
	"errors"
	"io"
	"reflect"

	"github.com/bcicen/jstream"
)

type JsonArrayStreamResult[T any] struct {
	Result *T
	Err    error
}

type JsonArrayStream[T any] struct {
	decoder *jstream.Decoder

	outputChan chan *JsonArrayStreamResult[T]
}

func NewJsonArrayStream[T any](r io.Reader) *JsonArrayStream[T] {
	return &JsonArrayStream[T]{
		decoder:    jstream.NewDecoder(r, 1),
		outputChan: make(chan *JsonArrayStreamResult[T]),
	}
}
func (s *JsonArrayStream[T]) Stream() chan *JsonArrayStreamResult[T] {
	go s.transformStream()
	return s.outputChan
}

func (s *JsonArrayStream[T]) transformStream() {
	defer close(s.outputChan)

	for v := range s.decoder.Stream() {
		t, err := s.transform(v)
		result := &JsonArrayStreamResult[T]{
			Result: t,
			Err:    err,
		}
		s.outputChan <- result
	}
}

func (s *JsonArrayStream[T]) transform(value *jstream.MetaValue) (*T, error) {
	valueMap := value.Value.(map[string]interface{})
	var outputObject T

	err := parseObject(valueMap, &outputObject)
	if err != nil {
		return nil, err
	}

	return &outputObject, nil
}

func parseObject(values map[string]interface{}, object interface{}) error {
	rv := reflect.ValueOf(object)
	rt := rv.Elem().Type()

	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return errors.New("invalid json parsing")
	}

	for i := 0; i < rt.NumField(); i++ {
		fieldDefinition := rt.Field(i)
		jsonName, jsonNameSpecified := fieldDefinition.Tag.Lookup("json")

		if !jsonNameSpecified || jsonName == "-" {
			continue
		}

		if value, found := values[jsonName]; found {
			fieldValue := rv.Elem().FieldByName(fieldDefinition.Name)

			err := parseValue(&fieldValue, value)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func parseValue(fieldValue *reflect.Value, value interface{}) error {
	if value == nil {
		return nil
	}

	fieldKind := fieldValue.Kind()

	switch fieldKind {
	case reflect.Ptr:
		pointerObject := reflect.New(fieldValue.Type().Elem())
		pointeeObject := pointerObject.Elem()

		err := parseValue(&pointeeObject, value)
		if err != nil {
			return err
		}

		fieldValue.Set(pointerObject)
	case reflect.Struct:
		innerValues := value.(map[string]interface{})
		innerObject := reflect.New(fieldValue.Type()).Interface()

		err := parseObject(innerValues, innerObject)
		if err != nil {
			return err
		}

		fieldValue.Set(reflect.ValueOf(innerObject).Elem())
	case reflect.Slice:
		innerValues := value.([]interface{})
		innerObjects := reflect.New(fieldValue.Type()).Elem()

		for _, innerValue := range innerValues {
			innerObject := reflect.New(fieldValue.Type().Elem())
			innerObjectElem := innerObject.Elem()

			err := parseValue(&innerObjectElem, innerValue)
			if err != nil {
				return err
			}

			innerObjects = reflect.Append(innerObjects, innerObject.Elem())
		}

		fieldValue.Set(innerObjects)
	default:
		fieldValue.Set(reflect.ValueOf(value).Convert(fieldValue.Type()))
	}

	return nil
}
