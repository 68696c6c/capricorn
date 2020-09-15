package src

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainGo(t *testing.T) {
	input := NewMainGo("/root/path", "github.com/68696c6c/test-example", "github.com/68696c6c/test-example/cmd")

	result := input.MustParse()

	expected := `package main

import (
	"os"

	"github.com/68696c6c/test-example/cmd"
)


func main() {
	if err := cmd.Root.Execute(); err != nil {
		println(err)
		os.Exit(1)
	}
}
`

	assert.Equal(t, expected, result)
}
