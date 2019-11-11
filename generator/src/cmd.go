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

	// @TODO Create make:migration command.

	return nil
}
