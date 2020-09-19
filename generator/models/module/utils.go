package module

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/spec"
)

func GetFixtureModule() Module {
	f := spec.GetFixtureSpec()
	return NewModuleFromSpec(f, true)
}

func GetFixtureResourceField(recNameKebob, fieldNameKebob string) ResourceField {
	return ResourceField{
		Key:      makeResourceKey(recNameKebob, fieldNameKebob),
		Relation: data.MakeInflection("relation"),
		Name:     data.MakeName(fieldNameKebob),
		TypeData: data.MakeTypeDataString(),
	}
}

const FixtureModuleYAML = ``
