package module

import (
	"fmt"
	"github.com/68696c6c/capricorn/generator/models"
	"github.com/68696c6c/capricorn/generator/utils"
	"github.com/jinzhu/inflection"
	"github.com/pkg/errors"
	"path/filepath"
	"strings"

	"github.com/68696c6c/capricorn/generator/models/templates"
	"github.com/68696c6c/capricorn/generator/models/templates/ops"
)

type Module struct {
	spec models.Spec

	Name models.Name        `yaml:"name"`
	Path templates.FileData `yaml:"path"`
	Ops  ops.Ops            `yaml:"ops"`

	Packages Packages `yaml:"packages"`

	Commands  []Command  `yaml:"commands,omitempty"`
	Resources []Resource `yaml:"resources,omitempty"`
	// Endpoints  []AppEndpoint  `yaml:"endpoints,omitempty"`
}

type Command struct {
	spec models.ConfigCommand
	Name models.Name
}

type resourceKey struct {
	Resource string
	Field    string
}

func (r resourceKey) String() string {
	if r.Field == "" {
		return r.Resource
	}
	return fmt.Sprintf("%s.%.s", r.Resource, r.Field)
}

func resourceKeyFromString(input string) resourceKey {
	parts := strings.Split(input, ".")
	if len(parts) > 1 {
		return resourceKey{
			Resource: parts[0],
			Field:    parts[1],
		}
	}
	if len(parts) == 1 {
		return resourceKey{
			Resource: parts[0],
		}
	}
	return resourceKey{}
}

func makeResourceKey(resource, field string) resourceKey {
	return resourceKey{
		Resource: resource,
		Field:    field,
	}
}

type AppResourceService struct {
	Name    models.Name   `yaml:"name"`
	Actions []models.Name `yaml:"actions"`
}

type Resource struct {
	spec       models.Resource
	Key        resourceKey        `yaml:"key"`
	Name       models.Name        `yaml:"name"`
	Fields     []AppResourceField `yaml:"fields,omitempty"`
	Controller AppResourceService `yaml:"controller,omitempty"`
	Repo       AppResourceService `yaml:"repo,omitempty"`
	Service    AppResourceService `yaml:"service,omitempty"`
}

type AppResourceField struct {
	spec        models.ResourceField
	Key         resourceKey       `yaml:"key"`
	Name        models.Name       `yaml:"name"`
	Type        string            `yaml:"type"`
	Index       *AppResourceIndex `yaml:"index"`
	IsRequired  bool              `yaml:"is_required"`
	IsPrimary   bool              `yaml:"is_primary"`
	IsGoatField bool              `yaml:"is_goat_field"`
}

type AppResourceIndex struct {
	Resource models.Name `yaml:"resource_name"`
	Field    models.Name `yaml:"field_name"`
	Unique   bool        `yaml:"unique"`
}

func ModuleFromSpec(spec models.Spec) (Module, error) {

	appName := makeName(spec.Module)
	resources := makeResources(spec.Resources)
	result := Module{
		spec:      spec,
		Name:      appName,
		Path:      makePath(spec.Module),
		Ops:       makeOps(appName),
		Packages:  makePackages(spec.Module, resources),
		Commands:  makeCommands(spec.Commands),
		Resources: resources,
	}

	return result, nil
}

func getProjectRootPath(specModule, projectPath string) (string, error) {
	rootPath := projectPath
	if projectPath == "" {
		projectPath, err := utils.GetProjectPath()
		if err != nil {
			return "", errors.Wrap(err, "failed to determine project root path")
		}
		rootPath = utils.JoinPath(projectPath, specModule)
	}
	return rootPath, nil
}

func makeName(specModule string) models.Name {
	moduleName := filepath.Base(specModule)
	return models.MakeName(moduleName)
}

func makePath(specModule string) templates.FileData {
	moduleName := filepath.Base(specModule)
	return templates.FileData{
		Base: moduleName,
		Full: specModule,
	}
}

func makeOps(appName models.Name) ops.Ops {
	return ops.Ops{
		Workdir:      appName.Kebob,
		AppHTTPAlias: appName.Kebob,
		MainDatabase: makeDatabase(appName),
	}
}

func makeDatabase(appName models.Name) ops.Database {
	return ops.Database{
		Host:     "db",
		Port:     "3306",
		Username: "root",
		Password: "secret",
		Name:     appName.Snake,
		Debug:    "1",
	}
}

func makeCommands(specCommands []models.ConfigCommand) []Command {
	var result []Command
	for _, c := range specCommands {
		cmd := Command{
			spec: c,
			Name: makeName(c.Name),
		}
		result = append(result, cmd)
	}
	return result
}

func makeResources(specResources []models.Resource) []Resource {
	// Need to know every model and how it relates to the other models.
	// This will let us know how to write table indexes in migrations, scaffold out dependency injection, etc.
	var result []Resource
	for _, r := range specResources {
		resource := makeResource(r)
		result = append(result, resource)
	}
	return result
}

func makeResource(specResource models.Resource) Resource {
	recName := models.MakeName(specResource.Name)
	recKey := makeResourceKey(recName.Kebob, "")

	crud := makeResourceCrud(specResource)

	result := Resource{
		spec:       specResource,
		Key:        recKey,
		Name:       recName,
		Fields:     makeResourceFields(specResource, recName, recKey),
		Controller: crud,
		Repo:       crud,
		Service:    makeResourceService(specResource, recName),
	}

	return result
}

func makeResourceFields(specResource models.Resource, recName models.Name, recKey resourceKey) []AppResourceField {
	result := []AppResourceField{
		{
			Key:         makeResourceKey(recKey.Resource, "id"),
			Name:        models.MakeName(recName.Snake + "_id"),
			Type:        "goat.ID",
			IsPrimary:   true,
			IsGoatField: true,
		},
		{
			Key:         makeResourceKey(recKey.Resource, "created_at"),
			Name:        models.MakeName("created_at"),
			Type:        "time.Time",
			IsGoatField: true,
		},
		{
			Key:         makeResourceKey(recKey.Resource, "updated_at"),
			Name:        models.MakeName("updated_at"),
			Type:        "*time.Time",
			IsGoatField: true,
		},
		{
			Key:         makeResourceKey(recKey.Resource, "deleted_at"),
			Name:        models.MakeName("deleted_at"),
			Type:        "*time.Time",
			IsGoatField: true,
		},
	}
	for _, f := range specResource.Fields {
		field := AppResourceField{
			spec:       f,
			Key:        makeResourceKey(recKey.Resource, f.Name),
			Name:       models.MakeName(f.Name),
			Type:       f.Type,
			IsRequired: f.Required,
		}
		result = append(result, field)
	}
	return result
}

func makeResourceCrud(specResource models.Resource) AppResourceService {
	recNamePlural := models.MakeName(inflection.Plural(specResource.Name))
	result := AppResourceService{
		Name: recNamePlural,
	}
	for _, a := range specResource.Actions {
		actionName := models.MakeName(actionNameFromString(a))
		result.Actions = append(result.Actions, actionName)
	}
	return result
}

func makeResourceService(specResource models.Resource, recName models.Name) AppResourceService {
	result := AppResourceService{
		Name: models.MakeName(recName.Snake + "_service"),
	}
	for _, a := range specResource.Custom {
		result.Actions = append(result.Actions, models.MakeName(a))
	}
	return result
}

func actionNameFromString(input string) string {
	switch input {
	case "index":
	case "list":
		return "list"
	case "read":
	case "view":
		return "view"
	case "create":
		return "create"
	case "update":
		return "update"
	case "delete":
		return "delete"
	}
	return ""
}
