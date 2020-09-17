package controllers

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/models/templates/src/controllers/handlers"
	"github.com/68696c6c/capricorn/generator/models/templates/src/utils"
)

// If generating a non-DDD app, response names are things like "userResponse" and do not need to be exported.
// For DDD apps, a response would be named something like "Response".
// This is because DDD responses are defined in their domain, but used by the router in the app http package as "users.Response".
// Model and Repo types are provided for the same reason.  In a non-DDD app, these will be models.ResourceName and
// repos.ResourceNameRepo instead of ResourceName and Repo.
type ControllerMeta struct {
	CreateRequestType    string
	UpdateRequestType    string
	ResourceResponseType string
	ListResponseType     string
	RepoType             string
	ModelType            string
	Exported             bool
}

type Controller struct {
	base            utils.Service
	controllerType  string
	constructorName string
	modelType       string
	methodMeta      handlers.Meta
}

func NewControllerFromMeta(meta utils.ServiceMeta, cMeta ControllerMeta) *Controller {
	// Non-DDD apps don't export controllers, DDD apps do.
	controllerType := meta.Name.Unexported
	constructorName := "new" + meta.Name.Exported
	if cMeta.Exported {
		controllerType = meta.Name.Exported
		constructorName = "New" + meta.Name.Exported
	}

	base := utils.NewService(meta, controllerType)
	return &Controller{
		base:            base,
		controllerType:  controllerType,
		constructorName: constructorName,
		modelType:       meta.ModelType,
		methodMeta: handlers.Meta{
			CreateRequestType: cMeta.CreateRequestType,
			UpdateRequestType: cMeta.UpdateRequestType,
			ViewResponseType:  cMeta.ResourceResponseType,
			ListResponseType:  cMeta.ListResponseType,
			Resource:          meta.Resource,
			Receiver:          base.Receiver,
			RepoField: golang.Value{
				Name: "repo",
				Type: cMeta.RepoType,
			},
			ErrorsField: golang.Value{
				Name: "errors",
				Type: "goat.ErrorHandler",
			},
			ContextValue: golang.Value{
				Name: "cx",
				Type: "*gin.Context",
			},
		},
	}
}

func (m *Controller) GetType() data.TypeData {
	return data.MakeTypeData(m.base.PackageData.Reference, m.controllerType)
}

func (m *Controller) MustGetFile() golang.File {
	if !m.base.Built {
		m.build()
	}
	return golang.File{
		Name:         m.base.FileData,
		Path:         m.base.PathData,
		Package:      m.base.PackageData,
		Imports:      m.base.Imports,
		InitFunction: golang.Function{},
		Consts:       []golang.Const{},
		Vars:         []golang.Var{},
		Interfaces:   m.base.Interfaces,
		Structs:      m.base.Structs,
		Functions:    m.base.Functions,
	}
}

func (m *Controller) build() {
	if m.base.Built {
		return
	}

	var imports golang.Imports
	var functions []golang.Function

	// We always need a controller struct.
	structs := []golang.Struct{
		{
			Name: m.controllerType,
			Fields: []golang.Field{
				{
					Name: m.methodMeta.RepoField.Name,
					Type: m.methodMeta.RepoField.Type,
				},
				{
					Name: m.methodMeta.ErrorsField.Name,
					Type: m.methodMeta.ErrorsField.Type,
				},
			},
		},
	}

	// Default functions.
	constructor := NewConstructor(
		m.constructorName,
		m.controllerType,
		m.methodMeta.ErrorsField.Name,
		m.methodMeta.RepoField.Name,
		m.methodMeta.RepoField.Type,
	)
	functions = append(functions, constructor.MustGetFunction())
	imports = golang.MergeImports(imports, constructor.GetImports())

	// CRUD functions.
	var needsViewResponse bool
	var needsListResponse bool
	for _, a := range m.base.Resource.Controller.Actions {
		switch a {

		case module.ResourceActionList:
			method := handlers.NewList(m.methodMeta)
			functions = append(functions, method.MustGetFunction())
			imports = golang.MergeImports(imports, method.GetImports())
			needsListResponse = true

		case module.ResourceActionView:
			method := handlers.NewView(m.methodMeta)
			functions = append(functions, method.MustGetFunction())
			imports = golang.MergeImports(imports, method.GetImports())
			needsViewResponse = true

		case module.ResourceActionCreate:
			method := handlers.NewCreate(m.methodMeta)
			functions = append(functions, method.MustGetFunction())
			imports = golang.MergeImports(imports, method.GetImports())
			structs = append(structs, NewRequestStruct(m.methodMeta.CreateRequestType, m.modelType))
			needsViewResponse = true

		case module.ResourceActionUpdate:
			method := handlers.NewUpdate(m.methodMeta)
			functions = append(functions, method.MustGetFunction())
			imports = golang.MergeImports(imports, method.GetImports())
			structs = append(structs, NewRequestStruct(m.methodMeta.UpdateRequestType, m.modelType))
			needsViewResponse = true

		case module.ResourceActionDelete:
			method := handlers.NewDelete(m.methodMeta)
			functions = append(functions, method.MustGetFunction())
			imports = golang.MergeImports(imports, method.GetImports())
			needsViewResponse = true

		}
	}

	if needsViewResponse {
		structs = append(structs, NewResourceResponse(m.methodMeta.ViewResponseType, m.modelType))
	}

	if needsListResponse {
		structs = append(structs, NewListResponse(m.methodMeta.ListResponseType, m.modelType))
	}

	m.base.Imports = imports
	m.base.Structs = structs
	m.base.Functions = functions
	m.base.Built = true
}
