package spec

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Spec struct {
	Name      string     `yaml:"name"`
	Module    string     `yaml:"module"`
	License   string     `yaml:"license"`
	Author    Author     `yaml:"author"`
	Resources []Resource `yaml:"resources"`
	Commands  []Command  `yaml:"commands"`
}

type Author struct {
	Name         string `yaml:"name"`
	Email        string `yaml:"email"`
	Organization string `yaml:"organization"`
}

type Resource struct {
	Name      string
	BelongsTo []string        `yaml:"belongs_to"`
	HasMany   []string        `yaml:"has_many"`
	Fields    []ResourceField `yaml:"fields"`
	Actions   []string        `yaml:"actions"`
	Custom    []string        `yaml:"custom"`
}

type ResourceField struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Required bool   `yaml:"required"`
	Unique   bool   `yaml:"unique"`
}

type Command struct {
	Name string   `yaml:"name"`
	Args []string `yaml:"args"`
	Use  string   `yaml:"use"`
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

	return spec, nil
}
