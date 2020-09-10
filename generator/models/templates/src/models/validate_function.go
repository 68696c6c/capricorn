package models

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/models/templates/src/validation"
	"github.com/68696c6c/capricorn/generator/utils"
)

var validateBodyTemplate = `
	return validation.ValidateStruct({{ .Receiver }},{{ .GetFields }}
	)
`

type Validate struct {
	Receiver string
	Single   data.Name
	Fields   []module.ResourceField
}

func (m Validate) GetFields() string {
	var result string
	for _, f := range m.Fields {
		field := validation.Field{
			Receiver: m.Receiver,
			Single:   m.Single,
			Field:    f,
		}
		result += field.MustParse()
	}
	return result
}

func (m Validate) MustMakeFunction() golang.Function {
	return golang.Function{
		Name:    "Validate",
		Imports: golang.Imports{},
		Arguments: []golang.Value{
			{
				Name: "db",
				Type: "*gorm.DB",
			},
		},
		ReturnValues: []golang.Value{
			{
				Type: "error",
			},
		},
		Receiver: golang.Value{
			Name: m.Receiver,
			Type: "*" + m.Single.Exported,
		},
		Body: m.MustParse(),
	}
}

func (m Validate) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_validate", validateBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
