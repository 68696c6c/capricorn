package src

import (
	"testing"

	"github.com/68696c6c/capricorn/generator/models/module"

	"github.com/stretchr/testify/assert"
)

func TestMainGo(t *testing.T) {
	f := module.GetFixtureModule()
	input := NewMainGo(f.Packages, "/root/path", "base/module")

	result := input.MustParse()
	println(result)

	expected := `package main

import (
	"os"

	"github.com/68696c6c/test-example/src/cmd"
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
