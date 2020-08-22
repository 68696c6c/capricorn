package project

import (
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/ops"
	"github.com/68696c6c/capricorn/generator/models/templates/src"
)

type Project struct {
	// Ops files
	AppEnv        ops.AppEnv        `yaml:"app_env"`
	AppEnvExample ops.AppEnv        `yaml:"app_env_example"`
	Gitignore     ops.Gitignore     `yaml:"gitignore"`
	Makefile      ops.Makefile      `yaml:"makefile"`
	Dockerfile    ops.Dockerfile    `yaml:"dockerfile"`
	DockerCompose ops.DockerCompose `yaml:"docker_compose"`
	// GoMod ops.GoMod `yaml:"go_mod"`

	SRC src.SRC `yaml:"src"`
}

func NewProjectFromModule(m module.Module, projectPath string, ddd bool) (Project, error) {
	// @TODO @CHECKPOINT working on generating project from a module.
	// rootPackage := m.Path.Full
	// rootPath, err := getProjectRootPath(projectPath, rootPackage)
	// if err != nil {
	// 	return Project{}, err
	// }
	//
	// var projectSRC src.SRC
	// if ddd {
	// 	projectSRC = src.NewSRCDDD(m, rootPath)
	// } else {
	//
	// }
	//
	// result := Project{
	// 	SRC: projectSRC,
	// }
	return Project{}, nil
}
