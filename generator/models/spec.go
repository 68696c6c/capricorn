package models

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Spec struct {
	Name      string
	Module    string
	License   string
	Author    Author
	Resources []Resource
	Commands  []ConfigCommand
}

type Author struct {
	Name         string
	Email        string
	Organization string
}

type Resource struct {
	Name      string
	BelongsTo []string `yaml:"belongs_to"`
	HasMany   []string `yaml:"has_many"`
	Fields    []ResourceField
	Actions   []string
	Custom    []string
}

type ResourceField struct {
	Name     string
	Type     string
	Required bool
}

type ConfigCommand struct {
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
