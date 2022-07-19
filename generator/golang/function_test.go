package golang

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunction_MustParse(t *testing.T) {
	f := Function{
		Name: "ExampleFunction",
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
	expected := `func (r ExampleStruct) ExampleFunction(arg1 string, arg2 int) (result string, err error) {
	return "result", nil
}`
	result := f.MustString()
	assert.Equal(t, expected, result)
}

func TestFunction_MustParse_NoReceiver(t *testing.T) {
	f := Function{
		Name: "ExampleFunction",
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
	expected := `func ExampleFunction(arg1 string, arg2 int) (result string, err error) {
	return "result", nil
}`
	result := f.MustString()
	assert.Equal(t, expected, result)
}

func TestFunction_MustParse_SingleNamedReturn(t *testing.T) {
	f := Function{
		Name: "ExampleFunction",
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
				Name: "err",
				Type: "error",
			},
		},
		Body: `return err`,
	}
	expected := `func (r ExampleStruct) ExampleFunction(arg1 string, arg2 int) (err error) {
	return err
}`
	result := f.MustString()
	assert.Equal(t, expected, result)
}

func TestFunction_MustParse_SingleReturn(t *testing.T) {
	f := Function{
		Name: "ExampleFunction",
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
				Type: "error",
			},
		},
		Body: `return nil`,
	}
	expected := `func (r ExampleStruct) ExampleFunction(arg1 string, arg2 int) error {
	return nil
}`
	result := f.MustString()
	assert.Equal(t, expected, result)
}

func TestFunction_MustParse_NoArgs(t *testing.T) {
	f := Function{
		Name: "ExampleFunction",
		Receiver: Value{
			Name: "r",
			Type: "ExampleStruct",
		},
		ReturnValues: []Value{
			{
				Type: "error",
			},
		},
		Body: `return nil`,
	}
	expected := `func (r ExampleStruct) ExampleFunction() error {
	return nil
}`
	result := f.MustString()
	assert.Equal(t, expected, result)
}
