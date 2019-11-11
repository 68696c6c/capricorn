package src

import (
	"github.com/68696c6c/capricorn/generator/models"
	"github.com/68696c6c/capricorn/generator/utils"

	"github.com/pkg/errors"
)

const mainTemplate = `
package main

import (
	"fmt"
	"os"

	"{{.Imports.CMD}}"
)

func main() {
	if err := cmd.Root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
`

func CreateProject(spec models.Project) error {
	err := utils.CreateDir(spec.Paths.Root)
	if err != nil {
		return errors.Wrapf(err, "failed to create project directory '%s'", spec.Paths.Root)
	}

	err = utils.GenerateFile(spec.Paths.Root, "main.go", mainTemplate, spec)
	if err != nil {
		return errors.Wrap(err, "failed to create main.go")
	}
	return nil
}
