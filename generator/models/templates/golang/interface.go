package golang

import "github.com/68696c6c/capricorn/generator/utils"

var interfaceTemplate = `type {{ .Name }} interface {
	{{- range $key, $value := .Functions }}
	{{ $value.GetSignature }}
	{{- end }}
}`

type Interface struct {
	Name      string     `yaml:"name"`
	Functions []Function `yaml:"functions"`
}

func (m Interface) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_interface", interfaceTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
