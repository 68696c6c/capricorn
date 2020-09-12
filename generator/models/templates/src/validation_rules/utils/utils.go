package utils

import (
	"fmt"

	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
)

func MakeRuleName(resourceSingleName data.Name, field module.ResourceField) (ruleName, constructorName string) {
	base := fmt.Sprintf("%s-%s-rule", resourceSingleName.Snake, field.Name.Snake)
	baseName := data.MakeName(base)
	return baseName.Unexported, "new" + baseName.Exported
}
