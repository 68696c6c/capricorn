package repo_methods

import (
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
	receiver     golang.Value
	ReceiverName string
	Single       data.Name
}

func NewFilter(meta MethodMeta) Filter {
	return Filter{
		receiver:     meta.Receiver,
		ReceiverName: meta.Receiver.Name,
		Single:       meta.Resource.Inflection.Single,
	}
}

func (m Filter) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         m.GetName(),
		Imports:      m.GetImports(),
		Receiver:     m.GetReceiver(),
		Arguments:    m.GetArgs(),
		ReturnValues: m.GetReturns(),
		Body:         m.MustParse(),
	}
}

func (m Filter) GetName() string {
	return "Filter"
}

func (m Filter) GetImports() golang.Imports {
	return golang.Imports{
		Standard: nil,
		App:      nil,
		Vendor:   []string{data.ImportGoat, data.ImportErrors, data.ImportQuery, data.ImportGorm},
	}
}

func (m Filter) GetReceiver() golang.Value {
	return m.receiver
}

func (m Filter) GetArgs() []golang.Value {
	return []golang.Value{
		{
			Name: "id",
			Type: "goat.ID",
		},
	}
}

func (m Filter) GetReturns() []golang.Value {
	return []golang.Value{
		{
			Name: "result",
			Type: "[]*" + m.Single.Exported,
		},
		{
			Name: "err",
			Type: "error",
		},
	}
}

func (m Filter) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_repo_filter", filterBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
