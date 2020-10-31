package repos

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/models/templates/src/repos/methods"
	"github.com/68696c6c/capricorn/generator/models/templates/src/utils"
)

type Repo struct {
	base               utils.Service
	implementationName string
	interfaceName      string
	methodMeta         methods.Meta
	constructor        golang.Function
}

func NewRepoFromMeta(meta utils.ServiceMeta) *Repo {
	receiverType := meta.Name.Exported + "Gorm"
	base := utils.NewService(meta, receiverType)
	return &Repo{
		base:               base,
		implementationName: receiverType,
		interfaceName:      meta.Name.Exported,
		methodMeta: methods.Meta{
			DBFieldName: "db",
			Resource:    meta.Resource,
			Receiver:    base.Receiver,
			ModelType:   meta.ModelType,
		},
	}
}

func (m *Repo) GetName() data.Name {
	return m.base.Name
}

func (m *Repo) GetConstructor() golang.Function {
	if !m.base.Built {
		m.build()
	}
	return m.constructor
}

func (m *Repo) GetInterfaceType() data.TypeData {
	return data.MakeTypeDataService(m.base.PackageData.Reference, m.interfaceName, m.base.Receiver.Name, false)
}

func (m *Repo) MustGetFile() golang.File {
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

func (m *Repo) build() {
	if m.base.Built {
		return
	}

	var imports golang.Imports
	var functions []golang.Function

	// Default functions.
	constructor := NewConstructor(m.interfaceName, m.implementationName, m.methodMeta.DBFieldName)
	m.constructor = constructor.MustGetFunction()
	functions = append(functions, constructor.MustGetFunction())
	imports = golang.MergeImports(imports, constructor.GetImports())

	getBaseQueryFunc := methods.NewBaseQuery(m.methodMeta)
	functions = append(functions, getBaseQueryFunc.MustGetFunction())
	imports = golang.MergeImports(imports, getBaseQueryFunc.GetImports())

	getFilteredQueryFunc := methods.NewBaseFilteredQuery(m.methodMeta)
	functions = append(functions, getFilteredQueryFunc.MustGetFunction())
	imports = golang.MergeImports(imports, getFilteredQueryFunc.GetImports())

	getPageQueryFunc := methods.NewBasePaginatedQuery(m.methodMeta)
	functions = append(functions, getPageQueryFunc.MustGetFunction())
	imports = golang.MergeImports(imports, getPageQueryFunc.GetImports())

	// CRUD functions.
	for _, a := range m.base.Resource.Repo.Actions {
		switch a {

		case module.ResourceActionList:
			method := methods.NewFilter(m.methodMeta)
			functions = append(functions, method.MustGetFunction())
			imports = golang.MergeImports(imports, method.GetImports())

		case module.ResourceActionView:
			method := methods.NewGetByID(m.methodMeta)
			functions = append(functions, method.MustGetFunction())
			imports = golang.MergeImports(imports, method.GetImports())

		case module.ResourceActionCreate:
			fallthrough
		case module.ResourceActionUpdate:
			method := methods.NewSave(m.methodMeta)
			functions = append(functions, method.MustGetFunction())
			imports = golang.MergeImports(imports, method.GetImports())

		case module.ResourceActionDelete:
			method := methods.NewDelete(m.methodMeta)
			functions = append(functions, method.MustGetFunction())
			imports = golang.MergeImports(imports, method.GetImports())

		}
	}

	m.base.Built = true
	m.base.Imports = imports
	m.base.Functions = functions
	m.base.Interfaces = []golang.Interface{
		{
			Name:      m.interfaceName,
			Functions: functions,
		},
	}
	m.base.Structs = []golang.Struct{
		{
			Name: m.implementationName,
			Fields: []golang.Field{
				{
					Name: m.methodMeta.DBFieldName,
					Type: "*gorm.DB",
				},
			},
		},
	}
	m.base.Built = true
}
