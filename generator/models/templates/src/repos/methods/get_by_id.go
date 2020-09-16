package methods

import (
	"fmt"

	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var getByIdBodyTemplate = `
	m := {{ .Single.Exported }}{
		Model: goat.Model{
			ID: id,
		},
	}
	errs := {{ .GetDbReference }}.First(&m).GetErrors()
	if len(errs) > 0 {
		return m, goat.ErrorsToError(errs)
	}
	return m, nil
`

type GetByID struct {
	name        string
	dbFieldName string
	receiver    golang.Value
	imports     golang.Imports
	args        []golang.Value
	returns     []golang.Value
	Single      data.Name
}

func NewGetByID(meta Meta) Method {
	return GetByID{
		name:        "GetByID",
		dbFieldName: meta.DBFieldName,
		receiver:    meta.Receiver,
		imports: golang.Imports{
			Standard: nil,
			App:      nil,
			Vendor:   []string{data.ImportGoat},
		},
		args: []golang.Value{
			{
				Name: "id",
				Type: "goat.ID",
			},
		},
		returns: []golang.Value{
			{
				Type: meta.ModelType,
			},
			{
				Type: "error",
			},
		},
		Single: meta.Resource.Inflection.Single,
	}
}

func (m GetByID) GetDbReference() string {
	return fmt.Sprintf("%s.%s", m.receiver.Name, m.dbFieldName)
}

func (m GetByID) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         m.name,
		Imports:      m.imports,
		Receiver:     m.receiver,
		Arguments:    m.args,
		ReturnValues: m.returns,
		Body:         m.MustParse(),
	}
}

func (m GetByID) GetImports() golang.Imports {
	return m.imports
}

func (m GetByID) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_repo_get_by_id", getByIdBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
