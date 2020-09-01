package handlers

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var listBodyTemplate = `
	q := query.NewQueryBuilder(c)

	result, errs := {{ .Receiver }}.repo.List(q)
	if len(errs) > 0 {
		{{ .Receiver }}.errors.HandleErrorsM(c, errs, "failed to list {{ .Plural.Space }}", goat.RespondServerError)
		return
	}

	errs = {{ .Receiver }}.repo.SetQueryTotal(q)
	if len(errs) > 0 {
		{{ .Receiver }}.errors.HandleErrorsM(c, errs, "failed to count {{ .Plural.Space }}", goat.RespondServerError)
		return
	}

	goat.RespondData(c, {{ .Response }}{result, q.Pagination})
`

type List struct {
	Receiver string
	Plural   data.Name
	Single   data.Name
	Response string
}

func GetListImports() golang.Imports {
	return golang.Imports{
		Standard: nil,
		App:      nil,
		Vendor:   []string{data.ImportGoat, data.ImportQuery, data.ImportGin},
	}
}

func (m List) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_handler_list", listBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
