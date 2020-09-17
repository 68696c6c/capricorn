package models

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/models/templates/src/validators"
	"github.com/68696c6c/capricorn/generator/utils"
)

var validateBodyTemplate = `
	return validation.ValidateStruct({{ .ReceiverName }},{{ .MustParseFields }}
	)
`

type Validate struct {
	name         string
	dbFieldName  string
	receiver     golang.Value
	imports      golang.Imports
	args         []golang.Value
	returns      []golang.Value
	fields       []*validators.ValidationField
	ReceiverName string
}

func NewValidate(meta validators.ValidationMeta, fields []*validators.ValidationField) Validate {
	return Validate{
		name:        "Validate",
		dbFieldName: meta.DBFieldName,
		receiver:    meta.ModelReceiver,
		imports: golang.Imports{
			Standard: nil,
			App:      nil,
			Vendor:   []string{data.ImportGoat, data.ImportValidation},
		},
		args: []golang.Value{
			{
				Name: "d",
				Type: "*gorm.DB",
			},
		},
		returns: []golang.Value{
			{
				Type: "error",
			},
		},
		fields:       fields,
		ReceiverName: meta.ReceiverName,
	}
}

func (m Validate) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         m.name,
		Imports:      m.imports,
		Receiver:     m.receiver,
		Arguments:    m.args,
		ReturnValues: m.returns,
		Body:         m.MustParse(),
	}
}

func (m Validate) GetImports() golang.Imports {
	return m.imports
}

func (m Validate) MustParseFields() string {
	var result string
	for _, f := range m.fields {
		result += f.MustParse()
	}
	return result
}

func (m Validate) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_validate", validateBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
