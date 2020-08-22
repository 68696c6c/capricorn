package project

import (
	"github.com/68696c6c/capricorn/generator/utils"

	"github.com/pkg/errors"
)

func getProjectRootPath(projectPath, projectModule string) (string, error) {
	rootPath := projectPath
	if projectPath == "" {
		projectPath, err := utils.GetProjectPath()
		if err != nil {
			return "", errors.Wrap(err, "failed to determine project root path")
		}
		rootPath = utils.JoinPath(projectPath, projectModule)
	}
	return rootPath, nil
}
