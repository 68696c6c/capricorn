package golang

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// isReferenceSliceOrPointerAndTrim tests.
func TestType_isReferenceSliceOrPointerAndTrim_scalar(t *testing.T) {
	input := "string"
	result, resultIsSlice, resultIsPointer := isReferenceSliceOrPointerAndTrim(input)
	assert.Equal(t, "string", result)
	assert.False(t, resultIsSlice)
	assert.False(t, resultIsPointer)
}

func TestType_isReferenceSliceOrPointerAndTrim_slice(t *testing.T) {
	input := "[]int"
	result, resultIsSlice, resultIsPointer := isReferenceSliceOrPointerAndTrim(input)
	assert.Equal(t, "int", result)
	assert.True(t, resultIsSlice)
	assert.False(t, resultIsPointer)
}

func TestType_isReferenceSliceOrPointerAndTrim_pointer(t *testing.T) {
	input := "*rune"
	result, resultIsSlice, resultIsPointer := isReferenceSliceOrPointerAndTrim(input)
	assert.Equal(t, "rune", result)
	assert.False(t, resultIsSlice)
	assert.True(t, resultIsPointer)
}

func TestType_isReferenceSliceOrPointerAndTrim_sliceOfPointers(t *testing.T) {
	input := "[]*time.Time"
	result, resultIsSlice, resultIsPointer := isReferenceSliceOrPointerAndTrim(input)
	assert.Equal(t, "time.Time", result)
	assert.True(t, resultIsSlice)
	assert.True(t, resultIsPointer)
}

// getPkgAndTypeFromReference tests.
func TestType_getPkgAndTypeFromReference_builtin(t *testing.T) {
	input := "string"
	resultPkg, resultType := getPkgAndTypeFromReference(input)
	assert.Equal(t, "", resultPkg)
	assert.Equal(t, "string", resultType)
}

func TestType_getPkgAndTypeFromReference_pkgType(t *testing.T) {
	input := "time.Time"
	resultPkg, resultType := getPkgAndTypeFromReference(input)
	assert.Equal(t, "time", resultPkg)
	assert.Equal(t, "Time", resultType)
}

// NewTypeFromReference tests.
func TestType_NewTypeDataFromReference_builtinScalar(t *testing.T) {
	input := "string"
	result := NewTypeFromReference(input)
	assert.Equal(t, input, result.Reference)
	assert.Equal(t, "", result.Package)
	assert.Equal(t, "string", result.Name)
	assert.Equal(t, "", result.ReceiverName)
	assert.Equal(t, "VARCHAR", result.DataType)
	assert.False(t, result.IsPointer)
	assert.False(t, result.IsSlice)
}

func TestType_NewTypeDataFromReference_builtinSlicePointer(t *testing.T) {
	input := "[]*int"
	result := NewTypeFromReference(input)
	assert.Equal(t, input, result.Reference)
	assert.Equal(t, "", result.Package)
	assert.Equal(t, "int", result.Name)
	assert.Equal(t, "", result.ReceiverName)
	assert.Equal(t, "INT", result.DataType)
	assert.True(t, result.IsPointer)
	assert.True(t, result.IsSlice)
}

func TestType_NewTypeDataFromReference_pkgTypeScalar(t *testing.T) {
	input := "time.Time"
	result := NewTypeFromReference(input)
	assert.Equal(t, input, result.Reference)
	assert.Equal(t, "time", result.Package)
	assert.Equal(t, "Time", result.Name)
	assert.Equal(t, "", result.ReceiverName)
	assert.Equal(t, "", result.DataType)
	assert.False(t, result.IsPointer)
	assert.False(t, result.IsSlice)
}

func TestType_NewTypeDataFromReference_pkgTypeSlicePointer(t *testing.T) {
	input := "[]*time.Time"
	result := NewTypeFromReference(input)
	assert.Equal(t, input, result.Reference)
	assert.Equal(t, "time", result.Package)
	assert.Equal(t, "Time", result.Name)
	assert.Equal(t, "", result.ReceiverName)
	assert.Equal(t, "", result.DataType)
	assert.True(t, result.IsPointer)
	assert.True(t, result.IsSlice)
}

// makeReceiverName tests.
func TestType_makeReceiverName(t *testing.T) {
	input := "exampleType"
	result := makeReceiverName(input)
	assert.Equal(t, "e", result)
}

func TestType_makeReceiverName_exported(t *testing.T) {
	input := "ExampleType"
	result := makeReceiverName(input)
	assert.Equal(t, "e", result)
}
