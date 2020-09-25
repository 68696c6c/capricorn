package data

import (
	"fmt"
	"strings"
)

// e.g. reference: pkgname.TypeName
// e.g. package: pkgname
// e.g. name: TypeName
// e.g. data_type: string
type TypeData struct {
	Reference    string `yaml:"reference,omitempty"`
	Package      string `yaml:"package,omitempty"`
	Name         string `yaml:"name,omitempty"`
	ReceiverName string `yaml:"receiver_name,omitempty"`
	DataType     string `yaml:"data_type,omitempty"`
	IsPointer    bool   `yaml:"is_pointer,omitempty"`
	IsSlice      bool   `yaml:"is_slice,omitempty"`
}

func NewTypeDataFromReference(reference string) *TypeData {
	var pkgName string
	var typeName string
	var isPointer bool
	var isSlice bool
	if strings.HasPrefix(reference, "[]") {
		isSlice = true
		reference = strings.TrimPrefix(reference, "[]")
	}
	if strings.HasPrefix(reference, "*") {
		isPointer = true
		reference = strings.TrimPrefix(reference, "*")
	}
	if strings.Contains(reference, ".") {
		parts := strings.Split(reference, ".")
		pkgName = parts[0]
		typeName = parts[1]
	} else {
		typeName = reference
	}
	return MakeTypeData(pkgName, typeName, "", isPointer, isSlice)
}

func makeReferenceName(pkgName, typeName string) string {
	ref := typeName
	if pkgName != "" {
		ref = fmt.Sprintf("%s.%s", pkgName, typeName)
	}
	return ref
}

func MakeTypeData(pkgName, typeName, dataType string, isPointer, isSlice bool) *TypeData {
	return &TypeData{
		Reference:    makeReferenceName(pkgName, typeName),
		Package:      pkgName,
		Name:         typeName,
		ReceiverName: typeName,
		DataType:     dataType,
		IsPointer:    isPointer,
		IsSlice:      isSlice,
	}
}

func MakeTypeDataService(pkgName, typeName, receiverName string, isPointer bool) TypeData {
	return TypeData{
		Reference:    makeReferenceName(pkgName, typeName),
		Package:      pkgName,
		Name:         typeName,
		ReceiverName: receiverName,
		DataType:     "",
		IsPointer:    isPointer,
		IsSlice:      false,
	}
}

func MakeTypeDataID() *TypeData {
	return MakeTypeData("goat", "ID", "BINARY(16) NOT NULL", false, false)
}

func MakeTypeDataCreatedAt() *TypeData {
	return MakeTypeData("time", "Time", "NOT NULL DEFAULT CURRENT_TIMESTAMP", false, false)
}

func MakeTypeDataUpdatedAt() *TypeData {
	return MakeTypeData("time", "Time", "NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP", true, false)
}

func MakeTypeDataDeletedAt() *TypeData {
	return MakeTypeData("time", "Time", "NULL DEFAULT NULL", true, false)
}

func MakeTypeDataString() *TypeData {
	return MakeTypeData("", "string", "VARCHAR(255)", false, false)
}

func MakeTypeDataBelongsTo(pkgName, typeName string) *TypeData {
	return MakeTypeData(pkgName, typeName, "", true, false)
}

func MakeTypeDataHasMany(pkgName, typeName string) *TypeData {
	return MakeTypeData(pkgName, typeName, "", true, true)
}
