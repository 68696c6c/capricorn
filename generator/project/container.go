package project

import (
	"fmt"
	"github.com/68696c6c/capricorn_rnd/generator/golang"
	"github.com/68696c6c/capricorn_rnd/generator/utils"
	"strings"
)

type ContainerTemplateData struct {
	TypeName             string
	VarName              string
	ReceiverName         string
	DbArgName            string
	LoggerArgName        string
	DbFieldName          string
	LoggerFieldName      string
	ErrorsFieldName      string
	InitializerName      string
	ResourceDeclarations string
}

const containerInit = `
 	if {{ .VarName }} != ({{ .TypeName }}{}) {
 		return {{ .VarName }}, nil
 	}

 	{{ .VarName }} = {{ .TypeName }}{
 		{{ .DbFieldName }}: {{ .DbArgName }},
 		{{ .LoggerFieldName }}: {{ .LoggerArgName }},
 		{{ .ErrorsFieldName }}: goat.NewErrorHandler({{ .LoggerArgName }}),
 		{{ .ResourceDeclarations }}
 	}

 	return {{ .VarName }}, nil`

func makeContainerInit(data ContainerTemplateData) *golang.Function {
	body, err := utils.ParseTemplateToString("tmp_template_containerInit", containerInit, data)
	if err != nil {
		panic(err)
	}
	return &golang.Function{
		Name: "Delete",
		Imports: golang.Imports{
			Vendor: []golang.Package{PkgGoat, PkgGorm, PkgLogrus},
		},
		Arguments: []golang.Value{
			{
				Name: data.DbArgName,
				Type: "*gorm.DB",
			},
			{
				Name: data.LoggerArgName,
				Type: "*logrus.Logger",
			},
		},
		Body: body,
	}
}

func MakeContainer(resourceFields []*golang.Field, containerFieldMap map[string]string, reposPkg golang.Package) (*golang.File, string, string) {
	var references []string
	for fieldName, constructorReference := range containerFieldMap {
		references = append(references, fmt.Sprintf("%s: %s,", fieldName, constructorReference))
	}
	data := ContainerTemplateData{
		TypeName:             "App",
		VarName:              "container",
		ReceiverName:         "a",
		DbArgName:            "db",
		LoggerArgName:        "logger",
		DbFieldName:          "DB",
		LoggerFieldName:      "Logger",
		ErrorsFieldName:      "Errors",
		InitializerName:      "GetApp",
		ResourceDeclarations: strings.Join(references, "\n"),
	}
	return golang.MakeGoFile("app").SetImports(golang.Imports{
		App: []golang.Package{reposPkg},
	}).SetVars([]*golang.Var{
		{
			Name: "container",
			Type: data.TypeName,
		},
	}).SetStructs([]*golang.Struct{
		{
			Name: data.TypeName,
			Fields: append([]*golang.Field{
				{
					Name: data.DbFieldName,
					Type: "*gorm.DB",
				},
				{
					Name: data.LoggerFieldName,
					Type: "*logrus.Logger",
				},
				{
					Name: data.ErrorsFieldName,
					Type: "goat.ErrorHandler",
				},
			}, resourceFields...),
		},
	}).SetFunctions([]*golang.Function{
		makeContainerInit(data),
	}), data.TypeName, data.InitializerName
}
