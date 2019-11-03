package project

import "github.com/pkg/errors"

const containerTemplate = `
package app

import (
	"github.com/68696c6c/goat"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

var container ServiceContainer

type ServiceContainer struct {
	DB     *gorm.DB
	Logger *logrus.Logger
}

func (a ServiceContainer) GetDB() *gorm.DB {
	return a.DB
}

func (a ServiceContainer) GetLogger() *logrus.Logger {
	return a.Logger
}

// Initializes the service container if it hasn't been already.
func GetApp() (goat.App, error) {
	if container != (ServiceContainer{}) {
		return container, nil
	}

	logger := goat.NewSTDOutLogger()

	db, err := goat.GetMainDB()
	if err != nil {
		return ServiceContainer{}, err
	}

	container = ServiceContainer{
		DB:     db,
		Logger: logger,
	}

	return container, nil
}

`

func CreateApp(spec Spec) error {
	err := createDir(spec.Paths.App)
	if err != nil {
		return errors.Wrapf(err, "failed to create app directory '%s'", spec.Paths.App)
	}

	// Create a service container.
	err = generateFile(spec.Paths.App, "container", containerTemplate, spec)
	if err != nil {
		return errors.Wrap(err, "failed to create container")
	}

	return nil
}
