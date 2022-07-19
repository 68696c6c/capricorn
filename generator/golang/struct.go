package golang

import "github.com/68696c6c/capricorn_rnd/generator/utils"

var structTemplate = `type {{ .Name }} struct {
	{{- range $key, $value := .Fields }}
	{{ $value.MustString }}
	{{- end }}
}`

type Struct struct {
	Name   string
	Fields []*Field
}

func (s *Struct) MustString() string {
	result, err := utils.ParseTemplateToString("tmp_template_struct", structTemplate, s)
	if err != nil {
		panic(err)
	}
	return result
}
