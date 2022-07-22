package project

import "github.com/68696c6c/capricorn_rnd/generator/golang"

func MakeInitialSeeder() *golang.File {
	return golang.MakeGoFile("initial").SetImports(golang.Imports{
		Vendor: []golang.Package{PkgGoat, PkgGorm},
	}).SetFunctions([]*golang.Function{
		{
			Name:    "Initial",
			Imports: golang.Imports{},
			Arguments: []golang.Value{
				{
					Name: "_",
					Type: "*gorm.DB",
				},
			},
			ReturnValues: []golang.Value{
				{
					Type: "error",
				},
			},
			Body: "return nil",
		},
	})
}
