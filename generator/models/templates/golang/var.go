package golang

import "github.com/68696c6c/capricorn/generator/utils"

var varTemplate = `var {{ .Name }} = {{ .Value }}`

// Represents a global var declaration.
type Var struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

func (m Var) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_var", varTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
