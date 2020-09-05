package repo_methods

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var saveBodyTemplate = `
	var errs []error
	if m.Model.ID.Valid() {
		errs = r.db.Save(m).GetErrors()
	} else {
		errs = r.db.Create(m).GetErrors()
	}
	if len(errs) > 0 {
		return goat.ErrorsToError(errs)
	}
	return nil
`

type Save struct {
	Receiver string
	Plural   data.Name
	Single   data.Name
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
