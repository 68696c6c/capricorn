package services

import (
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/models/templates/src/utils"
)

type Service struct {
	base               utils.Service
	implementationName string
	interfaceName      string
	receiver           golang.Value
}

func NewServiceFromMeta(meta utils.ServiceMeta) *Service {
	base := utils.NewService(meta, meta.Name.Exported)
	return &Service{
		base:               base,
		implementationName: meta.Name.Exported,
		interfaceName:      meta.Name.Exported + "Service",
		receiver:           base.Receiver,
	}
}

func (m *Service) MustGetFile() golang.File {
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

func (m *Service) build() {
	if m.base.Built {
		return
	}

	var imports golang.Imports
	var functions []golang.Function
	var iFunctions []golang.Function

	// Constructor; don't include as part of the interface.
	constructor := NewConstructor(m.interfaceName, m.implementationName)
	functions = append(functions, constructor.MustGetFunction())
	imports = golang.MergeImports(imports, constructor.GetImports())

	// Custom functions.
	for _, a := range m.base.Resource.Service.Actions {
		method := newMethod(a, m.receiver)
		imports = golang.MergeImports(imports, method.GetImports())
		functions = append(functions, method.MustGetFunction())
		iFunctions = append(iFunctions, method.MustGetFunction())
	}

	m.base.Imports = imports
	m.base.Functions = functions
	m.base.Interfaces = []golang.Interface{
		{
			Name:      m.interfaceName,
			Functions: iFunctions,
		},
	}
	m.base.Structs = []golang.Struct{
		{
			Name:   m.implementationName,
			Fields: []golang.Field{},
		},
	}
	m.base.Built = true
}
