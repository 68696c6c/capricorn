package cmd

import (
	"os"

	"github.com/68696c6c/capricorn/generator/models"
	"github.com/68696c6c/capricorn/generator/src"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	Capricorn.AddCommand(newProject)
}

var newProject = &cobra.Command{
	Use:   "new spec [specPath] [projectPath]",
	Short: "Creates a new Goat project from a YAML spec.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		logger := logrus.New()
		logger.SetLevel(logrus.DebugLevel)

		specFile := args[0]
		projectPath := args[1]
		spec, err := models.NewProject(specFile, projectPath)
		handleError(err)

		logger.Infof("creating project %s from config %s", spec.Spec.Module, specFile)

		// SRC
		err = src.CreateProject(spec)
		handleError(err)

		// Domains
		err = src.CreateDomains(&spec, logger)
		handleError(err)

		err = src.CreateDatabase(&spec, logger)
		handleError(err)

		// err = src.CreateModels(spec, logger)
		// handleError(err)
		//
		// err = src.CreateRepos(spec, logger)
		// handleError(err)

		err = src.CreateApp(spec, logger)
		handleError(err)

		err = src.CreateCMD(spec)
		handleError(err)

		err = src.CreateHTTP(spec, logger)
		handleError(err)

		err = src.FMT(spec.Paths.Root)
		handleError(err)

		// OPS
		err = src.CreateDocker(spec, logger)
		handleError(err)

		err = src.CreateMakefile(spec, logger)
		handleError(err)

		err = src.InitModule(spec)
		handleError(err)

		// logger.Infof("project spec: %s", spec.String())
		os.Exit(0)
	},
}
