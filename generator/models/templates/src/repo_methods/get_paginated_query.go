package repo_methods

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var applyPaginationToQueryBodyTemplate = `
	err := goat.ApplyPaginationToQuery(q, {{ .GetReceiverName }}.getBaseQuery())
	if err != nil {
		return errors.Wrap(err, "failed to set {{ .Single.Space }} query pagination")
	}
	return nil
`

type BasePaginatedQuery BaseQuery

func NewBasePaginatedQuery(meta MethodMeta) BasePaginatedQuery {
	return BasePaginatedQuery{
		dbFieldName: meta.DBFieldName,
		receiver:    meta.Receiver,
		Plural:      meta.Resource.Inflection.Plural,
		Single:      meta.Resource.Inflection.Single,
	}
}

func (m BasePaginatedQuery) GetReceiverName() string {
	return m.receiver.Name
}

func (m BasePaginatedQuery) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         m.GetName(),
		Imports:      m.GetImports(),
		Receiver:     m.GetReceiver(),
		Arguments:    m.GetArgs(),
		ReturnValues: m.GetReturns(),
		Body:         m.MustParse(),
	}
}

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
	return m.receiver
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
