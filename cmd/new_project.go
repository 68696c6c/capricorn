package cmd

import (
	"os"

	"github.com/68696c6c/capricorn/generator"

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
		spec, err := generator.NewSpecFromFilePath(specFile)
		handleError(err)

		logger.Infof("creating project %s from config %s", spec.ModuleName, specFile)

		err = generator.CreateProject(spec)
		handleError(err)

		err = generator.CreateApp(spec)
		handleError(err)

		err = generator.CreateCMD(spec)
		handleError(err)

		err = generator.CreateModels(spec, logger)
		handleError(err)

		err = generator.CreateRepos(spec, logger)
		handleError(err)

		err = generator.CreateHTTP(spec, logger)
		handleError(err)

		err = generator.FMT(spec.Paths.Root)
		handleError(err)

		err = generator.InitModule(spec.Paths.Root)
		handleError(err)

		logger.Debugf("project spec: %s", spec.String())
		os.Exit(0)
	},
}
