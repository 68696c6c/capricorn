package src

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/models/templates/src/models"
)

type Model struct {
	fileData    data.FileData
	pathData    data.PathData
	packageData data.PackageData
	name        data.Name
	resource    module.Resource
	receiver    golang.Value

	built            bool
	fields           []golang.Field
	validationFields []module.ResourceField
}

func newModelFromMeta(meta serviceMeta) *Model {
	fileData, pathData := data.MakeGoFileData(meta.packageData.GetImport(), meta.fileName)
	name := meta.resource.Inflection.Single
	receiver := golang.Value{
		Name: meta.receiverName,
		Type: "*" + name.Exported,
	}
	return &Model{
		fileData:    fileData,
		pathData:    pathData,
		packageData: meta.packageData,
		name:        name,
		resource:    meta.resource,
		receiver:    receiver,
	}
}

func (m *Model) GetType() data.TypeData {
	return data.MakeTypeData(m.packageData.Reference, m.name.Exported)
}

func (m *Model) GetValidationFields() []module.ResourceField {
	if !m.built {
		m.build()
	}
	return m.validationFields
}

func (m *Model) MustGetFile() golang.File {
	if !m.built {
		m.build()
	}
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

func (m Model) GetImports() golang.Imports {
	return golang.Imports{}
}

func (m Model) GetInit() golang.Function {
	return golang.Function{}
}

func (m Model) GetConsts() []golang.Const {
	return []golang.Const{}
}

func (m Model) GetVars() []golang.Var {
	return []golang.Var{}
}

func (m Model) GetInterfaces() []golang.Interface {
	return []golang.Interface{}
}

func (m Model) GetStructs() []golang.Struct {
	if !m.built {
		m.build()
	}

	return []golang.Struct{
		{
			Name:   m.name.Exported,
			Fields: m.fields,
		},
	}
}

func (m Model) MustGetFunctions() []golang.Function {
	if !m.built {
		m.build()
	}

	result := []golang.Function{m.getConstructor()}

	if len(m.validationFields) == 0 {
		return result
	}

	v := models.NewValidate(m.receiver, m.name, m.validationFields)

	result = append(result, v.MustGetFunction())

	return result
}

// @TODO this should be used for making model factories for use in tests, seeders, etc.
func (m Model) getConstructor() golang.Function {
	return golang.Function{}
}

func (m Model) build() {
	if m.built {
		return
	}

	// @TODO support different base models for soft and hard deletes, uuid or integer ids, different timestamp configurations, etc.
	m.fields = []golang.Field{
		{
			Type: "goat.Model",
		},
	}

	for _, f := range m.resource.Fields {
		field := golang.Field{
			Name: f.Name.Exported,
			Type: f.Type,
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
		m.fields = append(m.fields, field)

		if f.IsRequired || f.IsUnique {
			m.validationFields = append(m.validationFields, f)
		}
	}

	// @TODO does this add a line break between the fields?
	m.fields = append(m.fields, golang.Field{
		Name: "",
		Type: "",
	})

	for _, f := range m.resource.FieldsMeta.BelongsTo {
		field := golang.Field{
			Name: f.Name.Exported,
			Type: f.Type,
			Tags: []golang.Tag{
				{
					Key:    "json",
					Values: []string{f.Name.Snake, "omitempty"},
				},
			},
		}
		m.fields = append(m.fields, field)
	}

	for _, f := range m.resource.FieldsMeta.HasMany {
		field := golang.Field{
			Name: f.Name.Exported,
			Type: f.Type,
			Tags: []golang.Tag{
				{
					Key:    "json",
					Values: []string{f.Name.Snake, "omitempty"},
				},
			},
		}
		m.fields = append(m.fields, field)
	}

	m.built = true
}
