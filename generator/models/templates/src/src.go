package src

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"gopkg.in/yaml.v2"
)

// SRC represents the structure of a project src directory.
type SRC struct {
	_module  module.Module
	basePath string

	Package data.PackageData `yaml:"package"`
	Path    data.PathData    `yaml:"path"`

	App  App  `yaml:"app"`
	CMD  CMD  `yaml:"cmd"`
	HTTP HTTP `yaml:"http"`
}

func (m SRC) String() string {
	out, err := yaml.Marshal(&m)
	if err != nil {
		return "failed to marshal src to yaml"
	}
	return string(out)
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
	Validator      golang.File `yaml:"validator"`
	ValidatorTest  golang.File `yaml:"validator_test"`
}

type serviceMeta struct {
	name         data.Name
	receiverName string
	fileName     string
	resource     module.Resource
	packageData  data.PackageData
}

func NewSRCDDD(m module.Module, rootPath string) SRC {
	return SRC{
		_module:  m,
		basePath: rootPath,
		Path:     data.MakePathData(rootPath, m.Packages.SRC.Reference),
		Package:  m.Packages.SRC,
		App: App{
			// Container: makeContainer(c),
			Domains: makeDomains(m),
		},
	}
}

func makeDomains(m module.Module) []Domain {
	baseDomainPath := m.Packages.Domains.Path.Full

	var result []Domain
	for _, r := range m.Resources {
		result = append(result, makeDomain(r, baseDomainPath))
	}

	return result
}

func mergeImports(target, additional golang.Imports) golang.Imports {
	target.Standard = append(target.Standard, additional.Standard...)
	target.Vendor = append(target.Vendor, additional.Vendor...)
	target.App = append(target.App, additional.App...)
	return target
}

func makeDomain(r module.Resource, baseDomainPath string) Domain {

	// If this function is ever called, we are definitely generating a DDD app so name things accordingly.
	cName := data.MakeName("controller")
	rName := data.MakeName("repo")
	mName := data.MakeName("model")
	// sName := data.MakeName("service")
	vName := data.MakeName("validator")
	createRequestName := data.MakeName("create_request")
	updateRequestName := data.MakeName("update_request")
	viewResponseName := data.MakeName("response")
	listResponseName := data.MakeName("list_response")
	pkgData := data.MakePackageData(baseDomainPath, r.Inflection.Plural.Snake)

	model := newModelFromMeta(serviceMeta{
		receiverName: "m",
		fileName:     mName.Snake,
		resource:     r,
		packageData:  pkgData,
		name:         r.Inflection.Single,
	})

	validator := newValidatorFromMeta(validatorMeta{
		receiverName: "r",
		fileName:     vName.Snake,
		resource:     r,
		packageData:  pkgData,
		fields:       model.GetValidationFields(),
	})

	repo := newRepoFromMeta(serviceMeta{
		receiverName: "r",
		fileName:     rName.Snake,
		resource:     r,
		packageData:  pkgData,
		name:         rName,
	})

	controller := newControllerFromMeta(serviceMeta{
		receiverName: "c",
		fileName:     cName.Snake,
		resource:     r,
		packageData:  pkgData,
		name:         cName,
	}, controllerMeta{
		createRequestName: createRequestName.Exported,
		updateRequestName: updateRequestName.Exported,
		viewResponseName:  viewResponseName.Exported,
		listResponseName:  listResponseName.Exported,
		repoType:          repo.GetInterface().Name,
	})

	return Domain{
		Model:      model.MustGetFile(),
		Validator:  validator.MustGetFile(),
		Controller: controller.MustGetFile(),
		Repo:       repo.MustGetFile(),
		// Service: makeService(
		// 	serviceMeta{
		// 		resource:     r,
		// 		packageData:  pkgData,
		// 		name:         rName,
		// 		fileName:     rName.Snake,
		// 		receiverName: "s",
		// 	},
		// ),
	}
}
