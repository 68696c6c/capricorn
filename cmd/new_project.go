package cmd

import (
	"github.com/68696c6c/capricorn_rnd/generator/filesystem"
	"os"
	"path"

	"github.com/spf13/cobra"

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

			p, projectSrcDir := project.NewProjectDirFromSpec(projectSpec, projectSpec.Ops)
			filesystem.SetPaths(projectPath, p)
			filesystem.MustGenerate(p)

			baseImport := projectSpec.Module
			srcPath := path.Join(projectPath, projectSrcDir)

			projectTree := project.ProjectFromSpec(projectSpec)
			golang.SetPaths(srcPath, baseImport, projectTree)
			golang.MustGenerate(projectTree)

			err = utils.InitModule(srcPath, baseImport)
			if err != nil {
				panic(err)
			}

			os.Exit(0)
		},
	})
}
