package methods

import (
	"fmt"

	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var getFilteredQueryBodyTemplate = `
	result, err := q.ApplyToGorm({{ .ReceiverName }}.getBaseQuery())
	if err != nil {
		return result, err
	}
	return result, nil
`

type BaseFilteredQuery BaseQuery

func NewBaseFilteredQuery(meta Meta) Method {
	return BaseFilteredQuery{
		name:        "getFilteredQuery",
		dbFieldName: meta.DBFieldName,
		receiver:    meta.Receiver,
		imports: golang.Imports{
			Standard: []string{},
			App:      []string{},
			Vendor:   []string{data.ImportQuery, data.ImportGorm},
		},
		args: []golang.Value{
			{
				Name: "q",
				Type: "*query.Query",
			},
		},
		returns: []golang.Value{
			{
				Type: "*gorm.DB",
			},
			{
				Type: "error",
			},
		},
		ReceiverName: meta.Receiver.Name,
		Plural:       meta.Resource.Inflection.Plural,
		Single:       meta.Resource.Inflection.Single,
	}
}

func (m BaseFilteredQuery) GetDbReference() string {
	return fmt.Sprintf("%s.%s", m.receiver.Name, m.dbFieldName)
}

func (m BaseFilteredQuery) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         m.name,
		Imports:      m.imports,
		Receiver:     m.receiver,
		Arguments:    m.args,
		ReturnValues: m.returns,
		Body:         m.MustParse(),
	}
}

func (m BaseFilteredQuery) GetImports() golang.Imports {
	return m.imports
}

func (m BaseFilteredQuery) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_repo_filtered_query", getFilteredQueryBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
