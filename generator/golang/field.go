package golang

import (
	"fmt"
	"strings"

	"github.com/68696c6c/capricorn_rnd/generator/utils"
)

var structFieldTemplate = `{{ .Name }} {{ .Type }}{{ .GetTags }}`

type Field struct {
	Name string
	Type string
	Tags []*Tag
}

func (m *Field) AddTag(t *Tag) *Field {
	m.Tags = append(m.Tags, t)
	return m
}

func (m *Field) GetTags() string {
	var builtValues []string
	for _, v := range m.Tags {
		tagString := v.MustString()
		builtValues = append(builtValues, tagString)
	}
	if len(builtValues) == 0 {
		return ""
	}
	joinedValues := strings.Join(builtValues, " ")
	tags := strings.TrimSpace(joinedValues)
	return fmt.Sprintf(" `%s`", tags)
}

func (m *Field) MustString() string {
	result, err := utils.ParseTemplateToString("tmp_template_field", structFieldTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
