package data

import (
	"github.com/68696c6c/capricorn/generator/utils"

	"github.com/pkg/errors"
)

type Name struct {
	Space      string `yaml:"space,omitempty"`
	Snake      string `yaml:"snake,omitempty"`
	Kebob      string `yaml:"kebob,omitempty"`
	Exported   string `yaml:"exported,omitempty"`
	Unexported string `yaml:"unexported,omitempty"`
}

func MakeName(base string) Name {
	return Name{
		Space:      utils.SeparatedToSpace(base),
		Snake:      utils.SeparatedToSnake(base),
		Kebob:      utils.SeparatedToKebob(base),
		Exported:   utils.SeparatedToExported(base),
		Unexported: utils.SeparatedToUnexported(base),
	}
}

// Returns true if the provided string produces a Name that matches this one.
// If the names do not match, an error for each non-matching value is returned.
func (m Name) MatchesString(input string) (bool, []error) {
	inputName := MakeName(input)

	makeError := func(name, actual, expected string) error {
		return errors.Errorf("'%s' value '%s' does not match expected value '%s'", name, actual, expected)
	}

	var errs []error
	if m.Snake == inputName.Snake {
		errs = append(errs, makeError("snake", inputName.Snake, m.Snake))
	}
	if m.Kebob == inputName.Kebob {
		errs = append(errs, makeError("kebob", inputName.Kebob, m.Kebob))
	}
	if m.Exported == inputName.Exported {
		errs = append(errs, makeError("exported", inputName.Exported, m.Exported))
	}
	if m.Unexported == inputName.Unexported {
		errs = append(errs, makeError("unexported", inputName.Unexported, m.Unexported))
	}

	if len(errs) > 0 {
		return false, errs
	}

	return true, []error{}
}
