package golang

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStructFieldTag_MustParse(t *testing.T) {
	input := Tag{
		Key:    "example",
		Values: []string{"value1", "value2"},
	}
	expected := `example:"value1,value2"`
	result := input.MustParse()
	assert.Equal(t, expected, result)
}
