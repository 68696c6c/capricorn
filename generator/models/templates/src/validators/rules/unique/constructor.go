package unique

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/models/templates/src/validators/rules"
	"github.com/68696c6c/capricorn/generator/utils"
)

var ruleConstructorBodyTemplate = `
	return &{{ .RuleName }}{
		message: "{{ .Single.Space }} {{ .Field.Name.Space }} must be unique",
		{{ .DBFieldName }}:      {{ .DBArgName }},
	}
`

type constructor struct {
	name        string
	receiver    golang.Value
	imports     golang.Imports
	args        []golang.Value
	returns     []golang.Value
	RuleName    string
	Field       module.ResourceField
	Single      data.Name
	DBArgName   string
	DBFieldName string
}

func newConstructor(meta rules.RuleMeta) constructor {
	return constructor{
		name:     meta.ConstructorName,
		receiver: golang.Value{},
		imports: golang.Imports{
			Standard: []string{},
			App:      []string{},
			Vendor:   []string{data.ImportGorm},
		},
		args: []golang.Value{
			{
				Name: meta.DBArgName,
				Type: "*gorm.DB",
			},
		},
		returns: []golang.Value{
			{
				Type: meta.Receiver.Type,
			},
		},
		RuleName:    meta.RuleName,
		Field:       meta.Field,
		Single:      meta.Single,
		DBArgName:   meta.DBArgName,
		DBFieldName: meta.DBFieldName,
	}
}

func (m constructor) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         m.name,
		Imports:      m.imports,
		Receiver:     m.receiver,
		Arguments:    m.args,
		ReturnValues: m.returns,
		Body:         m.MustParse(),
	}
}

func (m constructor) GetImports() golang.Imports {
	return m.imports
}

func (m constructor) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_rule_unique_constructor_body", ruleConstructorBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
