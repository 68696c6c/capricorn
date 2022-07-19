package project

import (
	"fmt"
	"github.com/68696c6c/capricorn_rnd/generator/golang"
	"github.com/68696c6c/capricorn_rnd/generator/spec"
	"github.com/68696c6c/capricorn_rnd/generator/utils"
	"strings"
)

type HttpTemplateData struct {
	InitRouterName   string
	ContainerArgName string
	Routes           string
	appPackage       golang.Package
}

const httpRoutes = `
	router := goat.GetRouter()
	engine := router.GetEngine()

	engine.GET("/health", Health)
	engine.GET("/version", Version)
	api := engine.Group("/api")
	api.Use()
	{{ .Routes }}

	err := router.Run()
	if err != nil {
		goat.ExitError(errors.Wrap(err, "error starting server"))
	}`

type ControllerMeta struct {
	ConstructorName     string
	ResourceRequestName string
	RepoName            string
	Actions             []string
}

func MakeRoutes(controllerActionsMap map[string]ControllerMeta, containerReference string, appPkg golang.Package) (*golang.File, string) {
	containerArgName := "services"
	var routes []string
	for resourceName, controllerMeta := range controllerActionsMap {
		routesName := fmt.Sprintf("%sRoutes", resourceName)
		controllerName := fmt.Sprintf("%sController", resourceName)
		routes = append(routes, "{")
		routes = append(routes, fmt.Sprintf("%s := controllers.%s(%s.%s, %s.Errors)", controllerName, controllerMeta.ConstructorName, containerArgName, controllerMeta.RepoName, containerArgName))
		routes = append(routes, fmt.Sprintf(`%s := api.Group("/%s")`, routesName, resourceName))
		for _, action := range controllerMeta.Actions {
			switch action {
			case spec.ActionList:
				routes = append(routes, fmt.Sprintf(`%s.GET("", %s.List)`, routesName, controllerName))
				break
			case spec.ActionView:
				routes = append(routes, fmt.Sprintf(`%s.GET("/:id", %s.View)`, routesName, controllerName))
				break
			case spec.ActionCreate:
				routes = append(routes, fmt.Sprintf(`%s.POST("", goat.BindMiddleware(controllers.%s{}), %s.Create)`, routesName, controllerMeta.ResourceRequestName, controllerName))
				break
			case spec.ActionUpdate:
				routes = append(routes, fmt.Sprintf(`%s.PUT("/:id", goat.BindMiddleware(controllers.%s{}), %s.Update)`, routesName, controllerMeta.ResourceRequestName, controllerName))
				break
			case spec.ActionDelete:
				routes = append(routes, fmt.Sprintf(`%s.DELETE("/:id", %s.Delete)`, routesName, controllerName))
				break
			}
		}
		routes = append(routes, "}\n")
	}
	data := HttpTemplateData{
		InitRouterName:   "InitRouter",
		ContainerArgName: containerArgName,
		Routes:           strings.Join(routes, "\n"),
		appPackage:       appPkg,
	}
	return golang.MakeFile("routes").SetFunctions([]*golang.Function{
		{
			Name: data.InitRouterName,
			Arguments: []golang.Value{
				{
					Name: data.ContainerArgName,
					Type: containerReference,
				},
			},
			Body: utils.MustParse("tmp_template_httpRoutes", httpRoutes, data),
		},
	}), data.InitRouterName
}
