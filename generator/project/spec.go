package project

import (
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const (
	packageSRC    = "src"
	packageApp    = "app"
	packageCMD    = "cmd"
	packageHTTP   = "http"
	packageModels = "models"
	packageRepos  = "repos"
)

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

type AuthorConfig struct {
	Name  string
	Email string
}

type Paths struct {
	Root string
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
type field struct {
	Name      string
	FieldName string
	Type      string
	Required  bool
	Tag       string
}

type Model struct {
	Name       string
	StructName string
	Fields     []*field
	BelongsTo  []string `yaml:"belongs_to"`
	HasMany    []string `yaml:"has_many"`
}

// Repos
type Repo struct {
	Name                     string
	InterfaceName            string
	StructName               string
	Model                    string
	ModelStructName          string
	ModelsImportPath         string
	Methods                  []string
	MethodTemplates          []string
	InterfaceTemplateMethods []string
}

// HTTP
type middleware struct {
	name string
}

type handler struct {
	name         string
	dependencies string
	middlewares  []middleware
}

type controller struct {
	name       string
	structName string
	handlers   []handler
}

type HTTP struct {
	controllers []controller
	handlers    []handler
}

func NewSpecFromFilePath(filePath string) (Spec, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return Spec{}, errors.Wrap(err, "failed read yml spec file")
	}

	spec := Spec{}
	err = yaml.Unmarshal(file, &spec)
	if err != nil {
		return Spec{}, errors.Wrap(err, "failed parse project spec")
	}

	// Set the module name.
	spec.ModuleName = filepath.Base(spec.Module)

	// Get the absolute path to the project (e.g. within $GOPATH)
	rootPath, err := getProjectPath()
	if err != nil {
		return Spec{}, errors.Wrap(err, "failed to determine project path")
	}
	projectPath := joinPath(rootPath, spec.Module)

	// Set the project pacakge paths.
	spec.Paths = Paths{
		Root:     projectPath,
		Packages: newPackages(projectPath),
	}

	// Set the project package import paths.
	spec.Imports = Imports{
		Packages: newPackages(spec.Module),
	}

	return spec, nil
}

func newPackages(base string) Packages {
	srcPath := joinPath(base, packageSRC)
	return Packages{
		App:    joinPath(srcPath, packageApp),
		CMD:    joinPath(srcPath, packageCMD),
		HTTP:   joinPath(srcPath, packageHTTP),
		Models: joinPath(srcPath, packageModels),
		Repos:  joinPath(srcPath, packageRepos),
	}
}
