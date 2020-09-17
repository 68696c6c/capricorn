package unique

import (
	"fmt"

	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/models/templates/src/validators/rules"
	"github.com/68696c6c/capricorn/generator/utils"
)

var ruleValidateBodyTemplate = `
	{{ .Field.Name.Unexported }}, ok := value.({{ .Field.Type }})
	if !ok {
		return errors.New("invalid {{ .Single.Space }} {{ .Field.Name.Space }}")
	}

	query := {{ .GetDbReference }}.First(&{{ .Single.Exported }}{
		{{ .Field.Name.Exported }}: {{ .Field.Name.Unexported }},
	})
	if !query.RecordNotFound() {
		return errors.New("{{ .Single.Space }} {{ .Field.Name.Space }} already exists")
	}

	return nil
`

type validate struct {
	name        string
	receiver    golang.Value
	imports     golang.Imports
	args        []golang.Value
	returns     []golang.Value
	dbFieldName string
	Field       module.ResourceField
	Single      data.Name
}

func newValidate(meta rules.RuleMeta) validate {
	return validate{
		name:     "Validate",
		receiver: meta.Receiver,
		imports: golang.Imports{
			Standard: []string{},
			App:      []string{},
			Vendor:   []string{data.ImportErrors, data.ImportGorm},
		},
		args: []golang.Value{
			{
				Name: meta.DBArgName,
				Type: "*gorm.DB",
			},
		},
		returns: []golang.Value{
			{
				Type: "error",
			},
		},
		Field:       meta.Field,
		Single:      meta.Single,
		dbFieldName: meta.DBFieldName,
	}
}

func (m validate) GetDbReference() string {
	return fmt.Sprintf("%s.%s", m.receiver.Name, m.dbFieldName)
}

func (m validate) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         m.name,
		Imports:      m.imports,
		Receiver:     m.receiver,
		Arguments:    m.args,
		ReturnValues: m.returns,
		Body:         m.MustParse(),
	}
}

func (m validate) GetImports() golang.Imports {
	return m.imports
}

func (m validate) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_rule_unique_validate_body", ruleValidateBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
