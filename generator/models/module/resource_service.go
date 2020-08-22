package module

import (
	"github.com/68696c6c/capricorn/generator/models"
	"github.com/68696c6c/capricorn/generator/models/spec"

	"github.com/jinzhu/inflection"
)

type ResourceService struct {
	Name    models.Name   `yaml:"name"`
	Actions []models.Name `yaml:"actions"`
}

func makeResourceCrudService(specResource spec.Resource) ResourceService {
	recNamePlural := models.MakeName(inflection.Plural(specResource.Name))
	result := ResourceService{
		Name: recNamePlural,
	}
	for _, a := range specResource.Actions {
		actionName := models.MakeName(actionNameFromString(a))
		result.Actions = append(result.Actions, actionName)
	}
	return result
}

func makeResourceCustomService(specResource spec.Resource, recName models.Name) ResourceService {
	result := ResourceService{
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
