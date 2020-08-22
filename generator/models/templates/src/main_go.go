package src

import (
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/models/utils"
)

const templateMainFunc = `
	if err := cmd.Root.Execute(); err != nil {
		println(err)
		os.Exit(1)
	}`

func NewMainGo(modulePackages module.Packages, rootPath, rootPackage string) golang.File {
	fileData, pathData := utils.MakeGoFileData(rootPath, "main")
	return golang.File{
		Name:    fileData,
		Path:    pathData,
		Package: utils.MakePackageData(rootPackage, "main"),
		Imports: golang.Imports{
			Standard: []string{"os"},
			App:      []string{modulePackages.CMD.GetImport()},
		},
		Functions: []golang.Function{
			{
				Name: "main",
				Body: templateMainFunc,
			},
		},
	}
}
