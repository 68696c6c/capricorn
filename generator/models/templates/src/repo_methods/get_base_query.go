package repo_methods

import (
	"fmt"

	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var getBaseQueryBodyTemplate = `return {{ .GetDbReference }}.Model(&{{ .Single.Exported }}{})`

type BaseQuery struct {
	dbFieldName string
	receiver    golang.Value
	Plural      data.Name
	Single      data.Name
}

func NewBaseQuery(meta MethodMeta) BaseQuery {
	return BaseQuery{
		dbFieldName: meta.DBFieldName,
		receiver:    meta.Receiver,
		Plural:      meta.Resource.Inflection.Plural,
		Single:      meta.Resource.Inflection.Single,
	}
}

func (m BaseQuery) GetDbReference() string {
	return fmt.Sprintf("%s.%s", m.receiver.Name, m.dbFieldName)
}

func (m BaseQuery) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         m.GetName(),
		Imports:      m.GetImports(),
		Receiver:     m.GetReceiver(),
		Arguments:    m.GetArgs(),
		ReturnValues: m.GetReturns(),
		Body:         m.MustParse(),
	}
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

func (m BaseQuery) GetReceiver() golang.Value {
	return m.receiver
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
