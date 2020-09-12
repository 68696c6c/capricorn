package models

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var validateBodyTemplate = `
	return validation.ValidateStruct({{ .ReceiverName }},{{ .MustGetFields }}
	)
`

type Validate struct {
	ReceiverName string
	DB           string
	Single       data.Name
	Fields       []module.ResourceField
}

func NewValidate(receiverName string, singleName data.Name, fields []module.ResourceField) Validate {
	return Validate{
		ReceiverName: receiverName,
		Single:       singleName,
		Fields:       fields,
	}
}

func (m Validate) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         m.GetName(),
		Imports:      m.GetImports(),
		Arguments:    m.GetArgs(),
		ReturnValues: m.GetReturns(),
		Body:         m.MustParse(),
	}
}

func (m Validate) GetName() string {
	return "Validate"
}

func (m Validate) GetImports() golang.Imports {
	return golang.Imports{
		Standard: nil,
		App:      nil,
		Vendor:   []string{data.ImportGoat},
	}
}

func (m Validate) GetReceiver() golang.Value {
	return golang.Value{
		Name: "*" + m.ReceiverName,
	}
}

func (m Validate) GetArgs() []golang.Value {
	return []golang.Value{
		{
			Name: "db",
			Type: "*gorm.DB",
		},
	}
}

func (m Validate) GetReturns() []golang.Value {
	return []golang.Value{
		{
			Type: "error",
		},
	}
}

func (m Validate) MustGetFields() string {
	var result string
	for _, f := range m.Fields {
		field := ValidationField{
			Receiver: m.ReceiverName,
			Single:   m.Single,
			Field:    f,
		}
		result += field.MustParse()
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
