package project

import (
	"github.com/68696c6c/capricorn_rnd/generator/golang"
	"github.com/68696c6c/capricorn_rnd/generator/spec"
	"github.com/68696c6c/capricorn_rnd/generator/utils"
)

type CommandTemplateData struct {
	RootUse              string
	RootShort            string
	AuthorName           string
	AuthorEmail          string
	License              string
	RootCommandName      string
	InitializerReference string
	RouterReference      string
	appPackage           golang.Package
	httpPackage          golang.Package
}

const commandRootInit = `
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.SetDefault("author", "{{ .AuthorName }} <{{ .AuthorEmail }}>")
	viper.SetDefault("license", "{{ .License }}")`

const commandRootVar = `&cobra.Command{
	Use:   "{{ .RootUse }}",
	Short: "{{ .RootShort }}",
}`

func makeCommandRoot(data CommandTemplateData) *golang.File {
	return golang.MakeGoFile("root").SetImports(golang.Imports{
		Vendor: []golang.Package{PkgCobra, PkgViper},
	}).SetVars([]*golang.Var{
		{
			Name:  data.RootCommandName,
			Value: utils.MustParse("tmp_template_commandRootVar", commandRootVar, data),
		},
	}).SetFunctions([]*golang.Function{
		{
			Name: "init",
			Body: utils.MustParse("tmp_template_commandRootInit", commandRootInit, data),
		},
	})
}

const commandMigrate = `
	{{ .RootCommandName }}.AddCommand(&cobra.Command{
		Use:   "migrate",
		Short: "Runs the database migrations",
		Run: func(cmd *cobra.Command, args []string) {
			goat.Init()

			db, err := goat.GetMainDB()
			if err != nil {
				goat.ExitError(errors.Wrap(err, "error initializing migration connection"))
			}

			if err := goose.SetDialect("mysql"); err != nil {
				goat.ExitError(errors.Wrap(err, "error initializing goose"))
			}

			var arguments []string
			if len(args) > 1 {
				arguments = args[1:]
			}

			if err := goose.Run(args[0], db.DB(), ".", arguments...); err != nil {
				goat.ExitError(err)
			}

			goat.ExitSuccess()
		}
	})`

func makeCommandMigrate(data CommandTemplateData) *golang.File {
	return golang.MakeGoFile("migrate").SetImports(golang.Imports{
		Vendor: []golang.Package{PkgGoat, PkgCobra, PkgViper},
		App:    []golang.Package{data.appPackage, data.httpPackage},
	}).SetFunctions([]*golang.Function{
		{
			Name: "init",
			Body: utils.MustParse("tmp_template_commandMigrate", commandMigrate, data),
		},
	})
}

const commandSeed = `
	{{ .RootCommandName }}.AddCommand(&cobra.Command{
		Use:   "seed",
		Short: "Runs the database seeders",
		Run: func(cmd *cobra.Command, args []string) {
			goat.Init()

			logger := goat.GetLogger()
			logger.Info("seeding the database...")

			db, err := goat.GetMainDB()
			if err != nil {
				goat.ExitError(errors.Wrap(err, "error initializing seed connection"))
			}

			if err := seeders.Initial(db); err != nil {
				goat.ExitError(err)
			}

			goat.ExitSuccess()
		}
	})`

func makeCommandSeed(data CommandTemplateData) *golang.File {
	return golang.MakeGoFile("seed").SetImports(golang.Imports{
		Vendor: []golang.Package{PkgGoat, PkgCobra, PkgViper},
		App:    []golang.Package{data.appPackage, data.httpPackage},
	}).SetFunctions([]*golang.Function{
		{
			Name: "init",
			Body: utils.MustParse("tmp_template_commandSeed", commandSeed, data),
		},
	})
}

const commandServer = `
	{{ .RootCommandName }}.AddCommand(&cobra.Command{
		Use:   "server",
		Short: "Runs the web server",
		Run: func(cmd *cobra.Command, args []string) {
			goat.Init()

			logger := goat.GetLogger()

			db, err := goat.GetMainDB()
			if err != nil {
				goat.ExitError(errors.Wrap(err, "failed to initialize database connection"))
			}

			services, err := {{ .InitializerReference }}(db, logger)
			if err != nil {
				goat.ExitError(errors.Wrap(err, "failed to initialize service container"))
			}

			{{ .RouterReference }}(services)
		},
	})`

func makeCommandServer(data CommandTemplateData) *golang.File {
	return golang.MakeGoFile("server").SetImports(golang.Imports{
		Vendor: []golang.Package{PkgGoat, PkgCobra, PkgViper},
		App:    []golang.Package{data.appPackage, data.httpPackage},
	}).SetFunctions([]*golang.Function{
		{
			Name: "init",
			Body: utils.MustParse("tmp_template_commandServer", commandServer, data),
		},
	})
}

func MakeCommands(projectName, initializerReference, routerReference string, appPkg, httpPkg golang.Package, projectSpec spec.Spec) ([]*golang.File, string) {
	data := CommandTemplateData{
		RootUse:              projectName,
		RootShort:            "Root command for " + projectName,
		AuthorName:           projectSpec.Author.Name,
		AuthorEmail:          projectSpec.Author.Email,
		License:              projectSpec.License,
		RootCommandName:      "Root",
		InitializerReference: initializerReference,
		RouterReference:      routerReference,
		appPackage:           appPkg,
		httpPackage:          httpPkg,
	}
	return []*golang.File{
		makeCommandRoot(data),
		makeCommandMigrate(data),
		makeCommandSeed(data),
		makeCommandServer(data),
	}, data.RootCommandName
}
