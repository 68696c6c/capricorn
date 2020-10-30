package app

import (
	"fmt"
	"strings"

	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	utils2 "github.com/68696c6c/capricorn/generator/models/templates/src/utils"
	"github.com/68696c6c/capricorn/generator/utils"
)

var initializerBodyTemplate = `
	if {{ .SingletonName }} != ({{ .TypeName }}{}) {
		return {{ .SingletonName }}, nil
	}
	
	{{ .MustParseDeclarations }}

	{{ .SingletonName }} = {{ .TypeName }}{
		{{ .DBFieldName }}: {{ .DBArgName }},
		{{ .LoggerFieldName }}: {{ .LoggerArgName }},
		{{ .ErrorsFieldName }}: goat.NewErrorHandler({{ .LoggerArgName }}),
		{{ .MustParseFields }}
	}

	return {{ .SingletonName }}, nil
`

type Initializer struct {
	name            string
	receiver        golang.Value
	imports         golang.Imports
	args            []golang.Value
	returns         []golang.Value
	fields          []utils2.ContainerFieldMeta
	repoMap         map[string]string
	built           bool
	SingletonName   string
	TypeName        string
	LoggerArgName   string
	LoggerFieldName string
	DBArgName       string
	DBFieldName     string
	ErrorsFieldName string
}

func NewInitializer(c Container) *Initializer {
	dbArgName := "d"
	loggerArgName := "l"
	return &Initializer{
		name:     "GetApp",
		receiver: golang.Value{},
		imports: golang.Imports{
			Standard: nil,
			App:      nil,
			Vendor:   []string{data.ImportGorm},
		},
		args: []golang.Value{
			{
				Name: dbArgName,
				Type: "*gorm.DB",
			},
			{
				Name: loggerArgName,
				Type: "*logrus.Log",
			},
		},
		returns: []golang.Value{
			{
				Type: c.TypeData.Name,
			},
			{
				Type: "error",
			},
		},
		TypeName:        c.TypeData.Name,
		LoggerArgName:   loggerArgName,
		LoggerFieldName: c.loggerField.Name,
		DBArgName:       dbArgName,
		DBFieldName:     c.dbField.Name,
		ErrorsFieldName: c.errorsField.Name,
		fields:          c.fields,
	}
}

func (m *Initializer) build() {
	var imports []string
	repoMap := map[string]string{}
	for _, f := range m.fields {
		if f.ServiceType == utils2.ServiceTypeRepo {
			varName := f.Name.Unexported
			repoMap[f.DomainKey] = varName
		}
		imports = append(imports, f.PackageImport)
	}
	m.imports = golang.MergeImports(m.imports, golang.Imports{App: imports})
	m.repoMap = repoMap
	m.built = true
}

func handleRepoMapPanic(msg, domainKey string) {
	panic(msg + ", unable to find the name of the variable holding the repo argument for the service constructor.  domain key: " + domainKey)
}

func (m *Initializer) MustParseDeclarations() string {
	if !m.built {
		m.build()
	}
	var result []string
	for _, f := range m.fields {
		if f.ServiceType == utils2.ServiceTypeRepo {
			varName, ok := m.repoMap[f.DomainKey]
			if !ok {
				handleRepoMapPanic("failed to print container repo initialization", f.DomainKey)
			}
			l := fmt.Sprintf("%s := %s.%s(%s)", varName, f.TypeData.Package, f.Constructor.Name, m.DBArgName)
			result = append(result, l)
		}
	}
	return strings.Join(result, "\n")
}

func (m *Initializer) MustParseFields() string {
	if !m.built {
		m.build()
	}
	var result []string
	for _, f := range m.fields {
		// Repos have already been constructed in MustParseDeclarations so we only need to assign the variable to the field.
		if f.ServiceType == utils2.ServiceTypeRepo {
			result = append(result, fmt.Sprintf("%s: %s", f.Field.Name, f.Name.Unexported))
		}
		// Construct and assign the services.
		if f.ServiceType == utils2.ServiceTypeService {
			if f.ServiceType == utils2.ServiceTypeService {
				varName, ok := m.repoMap[f.DomainKey]
				if !ok {
					handleRepoMapPanic("failed to print container service initialization", f.DomainKey)
				}
				l := fmt.Sprintf("%s: %s.%s(%s),", f.Name.Unexported, f.TypeData.Package, f.Constructor.Name, varName)
				result = append(result, l)
			}
		}
	}
	return strings.Join(result, "\n")
}

func (m *Initializer) MustGetFunction() golang.Function {
	if !m.built {
		m.build()
	}
	return golang.Function{
		Name:         m.name,
		Imports:      m.imports,
		Receiver:     m.receiver,
		Arguments:    m.args,
		ReturnValues: m.returns,
		Body:         m.MustParse(),
	}
}

func (m *Initializer) GetImports() golang.Imports {
	return m.imports
}

func (m *Initializer) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_container_initializer", initializerBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
