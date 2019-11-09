package utils

import (
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const (
	pathDocker = "docker"
	pathOPS    = "ops"

	packageSRC    = "src"
	packageApp    = "app"
	packageCMD    = "cmd"
	packageHTTP   = "http"
	packageModels = "models"
	packageRepos  = "repos"
)

// Note: all structs and fields are exported so that the Spec.String() function can use yaml.Marshal to print the spec.

type Spec struct {
	Name       string
	Module     string
	ModuleName string
	License    string
	Author     AuthorConfig

	Paths   Paths
	Imports Imports

	Models []*Model
	Repos  []*Repo
	HTTP   HTTP
}

func (s Spec) String() string {
	out, err := yaml.Marshal(&s)
	if err != nil {
		return "failed to marshal spec to yaml"
	}
	return string(out)
}

type AuthorConfig struct {
	Name  string
	Email string
}

type Paths struct {
	Root   string
	Docker string
	Packages
}

type Imports struct {
	Packages
}

type Packages struct {
	App    string
	CMD    string
	HTTP   string
	Models string
	Repos  string
}

// Models
type Field struct {
	Name     string
	Type     string
	Required bool
	// Internal fields.
	FieldName string
	Tag       string
}

type Model struct {
	Name      string
	Fields    []*Field
	BelongsTo []string `yaml:"belongs_to"`
	HasMany   []string `yaml:"has_many"`
	// Internal fields.
	StructName string
}

// Repos
type Repo struct {
	Model   string
	Methods []string
	// Internal fields.
	Name                     string
	StructName               string
	InterfaceName            string
	ModelStructName          string
	ModelsImportPath         string
	MethodTemplates          []string
	InterfaceTemplateMethods []string
}

// HTTP
type HTTP struct {
	Middlewares     []*middleware
	Controllers     []*controller
	Routes          []RouteGroup
	RoutesTemplates []string
}

type middleware struct {
	Name string
}

type controller struct {
	Resource string
	Handlers []string

	// Internal fields.
	FileName         string
	StructNameUpper  string
	StructNameLower  string
	AppImportPath    string
	ModelsImportPath string

	HandlerTemplates  []string
	RequestTemplates  []string
	ResponseTemplates []string

	ResourceNameUpper       string
	ResourceNameUpperPlural string
	ResourceNameLower       string
	ResourceNameLowerPlural string

	Routes []Route
}

type RouteGroup struct {
	ControllerConstructor string
	ControllerName        string
	GroupName             string
	Routes                []Route
}

type Route struct {
	Method string
	URI    string
}

type handler struct {
	name         string
	dependencies string
	middlewares  []middleware
}

func NewSpecFromFilePath(filePath string) (Spec, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return Spec{}, errors.Wrap(err, "failed to read yaml spec file")
	}

	spec := Spec{}
	err = yaml.Unmarshal(file, &spec)
	if err != nil {
		return Spec{}, errors.Wrap(err, "failed to parse project spec")
	}

	// Set the module name.
	spec.ModuleName = filepath.Base(spec.Module)

	// Get the absolute path to the project (e.g. within $GOPATH)
	rootPath, err := GetProjectPath()
	if err != nil {
		return Spec{}, errors.Wrap(err, "failed to determine project path")
	}
	projectPath := JoinPath(rootPath, spec.Module)

	// Set the project pacakge paths.
	spec.Paths = Paths{
		Root:     projectPath,
		Docker:   JoinPath(projectPath, pathDocker),
		Packages: newPackages(projectPath),
	}

	// Set the project package import paths.
	spec.Imports = Imports{
		Packages: newPackages(spec.Module),
	}

	return spec, nil
}

func newPackages(base string) Packages {
	srcPath := JoinPath(base, packageSRC)
	return Packages{
		App:    JoinPath(srcPath, packageApp),
		CMD:    JoinPath(srcPath, packageCMD),
		HTTP:   JoinPath(srcPath, packageHTTP),
		Models: JoinPath(srcPath, packageModels),
		Repos:  JoinPath(srcPath, packageRepos),
	}
}
