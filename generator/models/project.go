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
	Package     string          `yaml:"package,omitempty"`
	Resource    ProjectResource `yaml:"resource,omitempty"`
	Name        string          `yaml:"name,omitempty"`
	Imports     []string        `yaml:"imports,omitempty"`
	Filename    string          `yaml:"filename,omitempty"`
	Constructor string          `yaml:"constructor,omitempty"`

	Fields    []Field  `yaml:"fields,omitempty"`
	BelongsTo []string `yaml:"belongs_to,omitempty"`
	HasMany   []string `yaml:"has_many,omitempty"`
}

type Field struct {
	Name     string
	Type     string
	Required bool
	Tag      string
}

type Repo struct {
	Package       string          `yaml:"package,omitempty"`
	Resource      ProjectResource `yaml:"resource,omitempty"`
	Name          Name            `yaml:"name,omitempty"`
	Imports       []string        `yaml:"imports,omitempty"`
	VendorImports []string        `yaml:"vendor_imports,omitempty"`
	Filename      string          `yaml:"filename,omitempty"`
	Constructor   string          `yaml:"constructor,omitempty"`
	Interface     string          `yaml:"interface,omitempty"`
	InterfaceName string

	Methods            []Method `yaml:"methods,omitempty"`
	MethodTemplates    []string `yaml:"-"`
	InterfaceTemplates []string `yaml:"-"`

	VarName string
}

type Method struct {
	Resource  ProjectResource
	Name      string
	Imports   []string
	Signature string
	Receiver  string
}

type Controller struct {
	Package       string          `yaml:"package,omitempty"`
	Resource      ProjectResource `yaml:"resource,omitempty"`
	Name          Name            `yaml:"name,omitempty"`
	Imports       []string        `yaml:"imports,omitempty"`
	VendorImports []string        `yaml:"vendor_imports,omitempty"`
	Filename      string          `yaml:"filename,omitempty"`
	Constructor   string          `yaml:"constructor,omitempty"`
	RepoName      string

	GroupName         string              `yaml:"group_name,omitempty"`
	Handlers          []Handler           `yaml:"handlers,omitempty"`
	HandlerTemplates  []string            `yaml:"-"`
	Requests          map[string]Request  `yaml:"requests,omitempty"`
	RequestTemplates  []string            `yaml:"-"`
	Responses         map[string]Response `yaml:"responses,omitempty"`
	ResponseTemplates []string            `yaml:"-"`
	RoutesTemplates   []string            `yaml:"-"`
}

type Request struct {
	Name  string
	Model string
}

type Response struct {
	Name  string
	Model string
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
	Root       string
	SRC        string
	OPS        string
	Docker     string
	App        string
	CMD        string
	Database   string
	HTTP       string
	Domains    string
	Repos      string
	Models     string
	Migrations string
}

type Service struct {
	Package     string          `yaml:"package,omitempty"`
	Resource    ProjectResource `yaml:"resource,omitempty"`
	Name        Name            `yaml:"name,omitempty"`
	Imports     []string        `yaml:"imports,omitempty"`
	Filename    string          `yaml:"filename,omitempty"`
	Constructor string          `yaml:"constructor,omitempty"`
	RepoName    string          `yaml:"repo_name,omitempty"`
	RepoArg     string          `yaml:"repo_arg,omitempty"`

	Methods            []Method `yaml:"methods,omitempty"`
	MethodTemplates    []string `yaml:"-"`
	InterfaceTemplates []string `yaml:"-"`
}

type Container struct {
	Repos []Repo `yaml:"repos,omitempty"`
}

type Project struct {
	Config Config
	Module Name

	Paths             Paths
	Imports           Paths
	Commands          []Command
	Domains           []Domain
	ReposWithServices []Repo `yaml:"repos_with_services,omitempty"` // domains that have a service which need repo injection
	DomainRepos       []Repo `yaml:"domain_repos,omitempty"`        // repos that do not need to be injected into a service

	Container   Container    `yaml:"container,omitempty"`
	Controllers []Controller `yaml:"controllers,omitempty"`
	Services    []Service    `yaml:"services,omitempty"`
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
	Commands  []ConfigCommand
}

type ConfigCommand struct {
	Name    string   `yaml:"name"`
	Args    []string `yaml:"args"`
	Use     string
	VarName string
}

type Command struct {
	AppImport string
	Name      Name
	Args      []string
	Use       string
	FileName  string
	VarName   string
}

type Domain struct {
	Import     string
	Name       string
	Model      Model
	Repo       Repo
	Controller Controller
	Service    Service
}

type Resource struct {
	Name      string
	BelongsTo []string `yaml:"belongs_to"`
	HasMany   []string `yaml:"has_many"`
	Fields    []ResourceField
	Actions   []string
	Custom    []string
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
	pathSRC        = "src"
	pathOPS        = "ops"
	pathDocker     = "docker"
	pathApp        = "app"
	pathCMD        = "cmd"
	pathDatabase   = "db"
	pathHTTP       = "http"
	pathDomains    = "resources"
	pathRepos      = "repos"
	pathModels     = "models"
	pathMigrations = "migrations"
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

	for _, c := range spec.Config.Commands {
		command := makeCommand(c, spec.Imports.App)
		spec.Commands = append(spec.Commands, command)
	}

	for _, r := range spec.Config.Resources {
		resource := makeProjectResource(r.Name)

		// If no methods were specified, default to all.
		if len(r.Actions) == 0 {
			r.Actions = []string{
				"list",
				"view",
				"create",
				"update",
				"delete",
			}
		}

		domain := Domain{
			Import: fmt.Sprintf("%s/%s", spec.Imports.Domains, resource.Plural.Snake),
			Name:   resource.Plural.Snake,
		}

		model := makeModel(resource, r, domain.Name, spec.Imports.Domains)
		domain.Model = model
		spec.Models = append(spec.Models, model)

		repo := makeRepo(resource, r, domain.Name)
		domain.Repo = repo
		spec.Repos = append(spec.Repos, repo)

		if serv := makeService(resource, r, domain.Name, repo); serv != nil {
			domain.Service = *serv
			spec.Services = append(spec.Services, *serv)
			spec.ReposWithServices = append(spec.ReposWithServices, repo)
		} else {
			spec.DomainRepos = append(spec.DomainRepos, repo)
		}

		controller := makeController(resource, r, domain.Name, repo)
		domain.Controller = controller

		spec.Domains = append(spec.Domains, domain)
	}

	return spec, nil
}

func makePaths(rootPath string) Paths {
	srcPath := utils.JoinPath(rootPath, pathSRC)
	dbPath := utils.JoinPath(srcPath, pathDatabase)
	return Paths{
		Root:       rootPath,
		SRC:        srcPath,
		OPS:        utils.JoinPath(rootPath, pathOPS),
		Docker:     utils.JoinPath(rootPath, pathDocker),
		App:        utils.JoinPath(srcPath, pathApp),
		CMD:        utils.JoinPath(srcPath, pathCMD),
		Database:   dbPath,
		HTTP:       utils.JoinPath(srcPath, pathHTTP),
		Domains:    utils.JoinPath(srcPath, pathApp),
		Repos:      utils.JoinPath(srcPath, pathRepos),
		Models:     utils.JoinPath(srcPath, pathModels),
		Migrations: utils.JoinPath(dbPath, pathMigrations),
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

func makeCommand(c ConfigCommand, appImport string) Command {
	cName := strings.Replace(c.Name, ":", "_", -1)
	commandName := MakeName(cName)
	return Command{
		AppImport: appImport,
		Name:      commandName,
		Use:       c.Name,
		Args:      c.Args,
		VarName:   commandName.Unexported,
		FileName:  commandName.Snake + ".go",
	}
}

func makeModel(r ProjectResource, config Resource, packageName, domainsImportBase string) Model {
	result := Model{
		Resource:    r,
		Name:        r.Single.Exported,
		Package:     packageName,
		Imports:     []string{},
		Filename:    "model.go",
		Constructor: "New",
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
				Name: rName.Exported,
				Type: "goat.ID",
			}
			field.Tag = fmt.Sprintf(`json:"%s"`, rName.Snake)
			fields = append(fields, field)
		}
	}

	for _, f := range config.Fields {
		t := MakeName(f.Name)
		rName := MakeName(inflection.Singular(t.Exported))
		field := Field{
			Name:     rName.Exported,
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

	if len(config.HasMany) > 0 {
		for _, r := range config.HasMany {
			t := MakeName(r)
			single := inflection.Singular(t.Exported)
			sName := MakeName(single)
			plural := inflection.Plural(t.Unexported)
			pName := MakeName(plural)
			field := Field{
				Name: pName.Exported,
				Type: fmt.Sprintf("[]*%s.%s", pName.Unexported, sName.Exported),
			}
			field.Tag = fmt.Sprintf(`json:"%s"`, pName.Snake)
			fields = append(fields, field)

			result.Imports = append(result.Imports, fmt.Sprintf("%s/%s", domainsImportBase, pName.Unexported))
		}
	}

	result.Fields = fields
	return result
}

func makeRepo(r ProjectResource, config Resource, packageName string) Repo {
	repoName := MakeName("repo_GORM")
	result := Repo{
		Resource:      r,
		Name:          repoName,
		Package:       packageName,
		Imports:       []string{},
		VendorImports: []string{"github.com/jinzhu/gorm"},
		Filename:      "repo.go",
		Constructor:   "NewRepo",
		Interface:     r.Plural.Exported + "Repo",
		InterfaceName: "Repo",
		VarName:       r.Plural.Unexported + "Repo",
	}

	// Build fields.
	var methods []Method
	saveDone := false
	viewDone := false

	makeView := func() {
		returnName := fmt.Sprintf("%s", r.Single.Exported)
		get := makeMethod(r, repoName, "GetByID", []string{"id goat.ID"}, []string{returnName, "[]error"})
		methods = append(methods, get)

		result.VendorImports = append(result.VendorImports, "github.com/68696c6c/goat")
		viewDone = true
	}

	for _, a := range config.Actions {

		switch a {
		case "create":
			fallthrough
		case "update":
			if saveDone {
				break
			}
			arg := fmt.Sprintf("m *%s", r.Single.Exported)
			save := makeMethod(r, repoName, "Save", []string{arg}, []string{"errs []error"})
			methods = append(methods, save)
			saveDone = true

			// Update takes a model as an argument, which implies the need to retrieve a model.
			if !viewDone {
				makeView()
			}

			break
		case "delete":
			arg := fmt.Sprintf("m *%s", r.Single.Exported)
			del := makeMethod(r, repoName, "Delete", []string{arg}, []string{"[]error"})
			methods = append(methods, del)
			break
		case "view":
			if !viewDone {
				makeView()
			}
			break
		case "list":
			returnName := fmt.Sprintf("m []*%s", r.Single.Exported)
			list := makeMethod(r, repoName, "List", []string{"q *query.Query"}, []string{returnName, "errs []error"})
			methods = append(methods, list)

			setTotal := makeMethod(r, repoName, "SetQueryTotal", []string{"q *query.Query"}, []string{"[]error"})
			methods = append(methods, setTotal)

			result.VendorImports = append(result.VendorImports, "github.com/68696c6c/goat")
			result.VendorImports = append(result.VendorImports, "github.com/68696c6c/goat/query")
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

func makeService(r ProjectResource, config Resource, packageName string, repo Repo) *Service {
	if len(config.Custom) == 0 {
		return nil
	}
	serviceName := MakeName(fmt.Sprintf("%sService", r.Single.Exported))
	service := Service{
		Resource:    r,
		Name:        serviceName,
		Package:     packageName,
		Imports:     []string{},
		Filename:    "service.go",
		Constructor: "NewService",
		RepoArg:     r.Plural.Unexported + "Repo",
		RepoName:    repo.InterfaceName,
	}
	for _, action := range config.Custom {
		methodName := MakeName(action)
		arg := fmt.Sprintf("m *%s", r.Single.Exported)
		save := makeMethod(r, serviceName, methodName.Exported, []string{arg}, []string{"err error"})
		service.Methods = append(service.Methods, save)
	}
	return &service
}

func makeController(r ProjectResource, config Resource, packageName string, repo Repo) Controller {
	controllerName := MakeName(r.Plural.Unexported)
	result := Controller{
		Resource:      r,
		Name:          controllerName,
		Package:       packageName,
		Imports:       []string{},
		VendorImports: []string{"github.com/68696c6c/goat", "github.com/gin-gonic/gin"},
		Filename:      "controller.go",
		Constructor:   "NewController",
		GroupName:     r.Plural.Unexported + "Routes",
		RepoName:      repo.Interface,
	}

	// Build handlers.
	var handlers []Handler
	requests := map[string]Request{}
	responses := map[string]Response{}

	responses["resource"] = Response{
		Name:  "ResourceResponse",
		Model: r.Single.Exported,
	}

	for _, a := range config.Actions {

		switch a {
		case "create":
			create := makeHandler(r, controllerName, "POST", "", "Create")
			handlers = append(handlers, create)
			requests["create"] = Request{
				Name:  "CreateRequest",
				Model: r.Single.Exported,
			}
			break
		case "update":
			update := makeHandler(r, controllerName, "PUT", "/:id", "Update")
			handlers = append(handlers, update)
			requests["update"] = Request{
				Name:  "UpdateRequest",
				Model: r.Single.Exported,
			}
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
			responses["list"] = Response{
				Name:  "ListResponse",
				Model: r.Single.Exported,
			}
			result.VendorImports = append(result.VendorImports, "github.com/68696c6c/goat/query")
			break
		}
	}

	result.Handlers = handlers
	result.Requests = requests
	result.Responses = responses
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
