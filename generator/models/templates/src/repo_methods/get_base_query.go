package repo_methods

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var getBaseQueryBodyTemplate = `return r.db.Model(&{{ .Single.Exported }}{})`

var getFilteredQueryBodyTemplate = `
	result, err := q.ApplyToGorm(r.getBaseQuery())
	if err != nil {
		return result, err
	}
	return result, nil
`

var applyPaginationToQueryBodyTemplate = `
	err := goat.ApplyPaginationToQuery(q, r.getBaseQuery())
	if err != nil {
		return errors.Wrap(err, "failed to set sites query pagination")
	}
	return nil
`

type Method interface {
	GetName() string
	MustParse() string
	GetImports() golang.Imports
	GetArgs() []golang.Value
	GetReturns() []golang.Value
}

type BaseQuery struct {
	Receiver string
	Plural   data.Name
	Single   data.Name
}

func (m BaseQuery) GetName() string {
	return "getBaseQuery"
}

func (m BaseQuery) GetImports() golang.Imports {
	return golang.Imports{
		Standard: []string{},
		App:      []string{},
		Vendor:   []string{data.ImportGoat, data.ImportErrors, data.ImportQuery, data.ImportGorm},
	}
}

func (m BaseQuery) GetArgs() []golang.Value {
	return []golang.Value{}
}

func (m BaseQuery) GetReturns() []golang.Value {
	return []golang.Value{
		{
			Type: "*gorm.DB",
		},
	}
}

func (m BaseQuery) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_repo_base_query", getBaseQueryBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
