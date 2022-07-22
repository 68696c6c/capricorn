package golang

import (
	"fmt"
	"strings"
)

type File struct {
	Name        string
	Extension   string
	Contents    string
	Package     string
	Imports     Imports
	Vars        []*Var
	TypeAliases []*Value
	Interfaces  []*Interface
	Structs     []*Struct
	Functions   []*Function
}

func (f *File) GetFullName() string {
	if f.Extension == "" {
		return f.Name
	}
	return fmt.Sprintf("%s.%s", f.Name, f.Extension)
}

func (f *File) SetName(name string) *File {
	f.Name = name
	return f
}

func (f *File) GetName() string {
	return f.Name
}

func (f *File) SetPackage(pkgName string) *File {
	f.Package = pkgName
	return f
}

func (f *File) GetPackage() string {
	return f.Package
}

func (f *File) SetImports(imports Imports) *File {
	f.Imports = imports
	return f
}

func (f *File) GetImports() Imports {
	return f.Imports
}

func (f *File) SetVars(vars []*Var) *File {
	f.Vars = vars
	return f
}

func (f *File) GetVars() []*Var {
	return f.Vars
}

func (f *File) SetTypeAliases(typeAliases []*Value) *File {
	f.TypeAliases = typeAliases
	return f
}

func (f *File) GetTypeAliases() []*Value {
	return f.TypeAliases
}

func (f *File) SetInterfaces(interfaces []*Interface) *File {
	f.Interfaces = interfaces
	return f
}

func (f *File) GetInterfaces() []*Interface {
	return f.Interfaces
}

func (f *File) SetStructs(structs []*Struct) *File {
	f.Structs = structs
	return f
}

func (f *File) GetStructs() []*Struct {
	return f.Structs
}

func (f *File) SetFunctions(functions []*Function) *File {
	f.Functions = functions
	return f
}

func (f *File) GetFunctions() []*Function {
	return f.Functions
}

func (f *File) MustStringVars() string {
	var result []string
	for _, v := range f.Vars {
		result = append(result, v.MustString())
	}
	return strings.Join(result, "\n\n")
}

func (f *File) MustStringTypeAliases() string {
	var result []string
	for _, v := range f.TypeAliases {
		result = append(result, "type "+v.MustString())
	}
	return strings.Join(result, "\n\n")
}

func (f *File) MustStringInterfaces() string {
	var result []string
	for _, v := range f.Interfaces {
		result = append(result, v.MustString())
	}
	return strings.Join(result, "\n\n")
}

func (f *File) MustStringStructs() string {
	var result []string
	for _, v := range f.Structs {
		result = append(result, v.MustString())
	}
	return strings.Join(result, "\n\n")
}

func (f *File) MustStringFunctions() string {
	var result []string
	for _, v := range f.Functions {
		result = append(result, v.MustString())
	}
	return strings.Join(result, "\n\n")
}

func (f *File) SetContents(contents string) *File {
	f.Contents = contents
	return f
}

func (f *File) GetContents() string {
	if f.Contents != "" {
		return f.Contents
	}
	var sections []string

	funcImports := Imports{
		Standard: []Package{},
		App:      []Package{},
		Vendor:   []Package{},
	}
	for _, function := range f.Functions {
		funcImports = MergeImports(funcImports, function.Imports)
	}
	finalImports := MergeImports(f.Imports, funcImports)

	if finalImports.HasImports() {
		sections = append(sections, finalImports.MustString())
	}
	// if f.InitFunction.Body != "" {
	// 	sections = append(sections, f.InitFunction.MustString())
	// }
	// if len(f.Consts) > 0 {
	// 	sections = append(sections, f.MustStringConsts())
	// }
	if len(f.Vars) > 0 {
		sections = append(sections, f.MustStringVars())
	}
	if len(f.TypeAliases) > 0 {
		sections = append(sections, f.MustStringTypeAliases())
	}
	if len(f.Interfaces) > 0 {
		sections = append(sections, f.MustStringInterfaces())
	}
	if len(f.Structs) > 0 {
		sections = append(sections, f.MustStringStructs())
	}
	if len(f.Functions) > 0 {
		sections = append(sections, f.MustStringFunctions())
	}

	result := []string{fmt.Sprintf("package %s\n", f.Package)}

	// Separate each section with an additional line break.
	result = append(result, strings.Join(sections, "\n\n\n"))

	return strings.Join(result, "\n") + "\n"
}

func MakeGoFile(name string) *File {
	return &File{
		Name:      name,
		Package:   "",
		Extension: "go",
	}
}

func MakeFile(name, ext string) *File {
	return &File{
		Name:      name,
		Extension: ext,
	}
}
