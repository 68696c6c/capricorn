package handlers

import (
	"fmt"

	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var listBodyTemplate = `
	q := query.NewQueryBuilder({{ .Context.Name }})

	result, errs := {{ .GetRepoReference }}.Filter(q)
	if len(errs) > 0 {
		{{ .GetErrorsReference }}.HandleErrorsM({{ .Context.Name }}, errs, "failed to list {{ .Plural.Space }}", goat.RespondServerError)
		return
	}

	goat.RespondData({{ .Context.Name }}, {{ .ResponseType }}{
		Data: result,
		Pagination: q.Pagination,
	})
`

type List struct {
	receiver     golang.Value
	repo         golang.Value
	errors       golang.Value
	Context      golang.Value
	Plural       data.Name
	ResponseType string
}

func NewList(meta MethodMeta) List {
	return List{
		receiver:     meta.Receiver,
		repo:         meta.RepoField,
		errors:       meta.ErrorsField,
		Context:      meta.ContextValue,
		Plural:       meta.Resource.Inflection.Plural,
		ResponseType: meta.ListResponseType,
	}
}

func (m List) GetRepoReference() string {
	return fmt.Sprintf("%s.%s", m.receiver.Name, m.repo.Name)
}

func (m List) GetErrorsReference() string {
	return fmt.Sprintf("%s.%s", m.receiver.Name, m.errors.Name)
}

func (m List) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         m.GetName(),
		Imports:      m.GetImports(),
		Receiver:     m.GetReceiver(),
		Arguments:    m.GetArgs(),
		ReturnValues: m.GetReturns(),
		Body:         m.MustParse(),
	}
}

func (m List) GetName() string {
	return "List"
}

func (m List) GetImports() golang.Imports {
	return golang.Imports{
		Standard: nil,
		App:      nil,
		Vendor:   []string{data.ImportGoat, data.ImportQuery, data.ImportGin},
	}
}

func (m List) GetReceiver() golang.Value {
	return m.receiver
}

func (m List) GetArgs() []golang.Value {
	return []golang.Value{m.Context}
}

func (m List) GetReturns() []golang.Value {
	return []golang.Value{}
}

func (m List) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_handler_list", listBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
