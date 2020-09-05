package repo_methods

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var deleteBodyTemplate = `
	errs :=  r.db.Delete(m).GetErrors()
	if len(errs) > 0 {
		return goat.ErrorsToError(errs)
	}
	return nil
`

type Delete struct {
	Receiver string
	Plural   data.Name
	Single   data.Name
}

func (m Delete) GetName() string {
	return "Delete"
}

func (m Delete) GetImports() golang.Imports {
	return golang.Imports{
		Standard: nil,
		App:      nil,
		Vendor:   []string{data.ImportGoat},
	}
}

func (m Delete) GetArgs() []golang.Value {
	return []golang.Value{
		{
			Name: "m",
			Type: "*" + m.Single.Exported,
		},
	}
}

func (m Delete) GetReturns() []golang.Value {
	return []golang.Value{
		{
			Type: "error",
		},
	}
}

func (m Delete) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_repo_delete", deleteBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
