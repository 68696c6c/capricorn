package golang

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeFileTestInterface(name string) (Interface, string) {
	f1, _ := makeInterfaceTestFunction("ExampleFunctionOne")
	f2, _ := makeInterfaceTestFunction("ExampleFunctionTwo")
	result := Interface{
		Name:      name,
		Functions: []Function{f1, f2},
	}
	expected := result.MustParse()
	return result, expected
}

func makeFileTestStruct(name string) (Struct, string) {
	field1, _ := makeStructTestField("ExportedField", "string")
	field2, _ := makeStructTestField("unexportedField", "int")
	result := Struct{
		Name:   name,
		Fields: []Field{field1, field2},
	}
	expected := result.MustParse()
	return result, expected
}

func TestFile_MustParse(t *testing.T) {
	interface1, expectedInterface1 := makeFileTestInterface("ExampleInterface1")
	interface2, expectedInterface2 := makeFileTestInterface("ExampleInterface2")
	struct1, expectedStruct1 := makeFileTestStruct("ExampleStruct1")
	struct2, expectedStruct2 := makeFileTestStruct("ExampleStruct2")
	func1, _ := makeInterfaceTestFunction("ExampleFunction1")
	func2, _ := makeInterfaceTestFunction("ExampleFunction2")

	expectedFunc1 := func1.MustParse()
	expectedFunc2 := func2.MustParse()

	input := File{
		Package: PackageData{
			Name: "package-name",
		},
		Imports: FileImports{
			Standard: []string{"standard"},
			App:      []string{"app"},
			Vendor:   []string{"vendor"},
		},
		InitFunction: Function{
			Name: "init",
			Body: "return",
		},
		Consts: []Value{
			{
				Name:  "const1",
				Value: `"const 1 value"`,
			},
			{
				Name:  "const2",
				Value: "1",
			},
		},
		Vars: []Value{
			{
				Name:  "var1",
				Value: `"var 1 value"`,
			},
			{
				Name:  "var2",
				Value: "1",
			},
		},
		Interfaces: []Interface{interface1, interface2},
		Structs:    []Struct{struct1, struct2},
		Functions:  []Function{func1, func2},
	}
	expected := fmt.Sprintf(`package package-name

import (
	"standard"

	"app"

	"vendor"
)


func init()  {
	return
}


const const1 = "const 1 value"

const const2 = 1


var var1 = "var 1 value"

var var2 = 1


%s

%s


%s

%s


%s

%s
`, expectedInterface1, expectedInterface2, expectedStruct1, expectedStruct2, expectedFunc1, expectedFunc2)

	result := input.MustParse()
	assert.Equal(t, expected, result)
}

func TestFile_MustParse_NoInit(t *testing.T) {
	interface1, expectedInterface1 := makeFileTestInterface("ExampleInterface1")
	interface2, expectedInterface2 := makeFileTestInterface("ExampleInterface2")
	struct1, expectedStruct1 := makeFileTestStruct("ExampleStruct1")
	struct2, expectedStruct2 := makeFileTestStruct("ExampleStruct2")
	func1, _ := makeInterfaceTestFunction("ExampleFunction1")
	func2, _ := makeInterfaceTestFunction("ExampleFunction2")

	expectedFunc1 := func1.MustParse()
	expectedFunc2 := func2.MustParse()

	input := File{
		Package: PackageData{
			Name: "package-name",
		},
		Imports: FileImports{
			Standard: []string{"standard"},
			App:      []string{"app"},
			Vendor:   []string{"vendor"},
		},
		Consts: []Value{
			{
				Name:  "const1",
				Value: `"const 1 value"`,
			},
			{
				Name:  "const2",
				Value: "1",
			},
		},
		Vars: []Value{
			{
				Name:  "var1",
				Value: `"var 1 value"`,
			},
			{
				Name:  "var2",
				Value: "1",
			},
		},
		Interfaces: []Interface{interface1, interface2},
		Structs:    []Struct{struct1, struct2},
		Functions:  []Function{func1, func2},
	}
	expected := fmt.Sprintf(`package package-name

import (
	"standard"

	"app"

	"vendor"
)


const const1 = "const 1 value"

const const2 = 1


var var1 = "var 1 value"

var var2 = 1


%s

%s


%s

%s


%s

%s
`, expectedInterface1, expectedInterface2, expectedStruct1, expectedStruct2, expectedFunc1, expectedFunc2)

	result := input.MustParse()
	assert.Equal(t, expected, result)
}

func TestFile_MustParse_NoConsts(t *testing.T) {
	interface1, expectedInterface1 := makeFileTestInterface("ExampleInterface1")
	interface2, expectedInterface2 := makeFileTestInterface("ExampleInterface2")
	struct1, expectedStruct1 := makeFileTestStruct("ExampleStruct1")
	struct2, expectedStruct2 := makeFileTestStruct("ExampleStruct2")
	func1, _ := makeInterfaceTestFunction("ExampleFunction1")
	func2, _ := makeInterfaceTestFunction("ExampleFunction2")

	expectedFunc1 := func1.MustParse()
	expectedFunc2 := func2.MustParse()

	input := File{
		Package: PackageData{
			Name: "package-name",
		},
		Imports: FileImports{
			Standard: []string{"standard"},
			App:      []string{"app"},
			Vendor:   []string{"vendor"},
		},
		InitFunction: Function{
			Name: "init",
			Body: "return",
		},
		Vars: []Value{
			{
				Name:  "var1",
				Value: `"var 1 value"`,
			},
			{
				Name:  "var2",
				Value: "1",
			},
		},
		Interfaces: []Interface{interface1, interface2},
		Structs:    []Struct{struct1, struct2},
		Functions:  []Function{func1, func2},
	}
	expected := fmt.Sprintf(`package package-name

import (
	"standard"

	"app"

	"vendor"
)


func init()  {
	return
}


var var1 = "var 1 value"

var var2 = 1


%s

%s


%s

%s


%s

%s
`, expectedInterface1, expectedInterface2, expectedStruct1, expectedStruct2, expectedFunc1, expectedFunc2)

	result := input.MustParse()
	assert.Equal(t, expected, result)
}

func TestFile_MustParse_NoConsts_NoVars(t *testing.T) {
	interface1, expectedInterface1 := makeFileTestInterface("ExampleInterface1")
	interface2, expectedInterface2 := makeFileTestInterface("ExampleInterface2")
	struct1, expectedStruct1 := makeFileTestStruct("ExampleStruct1")
	struct2, expectedStruct2 := makeFileTestStruct("ExampleStruct2")
	func1, _ := makeInterfaceTestFunction("ExampleFunction1")
	func2, _ := makeInterfaceTestFunction("ExampleFunction2")

	expectedFunc1 := func1.MustParse()
	expectedFunc2 := func2.MustParse()

	input := File{
		Package: PackageData{
			Name: "package-name",
		},
		Imports: FileImports{
			Standard: []string{"standard"},
			App:      []string{"app"},
			Vendor:   []string{"vendor"},
		},
		InitFunction: Function{
			Name: "init",
			Body: "return",
		},
		Interfaces: []Interface{interface1, interface2},
		Structs:    []Struct{struct1, struct2},
		Functions:  []Function{func1, func2},
	}
	expected := fmt.Sprintf(`package package-name

import (
	"standard"

	"app"

	"vendor"
)


func init()  {
	return
}


%s

%s


%s

%s


%s

%s
`, expectedInterface1, expectedInterface2, expectedStruct1, expectedStruct2, expectedFunc1, expectedFunc2)

	result := input.MustParse()
	assert.Equal(t, expected, result)
}

func TestFile_MustParse_NoInit_NoConsts_NoVars(t *testing.T) {
	interface1, expectedInterface1 := makeFileTestInterface("ExampleInterface1")
	interface2, expectedInterface2 := makeFileTestInterface("ExampleInterface2")
	struct1, expectedStruct1 := makeFileTestStruct("ExampleStruct1")
	struct2, expectedStruct2 := makeFileTestStruct("ExampleStruct2")
	func1, _ := makeInterfaceTestFunction("ExampleFunction1")
	func2, _ := makeInterfaceTestFunction("ExampleFunction2")

	expectedFunc1 := func1.MustParse()
	expectedFunc2 := func2.MustParse()

	input := File{
		Package: PackageData{
			Name: "package-name",
		},
		Imports: FileImports{
			Standard: []string{"standard"},
			App:      []string{"app"},
			Vendor:   []string{"vendor"},
		},
		Interfaces: []Interface{interface1, interface2},
		Structs:    []Struct{struct1, struct2},
		Functions:  []Function{func1, func2},
	}
	expected := fmt.Sprintf(`package package-name

import (
	"standard"

	"app"

	"vendor"
)


%s

%s


%s

%s


%s

%s
`, expectedInterface1, expectedInterface2, expectedStruct1, expectedStruct2, expectedFunc1, expectedFunc2)

	result := input.MustParse()
	assert.Equal(t, expected, result)
}
