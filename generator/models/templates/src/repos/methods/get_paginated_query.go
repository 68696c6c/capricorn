package methods

import (
	"fmt"

	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var applyPaginationToQueryBodyTemplate = `
	err := goat.ApplyPaginationToQuery(q, {{ .ReceiverName }}.getBaseQuery())
	if err != nil {
		return errors.Wrap(err, "failed to set {{ .Single.Space }} query pagination")
	}
	return nil
`

type BasePaginatedQuery BaseQuery

func NewBasePaginatedQuery(meta Meta) Method {
	return BasePaginatedQuery{
		name:        "applyPaginationToQuery",
		dbFieldName: meta.DBFieldName,
		receiver:    meta.Receiver,
		imports: golang.Imports{
			Standard: []string{},
			App:      []string{},
			Vendor:   []string{data.ImportGoat, data.ImportErrors, data.ImportQuery, data.ImportGorm},
		},
		args: []golang.Value{
			{
				Name: "q",
				Type: "*query.Query",
			},
		},
		returns: []golang.Value{
			{
				Type: "error",
			},
		},
		ReceiverName: meta.Receiver.Name,
		Plural:       meta.Resource.Inflection.Plural,
		Single:       meta.Resource.Inflection.Single,
	}
}

func (m BasePaginatedQuery) GetDbReference() string {
	return fmt.Sprintf("%s.%s", m.receiver.Name, m.dbFieldName)
}

func (m BasePaginatedQuery) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         m.name,
		Imports:      m.imports,
		Receiver:     m.receiver,
		Arguments:    m.args,
		ReturnValues: m.returns,
		Body:         m.MustParse(),
	}
}

func (m BasePaginatedQuery) GetImports() golang.Imports {
	return m.imports
}

func (m BasePaginatedQuery) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_repo_pagination_query", applyPaginationToQueryBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
