package validation_rules

import (
	"testing"

	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"

	"github.com/stretchr/testify/assert"
)

func TestMakeRuleName(t *testing.T) {
	recNameKebob := "example-resource"
	fNameKebob := "example-field"
	recName := data.MakeName(recNameKebob)
	field := module.GetFixtureResourceField(recNameKebob, fNameKebob)

	ruleResult, constResult := MakeRuleName(recName, field, "unique")

	assert.Equal(t, "exampleResourceExampleFieldUniqueRule", ruleResult)
	assert.Equal(t, "newExampleResourceExampleFieldUniqueRule", constResult)
}
