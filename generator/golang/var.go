package golang

import "github.com/68696c6c/capricorn_rnd/generator/utils"

var varTemplate = `var {{ .Name }}{{ .MustStringType }}{{ .MustStringValue }}`

// Represents a global var declaration.
type Var struct {
	Name  string
	Type  string
	Value string
}

func (v *Var) MustStringType() string {
	if v.Type != "" {
		return " " + v.Type
	}
	return ""
}

func (v *Var) MustStringValue() string {
	if v.Value != "" {
		return " = " + v.Value
	}
	return ""
}

func (v *Var) MustString() string {
	result, err := utils.ParseTemplateToString("tmp_template_var", varTemplate, v)
	if err != nil {
		panic(err)
	}
	return result
}
