package src

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/models/templates/src/controllers"
	"github.com/68696c6c/capricorn/generator/models/templates/src/models"
	"github.com/68696c6c/capricorn/generator/models/templates/src/repos"
	"github.com/68696c6c/capricorn/generator/models/templates/src/utils"
	"github.com/68696c6c/capricorn/generator/models/templates/src/validators"

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

func makeDomain(r module.Resource, baseDomainPath string) Domain {

	// If this function is ever called, we are definitely generating a DDD app so name things accordingly.
	cName := data.MakeName("controller")
	rName := data.MakeName("repo")
	mName := data.MakeName("model")
	// sName := data.MakeName("service")
	vName := data.MakeName("validator")
	validationReceiver := "r"
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
	}, validationReceiver)
	// In a non-DDD app, we would use the Reference instead of the Name.
	modelType := model.GetTypeName()

	validator := validators.NewValidatorFromMeta(utils.ServiceMeta{
		ReceiverName: validationReceiver,
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
	// In a non-DDD app, we would use the Reference instead of the Name.
	repoType := repo.GetInterfaceType()

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
		Repo:       repo.MustGetFile(),
		Controller: controller.MustGetFile(),
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
