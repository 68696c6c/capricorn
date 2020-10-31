package cmd

import (
	"fmt"
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

const serverInitFunctionTemplate = `
	{{ .RootCommandName }}.AddCommand(&cobra.Command{
		Use:   "{{ .Use }}",
		Short: "{{ .Short }}",
		Run: func(cmd *cobra.Command, args []string) {
			{{ .MustParseRun }}
		}
	})
`

const serverRunBodyTemplate = `
	goat.Init()

	logger := goat.GetLogger()

	db, err := goat.GetMainDB()
	if err != nil {
		goat.ExitError(errors.Wrap(err, "failed to initialize database connection"))
	}

	services, err := {{ .InitializerRef }}(db, logger)
	if err != nil {
		goat.ExitError(errors.Wrap(err, "failed to initialize service container"))
	}

	{{ .RouterRef }}(services)
`

type Server struct {
	FileData    data.FileData    `yaml:"file_data,omitempty"`
	PathData    data.PathData    `yaml:"path_data,omitempty"`
	PackageData data.PackageData `yaml:"package_data,omitempty"`

	RootCommandName string `yaml:"root_command_name,omitempty"`
	Use             string `yaml:"use,omitempty"`
	Short           string `yaml:"short,omitempty"`
	InitializerRef  string `yaml:"initializer_ref,omitempty"`
	RouterRef       string `yaml:"router_ref,omitempty"`

	imports      golang.Imports
	initFunction golang.Function
	built        bool
}

func NewServer(m module.Module, rootCmdName, appInitRef, routerRef string) *Server {
	appName := m.Name
	pkgData := m.Packages.CMD
	commandName := "server"
	fileData, pathData := data.MakeGoFileData(pkgData.GetImport(), commandName)
	return &Server{
		FileData:        fileData,
		PathData:        pathData,
		PackageData:     pkgData,
		RootCommandName: rootCmdName,
		InitializerRef:  appInitRef,
		RouterRef:       routerRef,
		Use:             commandName,
		Short:           fmt.Sprintf("Runs the %s web server", appName.Space),
		imports: golang.Imports{
			Standard: nil,
			App:      []string{m.Packages.App.GetImport(), m.Packages.HTTP.GetImport()},
			Vendor:   []string{data.ImportGoat, data.ImportErrors, data.ImportCobra},
		},
	}
}

func (m *Server) build() {
	if m.built {
		return
	}
	m.initFunction = golang.Function{
		Name: "init",
		Body: m.MustParseInit(),
	}
	m.built = true
}

func (m *Server) MustParseInit() string {
	result, err := utils.ParseTemplateToString("tmp_template_cmd_server_init", serverInitFunctionTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}

func (m *Server) MustParseRun() string {
	result, err := utils.ParseTemplateToString("tmp_template_cmd_server_run", serverRunBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}

func (m *Server) MustGetFile() golang.File {
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
