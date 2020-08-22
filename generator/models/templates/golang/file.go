package golang

import (
	"fmt"

	"github.com/68696c6c/capricorn/generator/models/utils"

	"strings"
)

type File struct {
	Name    utils.FileData    `yaml:"name"`
	Path    utils.PathData    `yaml:"path"`
	Package utils.PackageData `yaml:"package"`

	Imports      Imports     `yaml:"imports"`
	InitFunction Function    `yaml:"init_function"`
	Consts       []Const     `yaml:"consts"`
	Vars         []Var       `yaml:"vars"`
	Interfaces   []Interface `yaml:"interfaces"`
	Structs      []Struct    `yaml:"structs"`
	Functions    []Function  `yaml:"functions"`
}

func (m File) MustParseConsts() string {
	var result []string
	for _, v := range m.Consts {
		result = append(result, v.MustParse())
	}
	return strings.Join(result, "\n\n")
}

func (m File) MustParseVars() string {
	var result []string
	for _, v := range m.Vars {
		result = append(result, v.MustParse())
	}
	return strings.Join(result, "\n\n")
}

func (m File) MustParseInterfaces() string {
	var result []string
	for _, v := range m.Interfaces {
		result = append(result, v.MustParse())
	}
	return strings.Join(result, "\n\n")
}

func (m File) MustParseStructs() string {
	var result []string
	for _, v := range m.Structs {
		result = append(result, v.MustParse())
	}
	return strings.Join(result, "\n\n")
}

func (m File) MustParseFunctions() string {
	var result []string
	for _, v := range m.Functions {
		result = append(result, v.MustParse())
	}
	return strings.Join(result, "\n\n")
}

// This is only used for testing.
func (m File) MustParse() string {
	var sections []string

	if m.Imports.HasImports() {
		sections = append(sections, m.Imports.MustParse())
	}
	if m.InitFunction.Body != "" {
		sections = append(sections, m.InitFunction.MustParse())
	}
	if len(m.Consts) > 0 {
		sections = append(sections, m.MustParseConsts())
	}
	if len(m.Vars) > 0 {
		sections = append(sections, m.MustParseVars())
	}
	if len(m.Interfaces) > 0 {
		sections = append(sections, m.MustParseInterfaces())
	}
	if len(m.Structs) > 0 {
		sections = append(sections, m.MustParseStructs())
	}
	if len(m.Functions) > 0 {
		sections = append(sections, m.MustParseFunctions())
	}

	result := []string{fmt.Sprintf("package %s\n", m.Package.GetReference())}

	// Separate each section with an additional line break.
	result = append(result, strings.Join(sections, "\n\n\n"))

	return strings.Join(result, "\n") + "\n"
}