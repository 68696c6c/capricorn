package src

import (
	"github.com/68696c6c/capricorn/generator/models"
	"github.com/68696c6c/capricorn/generator/utils"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const containerTemplate = `
package app

import (
	"{{.Imports.Repos}}"

	"github.com/68696c6c/goat"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

var container ServiceContainer

type ServiceContainer struct {
	DB     *gorm.DB
	Logger *logrus.Logger
	Errors goat.ErrorHandler

{{- range $key, $value := .Repos }}
{{ $value.Interface }} repos.{{ $value.Interface }}
{{- end }}
}

func (a ServiceContainer) GetDB() *gorm.DB {
	return a.DB
}

func (a ServiceContainer) GetLogger() *logrus.Logger {
	return a.Logger
}

// Initializes the service container if it hasn't been already.
func GetApp(db *gorm.DB, logger *logrus.Logger) (ServiceContainer, error) {
	if container != (ServiceContainer{}) {
		return container, nil
	}

	container = ServiceContainer{
		DB:     db,
		Logger: logger,
		Errors: goat.NewErrorHandler(logger),
{{- range $key, $value := .Repos }}
{{ $value.Interface }}: repos.{{ $value.Constructor }}(db),
{{- end }}
	}

	return container, nil
}

`

func CreateApp(spec models.Project, logger *logrus.Logger) error {
	logPrefix := "CreateApp | "
	logger.Debug(logPrefix, "creating service container")
	//
	err := utils.CreateDir(spec.Paths.App)
	if err != nil {
		return errors.Wrapf(err, "failed to create app directory '%s'", spec.Paths.App)
	}

	// Create a service container.
	err = utils.GenerateFile(spec.Paths.App, "container.go", containerTemplate, spec)
	if err != nil {
		return errors.Wrap(err, "failed to create container")
	}

	return nil
}
