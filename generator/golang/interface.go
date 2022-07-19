package golang

import "github.com/68696c6c/capricorn_rnd/generator/utils"

var interfaceTemplate = `type {{ .Name }} interface {
	{{- range $key, $value := .Functions }}
	{{ $value.GetSignature }}
	{{- end }}
}`

type Interface struct {
	Name      string
	Functions []*Function
}

func (i *Interface) MustString() string {
	result, err := utils.ParseTemplateToString("tmp_template_interface", interfaceTemplate, i)
	if err != nil {
		panic(err)
	}
	return result
}
