package methods

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var initFunctionBodyTemplate = `
	goose.AddMigration({{ .UpName }}, {{ .DownName }})
`

type InitFunction struct {
	name     string
	receiver golang.Value
	imports  golang.Imports
	args     []golang.Value
	returns  []golang.Value
	UpName   string
	DownName string
}

func NewInitFunction(upName, downName string) InitFunction {
	return InitFunction{
		name:     "init",
		receiver: golang.Value{},
		imports: golang.Imports{
			Standard: nil,
			App:      nil,
			Vendor:   []string{data.ImportGoose},
		},
		args:     []golang.Value{},
		returns:  []golang.Value{},
		UpName:   upName,
		DownName: downName,
	}
}

func (m InitFunction) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         m.name,
		Imports:      m.imports,
		Receiver:     m.receiver,
		Arguments:    m.args,
		ReturnValues: m.returns,
		Body:         m.MustParse(),
	}
}

func (m InitFunction) GetImports() golang.Imports {
	return m.imports
}

func (m InitFunction) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_initial_migration_init_function", initFunctionBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
