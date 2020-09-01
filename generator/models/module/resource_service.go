package module

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/spec"
)

const ResourceActionList = "list"
const ResourceActionView = "view"
const ResourceActionCreate = "create"
const ResourceActionUpdate = "update"
const ResourceActionDelete = "delete"

type ResourceService struct {
	Name    data.Name `yaml:"name"`
	Actions []string  `yaml:"actions"`
}

func makeResourceCrudService(specResource spec.Resource, recName data.Inflection) ResourceService {
	result := ResourceService{
		Name: recName.Plural,
	}
	actions := specResource.Actions
	if len(actions) == 0 {
		actions = []string{
			ResourceActionList,
			ResourceActionView,
			ResourceActionCreate,
			ResourceActionUpdate,
			ResourceActionDelete,
		}
	}
	for _, a := range actions {
		result.Actions = append(result.Actions, actionNameFromString(a))
	}
	return result
}

func makeResourceCustomService(specResource spec.Resource, recName data.Inflection) ResourceService {
	result := ResourceService{
		Name: data.MakeName(recName.Plural.Snake + "_service"),
	}
	for _, a := range specResource.Custom {
		result.Actions = append(result.Actions, a)
	}
	return result
}

func actionNameFromString(input string) string {
	switch input {
	case "index":
	case "list":
		return ResourceActionList
	case "read":
	case "view":
		return ResourceActionView
	case "create":
		return ResourceActionCreate
	case "update":
		return ResourceActionUpdate
	case "delete":
		return ResourceActionDelete
	}
	return ""
}
