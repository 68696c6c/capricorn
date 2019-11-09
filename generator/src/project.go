package src

import (
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

func CreateProject(spec utils.Spec) error {
	err := utils.CreateDir(spec.Paths.Root)
	if err != nil {
		return errors.Wrapf(err, "failed to create project directory '%s'", spec.Paths.Root)
	}

	err = utils.GenerateGoFile(spec.Paths.Root, "main", mainTemplate, spec)
	if err != nil {
		return errors.Wrap(err, "failed to create docker-compose.yml")
	}
	return nil
}
