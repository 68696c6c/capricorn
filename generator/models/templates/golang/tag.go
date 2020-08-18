package golang

import (
	"strings"

	"github.com/68696c6c/capricorn/generator/utils"
)

var structFieldTagTemplate = `{{ .Key }}:"{{ .GetValues }}"`

type Tag struct {
	Key    string   `yaml:"key"`    // e.g. yaml
	Values []string `yaml:"values"` // e.g. "name", "omitempty"
}

func (m Tag) GetValues() string {
	var builtValues []string
	for _, v := range m.Values {
		builtValues = append(builtValues, v)
	}
	joinedValues := strings.Join(builtValues, ",")
	return strings.TrimSpace(joinedValues)
}

func (m Tag) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_tag", structFieldTagTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
