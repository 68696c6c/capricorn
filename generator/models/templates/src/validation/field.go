package validation

import (
	"strings"

	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/utils"
)

var fieldValidationsTemplate = `
validation.Field(&{{ .Receiver }}.{{ .Field.Name.Exported }}, {{ .GetRules }}),`

type Field struct {
	Receiver string
	Single   data.Name
	Field    module.ResourceField
}

func (m Field) GetRules() string {
	var rules []string
	if m.Field.IsRequired {
		rules = append(rules, "validation.Required")
	}
	if m.Field.IsUnique {
		rule := MakeUniqueRule("db", m.Single, m.Field.Name)
		rules = append(rules, rule.ConstructorName+"(db)")
	}
	return strings.Join(rules, ", ")
}

func (m Field) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_field_validations", fieldValidationsTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
