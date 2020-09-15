package src

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/models/templates/src/handlers"
)

type controllerMeta struct {
	createRequestName string
	updateRequestName string
	viewResponseName  string
	listResponseName  string
	repoType          string
}

type Controller struct {
	fileData    data.FileData
	pathData    data.PathData
	packageData data.PackageData
	name        data.Name
	resource    module.Resource
	receiver    golang.Value

	built      bool
	methodMeta handlers.MethodMeta
	structs    []golang.Struct
	functions  []golang.Function
	imports    golang.Imports
}

func newControllerFromMeta(meta serviceMeta, cMeta controllerMeta) *Controller {
	fileData, pathData := data.MakeGoFileData(meta.packageData.GetImport(), meta.fileName)
	receiver := golang.Value{
		Name: meta.receiverName,
		Type: meta.name.Exported,
	}
	return &Controller{
		fileData:    fileData,
		pathData:    pathData,
		packageData: meta.packageData,
		name:        meta.name,
		resource:    meta.resource,
		methodMeta: handlers.MethodMeta{
			CreateRequestType: cMeta.createRequestName,
			UpdateRequestType: cMeta.updateRequestName,
			ViewResponseType:  cMeta.viewResponseName,
			ListResponseType:  cMeta.listResponseName,
			Resource:          meta.resource,
			Receiver:          receiver,
			RepoField: golang.Value{
				Name: "repo",
				Type: cMeta.repoType,
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
	return golang.File{
		Name:         m.fileData,
		Path:         m.pathData,
		Package:      m.packageData,
		Imports:      m.GetImports(),
		InitFunction: m.GetInit(),
		Consts:       m.GetConsts(),
		Vars:         m.GetVars(),
		Interfaces:   m.GetInterfaces(),
		Structs:      m.GetStructs(),
		Functions:    m.MustGetFunctions(),
	}
}

func (m *Controller) GetImports() golang.Imports {
	if !m.built {
		m.build()
	}
	return m.imports
}

func (m *Controller) GetInit() golang.Function {
	return golang.Function{}
}

func (m *Controller) GetConsts() []golang.Const {
	return []golang.Const{}
}

func (m *Controller) GetVars() []golang.Var {
	return []golang.Var{}
}

func (m *Controller) GetInterfaces() []golang.Interface {
	return []golang.Interface{}
}

func (m *Controller) GetStructs() []golang.Struct {
	if !m.built {
		m.build()
	}
	return m.structs
}

func (m *Controller) MustGetFunctions() []golang.Function {
	if !m.built {
		m.build()
	}
	return m.functions
}

func (m *Controller) build() {
	var imports golang.Imports
	var methods []golang.Function

	// We always need a controller struct.
	structs := []golang.Struct{
		{
			Name: m.resource.Inflection.Plural.Exported,
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
	constructor := handlers.NewConstructor(
		m.name.Exported,
		m.methodMeta.ErrorsField.Name,
		m.methodMeta.RepoField.Name,
		m.methodMeta.RepoField.Type,
	)
	methods = append(methods, constructor.MustGetFunction())
	imports = mergeImports(imports, constructor.GetImports())

	// CRUD methods.
	var needsCreateRequest bool
	var needsUpdateRequest bool
	var needsViewResponse bool
	var needsListResponse bool
	for _, a := range m.resource.Controller.Actions {
		switch a {

		case module.ResourceActionList:
			method := handlers.NewList(m.methodMeta)
			methods = append(methods, method.MustGetFunction())
			imports = mergeImports(imports, method.GetImports())
			needsListResponse = true

		case module.ResourceActionView:
			method := handlers.NewView(m.methodMeta)
			methods = append(methods, method.MustGetFunction())
			imports = mergeImports(imports, method.GetImports())
			needsViewResponse = true

		case module.ResourceActionCreate:
			method := handlers.NewCreate(m.methodMeta)
			methods = append(methods, method.MustGetFunction())
			imports = mergeImports(imports, method.GetImports())
			needsCreateRequest = true
			needsViewResponse = true

		case module.ResourceActionUpdate:
			method := handlers.NewUpdate(m.methodMeta)
			methods = append(methods, method.MustGetFunction())
			imports = mergeImports(imports, method.GetImports())
			needsUpdateRequest = true
			needsViewResponse = true

		case module.ResourceActionDelete:
			method := handlers.NewDelete(m.methodMeta)
			methods = append(methods, method.MustGetFunction())
			imports = mergeImports(imports, method.GetImports())
			needsViewResponse = true

		}
	}

	if needsCreateRequest {
		structs = append(structs, golang.Struct{
			Name: m.methodMeta.CreateRequestType,
			Fields: []golang.Field{
				{
					Name: m.resource.Inflection.Single.Exported,
				},
			},
		})
	}

	if needsUpdateRequest {
		structs = append(structs, golang.Struct{
			Name: m.methodMeta.UpdateRequestType,
			Fields: []golang.Field{
				{
					Name: m.resource.Inflection.Single.Exported,
				},
			},
		})
	}

	if needsViewResponse {
		structs = append(structs, golang.Struct{
			Name: m.methodMeta.ViewResponseType,
			Fields: []golang.Field{
				{
					Name: m.resource.Inflection.Single.Exported,
				},
			},
		})
	}

	if needsListResponse {
		structs = append(structs, golang.Struct{
			Name: m.methodMeta.ListResponseType,
			Fields: []golang.Field{
				{
					Name: "Data",
					Type: "[]*" + m.resource.Inflection.Single.Exported,
					Tags: []golang.Tag{
						{
							Key:    "json",
							Values: []string{"data"},
						},
					},
				},
				{
					Type: "query.Pagination",
					Tags: []golang.Tag{
						{
							Key:    "json",
							Values: []string{"pagination"},
						},
					},
				},
			},
		})
	}

	m.structs = structs
	m.functions = methods
	m.imports = imports
	m.built = true
}

// // If generating a non-DDD app, response names are things like "userResponse" and do not need to be exported.
// // For DDD apps, a response would be named something like "Response".
// // This is because DDD responses are defined in their domain, but used by the router in the app http package as "users.Response".
// func makeController(meta serviceMeta, viewResponseName, listResponseName string) golang.File {
// 	fileData, pathData := data.MakeGoFileData(meta.packageData.GetImport(), meta.fileName)
// 	result := golang.File{
// 		Name:    fileData,
// 		Path:    pathData,
// 		Package: meta.packageData,
// 	}
//
// 	plural := meta.resource.Inflection.Plural
// 	single := meta.resource.Inflection.Single
//
// 	// @TODO need to make the repo struct
//
// 	for _, a := range meta.resource.Controller.Actions {
// 		switch a {
//
// 		case module.ResourceActionList:
// 			t := handlers.List{
// 				Receiver: meta.receiverName,
// 				Plural:   plural,
// 				Single:   single,
// 				Response: listResponseName,
// 			}
// 			h := makeHandler("List", t.MustParse(), handlers.GetListImports())
// 			result.Functions = append(result.Functions, h)
// 			result.Imports = mergeImports(result.Imports, h.Imports)
//
// 		case module.ResourceActionView:
// 			t := handlers.View{
// 				Receiver: meta.receiverName,
// 				Plural:   plural,
// 				Single:   single,
// 				Response: viewResponseName,
// 			}
// 			h := makeHandler("View", t.MustParse(), handlers.GetViewImports())
// 			result.Functions = append(result.Functions, h)
// 			result.Imports = mergeImports(result.Imports, h.Imports)
//
// 		case module.ResourceActionCreate:
// 			t := handlers.Create{
// 				Receiver: meta.receiverName,
// 				Plural:   plural,
// 				Single:   single,
// 				Response: viewResponseName,
// 			}
// 			h := makeHandler("Create", t.MustParse(), handlers.GetCreateImports())
// 			result.Functions = append(result.Functions, h)
// 			result.Imports = mergeImports(result.Imports, h.Imports)
//
// 		case module.ResourceActionUpdate:
// 			t := handlers.Update{
// 				Receiver: meta.receiverName,
// 				Plural:   plural,
// 				Single:   single,
// 				Response: viewResponseName,
// 			}
// 			h := makeHandler("Update", t.MustParse(), handlers.GetUpdateImports())
// 			result.Functions = append(result.Functions, h)
// 			result.Imports = mergeImports(result.Imports, h.Imports)
//
// 		case module.ResourceActionDelete:
// 			t := handlers.Delete{
// 				Receiver: meta.receiverName,
// 				Plural:   plural,
// 				Single:   single,
// 			}
// 			h := makeHandler("Delete", t.MustParse(), handlers.GetDeleteImports())
// 			result.Functions = append(result.Functions, h)
// 			result.Imports = mergeImports(result.Imports, h.Imports)
//
// 		}
// 	}
//
// 	return result
// }
//
// func makeHandler(name, body string, i golang.Imports) golang.Function {
// 	return golang.Function{
// 		Name:    name,
// 		Imports: i,
// 		Arguments: []golang.Value{
// 			{
// 				Name: "cx",
// 				Type: "*gin.Context",
// 			},
// 		},
// 		ReturnValues: []golang.Value{},
// 		Receiver: golang.Value{
// 			Name: "c",
// 			Type: "",
// 		},
// 		Body: body,
// 	}
// }
