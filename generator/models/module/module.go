package module

import (
	"path/filepath"

	"github.com/68696c6c/capricorn/generator/models"
	"github.com/68696c6c/capricorn/generator/models/spec"
	"github.com/68696c6c/capricorn/generator/models/templates/ops"
	"github.com/68696c6c/capricorn/generator/models/utils"

	"gopkg.in/yaml.v2"
)

type Module struct {
	_spec spec.Spec

	Name models.Name    `yaml:"name"`
	Path utils.PathData `yaml:"path"`
	Ops  ops.Ops        `yaml:"ops"`

	Packages Packages `yaml:"packages"`

	Commands  []Command  `yaml:"commands,omitempty"`
	Resources []Resource `yaml:"resources,omitempty"`
}

func (m Module) String() string {
	out, err := yaml.Marshal(&m)
	if err != nil {
		return "failed to marshal module to yaml"
	}
	return string(out)
}

func NewModuleFromSpec(s spec.Spec) Module {

	appName := makeModuleName(s.Module)
	result := Module{
		_spec:     s,
		Name:      appName,
		Path:      makePath(s.Module),
		Ops:       makeOps(appName),
		Packages:  makePackages(s.Module),
		Commands:  makeCommands(s.Commands),
		Resources: makeResources(s.Resources),
	}

	return result
}

func makeModuleName(specModule string) models.Name {
	moduleName := filepath.Base(specModule)
	return models.MakeName(moduleName)
}

func makePath(specModule string) utils.PathData {
	moduleName := filepath.Base(specModule)
	return utils.PathData{
		Base: moduleName,
		Full: specModule,
	}
}

func makeOps(appName models.Name) ops.Ops {
	return ops.Ops{
		Workdir:      appName.Kebob,
		AppHTTPAlias: appName.Kebob,
		MainDatabase: makeDatabase(appName),
	}
}

func makeDatabase(appName models.Name) ops.Database {
	return ops.Database{
		Host:     "db",
		Port:     "3306",
		Username: "root",
		Password: "secret",
		Name:     appName.Snake,
		Debug:    "1",
	}
}
