package cmd

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

const rootInitFunctionTemplate = `
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.SetDefault("author", "{{ .AuthorName }} <{{ .AuthorEmail }}>")
	viper.SetDefault("license", "{{ .License }}")
`

const rootVarTemplate = `&cobra.Command{
	Use:   "{{ .Use }}",
	Short: "{{ .Short }}",
}`

type Root struct {
	FileData     data.FileData    `yaml:"file_data,omitempty"`
	PathData     data.PathData    `yaml:"path_data,omitempty"`
	PackageData  data.PackageData `yaml:"package_data,omitempty"`
	Name         string           `yaml:"name,omitempty"`
	Use          string           `yaml:"use,omitempty"`
	Short        string           `yaml:"short,omitempty"`
	AuthorName   string           `yaml:"author_name,omitempty"`
	AuthorEmail  string           `yaml:"author_email,omitempty"`
	License      string           `yaml:"license,omitempty"`
	imports      golang.Imports
	initFunction golang.Function
	vars         []golang.Var
	built        bool
}

func NewRoot(m module.Module) *Root {
	name := "root"
	appName := m.Name
	pkgData := m.Packages.CMD
	author := m.GetAuthor()
	fileData, pathData := data.MakeGoFileData(pkgData.GetImport(), name)
	return &Root{
		FileData:    fileData,
		PathData:    pathData,
		PackageData: pkgData,
		Name:        name,
		Use:         appName.Kebob,
		Short:       "Root command for " + appName.Space,
		AuthorName:  author.Name,
		AuthorEmail: author.Email,
		License:     m.GetLicense(),
		imports: golang.Imports{
			Standard: []string{"strings"},
			App:      nil,
			Vendor:   []string{data.ImportCobra, data.ImportViper},
		},
	}
}

func (m *Root) build() {
	if m.built {
		return
	}
	m.initFunction = golang.Function{
		Name: "init",
		Body: m.MustParseInit(),
	}
	m.vars = []golang.Var{
		{
			Name:  "Root",
			Value: m.MustParseVar(),
		},
	}
	m.built = true
}

func (m *Root) MustParseInit() string {
	result, err := utils.ParseTemplateToString("tmp_template_cmd_root_init", rootInitFunctionTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}

func (m *Root) MustParseVar() string {
	result, err := utils.ParseTemplateToString("tmp_template_cmd_root_var", rootVarTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}

func (m *Root) MustGetFile() golang.File {
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
		Vars:         m.vars,
		Interfaces:   []golang.Interface{},
		Structs:      []golang.Struct{},
		Functions:    []golang.Function{},
	}
}
