package golang

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVar_MustParse(t *testing.T) {
	input := Var{
		Name:  "example",
		Value: `"example value,omitempty"`,
	}
	expected := `var example = "example value,omitempty"`
	result := input.MustParse()
	assert.Equal(t, expected, result)
}
