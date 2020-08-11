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
	Use:   "{{.Module.Kebob}}",
	Short: "Root command for {{.Config.Name}}",
}

func init() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.SetDefault("author", "{{.Config.Author.Name}} <{{.Config.Author.Email}}>")
	viper.SetDefault("license", "{{.Config.License}}")
}`

const serverTemplate = `
package cmd

import (
	"{{.Imports.App}}"
	"{{.Imports.HTTP}}"

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

const migrateTemplate = `
package cmd

import (
	"strings"

	"github.com/68696c6c/goat"
	"github.com/68696c6c/goose"
	"github.com/spf13/cobra"
)

func init() {
	Root.AddCommand(migrateCommand)
	dryRun = migrateCommand.Flags().BoolP("dry-run", "d", false, "Only report what would have been done")
}

var dryRun *bool

var migrateCommand = &cobra.Command{
	Use:   "migrate [action] [--dry-run]",
	Short: "Runs the SQL migrations.",
	Long: "Runs the SQL migrations.  Valid actions are: 'up' (default), 'drop', 'reset', and 'install'.",
	Run: func(cmd *cobra.Command, args []string) {
		connection, err := goat.GetMigrationDB()
		if err != nil {
			goat.ExitError(err)
		}
		schema, err := goat.GetSchema(connection)

		// Inform Goose of the current environment.
		// @TODO yuck.
		if err := goat.ErrorIfProd(); err != nil {
			goose.SetEnvProduction(true)
		} else {
			goose.SetEnvProduction(false)
		}

		// Only allow "up" and "install" operations in production.
		allowed := []string{goose.MigrateOperationInstall, goose.MigrateOperationUp}
		goose.SetProductionOperations(allowed)

		// Perform the migration operation.
		migrated, dropped, err := goose.HandleMigrate(schema, args, dryRun)

		dmsg := strings.Join(dropped, "\n")
		println("dropped tables: \n" + dmsg)

		mmsg := strings.Join(migrated, "\n")
		println("migrated tables: \n" + mmsg)

		if err != nil {
			goat.ExitError(err)
		} else {
			goat.ExitSuccess()
		}
	},
}`

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

	// @TODO Create migrate command.
	err = utils.GenerateFile(spec.Paths.CMD, "migrate.go", migrateTemplate, spec)
	if err != nil {
		return errors.Wrap(err, "failed to create migrate command")
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
