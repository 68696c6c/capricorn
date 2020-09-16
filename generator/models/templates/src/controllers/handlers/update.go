package handlers

import (
	"fmt"

	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var updateBodyTemplate = `
	i := {{ .Context.Name }}.Param("id")
	id, err := goat.ParseID(i)
	if err != nil {
		{{ .GetErrorsReference }}.HandleErrorM({{ .Context.Name }}, err, "failed to parse id: "+i, goat.RespondBadRequestError)
		return
	}

	// @TODO replace this block with an existence validator and build "not found" handling into the repo.
	_, errs := {{ .GetRepoReference }}.GetByID(id)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			{{ .GetErrorsReference }}.HandleMessage({{ .Context.Name }}, "{{ .Single.Space }} does not exist", goat.RespondNotFoundError)
			return
		} else {
			{{ .GetErrorsReference }}.HandleErrorsM({{ .Context.Name }}, errs, "failed to get {{ .Single.Space }}", goat.RespondServerError)
			return
		}
	}

	req, ok := goat.GetRequest({{ .Context.Name }}).(*UpdateRequest)
	if !ok {
		{{ .GetErrorsReference }}.HandleMessage({{ .Context.Name }}, "failed to get request", goat.RespondBadRequestError)
		return
	}

	// @TODO generate model factories.
	// @TODO generate model validators.
	errs = {{ .GetRepoReference }}.Save(&req.{{ .Single.Exported }})
	if len(errs) > 0 {
		{{ .GetErrorsReference }}.HandleErrorsM({{ .Context.Name }}, errs, "failed to save {{ .Single.Space }}", goat.RespondServerError)
		return
	}

	goat.RespondCreated({{ .Context.Name }}, {{ .ResponseType }}{req.{{ .Single.Exported }}})
`

type Update struct {
	receiver     golang.Value
	repo         golang.Value
	errors       golang.Value
	Context      golang.Value
	Single       data.Name
	ResponseType string
}

func NewUpdate(meta Meta) Update {
	return Update{
		receiver:     meta.Receiver,
		repo:         meta.RepoField,
		errors:       meta.ErrorsField,
		Context:      meta.ContextValue,
		Single:       meta.Resource.Inflection.Single,
		ResponseType: meta.ViewResponseType,
	}
}

func (m Update) GetRepoReference() string {
	return fmt.Sprintf("%s.%s", m.receiver.Name, m.repo.Name)
}

func (m Update) GetErrorsReference() string {
	return fmt.Sprintf("%s.%s", m.receiver.Name, m.errors.Name)
}

func (m Update) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         "Update",
		Imports:      m.GetImports(),
		Receiver:     m.receiver,
		Arguments:    []golang.Value{m.Context},
		ReturnValues: []golang.Value{},
		Body:         m.MustParse(),
	}
}

func (m Update) GetImports() golang.Imports {
	return golang.Imports{
		Standard: nil,
		App:      nil,
		Vendor:   []string{data.ImportGoat, data.ImportGin},
	}
}

func (m Update) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_handler_update", updateBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
