package golang

import "github.com/68696c6c/capricorn/generator/utils"

var varTemplate = `var {{ .Name }} = {{ .Value }}`

// Represents a global var declaration.
type Var struct {
	Name  string `yaml:"name,omitempty"`
	Value string `yaml:"value,omitempty"`
}

func (m Var) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_var", varTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
