package src

import (
	"github.com/68696c6c/capricorn/generator/utils"

	"github.com/jinzhu/inflection"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
	"github.com/gin-gonic/gin"
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

const updateRequestTemplate = `
type update{{.ResourceNameUpper}}Request struct {
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
	m := req.{{.ResourceNameUpper}}
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

	_, errs := h.app.{{.ResourceNameUpperPlural}}Repo.GetByID(id)
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
	errs = h.app.{{.ResourceNameUpperPlural}}Repo.Save(&req.{{.ResourceNameUpper}})
	if len(errs) > 0 {
		h.app.Errors.HandleErrorsM(c, errs, "failed to save {{.ResourceNameLower}}", goat.RespondServerError)
		return
	}

	goat.RespondData(c, {{.ResourceNameLower}}Response{req.{{.ResourceNameUpper}}})
}`

const handlerGetTemplate = `
func (h {{.StructNameLower}}) Get(c *gin.Context) {
	i := c.Param("id")
	id, err := goat.ParseID(i)
	if err != nil {
		h.app.Errors.HandleErrorM(c, err, "failed to parse id: "+i, goat.RespondBadRequestError)
		return
	}

	m, errs := h.app.{{.ResourceNameUpperPlural}}Repo.GetByID(id)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			h.app.Errors.HandleMessage(c, "{{.ResourceNameLower}} does not exist", goat.RespondNotFoundError)
			return
		} else {
			h.app.Errors.HandleErrorsM(c, errs, "failed to get {{.ResourceNameLower}}", goat.RespondServerError)
			return
		}
	}

	goat.RespondData(c, {{.ResourceNameLower}}Response{*m})
}`

const handlerListTemplate = `
func (h {{.StructNameLower}}) List(c *gin.Context) {
	q := query.NewQueryBuilder(c)

	result, errs := h.app.{{.ResourceNameUpperPlural}}Repo.List(q)
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

	m, errs := h.app.{{.ResourceNameUpperPlural}}Repo.GetByID(id)
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
	errs = h.app.{{.ResourceNameUpperPlural}}Repo.Delete(m)
	if len(errs) > 0 {
		h.app.Errors.HandleErrorsM(c, errs, "failed to delete {{.ResourceNameLower}}", goat.RespondServerError)
		return
	}

	goat.RespondData(c, {{.ResourceNameLower}}Response{*m})
}`

const routesTemplate = `
package http

import (
	"{{.Imports.Packages.App}}"

	"github.com/68696c6c/goat"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func InitRouter(services app.ServiceContainer) {
	router := goat.GetRouter()
	engine := router.GetEngine()

	engine.GET("/health", Health)
	engine.GET("/version", Version)
	api := engine.Group("/api")
	api.Use()
	{
		{{- range $key, $value := .HTTP.RoutesTemplates }}
		{{ $value }}
		{{- end }}
	}

	err := router.Run()
	if err != nil {
		goat.ExitError(errors.Wrap(err, "error starting server"))
	}
}

func Health(c *gin.Context) {
	goat.RespondMessage(c, "ok")
}

func Version(c *gin.Context) {
	// @TODO show version.
	goat.RespondMessage(c, "something helpful here")
}
`

const routeGroupTemplate = `
		{{.ControllerName}} := {{.ControllerConstructor}}(services)
		{{.GroupName}} := api.Group("/{{.GroupName}}")
		{{.GroupName}}.GET("", {{.ControllerName}}.List)
		{{.GroupName}}.GET("/:id", {{.ControllerName}}.Get)
		{{.GroupName}}.POST("", goat.BindMiddleware({{.CreateRequestStructName}}{}), {{.ControllerName}}.Create)
		{{.GroupName}}.PUT("/:id", goat.BindMiddleware({{.UpdateRequestStructName}}{}), {{.ControllerName}}.Update)
		{{.GroupName}}.DELETE("/:id", {{.ControllerName}}.Delete)
`

func CreateHTTP(spec *utils.Spec, logger *logrus.Logger) error {
	logPrefix := "CreateHTTP | "

	err := utils.CreateDir(spec.Paths.HTTP)
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

		resourceUpper := utils.SeparatedToExported(single)
		resourceUpperPlural := utils.SeparatedToExported(plural)

		resourceLower := utils.SeparatedToUnexported(single)
		resourceLowerPlural := utils.SeparatedToUnexported(plural)

		upperControllerName := resourceUpperPlural + "Controller"
		lowerControllerName := resourceLowerPlural + "Controller"

		c.FileName = plural
		c.StructNameUpper = upperControllerName
		c.StructNameLower = lowerControllerName
		c.ResourceNameUpper = resourceUpper
		c.ResourceNameUpperPlural = resourceUpperPlural
		c.ResourceNameLower = resourceLower
		c.ResourceNameLowerPlural = resourceLowerPlural

		// Create requests.
		createRequest, err := utils.ParseTemplateToString("create_request", createRequestTemplate, c)
		if err != nil {
			return errors.Wrap(err, "failed to generate controller request 'create'")
		}
		c.RequestTemplates = append(c.RequestTemplates, createRequest)

		updateRequest, err := utils.ParseTemplateToString("update_request", updateRequestTemplate, c)
		if err != nil {
			return errors.Wrap(err, "failed to generate controller request 'update'")
		}
		c.RequestTemplates = append(c.RequestTemplates, updateRequest)

		// Create responses.
		getResponse, err := utils.ParseTemplateToString("get_response", resourceResponseTemplate, c)
		if err != nil {
			return errors.Wrap(err, "failed to generate controller response 'get'")
		}
		c.ResponseTemplates = append(c.ResponseTemplates, getResponse)

		listResponse, err := utils.ParseTemplateToString("list_response", listResponseTemplate, c)
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
				handler, err := utils.ParseTemplateToString("handler_create", handlerCreateTemplate, c)
				if err != nil {
					return errors.Wrap(err, "failed to generate controller handler 'create'")
				}
				c.HandlerTemplates = append(c.HandlerTemplates, handler)
				c.Routes = append(c.Routes, utils.Route{
					Method: "POST",
					URI:    "",
				})

			case handlerUpdate:
				handler, err := utils.ParseTemplateToString("handler_update", handlerUpdateTemplate, c)
				if err != nil {
					return errors.Wrap(err, "failed to generate controller handler 'update'")
				}
				c.HandlerTemplates = append(c.HandlerTemplates, handler)
				c.Routes = append(c.Routes, utils.Route{
					Method: "PUT",
					URI:    "/:id",
				})

			case handlerGet:
				handler, err := utils.ParseTemplateToString("handler_get", handlerGetTemplate, c)
				if err != nil {
					return errors.Wrap(err, "failed to generate controller handler 'get'")
				}
				c.HandlerTemplates = append(c.HandlerTemplates, handler)
				c.Routes = append(c.Routes, utils.Route{
					Method: "GET",
					URI:    "/:id",
				})

			case handlerList:
				handler, err := utils.ParseTemplateToString("handler_list", handlerListTemplate, c)
				if err != nil {
					return errors.Wrap(err, "failed to generate controller handler 'list'")
				}
				c.HandlerTemplates = append(c.HandlerTemplates, handler)
				c.Routes = append(c.Routes, utils.Route{
					Method: "GET",
					URI:    "",
				})

			case handlerDelete:
				handler, err := utils.ParseTemplateToString("handler_delete", handlerDeleteTemplate, c)
				if err != nil {
					return errors.Wrap(err, "failed to generate controller handler 'delete'")
				}
				c.HandlerTemplates = append(c.HandlerTemplates, handler)
				c.Routes = append(c.Routes, utils.Route{
					Method: "DELETE",
					URI:    "/:id",
				})
			}
		}

		group := utils.RouteGroup{
			ControllerConstructor:   "new" + upperControllerName,
			ControllerName:          lowerControllerName,
			GroupName:               resourceLowerPlural,
			Routes:                  c.Routes,
			CreateRequestStructName: "create" + resourceUpper + "Request",
			UpdateRequestStructName: "update" + resourceUpper + "Request",
		}
		logger.Debug(logPrefix, "route group: ", group)
		spec.HTTP.Routes = append(spec.HTTP.Routes, group)

		routeGroup, err := utils.ParseTemplateToString("route_group", routeGroupTemplate, group)
		if err != nil {
			return errors.Wrap(err, "failed to generate route group")
		}
		logger.Debug(logPrefix, "routeGroup: ", routeGroup)
		spec.HTTP.RoutesTemplates = append(spec.HTTP.RoutesTemplates, routeGroup)

		err = utils.GenerateGoFile(spec.Paths.HTTP, c.FileName, controllerTemplate, *c)
		if err != nil {
			return errors.Wrap(err, "failed to generate controller")
		}
	}

	// Create routes.
	err = utils.GenerateGoFile(spec.Paths.HTTP, "routes", routesTemplate, spec)
	if err != nil {
		return errors.Wrap(err, "failed to generate controller")
	}

	return nil
}
