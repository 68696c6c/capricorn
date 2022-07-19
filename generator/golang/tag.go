package golang

import (
	"strings"

	"github.com/68696c6c/capricorn_rnd/generator/utils"
)

var structFieldTagTemplate = `{{ .Key }}:"{{ .GetValues }}"`

type Tag struct {
	Key    string   // e.g. json
	Values []string // e.g. "name", "omitempty"
}

func (m *Tag) GetValues() string {
	var builtValues []string
	for _, v := range m.Values {
		builtValues = append(builtValues, v)
	}
	joinedValues := strings.Join(builtValues, ",")
	return strings.TrimSpace(joinedValues)
}

func (m *Tag) MustString() string {
	result, err := utils.ParseTemplateToString("tmp_template_tag", structFieldTagTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
