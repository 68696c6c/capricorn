package src

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
)

type Service interface {
	GetPackageData() data.PackageData
	GetName() data.Name
	GetInterface() golang.Interface
	GetConstructor() golang.Function
}

type Method interface {
	GetName() string
	MustParse() string
	GetImports() golang.Imports
	GetReceiver() golang.Value
	GetArgs() []golang.Value
	GetReturns() []golang.Value
}

type SourceFunction interface {
	MustGetFunction() golang.Function
}
