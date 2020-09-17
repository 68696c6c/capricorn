package models

import (
	"strings"

	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/src/models/validation_rules"
	"github.com/68696c6c/capricorn/generator/models/templates/src/models/validation_rules/required"
	"github.com/68696c6c/capricorn/generator/models/templates/src/models/validation_rules/unique"
	"github.com/68696c6c/capricorn/generator/utils"
)

var fieldValidationsTemplate = `
validation.Field(&{{ .ReceiverName }}.{{ .Field.Name.Exported }}, {{ .RenderRules }}),`

type ValidationField struct {
	dbFieldName  string
	single       data.Name
	rules        []validation_rules.Rule
	ReceiverName string
	Field        module.ResourceField

	built bool
}

func NewValidationField(meta ValidationMeta, field module.ResourceField) *ValidationField {
	return &ValidationField{
		dbFieldName:  meta.DBFieldName,
		single:       meta.ModelName,
		ReceiverName: meta.Receiver.Name,
		Field:        field,
	}
}

func (m *ValidationField) RenderRules() string {
	if !m.built {
		m.build()
	}
	var rules []string
	for _, r := range m.rules {
		rules = append(rules, r.GetUsage())
	}
	return strings.Join(rules, ", ")
}

func (m *ValidationField) GetRules() []validation_rules.Rule {
	if !m.built {
		m.build()
	}
	return m.rules
}

func (m *ValidationField) build() {
	if m.built {
		return
	}

	var rules []validation_rules.Rule

	if m.Field.IsRequired {
		rules = append(rules, required.NewRule())
	}

	if m.Field.IsUnique {
		rule := unique.NewRule(m.dbFieldName, m.ReceiverName, m.single, m.Field)
		rules = append(rules, rule)
	}

	m.rules = rules
	m.built = true
}

func (m *ValidationField) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_field_validations", fieldValidationsTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
