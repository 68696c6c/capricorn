package module

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/spec"
)

type Command struct {
	_spec spec.Command
	Name  data.Name `yaml:"name,omitempty"`
}

func makeCommands(specCommands []spec.Command) []Command {
	var result []Command
	for _, c := range specCommands {
		cmd := Command{
			_spec: c,
			Name:  data.MakeName(c.Name),
		}
		result = append(result, cmd)
	}
	return result
}
