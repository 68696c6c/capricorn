package project

import (
	"fmt"
	"path"
	"strings"

	"github.com/68696c6c/girraph"

	"github.com/68696c6c/capricorn_rnd/generator/golang"
	"github.com/68696c6c/capricorn_rnd/generator/spec"
	"github.com/68696c6c/capricorn_rnd/generator/utils"
)

func FromSpec(projectSpec spec.Spec) (girraph.Tree[golang.Package], string) {
	projectDir := golang.MakePackageNode(projectSpec.Name)
	projectDir.GetMeta().SetFiles([]*golang.File{
		makeAppEnv(projectSpec.Ops),
		makeAppTemplateEnv(projectSpec.Ops),
		makeDockerCompose(projectSpec.Ops),
		makeDockerfile(projectSpec.Ops),
		makeMakefile(projectSpec.Ops),
		golang.MakeFile(".gitignore", "").SetContents(gitignore),
	})

	srcDir, srcPath := makeSrc(projectSpec)
	projectDir.SetChildren([]girraph.Tree[golang.Package]{
		makeOps(),
		srcDir,
	})

	return projectDir, srcPath
}

func newProject(baseImport string) girraph.Tree[golang.Package] {
	result := golang.MakePackageNode("")
	result.GetMeta().SetReference("main").SetImport(baseImport)
	return result
}

func makeSrc(projectSpec spec.Spec) (girraph.Tree[golang.Package], string) {
	srcDir := golang.MakePackageNode("src")

	project := newProject(projectSpec.Module)

	app := golang.MakePackageNode("app")

	enumMap := make(map[string]string)
	enums := golang.MakePackageNode("enums")
	for _, e := range projectSpec.Enums {
		enum, enumType := MakeEnum(e)
		enumMap[enum.Name] = fmt.Sprintf("enums.%s", enumType)
		enums.GetMeta().AddFile(enum)
	}

	models := golang.MakePackageNode("models")
	repos := golang.MakePackageNode("repos")
	controllers := golang.MakePackageNode("controllers")
	var containerFields []*golang.Field
	var modelReferences []string
	containerFieldMap := make(map[string]string)
	controllerActionsMap := make(map[string]ControllerMeta)
	for _, r := range projectSpec.Resources {
		recName := utils.MakeInflection(r.Name)
		actions := r.Actions
		if len(actions) == 0 {
			actions = spec.ActionsAll
		}
		model, modelName := MakeModel(recName, r, enumMap, enums.GetMeta())
		modelReference := fmt.Sprintf("models.%s", modelName)
		modelReferences = append(modelReferences, modelReference)
		models.GetMeta().AddFile(model)
		repo, repoName := MakeRepo(recName, actions, modelReference, models.GetMeta())
		repoReference := fmt.Sprintf("repos.%s", repoName)
		repoConstructorReference := fmt.Sprintf("repos.New%s", repoName)
		repos.GetMeta().AddFile(repo)
		containerFields = append(containerFields, &golang.Field{
			Type: repoReference,
		})
		containerFieldMap[repoName] = repoConstructorReference
		controller, controllerConstructorName, requestName := MakeController(recName, actions, modelName, modelReference, repoReference, models.GetMeta(), repos.GetMeta())
		controllers.GetMeta().AddFile(controller)
		controllerActionsMap[repoName] = ControllerMeta{
			ConstructorName:     controllerConstructorName,
			ResourceRequestName: requestName,
			RepoName:            repoName,
			Actions:             actions,
		}
	}

	container, containerTypeName, initFuncName := MakeContainer(containerFields, containerFieldMap, repos.GetMeta())
	containerReference := fmt.Sprintf("app.%s", containerTypeName)
	app.GetMeta().SetFiles([]*golang.File{
		container,
	})
	services := golang.MakePackageNode("services")
	app.SetChildren([]girraph.Tree[golang.Package]{
		controllers,
		enums,
		models,
		repos,
		services,
	})

	database := golang.MakePackageNode("database")
	migrations := golang.MakePackageNode("migrations")
	migrations.GetMeta().SetFiles([]*golang.File{
		MakeInitialMigration(modelReferences, models.GetMeta()),
	})
	seeders := golang.MakePackageNode("seeders")
	seeders.GetMeta().SetFiles([]*golang.File{
		MakeInitialSeeder(),
	})
	database.SetChildren([]girraph.Tree[golang.Package]{
		migrations,
		seeders,
	})

	http := golang.MakePackageNode("http")
	routes, routerConstructorName := MakeRoutes(controllerActionsMap, containerReference, app.GetMeta())
	http.GetMeta().SetFiles([]*golang.File{
		routes,
	})
	routerReference := fmt.Sprintf("http.%s", routerConstructorName)

	parts := strings.Split(projectSpec.Module, "/")
	projectModuleName := parts[len(parts)-1]
	projectName := utils.MakeInflection(projectModuleName).Single.Kebob

	cmdDir := golang.MakePackageNode("cmd")
	initReference := fmt.Sprintf("app.%s", initFuncName)
	commands, rootCommandName := MakeCommands(projectName, initReference, routerReference, app.GetMeta(), http.GetMeta(), projectSpec)
	cmdDir.GetMeta().SetFiles(commands)
	for _, c := range projectSpec.Commands {
		cmdName := strings.Replace(c.Name, ":", "_", -1)
		cmdDir.GetMeta().AddFile(golang.MakeGoFile(cmdName))
	}
	rootCommandReference := fmt.Sprintf("cmd.%s", rootCommandName)
	project.GetMeta().SetFiles([]*golang.File{
		MakeMain(rootCommandReference, cmdDir.GetMeta()),
	})

	project.SetChildren([]girraph.Tree[golang.Package]{
		app,
		cmdDir,
		database,
		http,
	})

	srcDir.SetChildren([]girraph.Tree[golang.Package]{
		project,
	})

	golang.SetImports(projectSpec.Module, project)

	return srcDir, path.Join(projectSpec.Name, "src")
}

func MakeMain(rootCommandReference string, cmdPkg golang.Package) *golang.File {
	body := fmt.Sprintf(`
	if err := %s.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}`, rootCommandReference)
	return golang.MakeGoFile("main").SetFunctions([]*golang.Function{
		{
			Name: "main",
			Imports: golang.Imports{
				Standard: []golang.Package{PkgStdFmt, PkgStdOs},
				App:      []golang.Package{cmdPkg},
			},
			Body: body,
		},
	})
}
