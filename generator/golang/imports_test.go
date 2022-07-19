package golang

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	PkgStdOne = MakePackage("standard-one")
	PkgStdTwo = MakePackage("standard-two")

	PkgVendorOne = MakePackage("vendor-one")
	PkgVendorTwo = MakePackage("vendor-two")

	PkgAppOne = MakePackage("app-one")
	PkgAppTwo = MakePackage("app-two")

	PkgOne   = MakePackage("one")
	PkgTwo   = MakePackage("two")
	PkgThree = MakePackage("three")
	PkgFour  = MakePackage("four")
)

func TestImports_MustParse(t *testing.T) {
	input := Imports{
		Standard: []Package{PkgStdOne, PkgStdTwo},
		Vendor:   []Package{PkgVendorOne, PkgVendorTwo},
		App:      []Package{PkgAppOne, PkgAppTwo},
	}

	result := input.MustString()
	expected := `import (
	"standard-one"
	"standard-two"

	"vendor-one"
	"vendor-two"

	"app-one"
	"app-two"
)`

	assert.Equal(t, expected, result)
}

func TestImports_MustParse_None(t *testing.T) {
	input := Imports{}

	result := input.MustString()
	expected := ""

	assert.Equal(t, expected, result)
}

func TestImports_MustParse_NoStandard(t *testing.T) {
	input := Imports{
		Vendor: []Package{PkgVendorOne, PkgVendorTwo},
		App:    []Package{PkgAppOne, PkgAppTwo},
	}

	result := input.MustString()
	expected := `import (
	"vendor-one"
	"vendor-two"

	"app-one"
	"app-two"
)`

	assert.Equal(t, expected, result)
}

func TestImports_MustParse_NoApp(t *testing.T) {
	input := Imports{
		Standard: []Package{PkgStdOne, PkgStdTwo},
		Vendor:   []Package{PkgVendorOne, PkgVendorTwo},
	}

	result := input.MustString()
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
		Standard: []Package{PkgStdOne, PkgStdTwo},
		App:      []Package{PkgAppOne, PkgAppTwo},
	}

	result := input.MustString()
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
		Standard: []Package{PkgOne},
		App:      []Package{PkgOne, PkgTwo},
		Vendor:   []Package{PkgOne, PkgTwo, PkgThree},
	}

	additional := Imports{
		Standard: []Package{},
		App:      []Package{},
		Vendor:   []Package{PkgOne, PkgTwo, PkgThree, PkgFour},
	}

	stack = MergeImports(stack, additional)

	assert.Len(t, stack.Standard, 1)
	assert.Len(t, stack.App, 2)
	assert.Len(t, stack.Vendor, 4)
}
