package models

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/68696c6c/capricorn/generator/utils"

	"github.com/jinzhu/inflection"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Name struct {
	Snake      string
	Kebob      string
	Exported   string
	Unexported string
	Package    string
}

type ProjectResource struct {
	Single Name
	Plural Name
}

type Model struct {
	Resource    ProjectResource `yaml:"resource,omitempty"`
	Name        Name            `yaml:"name,omitempty"`
	Imports     []string        `yaml:"imports,omitempty"`
	Filename    string          `yaml:"filename,omitempty"`
	Constructor string          `yaml:"constructor,omitempty"`

	Fields    []Field  `yaml:"fields,omitempty"`
	BelongsTo []string `yaml:"belongs_to,omitempty"`
	HasMany   []string `yaml:"has_many,omitempty"`
}

type Field struct {
	Name     Name
	Type     string
	Required bool
	Tag      string
}

type Repo struct {
	Resource    ProjectResource `yaml:"resource,omitempty"`
	Name        Name            `yaml:"name,omitempty"`
	Imports     []string        `yaml:"imports,omitempty"`
	Filename    string          `yaml:"filename,omitempty"`
	Constructor string          `yaml:"constructor,omitempty"`

	Methods []Method `yaml:"methods,omitempty"`
}

type Method struct {
	Resource  ProjectResource
	Name      string
	Imports   []string
	Signature string
	Receiver  string
}

type Controller struct {
	Resource    ProjectResource `yaml:"resource,omitempty"`
	Name        Name            `yaml:"name,omitempty"`
	Imports     []string        `yaml:"imports,omitempty"`
	Filename    string          `yaml:"filename,omitempty"`
	Constructor string          `yaml:"constructor,omitempty"`

	GroupName string    `yaml:"group_name,omitempty"`
	Handlers  []Handler `yaml:"handlers,omitempty"`
}

type Handler struct {
	Resource    ProjectResource
	Name        string
	Imports     []string
	Filename    string
	Constructor string

	Signature   string
	Receiver    string
	URI         string
	Action      string
	Middlewares []Middleware
}

type Middleware struct {
	Resource    ProjectResource
	Name        Name
	Imports     []string
	Filename    string
	Constructor string

	Parameters []string
}

type Paths struct {
	Root   string
	SRC    string
	OPS    string
	Docker string
	App    string
	HTTP   string
	Repos  string
	Models string
}

type Container struct {
	Repos []Repo `yaml:"repos,omitempty"`
}

type Project struct {
	Config Config
	Module Name

	Paths   Paths
	Imports Paths

	Container   Container    `yaml:"container,omitempty"`
	Controllers []Controller `yaml:"controllers,omitempty"`
	Repos       []Repo       `yaml:"repos,omitempty"`
	Models      []Model      `yaml:"models,omitempty"`
}

func (s Project) String() string {
	out, err := yaml.Marshal(&s)
	if err != nil {
		return "failed to marshal spec to yaml"
	}
	return string(out)
}

type Config struct {
	Name      string
	Module    string
	License   string
	Author    Author
	Resources []Resource
}

type Resource struct {
	Name      string
	BelongsTo []string `yaml:"belongs_to"`
	HasMany   []string `yaml:"has_many"`
	Fields    []ResourceField
	Actions   []string
}

type ResourceField struct {
	Name     string
	Type     string
	Required bool
}

type Author struct {
	Name         string
	Email        string
	Organization string
}

const (
	pathSRC    = "src"
	pathOPS    = "ops"
	pathDocker = "docker"
	pathApp    = "app"
	pathHTTP   = "http"
	pathRepos  = "repos"
	pathModels = "models"
)

func NewProject(filePath string) (Project, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return Project{}, errors.Wrap(err, "failed to read yaml spec file")
	}

	config := Config{}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return Project{}, errors.Wrap(err, "failed to parse project spec")
	}

	spec := Project{
		Config: config,
	}

	kebob := filepath.Base(config.Module)
	spec.Module = Name{
		Snake:      utils.SeparatedToSnake(kebob),
		Kebob:      kebob,
		Exported:   utils.SeparatedToExported(kebob),
		Unexported: utils.SeparatedToUnexported(kebob),
		Package:    config.Module,
	}

	projectPath, err := utils.GetProjectPath()
	if err != nil {
		return Project{}, errors.Wrap(err, "failed to determine project path")
	}
	rootPath := utils.JoinPath(projectPath, config.Module)

	spec.Paths = makePaths(rootPath)
	spec.Imports = makePaths(config.Module)

	for _, r := range spec.Config.Resources {
		resource := makeProjectResource(r.Name)

		model := makeModel(resource, r)
		spec.Models = append(spec.Models, model)

		repo := makeRepo(resource, r)
		spec.Repos = append(spec.Repos, repo)

		controller := makeController(resource, r)
		spec.Controllers = append(spec.Controllers, controller)
	}

	return spec, nil
}

func makePaths(rootPath string) Paths {
	srcPath := utils.JoinPath(rootPath, pathSRC)
	return Paths{
		Root:   rootPath,
		SRC:    srcPath,
		OPS:    utils.JoinPath(rootPath, pathOPS),
		Docker: utils.JoinPath(rootPath, pathDocker),
		App:    utils.JoinPath(srcPath, pathApp),
		HTTP:   utils.JoinPath(srcPath, pathHTTP),
		Repos:  utils.JoinPath(srcPath, pathRepos),
		Models: utils.JoinPath(srcPath, pathModels),
	}
}

func makeProjectResource(name string) ProjectResource {
	single := inflection.Singular(name)
	plural := inflection.Plural(name)
	return ProjectResource{
		Single: MakeName(single),
		Plural: MakeName(plural),
	}
}

func makeModel(r ProjectResource, config Resource) Model {
	result := Model{
		Resource:    r,
		Name:        MakeName(r.Single.Exported),
		Imports:     []string{},
		Filename:    r.Single.Kebob + ".go",
		Constructor: "New" + r.Single.Exported,
		HasMany:     config.HasMany,
		BelongsTo:   config.BelongsTo,
	}

	// Build fields.
	var fields []Field

	if len(config.BelongsTo) > 0 {
		for _, r := range config.BelongsTo {
			t := MakeName(r)
			rName := MakeName(fmt.Sprintf("%s_id", t.Exported))
			field := Field{
				Name: rName,
				Type: "goat.ID",
			}
			var extra string
			// if r.Required {
			if true {
				extra = ` binding:"required"`
			}
			field.Tag = fmt.Sprintf(`json:"%s"%s`, rName.Snake, extra)
			fields = append(fields, field)
		}
	}

	if len(config.HasMany) > 0 {
		for _, r := range config.HasMany {
			t := MakeName(r)
			field := Field{
				Name: MakeName(t.Unexported),
				Type: fmt.Sprintf("[]*%s", t.Exported),
			}
			fields = append(fields, field)
		}
	}

	for _, f := range config.Fields {
		field := Field{
			Name:     MakeName(f.Name),
			Type:     f.Type,
			Required: f.Required,
		}
		var extra string
		if f.Required {
			extra = ` binding:"required"`
		}
		field.Tag = fmt.Sprintf(`json:"%s"%s`, f.Name, extra)
		fields = append(fields, field)
	}

	result.Fields = fields
	return result
}

func makeRepo(r ProjectResource, config Resource) Repo {
	repoName := MakeName(r.Plural.Exported + "_repo")
	result := Repo{
		Resource:    r,
		Name:        repoName,
		Imports:     []string{},
		Filename:    r.Plural.Kebob + "_repo.go",
		Constructor: "New" + r.Plural.Exported + "Repo",
	}

	// Build fields.
	var methods []Method

	for _, a := range config.Actions {

		switch a {
		case "create":
			fallthrough
		case "update":
			arg := fmt.Sprintf("m *models.%s", r.Single.Exported)
			save := makeMethod(r, repoName, "Save", []string{arg}, []string{"[]error"})
			methods = append(methods, save)
			break
		case "delete":
			del := makeMethod(r, repoName, "Delete", []string{"id goat.ID"}, []string{"[]error"})
			methods = append(methods, del)
			break
		case "view":
			result := fmt.Sprintf("models.%s", r.Single.Exported)
			get := makeMethod(r, repoName, "GetByID", []string{"id goat.ID"}, []string{result, "[]error"})
			methods = append(methods, get)
			break
		case "list":
			result := fmt.Sprintf("[]*models.%s", r.Single.Exported)
			list := makeMethod(r, repoName, "List", []string{"q *query.Query"}, []string{result, "[]error"})
			methods = append(methods, list)

			setTotal := makeMethod(r, repoName, "SetQueryTotal", []string{"q *query.Query"}, []string{"[]error"})
			methods = append(methods, setTotal)
			break
		}
	}

	result.Methods = methods
	return result
}

func makeMethod(r ProjectResource, repoName Name, name string, parameters, returns []string) Method {
	sig := fmt.Sprintf("%s(%s) (%s)", name, strings.Join(parameters, ", "), strings.Join(returns, ", "))
	return Method{
		Resource:  r,
		Name:      name,
		Imports:   []string{},
		Signature: sig,
		Receiver:  repoName.Exported,
	}
}

func makeController(r ProjectResource, config Resource) Controller {
	controllerName := MakeName(r.Plural.Exported)
	result := Controller{
		Resource:    r,
		Name:        controllerName,
		Imports:     []string{},
		Filename:    r.Plural.Kebob + ".go",
		Constructor: "new" + r.Plural.Exported + "Controller",
	}

	// Build fields.
	var handlers []Handler

	for _, a := range config.Actions {

		switch a {
		case "create":
			create := makeHandler(r, controllerName, "POST", "", "Create")
			handlers = append(handlers, create)
			break
		case "update":
			update := makeHandler(r, controllerName, "PUT", "/:id", "Update")
			handlers = append(handlers, update)
			break
		case "delete":
			del := makeHandler(r, controllerName, "DELETE", "/:id", "Delete")
			handlers = append(handlers, del)
			break
		case "view":
			view := makeHandler(r, controllerName, "GET", "/:id", "View")
			handlers = append(handlers, view)
			break
		case "list":
			list := makeHandler(r, controllerName, "GET", "", "List")
			handlers = append(handlers, list)
			break
		}
	}

	result.Handlers = handlers
	return result
}

func makeHandler(r ProjectResource, controllerName Name, action, uri, name string) Handler {
	sig := fmt.Sprintf("%s(c *gin.Context)", name)
	return Handler{
		Resource:    r,
		Name:        name,
		Imports:     []string{},
		URI:         uri,
		Action:      action,
		Signature:   sig,
		Receiver:    controllerName.Exported,
		Middlewares: []Middleware{},
	}
}

func MakeName(base string) Name {
	return Name{
		Snake:      utils.SeparatedToSnake(base),
		Kebob:      utils.SeparatedToKebob(base),
		Exported:   utils.SeparatedToExported(base),
		Unexported: utils.SeparatedToUnexported(base),
		Package:    "",
	}
}
