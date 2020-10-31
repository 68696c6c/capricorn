package cmd

import (
	"fmt"
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

const migrateInitFunctionTemplate = `
	{{ .RootCommandName }}.AddCommand(&cobra.Command{
		Use:   "{{ .Use }}",
		Short: "{{ .Short }}",
		Run: func(cmd *cobra.Command, args []string) {
			{{ .MustParseRun }}
		}
	})
`

const migrateRunBodyTemplate = `
	goat.Init()

	db, err := goat.GetMainDB()
	if err != nil {
		goat.ExitError(errors.Wrap(err, "error initializing migration connection"))
	}

	if err := goose.SetDialect("mysql"); err != nil {
		goat.ExitError(errors.Wrap(err, "error initializing goose"))
	}

	var arguments []string
	if len(args) > 1 {
		arguments = args[1:]
	}

	if err := goose.Run(args[0], db.DB(), ".", arguments...); err != nil {
		goat.ExitError(err)
	}

	goat.ExitSuccess()
`

type Migrate struct {
	FileData    data.FileData    `yaml:"file_data,omitempty"`
	PathData    data.PathData    `yaml:"path_data,omitempty"`
	PackageData data.PackageData `yaml:"package_data,omitempty"`

	RootCommandName string `yaml:"root_command_name,omitempty"`
	Use             string `yaml:"use,omitempty"`
	Short           string `yaml:"short,omitempty"`

	imports      golang.Imports
	initFunction golang.Function
	built        bool
}

func NewMigrate(m module.Module, rootCmdName string) *Migrate {
	appName := m.Name
	pkgData := m.Packages.CMD
	commandName := "migrate"
	fileData, pathData := data.MakeGoFileData(pkgData.GetImport(), commandName)
	return &Migrate{
		FileData:        fileData,
		PathData:        pathData,
		PackageData:     pkgData,
		RootCommandName: rootCmdName,
		Use:             commandName,
		Short:           fmt.Sprintf("Runs the %s migrations", appName.Space),
		imports: golang.Imports{
			Standard: nil,
			App:      []string{fmt.Sprintf("_ \"%s\"", m.Packages.Migrations)},
			Vendor:   []string{data.ImportGoat, data.ImportSqlDriver, data.ImportErrors, data.ImportGoose, data.ImportCobra},
		},
	}
}

func (m *Migrate) build() {
	if m.built {
		return
	}
	m.initFunction = golang.Function{
		Name: "init",
		Body: m.MustParseInit(),
	}
	m.built = true
}

func (m *Migrate) MustParseInit() string {
	result, err := utils.ParseTemplateToString("tmp_template_cmd_migrate_init", migrateInitFunctionTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}

func (m *Migrate) MustParseRun() string {
	result, err := utils.ParseTemplateToString("tmp_template_cmd_migrate_run", migrateRunBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}

func (m *Migrate) MustGetFile() golang.File {
	if !m.built {
		m.build()
	}
	return golang.File{
		Name:         m.FileData,
		Path:         m.PathData,
		Package:      m.PackageData,
		Imports:      m.imports,
		InitFunction: m.initFunction,
		Consts:       []golang.Const{},
		Vars:         []golang.Var{},
		Interfaces:   []golang.Interface{},
		Structs:      []golang.Struct{},
		Functions:    []golang.Function{},
	}
}
