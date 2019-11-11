package src

import (
	"fmt"

	"github.com/68696c6c/capricorn/generator/models"
	"github.com/68696c6c/capricorn/generator/utils"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const modelTemplate = `
package models

import "github.com/68696c6c/goat"

{{ $tick := "` + "`" + `" }}
type {{ .Name }} struct {
	goat.Model
	{{- range $key, $value := .Fields }}
	{{ $value.Name }} {{ $value.Type }} {{ $tick }}{{ $value.Tag }}{{ $tick }}
	{{- end }}
}
`

func CreateModels(spec models.Project, logger *logrus.Logger) error {
	logPrefix := "CreateModels | "

	err := utils.CreateDir(spec.Paths.Models)
	if err != nil {
		return errors.Wrapf(err, "failed to create models directory '%s'", spec.Paths.Models)
	}

	for _, m := range spec.Models {
		logger.Infof(logPrefix, fmt.Sprintf("creating model %s", m.Filename))
		err = utils.GenerateFile(spec.Paths.Models, m.Filename, modelTemplate, m)
		if err != nil {
			return errors.Wrap(err, "failed to generate model")
		}
	}

	return nil
}
