package repo_methods

import (
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
	errs := r.db.First(&m).GetErrors()
	if len(errs) > 0 {
		return m, goat.ErrorsToError(errs)
	}
	return m, nil
`

type GetByID struct {
	Receiver string
	Plural   data.Name
	Single   data.Name
}

func (m GetByID) GetName() string {
	return "GetByID"
}

func (m GetByID) GetImports() golang.Imports {
	return golang.Imports{
		Standard: nil,
		App:      nil,
		Vendor:   []string{data.ImportGoat},
	}
}

func (m GetByID) GetArgs() []golang.Value {
	return []golang.Value{
		{
			Name: "id",
			Type: "goat.ID",
		},
	}
}

func (m GetByID) GetReturns() []golang.Value {
	return []golang.Value{
		{
			Type: m.Single.Exported,
		},
		{
			Type: "error",
		},
	}
}

func (m GetByID) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_repo_get_by_id", getByIdBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
