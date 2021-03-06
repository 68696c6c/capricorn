package utils

import (
	"bytes"
	"fmt"
	"go/build"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

func JoinPath(parts ...string) string {
	return strings.Join(parts, "/")
}

func GetProjectPath() (projectPath string, err error) {
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		goPath = build.Default.GOPATH
	}
	path := JoinPath(goPath, "src")
	projectPath, err = filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return
}

func CreateDir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func CopyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func AppendFileText(fileName, text string) error {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return errors.Wrap(err, "failed to open file for appending")
	}
	defer f.Close()
	if _, err = f.WriteString(text); err != nil {
		return errors.Wrap(err, "failed to write to file")
	}
	return nil
}

func GenerateFile(basePath, fileName, fileTemplate string, data interface{}) error {
	t := template.Must(template.New(fileName).Parse(fileTemplate))

	filePath := fmt.Sprintf("%s/%s", basePath, fileName)
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

func ParseTemplateToString(name, temp string, data interface{}) (string, error) {
	var tpl bytes.Buffer
	t := template.Must(template.New(name).Parse(temp))
	err := t.Execute(&tpl, data)
	if err != nil {
		return "", errors.Wrapf(err, "failed parse template '%s'", name)
	}
	return tpl.String(), nil
}

func separatedToCamel(input string, leadingCap bool) string {
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
				if v == '_' || v == '-' {
					isToUpper = true
				} else {
					output += string(v)
				}
			}
		}
	}
	return output
}

func SeparatedToUnexported(input string) string {
	return separatedToCamel(input, false)
}

func SeparatedToExported(input string) string {
	return separatedToCamel(input, true)
}

func separatedToSeparated(input string, separator rune) string {
	var output string
	for _, v := range input {
		if v == '_' || v == '-' {
			output += string(separator)
		} else {
			output += strings.ToLower(string(v))
		}
	}
	return output
}

func SeparatedToSnake(input string) string {
	return separatedToSeparated(input, '_')
}

func SeparatedToKebob(input string) string {
	return separatedToSeparated(input, '-')
}
