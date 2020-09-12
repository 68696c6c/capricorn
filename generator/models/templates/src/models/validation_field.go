package models

import (
	"strings"

	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/src/validation_rules/unique"
	"github.com/68696c6c/capricorn/generator/utils"
)

var fieldValidationsTemplate = `
validation.Field(&{{ .Receiver }}.{{ .Field.Name.Exported }}, {{ .GetRules }}),`

type ValidationField struct {
	Receiver string
	DB       string
	Single   data.Name
	Field    module.ResourceField
}

func (m ValidationField) GetRules() string {
	var rules []string
	if m.Field.IsRequired {
		rules = append(rules, "validation.Required")
	}
	if m.Field.IsUnique {
		rule := unique.NewRule(m.DB, m.Receiver, m.Single, m.Field)
		rules = append(rules, rule.GetConstructorCall())
	}
	return strings.Join(rules, ", ")
}

func (m ValidationField) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_field_validations", fieldValidationsTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
