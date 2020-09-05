package repo_methods

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

type BaseFilteredQuery BaseQuery

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
