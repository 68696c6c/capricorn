package methods

import (
	"fmt"

	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var saveBodyTemplate = `
	var errs []error
	if m.Model.ID.Valid() {
		errs = {{ .GetDbReference }}.Save(m).GetErrors()
	} else {
		errs = {{ .GetDbReference }}.Create(m).GetErrors()
	}
	if len(errs) > 0 {
		return goat.ErrorsToError(errs)
	}
	return nil
`

type Save struct {
	name        string
	dbFieldName string
	receiver    golang.Value
	imports     golang.Imports
	args        []golang.Value
	returns     []golang.Value
	Single      data.Name
}

func NewSave(meta Meta) Method {
	return Save{
		name:        "Save",
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

func (m Save) GetDbReference() string {
	return fmt.Sprintf("%s.%s", m.receiver.Name, m.dbFieldName)
}

func (m Save) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         m.name,
		Imports:      m.imports,
		Receiver:     m.receiver,
		Arguments:    m.args,
		ReturnValues: m.returns,
		Body:         m.MustParse(),
	}
}

func (m Save) GetImports() golang.Imports {
	return m.imports
}

func (m Save) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_repo_save", saveBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
