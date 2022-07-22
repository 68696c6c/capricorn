package project

import (
	"fmt"
	"strings"

	"github.com/68696c6c/capricorn_rnd/generator/golang"
	"github.com/68696c6c/capricorn_rnd/generator/spec"
	"github.com/68696c6c/capricorn_rnd/generator/utils"
)

type EnumTemplateData struct {
	TypeName               string
	TypeNameReadable       string
	ReceiverName           string
	FromStringFunctionName string
	ValueArgName           string
	Values                 string
}

const enumFromString = `
	values := []{{ .TypeName }}{
{{ .Values }}
	}
	for _, v := range values {
		if string(v) == {{ .ValueArgName }} {
			return {{ .TypeName }}({{ .ValueArgName }}), nil
		}
	}
	return "", errors.Errorf("'%s' is not a valid {{ .TypeNameReadable }}", {{ .ValueArgName }})`

func makeEnumFuncFromString(data EnumTemplateData) *golang.Function {
	body, err := utils.ParseTemplateToString("tmp_template_enumFromStringFuncBody", enumFromString, data)
	if err != nil {
		panic(err)
	}
	return &golang.Function{
		Name: data.FromStringFunctionName,
		Imports: golang.Imports{
			Vendor: []golang.Package{PkgErrors},
		},
		Arguments: []golang.Value{
			{
				Name: data.ValueArgName,
				Type: "string",
			},
		},
		ReturnValues: []golang.Value{
			{
				Type: data.TypeName,
			},
			{
				Type: "error",
			},
		},
		Body: body,
	}
}

const enumScanFuncBody = `
	stringValue := fmt.Sprintf("%v", {{ .ValueArgName }})
	result, err := {{ .FromStringFunctionName }}(stringValue)
	if err != nil {
		return err
	}
	*{{ .ReceiverName }} = result
	return nil`

func makeEnumFuncScan(data EnumTemplateData) *golang.Function {
	body, err := utils.ParseTemplateToString("tmp_template_enumScanFuncBody", enumScanFuncBody, data)
	if err != nil {
		panic(err)
	}
	return &golang.Function{
		Name: "Scan",
		Imports: golang.Imports{
			Standard: []golang.Package{PkgStdFmt},
		},
		Receiver: golang.Value{
			Name: data.ReceiverName,
			Type: "*" + data.TypeName,
		},
		Arguments: []golang.Value{
			{
				Name: data.ValueArgName,
				Type: "any",
			},
		},
		ReturnValues: []golang.Value{
			{
				Type: "error",
			},
		},
		Body: body,
	}
}

const enumValueFuncBody = `return string({{ .ReceiverName }}), nil`

func makeEnumFuncValue(data EnumTemplateData) *golang.Function {
	body, err := utils.ParseTemplateToString("tmp_template_enumValueFuncBody", enumValueFuncBody, data)
	if err != nil {
		panic(err)
	}
	return &golang.Function{
		Name: "Value",
		Imports: golang.Imports{
			Standard: []golang.Package{PkgStdSqlDriver},
		},
		Receiver: golang.Value{
			Name: data.ReceiverName,
			Type: data.TypeName,
		},
		ReturnValues: []golang.Value{
			{
				Type: "driver.Value",
			},
			{
				Type: "error",
			},
		},
		Body: body,
	}
}

const enumStringFuncBody = `return string(t)`

func makeEnumFuncString(data EnumTemplateData) *golang.Function {
	body, err := utils.ParseTemplateToString("tmp_template_enumStringFuncBody", enumStringFuncBody, data)
	if err != nil {
		panic(err)
	}
	return &golang.Function{
		Name: "String",
		Receiver: golang.Value{
			Name: data.ReceiverName,
			Type: data.TypeName,
		},
		ReturnValues: []golang.Value{
			{
				Type: "string",
			},
		},
		Body: body,
	}
}

func MakeEnum(enum spec.Enum) (*golang.File, string) {
	name := utils.MakeInflection(enum.Name).Single

	var values []string
	for _, v := range enum.Values {
		values = append(values, "\t\t"+fmt.Sprintf(`"%s",`, v))
	}

	data := EnumTemplateData{
		TypeName:               name.Exported,
		TypeNameReadable:       name.Space,
		ReceiverName:           "t",
		FromStringFunctionName: fmt.Sprintf("%sFromString", name.Exported),
		ValueArgName:           "value",
		Values:                 strings.Join(values, "\n"),
	}

	return golang.MakeGoFile(name.Snake).SetTypeAliases([]*golang.Value{
		{
			Name: name.Exported,
			Type: "string",
		},
	}).SetFunctions([]*golang.Function{
		makeEnumFuncFromString(data),
		makeEnumFuncScan(data),
		makeEnumFuncValue(data),
		makeEnumFuncString(data),
	}), name.Exported
}
