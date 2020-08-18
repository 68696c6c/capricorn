package golang

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeInterfaceTestFunction(name string) (Function, string) {
	f := Function{
		Name: name,
		Receiver: Value{
			Name: "r",
			Type: "ExampleStruct",
		},
		Arguments: []Value{
			{
				Name: "arg1",
				Type: "string",
			},
			{
				Name: "arg2",
				Type: "int",
			},
		},
		ReturnValues: []Value{
			{
				Name: "result",
				Type: "string",
			},
			{
				Name: "err",
				Type: "error",
			},
		},
		Body: `return "result", nil`,
	}
	expected := f.GetSignature()
	return f, expected
}

func TestInterface_MustParse(t *testing.T) {
	f1, f1Expected := makeInterfaceTestFunction("ExampleFunctionOne")
	f2, f2Expected := makeInterfaceTestFunction("ExampleFunctionTwo")

	i := Interface{
		Name:      "ExampleInterface",
		Functions: []Function{f1, f2},
	}
	expected := fmt.Sprintf(`type ExampleInterface interface {
	%s
	%s
}`, f1Expected, f2Expected)

	result := i.MustParse()
	assert.Equal(t, expected, result)
}
