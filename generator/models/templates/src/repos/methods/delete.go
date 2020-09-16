package methods

import (
	"fmt"

	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var deleteBodyTemplate = `
	errs :=  {{ .GetDbReference }}.Delete(m).GetErrors()
	if len(errs) > 0 {
		return goat.ErrorsToError(errs)
	}
	return nil
`

type Delete struct {
	name        string
	dbFieldName string
	receiver    golang.Value
	imports     golang.Imports
	args        []golang.Value
	returns     []golang.Value
	Single      data.Name
}

func NewDelete(meta Meta) Method {
	return Delete{
		name:        "Delete",
		dbFieldName: meta.DBFieldName,
		receiver:    meta.Receiver,
		imports: golang.Imports{
			Standard: nil,
			App:      nil,
			Vendor:   []string{data.ImportGoat},
		},
		args: []golang.Value{
			{
				Name: "m",
				Type: "*" + meta.ModelType,
			},
		},
		returns: []golang.Value{
			{
				Type: "error",
			},
		},
		Single: meta.Resource.Inflection.Single,
	}
}

func (m Delete) GetDbReference() string {
	return fmt.Sprintf("%s.%s", m.receiver.Name, m.dbFieldName)
}

func (m Delete) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         m.name,
		Imports:      m.imports,
		Receiver:     m.receiver,
		Arguments:    m.args,
		ReturnValues: m.returns,
		Body:         m.MustParse(),
	}
}

func (m Delete) GetImports() golang.Imports {
	return m.imports
}

func (m Delete) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_repo_delete", deleteBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
