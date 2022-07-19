package golang

import (
	"strings"

	"github.com/68696c6c/capricorn_rnd/generator/utils"
)

var valueTemplate = `{{ .Name }} {{ .Type }}`

// Represents an argument or return value.
type Value struct {
	Name string
	Type string
}

func (v Value) MustString() string {
	result, err := utils.ParseTemplateToString("tmp_template_value", valueTemplate, v)
	if err != nil {
		panic(err)
	}
	return result
}

func getJoinedValueString(values []Value) string {
	var builtValues []string
	for _, v := range values {
		builtValues = append(builtValues, v.MustString())
	}
	joinedValues := strings.Join(builtValues, ", ")
	return strings.TrimSpace(joinedValues)
}
