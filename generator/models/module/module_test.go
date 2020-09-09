package module

import (
	"testing"

	"github.com/68696c6c/capricorn/generator/models/spec"

	"github.com/stretchr/testify/assert"
)

func TestModule_ModuleFromSpec(t *testing.T) {
	f := spec.GetFixtureSpec()

	result := NewModuleFromSpec(f, true)
	resultYAML := result.String()
	println(resultYAML)

	assert.Equal(t, spec.FixtureSpecYAML, resultYAML)
}
