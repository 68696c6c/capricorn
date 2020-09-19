package module

import (
	"fmt"

	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/spec"
)

type ResourceField struct {
	_spec      spec.ResourceField
	Key        resourceKey     `yaml:"key,omitempty"`
	Relation   data.Inflection `yaml:"relation,omitempty"`
	Name       data.Name       `yaml:"name,omitempty"`
	TypeData   *data.TypeData  `yaml:"type_data,omitempty"`
	IsRequired bool            `yaml:"is_required,omitempty"`
	IsUnique   bool            `yaml:"is_unique,omitempty"`
	IsIndexed  bool            `yaml:"is_indexed,omitempty"`
	IsPrimary  bool            `yaml:"is_primary,omitempty"`
}

type ResourceFields struct {
	Database  []*ResourceField `yaml:"goat,omitempty"`  // fields that exist in the database
	Model     []*ResourceField `yaml:"model,omitempty"` // fields that will be written into the model struct
	BelongsTo []*ResourceField `yaml:"belongs_to,omitempty"`
	HasMany   []*ResourceField `yaml:"has_many,omitempty"`
	Unique    []*ResourceField `yaml:"unique,omitempty"`
	Indexed   []*ResourceField `yaml:"indexed,omitempty"`
}

type resourceKey struct {
	Resource string `yaml:"resource,omitempty"`
	Field    string `yaml:"field,omitempty"`
}

func (r resourceKey) String() string {
	if r.Field == "" {
		return r.Resource
	}
	return fmt.Sprintf("%s.%.s", r.Resource, r.Field)
}

func makeResourceFields(specResource spec.Resource, recKey resourceKey, ddd bool) ResourceFields {
	result := ResourceFields{
		Database: []*ResourceField{
			{
				Key:       makeResourceKey(recKey.Resource, "id"),
				Name:      data.MakeName("id"),
				TypeData:  data.MakeTypeDataID(),
				IsPrimary: true,
			},
			{
				Key:      makeResourceKey(recKey.Resource, "created_at"),
				Name:     data.MakeName("created_at"),
				TypeData: data.MakeTypeDataCreatedAt(),
			},
			{
				Key:      makeResourceKey(recKey.Resource, "updated_at"),
				Name:     data.MakeName("updated_at"),
				TypeData: data.MakeTypeDataUpdatedAt(),
			},
			{
				Key:      makeResourceKey(recKey.Resource, "deleted_at"),
				Name:     data.MakeName("deleted_at"),
				TypeData: data.MakeTypeDataDeletedAt(),
			},
		},
	}

	for _, rel := range specResource.BelongsTo {
		relName := data.MakeInflection(rel)
		relFieldName := relName.Single.Snake

		fieldName := relFieldName + "_id"
		field := &ResourceField{
			Key:      makeResourceKey(recKey.Resource, fieldName),
			Name:     data.MakeName(fieldName),
			TypeData: data.MakeTypeDataID(),
		}

		// These are the foreign key field.
		result.Database = append(result.Database, field)
		result.Model = append(result.Model, field)

		// This is the field that GORM will hydrate the relation into.
		pkgName, tName := getRelationPkgAndName(relName, ddd)
		result.BelongsTo = append(result.BelongsTo, &ResourceField{
			Key:      makeResourceKey(recKey.Resource, relFieldName),
			Name:     data.MakeName(relFieldName),
			TypeData: data.MakeTypeDataBelongsTo(pkgName, tName),
		})
	}

	for _, f := range specResource.Fields {
		field := &ResourceField{
			_spec:      *f,
			Key:        makeResourceKey(recKey.Resource, f.Name),
			Name:       data.MakeName(f.Name),
			TypeData:   f.GetTypeData(),
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
		pkgName, tName := getRelationPkgAndName(relName, ddd)
		field := &ResourceField{
			Key:      makeResourceKey(recKey.Resource, relName.Plural.Exported),
			Relation: relName,
			Name:     relName.Plural,
			TypeData: data.MakeTypeDataHasMany(pkgName, tName),
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

func getRelationPkgAndName(relInflection data.Inflection, ddd bool) (string, string) {
	var pkgName string
	var tName string
	if ddd {
		pkgName = relInflection.Plural.Snake
		tName = relInflection.Single.Exported
	} else {
		tName = relInflection.Single.Exported
	}
	return pkgName, tName
}
