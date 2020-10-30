package services

import (
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var methodBodyTemplate = `
	return
`

type methodMeta struct {
	name     string
	receiver golang.Value
}

type Method struct {
	name     string
	receiver golang.Value
	imports  golang.Imports
	args     []golang.Value
	returns  []golang.Value
}

func newMethod(name string, receiver golang.Value) Method {
	return Method{
		name:     name,
		receiver: receiver,
		imports: golang.Imports{
			Standard: nil,
			App:      nil,
			Vendor:   nil,
		},
		args:    []golang.Value{},
		returns: []golang.Value{},
	}
}

func (m Method) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         m.name,
		Imports:      m.imports,
		Receiver:     m.receiver,
		Arguments:    m.args,
		ReturnValues: m.returns,
		Body:         m.MustParse(),
	}
}

func (m Method) GetImports() golang.Imports {
	return m.imports
}

func (m Method) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_service_method", methodBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
