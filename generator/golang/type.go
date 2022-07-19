package golang

import (
	"fmt"
	"strings"
)

// @TODO decimal.Decimal = DECIMAL(18, 4)

var builtinToDataTypeMap = map[string]string{
	"bool":       "BOOLEAN",
	"byte":       "VARCHAR(1)",
	"complex128": "BIGINT",
	"complex64":  "BIGINT",
	"error":      "VARCHAR",
	"float32":    "DOUBLE",
	"float64":    "DOUBLE",
	"int":        "INT",
	"int8":       "INT",
	"int16":      "INT",
	"int32":      "INT",
	"int64":      "INT",
	"rune":       "VARCHAR(1)",
	"string":     "VARCHAR",
	"uint":       "INT",
	"uint8":      "INT",
	"uint16":     "INT",
	"uint32":     "INT",
	"uint64":     "INT",
	"uintptr":    "BIGINT",
}

// Type holds all the information needed to describe a Go type in any context.
// e.g. reference: pkgname.TypeName
// e.g. package: pkgname
// e.g. name: TypeName
// e.g. data_type: string
type Type struct {
	Reference    string `yaml:"reference,omitempty"`
	Package      string `yaml:"package,omitempty"`
	Name         string `yaml:"name,omitempty"`
	ReceiverName string `yaml:"receiver_name,omitempty"`
	DataType     string `yaml:"data_type,omitempty"`
	IsPointer    bool   `yaml:"is_pointer,omitempty"`
	IsSlice      bool   `yaml:"is_slice,omitempty"`
}

func NewTypeFromReference(reference string) *Type {
	var receiverName string
	var dataType string

	trimmed, isSlice, isPointer := isReferenceSliceOrPointerAndTrim(reference)
	pkgName, typeName := getPkgAndTypeFromReference(trimmed)

	// Any slice or pointer syntax has been stripped from the reference now so we can check if it's a builtin.
	if dt, ok := isBuiltIn(typeName); ok {
		dataType = dt
		receiverName = ""
	}

	return &Type{
		Reference:    reference,
		Package:      pkgName,
		Name:         typeName,
		ReceiverName: receiverName,
		DataType:     dataType,
		IsPointer:    isPointer,
		IsSlice:      isSlice,
	}
}

func MakeType(pkgName, typeName, dataType string, isPointer, isSlice bool) *Type {
	return &Type{
		Reference:    makeReferenceName(pkgName, typeName),
		Package:      pkgName,
		Name:         typeName,
		ReceiverName: makeReceiverName(typeName),
		DataType:     dataType,
		IsPointer:    isPointer,
		IsSlice:      isSlice,
	}
}

func MakeTypeID() *Type {
	return MakeType("goat", "ID", "BINARY(16) NOT NULL", false, false)
}

func MakeTypeCreatedAt() *Type {
	return MakeType("time", "Time", "NOT NULL DEFAULT CURRENT_TIMESTAMP", false, false)
}

func MakeTypeUpdatedAt() *Type {
	return MakeType("time", "Time", "NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP", true, false)
}

func MakeTypeDeletedAt() *Type {
	return MakeType("time", "Time", "NULL DEFAULT NULL", true, false)
}

func MakeTypeString() *Type {
	return MakeType("", "string", "VARCHAR(255)", false, false)
}

func MakeTypeBelongsTo(pkgName, typeName string) *Type {
	return MakeType(pkgName, typeName, "", true, false)
}

func MakeTypeHasMany(pkgName, typeName string) *Type {
	return MakeType(pkgName, typeName, "", true, true)
}

// Returns the provided reference string with any pointer or slice prefixes removed.
// Also returns boolean values indicating whether the reference was determined to be a pointer or slice.
// This function checks for both pointer and slice references because the checks for pointers and slices need to be done
// in the correct order.  I.e., the [] needs to be trimmed before we can check if the string starts with *.
func isReferenceSliceOrPointerAndTrim(reference string) (trimmedReference string, isSlice, isPointer bool) {
	result := reference
	if strings.HasPrefix(result, "[]") {
		isSlice = true
		result = strings.TrimPrefix(result, "[]")
	}
	if strings.HasPrefix(result, "*") {
		isPointer = true
		result = strings.TrimPrefix(result, "*")
	}
	return result, isSlice, isPointer
}

func getPkgAndTypeFromReference(trimmedReference string) (pkgName, typeName string) {
	if strings.Contains(trimmedReference, ".") {
		parts := strings.Split(trimmedReference, ".")
		return parts[0], parts[1]
	}
	return "", trimmedReference
}

func isBuiltIn(reference string) (dataType string, isBuiltin bool) {
	dataType, isBuiltin = builtinToDataTypeMap[reference]
	return
}

func makeReferenceName(pkgName, typeName string) string {
	ref := typeName
	if pkgName != "" {
		ref = fmt.Sprintf("%s.%s", pkgName, typeName)
	}
	return ref
}

func makeReceiverName(typeName string) string {
	if len(typeName) < 1 {
		return ""
	}
	return strings.ToLower(typeName[0:1])
}
