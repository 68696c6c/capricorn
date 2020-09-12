package src

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/models/templates/src/models"
)

type modelMeta struct {
	receiverName string
	fileName     string
	resource     module.Resource
	packageData  data.PackageData
}

type Model struct {
	receiverName     string
	resource         module.Resource
	single           data.Name
	fieldsSet        bool
	fields           []golang.Field
	validationFields []module.ResourceField
	fileData         data.FileData
	pathData         data.PathData
	packageData      data.PackageData
}

func newModelFromMeta(meta modelMeta) Model {
	fileData, pathData := data.MakeGoFileData(meta.packageData.GetImport(), meta.fileName)
	single := meta.resource.Inflection.Single
	return Model{
		receiverName: meta.receiverName,
		resource:     meta.resource,
		single:       single,
		fileData:     fileData,
		pathData:     pathData,
		packageData:  meta.packageData,
	}
}

func (m Model) GetValidationFields() []module.ResourceField {
	if !m.fieldsSet {
		m.setFields()
	}
	return m.validationFields
}

func (m Model) MustGetFile() golang.File {
	return golang.File{
		Name:      m.fileData,
		Path:      m.pathData,
		Package:   m.packageData,
		Structs:   m.GetStructs(),
		Functions: m.MustGetFunctions(),
	}
}

func (m Model) GetStructs() []golang.Struct {
	if !m.fieldsSet {
		m.setFields()
	}

	return []golang.Struct{
		{
			Name:   m.single.Exported,
			Fields: m.fields,
		},
	}
}

func (m Model) MustGetFunctions() []golang.Function {
	if !m.fieldsSet {
		m.setFields()
	}

	result := []golang.Function{m.getConstructor()}

	if len(m.validationFields) == 0 {
		return result
	}

	v := models.NewValidate(m.receiverName, m.single, m.validationFields)

	result = append(result, v.MustGetFunction())

	return result
}

// @TODO this should be used for making model factories for use in tests, seeders, etc.
func (m Model) getConstructor() golang.Function {
	return golang.Function{}
}

func (m Model) setFields() {
	if m.fieldsSet {
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

	m.fieldsSet = true
}
