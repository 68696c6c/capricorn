package required

import (
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/models/templates/src/validators/rules"
)

type Required struct{}

func NewRule() rules.Rule {
	return Required{}
}

func (m Required) GetUsage() string {
	return "validation.Required"
}

func (m Required) GetStructs() []golang.Struct {
	return []golang.Struct{}
}

func (m Required) MustGetFunctions() []golang.Function {
	return []golang.Function{}
}
