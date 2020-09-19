package models

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/models/templates/src/utils"
	"github.com/68696c6c/capricorn/generator/models/templates/src/validators"
)

type Model struct {
	base             utils.Service
	validationMeta   validators.ValidationMeta
	validationFields []*validators.ValidationField
}

func NewModelFromMeta(meta utils.ServiceMeta, validationReceiver string) *Model {
	base := utils.NewService(meta, "*"+meta.Name.Exported)
	return &Model{
		base: base,
		validationMeta: validators.ValidationMeta{
			DBFieldName:   "d",
			ReceiverName:  validationReceiver,
			ModelName:     base.Name,
			ModelReceiver: base.Receiver,
		},
	}
}

func (m *Model) GetType() data.TypeData {
	return data.MakeTypeData(m.base.PackageData.Reference, m.base.Name.Exported)
}

func (m *Model) GetValidationFields() []*validators.ValidationField {
	if !m.base.Built {
		m.build()
	}
	return m.validationFields
}

func (m *Model) MustGetFile() golang.File {
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

func (m *Model) build() {
	if m.base.Built {
		return
	}

	// @TODO include a model factory for use in tests, seeders, etc.
	var functions []golang.Function
	var imports golang.Imports
	var validationFields []*validators.ValidationField

	// @TODO support different base models for soft and hard deletes, uuid or integer ids, different timestamp configurations, etc.
	fields := []golang.Field{
		{
			Type: "goat.Model",
		},
	}

	for _, f := range m.base.Resource.Fields {
		field := golang.Field{
			Name: f.Name.Exported,
			Type: f.TypeData,
			Tags: []golang.Tag{
				{
					Key:    "json",
					Values: []string{f.Name.Snake},
				},
			},
		}
		if f.IsRequired {
			field.Tags = append(field.Tags, golang.Tag{
				Key:    "binding",
				Values: []string{"required"},
			})
		}
		fields = append(fields, field)

		if f.IsRequired || f.IsUnique {
			v := validators.NewValidationField(m.validationMeta, f)
			validationFields = append(validationFields, v)
		}
	}

	// @TODO does this add a line break between the fields?
	fields = append(fields, golang.Field{
		Name: "",
		Type: "",
	})

	for _, f := range m.base.Resource.FieldsMeta.BelongsTo {
		field := golang.Field{
			Name: f.Name.Exported,
			Type: f.TypeData,
			Tags: []golang.Tag{
				{
					Key:    "json",
					Values: []string{f.Name.Snake, "omitempty"},
				},
			},
		}
		fields = append(fields, field)
	}

	for _, f := range m.base.Resource.FieldsMeta.HasMany {
		field := golang.Field{
			Name: f.Name.Exported,
			Type: f.TypeData,
			Tags: []golang.Tag{
				{
					Key:    "json",
					Values: []string{f.Name.Snake, "omitempty"},
				},
			},
		}
		fields = append(fields, field)
	}

	if len(validationFields) > 0 {
		v := NewValidate(m.validationMeta, validationFields)
		functions = append(functions, v.MustGetFunction())
		imports = golang.MergeImports(imports, v.GetImports())
	}

	m.validationFields = validationFields
	m.base.Imports = imports
	m.base.Structs = []golang.Struct{
		{
			Name:   m.base.Name.Exported,
			Fields: fields,
		},
	}
	m.base.Functions = functions
	m.base.Built = true
}
