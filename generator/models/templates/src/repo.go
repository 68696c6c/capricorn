package src

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	repos "github.com/68696c6c/capricorn/generator/models/templates/src/repo_methods"
)

type Repo struct {
	fileData    data.FileData
	pathData    data.PathData
	packageData data.PackageData
	name        data.Name
	resource    module.Resource
	receiver    golang.Value

	built      bool
	methodMeta repos.MethodMeta
	functions  []golang.Function
	imports    golang.Imports

	dbFieldName string
}

func newRepoFromMeta(meta serviceMeta) *Repo {
	fileData, pathData := data.MakeGoFileData(meta.packageData.GetImport(), meta.fileName)
	receiver := golang.Value{
		Name: meta.receiverName,
		Type: meta.name.Exported + "Gorm",
	}
	return &Repo{
		fileData:    fileData,
		pathData:    pathData,
		packageData: meta.packageData,
		name:        meta.name,
		resource:    meta.resource,
		methodMeta: repos.MethodMeta{
			DBFieldName: "db",
			Resource:    meta.resource,
			Receiver:    receiver,
		},
		dbFieldName: "db",
	}
}

func (m *Repo) GetInterface() golang.Interface {
	if !m.built {
		m.build()
	}
	return golang.Interface{
		Name:      m.name.Exported,
		Functions: m.functions,
	}
}

func (m *Repo) MustGetFile() golang.File {
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

func (m *Repo) GetImports() golang.Imports {
	if !m.built {
		m.build()
	}
	return m.imports
}

func (m *Repo) GetInit() golang.Function {
	return golang.Function{}
}

func (m *Repo) GetConsts() []golang.Const {
	return []golang.Const{}
}

func (m *Repo) GetVars() []golang.Var {
	return []golang.Var{}
}

func (m *Repo) GetInterfaces() []golang.Interface {
	return []golang.Interface{m.GetInterface()}
}

func (m *Repo) GetStructs() []golang.Struct {
	return []golang.Struct{
		{
			Name: m.resource.Inflection.Single.Exported,
			Fields: []golang.Field{
				{
					Name: m.dbFieldName,
					Type: "*gorm.DB",
				},
			},
		},
	}
}

func (m *Repo) MustGetFunctions() []golang.Function {
	if !m.built {
		m.build()
	}
	return m.functions
}

func (m *Repo) build() {
	var imports golang.Imports
	var methods []golang.Function

	// Default methods.
	constructor := repos.NewConstructor(m.name.Exported, m.receiver.Type)
	methods = append(methods, constructor.MustGetFunction())
	imports = mergeImports(imports, constructor.GetImports())

	getBaseQueryFunc := repos.NewBaseQuery(m.methodMeta)
	methods = append(methods, getBaseQueryFunc.MustGetFunction())
	imports = mergeImports(imports, getBaseQueryFunc.GetImports())

	getFilteredQueryFunc := repos.NewBaseFilteredQuery(m.methodMeta)
	methods = append(methods, getFilteredQueryFunc.MustGetFunction())
	imports = mergeImports(imports, getFilteredQueryFunc.GetImports())

	getPageQueryFunc := repos.NewBasePaginatedQuery(m.methodMeta)
	methods = append(methods, getPageQueryFunc.MustGetFunction())
	imports = mergeImports(imports, getPageQueryFunc.GetImports())

	// CRUD methods.
	for _, a := range m.resource.Repo.Actions {
		switch a {

		case module.ResourceActionList:
			method := repos.NewFilter(m.methodMeta)
			methods = append(methods, method.MustGetFunction())
			imports = mergeImports(imports, method.GetImports())

		case module.ResourceActionView:
			method := repos.NewGetByID(m.methodMeta)
			methods = append(methods, method.MustGetFunction())
			imports = mergeImports(imports, method.GetImports())

		case module.ResourceActionCreate:
			fallthrough
		case module.ResourceActionUpdate:
			method := repos.NewSave(m.methodMeta)
			methods = append(methods, method.MustGetFunction())
			imports = mergeImports(imports, method.GetImports())

		case module.ResourceActionDelete:
			method := repos.NewDelete(m.methodMeta)
			methods = append(methods, method.MustGetFunction())
			imports = mergeImports(imports, method.GetImports())

		}
	}

	m.functions = methods
	m.imports = imports
	m.built = true
}
