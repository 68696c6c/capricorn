package module

import (
	"path/filepath"

	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/spec"
	"github.com/68696c6c/capricorn/generator/models/templates/ops"

	"gopkg.in/yaml.v2"
)

// A Module describes all of a project's resources that need to generated and their relationships to each other.
// Since the relation of Go packages is affected by the file structure, the Module model is also the decision point
// for DDD or MVC pattern.
type Module struct {
	_spec spec.Spec

	Name    data.Name        `yaml:"name,omitempty"`
	Package data.PackageData `yaml:"package,omitempty"`
	Ops     ops.Ops          `yaml:"ops,omitempty"`

	Packages Packages `yaml:"packages,omitempty"`

	Commands  []Command       `yaml:"commands,omitempty"`
	Enums     map[string]Enum `yaml:"enums,omitempty"`
	Resources []Resource      `yaml:"resources,omitempty"`
}

func (m Module) String() string {
	out, err := yaml.Marshal(&m)
	if err != nil {
		return "failed to marshal module to yaml"
	}
	return string(out)
}

func (m Module) GetAuthor() spec.Author {
	return m._spec.Author
}

func (m Module) GetLicense() string {
	return m._spec.License
}

func NewModuleFromSpec(s spec.Spec, ddd bool) Module {

	pkgData := makeModulePackage(s.Module)
	modulePackages := makePackages(s.Module, ddd)
	enumMap := makeEnums(s.Enums, modulePackages.Enums)
	result := Module{
		_spec:     s,
		Name:      pkgData.Name,
		Package:   pkgData,
		Ops:       makeOps(s.Ops, pkgData.Name),
		Packages:  modulePackages,
		Commands:  makeCommands(s.Commands),
		Enums:     enumMap,
		Resources: makeResources(s.Resources, enumMap, ddd),
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

func makeOps(specOps ops.Ops, appName data.Name) ops.Ops {
	workdir := appName.Kebob
	if specOps.Workdir != "" {
		workdir = specOps.Workdir
	}
	alias := appName.Kebob + ".local"
	if specOps.AppHTTPAlias != "" {
		alias = specOps.AppHTTPAlias
	}
	return ops.Ops{
		Workdir:      workdir,
		AppHTTPAlias: alias,
		MainDatabase: makeDatabase(specOps.MainDatabase, appName),
	}
}

func makeDatabase(specDatabase ops.Database, appName data.Name) ops.Database {
	host := "db"
	if specDatabase.Host != "" {
		host = specDatabase.Host
	}
	port := "3306"
	if specDatabase.Port != "" {
		port = specDatabase.Port
	}
	username := "root"
	if specDatabase.Username != "" {
		username = specDatabase.Username
	}
	password := "secret"
	if specDatabase.Password != "" {
		password = specDatabase.Password
	}
	name := appName.Snake
	if specDatabase.Name != "" {
		name = specDatabase.Name
	}
	debug := "1"
	if specDatabase.Debug != "" {
		debug = specDatabase.Debug
	}
	return ops.Database{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Name:     name,
		Debug:    debug,
	}
}
