package golang

import "github.com/68696c6c/capricorn/generator/utils"

var constTemplate = `const {{ .Name }} = {{ .Value }}`

// Represents a global const declaration.
type Const struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

func (m Const) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_const", constTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
