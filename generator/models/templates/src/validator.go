package src

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/models/templates/src/validation"
)

func makeValidator(resource module.Resource, fields []module.ResourceField, pkgData data.PackageData, fileName string) golang.File {
	fileData, pathData := data.MakeGoFileData(pkgData.GetImport(), fileName)
	result := golang.File{
		Name:    fileData,
		Path:    pathData,
		Package: pkgData,
	}

	single := resource.Inflection.Single

	for _, f := range fields {
		if f.IsUnique {
			rule := validation.MakeUniqueRule("db", single, f.Name)
			result.Functions = append(result.Functions, rule.MustMakeFunction())
		}
	}

	return result
}
