package controllers

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
	receiver        golang.Value
	imports         golang.Imports
	args            []golang.Value
	returns         []golang.Value
	StructName      string
	RepoFieldName   string
	ErrorsFieldName string
}

// Non-DDD apps don't export controllers, DDD apps do.  Therefore, the name of the constructor will depend on what kind
// of app we are generating and that decision happens before this function is called.
func NewConstructor(name, controllerType, errorsFieldName, repoFieldName, repoType string) Constructor {
	return Constructor{
		name:     name,
		receiver: golang.Value{},
		imports: golang.Imports{
			Standard: nil,
			App:      nil,
			Vendor:   []string{data.ImportGorm},
		},
		args: []golang.Value{
			{
				Name: repoFieldName,
				Type: repoType,
			},
			{
				Name: errorsFieldName,
				Type: "goat.ErrorHandler",
			},
		},
		returns: []golang.Value{
			{
				Type: controllerType,
			},
		},
		StructName:      controllerType,
		RepoFieldName:   repoFieldName,
		ErrorsFieldName: errorsFieldName,
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
	result, err := utils.ParseTemplateToString("tmp_template_controller_constructor", constructorBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
