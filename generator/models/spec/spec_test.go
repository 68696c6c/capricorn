package spec

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSpec(t *testing.T) {
	path, err := os.Getwd()
	require.Nil(t, err)
	println(path)
	result, err := NewSpec("tmp.yml")
	require.Nil(t, err)
	resultYAML := result.String()
	println(resultYAML)
	// for _, e := range result.Enums {
	// 	for key, value := range e {
	// 		println(key)
	// 		println(value.Description)
	// 	}
	// }
	assert.True(t, false)
}
