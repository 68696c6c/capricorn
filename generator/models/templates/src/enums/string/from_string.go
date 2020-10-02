package string

import (
	"fmt"
	"strings"

	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var fromStringBodyTemplate = `
	values := []{{ .TypeName }}{
{{ .Values }}
	}
	for _, v := range values {
		if string(v) == {{ .ValueArgName }} {
			return {{ .TypeName }}({{ .ValueArgName }}), nil
		}
	}
	return "", errors.Errorf("'%s' is not a valid {{ .TypeNameReadable }}", s)
`

type FromString struct {
	name             string
	imports          golang.Imports
	receiver         golang.Value
	args             []golang.Value
	returns          []golang.Value
	Values           string
	TypeName         string
	TypeNameReadable string
	ValueArgName     string
}

// We can accept a string typeName rather than a data.Name because if the name of the type is not exported, the name of
// this function doesn't need to be either.
func NewFromString(typeName, typeNameReadable string, values []string) FromString {
	valueArgName := "s"
	var vv []string
	for _, v := range values {
		vv = append(vv, "\t\t"+fmt.Sprintf(`"%s",`, v))
	}
	return FromString{
		name: typeName + "FromString",
		imports: golang.Imports{
			Standard: nil,
			App:      nil,
			Vendor:   []string{data.ImportErrors},
		},
		receiver: golang.Value{},
		args: []golang.Value{
			{
				Name: valueArgName,
				Type: "string",
			},
		},
		returns: []golang.Value{
			{
				Type: typeName,
			},
			{
				Type: "error",
			},
		},
		Values:           strings.Join(vv, "\n"),
		TypeName:         typeName,
		TypeNameReadable: typeNameReadable,
		ValueArgName:     valueArgName,
	}
}

func (m FromString) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         m.name,
		Imports:      m.imports,
		Receiver:     m.receiver,
		Arguments:    m.args,
		ReturnValues: m.returns,
		Body:         m.MustParse(),
	}
}

func (m FromString) GetImports() golang.Imports {
	return m.imports
}

func (m FromString) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_enum_from_string", fromStringBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
