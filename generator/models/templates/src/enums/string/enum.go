package string

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/models/templates/src/enums/meta"
)

type Enum struct {
	typeNameReadable string
	values           []string
	typeAlias        golang.Value
	receiver         golang.Value
}

func NewEnumFromMeta(m meta.Meta) *Enum {
	return &Enum{
		typeNameReadable: m.TypeNameReadable,
		values:           m.Values,
		typeAlias: golang.Value{
			Name: m.TypeName,
			Type: "string",
		},
		receiver: golang.Value{
			Name: m.ReceiverName,
			Type: m.TypeName,
		},
	}
}

func (m *Enum) MustGetFile() golang.File {
	var imports golang.Imports
	var functions []golang.Function

	fromStringFunc := NewFromString(m.typeAlias.Name, m.typeNameReadable, m.values)
	functions = append(functions, fromStringFunc.MustGetFunction())
	imports = golang.MergeImports(imports, fromStringFunc.GetImports())

	stringFunc := NewString(m.receiver)
	functions = append(functions, stringFunc.MustGetFunction())
	imports = golang.MergeImports(imports, stringFunc.GetImports())

	scanFunc := NewScan(m.receiver, fromStringFunc.name)
	functions = append(functions, scanFunc.MustGetFunction())
	imports = golang.MergeImports(imports, scanFunc.GetImports())

	valueFunc := NewValue(m.receiver)
	functions = append(functions, valueFunc.MustGetFunction())
	imports = golang.MergeImports(imports, valueFunc.GetImports())

	return golang.File{
		Name:         data.FileData{},
		Path:         data.PathData{},
		Package:      data.PackageData{},
		Imports:      imports,
		InitFunction: golang.Function{},
		Consts:       nil,
		Vars:         nil,
		Interfaces:   nil,
		TypeAliases:  []golang.Value{m.typeAlias},
		Structs:      nil,
		Functions:    functions,
	}
}
