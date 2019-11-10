package src

import (
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
	Use:   "{{.ModuleName}}",
	Short: "Root command for {{.Name}}",
}

func init() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.SetDefault("author", "{{.Author.Name}} <{{.Author.Email}}>")
	viper.SetDefault("license", "{{.License}}")
}

`

const serverTemplate = `
package cmd

import (
	"{{.Imports.Packages.App}}"
	"{{.Imports.Packages.HTTP}}"

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
}

`

func CreateCMD(spec utils.Spec) error {
	err := utils.CreateDir(spec.Paths.CMD)
	if err != nil {
		return errors.Wrapf(err, "failed to create cmd directory '%s'", spec.Paths.CMD)
	}

	// Create root command.
	err = utils.GenerateGoFile(spec.Paths.CMD, "root", rootTemplate, spec)
	if err != nil {
		return errors.Wrap(err, "failed to create root command")
	}

	// Create server command.
	err = utils.GenerateGoFile(spec.Paths.CMD, "server", serverTemplate, spec)
	if err != nil {
		return errors.Wrap(err, "failed to create server command")
	}

	// @TODO Create migrate command.

	// @TODO Create make:migration command.

	return nil
}
