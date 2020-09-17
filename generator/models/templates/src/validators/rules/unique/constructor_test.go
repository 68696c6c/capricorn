package unique

import (
	"testing"

	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/src/validators/rules"

	"github.com/stretchr/testify/assert"
)

func TestConstructor_MustParse(t *testing.T) {
	recNameKebob := "example-resource"
	fNameKebob := "example-field"
	recName := data.MakeName(recNameKebob)

	field := module.GetFixtureResourceField(recNameKebob, fNameKebob)
	field.IsRequired = true
	field.IsUnique = true

	c := newConstructor(rules.RuleMeta{
		RuleName:        "ruleName",
		ConstructorName: "newRuleName",
		DBArgName:       "d",
		DBFieldName:     "db",
		Single:          recName,
		Field:           field,
	})
	result := c.MustParse()

	assert.Equal(t, fixtureConstructor, result)
}

const fixtureConstructor = `
	return &ruleName{
		message: "example resource example field must be unique",
		db:      d,
	}
`
