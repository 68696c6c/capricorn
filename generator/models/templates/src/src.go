package src

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/models/templates/src/handlers"

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
		cName := data.MakeName("controller")
		viewResponseName := data.MakeName("response")
		listResponseName := data.MakeName("list_response")
		domainPKG := data.MakePackageData(baseDomainPath, r.Inflection.Plural.Snake)

		d := Domain{
			Controller: makeController(controllerMeta{
				name:             cName,
				receiverName:     "c",
				viewResponseName: viewResponseName.Exported,
				listResponseName: listResponseName.Exported,
				fileName:         cName.Snake,
				resource:         r,
				packageData:      domainPKG,
			}),
		}
		result = append(result, d)
	}
	return result
}

type controllerMeta struct {
	name             data.Name
	receiverName     string
	viewResponseName string
	listResponseName string
	fileName         string
	resource         module.Resource
	packageData      data.PackageData
}

func makeController(c controllerMeta) golang.File {
	fileData, pathData := data.MakeGoFileData(c.packageData.GetImport(), c.fileName)
	result := golang.File{
		Name:    fileData,
		Path:    pathData,
		Package: c.packageData,
	}

	plural := c.resource.Inflection.Plural
	single := c.resource.Inflection.Single

	for _, a := range c.resource.Controller.Actions {
		switch a {

		case module.ResourceActionList:
			t := handlers.List{
				Receiver: c.receiverName,
				Plural:   plural,
				Single:   single,
				Response: c.listResponseName,
			}
			h := makeHandler("List", t.MustParse(), handlers.GetListImports())
			result.Functions = append(result.Functions, h)
			result.Imports = mergeImports(result.Imports, h.Imports)

		case module.ResourceActionView:
			t := handlers.View{
				Receiver: c.receiverName,
				Plural:   plural,
				Single:   single,
				Response: c.listResponseName,
			}
			h := makeHandler("View", t.MustParse(), handlers.GetViewImports())
			result.Functions = append(result.Functions, h)
			result.Imports = mergeImports(result.Imports, h.Imports)

		case module.ResourceActionCreate:
			t := handlers.Create{
				Receiver: c.receiverName,
				Plural:   plural,
				Single:   single,
				Response: c.listResponseName,
			}
			h := makeHandler("Create", t.MustParse(), handlers.GetCreateImports())
			result.Functions = append(result.Functions, h)
			result.Imports = mergeImports(result.Imports, h.Imports)

		case module.ResourceActionUpdate:
			t := handlers.Update{
				Receiver: c.receiverName,
				Plural:   plural,
				Single:   single,
				Response: c.listResponseName,
			}
			h := makeHandler("Update", t.MustParse(), handlers.GetUpdateImports())
			result.Functions = append(result.Functions, h)
			result.Imports = mergeImports(result.Imports, h.Imports)

		case module.ResourceActionDelete:
			t := handlers.Delete{
				Receiver: c.receiverName,
				Plural:   plural,
				Single:   single,
				Response: c.listResponseName,
			}
			h := makeHandler("Delete", t.MustParse(), handlers.GetDeleteImports())
			result.Functions = append(result.Functions, h)
			result.Imports = mergeImports(result.Imports, h.Imports)
		}
	}

	return result
}

func mergeImports(target, additional golang.Imports) golang.Imports {
	target.Standard = append(target.Standard, additional.Standard...)
	target.Vendor = append(target.Vendor, additional.Vendor...)
	target.App = append(target.App, additional.App...)
	return target
}

func makeHandler(name, body string, i golang.Imports) golang.Function {
	return golang.Function{
		Name:    name,
		Imports: i,
		Arguments: []golang.Value{
			{
				Name: "cx",
				Type: "*gin.Context",
			},
		},
		ReturnValues: []golang.Value{},
		Receiver: golang.Value{
			Name: "c",
			Type: "",
		},
		Body: body,
	}
}
