package module

import (
	"fmt"
	"github.com/68696c6c/capricorn/generator/models/data"

	"github.com/68696c6c/capricorn/generator/models/spec"
)

type ResourceField struct {
	_spec       spec.ResourceField
	Key         resourceKey    `yaml:"key"`
	Name        data.Name      `yaml:"name"`
	Type        string         `yaml:"type"`
	Index       *ResourceIndex `yaml:"index"`
	IsRequired  bool           `yaml:"is_required"`
	IsPrimary   bool           `yaml:"is_primary"`
	IsGoatField bool           `yaml:"is_goat_field"`
}

type ResourceIndex struct {
	Resource data.Name `yaml:"resource_name"`
	Field    data.Name `yaml:"field_name"`
	Unique   bool      `yaml:"unique"`
}

type resourceKey struct {
	Resource string `yaml:"resource"`
	Field    string `yaml:"field"`
}

func (r resourceKey) String() string {
	if r.Field == "" {
		return r.Resource
	}
	return fmt.Sprintf("%s.%.s", r.Resource, r.Field)
}

func makeResourceFields(specResource spec.Resource, recKey resourceKey) []ResourceField {
	result := []ResourceField{
		{
			Key:         makeResourceKey(recKey.Resource, "id"),
			Name:        data.MakeName("id"),
			Type:        "goat.ID",
			IsPrimary:   true,
			IsGoatField: true,
		},
		{
			Key:         makeResourceKey(recKey.Resource, "created_at"),
			Name:        data.MakeName("created_at"),
			Type:        "time.Time",
			IsGoatField: true,
		},
		{
			Key:         makeResourceKey(recKey.Resource, "updated_at"),
			Name:        data.MakeName("updated_at"),
			Type:        "*time.Time",
			IsGoatField: true,
		},
		{
			Key:         makeResourceKey(recKey.Resource, "deleted_at"),
			Name:        data.MakeName("deleted_at"),
			Type:        "*time.Time",
			IsGoatField: true,
		},
	}
	for _, f := range specResource.Fields {
		field := ResourceField{
			_spec:      f,
			Key:        makeResourceKey(recKey.Resource, f.Name),
			Name:       data.MakeName(f.Name),
			Type:       f.Type,
			IsRequired: f.Required,
		}
		result = append(result, field)
	}
	return result
}

func makeResourceKey(resource, field string) resourceKey {
	return resourceKey{
		Resource: resource,
		Field:    field,
	}
}

// func resourceKeyFromString(input string) resourceKey {
// 	parts := strings.Split(input, ".")
// 	if len(parts) > 1 {
// 		return resourceKey{
// 			Resource: parts[0],
// 			Field:    parts[1],
// 		}
// 	}
// 	if len(parts) == 1 {
// 		return resourceKey{
// 			Resource: parts[0],
// 		}
// 	}
// 	return resourceKey{}
// }
