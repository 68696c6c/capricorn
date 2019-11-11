package src

import (
	"fmt"

	"github.com/68696c6c/capricorn/generator/models"
	"github.com/68696c6c/capricorn/generator/utils"

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
	{{- range $key, $value := .Imports }}
	"{{ $value }}"
	{{- end }}

	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/query"
	"github.com/gin-gonic/gin"
)
{{ $tick := "` + "`" + `" }}
type {{.Name.Exported}} struct {
	app app.ServiceContainer
}

func {{.Constructor}}(a app.ServiceContainer) {{.Name.Exported}} {
	return {{.Name.Exported}}{
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
type {{.Name}} struct {
	models.{{.Model}}
}`

const updateRequestTemplate = `
type {{.Name}} struct {
	models.{{.Model}}
}`

const viewResponseTemplate = `
type {{.Name}} struct {
	models.{{.Model}}
}`

const listResponseTemplate = `
{{ $tick := "` + "`" + `" }}
type {{.Name}} struct {
	Data             []*models.{{.Model}} {{ $tick }}json:"data"{{ $tick }}
	query.Pagination {{ $tick }}json:"pagination"{{ $tick }}
}`

const handlerCreateTemplate = `
func (h {{.Receiver}}) {{.Signature}} {
	req, ok := goat.GetRequest(c).(*create{{.Resource.Single.Exported}}Request)
	if !ok {
		h.app.Errors.HandleMessage(c, "failed to get request", goat.RespondBadRequestError)
		return
	}

	// @TODO generate model factories.
	// @TODO generate model validators.
	m := req.{{.Resource.Single.Exported}}
	errs := h.app.{{.Resource.Plural.Exported}}Repo.Save(&m)
	if len(errs) > 0 {
		h.app.Errors.HandleErrorsM(c, errs, "failed to save {{.Resource.Single.Snake}}", goat.RespondServerError)
		return
	}

	goat.RespondCreated(c, {{.Resource.Single.Unexported}}Response{m})
}`

const handlerUpdateTemplate = `
func (h {{.Receiver}}) {{.Signature}} {
	i := c.Param("id")
	id, err := goat.ParseID(i)
	if err != nil {
		h.app.Errors.HandleErrorM(c, err, "failed to parse id: "+i, goat.RespondBadRequestError)
		return
	}

	_, errs := h.app.{{.Resource.Plural.Exported}}Repo.GetByID(id)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			h.app.Errors.HandleMessage(c, "{{.Resource.Single.Snake}} does not exist", goat.RespondNotFoundError)
			return
		} else {
			h.app.Errors.HandleErrorsM(c, errs, "failed to get {{.Resource.Single.Snake}}", goat.RespondServerError)
			return
		}
	}

	req, ok := goat.GetRequest(c).(*update{{.Resource.Single.Exported}}Request)
	if !ok {
		h.app.Errors.HandleMessage(c, "failed to get request", goat.RespondBadRequestError)
		return
	}

	// @TODO generate model factories.
	// @TODO generate model validators.
	errs = h.app.{{.Resource.Plural.Exported}}Repo.Save(&req.{{.Resource.Single.Exported}})
	if len(errs) > 0 {
		h.app.Errors.HandleErrorsM(c, errs, "failed to save {{.Resource.Single.Snake}}", goat.RespondServerError)
		return
	}

	goat.RespondCreated(c, {{.Resource.Single.Unexported}}Response{req.{{.Resource.Single.Exported}}})
}`

const handlerViewTemplate = `
func (h {{.Receiver}}) {{.Signature}} {
	i := c.Param("id")
	id, err := goat.ParseID(i)
	if err != nil {
		h.app.Errors.HandleErrorM(c, err, "failed to parse id: "+i, goat.RespondBadRequestError)
		return
	}

	m, errs := h.app.{{.Resource.Plural.Exported}}Repo.GetByID(id)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			h.app.Errors.HandleMessage(c, "{{.Resource.Single.Snake}} does not exist", goat.RespondNotFoundError)
			return
		} else {
			h.app.Errors.HandleErrorsM(c, errs, "failed to get {{.Resource.Single.Snake}}", goat.RespondServerError)
			return
		}
	}

	goat.RespondData(c, {{.Resource.Single.Unexported}}Response{m})
}`

const handlerListTemplate = `
func (h {{.Receiver}}) {{.Signature}} {
	q := query.NewQueryBuilder(c)

	result, errs := h.app.{{.Resource.Plural.Exported}}Repo.List(q)
	if len(errs) > 0 {
		h.app.Errors.HandleErrorsM(c, errs, "failed to get {{.Resource.Single.Snake}}", goat.RespondServerError)
		return
	}

	errs = h.app.{{.Resource.Plural.Exported}}Repo.SetQueryTotal(q)
	if len(errs) > 0 {
		h.app.Errors.HandleErrorsM(c, errs, "failed to count {{.Resource.Single.Snake}}", goat.RespondServerError)
		return
	}

	goat.RespondData(c, {{.Resource.Plural.Unexported}}Response{result, q.Pagination})
}
`

const handlerDeleteTemplate = `
func (h {{.Receiver}}) {{.Signature}} {
	i := c.Param("id")
	id, err := goat.ParseID(i)
	if err != nil {
		h.app.Errors.HandleErrorM(c, err, "failed to parse id: "+i, goat.RespondBadRequestError)
		return
	}

	m, errs := h.app.{{.Resource.Plural.Exported}}Repo.GetByID(id)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			h.app.Errors.HandleMessage(c, "{{.Resource.Single.Snake}} does not exist", goat.RespondNotFoundError)
			return
		} else {
			h.app.Errors.HandleErrorsM(c, errs, "failed to get {{.Resource.Single.Snake}}", goat.RespondServerError)
			return
		}
	}

	// @TODO generate model factories.
	// @TODO generate model validators.
	errs = h.app.{{.Resource.Plural.Exported}}Repo.Delete(m)
	if len(errs) > 0 {
		h.app.Errors.HandleErrorsM(c, errs, "failed to delete {{.Resource.Single.Snake}}", goat.RespondServerError)
		return
	}

	goat.RespondData(c, {{.Resource.Single.Unexported}}Response{*m})
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

func CreateHTTP(spec models.Project, logger *logrus.Logger) error {
	logPrefix := "CreateHTTP | "

	err := utils.CreateDir(spec.Paths.HTTP)
	if err != nil {
		return errors.Wrapf(err, "failed to create http directory '%s'", spec.Paths.HTTP)
	}

	// Create middlewares.
	// @TODO dewit
	// for _, m := range spec.HTTP.Middlewares {
	// 	logger.Debug(logPrefix, "middleware: ", m.Name)
	// }

	// Create controllers.
	for _, c := range spec.Controllers {
		logger.Debug(logPrefix, "controller: ", c.Filename)

		// Create requests.
		if createRequest, ok := c.Requests["create"]; ok {
			rt, err := utils.ParseTemplateToString("create_request", createRequestTemplate, createRequest)
			if err != nil {
				return errors.Wrap(err, "failed to generate controller request 'create'")
			}
			c.RequestTemplates = append(c.RequestTemplates, rt)
		}

		if updateRequest, ok := c.Requests["update"]; ok {
			rt, err := utils.ParseTemplateToString("update_request", updateRequestTemplate, updateRequest)
			if err != nil {
				return errors.Wrap(err, "failed to generate controller request 'update'")
			}
			c.RequestTemplates = append(c.RequestTemplates, rt)
		}

		// Create responses.
		if viewResponse, ok := c.Responses["view"]; ok {
			rt, err := utils.ParseTemplateToString("view_response", viewResponseTemplate, viewResponse)
			if err != nil {
				return errors.Wrap(err, "failed to generate controller response 'view'")
			}
			c.ResponseTemplates = append(c.ResponseTemplates, rt)
		}

		if listResponse, ok := c.Responses["list"]; ok {
			rt, err := utils.ParseTemplateToString("list_response", listResponseTemplate, listResponse)
			if err != nil {
				return errors.Wrap(err, "failed to generate controller response 'list'")
			}
			c.ResponseTemplates = append(c.ResponseTemplates, rt)
		}

		// Create handlers.
		for _, h := range c.Handlers {
			logger.Infof(logPrefix, fmt.Sprintf("creating controller handler %s", h.Name))

			switch h.Name {
			case "Create":
				ht, err := utils.ParseTemplateToString("controller_create", handlerCreateTemplate, h)
				if err != nil {
					return errors.Wrap(err, "failed to generate controller handler 'Create'")
				}
				c.HandlerTemplates = append(c.HandlerTemplates, ht)

			case "Update":
				ht, err := utils.ParseTemplateToString("controller_update", handlerUpdateTemplate, h)
				if err != nil {
					return errors.Wrap(err, "failed to generate controller handler 'Update'")
				}
				c.HandlerTemplates = append(c.HandlerTemplates, ht)
				// c.RequestTemplates = append(c.RequestTemplates, h.Signature)

			case "List":
				ht, err := utils.ParseTemplateToString("controller_list", handlerListTemplate, h)
				if err != nil {
					return errors.Wrap(err, "failed to generate controller handler 'List'")
				}
				c.HandlerTemplates = append(c.HandlerTemplates, ht)

			case "View":
				ht, err := utils.ParseTemplateToString("controller_view", handlerViewTemplate, h)
				if err != nil {
					return errors.Wrap(err, "failed to generate controller handler 'View'")
				}
				c.HandlerTemplates = append(c.HandlerTemplates, ht)

			case "Delete":
				ht, err := utils.ParseTemplateToString("controller_delete", handlerDeleteTemplate, h)
				if err != nil {
					return errors.Wrap(err, "failed to generate controller handler 'Delete'")
				}
				c.HandlerTemplates = append(c.HandlerTemplates, ht)
			}
		}

		// group := utils.RouteGroup{
		// 	ControllerConstructor:   "new" + upperControllerName,
		// 	ControllerName:          lowerControllerName,
		// 	GroupName:               resourceLowerPlural,
		// 	Routes:                  c.Routes,
		// 	CreateRequestStructName: "create" + resourceUpper + "Request",
		// 	UpdateRequestStructName: "update" + resourceUpper + "Request",
		// }
		// logger.Debug(logPrefix, "route group: ", group)
		// spec.HTTP.Routes = append(spec.HTTP.Routes, group)
		//
		// routeGroup, err := utils.ParseTemplateToString("route_group", routeGroupTemplate, group)
		// if err != nil {
		// 	return errors.Wrap(err, "failed to generate route group")
		// }
		// logger.Debug(logPrefix, "routeGroup: ", routeGroup)
		// spec.HTTP.RoutesTemplates = append(spec.HTTP.RoutesTemplates, routeGroup)

		err = utils.GenerateGoFile(spec.Paths.HTTP, c.Filename, controllerTemplate, c)
		if err != nil {
			return errors.Wrap(err, "failed to generate controller")
		}
	}

	// Create routes.
	// err = utils.GenerateGoFile(spec.Paths.HTTP, "routes", routesTemplate, spec)
	// if err != nil {
	// 	return errors.Wrap(err, "failed to generate controller")
	// }

	return nil
}
