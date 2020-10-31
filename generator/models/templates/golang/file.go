package golang

import (
	"fmt"
	"github.com/68696c6c/capricorn/generator/utils"
	"strings"

	"github.com/68696c6c/capricorn/generator/models/data"
)

type SourceFile interface {
	MustGetFile() File
}

type File struct {
	Name    data.FileData    `yaml:"name,omitempty"`
	Path    data.PathData    `yaml:"path,omitempty"`
	Package data.PackageData `yaml:"package,omitempty"`

	Imports      Imports     `yaml:"imports,omitempty"`
	InitFunction Function    `yaml:"init_function,omitempty"`
	Consts       []Const     `yaml:"consts,omitempty"`
	Vars         []Var       `yaml:"vars,omitempty"`
	Interfaces   []Interface `yaml:"interfaces,omitempty"`
	TypeAliases  []Value     `yaml:"type_aliases,omitempty"`
	Structs      []Struct    `yaml:"structs,omitempty"`
	Functions    []Function  `yaml:"functions,omitempty"`
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

func (m File) MustParseTypeAliases() string {
	var result []string
	for _, v := range m.TypeAliases {
		result = append(result, "type "+v.MustParse())
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
	if len(m.TypeAliases) > 0 {
		sections = append(sections, m.MustParseTypeAliases())
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

func (m File) Generate() error {
	err := utils.WriteFile(m.Path.Base, m.Name.Full, m.MustParse())
	if err != nil {
		return err
	}
	return nil
}
