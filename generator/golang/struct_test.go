package golang

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeStructTestField(name, fieldType string) (*Field, string) {
	tag, _ := makeStructFieldTestTag("tag")
	f := Field{
		Name: name,
		Type: fieldType,
		Tags: []*Tag{tag},
	}
	expected := f.MustString()
	return &f, expected
}

func TestStruct_MustParse(t *testing.T) {
	field1, expectedField1 := makeStructTestField("ExportedField", "string")
	field2, expectedField2 := makeStructTestField("unexportedField", "int")
	input := Struct{
		Name:   "ExampleStruct",
		Fields: []*Field{field1, field2},
	}
	expected := fmt.Sprintf(`type ExampleStruct struct {
	%s
	%s
}`, expectedField1, expectedField2)
	result := input.MustString()
	assert.Equal(t, expected, result)
}
