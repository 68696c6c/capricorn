package golang

import "github.com/68696c6c/capricorn/generator/utils"

var templateFile = `package {{ .Package.Name }}

import (
	{{- range $key, $value := .Imports.Standard }}
	"{{ $value }}"
	{{- end }}{{ println }}

	{{- range $key, $value := .Imports.App }}
	"{{ $value }}"
	{{- end }}{{ println }}

	{{- range $key, $value := .Imports.Vendor }}
	"{{ $value }}"
	{{- end }}
)

{{- if .InitFunction.Body }}
{{- println }}
{{- println }}
{{ .InitFunction.MustParse }}
{{- println }}
{{- else }}
{{- println }}
{{- end }}

{{- $length := len .Consts }}
{{- if gt $length 0 }}
{{- println }}
{{- range $key, $value := .Consts }}
const {{ $value.Name }} = {{ $value.Value }}
{{- println }}
{{- end }}
{{- end }}

{{- $length := len .Vars }}
{{- if gt $length 0 }}
{{- println }}
{{- range $key, $value := .Vars }}
var {{ $value.Name }} = {{ $value.Value }}
{{- println }}
{{- end }}
{{- end }}

{{- $length := len .Interfaces }}
{{- if gt $length 0 }}
{{- println }}
{{- range $key, $value := .Interfaces }}
{{ $value.MustParse }}
{{- println }}
{{- end }}
{{- end }}

{{- $length := len .Interfaces }}
{{- if gt $length 0 }}
{{- println }}
{{- range $key, $value := .Structs }}
{{ $value.MustParse }}
{{- println }}
{{- end }}
{{- end }}

{{- $length := len .Interfaces }}
{{- if gt $length 0 }}
{{- println }}
{{- range $key, $value := .Functions }}
{{ $value.MustParse }}
{{- println }}
{{- end }}
{{- end }}`

type File struct {
	Name    FileData    `yaml:"name"`
	Path    FileData    `yaml:"path"`
	Package PackageData `yaml:"package"`

	Imports      FileImports `yaml:"imports"`
	InitFunction Function    `yaml:"init_function"`
	Consts       []Value     `yaml:"consts"`
	Vars         []Value     `yaml:"vars"`
	Interfaces   []Interface `yaml:"interfaces"`
	Structs      []Struct    `yaml:"structs"`
	Functions    []Function  `yaml:"functions"`
}

type FileData struct {
	// e.g. path: src/app/domain/example.go
	// e.g. file: example.go
	Full string `yaml:"full"`

	// e.g. path: src/app/domain/
	// e.g. file: example
	Base string `yaml:"base"`
}

type PackageData struct {
	Name   string `yaml:"name"`   // e.g. domain
	Module string `yaml:"module"` // e.g. github.com/example/src/app/domain
	Local  string `yaml:"local"`  // e.g. src/app/domain
}

type FileImports struct {
	Standard []string `yaml:"standard"`
	App      []string `yaml:"app"`
	Vendor   []string `yaml:"vendor"`
}

type Value struct {
	Name  string `yaml:"name"`
	Type  string `yaml:"type"`
	Value string `yml:"value"`
}

func (m File) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_file", templateFile, m)
	if err != nil {
		panic(err)
	}
	return result
}
