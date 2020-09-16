package methods

import (
	"fmt"

	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var filterBodyTemplate = `
	dataQuery, err := {{ .ReceiverName }}.getFilteredQuery(q)
	if err != nil {
		return result, errors.Wrap(err, "failed to build filter sites query")
	}

	errs := dataQuery.Find(&result).GetErrors()
	if len(errs) > 0 && goat.ErrorsBesidesRecordNotFound(errs) {
		err := goat.ErrorsToError(errs)
		return result, errors.Wrap(err, "failed to execute filter sites data query")
	}

	if err := {{ .ReceiverName }}.applyPaginationToQuery(q); err != nil {
		return result, err
	}

	return result, nil
`

type Filter struct {
	name         string
	dbFieldName  string
	receiver     golang.Value
	imports      golang.Imports
	args         []golang.Value
	returns      []golang.Value
	ReceiverName string
	Single       data.Name
}

func NewFilter(meta Meta) Method {
	return Filter{
		name:        "Filter",
		dbFieldName: meta.DBFieldName,
		receiver:    meta.Receiver,
		imports: golang.Imports{
			Standard: nil,
			App:      nil,
			Vendor:   []string{data.ImportGoat, data.ImportErrors, data.ImportQuery, data.ImportGorm},
		},
		args: []golang.Value{
			{
				Name: "id",
				Type: "goat.ID",
			},
		},
		returns: []golang.Value{
			{
				Name: "result",
				Type: "[]*" + meta.ModelType,
			},
			{
				Name: "err",
				Type: "error",
			},
		},
		ReceiverName: meta.Receiver.Name,
		Single:       meta.Resource.Inflection.Single,
	}
}

func (m Filter) GetDbReference() string {
	return fmt.Sprintf("%s.%s", m.receiver.Name, m.dbFieldName)
}

func (m Filter) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         m.name,
		Imports:      m.imports,
		Receiver:     m.receiver,
		Arguments:    m.args,
		ReturnValues: m.returns,
		Body:         m.MustParse(),
	}
}

func (m Filter) GetImports() golang.Imports {
	return m.imports
}

func (m Filter) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_repo_filter", filterBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
