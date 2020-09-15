package handlers

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var constructorBodyTemplate = `
	return {{ .StructName }}{
		{{ .RepoFieldName }}: {{ .RepoFieldName }},
		{{ .ErrorsFieldName }}: {{ .ErrorsFieldName }},
	}
`

type Constructor struct {
	name            string
	repoType        string
	StructName      string
	RepoFieldName   string
	ErrorsFieldName string
}

func NewConstructor(name, errorsFieldName, repoFieldName, repoType string) Constructor {
	return Constructor{
		name:            "New" + name,
		repoType:        repoType,
		StructName:      name,
		RepoFieldName:   repoFieldName,
		ErrorsFieldName: errorsFieldName,
	}
}

func (m Constructor) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         m.GetName(),
		Imports:      m.GetImports(),
		Receiver:     m.GetReceiver(),
		Arguments:    m.GetArgs(),
		ReturnValues: m.GetReturns(),
		Body:         m.MustParse(),
	}
}

func (m Constructor) GetName() string {
	return m.name
}

func (m Constructor) GetImports() golang.Imports {
	return golang.Imports{
		Standard: nil,
		App:      nil,
		Vendor:   []string{data.ImportGorm},
	}
}

func (m Constructor) GetReceiver() golang.Value {
	return golang.Value{}
}

func (m Constructor) GetArgs() []golang.Value {
	return []golang.Value{
		{
			Name: m.RepoFieldName,
			Type: m.repoType,
		},
		{
			Name: m.ErrorsFieldName,
			Type: "goat.ErrorHandler",
		},
	}
}

func (m Constructor) GetReturns() []golang.Value {
	return []golang.Value{
		{
			Type: m.StructName,
		},
	}
}

func (m Constructor) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_controller_constructor", constructorBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
