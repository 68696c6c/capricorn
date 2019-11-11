package cmd

import (
	"os"

	"github.com/68696c6c/capricorn/generator/models"
	"github.com/68696c6c/capricorn/generator/ops"
	"github.com/68696c6c/capricorn/generator/src"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	Capricorn.AddCommand(newProject)
}

var newProject = &cobra.Command{
	Use:   "new spec",
	Short: "Creates a new Goat project from a YAML spec.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logger := logrus.New()
		logger.SetLevel(logrus.DebugLevel)

		specFile := args[0]
		spec, err := models.NewProject(specFile)
		handleError(err)

		logger.Infof("creating project %s from config %s", spec.Config.Module, specFile)

		// SRC
		err = src.CreateProject(spec)
		handleError(err)

		err = src.CreateModels(spec, logger)
		handleError(err)

		err = src.CreateRepos(spec, logger)
		handleError(err)

		err = src.CreateApp(spec, logger)
		handleError(err)

		err = src.CreateCMD(spec)
		handleError(err)

		err = src.CreateHTTP(&spec, logger)
		handleError(err)

		err = src.FMT(spec.Paths.Root)
		handleError(err)

		// OPS
		err = ops.CreateDocker(spec, logger)
		handleError(err)

		err = ops.CreateMakefile(spec, logger)
		handleError(err)

		err = ops.InitModule(spec.Paths.Root)
		handleError(err)

		logger.Infof("project spec: %s", spec.String())
		os.Exit(0)
	},
}
