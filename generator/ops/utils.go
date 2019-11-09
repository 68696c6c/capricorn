package ops

import (
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

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

	cmd = exec.Command("make", "deps")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed to install dependancies")
	}

	return nil
}
