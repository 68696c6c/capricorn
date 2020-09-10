package src

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/models/templates/src/models"
)

func makeModel(resource module.Resource, pkgData data.PackageData, fileName string) (golang.File, []module.ResourceField) {
	fileData, pathData := data.MakeGoFileData(pkgData.GetImport(), fileName)
	result := golang.File{
		Name:    fileData,
		Path:    pathData,
		Package: pkgData,
	}

	single := resource.Inflection.Single

	m := golang.Struct{
		Name: single.Exported,
		Fields: []golang.Field{
			{
				Type: "goat.Model",
			},
		},
	}

	var validationFields []module.ResourceField
	for _, f := range resource.Fields {
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
		m.Fields = append(m.Fields, field)

		if f.IsRequired || f.IsUnique {
			validationFields = append(validationFields, f)
		}
	}

	// @TODO does this add a line break between the fields?
	m.Fields = append(m.Fields, golang.Field{
		Name: "",
		Type: "",
	})

	for _, f := range resource.FieldsMeta.BelongsTo {
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
		m.Fields = append(m.Fields, field)
	}

	for _, f := range resource.FieldsMeta.HasMany {
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
		m.Fields = append(m.Fields, field)
	}

	result.Structs = []golang.Struct{m}

	if len(validationFields) > 0 {
		v := models.Validate{
			Receiver: "m",
			Single:   single,
			Fields:   validationFields,
		}
		result.Functions = append(result.Functions, v.MustMakeFunction())
	}

	return result, validationFields
}
