package src

import (
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

func FMT(path string) error {
	err := os.Chdir(path)
	if err != nil {
		return errors.Wrap(err, "failed to navigate to dir to format")
	}

	cmd := exec.Command("gofmt", "-w", "-s", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed to format dir")
	}

	return nil
}
