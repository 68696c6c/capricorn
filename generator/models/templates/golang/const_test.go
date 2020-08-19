package golang

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConst_MustParse(t *testing.T) {
	input := Const{
		Name:  "example",
		Value: `"example value"`,
	}
	expected := `const example = "example value"`
	result := input.MustParse()
	assert.Equal(t, expected, result)
}
