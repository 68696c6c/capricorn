package string

import (
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var stringBodyTemplate = `return string(t)`

type String struct {
	name     string
	imports  golang.Imports
	receiver golang.Value
	args     []golang.Value
	returns  []golang.Value
}

func NewString(receiver golang.Value) String {
	return String{
		name: "String",
		imports: golang.Imports{
			Standard: nil,
			App:      nil,
			Vendor:   nil,
		},
		receiver: receiver,
		args:     []golang.Value{},
		returns: []golang.Value{
			{
				Type: "string",
			},
		},
	}
}

func (m String) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         m.name,
		Imports:      m.imports,
		Receiver:     m.receiver,
		Arguments:    m.args,
		ReturnValues: m.returns,
		Body:         m.MustParse(),
	}
}

func (m String) GetImports() golang.Imports {
	return m.imports
}

func (m String) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_enum_string", stringBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
