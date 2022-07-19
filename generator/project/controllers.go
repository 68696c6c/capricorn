package project

import (
	"fmt"

	"github.com/68696c6c/capricorn_rnd/generator/golang"
	"github.com/68696c6c/capricorn_rnd/generator/spec"
	"github.com/68696c6c/capricorn_rnd/generator/utils"
)

type ControllerTemplateData struct {
	StructName           string
	ConstructorName      string
	ReceiverName         string
	RepoFieldName        string
	ErrorsFieldName      string
	ModelReference       string
	ModelTypeName        string
	RepoReference        string
	ContextArgName       string
	ResourceReadableName string
	RequestTypeName      string
	ListResponseTypeName string
	ViewResponseTypeName string
}

const controllerConstructor = `
	return {{ .StructName }}{
		{{ .RepoFieldName }}: {{ .RepoFieldName }},
		{{ .ErrorsFieldName }}: {{ .ErrorsFieldName }},
	}`

func makeControllerConstructor(data ControllerTemplateData) *golang.Function {
	body, err := utils.ParseTemplateToString("tmp_template_controllerConstructor", controllerConstructor, data)
	if err != nil {
		panic(err)
	}
	return &golang.Function{
		Name: data.ConstructorName,
		Imports: golang.Imports{
			Vendor: []golang.Package{PkgGoat, PkgGorm},
		},
		Arguments: []golang.Value{
			{
				Name: data.RepoFieldName,
				Type: data.RepoReference,
			},
			{
				Name: data.ErrorsFieldName,
				Type: "goat.ErrorHandler",
			},
		},
		ReturnValues: []golang.Value{
			{
				Type: data.StructName,
			},
		},
		Body: body,
	}
}

const controllerList = `
	q := query.NewQueryBuilder({{ .ContextArgName }})

	result, errs := {{ .ReceiverName }}.{{ .RepoFieldName }}.Filter(q)
	if len(errs) > 0 {
		{{ .ReceiverName }}.{{ .ErrorsFieldName }}.HandleErrorsM({{ .ContextArgName }}, errs, "failed to list {{ .ResourceReadableName }}", goat.RespondServerError)
		return
	}

	goat.RespondData({{ .ContextArgName }}, {{ .ListResponseTypeName }}{
		Data: result,
		Pagination: q.Pagination,
	})`

func makeControllerList(data ControllerTemplateData) *golang.Function {
	body, err := utils.ParseTemplateToString("tmp_template_controllerList", controllerList, data)
	if err != nil {
		panic(err)
	}
	return &golang.Function{
		Name: "List",
		Imports: golang.Imports{
			Vendor: []golang.Package{PkgGoat, PkgQuery, PkgGin},
		},
		Receiver: golang.Value{
			Name: data.ReceiverName,
			Type: data.StructName,
		},
		Arguments: []golang.Value{
			{
				Name: data.ContextArgName,
				Type: "*gin.Context",
			},
		},
		Body: body,
	}
}

const controllerView = `
	i := {{ .ContextArgName }}.Param("id")
	id, err := goat.ParseID(i)
	if err != nil {
		{{ .ReceiverName }}.{{ .ErrorsFieldName }}.HandleErrorM({{ .ContextArgName }}, err, "failed to parse id: "+i, goat.RespondBadRequestError)
		return
	}

	m, errs := {{ .ReceiverName }}.{{ .RepoFieldName }}.GetByID(id)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			{{ .ReceiverName }}.{{ .ErrorsFieldName }}.HandleMessage({{ .ContextArgName }}, "{{ .ResourceReadableName }} does not exist", goat.RespondNotFoundError)
			return
		} else {
			{{ .ReceiverName }}.{{ .ErrorsFieldName }}.HandleErrorsM({{ .ContextArgName }}, errs, "failed to get {{ .ResourceReadableName }}", goat.RespondServerError)
			return
		}
	}

	goat.RespondData({{ .ContextArgName }}, {{ .ViewResponseTypeName }}{m})`

func makeControllerView(data ControllerTemplateData) *golang.Function {
	body, err := utils.ParseTemplateToString("tmp_template_controllerView", controllerView, data)
	if err != nil {
		panic(err)
	}
	return &golang.Function{
		Name: "View",
		Imports: golang.Imports{
			Vendor: []golang.Package{PkgGoat, PkgGin},
		},
		Receiver: golang.Value{
			Name: data.ReceiverName,
			Type: data.StructName,
		},
		Arguments: []golang.Value{
			{
				Name: data.ContextArgName,
				Type: "*gin.Context",
			},
		},
		Body: body,
	}
}

const controllerCreate = `
	req, ok := goat.GetRequest({{ .ContextArgName }}).(*{{ .RequestTypeName }})
	if !ok {
		{{ .ReceiverName }}.{{ .ErrorsFieldName }}.HandleMessage(c, "failed to get request", goat.RespondBadRequestError)
		return
	}

	// @TODO generate model factories.
	// @TODO generate model validators.
	m := req.{{ .ModelTypeName }}
	errs := {{ .ReceiverName }}.{{ .RepoFieldName }}.Save(&m)
	if len(errs) > 0 {
		{{ .ReceiverName }}.{{ .ErrorsFieldName }}.HandleErrorsM({{ .ContextArgName }}, errs, "failed to save {{ .ResourceReadableName }}", goat.RespondServerError)
		return
	}

	goat.RespondCreated({{ .ContextArgName }}, {{ .ViewResponseTypeName }}{m})`

func makeControllerCreate(data ControllerTemplateData) *golang.Function {
	body, err := utils.ParseTemplateToString("tmp_template_controllerCreate", controllerCreate, data)
	if err != nil {
		panic(err)
	}
	return &golang.Function{
		Name: "Create",
		Imports: golang.Imports{
			Vendor: []golang.Package{PkgGoat, PkgGin},
		},
		Receiver: golang.Value{
			Name: data.ReceiverName,
			Type: data.StructName,
		},
		Arguments: []golang.Value{
			{
				Name: data.ContextArgName,
				Type: "*gin.Context",
			},
		},
		Body: body,
	}
}

const controllerUpdate = `
	i := {{ .ContextArgName }}.Param("id")
	id, err := goat.ParseID(i)
	if err != nil {
		{{ .ReceiverName }}.{{ .ErrorsFieldName }}.HandleErrorM({{ .ContextArgName }}, err, "failed to parse id: "+i, goat.RespondBadRequestError)
		return
	}

	// @TODO replace this block with an existence validator and build "not found" handling into the repo.
	_, errs := {{ .ReceiverName }}.{{ .RepoFieldName }}.GetByID(id)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			{{ .ReceiverName }}.{{ .ErrorsFieldName }}.HandleMessage({{ .ContextArgName }}, "{{ .ResourceReadableName }} does not exist", goat.RespondNotFoundError)
			return
		} else {
			{{ .ReceiverName }}.{{ .ErrorsFieldName }}.HandleErrorsM({{ .ContextArgName }}, errs, "failed to get {{ .ResourceReadableName }}", goat.RespondServerError)
			return
		}
	}

	req, ok := goat.GetRequest({{ .ContextArgName }}).(*UpdateRequest)
	if !ok {
		{{ .ReceiverName }}.{{ .ErrorsFieldName }}.HandleMessage({{ .ContextArgName }}, "failed to get request", goat.RespondBadRequestError)
		return
	}

	// @TODO generate model factories.
	// @TODO generate model validators.
	errs = {{ .ReceiverName }}.{{ .RepoFieldName }}.Save(&req.{{ .ModelTypeName }})
	if len(errs) > 0 {
		{{ .ReceiverName }}.{{ .ErrorsFieldName }}.HandleErrorsM({{ .ContextArgName }}, errs, "failed to save {{ .ResourceReadableName }}", goat.RespondServerError)
		return
	}

	goat.RespondCreated({{ .ContextArgName }}, {{ .ViewResponseTypeName }}{req.{{ .ModelTypeName }}})`

func makeControllerUpdate(data ControllerTemplateData) *golang.Function {
	body, err := utils.ParseTemplateToString("tmp_template_controllerUpdate", controllerUpdate, data)
	if err != nil {
		panic(err)
	}
	return &golang.Function{
		Name: "Update",
		Imports: golang.Imports{
			Vendor: []golang.Package{PkgGoat, PkgGin},
		},
		Receiver: golang.Value{
			Name: data.ReceiverName,
			Type: data.StructName,
		},
		Arguments: []golang.Value{
			{
				Name: data.ContextArgName,
				Type: "*gin.Context",
			},
		},
		Body: body,
	}
}

const controllerDelete = `
	i := {{ .ContextArgName }}.Param("id")
	id, err := goat.ParseID(i)
	if err != nil {
		{{ .ReceiverName }}.{{ .ErrorsFieldName }}.HandleErrorM({{ .ContextArgName }}, err, "failed to parse id: "+i, goat.RespondBadRequestError)
		return
	}

	m, errs := {{ .ReceiverName }}.{{ .RepoFieldName }}.GetByID(id)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			{{ .ReceiverName }}.{{ .ErrorsFieldName }}.HandleMessage({{ .ContextArgName }}, "{{ .ResourceReadableName }} does not exist", goat.RespondNotFoundError)
			return
		} else {
			{{ .ReceiverName }}.{{ .ErrorsFieldName }}.HandleErrorsM({{ .ContextArgName }}, errs, "failed to get {{ .ResourceReadableName }}", goat.RespondServerError)
			return
		}
	}

	// @TODO generate model factories.
	// @TODO generate model validators.
	errs = {{ .ReceiverName }}.{{ .RepoFieldName }}.Delete(&m)
	if len(errs) > 0 {
		{{ .ReceiverName }}.{{ .ErrorsFieldName }}.HandleErrorsM({{ .ContextArgName }}, errs, "failed to delete {{ .ResourceReadableName }}", goat.RespondServerError)
		return
	}

	goat.RespondValid({{ .ContextArgName }})`

func makeControllerDelete(data ControllerTemplateData) *golang.Function {
	body, err := utils.ParseTemplateToString("tmp_template_controllerDelete", controllerDelete, data)
	if err != nil {
		panic(err)
	}
	return &golang.Function{
		Name: "Delete",
		Imports: golang.Imports{
			Vendor: []golang.Package{PkgGoat, PkgGin},
		},
		Receiver: golang.Value{
			Name: data.ReceiverName,
			Type: data.StructName,
		},
		Arguments: []golang.Value{
			{
				Name: data.ContextArgName,
				Type: "*gin.Context",
			},
		},
		Body: body,
	}
}

func MakeController(name utils.Inflection, actions []string, modelTypeName, modelReference, repoReference string, modelsPkg, reposPkg golang.Package) (*golang.File, string, string) {
	data := ControllerTemplateData{
		StructName:           fmt.Sprintf("%sController", name.Plural.Exported),
		ConstructorName:      fmt.Sprintf("New%sController", name.Plural.Exported),
		ReceiverName:         "c",
		RepoFieldName:        "repo",
		ErrorsFieldName:      "errors",
		ModelReference:       modelReference,
		ModelTypeName:        modelTypeName,
		RepoReference:        repoReference,
		ContextArgName:       "cx",
		ResourceReadableName: name.Plural.Space,
		RequestTypeName:      fmt.Sprintf("%sRequest", name.Single.Exported),
		ListResponseTypeName: fmt.Sprintf("%sResponse", name.Plural.Exported),
		ViewResponseTypeName: fmt.Sprintf("%sResponse", name.Single.Exported),
	}
	structs := []*golang.Struct{
		{
			Name: data.StructName,
			Fields: []*golang.Field{
				{
					Name: data.RepoFieldName,
					Type: data.RepoReference,
				},
				{
					Name: data.ErrorsFieldName,
					Type: "goat.ErrorHandler",
				},
			},
		},
		{
			Name: data.RequestTypeName,
			Fields: []*golang.Field{
				{
					Name: data.ModelReference,
				},
			},
		},
	}
	functions := []*golang.Function{
		makeControllerConstructor(data),
	}
	var needsViewResponse bool
	for _, action := range actions {
		switch action {
		case spec.ActionList:
			functions = append(functions, makeControllerList(data))
			structs = append(structs, &golang.Struct{
				Name: data.ListResponseTypeName,
				Fields: []*golang.Field{
					{
						Name: "Data",
						Type: "[]*" + data.ModelReference,
						Tags: []*golang.Tag{
							{
								Key:    "json",
								Values: []string{"data"},
							},
						},
					},
					{
						Type: "query.Pagination",
						Tags: []*golang.Tag{
							{
								Key:    "json",
								Values: []string{"pagination"},
							},
						},
					},
				},
			})
			break
		case spec.ActionView:
			functions = append(functions, makeControllerView(data))
			needsViewResponse = true
			break
		case spec.ActionCreate:
			functions = append(functions, makeControllerCreate(data))
			needsViewResponse = true
			break
		case spec.ActionUpdate:
			functions = append(functions, makeControllerUpdate(data))
			needsViewResponse = true
			break
		case spec.ActionDelete:
			functions = append(functions, makeControllerDelete(data))
			needsViewResponse = true
			break
		}
	}
	if needsViewResponse {
		structs = append(structs, &golang.Struct{
			Name: data.ViewResponseTypeName,
			Fields: []*golang.Field{
				{
					Name: data.ModelReference,
				},
			},
		})
	}
	return golang.MakeFile(name.Plural.Snake).SetImports(golang.Imports{
		App: []golang.Package{modelsPkg, reposPkg},
	}).SetStructs(structs).SetFunctions(functions), data.ConstructorName, data.RequestTypeName
}
