package spec

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Spec struct {
	Name      string     `yaml:"name,omitempty"`
	Module    string     `yaml:"module,omitempty"`
	License   string     `yaml:"license,omitempty"`
	Author    Author     `yaml:"author,omitempty"`
	Commands  []Command  `yaml:"commands,omitempty"`
	Enums     []Enum     `yaml:"enums,omitempty"`
	Resources []Resource `yaml:"resources,omitempty"`
}

type Author struct {
	Name         string `yaml:"name,omitempty,omitempty"`
	Email        string `yaml:"email,omitempty,omitempty"`
	Organization string `yaml:"organization,omitempty"`
}

type Command struct {
	Name string   `yaml:"name,omitempty"`
	Args []string `yaml:"args,omitempty"`
	Use  string   `yaml:"use,omitempty"`
}

type Enum struct {
	Name        string   `yaml:"name,omitempty"`
	Description string   `yaml:"description,omitempty"`
	Type        string   `yaml:"type,omitempty"`
	Values      []string `yaml:"values,omitempty"`
}

type Resource struct {
	Name      string
	BelongsTo []string         `yaml:"belongs_to,omitempty"`
	HasMany   []string         `yaml:"has_many,omitempty"`
	Fields    []*ResourceField `yaml:"fields,omitempty"`
	Actions   []string         `yaml:"actions,omitempty"`
	Custom    []string         `yaml:"custom,omitempty"`
}

type ResourceField struct {
	Name     string `yaml:"name,omitempty"`
	Type     string `yaml:"type,omitempty"`
	Enum     string `yaml:"enum,omitempty"`
	Required bool   `yaml:"required,omitempty"`
	Unique   bool   `yaml:"unique,omitempty"`
	Indexed  bool   `yaml:"indexed,omitempty"`
}

func (m *ResourceField) GetTypeData() *data.TypeData {
	if m.Type == m.Enum {
		return nil
	}
	return data.NewTypeDataFromReference(m.Type)
}

func (m Spec) String() string {
	out, err := yaml.Marshal(&m)
	if err != nil {
		return "failed to marshal spec to yaml"
	}
	return string(out)
}

func (m *ResourceField) UnmarshalYAML(unmarshal func(interface{}) error) error {
	result := map[string]string{}
	err := unmarshal(result)
	if err != nil {
		return err
	}

	var resultName string
	var resultType string
	var resultEnum string
	var resultRequired bool
	var resultUnique bool
	var resultIndexed bool
	for key, value := range result {
		println("result key", key)
		switch key {
		case "name":
			resultName = value
		case "type":
			resultType = value
		case "required":
			println("result required", value)
			if value == "true" {
				resultRequired = true
			}
		case "enum":
			resultEnum = value
		case "unique":
			resultUnique = true
		case "indexed":
			resultIndexed = true
		}
	}

	if resultEnum != "" {
		resultType = resultEnum
	}

	*m = ResourceField{
		Enum:     resultEnum,
		Name:     resultName,
		Type:     resultType,
		Required: resultRequired,
		Unique:   resultUnique,
		Indexed:  resultIndexed,
	}

	return nil
}

func NewSpec(filePath string) (Spec, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return Spec{}, errors.Wrap(err, "failed to read spec file")
	}

	result := Spec{}
	err = yaml.Unmarshal(file, &result)
	if err != nil {
		return Spec{}, errors.Wrap(err, "failed to unmarshal spec")
	}

	enumMap := map[string]Enum{}
	for _, e := range result.Enums {
		enumMap[e.Name] = e
	}

	for _, r := range result.Resources {
		for _, f := range r.Fields {
			if f.Type == f.Enum {
				e, ok := enumMap[f.Enum]
				if !ok {
					return Spec{}, errors.New("failed to map enum field")
				}
				f.Type = e.Type
			}
		}
	}

	return result, nil
}
