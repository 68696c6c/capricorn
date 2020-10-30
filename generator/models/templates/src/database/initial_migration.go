package database

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/models/templates/src/database/methods"
)

type MigrationMeta struct {
	PackageData data.PackageData
	AppImports  []string
	ModelRefs   []string
}

type InitialMigration struct {
	fileData    data.FileData
	pathData    data.PathData
	packageData data.PackageData
	appImports  []string
	modelRefs   []string

	imports      golang.Imports
	initFunction golang.Function
	functions    []golang.Function

	built bool
}

func NewInitialMigration(meta MigrationMeta) InitialMigration {
	fileName := makeGooseMigrationName("initial_migration")
	pkgData := meta.PackageData
	fileData, pathData := data.MakeGoFileData(pkgData.GetImport(), fileName)
	return InitialMigration{
		fileData:    fileData,
		pathData:    pathData,
		packageData: meta.PackageData,
		appImports:  meta.AppImports,
		modelRefs:   meta.ModelRefs,
	}
}

func (m *InitialMigration) MustGetFile() golang.File {
	if !m.built {
		m.build()
	}
	return golang.File{
		Name:         m.fileData,
		Path:         m.pathData,
		Package:      m.packageData,
		Imports:      m.imports,
		InitFunction: m.initFunction,
		Consts:       []golang.Const{},
		Vars:         []golang.Var{},
		Interfaces:   []golang.Interface{},
		TypeAliases:  []golang.Value{},
		Structs:      []golang.Struct{},
		Functions:    m.functions,
	}
}

func (m *InitialMigration) build() {
	if m.built {
		return
	}

	var imports golang.Imports
	var functions []golang.Function

	up := methods.NewUp(m.appImports, m.modelRefs)
	functions = append(functions, up.MustGetFunction())
	imports = golang.MergeImports(imports, up.GetImports())

	down := methods.NewDown(m.appImports, m.modelRefs)
	functions = append(functions, down.MustGetFunction())
	imports = golang.MergeImports(imports, down.GetImports())

	initFunction := methods.NewInitFunction(up.Name, down.Name)
	m.initFunction = initFunction.MustGetFunction()
	imports = golang.MergeImports(imports, initFunction.GetImports())

	m.imports = imports
	m.functions = functions
	m.built = true
}
