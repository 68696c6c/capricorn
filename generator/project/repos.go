package project

import (
	"fmt"

	"github.com/68696c6c/capricorn_rnd/generator/golang"
	"github.com/68696c6c/capricorn_rnd/generator/spec"
	"github.com/68696c6c/capricorn_rnd/generator/utils"
)

type RepoTemplateData struct {
	InterfaceName         string
	StructName            string
	DbFieldName           string
	DbArgName             string
	ReceiverName          string
	ModelReference        string
	BaseQueryFuncName     string
	FilteredQueryFuncName string
	PaginateQueryFuncName string
	ResourceReadableName  string
	QueryArgName          string
}

const repoConstructor = `
	return {{ .StructName }}{
		{{ .DbFieldName }}: {{ .DbArgName }},
	}`

func makeRepoConstructor(data RepoTemplateData) *golang.Function {
	body, err := utils.ParseTemplateToString("tmp_template_repoConstructor", repoConstructor, data)
	if err != nil {
		panic(err)
	}
	return &golang.Function{
		Name: fmt.Sprintf("New%s", data.InterfaceName),
		Imports: golang.Imports{
			Vendor: []golang.Package{PkgGorm},
		},
		Arguments: []golang.Value{
			{
				Name: data.DbArgName,
				Type: "*gorm.DB",
			},
		},
		ReturnValues: []golang.Value{
			{
				Type: data.InterfaceName,
			},
		},
		Body: body,
	}
}

const repoBaseQuery = `return {{ .ReceiverName }}.{{ .DbFieldName }}.Model(&{{ .ModelReference }}{})`

func makeRepoBaseQuery(data RepoTemplateData) *golang.Function {
	body, err := utils.ParseTemplateToString("tmp_template_repoBaseQuery", repoBaseQuery, data)
	if err != nil {
		panic(err)
	}
	return &golang.Function{
		Name: data.BaseQueryFuncName,
		Imports: golang.Imports{
			Vendor: []golang.Package{PkgGoat, PkgErrors, PkgQuery, PkgGorm},
		},
		Receiver: golang.Value{
			Name: data.ReceiverName,
			Type: data.StructName,
		},
		ReturnValues: []golang.Value{
			{
				Type: "*gorm.DB",
			},
		},
		Body: body,
	}
}

const repoFilteredQuery = `
	result, err := {{ .QueryArgName }}.ApplyToGorm({{ .ReceiverName }}.{{ .BaseQueryFuncName }}())
	if err != nil {
		return result, err
	}
	return result, nil`

func makeRepoFilteredQuery(data RepoTemplateData) *golang.Function {
	body, err := utils.ParseTemplateToString("tmp_template_repoFilteredQuery", repoFilteredQuery, data)
	if err != nil {
		panic(err)
	}
	return &golang.Function{
		Name: data.FilteredQueryFuncName,
		Imports: golang.Imports{
			Vendor: []golang.Package{PkgQuery, PkgGorm},
		},
		Receiver: golang.Value{
			Name: data.ReceiverName,
			Type: data.StructName,
		},
		Arguments: []golang.Value{
			{
				Name: data.QueryArgName,
				Type: "*query.Query",
			},
		},
		ReturnValues: []golang.Value{
			{
				Type: "*gorm.DB",
			},
			{
				Type: "error",
			},
		},
		Body: body,
	}
}

const repoPaginateQuery = `
	err := goat.ApplyPaginationToQuery({{ .QueryArgName }}, {{ .ReceiverName }}.{{ .BaseQueryFuncName }}())
	if err != nil {
		return errors.Wrap(err, "failed to set {{ .ResourceReadableName }} query pagination")
	}
	return nil`

func makeRepoPaginatedQuery(data RepoTemplateData) *golang.Function {
	body, err := utils.ParseTemplateToString("tmp_template_repoPaginateQuery", repoPaginateQuery, data)
	if err != nil {
		panic(err)
	}
	return &golang.Function{
		Name: data.PaginateQueryFuncName,
		Imports: golang.Imports{
			Vendor: []golang.Package{PkgGoat, PkgErrors, PkgQuery, PkgGorm},
		},
		Receiver: golang.Value{
			Name: data.ReceiverName,
			Type: data.StructName,
		},
		Arguments: []golang.Value{
			{
				Name: data.QueryArgName,
				Type: "*query.Query",
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

const repoFilter = `
	dataQuery, err := {{ .ReceiverName }}.{{ .FilteredQueryFuncName }}({{ .QueryArgName }})
	if err != nil {
		return result, errors.Wrap(err, "failed to build filter {{ .ResourceReadableName }} query")
	}

	errs := dataQuery.Find(&result).GetErrors()
	if len(errs) > 0 && goat.ErrorsBesidesRecordNotFound(errs) {
		err := goat.ErrorsToError(errs)
		return result, errors.Wrap(err, "failed to execute filter {{ .ResourceReadableName }} query")
	}

	if err := {{ .ReceiverName }}.{{ .PaginateQueryFuncName }}({{ .QueryArgName }}); err != nil {
		return result, err
	}

	return result, nil`

func makeRepoFilter(data RepoTemplateData) *golang.Function {
	body, err := utils.ParseTemplateToString("tmp_template_repoFilter", repoFilter, data)
	if err != nil {
		panic(err)
	}
	return &golang.Function{
		Name: "Filter",
		Imports: golang.Imports{
			Vendor: []golang.Package{PkgGoat, PkgErrors, PkgQuery, PkgGorm},
		},
		Receiver: golang.Value{
			Name: data.ReceiverName,
			Type: data.StructName,
		},
		Arguments: []golang.Value{
			{
				Name: data.QueryArgName,
				Type: "*query.Query",
			},
		},
		ReturnValues: []golang.Value{
			{
				Name: "result",
				Type: "[]*" + data.ModelReference,
			},
			{
				Name: "err",
				Type: "error",
			},
		},
		Body: body,
	}
}

const repoGetById = `
	m := {{ .ModelReference }}{
		Model: goat.Model{
			ID: id,
		},
	}
	errs := {{ .ReceiverName }}.{{ .DbFieldName }}.First(&m).GetErrors()
	if len(errs) > 0 {
		return m, goat.ErrorsToError(errs)
	}
	return m, nil`

func makeRepoGetById(data RepoTemplateData) *golang.Function {
	body, err := utils.ParseTemplateToString("tmp_template_repoGetById", repoGetById, data)
	if err != nil {
		panic(err)
	}
	return &golang.Function{
		Name: "GetById",
		Imports: golang.Imports{
			Vendor: []golang.Package{PkgGoat},
		},
		Receiver: golang.Value{
			Name: data.ReceiverName,
			Type: data.StructName,
		},
		Arguments: []golang.Value{
			{
				Name: "id",
				Type: "goat.ID",
			},
		},
		ReturnValues: []golang.Value{
			{
				Type: data.ModelReference,
			},
			{
				Type: "error",
			},
		},
		Body: body,
	}
}

const repoSave = `
	var errs []error
	if m.Model.ID.Valid() {
		errs = {{ .ReceiverName }}.{{ .DbFieldName }}.Save(m).GetErrors()
	} else {
		errs = {{ .ReceiverName }}.{{ .DbFieldName }}.Create(m).GetErrors()
	}
	if len(errs) > 0 {
		return goat.ErrorsToError(errs)
	}
	return nil`

func makeRepoSave(data RepoTemplateData) *golang.Function {
	body, err := utils.ParseTemplateToString("tmp_template_repoSave", repoSave, data)
	if err != nil {
		panic(err)
	}
	return &golang.Function{
		Name: "Save",
		Imports: golang.Imports{
			Vendor: []golang.Package{PkgGoat},
		},
		Receiver: golang.Value{
			Name: data.ReceiverName,
			Type: data.StructName,
		},
		Arguments: []golang.Value{
			{
				Name: "m",
				Type: "*" + data.ModelReference,
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

const repoDelete = `
	errs :=  {{ .ReceiverName }}.{{ .DbFieldName }}.Delete(m).GetErrors()
	if len(errs) > 0 {
		return goat.ErrorsToError(errs)
	}
	return nil`

func makeRepoDelete(data RepoTemplateData) *golang.Function {
	body, err := utils.ParseTemplateToString("tmp_template_repoDelete", repoDelete, data)
	if err != nil {
		panic(err)
	}
	return &golang.Function{
		Name: "Delete",
		Imports: golang.Imports{
			Vendor: []golang.Package{PkgGoat},
		},
		Receiver: golang.Value{
			Name: data.ReceiverName,
			Type: data.StructName,
		},
		Arguments: []golang.Value{
			{
				Name: "m",
				Type: "*" + data.ModelReference,
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

func MakeRepo(name utils.Inflection, actions []string, modelReference string, modelsPkg golang.Package) (*golang.File, string) {
	data := RepoTemplateData{
		InterfaceName:         fmt.Sprintf("%sRepo", name.Plural.Exported),
		StructName:            fmt.Sprintf("%sRepoGorm", name.Plural.Unexported),
		DbFieldName:           "db",
		DbArgName:             "db",
		QueryArgName:          "q",
		ReceiverName:          "r",
		ModelReference:        modelReference,
		BaseQueryFuncName:     "getBaseQuery",
		FilteredQueryFuncName: "getFilteredQuery",
		PaginateQueryFuncName: "ApplyPaginationToQuery",
		ResourceReadableName:  name.Plural.Space,
	}
	functions := []*golang.Function{
		makeRepoPaginatedQuery(data),
	}
	var saveHandled bool
	for _, action := range actions {
		switch action {
		case spec.ActionList:
			functions = append(functions, makeRepoFilter(data))
			break
		case spec.ActionView:
			functions = append(functions, makeRepoGetById(data))
			break
		case spec.ActionCreate:
			fallthrough
		case spec.ActionUpdate:
			if !saveHandled {
				functions = append(functions, makeRepoSave(data))
				saveHandled = true
			}
			break
		case spec.ActionDelete:
			functions = append(functions, makeRepoDelete(data))
			break
		}
	}
	return golang.MakeFile(name.Plural.Snake).SetImports(golang.Imports{
		App: []golang.Package{modelsPkg},
	}).SetInterfaces([]*golang.Interface{
		{
			Name:      data.InterfaceName,
			Functions: functions,
		},
	}).SetStructs([]*golang.Struct{
		{
			Name: data.StructName,
			Fields: []*golang.Field{
				{
					Name: data.DbFieldName,
					Type: "*gorm.DB",
				},
			},
		},
	}).SetFunctions(append([]*golang.Function{
		makeRepoConstructor(data),
		makeRepoBaseQuery(data),
		makeRepoFilteredQuery(data),
	}, functions...)), data.InterfaceName
}
