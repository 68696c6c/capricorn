package repo_methods

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
	dbFieldName string
	receiver    golang.Value
	Single      data.Name
}

func NewSave(meta MethodMeta) Save {
	return Save{
		dbFieldName: meta.DBFieldName,
		receiver:    meta.Receiver,
		Single:      meta.Resource.Inflection.Single,
	}
}

func (m Save) GetDbReference() string {
	return fmt.Sprintf("%s.%s", m.receiver.Name, m.dbFieldName)
}

func (m Save) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         m.GetName(),
		Imports:      m.GetImports(),
		Receiver:     m.GetReceiver(),
		Arguments:    m.GetArgs(),
		ReturnValues: m.GetReturns(),
		Body:         m.MustParse(),
	}
}

func (m Save) GetName() string {
	return "Save"
}

func (m Save) GetImports() golang.Imports {
	return golang.Imports{
		Standard: nil,
		App:      nil,
		Vendor:   []string{data.ImportGoat},
	}
}

func (m Save) GetReceiver() golang.Value {
	return m.receiver
}

func (m Save) GetArgs() []golang.Value {
	return []golang.Value{
		{
			Name: "m",
			Type: "*" + m.Single.Exported,
		},
	}
}

func (m Save) GetReturns() []golang.Value {
	return []golang.Value{
		{
			Type: "error",
		},
	}
}

func (m Save) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_repo_save", saveBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
