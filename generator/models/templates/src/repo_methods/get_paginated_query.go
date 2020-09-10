package repo_methods

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var applyPaginationToQueryBodyTemplate = `
	err := goat.ApplyPaginationToQuery(q, r.getBaseQuery())
	if err != nil {
		return errors.Wrap(err, "failed to set sites query pagination")
	}
	return nil
`

type BasePaginatedQuery BaseQuery

func (m BasePaginatedQuery) GetName() string {
	return "applyPaginationToQuery"
}

func (m BasePaginatedQuery) GetImports() golang.Imports {
	return golang.Imports{
		Standard: []string{},
		App:      []string{},
		Vendor:   []string{data.ImportGoat, data.ImportErrors, data.ImportQuery, data.ImportGorm},
	}
}

func (m BasePaginatedQuery) GetReceiver() golang.Value {
	return golang.Value{
		Name: m.Receiver,
	}
}

func (m BasePaginatedQuery) GetArgs() []golang.Value {
	return []golang.Value{
		{
			Name: "q",
			Type: "*query.Query",
		},
	}
}

func (m BasePaginatedQuery) GetReturns() []golang.Value {
	return []golang.Value{
		{
			Type: "error",
		},
	}
}

func (m BasePaginatedQuery) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_repo_pagination_query", applyPaginationToQueryBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
