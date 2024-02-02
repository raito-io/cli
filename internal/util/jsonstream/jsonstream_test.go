package jsonstream

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type internalType struct {
	InternalData string `json:"internalData"`
}

type objectType struct {
	OtherObject internalType `json:"otherObject"`
	SomeString  string       `json:"someString"`
	SomeInt     int          `json:"someInt"`
	ArrayObject []string     `json:"someArray"`
}

type objectTypeWithPointer struct {
	OtherObject  *internalType `json:"otherObject"`
	SomeString   *string       `json:"someString"`
	SomeInt      *int          `json:"someInt"`
	SomeNilSting *string       `json:"SomeNilSting"`
}

func TestJsonStream(t *testing.T) {
	// Given
	jsonObject := []objectType{
		{
			OtherObject: internalType{
				InternalData: "Bla",
			},
			SomeInt:     4,
			SomeString:  "Foo",
			ArrayObject: []string{"arrayObject1", "arrayObject2"},
		},
		{
			OtherObject: internalType{
				InternalData: "AnotherBla",
			},
			SomeInt:    7,
			SomeString: "AnotherFoo",
		},
	}

	rawJson, err := json.Marshal(jsonObject)
	assert.NoError(t, err)

	reader := bytes.NewReader(rawJson)

	jsonSteamer := NewJsonArrayStream[objectType](reader)

	objects := make([]objectType, 0)

	objectChannel := jsonSteamer.Stream()
	i := 0
	for object := range objectChannel {
		assert.NotNil(t, object)
		assert.NoError(t, object.Err)
		assert.NotNil(t, object.Result)
		i++
		objects = append(objects, *object.Result)
	}

	assert.Equal(t, 2, i)
	assert.Equal(t, jsonObject, objects)
}

func TestJsonStreamWithPointers(t *testing.T) {
	//Given
	ohterObject := internalType{
		InternalData: "Bla",
	}
	someString := "someString"
	someInt := 42

	jsonObject := objectTypeWithPointer{
		OtherObject:  &ohterObject,
		SomeString:   &someString,
		SomeInt:      &someInt,
		SomeNilSting: nil,
	}

	rawJson, err := json.Marshal([]objectTypeWithPointer{jsonObject})
	assert.NoError(t, err)

	reader := bytes.NewReader(rawJson)

	jsonSteamer := NewJsonArrayStream[objectTypeWithPointer](reader)

	objects := make([]objectTypeWithPointer, 0)

	objectChannel := jsonSteamer.Stream()
	i := 0
	for object := range objectChannel {
		assert.NotNil(t, object)
		assert.NoError(t, object.Err)
		assert.NotNil(t, object.Result)
		i++
		objects = append(objects, *object.Result)
	}

	assert.Equal(t, 1, i)
	assert.Equal(t, jsonObject, objects[0])

}
