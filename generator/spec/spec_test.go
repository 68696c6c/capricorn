package spec

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSpec_Unmarshal(t *testing.T) {
	result := Spec{}
	err := yaml.Unmarshal(GetFixtureInput(), &result)
	require.Nil(t, err)
	fixture := GetFixtureSpec()
	assert.Equal(t, fmt.Sprintf("%v", fixture), fmt.Sprintf("%v", result))
}
