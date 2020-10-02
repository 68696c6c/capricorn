package string

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var valueBodyTemplate = `return string(t), nil`

type Value struct {
	name         string
	imports      golang.Imports
	receiver     golang.Value
	args         []golang.Value
	returns      []golang.Value
	ReceiverName string
	Single       data.Name
}

func NewValue(receiver golang.Value) Value {
	return Value{
		name: "Value",
		imports: golang.Imports{
			Standard: []string{"database/sql/driver"},
			App:      nil,
			Vendor:   nil,
		},
		receiver: receiver,
		args:     []golang.Value{},
		returns: []golang.Value{
			{
				Type: "driver.Value",
			},
			{
				Type: "error",
			},
		},
	}
}

func (m Value) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         m.name,
		Imports:      m.imports,
		Receiver:     m.receiver,
		Arguments:    m.args,
		ReturnValues: m.returns,
		Body:         m.MustParse(),
	}
}

func (m Value) GetImports() golang.Imports {
	return m.imports
}

func (m Value) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_enum_value", valueBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
