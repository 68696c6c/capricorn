package string

import (
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var scanBodyTemplate = `
	stringValue := fmt.Sprintf("%v", {{ .ValueArgName }})
	result, err := {{ .FromStringFunctionName }}(stringValue)
	if err != nil {
		return err
	}
	*t = result
	return nil
`

type Scan struct {
	name                   string
	imports                golang.Imports
	receiver               golang.Value
	args                   []golang.Value
	returns                []golang.Value
	FromStringFunctionName string
	ValueArgName           string
}

func NewScan(receiver golang.Value, fromStringFunctionName string) Scan {
	valueArgName := "value"
	return Scan{
		name: "Scan",
		imports: golang.Imports{
			Standard: []string{"fmt"},
			App:      nil,
			Vendor:   nil,
		},
		receiver: golang.Value{
			Name: receiver.Name,
			Type: "*" + receiver.Type,
		},
		args: []golang.Value{
			{
				Name: valueArgName,
				Type: "interface{}",
			},
		},
		returns: []golang.Value{
			{
				Type: "error",
			},
		},
		FromStringFunctionName: fromStringFunctionName,
		ValueArgName:           valueArgName,
	}
}

func (m Scan) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         m.name,
		Imports:      m.imports,
		Receiver:     m.receiver,
		Arguments:    m.args,
		ReturnValues: m.returns,
		Body:         m.MustParse(),
	}
}

func (m Scan) GetImports() golang.Imports {
	return m.imports
}

func (m Scan) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_enum_scan", scanBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
