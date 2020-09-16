package handlers

import (
	"fmt"

	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var deleteBodyTemplate = `
	i := {{ .Context.Name }}.Param("id")
	id, err := goat.ParseID(i)
	if err != nil {
		{{ .GetErrorsReference }}.HandleErrorM({{ .Context.Name }}, err, "failed to parse id: "+i, goat.RespondBadRequestError)
		return
	}

	m, errs := {{ .GetRepoReference }}.GetByID(id)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			{{ .GetErrorsReference }}.HandleMessage({{ .Context.Name }}, "{{ .Single.Space }} does not exist", goat.RespondNotFoundError)
			return
		} else {
			{{ .GetErrorsReference }}.HandleErrorsM({{ .Context.Name }}, errs, "failed to get {{ .Single.Space }}", goat.RespondServerError)
			return
		}
	}

	// @TODO generate model factories.
	// @TODO generate model validators.
	errs = {{ .GetRepoReference }}.Delete(&m)
	if len(errs) > 0 {
		{{ .GetErrorsReference }}.HandleErrorsM({{ .Context.Name }}, errs, "failed to delete {{ .Single.Space }}", goat.RespondServerError)
		return
	}

	goat.RespondValid({{ .Context.Name }})
`

type Delete struct {
	receiver golang.Value
	repo     golang.Value
	errors   golang.Value
	Context  golang.Value
	Single   data.Name
}

func NewDelete(meta Meta) Delete {
	return Delete{
		receiver: meta.Receiver,
		repo:     meta.RepoField,
		errors:   meta.ErrorsField,
		Context:  meta.ContextValue,
		Single:   meta.Resource.Inflection.Single,
	}
}

func (m Delete) GetRepoReference() string {
	return fmt.Sprintf("%s.%s", m.receiver.Name, m.repo.Name)
}

func (m Delete) GetErrorsReference() string {
	return fmt.Sprintf("%s.%s", m.receiver.Name, m.errors.Name)
}

func (m Delete) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         "Delete",
		Imports:      m.GetImports(),
		Receiver:     m.receiver,
		Arguments:    []golang.Value{m.Context},
		ReturnValues: []golang.Value{},
		Body:         m.MustParse(),
	}
}

func (m Delete) GetImports() golang.Imports {
	return golang.Imports{
		Standard: nil,
		App:      nil,
		Vendor:   []string{data.ImportGoat, data.ImportGin},
	}
}

func (m Delete) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_handler_delete", deleteBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}