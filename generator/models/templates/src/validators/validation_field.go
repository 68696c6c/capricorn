package validators

import (
	"strings"

	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/src/validators/rules"
	"github.com/68696c6c/capricorn/generator/models/templates/src/validators/rules/required"
	"github.com/68696c6c/capricorn/generator/models/templates/src/validators/rules/unique"
	"github.com/68696c6c/capricorn/generator/utils"
)

var fieldValidationsTemplate = `
validation.Field(&{{ .ReceiverName }}.{{ .Field.Name.Exported }}, {{ .RenderRules }}),`

type ValidationField struct {
	dbFieldName  string
	single       data.Name
	rules        []rules.Rule
	ReceiverName string
	Field        module.ResourceField

	built bool
}

func NewValidationField(meta ValidationMeta, field module.ResourceField) *ValidationField {
	return &ValidationField{
		dbFieldName:  meta.DBFieldName,
		single:       meta.ModelName,
		ReceiverName: meta.ReceiverName,
		Field:        field,
	}
}

func (m *ValidationField) RenderRules() string {
	if !m.built {
		m.build()
	}
	var fieldRules []string
	for _, r := range m.rules {
		fieldRules = append(fieldRules, r.GetUsage())
	}
	return strings.Join(fieldRules, ", ")
}

func (m *ValidationField) GetRules() []rules.Rule {
	if !m.built {
		m.build()
	}
	return m.rules
}

func (m *ValidationField) build() {
	if m.built {
		return
	}

	var fieldRules []rules.Rule

	if m.Field.IsRequired {
		fieldRules = append(fieldRules, required.NewRule())
	}

	if m.Field.IsUnique {
		rule := unique.NewRule(m.dbFieldName, m.ReceiverName, m.single, m.Field)
		fieldRules = append(fieldRules, rule)
	}

	m.rules = fieldRules
	m.built = true
}

func (m *ValidationField) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_field_validations", fieldValidationsTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
