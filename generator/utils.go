package generator

import (
	"bytes"
	"fmt"
	"go/build"
	"os"
	"os/exec"
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

	cmd := exec.Command("go", "mod", "init")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed to init go modules")
	}

	cmd = exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed to run go mod tidy")
	}

	cmd = exec.Command("go", "mod", "vendor")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed to run go mod vendor")
	}

	return nil
}
