package generator

import (
	"bytes"
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

func joinPath(parts ...string) string {
	return strings.Join(parts, "/")
}

func getProjectPath() (projectPath string, err error) {
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		goPath = build.Default.GOPATH
	}
	path := joinPath(goPath, "src")
	projectPath, err = filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return
}

func createDir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func generateFile(basePath, fileName, fileTemplate string, data interface{}) error {
	t := template.Must(template.New(fileName).Parse(fileTemplate))

	filePath := fmt.Sprintf("%s/%s.go", basePath, fileName)
	f, err := os.Create(filePath)
	if err != nil {
		return errors.Wrapf(err, "failed create file '%s'", filePath)
	}

	err = t.Execute(f, data)
	if err != nil {
		return errors.Wrapf(err, "failed write file '%s'", filePath)
	}

	err = f.Close()
	if err != nil {
		return errors.Wrapf(err, "failed to close file '%s'", filePath)
	}

	return nil
}

func parseTemplateToString(name, temp string, data interface{}) (string, error) {
	var tpl bytes.Buffer
	t := template.Must(template.New(name).Parse(temp))
	err := t.Execute(&tpl, data)
	if err != nil {
		return "", errors.Wrapf(err, "failed parse template '%s'", name)
	}
	return tpl.String(), nil
}

func snakeToCamel(input string, leadingCap bool) string {
	isToUpper := false
	var output string
	for k, v := range input {
		if k == 0 && leadingCap {
			output = strings.ToUpper(string(input[0]))
		} else {
			if isToUpper {
				output += strings.ToUpper(string(v))
				isToUpper = false
			} else {
				if v == '_' {
					isToUpper = true
				} else {
					output += string(v)
				}
			}
		}
	}
	return output
}

func snakeToUnexportedName(input string) string {
	return snakeToCamel(input, false)
}

func snakeToExportedName(input string) string {
	return snakeToCamel(input, true)
}
