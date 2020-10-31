package project

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/spec"
	"github.com/68696c6c/capricorn/generator/models/templates/ops"
	"github.com/68696c6c/capricorn/generator/models/templates/src"
	"github.com/68696c6c/capricorn/generator/utils"
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

	rootPath string
}

func NewProject(filePath, projectPath string) (Project, error) {
	projectSpec, err := spec.NewSpec(filePath)
	if err != nil {
		return Project{}, err
	}
	mod := module.NewModuleFromSpec(projectSpec, true)
	result, err := newProjectFromModule(mod, projectPath, "", true)
	if err != nil {
		return Project{}, err
	}
	return result, nil
}

func (m Project) MustGenerate() {
	utils.PanicError(utils.CreateDir(m.rootPath))
	utils.PanicError(m.AppEnv.Generate())
	utils.PanicError(m.AppEnvExample.Generate())
	utils.PanicError(m.Gitignore.Generate())
	utils.PanicError(m.Makefile.Generate())
	utils.PanicError(m.Dockerfile.Generate())
	utils.PanicError(m.DockerCompose.Generate())
	m.SRC.MustGenerate()
}

func newProjectFromModule(m module.Module, projectPath, timestamp string, ddd bool) (Project, error) {
	rootPackage := m.Package.Reference
	rootPath, err := getProjectRootPath(projectPath, rootPackage)
	if err != nil {
		return Project{}, err
	}

	var projectSRC src.SRC
	if ddd {
		projectSRC = src.NewSRCDDD(m, rootPath, timestamp)
	} else {
		panic("not implemented yet")
	}

	result := Project{
		AppEnv:        makeEnvFile(m.Ops, rootPath, ".app.env"),
		AppEnvExample: makeEnvFile(m.Ops, rootPath, ".app.example.env"),
		Gitignore:     makeGitignore(m.Ops, rootPath),
		Makefile:      makeMakefile(m.Ops, rootPath),
		Dockerfile:    makeDockerfile(m.Ops, rootPath),
		DockerCompose: makeDockerCompose(m.Ops, rootPath),
		SRC:           projectSRC,
		rootPath:      rootPath,
	}

	return result, nil
}

func makeEnvFile(moduleOps ops.Ops, rootPath, filename string) ops.AppEnv {
	fileData, pathData := data.MakeGoFileData(rootPath, filename)
	return ops.AppEnv{
		Name: fileData,
		Path: pathData,
		Data: moduleOps,
	}
}

func makeGitignore(moduleOps ops.Ops, rootPath string) ops.Gitignore {
	fileData, pathData := data.MakeGoFileData(rootPath, ".gitignore")
	return ops.Gitignore{
		Name: fileData,
		Path: pathData,
		Data: moduleOps,
	}
}

func makeMakefile(moduleOps ops.Ops, rootPath string) ops.Makefile {
	fileData, pathData := data.MakeGoFileData(rootPath, "Makefile")
	return ops.Makefile{
		Name: fileData,
		Path: pathData,
		Data: moduleOps,
	}
}

func makeDockerfile(moduleOps ops.Ops, rootPath string) ops.Dockerfile {
	fileData, pathData := data.MakeGoFileData(rootPath, "Dockerfile")
	return ops.Dockerfile{
		Name: fileData,
		Path: pathData,
		Data: moduleOps,
	}
}

func makeDockerCompose(moduleOps ops.Ops, rootPath string) ops.DockerCompose {
	fileData, pathData := data.MakeGoFileData(rootPath, "docker-compose.yml")
	return ops.DockerCompose{
		Name: fileData,
		Path: pathData,
		Data: moduleOps,
	}
}
