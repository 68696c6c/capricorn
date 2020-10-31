package src

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/models/templates/src/app"
	"github.com/68696c6c/capricorn/generator/models/templates/src/controllers"
	"github.com/68696c6c/capricorn/generator/models/templates/src/database"
	enumMeta "github.com/68696c6c/capricorn/generator/models/templates/src/enums/meta"
	enumString "github.com/68696c6c/capricorn/generator/models/templates/src/enums/string"
	"github.com/68696c6c/capricorn/generator/models/templates/src/models"
	"github.com/68696c6c/capricorn/generator/models/templates/src/repos"
	"github.com/68696c6c/capricorn/generator/models/templates/src/services"
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
	DB   DB          `yaml:"db,omitempty"`
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
	Container golang.File   `yaml:"container,omitempty"`
	Enums     []golang.File `yaml:"enums,omitempty"`
	Domains   []Domain      `yaml:"domains,omitempty"`
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
	PackageData     data.PackageData
	Controller      golang.File `yaml:"controller,omitempty"`
	ControllerTest  golang.File `yaml:"controller_test,omitempty"`
	Repo            golang.File `yaml:"repo,omitempty"`
	RepoTest        golang.File `yaml:"repo_test,omitempty"`
	Model           golang.File `yaml:"model,omitempty"`
	ModelTest       golang.File `yaml:"model_test,omitempty"`
	Service         golang.File `yaml:"service,omitempty"`
	ServiceTest     golang.File `yaml:"service_test,omitempty"`
	Validator       golang.File `yaml:"validator,omitempty"`
	ValidatorTest   golang.File `yaml:"validator_test,omitempty"`
	modelType       data.TypeData
	containerFields []utils.ContainerFieldMeta
}

func NewSRCDDD(m module.Module, rootPath, timestamp string) SRC {
	domains := makeDomains(m)
	return SRC{
		_module:  m,
		basePath: rootPath,
		Path:     data.MakePathData(rootPath, m.Packages.SRC.Reference),
		Package:  m.Packages.SRC,
		App: App{
			Enums:     makeEnums(m),
			Container: makeContainer(m, domains),
			Domains:   domains,
		},
		DB: DB{
			Migrations: makeInitialMigrations(m, domains, timestamp),
		},
		Main: NewMainGo(rootPath, m.Packages.SRC.GetImport(), m.Packages.CMD.GetImport()),
	}
}

func makeContainer(m module.Module, domains []Domain) golang.File {
	var fields []utils.ContainerFieldMeta
	for _, d := range domains {
		fields = append(fields, d.containerFields...)
	}
	name := data.MakeName("app")
	container := app.NewContainer(app.Meta{
		ReceiverName: "a",
		FileName:     name.Snake,
		PackageData:  m.Packages.App,
		Name:         name,
	}, fields)
	return container.MustGetFile()
}

func makeInitialMigrations(m module.Module, domains []Domain, version string) []golang.File {
	var imports []string
	var modelRefs []string
	for _, d := range domains {
		imports = append(imports, d.PackageData.GetImport())
		modelRefs = append(modelRefs, d.modelType.Reference)
	}

	mig := database.NewInitialMigration(database.MigrationMeta{
		PackageData: m.Packages.Migrations,
		AppImports:  imports,
		ModelRefs:   modelRefs,
	}, version)

	return []golang.File{mig.MustGetFile()}
}

func makeEnums(m module.Module) []golang.File {
	var result []golang.File
	for _, e := range m.Enums {
		typeName := e.Inflection.Single
		pkgData := m.Packages.Enums

		switch e.EnumType {
		case "string":
			fileData, pathData := data.MakeGoFileData(pkgData.GetImport(), e.Inflection.Single.Snake)
			enumFile := enumString.NewEnumFromMeta(enumMeta.Meta{
				FileData:         fileData,
				PathData:         pathData,
				PackageData:      pkgData,
				TypeName:         typeName.Exported,
				TypeNameReadable: typeName.Space,
				ReceiverName:     "e",
				Values:           e.Values,
			})
			result = append(result, enumFile.MustGetFile())
		}
	}

	return result
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
	vName := data.MakeName("validator")
	sName := data.MakeName("service")
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
	modelType := model.GetType()

	validator := validators.NewValidatorFromMeta(utils.ServiceMeta{
		ReceiverName: validationReceiver,
		FileName:     vName.Snake,
		Resource:     r,
		PackageData:  pkgData,
	}, model.GetValidationFields())

	unscopedRepoName := r.Inflection.Plural.Exported + "Repo"
	repo := repos.NewRepoFromMeta(utils.ServiceMeta{
		ReceiverName: "r",
		FileName:     rName.Snake,
		Resource:     r,
		PackageData:  pkgData,
		Name:         rName,
		// In a non-DDD app, we would use the Reference instead of the Name.
		ModelType: modelType.Name,
	})
	repoType := repo.GetInterfaceType()

	controller := controllers.NewControllerFromMeta(utils.ServiceMeta{
		ReceiverName: "c",
		FileName:     cName.Snake,
		Resource:     r,
		PackageData:  pkgData,
		Name:         cName,
		// In a non-DDD app, we would use the Reference instead of the Name.
		ModelType: modelType.Name,
	}, controllers.ControllerMeta{
		CreateRequestType:    createRequestName.Exported,
		UpdateRequestType:    updateRequestName.Exported,
		ResourceResponseType: viewResponseName.Exported,
		ListResponseType:     listResponseName.Exported,
		RepoType:             repoType,
		Exported:             true,
	})

	unscopedServiceName := r.Inflection.Plural.Exported + "Service"
	service := services.NewServiceFromMeta(utils.ServiceMeta{
		ReceiverName: "s",
		FileName:     sName.Snake,
		Resource:     r,
		PackageData:  pkgData,
		Name:         sName,
	}, services.ServiceMeta{
		ImplementationName: unscopedServiceName,
		// In a non-DDD app, we would use repoType.Reference instead of repoType.Name since the repo would be in a different package.
		RepoTypeRef: repoType.Name,
		RepoName:    repo.GetName(),
	})
	serviceType := service.GetInterfaceType()

	domainKey := r.Inflection.Plural.Kebob
	return Domain{
		PackageData: pkgData,
		Model:       model.MustGetFile(),
		Validator:   validator.MustGetFile(),
		Repo:        repo.MustGetFile(),
		Controller:  controller.MustGetFile(),
		Service:     service.MustGetFile(),
		modelType:   modelType,
		containerFields: []utils.ContainerFieldMeta{
			{
				DomainKey:     domainKey,
				PackageImport: pkgData.GetImport(),
				Name:          repo.GetName(),
				Constructor:   repo.GetConstructor(),
				TypeData:      repoType,
				ServiceType:   utils.ServiceTypeRepo,
				Field: golang.Field{
					Name: unscopedRepoName,
					Type: repoType.Reference,
				},
			},
			{
				DomainKey:     domainKey,
				PackageImport: pkgData.GetImport(),
				Name:          service.GetName(),
				Constructor:   service.GetConstructor(),
				TypeData:      serviceType,
				ServiceType:   utils.ServiceTypeService,
				Field: golang.Field{
					Name: unscopedServiceName,
					Type: serviceType.Reference,
				},
			},
		},
	}
}
