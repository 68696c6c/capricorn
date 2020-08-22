package module

import (
	"github.com/68696c6c/capricorn/generator/models"
	"github.com/68696c6c/capricorn/generator/models/spec"
)

type Resource struct {
	_spec      spec.Resource
	Key        resourceKey     `yaml:"key"`
	Name       models.Name     `yaml:"name"`
	Fields     []ResourceField `yaml:"fields,omitempty"`
	Controller ResourceService `yaml:"controller,omitempty"`
	Repo       ResourceService `yaml:"repo,omitempty"`
	Service    ResourceService `yaml:"service,omitempty"`
}

func makeResources(specResources []spec.Resource) []Resource {
	// Need to know every model and how it relates to the other models.
	// This will let us know how to write table indexes in migrations, scaffold out dependency injection, etc.
	var result []Resource
	for _, r := range specResources {
		resource := makeResource(r)
		result = append(result, resource)
	}
	return result
}

func makeResource(specResource spec.Resource) Resource {
	recName := models.MakeName(specResource.Name)
	recKey := makeResourceKey(recName.Kebob, "")

	crud := makeResourceCrudService(specResource)

	result := Resource{
		_spec:      specResource,
		Key:        recKey,
		Name:       recName,
		Fields:     makeResourceFields(specResource, recName, recKey),
		Controller: crud,
		Repo:       crud,
		Service:    makeResourceCustomService(specResource, recName),
	}

	return result
}
