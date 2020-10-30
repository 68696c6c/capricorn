package methods

import (
	"fmt"
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
	"strings"
)

var upBodyTemplate = `
	goat.Init()

	{{ .DBFieldName }}, err := goat.GetMigrationDB()
	if err != nil {
		return errors.Wrap(err, "failed to initialize migration connection")
	}
	{{ .MustParseTables }}

	return nil
`

type Up struct {
	Name     string
	receiver golang.Value
	imports  golang.Imports
	args     []golang.Value
	returns  []golang.Value

	DBFieldName string
	models      []string
}

func NewUp(appImports, modelRefs []string) Up {
	return Up{
		Name:     "upInitialMigration",
		receiver: golang.Value{},
		imports: golang.Imports{
			Standard: []string{"database/sql"},
			App:      appImports,
			Vendor:   []string{data.ImportGoat, data.ImportErrors},
		},
		args: []golang.Value{
			{
				Name: "tx",
				Type: "*sql.Tx",
			},
		},
		returns: []golang.Value{
			{
				Type: "error",
			},
		},
		DBFieldName: "db",
		models:      modelRefs,
	}
}

func (m Up) MustParseTables() string {
	var result []string
	for _, ref := range m.models {
		result = append(result, fmt.Sprintf("%s.AutoMigrate(&%s{})", m.DBFieldName, ref))
	}
	return strings.Join(result, "\n")
}

func (m Up) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         m.Name,
		Imports:      m.imports,
		Receiver:     m.receiver,
		Arguments:    m.args,
		ReturnValues: m.returns,
		Body:         m.MustParse(),
	}
}

func (m Up) GetImports() golang.Imports {
	return m.imports
}

func (m Up) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_initial_migration_up", upBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
