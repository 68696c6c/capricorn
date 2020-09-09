package module

import (
	"path/filepath"

	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/spec"
	"github.com/68696c6c/capricorn/generator/models/templates/ops"

	"gopkg.in/yaml.v2"
)

type Module struct {
	_spec spec.Spec

	Name    data.Name        `yaml:"name"`
	Package data.PackageData `yaml:"package"`
	Ops     ops.Ops          `yaml:"ops"`

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

func NewModuleFromSpec(s spec.Spec, ddd bool) Module {

	pkgData := makeModulePackage(s.Module)
	result := Module{
		_spec:     s,
		Name:      pkgData.Name,
		Package:   pkgData,
		Ops:       makeOps(pkgData.Name),
		Packages:  makePackages(s.Module),
		Commands:  makeCommands(s.Commands),
		Resources: makeResources(s.Resources, ddd),
	}

	return result
}

func makeModulePackage(specModule string) data.PackageData {
	moduleBase := filepath.Dir(specModule)
	moduleName := filepath.Base(specModule)
	pkgData := data.MakePackageData(moduleBase, moduleName)

	// Packages are usually referenced by their snake name, but for the top-level module we want to use the exact name the user provided.
	pkgData.Reference = moduleName

	return pkgData
}

func makeOps(appName data.Name) ops.Ops {
	return ops.Ops{
		Workdir:      appName.Kebob,
		AppHTTPAlias: appName.Kebob,
		MainDatabase: makeDatabase(appName),
	}
}

func makeDatabase(appName data.Name) ops.Database {
	return ops.Database{
		Host:     "db",
		Port:     "3306",
		Username: "root",
		Password: "secret",
		Name:     appName.Snake,
		Debug:    "1",
	}
}
