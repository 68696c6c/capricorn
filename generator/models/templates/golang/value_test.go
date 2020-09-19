package golang

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValue_MustParse(t *testing.T) {
	input := Value{
		Name: "example",
		Type: "string",
	}
	expected := "example string"
	result := input.MustParse()
	assert.Equal(t, expected, result)
}

func TestValue_getJoinedValueString_Single(t *testing.T) {
	input := []Value{
		{
			Name: "value",
			Type: "int",
		},
	}
	expected := "value int"
	result := getJoinedValueString(input)
	assert.Equal(t, expected, result)
}

func TestValue_getJoinedValueString_Multiple(t *testing.T) {
	input := []Value{
		{
			Name: "arg1",
			Type: "*types.TypeData",
		},
		{
			Name: "arg2",
			Type: "string",
		},
	}
	expected := "arg1 *types.TypeData, arg2 string"
	result := getJoinedValueString(input)
	assert.Equal(t, expected, result)
}
