package golang

import "github.com/68696c6c/capricorn/generator/utils"

var structTemplate = `type {{ .Name }} struct {
	{{- range $key, $value := .Fields }}
	{{ $value.MustParse }}
	{{- end }}
}`

type Struct struct {
	Name   string  `yaml:"name"`
	Fields []Field `yaml:"fields"`
}

func (m Struct) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_struct", structTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
