package project

import "github.com/pkg/errors"

func CreateProject(spec Spec) error {
	err := createDir(spec.Paths.Root)
	if err != nil {
		return errors.Wrapf(err, "failed to create project directory '%s'", spec.Paths.Root)
	}
	return nil
}
