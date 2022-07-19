package golang

import (
	"fmt"
	"strings"

	"github.com/68696c6c/capricorn_rnd/generator/utils"
)

var functionTemplate = `func {{ .GetReceiver }}{{ .GetSignature }} {
{{ .GetBody }}
}`

type Function struct {
	Name         string  `yaml:"name,omitempty"`
	Imports      Imports `yaml:"imports,omitempty"` // Any imports that this function requires.
	Arguments    []Value `yaml:"arguments,omitempty"`
	ReturnValues []Value `yaml:"return_values,omitempty"`
	Receiver     Value   `yaml:"receiver,omitempty"`
	Body         string  `yaml:"body,omitempty"` // The actual function code.
}

func (f *Function) GetReceiver() string {
	r := fmt.Sprintf("%s %s", f.Receiver.Name, f.Receiver.Type)
	r = strings.TrimSpace(r)
	if r != "" {
		return fmt.Sprintf("(%s) ", r)
	}
	return ""
}

func (f *Function) hasNamedReturn() bool {
	if len(f.ReturnValues) > 0 {
		for _, v := range f.ReturnValues {
			if v.Name != "" {
				return true
			}
		}
	}
	return false
}

func (f *Function) GetSignature() string {
	args := getJoinedValueString(f.Arguments)
	returns := getJoinedValueString(f.ReturnValues)

	if len(f.ReturnValues) > 1 || f.hasNamedReturn() {
		returns = fmt.Sprintf("(%s)", returns)
	}

	result := fmt.Sprintf("%s(%s) %s", f.Name, args, returns)

	return strings.TrimSpace(result)
}

func (f *Function) GetBody() string {
	result := strings.TrimLeft(f.Body, "\n")
	result = strings.TrimLeft(result, "\t")
	result = fmt.Sprintf("\t%s", result)
	return result
}

func (f *Function) MustString() string {
	if f.Name == "" {
		return ""
	}
	result, err := utils.ParseTemplateToString("tmp_template_function", functionTemplate, f)
	if err != nil {
		panic(err)
	}
	return result
}
