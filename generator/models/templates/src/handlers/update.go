package handlers

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var updateBodyTemplate = `
	i := c.Param("id")
	id, err := goat.ParseID(i)
	if err != nil {
		{{ .Receiver }}.errors.HandleErrorM(c, err, "failed to parse id: "+i, goat.RespondBadRequestError)
		return
	}

	_, errs := {{ .Receiver }}.repo.GetByID(id)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			{{ .Receiver }}.errors.HandleMessage(c, "{{ .Single.Space }} does not exist", goat.RespondNotFoundError)
			return
		} else {
			{{ .Receiver }}.errors.HandleErrorsM(c, errs, "failed to get {{ .Single.Space }}", goat.RespondServerError)
			return
		}
	}

	req, ok := goat.GetRequest(c).(*UpdateRequest)
	if !ok {
		{{ .Receiver }}.errors.HandleMessage(c, "failed to get request", goat.RespondBadRequestError)
		return
	}

	// @TODO generate model factories.
	// @TODO generate model validators.
	errs = {{ .Receiver }}.repo.Save(&req.{{ .Single.Exported }})
	if len(errs) > 0 {
		{{ .Receiver }}.errors.HandleErrorsM(c, errs, "failed to save {{ .Single.Space }}", goat.RespondServerError)
		return
	}

	goat.RespondCreated(c, {{ .Response }}{req.{{ .Single.Exported }}})
`

type Update struct {
	Receiver string
	Plural   data.Name
	Single   data.Name
	Response string
}

func GetUpdateImports() golang.Imports {
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
