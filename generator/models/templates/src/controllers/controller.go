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
	fileData    data.FileData
	pathData    data.PathData
	packageData data.PackageData
	name        data.Name
	resource    module.Resource
	receiver    golang.Value

	built      bool
	exported   bool
	methodMeta handlers.Meta
	structs    []golang.Struct
	functions  []golang.Function
	imports    golang.Imports
	modelType  string
}

func NewControllerFromMeta(meta utils.ServiceMeta, cMeta ControllerMeta) *Controller {
	fileData, pathData := data.MakeGoFileData(meta.PackageData.GetImport(), meta.FileName)
	receiver := golang.Value{
		Name: meta.ReceiverName,
		Type: meta.Name.Exported,
	}
	return &Controller{
		fileData:    fileData,
		pathData:    pathData,
		packageData: meta.PackageData,
		name:        meta.Name,
		resource:    meta.Resource,
		modelType:   cMeta.ModelType,
		exported:    cMeta.Exported,
		methodMeta: handlers.Meta{
			CreateRequestType: cMeta.CreateRequestType,
			UpdateRequestType: cMeta.UpdateRequestType,
			ViewResponseType:  cMeta.ResourceResponseType,
			ListResponseType:  cMeta.ListResponseType,
			Resource:          meta.Resource,
			Receiver:          receiver,
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

func (m *Controller) MustGetFile() golang.File {
	if !m.built {
		m.build()
	}
	return golang.File{
		Name:         m.fileData,
		Path:         m.pathData,
		Package:      m.packageData,
		Imports:      m.imports,
		InitFunction: golang.Function{},
		Consts:       []golang.Const{},
		Vars:         []golang.Var{},
		Interfaces:   []golang.Interface{},
		Structs:      m.structs,
		Functions:    m.functions,
	}
}

func (m *Controller) build() {
	var imports golang.Imports
	var methods []golang.Function

	// Non-DDD apps don't export controllers, DDD apps do.
	controllerType := m.name.Unexported
	constructorName := "new" + m.name.Exported
	if m.exported {
		controllerType = m.name.Exported
		constructorName = "New" + m.name.Exported
	}

	// We always need a controller struct.
	structs := []golang.Struct{
		{
			Name: controllerType,
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

	// Default methods.
	constructor := NewConstructor(
		constructorName,
		controllerType,
		m.methodMeta.ErrorsField.Name,
		m.methodMeta.RepoField.Name,
		m.methodMeta.RepoField.Type,
	)
	methods = append(methods, constructor.MustGetFunction())
	imports = golang.MergeImports(imports, constructor.GetImports())

	// CRUD methods.
	var needsViewResponse bool
	var needsListResponse bool
	for _, a := range m.resource.Controller.Actions {
		switch a {

		case module.ResourceActionList:
			method := handlers.NewList(m.methodMeta)
			methods = append(methods, method.MustGetFunction())
			imports = golang.MergeImports(imports, method.GetImports())
			needsListResponse = true

		case module.ResourceActionView:
			method := handlers.NewView(m.methodMeta)
			methods = append(methods, method.MustGetFunction())
			imports = golang.MergeImports(imports, method.GetImports())
			needsViewResponse = true

		case module.ResourceActionCreate:
			method := handlers.NewCreate(m.methodMeta)
			methods = append(methods, method.MustGetFunction())
			imports = golang.MergeImports(imports, method.GetImports())
			structs = append(structs, NewRequestStruct(m.methodMeta.CreateRequestType, m.modelType))
			needsViewResponse = true

		case module.ResourceActionUpdate:
			method := handlers.NewUpdate(m.methodMeta)
			methods = append(methods, method.MustGetFunction())
			imports = golang.MergeImports(imports, method.GetImports())
			structs = append(structs, NewRequestStruct(m.methodMeta.UpdateRequestType, m.modelType))
			needsViewResponse = true

		case module.ResourceActionDelete:
			method := handlers.NewDelete(m.methodMeta)
			methods = append(methods, method.MustGetFunction())
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

	m.structs = structs
	m.functions = methods
	m.imports = imports
	m.built = true
}
