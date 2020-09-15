package repo_methods

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var getFilteredQueryBodyTemplate = `
	result, err := q.ApplyToGorm({{ .GetReceiverName }}.getBaseQuery())
	if err != nil {
		return result, err
	}
	return result, nil
`

type BaseFilteredQuery BaseQuery

func NewBaseFilteredQuery(meta MethodMeta) BaseFilteredQuery {
	return BaseFilteredQuery{
		dbFieldName: meta.DBFieldName,
		receiver:    meta.Receiver,
		Plural:      meta.Resource.Inflection.Plural,
		Single:      meta.Resource.Inflection.Single,
	}
}

func (m BaseFilteredQuery) GetReceiverName() string {
	return m.receiver.Name
}

func (m BaseFilteredQuery) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         m.GetName(),
		Imports:      m.GetImports(),
		Receiver:     m.GetReceiver(),
		Arguments:    m.GetArgs(),
		ReturnValues: m.GetReturns(),
		Body:         m.MustParse(),
	}
}

func (m BaseFilteredQuery) GetName() string {
	return "getFilteredQuery"
}

func (m BaseFilteredQuery) GetImports() golang.Imports {
	return golang.Imports{
		Standard: []string{},
		App:      []string{},
		Vendor:   []string{data.ImportQuery, data.ImportGorm},
	}
}

func (m BaseFilteredQuery) GetReceiver() golang.Value {
	return m.receiver
}

func (m BaseFilteredQuery) GetArgs() []golang.Value {
	return []golang.Value{
		{
			Name: "q",
			Type: "*query.Query",
		},
	}
}

func (m BaseFilteredQuery) GetReturns() []golang.Value {
	return []golang.Value{
		{
			Type: "*gorm.DB",
		},
		{
			Type: "error",
		},
	}
}

func (m BaseFilteredQuery) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_repo_filtered_query", getFilteredQueryBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
