package src

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/models/templates/src/handlers"
)

// If generating a non-DDD app, response names are things like "userResponse" and do not need to be exported.
// For DDD apps, a response would be named something like "Response".
// This is because DDD responses are defined in their domain, but used by the router in the app http package as "users.Response".
func makeController(meta serviceMeta, viewResponseName, listResponseName string) golang.File {
	fileData, pathData := data.MakeGoFileData(meta.packageData.GetImport(), meta.fileName)
	result := golang.File{
		Name:    fileData,
		Path:    pathData,
		Package: meta.packageData,
	}

	plural := meta.resource.Inflection.Plural
	single := meta.resource.Inflection.Single

	// @TODO need to make the repo struct

	for _, a := range meta.resource.Controller.Actions {
		switch a {

		case module.ResourceActionList:
			t := handlers.List{
				Receiver: meta.receiverName,
				Plural:   plural,
				Single:   single,
				Response: listResponseName,
			}
			h := makeHandler("List", t.MustParse(), handlers.GetListImports())
			result.Functions = append(result.Functions, h)
			result.Imports = mergeImports(result.Imports, h.Imports)

		case module.ResourceActionView:
			t := handlers.View{
				Receiver: meta.receiverName,
				Plural:   plural,
				Single:   single,
				Response: viewResponseName,
			}
			h := makeHandler("View", t.MustParse(), handlers.GetViewImports())
			result.Functions = append(result.Functions, h)
			result.Imports = mergeImports(result.Imports, h.Imports)

		case module.ResourceActionCreate:
			t := handlers.Create{
				Receiver: meta.receiverName,
				Plural:   plural,
				Single:   single,
				Response: viewResponseName,
			}
			h := makeHandler("Create", t.MustParse(), handlers.GetCreateImports())
			result.Functions = append(result.Functions, h)
			result.Imports = mergeImports(result.Imports, h.Imports)

		case module.ResourceActionUpdate:
			t := handlers.Update{
				Receiver: meta.receiverName,
				Plural:   plural,
				Single:   single,
				Response: viewResponseName,
			}
			h := makeHandler("Update", t.MustParse(), handlers.GetUpdateImports())
			result.Functions = append(result.Functions, h)
			result.Imports = mergeImports(result.Imports, h.Imports)

		case module.ResourceActionDelete:
			t := handlers.Delete{
				Receiver: meta.receiverName,
				Plural:   plural,
				Single:   single,
			}
			h := makeHandler("Delete", t.MustParse(), handlers.GetDeleteImports())
			result.Functions = append(result.Functions, h)
			result.Imports = mergeImports(result.Imports, h.Imports)

		}
	}

	return result
}

func makeHandler(name, body string, i golang.Imports) golang.Function {
	return golang.Function{
		Name:    name,
		Imports: i,
		Arguments: []golang.Value{
			{
				Name: "cx",
				Type: "*gin.Context",
			},
		},
		ReturnValues: []golang.Value{},
		Receiver: golang.Value{
			Name: "c",
			Type: "",
		},
		Body: body,
	}
}
