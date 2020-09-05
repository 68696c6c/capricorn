package handlers

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var deleteBodyTemplate = `
	i := c.Param("id")
	id, err := goat.ParseID(i)
	if err != nil {
		{{ .Receiver }}.errors.HandleErrorM(c, err, "failed to parse id: "+i, goat.RespondBadRequestError)
		return
	}

	m, errs := {{ .Receiver }}.repo.GetByID(id)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			{{ .Receiver }}.errors.HandleMessage(c, "{{ .Single.Space }} does not exist", goat.RespondNotFoundError)
			return
		} else {
			{{ .Receiver }}.errors.HandleErrorsM(c, errs, "failed to get {{ .Single.Space }}", goat.RespondServerError)
			return
		}
	}

	// @TODO generate model factories.
	// @TODO generate model validators.
	errs = {{ .Receiver }}.repo.Delete(&m)
	if len(errs) > 0 {
		{{ .Receiver }}.errors.HandleErrorsM(c, errs, "failed to delete {{ .Single.Space }}", goat.RespondServerError)
		return
	}

	goat.RespondValid(c)
`

type Delete struct {
	Receiver string
	Plural   data.Name
	Single   data.Name
}

func GetDeleteImports() golang.Imports {
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
