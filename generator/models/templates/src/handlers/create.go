package handlers

import (
	"fmt"
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var createBodyTemplate = `
	req, ok := goat.GetRequest({{ .Context.Name }}).(*{{ .RequestType }})
	if !ok {
		{{ .GetErrorsReference }}.HandleMessage(c, "failed to get request", goat.RespondBadRequestError)
		return
	}

	// @TODO generate model factories.
	// @TODO generate model validators.
	m := req.{{ .Single.Exported }}
	errs := {{ .GetRepoReference }}.Save(&m)
	if len(errs) > 0 {
		{{ .GetErrorsReference }}.HandleErrorsM({{ .Context.Name }}, errs, "failed to save {{ .Single.Space }}", goat.RespondServerError)
		return
	}

	goat.RespondCreated({{ .Context.Name }}, {{ .ResponseType }}{m})
`

type Create struct {
	receiver     golang.Value
	repo         golang.Value
	errors       golang.Value
	Context      golang.Value
	Single       data.Name
	RequestType  string
	ResponseType string
}

func NewCreate(meta MethodMeta) Create {
	return Create{
		receiver:     meta.Receiver,
		repo:         meta.RepoField,
		errors:       meta.ErrorsField,
		Context:      meta.ContextValue,
		Single:       meta.Resource.Inflection.Single,
		RequestType:  meta.CreateRequestType,
		ResponseType: meta.ViewResponseType,
	}
}

func (m Create) GetRepoReference() string {
	return fmt.Sprintf("%s.%s", m.receiver.Name, m.repo.Name)
}

func (m Create) GetErrorsReference() string {
	return fmt.Sprintf("%s.%s", m.receiver.Name, m.errors.Name)
}

func (m Create) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         m.GetName(),
		Imports:      m.GetImports(),
		Receiver:     m.GetReceiver(),
		Arguments:    m.GetArgs(),
		ReturnValues: m.GetReturns(),
		Body:         m.MustParse(),
	}
}

func (m Create) GetName() string {
	return "Create"
}

func (m Create) GetImports() golang.Imports {
	return golang.Imports{
		Standard: nil,
		App:      nil,
		Vendor:   []string{data.ImportGoat, data.ImportGin},
	}
}

func (m Create) GetReceiver() golang.Value {
	return m.receiver
}

func (m Create) GetArgs() []golang.Value {
	return []golang.Value{m.Context}
}

func (m Create) GetReturns() []golang.Value {
	return []golang.Value{}
}

func (m Create) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_handler_create", createBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
