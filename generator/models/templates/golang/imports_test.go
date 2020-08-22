package golang

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImports_MustParse(t *testing.T) {
	input := Imports{
		Standard: []string{"standard-one", "standard-two"},
		App:      []string{"app-one", "app-two"},
		Vendor:   []string{"vendor-one", "vendor-two"},
	}

	result := input.MustParse()
	expected := `import (
	"standard-one"
	"standard-two"

	"app-one"
	"app-two"

	"vendor-one"
	"vendor-two"
)`

	assert.Equal(t, expected, result)
}

func TestImports_MustParse_None(t *testing.T) {
	input := Imports{}

	result := input.MustParse()
	expected := ""

	assert.Equal(t, expected, result)
}

func TestImports_MustParse_NoStandard(t *testing.T) {
	input := Imports{
		App:    []string{"app-one", "app-two"},
		Vendor: []string{"vendor-one", "vendor-two"},
	}

	result := input.MustParse()
	expected := `import (
	"app-one"
	"app-two"

	"vendor-one"
	"vendor-two"
)`

	assert.Equal(t, expected, result)
}

func TestImports_MustParse_NoApp(t *testing.T) {
	input := Imports{
		Standard: []string{"standard-one", "standard-two"},
		Vendor:   []string{"vendor-one", "vendor-two"},
	}

	result := input.MustParse()
	expected := `import (
	"standard-one"
	"standard-two"

	"vendor-one"
	"vendor-two"
)`

	assert.Equal(t, expected, result)
}

func TestImports_MustParse_NoVendor(t *testing.T) {
	input := Imports{
		Standard: []string{"standard-one", "standard-two"},
		App:      []string{"app-one", "app-two"},
	}

	result := input.MustParse()
	expected := `import (
	"standard-one"
	"standard-two"

	"app-one"
	"app-two"
)`

	assert.Equal(t, expected, result)
}
