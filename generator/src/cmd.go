package src

import (
	"github.com/68696c6c/capricorn/generator/models"
	"github.com/68696c6c/capricorn/generator/utils"
	"github.com/pkg/errors"
)

const rootTemplate = `
package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Root = &cobra.Command{
	Use:   "{{ .Module.Kebob }}",
	Short: "Root command for {{ .Config.Name }}",
}

func init() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.SetDefault("author", "{{ .Config.Author.Name }} <{{ .Config.Author.Email }}>")
	viper.SetDefault("license", "{{ .Config.License }}")
}`

const serverTemplate = `
package cmd

import (
	"{{ .Imports.App }}"
	"{{ .Imports.HTTP }}"

	"github.com/68696c6c/goat"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	Root.AddCommand(serverCommand)
}

var serverCommand = &cobra.Command{
	Use:   "server",
	Short: "Runs the API web server.",
	Long:  "Runs the API web server.",
	Run: func(cmd *cobra.Command, args []string) {
		goat.Init()

		logger := goat.GetLogger()

		db, err := goat.GetMainDB()
		if err != nil {
			goat.ExitError(errors.Wrap(err, "failed to initialize database connection"))
		}

		services, err := app.GetApp(db, logger)
		if err != nil {
			goat.ExitError(errors.Wrap(err, "failed to initialize service container"))
		}

		http.InitRouter(services)
	},
}`

const migrateTemplate = `package cmd

import (
	_ "{{ .Imports.Migrations }}"

	"github.com/68696c6c/goat"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/pressly/goose"
	"github.com/spf13/cobra"
)

func init() {
	Root.AddCommand(migrateCommand)
}

var migrateCommand = &cobra.Command{
	Use:   "migrate",
	Short: "goose migrations (go run main.go migrate up)",
	Args:  cobra.MinimumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
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
	},
{{ $tick := "` + "`" + `" }}
	Example: {{ $tick }}
Usage: app migrate [OPTIONS] COMMAND

Drivers:
postgres
mysql
sqlite3
redshift

Commands:
up                   Migrate the DB to the most recent version available
up-to VERSION        Migrate the DB to a specific VERSION
down                 Roll back the version by 1
down-to VERSION      Roll back to a specific VERSION
redo                 Re-run the latest migration
status               Dump the migration status for the current DB
version              Print the current version of the database
create NAME [sql|go] Creates new migration file with the current timestamp

Examples:
app migrate status
app migrate create init sql
app migrate create add_some_column sql
app migrate create fetch_user_data go
app migrate up

app migrate status{{ $tick }},
}
`

const seedTemplate = `package cmd

import (
	"{{ .Imports.Seeders }}"

	"github.com/68696c6c/goat"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	Root.AddCommand(seedCommand)
}

var seedCommand = &cobra.Command{
	Use:   "seed",
	Short: "Seed the main database with some starting data.",
	Run: func(_ *cobra.Command, args []string) {
		goat.Init()

		logger := goat.GetLogger()
		logger.Info("seeding to database...")

		db, err := goat.GetMainDB()
		if err != nil {
			goat.ExitError(errors.Wrap(err, "error initializing seed connection"))
		}

		if err := seeders.Initial(db); err != nil {
			goat.ExitError(err)
		}

		goat.ExitSuccess()
	},
}
`

const genericTemplate = `
package cmd

import (
	"{{ .AppImport }}"

	"github.com/68696c6c/goat"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	Root.AddCommand({{ .VarName }})
}

var {{ .VarName }} = &cobra.Command{
	Use:   "{{ .Use }}",
	Short: "todo",
	Long:  "todo",
	Run: func(cmd *cobra.Command, args []string) {
		goat.Init()

		logger := goat.GetLogger()

		db, err := goat.GetMainDB()
		if err != nil {
			goat.ExitError(errors.Wrap(err, "failed to initialize database connection"))
		}

		_, err = app.GetApp(db, logger)
		if err != nil {
			goat.ExitError(errors.Wrap(err, "failed to initialize service container"))
		}

		goat.ExitError(errors.New("todo"))
	},
}`

func CreateCMD(spec models.Project) error {
	err := utils.CreateDir(spec.Paths.CMD)
	if err != nil {
		return errors.Wrapf(err, "failed to create cmd directory '%s'", spec.Paths.CMD)
	}

	// Create root command.
	err = utils.GenerateFile(spec.Paths.CMD, "root.go", rootTemplate, spec)
	if err != nil {
		return errors.Wrap(err, "failed to create root command")
	}

	// Create server command.
	err = utils.GenerateFile(spec.Paths.CMD, "server.go", serverTemplate, spec)
	if err != nil {
		return errors.Wrap(err, "failed to create server command")
	}

	// Create migrate command.
	err = utils.GenerateFile(spec.Paths.CMD, "migrate.go", migrateTemplate, spec)
	if err != nil {
		return errors.Wrap(err, "failed to create migrate command")
	}

	// Create seed command.
	err = utils.GenerateFile(spec.Paths.CMD, "seed.go", seedTemplate, spec)
	if err != nil {
		return errors.Wrap(err, "failed to create seed command")
	}

	// @TODO Create make:migration command.

	for _, c := range spec.Commands {
		err = utils.GenerateFile(spec.Paths.CMD, c.FileName, genericTemplate, c)
		if err != nil {
			return errors.Wrapf(err, "failed to create command '%s'", c.Name)
		}
	}

	return nil
}
