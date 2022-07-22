package project

import (
	"fmt"
	"strings"

	"github.com/68696c6c/capricorn_rnd/generator/golang"
	"github.com/68696c6c/capricorn_rnd/generator/spec"
	"github.com/68696c6c/capricorn_rnd/generator/utils"
)

func MakeModel(name utils.Inflection, resource spec.Resource, enumMap map[string]string, enumsPkg golang.Package) (*golang.File, string) {
	modelName := name.Single.Exported
	fields := []*golang.Field{
		{
			Type: "goat.Model",
		},
	}

	for _, f := range resource.BelongsTo {
		parentName := utils.MakeInflection(f).Single.Snake
		fieldName := utils.MakeInflection(fmt.Sprintf("%s_id", parentName)).Single
		field := &golang.Field{
			Name: fieldName.Exported,
			Type: "goat.ID",
			Tags: []*golang.Tag{
				{
					Key:    "json",
					Values: []string{fieldName.Snake, "omitempty"},
				},
			},
		}
		fields = append(fields, field)
	}

	var hasEnums bool
	for _, f := range resource.Fields {
		fieldType := f.Type
		if strings.HasPrefix(fieldType, "enum:") {
			ft := strings.ReplaceAll(fieldType, "enum:", "")
			fft, ok := enumMap[ft]
			if !ok {
				panic("failed to set model enum field type")
			}
			fieldType = fft
			hasEnums = true
		}

		fieldName := utils.MakeInflection(f.Name).Single
		field := &golang.Field{
			Name: fieldName.Exported,
			Type: fieldType,
		}
		if f.Required {
			field.AddTag(&golang.Tag{
				Key:    "binding",
				Values: []string{"required"},
			})
		}
		field.AddTag(&golang.Tag{
			Key:    "json",
			Values: []string{fieldName.Snake},
		})
		fields = append(fields, field)
	}

	for _, f := range resource.BelongsTo {
		fieldName := utils.MakeInflection(f).Single
		field := &golang.Field{
			Name: fieldName.Exported,
			Type: "*" + fieldName.Exported,
			Tags: []*golang.Tag{
				{
					Key:    "json",
					Values: []string{fieldName.Snake, "omitempty"},
				},
			},
		}
		fields = append(fields, field)
	}

	for _, f := range resource.HasMany {
		fieldName := utils.MakeInflection(f)
		field := &golang.Field{
			Name: fieldName.Plural.Exported,
			Type: "[]*" + fieldName.Single.Exported,
			Tags: []*golang.Tag{
				{
					Key:    "json",
					Values: []string{fieldName.Plural.Snake, "omitempty"},
				},
			},
		}
		fields = append(fields, field)
	}

	var appImports []golang.Package
	if hasEnums {
		appImports = append(appImports, enumsPkg)
	}

	return golang.MakeGoFile(name.Single.Space).SetImports(golang.Imports{
		Vendor: []golang.Package{PkgGoat},
		App:    appImports,
	}).SetStructs([]*golang.Struct{
		{
			Name:   modelName,
			Fields: fields,
		},
	}), modelName
}
