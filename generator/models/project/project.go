package project

import (
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/ops"
	"github.com/68696c6c/capricorn/generator/models/templates/src"
)

type Project struct {
	// Ops files
	AppEnv        ops.AppEnv        `yaml:"app_env,omitempty"`
	AppEnvExample ops.AppEnv        `yaml:"app_env_example,omitempty"`
	Gitignore     ops.Gitignore     `yaml:"gitignore,omitempty"`
	Makefile      ops.Makefile      `yaml:"makefile,omitempty"`
	Dockerfile    ops.Dockerfile    `yaml:"dockerfile,omitempty"`
	DockerCompose ops.DockerCompose `yaml:"docker_compose,omitempty"`
	// GoMod ops.GoMod `yaml:"go_mod,omitempty"`

	SRC src.SRC `yaml:"src,omitempty"`
}

func NewProjectFromModule(m module.Module, projectPath string, ddd bool) (Project, error) {
	// @TODO @CHECKPOINT working on generating project from a module.
	rootPackage := m.Package.Reference
	rootPath, err := getProjectRootPath(projectPath, rootPackage)
	if err != nil {
		return Project{}, err
	}

	var projectSRC src.SRC
	if ddd {
		projectSRC = src.NewSRCDDD(m, rootPath)
	} else {

	}

	result := Project{
		SRC: projectSRC,
	}

	return result, nil
}
