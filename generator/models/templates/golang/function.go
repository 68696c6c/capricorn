package golang

import (
	"fmt"
	"strings"

	"github.com/68696c6c/capricorn/generator/utils"
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

func (m Function) GetReceiver() string {
	r := fmt.Sprintf("%s %s", m.Receiver.Name, m.Receiver.Type)
	r = strings.TrimSpace(r)
	if r != "" {
		return fmt.Sprintf("(%s) ", r)
	}
	return ""
}

func (m Function) GetSignature() string {
	args := getJoinedValueString(m.Arguments)
	returns := getJoinedValueString(m.ReturnValues)

	if len(m.ReturnValues) > 0 {
		for _, v := range m.ReturnValues {
			if v.Name != "" {
				returns = fmt.Sprintf("(%s)", returns)
				break
			}
		}
	}

	result := fmt.Sprintf("%s(%s) %s", m.Name, args, returns)

	return strings.TrimSpace(result)
}

func (m Function) GetBody() string {
	result := strings.TrimLeft(m.Body, "\n")
	result = strings.TrimLeft(result, "\t")
	result = fmt.Sprintf("\t%s", result)
	return result
}

func (m Function) MustParse() string {
	if m.Name == "" {
		return ""
	}
	result, err := utils.ParseTemplateToString("tmp_template_function", functionTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
