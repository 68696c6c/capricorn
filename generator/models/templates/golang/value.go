package golang

import (
	"strings"

	"github.com/68696c6c/capricorn/generator/utils"
)

var valueTemplate = `{{ .Name }} {{ .TypeData }}`

// Represents an argument or return value.
type Value struct {
	Name string `yaml:"name,omitempty"`
	Type string `yaml:"type,omitempty"`
}

func (m Value) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_value", valueTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}

func getJoinedValueString(values []Value) string {
	var builtValues []string
	for _, v := range values {
		builtValues = append(builtValues, v.MustParse())
	}
	joinedValues := strings.Join(builtValues, ", ")
	return strings.TrimSpace(joinedValues)
}
