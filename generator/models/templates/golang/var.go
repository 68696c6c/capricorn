package golang

import "github.com/68696c6c/capricorn/generator/utils"

var varTemplate = `var {{ .Name }}{{ .MustParseType }}{{ .MustParseValue }}`

// Represents a global var declaration.
type Var struct {
	Name  string `yaml:"name,omitempty"`
	Type  string `yaml:"type,omitempty"`
	Value string `yaml:"value,omitempty"`
}

func (m Var) MustParseType() string {
	if m.Type != "" {
		return " " + m.Type
	}
	return ""
}

func (m Var) MustParseValue() string {
	if m.Value != "" {
		return " = " + m.Value
	}
	return ""
}

func (m Var) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_var", varTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
