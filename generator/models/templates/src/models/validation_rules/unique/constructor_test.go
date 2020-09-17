package unique

import (
	"testing"

	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"

	"github.com/stretchr/testify/assert"
)

func TestConstructor_MustParse(t *testing.T) {
	recNameKebob := "example-resource"
	fNameKebob := "example-field"
	recName := data.MakeName(recNameKebob)

	field := module.GetFixtureResourceField(recNameKebob, fNameKebob)
	field.IsRequired = true
	field.IsUnique = true

	input := constructor{
		RuleName: "ruleName",
		Name:     "NewRule",
		DB:       "db",
		Field:    field,
		Single:   recName,
	}
	result := input.MustParse()

	assert.Equal(t, fixtureConstructor, result)
}

const fixtureConstructor = `
	return &ruleName{
		message: "example resource example field must be unique",
		db:      db,
	}
`
