package src

import (
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

type SRC struct {
	_module  module.Module
	basePath string
	path     templates.PathData
	pkgData  golang.PackageData
	Main     golang.File `yaml:"main"`
	App      App         `yaml:"app"`
	CMD      CMD         `yaml:"cmd"`
	HTTP     HTTP        `yaml:"http"`
}

type App struct {
	Container golang.File `yaml:"container"`
	Domains   []Domain    `yaml:"domains"`
}

type CMD struct {
	Root    golang.File   `yaml:"root"`
	Server  golang.File   `yaml:"server"`
	Migrate golang.File   `yaml:"migrate"`
	Seed    golang.File   `yaml:"seed"`
	Custom  []golang.File `yaml:"custom"`
}

type DB struct {
	Migrations []golang.File `yaml:"migrations"`
	Seeders    []golang.File `yaml:"seeders"`
}

type HTTP struct {
	Routes golang.File `yaml:"routes"`
}

type Domain struct {
	Controller     golang.File `yaml:"controller"`
	ControllerTest golang.File `yaml:"controller_test"`
	Repo           golang.File `yaml:"repo"`
	RepoTest       golang.File `yaml:"repo_test"`
	Model          golang.File `yaml:"model"`
	ModelTest      golang.File `yaml:"model_test"`
	Service        golang.File `yaml:"service"`
	ServiceTest    golang.File `yaml:"service_test"`
}

func NewSRCDDD(m module.Module, rootPath string) SRC {
	srcPath := m.Packages.SRC.Path // Base = src, Full = github.com/user/example/src

	result := SRC{
		_module:  m,
		basePath: rootPath,
		path: templates.PathData{
			Full: utils.JoinPath(rootPath, srcPath.Base),
			Base: rootPath,
		},
		pkgData: golang.PackageData{
			Name:   srcPath.Base,
			Module: utils.JoinPath(rootPath, srcPath.Base),
		},
	}
	result.makeDomains()
	// return SRC{
	// 	Main: NewMainGo(modulePackages, rootPath, rootPackage),
	// 	App: App{
	// 		// Container: makeContainer(c),
	// 		Domains:   makeDomains(),
	// 	},
	// }
	return result
}

func (m SRC) makeDomains() []Domain {
	// baseDomainPath := m._module.Packages.Domains.Path.Full

	var result []Domain
	// for _, r := range m._module.Resources {
	// 	plural := inflection.Plural(r.Name.Snake)
	// 	domainPKG := golang.PackageData{
	// 		Name:   plural,
	// 		Module: utils.JoinPath(baseDomainPath, plural),
	// 	}
	//
	// 	controllerName := models.MakeName("controller")
	//
	// 	d := Domain{
	// 		Controller: makeController(controllerName.Snake, domainPKG, r),
	// 	}
	// }
	return result
}

func makeController(baseFileName string, pkgData golang.PackageData, resource module.Resource) golang.File {
	fileData, pathData := templates.MakeGoFileData(pkgData.Module, baseFileName)
	return golang.File{
		Name:    fileData,
		Path:    pathData,
		Package: pkgData,
		Imports: golang.Imports{
			Standard: nil,
			App:      nil,
			Vendor:   nil,
		},
		InitFunction: golang.Function{
			Name:         "",
			Imports:      nil,
			Arguments:    nil,
			ReturnValues: nil,
			Receiver: golang.Value{
				Name: "",
				Type: "",
			},
			Body: "",
		},
		Consts:     nil,
		Vars:       nil,
		Interfaces: nil,
		Structs:    nil,
		Functions:  nil,
	}
}
