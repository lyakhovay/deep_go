package main

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize(data any) string {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return ""
	}
	result := make([]string, 0, v.NumField())
	for i := range v.NumField() {
		tag, ok := v.Type().Field(i).Tag.Lookup("properties")
		if !ok {
			continue
		}
		fieldToWrite := fieldToString(tag, v.Field(i))
		if len(fieldToWrite) > 0 {
			result = append(result, fieldToWrite)
		}
	}
	return strings.Join(result, "\n")
}

func getPropertyTag(tag string) (string, bool) {
	var omitempty bool
	tagParts := strings.Split(tag, ",")
	newTag := make([]string, 0, len(tagParts))
	for _, tagPart := range tagParts {
		if tagPart == "omitempty" {
			omitempty = true
			continue
		}
		newTag = append(newTag, tagPart)
	}
	return strings.Join(newTag, ","), omitempty
}

func fieldToString(tag string, field reflect.Value) string {
	newTag, omitempty := getPropertyTag(tag)
	if omitempty && field.IsZero() {
		return ""
	}
	str := strings.Builder{}
	str.WriteString(newTag)
	str.WriteString("=")
	switch field.Kind() {
	case reflect.String:
		str.WriteString(field.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		str.WriteString(strconv.FormatInt(field.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		str.WriteString(strconv.FormatUint(field.Uint(), 10))
	case reflect.Float32, reflect.Float64:
		str.WriteString(strconv.FormatFloat(field.Float(), 'g', -1, 64))
	case reflect.Bool:
		str.WriteString(strconv.FormatBool(field.Bool()))
	default:
		panic("unsupported field type")
	}
	return str.String()
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
