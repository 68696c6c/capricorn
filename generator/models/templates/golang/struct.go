package golang

import "github.com/68696c6c/capricorn/generator/utils"

var structTemplate = `type {{ .Name }} struct {
	{{- range $key, $value := .Fields }}
	{{ $value.MustParse }}
	{{- end }}
}`

type Struct struct {
	Name   string  `yaml:"name,omitempty"`
	Fields []Field `yaml:"fields,omitempty"`
}

func (m Struct) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_struct", structTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
