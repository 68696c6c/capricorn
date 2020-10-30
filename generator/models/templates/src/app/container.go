package app

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/models/templates/src/utils"
)

type Meta struct {
	Name         data.Name
	ReceiverName string
	FileName     string
	PackageData  data.PackageData
}

type ContainerFieldMeta struct {
}

type Container struct {
	FileData    data.FileData
	PathData    data.PathData
	PackageData data.PackageData
	TypeData    data.TypeData
	Name        data.Name
	Receiver    golang.Value

	Imports    golang.Imports
	Vars       []golang.Var
	Interfaces []golang.Interface
	Structs    []golang.Struct
	Functions  []golang.Function

	singletonName string
	fields        []utils.ContainerFieldMeta
	dbField       golang.Field
	errorsField   golang.Field
	loggerField   golang.Field
	built         bool
}

func NewContainer(meta Meta, fields []utils.ContainerFieldMeta) Container {
	pkgData := meta.PackageData
	fileData, pathData := data.MakeGoFileData(pkgData.GetImport(), meta.FileName)
	return Container{
		FileData:    fileData,
		PathData:    pathData,
		PackageData: meta.PackageData,
		TypeData:    data.MakeTypeDataService(pkgData.Reference, meta.Name.Exported, meta.ReceiverName, false),
		Name:        meta.Name,
		Receiver: golang.Value{
			Name: meta.ReceiverName,
			Type: meta.Name.Exported,
		},
		singletonName: "container",
		fields:        fields,
		dbField: golang.Field{
			Name: "DB",
			Type: "*gorm.DB",
		},
		loggerField: golang.Field{
			Name: "Logger",
			Type: "*logrus.Logger",
		},
		errorsField: golang.Field{
			Name: "Errors",
			Type: "goat.ErrorHandler",
		},
	}
}

func (m *Container) MustGetFile() golang.File {
	if !m.built {
		m.build()
	}
	return golang.File{
		Name:         m.FileData,
		Path:         m.PathData,
		Package:      m.PackageData,
		Imports:      m.Imports,
		InitFunction: golang.Function{},
		Consts:       []golang.Const{},
		Vars:         m.Vars,
		Interfaces:   m.Interfaces,
		Structs:      m.Structs,
		Functions:    m.Functions,
	}
}

func (m *Container) build() {
	if m.built {
		return
	}

	imports := golang.Imports{
		Standard: []string{},
		App:      []string{},
		Vendor:   []string{data.ImportGoat, data.ImportGorm, data.ImportLogrus},
	}
	var functions []golang.Function

	initializer := NewInitializer(*m)
	functions = append(functions, initializer.MustGetFunction())
	imports = golang.MergeImports(imports, initializer.GetImports())

	var fields []golang.Field
	for _, f := range m.fields {
		fields = append(fields, f.Field)
	}

	m.Imports = imports
	m.Vars = []golang.Var{
		{
			Name: m.singletonName,
			Type: m.TypeData.Name,
		},
	}
	m.Functions = functions
	m.Structs = []golang.Struct{
		{
			Name:   m.TypeData.Name,
			Fields: fields,
		},
	}
	m.built = true
}
