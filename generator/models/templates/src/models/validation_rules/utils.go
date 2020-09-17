package validation_rules

import (
	"fmt"

	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
)

type Rule interface {
	GetUsage() string
	GetStructs() []golang.Struct
	MustGetFunctions() []golang.Function
}

func MakeRuleName(resourceSingleName data.Name, field module.ResourceField, ruleType string) (ruleName, constructorName string) {
	base := fmt.Sprintf("%s-%s-%s-rule", resourceSingleName.Kebob, field.Name.Kebob, ruleType)
	baseName := data.MakeName(base)
	return baseName.Unexported, "new" + baseName.Exported
}
