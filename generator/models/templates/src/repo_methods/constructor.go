package repo_methods

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
	interfaceName string
	name          string
	StructName    string
	DBArgName     string
	DBFieldName   string
}

func NewConstructor(interfaceName, implementationName string) Constructor {
	return Constructor{
		name:          "New" + interfaceName,
		interfaceName: interfaceName,
		StructName:    implementationName,
		DBArgName:     "d",
		DBFieldName:   "db",
	}
}

func (m Constructor) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         m.GetName(),
		Imports:      m.GetImports(),
		Receiver:     m.GetReceiver(),
		Arguments:    m.GetArgs(),
		ReturnValues: m.GetReturns(),
		Body:         m.MustParse(),
	}
}

func (m Constructor) GetName() string {
	return m.name
}

func (m Constructor) GetImports() golang.Imports {
	return golang.Imports{
		Standard: nil,
		App:      nil,
		Vendor:   []string{data.ImportGorm},
	}
}

func (m Constructor) GetReceiver() golang.Value {
	return golang.Value{}
}

func (m Constructor) GetArgs() []golang.Value {
	return []golang.Value{
		{
			Name: m.DBArgName,
			Type: "*gorm.DB",
		},
	}
}

func (m Constructor) GetReturns() []golang.Value {
	return []golang.Value{
		{
			Type: m.interfaceName,
		},
	}
}

func (m Constructor) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_repo_constructor", constructorBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
