package golang

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeStructFieldTestTag(name string) (Tag, string) {
	tag := Tag{
		Key:    name,
		Values: []string{"value1", "value2"},
	}
	expected := fmt.Sprintf(`%s:"value1,value2"`, name)
	return tag, expected
}

func TestField_MustParse(t *testing.T) {
	tag1, expectedTag1 := makeStructFieldTestTag("tag1")
	tag2, expectedTag2 := makeStructFieldTestTag("tag2")
	input := Field{
		Name: "ExampleField",
		Type: "string",
		Tags: []Tag{tag1, tag2},
	}
	expected := fmt.Sprintf("ExampleField string `%s %s`", expectedTag1, expectedTag2)
	result := input.MustParse()
	assert.Equal(t, expected, result)
}
