package module

import (
	"fmt"

	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/spec"
)

type Resource struct {
	_spec      spec.Resource
	Key        resourceKey     `yaml:"key"`
	Inflection data.Inflection `yaml:"inflection"`
	Fields     []ResourceField `yaml:"fields,omitempty"`
	Controller ResourceService `yaml:"controller,omitempty"`
	Repo       ResourceService `yaml:"repo,omitempty"`
	Service    ResourceService `yaml:"service,omitempty"`
	FieldsMeta ResourceFields  `yaml:"fields_meta,omitempty"`
	Indexes    Indexes         `yaml:"indexes,omitempty"`
}

func (m Resource) GetPrimaryField() (ResourceField, bool) {
	for _, f := range m.Fields {
		if f.Name.Snake == "id" {
			return f, true
		}
	}
	return ResourceField{}, false
}

type Indexes struct {
	Primary Index   `yaml:"primary,omitempty"`
	Unique  []Index `yaml:"unique,omitempty"`
	Fields  []Index `yaml:"fields,omitempty"`
	Foreign []Index `yaml:"foreign,omitempty"`
}

type Index struct {
	Key               string    `yaml:"key"`
	FieldName         data.Name `yaml:"field_name"`
	ResourceName      data.Name `yaml:"resource_name"`
	ResourceFieldName data.Name `yaml:"resource_field_name"`
}

func makeResources(specResources []spec.Resource, ddd bool) []Resource {
	// Need to know every model and how it relates to the other models.
	// This will let us know how to write table indexes in migrations, scaffold out dependency injection, etc.
	var result []Resource
	resourceMap := map[string]Resource{}
	for _, r := range specResources {
		resource := makeResource(r, ddd)
		result = append(result, resource)
		resourceMap[resource.Inflection.Single.Snake] = resource
	}

	// Now that we have built each resource, we can build the relations and indexes.
	// We cannot build fields for the relations yet because the naming of resources depends on whether this is a DDD app
	// or not and that decision happens at the file generation level.
	for _, rec := range result {
		var pkName data.Name
		if pk, ok := rec.GetPrimaryField(); ok {
			pkName = pk.Name
		} else {
			pkName = data.MakeName("id")
		}

		indexes := Indexes{
			Primary: Index{
				FieldName: pkName,
			},
		}

		for _, f := range rec.FieldsMeta.BelongsTo {
			indexes.Foreign = append(indexes.Foreign, Index{
				FieldName:         f.Name,
				ResourceName:      f.Relation.Plural,
				ResourceFieldName: data.MakeName("id"),
			})
		}

		for _, f := range rec.FieldsMeta.Unique {
			indexes.Unique = append(indexes.Unique, Index{
				Key:       fmt.Sprintf("%s_%s_unique", rec.Inflection.Plural.Snake, f.Name.Snake),
				FieldName: f.Name,
			})
		}

		for _, f := range rec.FieldsMeta.Indexed {
			indexes.Fields = append(indexes.Fields, Index{
				Key:       fmt.Sprintf("%s_%s_index", rec.Inflection.Plural.Snake, f.Name.Snake),
				FieldName: f.Name,
			})
		}

	}

	return result
}

func makeResource(specResource spec.Resource, ddd bool) Resource {
	recName := data.MakeInflection(specResource.Name)
	recKey := makeResourceKey(recName.Single.Kebob, "")

	crud := makeResourceCrudService(specResource, recName)
	fields := makeResourceFields(specResource, recKey, ddd)

	result := Resource{
		_spec:      specResource,
		Key:        recKey,
		Inflection: recName,
		Fields:     fields.Model,
		Controller: crud,
		Repo:       crud,
		Service:    makeResourceCustomService(specResource, recName),
		FieldsMeta: fields,
	}

	return result
}
