package handlers

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var createBodyTemplate = `
	req, ok := goat.GetRequest(c).(*CreateRequest)
	if !ok {
		{{ .Receiver }}.errors.HandleMessage(c, "failed to get request", goat.RespondBadRequestError)
		return
	}

	// @TODO generate model factories.
	// @TODO generate model validators.
	m := req.{{ .Single.Exported }}
	errs := {{ .Receiver }}.repo.Save(&m)
	if len(errs) > 0 {
		{{ .Receiver }}.errors.HandleErrorsM(c, errs, "failed to save {{ .Single.Space }}", goat.RespondServerError)
		return
	}

	goat.RespondCreated(c, {{ .Response }}{m})
`

type Create struct {
	Receiver string
	Plural   data.Name
	Single   data.Name
	Response string
}

func GetCreateImports() golang.Imports {
	return golang.Imports{
		Standard: nil,
		App:      nil,
		Vendor:   []string{data.ImportGoat, data.ImportGin},
	}
}

func (m Create) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_handler_create", createBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
