package validators

import (
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/models/templates/src/utils"
)

type Validator struct {
	base   utils.Service
	fields []*ValidationField
}

func NewValidatorFromMeta(meta utils.ServiceMeta, fields []*ValidationField) *Validator {
	base := utils.NewService(meta, meta.ReceiverName)
	return &Validator{
		base:   base,
		fields: fields,
	}
}

func (m *Validator) MustGetFile() golang.File {
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

func (m *Validator) build() {
	if m.base.Built {
		return
	}

	var structs []golang.Struct
	var functions []golang.Function

	for _, f := range m.fields {
		for _, r := range f.GetRules() {
			structs = append(structs, r.GetStructs()...)
			functions = append(functions, r.MustGetFunctions()...)
		}
	}

	m.base.Structs = structs
	m.base.Functions = functions
	m.base.Built = true
}
