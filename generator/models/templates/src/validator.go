package src

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/models/templates/src/validation_rules/unique"
)

type validatorMeta struct {
	receiverName string
	fileName     string
	resource     module.Resource
	packageData  data.PackageData
	fields       []module.ResourceField
}

type Validator struct {
	db           string
	receiverName string
	single       data.Name
	fields       []module.ResourceField
	fileData     data.FileData
	pathData     data.PathData
	packageData  data.PackageData
}

func newValidatorFromMeta(meta validatorMeta) Validator {
	fileData, pathData := data.MakeGoFileData(meta.packageData.GetImport(), meta.fileName)
	single := meta.resource.Inflection.Single
	return Validator{
		db:           "db",
		receiverName: meta.receiverName,
		single:       single,
		fileData:     fileData,
		pathData:     pathData,
		packageData:  meta.packageData,
	}
}

func (m Validator) MustGetFile() golang.File {
	return golang.File{
		Name:      m.fileData,
		Path:      m.pathData,
		Package:   m.packageData,
		Structs:   m.GetStructs(),
		Functions: m.MustGetFunctions(),
	}
}

func (m Validator) GetStructs() []golang.Struct {
	var result []golang.Struct

	for _, f := range m.fields {
		if f.IsUnique {
			rule := unique.NewRule(m.db, m.receiverName, m.single, f)
			result = append(result, rule.GetStructs()...)
		}
	}

	return result
}

func (m Validator) MustGetFunctions() []golang.Function {
	var result []golang.Function

	for _, f := range m.fields {
		if f.IsUnique {
			rule := unique.NewRule(m.db, m.receiverName, m.single, f)
			result = append(result, rule.MustGetFunctions()...)
		}
	}

	return result
}
