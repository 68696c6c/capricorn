package src

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
)

func makeModel(meta serviceMeta) golang.File {
	fileData, pathData := data.MakeGoFileData(meta.packageData.GetImport(), meta.fileName)
	result := golang.File{
		Name:    fileData,
		Path:    pathData,
		Package: meta.packageData,
	}

	single := meta.resource.Inflection.Single

	m := golang.Struct{
		Name: single.Exported,
		Fields: []golang.Field{
			{
				Type: "goat.Model",
			},
		},
	}

	for _, f := range meta.resource.Fields {
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
	}

	// @TODO does this add a line break between the fields?
	m.Fields = append(m.Fields, golang.Field{
		Name: "",
		Type: "",
	})

	for _, f := range meta.resource.FieldsMeta.BelongsTo {
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

	for _, f := range meta.resource.FieldsMeta.HasMany {
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

	// @TODO build validators for any unique fields.

	return result
}
