package src

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	repos "github.com/68696c6c/capricorn/generator/models/templates/src/repo_methods"
)

func makeRepo(meta serviceMeta) golang.File {
	fileData, pathData := data.MakeGoFileData(meta.packageData.GetImport(), meta.fileName)
	result := golang.File{
		Name:    fileData,
		Path:    pathData,
		Package: meta.packageData,
	}

	plural := meta.resource.Inflection.Plural
	single := meta.resource.Inflection.Single

	// Default methods.
	getBaseQueryFunc := makeRepoMethod(repos.BaseQuery{
		Receiver: meta.receiverName,
		Plural:   plural,
		Single:   single,
	})
	result.Functions = append(result.Functions, getBaseQueryFunc)
	result.Imports = mergeImports(result.Imports, getBaseQueryFunc.Imports)

	getFilteredQueryFunc := makeRepoMethod(repos.BaseFilteredQuery{
		Receiver: meta.receiverName,
		Plural:   plural,
		Single:   single,
	})
	result.Functions = append(result.Functions, getFilteredQueryFunc)
	result.Imports = mergeImports(result.Imports, getFilteredQueryFunc.Imports)

	getPageQueryFunc := makeRepoMethod(repos.BasePaginatedQuery{
		Receiver: meta.receiverName,
		Plural:   plural,
		Single:   single,
	})
	result.Functions = append(result.Functions, getPageQueryFunc)
	result.Imports = mergeImports(result.Imports, getPageQueryFunc.Imports)

	// CRUD methods.
	for _, a := range meta.resource.Repo.Actions {
		switch a {

		case module.ResourceActionList:
			m := makeRepoMethod(repos.Filter{
				Receiver: meta.receiverName,
				Plural:   plural,
				Single:   single,
			})
			result.Functions = append(result.Functions, m)
			result.Imports = mergeImports(result.Imports, m.Imports)

		case module.ResourceActionView:
			m := makeRepoMethod(repos.GetByID{
				Receiver: meta.receiverName,
				Plural:   plural,
				Single:   single,
			})
			result.Functions = append(result.Functions, m)
			result.Imports = mergeImports(result.Imports, m.Imports)

		case module.ResourceActionCreate:
			fallthrough
		case module.ResourceActionUpdate:
			m := makeRepoMethod(repos.Save{
				Receiver: meta.receiverName,
				Plural:   plural,
				Single:   single,
			})
			result.Functions = append(result.Functions, m)
			result.Imports = mergeImports(result.Imports, m.Imports)

		case module.ResourceActionDelete:
			m := makeRepoMethod(repos.Delete{
				Receiver: meta.receiverName,
				Plural:   plural,
				Single:   single,
			})
			result.Functions = append(result.Functions, m)
			result.Imports = mergeImports(result.Imports, m.Imports)

		}
	}

	return result
}

func makeRepoMethod(t repos.Method) golang.Function {
	return golang.Function{
		Name:         t.GetName(),
		Imports:      t.GetImports(),
		Arguments:    t.GetArgs(),
		ReturnValues: t.GetReturns(),
		Receiver: golang.Value{
			Name: "r",
			Type: "",
		},
		Body: t.MustParse(),
	}
}
