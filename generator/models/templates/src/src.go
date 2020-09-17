package src

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/models/templates/src/controllers"
	"github.com/68696c6c/capricorn/generator/models/templates/src/models"
	"github.com/68696c6c/capricorn/generator/models/templates/src/repos"
	"github.com/68696c6c/capricorn/generator/models/templates/src/utils"
	"gopkg.in/yaml.v2"
)

// SRC represents the structure of a project src directory.
type SRC struct {
	_module  module.Module
	basePath string

	Package data.PackageData `yaml:"package,omitempty"`
	Path    data.PathData    `yaml:"path,omitempty"`

	App  App         `yaml:"app,omitempty"`
	CMD  CMD         `yaml:"cmd,omitempty"`
	HTTP HTTP        `yaml:"http,omitempty"`
	Main golang.File `yaml:"main,omitempty"`
}

func (m SRC) String() string {
	out, err := yaml.Marshal(&m)
	if err != nil {
		return "failed to marshal src to yaml"
	}
	return string(out)
}

type App struct {
	Container golang.File `yaml:"container,omitempty"`
	Domains   []Domain    `yaml:"domains,omitempty"`
}

type CMD struct {
	Root    golang.File   `yaml:"root,omitempty"`
	Server  golang.File   `yaml:"server,omitempty"`
	Migrate golang.File   `yaml:"migrate,omitempty"`
	Seed    golang.File   `yaml:"seed,omitempty"`
	Custom  []golang.File `yaml:"custom,omitempty"`
}

type DB struct {
	Migrations []golang.File `yaml:"migrations,omitempty"`
	Seeders    []golang.File `yaml:"seeders,omitempty"`
}

type HTTP struct {
	Routes golang.File `yaml:"routes,omitempty"`
}

type Domain struct {
	Controller     golang.File `yaml:"controller,omitempty"`
	ControllerTest golang.File `yaml:"controller_test,omitempty"`
	Repo           golang.File `yaml:"repo,omitempty"`
	RepoTest       golang.File `yaml:"repo_test,omitempty"`
	Model          golang.File `yaml:"model,omitempty"`
	ModelTest      golang.File `yaml:"model_test,omitempty"`
	Service        golang.File `yaml:"service,omitempty"`
	ServiceTest    golang.File `yaml:"service_test,omitempty"`
	Validator      golang.File `yaml:"validator,omitempty"`
	ValidatorTest  golang.File `yaml:"validator_test,omitempty"`
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
		Main: NewMainGo(rootPath, m.Packages.SRC.GetImport(), m.Packages.CMD.GetImport()),
	}
}

func makeDomains(m module.Module) []Domain {
	baseDomainPath := m.Packages.Domains.GetImport()

	var result []Domain
	for _, r := range m.Resources {
		result = append(result, makeDomain(r, baseDomainPath))
	}

	return result
}

func removeDuplicateStrings(items []string) []string {
	keys := make(map[string]bool)
	var result []string
	for _, i := range items {
		if _, ok := keys[i]; !ok {
			keys[i] = true
			result = append(result, i)
		}
	}
	return result
}

func mergeImports(target, additional golang.Imports) golang.Imports {
	target.Standard = append(target.Standard, additional.Standard...)
	target.App = append(target.App, additional.App...)
	target.Vendor = append(target.Vendor, additional.Vendor...)
	return golang.Imports{
		Standard: removeDuplicateStrings(target.Standard),
		App:      removeDuplicateStrings(target.App),
		Vendor:   removeDuplicateStrings(target.Vendor),
	}
}

// type Meta struct {
// 	ControllerName       data.Name
//
// 	RepoName             data.Name
// 	RepoType             string
//
// 	ModelName            data.Name
// 	ModelType            string
//
// 	ValidatorName        data.Name
//
// 	CreateRequestName    data.Name
// 	UpdateRequestName    data.Name
// 	ResourceResponseName data.Name
// 	ListResponseName     data.Name
// 	ExportController     bool
// }

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

	model := models.NewModelFromMeta(utils.ServiceMeta{
		ReceiverName: "m",
		FileName:     mName.Snake,
		Resource:     r,
		PackageData:  pkgData,
		Name:         r.Inflection.Single,
	})
	// In a non-DDD app, we would use .Reference
	modelType := model.GetType().Type

	validator := models.NewValidatorFromMeta(utils.ServiceMeta{
		ReceiverName: "r",
		FileName:     vName.Snake,
		Resource:     r,
		PackageData:  pkgData,
	}, model.GetValidationFields())

	repo := repos.NewRepoFromMeta(utils.ServiceMeta{
		ReceiverName: "r",
		FileName:     rName.Snake,
		Resource:     r,
		PackageData:  pkgData,
		Name:         rName,
		ModelType:    modelType,
	})
	// In a non-DDD app, we would use .Reference
	repoType := repo.GetType().Type

	controller := controllers.NewControllerFromMeta(utils.ServiceMeta{
		ReceiverName: "c",
		FileName:     cName.Snake,
		Resource:     r,
		PackageData:  pkgData,
		Name:         cName,
		ModelType:    modelType,
	}, controllers.ControllerMeta{
		CreateRequestType:    createRequestName.Exported,
		UpdateRequestType:    updateRequestName.Exported,
		ResourceResponseType: viewResponseName.Exported,
		ListResponseType:     listResponseName.Exported,
		RepoType:             repoType,
		Exported:             true,
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
