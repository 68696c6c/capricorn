package src

import (
	"os"
	"os/exec"

	"github.com/68696c6c/capricorn/generator/models"
	"github.com/68696c6c/capricorn/generator/models/templates/ops"
	"github.com/68696c6c/capricorn/generator/utils"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func CreateDocker(spec models.Project, logger *logrus.Logger) error {
	logPrefix := "CreateDocker | "
	logger.Debug(logPrefix, "generating docker files")

	err := utils.CreateDir(spec.Paths.Docker)
	if err != nil {
		return errors.Wrapf(err, "failed to create project directory '%s'", spec.Paths.Docker)
	}

	err = utils.GenerateFile(spec.Paths.Docker, "Dockerfile", ops.DockerfileTemplate, spec)
	if err != nil {
		return errors.Wrap(err, "failed to create Dockerfile")
	}

	err = utils.GenerateFile(spec.Paths.Root, "docker-compose.yml", ops.DockerComposeTemplate, spec)
	if err != nil {
		return errors.Wrap(err, "failed to create docker-compose.yml")
	}

	err = utils.GenerateFile(spec.Paths.Root, ".app.env", ops.AppEnvTemplate, spec)
	if err != nil {
		return errors.Wrap(err, "failed to create .app.env")
	}

	err = utils.GenerateFile(spec.Paths.Root, ".app.example.env", ops.AppEnvTemplate, spec)
	if err != nil {
		return errors.Wrap(err, "failed to create .app.env")
	}

	return nil
}

func CreateMakefile(spec models.Project, logger *logrus.Logger) error {
	logPrefix := "CreateMakefile | "
	logger.Debug(logPrefix, "generating makefile")

	err := utils.GenerateFile(spec.Paths.Root, "Makefile", ops.MakefileTemplate, spec)
	if err != nil {
		return errors.Wrap(err, "failed to create Makefile")
	}

	return nil
}

func InitModule(spec models.Project) error {
	path := spec.Paths.Root
	err := os.Chdir(path)
	if err != nil {
		return errors.Wrap(err, "failed to navigate to dir to initialize module")
	}

	err = os.Setenv("GO111MODULE", "on")
	if err != nil {
		return errors.Wrap(err, "failed to enable go modules")
	}
	defer os.Unsetenv("GO111MODULE")

	cmd := exec.Command("make", "image-local")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed to make image")
	}

	cmd = exec.Command("go", "mod", "init", spec.Spec.Module)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed to init go modules")
	}

	err = utils.AppendFileText(path+"/go.mod", `
replace github.com/ugorji/go v1.1.4 => github.com/ugorji/go v0.0.0-20190204201341-e444a5086c43
`)
	if err != nil {
		return errors.Wrap(err, "failed to update go.mod")
	}

	cmd = exec.Command("make", "deps")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed to install dependencies")
	}

	err = utils.GenerateFile(path, ".gitignore", ops.GitignoreTemplate, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create .gitignore")
	}

	return nil
}
