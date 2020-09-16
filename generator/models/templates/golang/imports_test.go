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

func Test_MergeImports(t *testing.T) {
	stack := Imports{
		Standard: []string{"one"},
		App:      []string{"one", "two"},
		Vendor:   []string{"one", "two", "three"},
	}

	additional := Imports{
		Standard: []string{},
		App:      []string{},
		Vendor:   []string{"one", "two", "three", "four"},
	}

	stack = MergeImports(stack, additional)

	assert.Len(t, stack.Standard, 1)
	assert.Len(t, stack.App, 2)
	assert.Len(t, stack.Vendor, 4)
}
