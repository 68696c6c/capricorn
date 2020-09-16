package methods

import (
	"fmt"

	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var getBaseQueryBodyTemplate = `return {{ .GetDbReference }}.Model(&{{ .Single.Exported }}{})`

type BaseQuery struct {
	name         string
	dbFieldName  string
	receiver     golang.Value
	imports      golang.Imports
	args         []golang.Value
	returns      []golang.Value
	ReceiverName string
	Plural       data.Name
	Single       data.Name
}

func NewBaseQuery(meta Meta) Method {
	return BaseQuery{
		name:        "getBaseQuery",
		dbFieldName: meta.DBFieldName,
		receiver:    meta.Receiver,
		imports: golang.Imports{
			Standard: []string{},
			App:      []string{},
			Vendor:   []string{data.ImportGoat, data.ImportErrors, data.ImportQuery, data.ImportGorm},
		},
		args: []golang.Value{},
		returns: []golang.Value{
			{
				Type: "*gorm.DB",
			},
		},
		ReceiverName: meta.Receiver.Name,
		Plural:       meta.Resource.Inflection.Plural,
		Single:       meta.Resource.Inflection.Single,
	}
}

func (m BaseQuery) GetDbReference() string {
	return fmt.Sprintf("%s.%s", m.receiver.Name, m.dbFieldName)
}

func (m BaseQuery) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         m.name,
		Imports:      m.imports,
		Receiver:     m.receiver,
		Arguments:    m.args,
		ReturnValues: m.returns,
		Body:         m.MustParse(),
	}
}

func (m BaseQuery) GetImports() golang.Imports {
	return m.imports
}

func (m BaseQuery) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_repo_base_query", getBaseQueryBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
