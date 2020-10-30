package services

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var constructorBodyTemplate = `
	return &{{ .StructName }}{
		{{ .RepoFieldName }}: {{ .RepoArgName }},
	}
`

type Constructor struct {
	name          string
	receiver      golang.Value
	imports       golang.Imports
	args          []golang.Value
	returns       []golang.Value
	StructName    string
	RepoFieldName string
	RepoArgName   string
}

func NewConstructor(interfaceName, implementationName, repoTypeRef string, repoName data.Name) Constructor {
	repoArgName := repoName.Unexported
	return Constructor{
		name:     "New" + interfaceName,
		receiver: golang.Value{},
		imports: golang.Imports{
			Standard: nil,
			App:      nil,
			Vendor:   nil,
		},
		args: []golang.Value{
			{
				Name: repoArgName,
				Type: repoTypeRef,
			},
		},
		returns: []golang.Value{
			{
				Type: interfaceName,
			},
		},
		StructName:    implementationName,
		RepoFieldName: repoName.Unexported,
		RepoArgName:   repoArgName,
	}
}

func (m Constructor) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         m.name,
		Imports:      m.imports,
		Receiver:     m.receiver,
		Arguments:    m.args,
		ReturnValues: m.returns,
		Body:         m.MustParse(),
	}
}

func (m Constructor) GetImports() golang.Imports {
	return m.imports
}

func (m Constructor) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_service_constructor", constructorBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
