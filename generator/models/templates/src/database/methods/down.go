package methods

import (
	"fmt"
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
	"strings"
)

var downBodyTemplate = `
	goat.Init()

	{{ .DBFieldName }}, err := goat.GetMigrationDB()
	if err != nil {
		return errors.Wrap(err, "failed to initialize migration connection")
	}
	{{ .MustParseTables }}

	return nil
`

type Down struct {
	Name     string
	receiver golang.Value
	imports  golang.Imports
	args     []golang.Value
	returns  []golang.Value

	DBFieldName string
	models      []string
}

func NewDown(appImports, modelRefs []string) Down {
	return Down{
		Name:     "downInitialMigration",
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

func (m Down) MustParseTables() string {
	var result []string
	for _, ref := range m.models {
		result = append(result, fmt.Sprintf("%s.DropTable(&%s{})", m.DBFieldName, ref))
	}
	return strings.Join(result, "\n")
}

func (m Down) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         m.Name,
		Imports:      m.imports,
		Receiver:     m.receiver,
		Arguments:    m.args,
		ReturnValues: m.returns,
		Body:         m.MustParse(),
	}
}

func (m Down) GetImports() golang.Imports {
	return m.imports
}

func (m Down) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_initial_migration_down", downBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
