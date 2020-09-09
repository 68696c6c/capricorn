package module

import (
	"fmt"

	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/spec"
)

type ResourceField struct {
	_spec      spec.ResourceField
	Key        resourceKey     `yaml:"key"`
	Relation   data.Inflection `yaml:"relation"`
	Name       data.Name       `yaml:"name"`
	Type       string          `yaml:"type"`
	IsRequired bool            `yaml:"is_required"`
	IsUnique   bool            `yaml:"is_unique"`
	IsIndexed  bool            `yaml:"is_indexed"`
	IsPrimary  bool            `yaml:"is_primary"`
}

type ResourceFields struct {
	Database  []ResourceField `yaml:"goat,omitempty"`  // fields that exist in the database
	Model     []ResourceField `yaml:"model,omitempty"` // fields that will be written into the model struct
	BelongsTo []ResourceField `yaml:"belongs_to,omitempty"`
	HasMany   []ResourceField `yaml:"has_many,omitempty"`
	Unique    []ResourceField `yaml:"unique,omitempty"`
	Indexed   []ResourceField `yaml:"indexed,omitempty"`
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

func makeResourceFields(specResource spec.Resource, recKey resourceKey, ddd bool) ResourceFields {
	result := ResourceFields{
		Database: []ResourceField{
			{
				Key:       makeResourceKey(recKey.Resource, "id"),
				Name:      data.MakeName("id"),
				Type:      "goat.ID",
				IsPrimary: true,
			},
			{
				Key:  makeResourceKey(recKey.Resource, "created_at"),
				Name: data.MakeName("created_at"),
				Type: "time.Time",
			},
			{
				Key:  makeResourceKey(recKey.Resource, "updated_at"),
				Name: data.MakeName("updated_at"),
				Type: "*time.Time",
			},
			{
				Key:  makeResourceKey(recKey.Resource, "deleted_at"),
				Name: data.MakeName("deleted_at"),
				Type: "*time.Time",
			},
		},
	}

	for _, rel := range specResource.BelongsTo {
		relName := data.MakeInflection(rel)
		relFieldName := relName.Single.Snake

		fieldName := relFieldName + "_id"
		field := ResourceField{
			Key:  makeResourceKey(recKey.Resource, fieldName),
			Name: data.MakeName(fieldName),
			Type: "goat.ID",
		}

		// These are the foreign key field.
		result.Database = append(result.Database, field)
		result.Model = append(result.Model, field)

		// This is the field that GORM will hydrate the relation into.

		result.BelongsTo = append(result.BelongsTo, ResourceField{
			Key:  makeResourceKey(recKey.Resource, relFieldName),
			Name: data.MakeName(relFieldName),
			Type: makeBelongsToType(relName, ddd),
		})
	}

	for _, f := range specResource.Fields {
		field := ResourceField{
			_spec:      f,
			Key:        makeResourceKey(recKey.Resource, f.Name),
			Name:       data.MakeName(f.Name),
			Type:       f.Type,
			IsRequired: f.Required,
			IsUnique:   f.Unique,
			IsIndexed:  f.Indexed,
		}
		result.Database = append(result.Database, field)
		result.Model = append(result.Model, field)
		if f.Unique {
			result.Unique = append(result.Unique, field)
		}
		if f.Indexed {
			result.Indexed = append(result.Indexed, field)
		}
	}

	for _, rel := range specResource.HasMany {
		relName := data.MakeInflection(rel)
		field := ResourceField{
			Key:      makeResourceKey(recKey.Resource, relName.Plural.Exported),
			Relation: relName,
			Name:     relName.Plural,
			Type:     makeHasManyType(relName, ddd),
		}

		// This is the field that GORM will hydrate the relation into.
		// HasMany relations built into the parent model, not the child, so nothing else to do for now.
		result.HasMany = append(result.HasMany, field)
	}

	return result
}

func makeResourceKey(resource, field string) resourceKey {
	return resourceKey{
		Resource: resource,
		Field:    field,
	}
}

func makeHasManyType(relInflection data.Inflection, ddd bool) string {
	var t string
	if ddd {
		t = fmt.Sprintf("%s.%s", relInflection.Plural.Snake, relInflection.Single.Exported)
	} else {
		t = relInflection.Single.Exported
	}
	return fmt.Sprintf("[]*%s", t)
}

func makeBelongsToType(relInflection data.Inflection, ddd bool) string {
	var t string
	if ddd {
		t = fmt.Sprintf("%s.%s", relInflection.Plural.Snake, relInflection.Single.Exported)
	} else {
		t = relInflection.Single.Exported
	}
	return fmt.Sprintf("*%s", t)
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
