package data

import (
	"fmt"

	"github.com/68696c6c/capricorn/generator/utils"
)

const ImportGoat = "github.com/68696c6c/goat"
const ImportQuery = "github.com/68696c6c/goat/query"
const ImportGin = "github.com/gin-gonic/gin"

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
	Base string `yaml:"base"`
	Full string `yaml:"full"`
}

// e.g. base: src/app/domain/
// e.g. full: src/app/domain/example.go
type PathData FileData

// e.g. reference: pkgname
// e.g. base: github.com/user/example/src
// e.g. full: github.com/user/example/src/pkgname
type PackageData struct {
	Reference string   `yaml:"reference"`
	Name      Name     `yaml:"name"`
	Path      PathData `yaml:"path"`
}

func (m PackageData) GetImport() string {
	return m.Path.Full
}

func (m PackageData) GetReference() string {
	return m.Reference
}

func MakePathData(basePath, name string) PathData {
	return PathData{
		Base: basePath,
		Full: utils.JoinPath(basePath, name),
	}
}

func MakeGoFileData(basePath, fileBaseName string) (FileData, PathData) {
	f := FileData{
		Base: fileBaseName,
		Full: fmt.Sprintf("%s.go", fileBaseName),
	}
	p := MakePathData(basePath, f.Full)
	return f, p
}

func MakePackageData(pkgBase, pkgName string) PackageData {
	name := MakeName(pkgName)
	return PackageData{
		Name:      name,
		Reference: name.Snake,
		Path:      MakePathData(pkgBase, pkgName),
	}
}