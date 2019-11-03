package generator

import (
	"github.com/jinzhu/inflection"
	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"
)

const (
	handlerCreate = "create"
	handlerUpdate = "update"
	handlerGet    = "get"
	handlerList   = "list"
	handlerDelete = "delete"
)

const controllerTemplate = `
package http

import (
	"{{.AppImportPath}}"
	"{{.ModelsImportPath}}"

	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/query"
)

type {{.StructNameLower}} struct {
	app app.ServiceContainer
}

func new{{.StructNameUpper}}(a app.ServiceContainer) {{.StructNameLower}} {
	return {{.StructNameLower}}{
		app: a,
	}
}

{{- range $key, $value := .RequestTemplates }}
{{ $value }}
{{- end }}

{{- range $key, $value := .ResponseTemplates }}
{{ $value }}
{{- end }}

{{- range $key, $value := .HandlerTemplates }}
{{ $value }}
{{- end }}

`

const createRequestTemplate = `
type create{{.ResourceNameUpper}}Request struct {
	models.{{.ResourceNameUpper}}
}`

const resourceResponseTemplate = `
type {{.ResourceNameLower}}Response struct {
	models.{{.ResourceNameUpper}}
}`

const listResponseTemplate = `
{{ $tick := "` + "`" + `" }}
type {{.ResourceNameLowerPlural}}Response struct {
	Data             []*models.{{.ResourceNameUpper}} {{ $tick }}json:"data"{{ $tick }}
	query.Pagination {{ $tick }}json:"pagination"{{ $tick }}
}`

const handlerCreateTemplate = `
func (h {{.StructNameLower}}) Create(c *gin.Context) {
	req, ok := goat.GetRequest(c).(*create{{.ResourceNameUpper}}Request)
	if !ok {
		h.app.Errors.HandleMessage(c, "failed to get request", goat.RespondBadRequestError)
		return
	}

	// @TODO generate model factories.
	// @TODO generate model validators.
	m := models.{{.ResourceNameUpper}}{}
	errs := h.app.{{.ResourceNameUpperPlural}}Repo.Save(&m)
	if len(errs) > 0 {
		h.app.Errors.HandleErrorsM(c, errs, "failed to save {{.ResourceNameLower}}", goat.RespondServerError)
		return
	}

	goat.RespondCreated(c, {{.ResourceNameLower}}Response{m})
}`

const handlerUpdateTemplate = `
func (h {{.StructNameLower}}) Update(c *gin.Context) {
	i := c.Param("id")
	id, err := goat.ParseID(i)
	if err != nil {
		h.app.Errors.HandleErrorM(c, err, "failed to parse id: "+i, goat.RespondBadRequestError)
		return
	}

	m, errs := h.app.{{.ResourceNameUpperPlural}}Repo.GetByID(id, true)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			h.app.Errors.HandleMessage(c, "{{.ResourceNameLower}} does not exist", goat.RespondNotFoundError)
			return
		} else {
			h.app.Errors.HandleErrorsM(c, errs, "failed to get {{.ResourceNameLower}}", goat.RespondServerError)
			return
		}
	}

	req, ok := goat.GetRequest(c).(*create{{.ResourceNameUpper}}Request)
	if !ok {
		h.app.Errors.HandleMessage(c, "failed to get request", goat.RespondBadRequestError)
		return
	}

	// @TODO generate model factories.
	// @TODO generate model validators.
	errs := h.app.{{.ResourceNameUpperPlural}}Repo.Save(&m)
	if len(errs) > 0 {
		h.app.Errors.HandleErrorsM(c, errs, "failed to save {{.ResourceNameLower}}", goat.RespondServerError)
		return
	}

	goat.RespondData(c, {{.ResourceNameLower}}Response{m})
}`

const handlerGetTemplate = `
func (h {{.StructNameLower}}) Get(c *gin.Context) {
	i := c.Param("id")
	id, err := goat.ParseID(i)
	if err != nil {
		h.app.Errors.HandleErrorM(c, err, "failed to parse id: "+i, goat.RespondBadRequestError)
		return
	}

	m, errs := h.app.{{.ResourceNameUpperPlural}}Repo.GetByID(id, true)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			h.app.Errors.HandleMessage(c, "{{.ResourceNameLower}} does not exist", goat.RespondNotFoundError)
			return
		} else {
			h.app.Errors.HandleErrorsM(c, errs, "failed to get {{.ResourceNameLower}}", goat.RespondServerError)
			return
		}
	}

	goat.RespondData(c, {{.ResourceNameLower}}Response{m})
}`

const handlerListTemplate = `
func (h {{.StructNameLower}}) List(c *gin.Context) {
	q := query.NewQueryBuilder(c)

	result, errs := h.app.{{.ResourceNameUpperPlural}}Repo.GetAll(q)
	if len(errs) > 0 {
		h.app.Errors.HandleErrorsM(c, errs, "failed to get {{.ResourceNameLowerPlural}}", goat.RespondServerError)
		return
	}

	errs = h.app.{{.ResourceNameUpperPlural}}Repo.SetQueryTotal(q)
	if len(errs) > 0 {
		h.app.Errors.HandleErrorsM(c, errs, "failed to count {{.ResourceNameLowerPlural}}", goat.RespondServerError)
		return
	}

	goat.RespondData(c, {{.ResourceNameLowerPlural}}Response{result, q.Pagination})
}
`

const handlerDeleteTemplate = `
func (h {{.StructNameLower}}) Delete(c *gin.Context) {
	i := c.Param("id")
	id, err := goat.ParseID(i)
	if err != nil {
		h.app.Errors.HandleErrorM(c, err, "failed to parse id: "+i, goat.RespondBadRequestError)
		return
	}

	m, errs := h.app.{{.ResourceNameUpperPlural}}Repo.GetByID(id, true)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			h.app.Errors.HandleMessage(c, "{{.ResourceNameLower}} does not exist", goat.RespondNotFoundError)
			return
		} else {
			h.app.Errors.HandleErrorsM(c, errs, "failed to get {{.ResourceNameLower}}", goat.RespondServerError)
			return
		}
	}

	// @TODO generate model factories.
	// @TODO generate model validators.
	errs := h.app.{{.ResourceNameUpperPlural}}Repo.Delete(&m)
	if len(errs) > 0 {
		h.app.Errors.HandleErrorsM(c, errs, "failed to delete {{.ResourceNameLower}}", goat.RespondServerError)
		return
	}

	goat.RespondData(c, {{.ResourceNameLower}}Response{m})
}`

func CreateHTTP(spec Spec, logger *logrus.Logger) error {
	logPrefix := "CreateHTTP | "

	err := createDir(spec.Paths.HTTP)
	if err != nil {
		return errors.Wrapf(err, "failed to create http directory '%s'", spec.Paths.HTTP)
	}

	// Create middlewares.
	for _, m := range spec.HTTP.Middlewares {
		logger.Debug(logPrefix, "middleware: ", m.Name)
	}

	// Create controllers.
	for _, c := range spec.HTTP.Controllers {
		logger.Debug(logPrefix, "controller: ", c.Resource)

		c.AppImportPath = spec.Imports.App
		c.ModelsImportPath = spec.Imports.Models

		single := inflection.Singular(c.Resource)
		plural := inflection.Plural(c.Resource)

		resourceUpper := snakeToExportedName(single)
		resourceUpperPlural := snakeToExportedName(plural)

		resourceLower := snakeToUnexportedName(single)
		resourceLowerPlural := snakeToUnexportedName(plural)

		c.FileName = plural
		c.StructNameUpper = resourceUpperPlural + "Controller"
		c.StructNameLower = resourceLowerPlural + "Controller"
		c.ResourceNameUpper = resourceUpper
		c.ResourceNameUpperPlural = resourceUpperPlural
		c.ResourceNameLower = resourceLower
		c.ResourceNameLowerPlural = resourceLowerPlural

		// Create requests.
		request, err := parseTemplateToString("create_request", createRequestTemplate, c)
		if err != nil {
			return errors.Wrap(err, "failed to generate controller request 'create'")
		}
		c.RequestTemplates = append(c.RequestTemplates, request)

		// Create responses.
		getResponse, err := parseTemplateToString("get_response", resourceResponseTemplate, c)
		if err != nil {
			return errors.Wrap(err, "failed to generate controller response 'get'")
		}
		c.ResponseTemplates = append(c.ResponseTemplates, getResponse)

		listResponse, err := parseTemplateToString("list_response", listResponseTemplate, c)
		if err != nil {
			return errors.Wrap(err, "failed to generate controller response 'list'")
		}
		c.ResponseTemplates = append(c.ResponseTemplates, listResponse)

		// If no methods were specified, default to all.
		if len(c.Handlers) == 0 {
			c.Handlers = []string{
				handlerCreate,
				handlerUpdate,
				handlerGet,
				handlerList,
				handlerDelete,
			}
		}

		// Create handlers.
		for _, h := range c.Handlers {
			logger.Debug(logPrefix, "handler: ", h)

			switch h {
			case handlerCreate:
				handler, err := parseTemplateToString("handler_create", handlerCreateTemplate, c)
				if err != nil {
					return errors.Wrap(err, "failed to generate controller handler 'create'")
				}
				c.HandlerTemplates = append(c.HandlerTemplates, handler)

			case handlerUpdate:
				handler, err := parseTemplateToString("handler_update", handlerUpdateTemplate, c)
				if err != nil {
					return errors.Wrap(err, "failed to generate controller handler 'update'")
				}
				c.HandlerTemplates = append(c.HandlerTemplates, handler)

			case handlerGet:
				handler, err := parseTemplateToString("handler_get", handlerGetTemplate, c)
				if err != nil {
					return errors.Wrap(err, "failed to generate controller handler 'get'")
				}
				c.HandlerTemplates = append(c.HandlerTemplates, handler)

			case handlerList:
				handler, err := parseTemplateToString("handler_list", handlerListTemplate, c)
				if err != nil {
					return errors.Wrap(err, "failed to generate controller handler 'list'")
				}
				c.HandlerTemplates = append(c.HandlerTemplates, handler)

			case handlerDelete:
				handler, err := parseTemplateToString("handler_delete", handlerDeleteTemplate, c)
				if err != nil {
					return errors.Wrap(err, "failed to generate controller handler 'delete'")
				}
				c.HandlerTemplates = append(c.HandlerTemplates, handler)
			}
		}

		err = generateFile(spec.Paths.HTTP, c.FileName, controllerTemplate, *c)
		if err != nil {
			return errors.Wrap(err, "failed to generate controller")
		}
	}

	return nil
}
