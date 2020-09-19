package module

import (
	"testing"

	"github.com/68696c6c/capricorn/generator/models/spec"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestModule_ModuleFromSpec(t *testing.T) {
	s := spec.Spec{}
	err := yaml.Unmarshal(spec.GetFixtureInput(), &s)
	require.Nil(t, err)

	result := NewModuleFromSpec(s, true)
	resultYAML := result.String()
	println(resultYAML)

	assert.Equal(t, spec.FixtureSpecYAML, resultYAML)
}
