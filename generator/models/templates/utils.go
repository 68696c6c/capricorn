package templates

import (
	"fmt"

	"github.com/68696c6c/capricorn/generator/utils"
)

type Template interface {
	Generate() error
}

// MustParse is for rendering a template inside of another template.  Since it is called inside of a template, there is
// no way to handle an error, leaving no option but to panic, hence "must" in the name.
type SubTemplate interface {
	MustParse() string
}

// e.g. base: example
// e.g. full: example.go
type FileData struct {
	Full string `yaml:"full"`
	Base string `yaml:"base"`
}

// e.g. base: src/app/domain/
// e.g. full: src/app/domain/example.go
type PathData FileData

func MakeGoFileData(basePath, fileBaseName string) (FileData, PathData) {
	f := FileData{
		Full: fmt.Sprintf("%s.go", fileBaseName),
		Base: fileBaseName,
	}
	p := PathData{
		Full: utils.JoinPath(basePath, f.Full),
		Base: f.Full,
	}
	return f, p
}
