package module

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/spec"
)

type Enum struct {
	_spec       spec.Enum
	Inflection  data.Inflection `yaml:"inflection,omitempty"`
	TypeData    *data.TypeData  `yaml:"type_data,omitempty"`
	Description string          `yaml:"description,omitempty"`
	Values      []string        `yaml:"values,omitempty"`
}

func makeEnums(specEnums []spec.Enum, enumPkgData data.PackageData) map[string]Enum {
	result := map[string]Enum{}
	pkgName := enumPkgData.Reference
	for _, e := range specEnums {
		inflection := data.MakeInflection(e.Name)
		result[inflection.Single.Snake] = Enum{
			_spec:       e,
			Inflection:  inflection,
			TypeData:    data.MakeTypeData(pkgName, inflection.Single.Exported, e.Type, false, false),
			Description: e.Description,
			Values:      e.Values,
		}
	}
	return result
}
