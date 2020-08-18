package golang

import (
	"fmt"
	"strings"

	"github.com/68696c6c/capricorn/generator/utils"
)

var structFieldTemplate = `{{ .Name }} {{ .Type }}{{ .GetTags }}`

type Field struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Tags []Tag  `yaml:"tags"`
}

func (m Field) GetTags() string {
	var builtValues []string
	for _, v := range m.Tags {
		tagString := v.MustParse()
		builtValues = append(builtValues, tagString)
	}
	if len(builtValues) == 0 {
		return ""
	}
	joinedValues := strings.Join(builtValues, " ")
	tags := strings.TrimSpace(joinedValues)
	return fmt.Sprintf(" `%s`", tags)
}

func (m Field) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_field", structFieldTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
