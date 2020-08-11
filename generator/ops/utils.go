package ops

import (
	"os"
	"os/exec"

	"github.com/68696c6c/capricorn/generator/utils"

	"github.com/pkg/errors"
)

const gitIgnoreTemplate = `
.idea
vendor
.app.env
}`

func InitModule(path string) error {
	err := os.Chdir(path)
	if err != nil {
		return errors.Wrap(err, "failed to navigate to dir to initialize module")
	}

	err = os.Setenv("GO111MODULE", "on")
	if err != nil {
		return errors.Wrap(err, "failed to enable go modules")
	}
	defer os.Unsetenv("GO111MODULE")

	cmd := exec.Command("make", "image")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed to make image")
	}

	cmd = exec.Command("make", "init")
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

	err = utils.GenerateFile(path, ".gitignore", gitIgnoreTemplate, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create .gitignore")
	}

	return nil
}
