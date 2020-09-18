package spec

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// type EnumData map[string]map[string]string
type EnumData map[string]*Enum

// func (m *Enum) UnmarshalYAML(unmarshal func(interface{}) error) error {
//
// }

func (m *EnumData) UnmarshalYAML(unmarshal func(interface{}) error) error {
	result := EnumData{}
	var enumName string
	var enumType string
	var enumDescription string
	var enumValues []string
	t := make(map[string]map[string]string)
	err := unmarshal(t)
	if err != nil {
		// enumType =
		println("---ERROR", fmt.Sprintf("%v", t))
		tv := make(map[string]map[string][]string)
		err := unmarshal(tv)
		if err != nil {
			println("duh")
			for key, value := range tv {
				println("TV KEY", key)
				for k, v := range value {
					println("-------TV FIELD", k, fmt.Sprintf("%v", v))
					switch k {
					case "values":
						enumValues = v
					}
				}
			}
		}
	}

	for key, value := range t {
		enumName = key
		println("-----ENUM NAME", enumName)
		for k, v := range value {
			println("-------ENUM FIELD", k)
			switch k {
			case "description":
				enumDescription = v
			case "type":
				enumType = v
			}
		}
		result[enumName] = &Enum{
			Name:        enumName,
			Description: enumDescription,
			Type:        enumType,
			Values:      enumValues,
		}
	}

	*m = result

	for key, value := range *m {
		println("final enum data", key, fmt.Sprintf("%v", value))
	}

	return nil
}

type ResourceFields []ResourceField

func (m *ResourceFields) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var result ResourceFields
	var t []map[string]string
	err := unmarshal(&t)
	if err != nil {
		println("ResourceFields error", fmt.Sprintf("%v", t))
		for _, field := range t {
			var fieldName string
			var fieldType string
			for key, value := range field {
				println("ResourceFields field", key, value)
				switch key {
				case "name":
					fieldName = value
				case "type":
					fieldType = value
				}
			}
			result = append(result, ResourceField{
				Name:     fieldName,
				Type:     fieldType,
				Required: false,
				Unique:   false,
				Indexed:  false,
			})
		}
	}

	*m = result

	for _, value := range *m {
		println("final ResourceFields data", fmt.Sprintf("%v", value))
	}

	return nil
}

// func (m *ResourceField) UnmarshalYAML(unmarshal func(interface{}) error) error {
// 	println()
// 	println(fmt.Sprintf("ResourceField: %v", m))
// 	var fieldName string
// 	var fieldType string
//
// 	// println(m.Name, "|", m.Type)
// 	t := make(map[string]string)
// 	err := unmarshal(&t)
// 	if err != nil {
// 		println("error", fmt.Sprintf("%v", t))
//
// 		enumData := EnumData{}
// 		err := unmarshal(&enumData)
// 		if err != nil {
// 			return err
// 		}
// 		println("!!!", fmt.Sprintf("%v", enumData))
// 		for enumKey, enum := range enumData {
// 			println("!!!! enumData", enumKey, "|", fmt.Sprintf("%v", enum))
// 			// switch enumKey {
// 			// case "name":
// 			// 	m.Name = enumKey
// 			// case "type":
// 			m.Type = enum.GetType()
// 			// }
// 		}
//
// 		// t2 := make(map[string]*Enum)
// 		// err := unmarshal(&t2)
// 		// if err != nil {
// 		// 	for tk, tv := range t2 {
// 		// 		println("!!!! TK TYPE", tk, "|", tv.GetType())
// 		// 		switch tk {
// 		// 		case "name":
// 		// 			m.Name = tk
// 		// 		case "type":
// 		// 			m.Type = tv.GetType()
// 		// 		}
// 		// 	}
// 		// }
//
// 	} else {
// 		println(fmt.Sprintf("%v", t))
// 		fieldName = t["name"]
// 		fieldType = t["type"]
// 	}
// 	m.Name = fieldName
// 	m.Type = fieldType
// 	println("result", fmt.Sprintf("%v", m))
// 	return nil
// }

type Spec struct {
	enums   []Enum
	Name    string     `yaml:"name,omitempty"`
	Module  string     `yaml:"module,omitempty"`
	License string     `yaml:"license,omitempty"`
	Author  Author     `yaml:"author,omitempty"`
	Enums   []EnumData `yaml:"enums,omitempty"`
	// Enums     []map[string]*Enum `yaml:"enums,omitempty"`
	Commands  []Command  `yaml:"commands,omitempty"`
	Resources []Resource `yaml:"resources,omitempty"`
}

// func (m *Enum) UnmarshalYAML(unmarshal func(interface{}) error) error {
//
// 	err := yaml.Unmarshal(file, &spec)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

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
	Name        string
	Description string   `yaml:"description,omitempty"`
	Type        string   `yaml:"type,omitempty"`
	Values      []string `yaml:"values,omitempty"`
}

func (m *Enum) GetType() string {
	return m.Name
}

type Resource struct {
	Name      string
	BelongsTo []string       `yaml:"belongs_to,omitempty"`
	HasMany   []string       `yaml:"has_many,omitempty"`
	Fields    ResourceFields `yaml:"fields,omitempty"`
	Actions   []string       `yaml:"actions,omitempty"`
	Custom    []string       `yaml:"custom,omitempty"`
}

type ResourceFieldYAML map[string]string

type ResourceField struct {
	Name     string `yaml:"name,omitempty"`
	Type     string `yaml:"type,omitempty"`
	Required bool   `yaml:"required,omitempty"`
	Unique   bool   `yaml:"unique,omitempty"`
	Indexed  bool   `yaml:"indexed,omitempty"`
}

// func (m ResourceField) String() string {
// 	return fmt.Sprintf("%s: %s", m.Name, m.Type)
// }

func (m Spec) String() string {
	out, err := yaml.Marshal(&m)
	if err != nil {
		return "failed to marshal spec to yaml"
	}
	return string(out)
}

func NewSpec(filePath string) (Spec, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return Spec{}, errors.Wrap(err, "failed to read spec file")
	}

	spec := Spec{}
	err = yaml.Unmarshal(file, &spec)
	if err != nil {
		return Spec{}, errors.Wrap(err, "failed to unmarshal spec")
	}

	for _, r := range spec.Resources {
		println(r.Fields)
	}

	return spec, nil
}
