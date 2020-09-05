package repo_methods

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var filterBodyTemplate = `
	dataQuery, err := r.getFilteredQuery(q)
	if err != nil {
		return result, errors.Wrap(err, "failed to build filter sites query")
	}

	errs := dataQuery.Find(&result).GetErrors()
	if len(errs) > 0 && goat.ErrorsBesidesRecordNotFound(errs) {
		err := goat.ErrorsToError(errs)
		return result, errors.Wrap(err, "failed to execute filter sites data query")
	}

	if err := r.applyPaginationToQuery(q); err != nil {
		return result, err
	}

	return result, nil
`

type Filter struct {
	Receiver string
	Plural   data.Name
	Single   data.Name
}

func (m Filter) GetName() string {
	return "Filter"
}

func (m Filter) GetImports() golang.Imports {
	return golang.Imports{
		Standard: nil,
		App:      nil,
		Vendor:   []string{data.ImportGoat, data.ImportErrors, data.ImportQuery, data.ImportGorm},
	}
}

func (m Filter) GetArgs() []golang.Value {
	return []golang.Value{
		{
			Name: "id",
			Type: "goat.ID",
		},
	}
}

func (m Filter) GetReturns() []golang.Value {
	return []golang.Value{
		{
			Name: "result",
			Type: "[]*" + m.Single.Exported,
		},
		{
			Name: "err",
			Type: "error",
		},
	}
}

func (m Filter) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_repo_filter", filterBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
