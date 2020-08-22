package src

import (
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
)

const templateMainFunc = `
	if err := cmd.Root.Execute(); err != nil {
		println(err)
		os.Exit(1)
	}`

func NewMainGo(modulePackages module.Packages, rootPath, rootPackage string) golang.File {
	fileData, pathData := templates.MakeGoFileData(rootPath, "main")
	return golang.File{
		Name:    fileData,
		Path:    pathData,
		Package: golang.MakePackageData(rootPath, rootPackage, "main"),
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
