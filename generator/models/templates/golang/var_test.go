package golang

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVar_MustParse(t *testing.T) {
	input := Var{
		Name:  "example",
		Value: `"example value"`,
	}
	expected := `var example = "example value"`
	result := input.MustParse()
	assert.Equal(t, expected, result)
}
