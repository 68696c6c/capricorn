package repos

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var constructorBodyTemplate = `
	return {{ .StructName }}{
		{{ .DBFieldName }}: {{ .DBArgName }},
	}
`

type Constructor struct {
	name        string
	receiver    golang.Value
	imports     golang.Imports
	args        []golang.Value
	returns     []golang.Value
	StructName  string
	DBArgName   string
	DBFieldName string
}

func NewConstructor(interfaceName, implementationName, dbFieldName string) Constructor {
	dbArgName := "d"
	return Constructor{
		name:     "New" + interfaceName,
		receiver: golang.Value{},
		imports: golang.Imports{
			Standard: nil,
			App:      nil,
			Vendor:   []string{data.ImportGorm},
		},
		args: []golang.Value{
			{
				Name: dbArgName,
				Type: "*gorm.DB",
			},
		},
		returns: []golang.Value{
			{
				Type: interfaceName,
			},
		},
		StructName:  implementationName,
		DBArgName:   dbArgName,
		DBFieldName: dbFieldName,
	}
}

func (m Constructor) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         m.name,
		Imports:      m.imports,
		Receiver:     m.receiver,
		Arguments:    m.args,
		ReturnValues: m.returns,
		Body:         m.MustParse(),
	}
}

func (m Constructor) GetImports() golang.Imports {
	return m.imports
}

func (m Constructor) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_repo_constructor", constructorBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
