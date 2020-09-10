package validation

import (
	"fmt"

	"github.com/68696c6c/capricorn/generator/models/data"
)

func makeRuleConstructorName(singleName, fieldName data.Name) string {
	return fmt.Sprintf("new%s%sUniqueRule", singleName.Exported, fieldName.Exported)
}
