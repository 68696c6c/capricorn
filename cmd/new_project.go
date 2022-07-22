package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"path"

	"github.com/68696c6c/capricorn_rnd/generator/golang"
	"github.com/68696c6c/capricorn_rnd/generator/project"
	"github.com/68696c6c/capricorn_rnd/generator/spec"
	"github.com/68696c6c/capricorn_rnd/generator/utils"
)

func init() {
	Capricorn.AddCommand(&cobra.Command{
		Use:   "new [specPath] [projectPath]",
		Short: "Creates a new Goat project from a YAML spec.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			specFile := args[0]
			projectPath := args[1]

			println("specFile", specFile)
			println("projectPath", projectPath)

			projectSpec, err := spec.NewSpec(specFile)
			handleError(err)

			projectTree, projectSrcDir := project.FromSpec(projectSpec)

			golang.SetPaths(projectPath, projectTree)
			golang.MustGenerate(projectTree)

			srcPath := path.Join(projectPath, projectSrcDir)
			err = utils.InitModule(srcPath, projectSpec.Module)
			if err != nil {
				panic(err)
			}

			os.Exit(0)
		},
	})
}
