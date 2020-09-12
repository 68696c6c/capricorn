package src

import (
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
)

type Method interface {
	GetName() string
	MustParse() string
	GetImports() golang.Imports
	GetReceiver() golang.Value
	GetArgs() []golang.Value
	GetReturns() []golang.Value
}

type SourceFile interface {
	MustGetFile() golang.File
	GetStructs() []golang.Struct
	MustGetFunctions() []golang.Function
}

type SourceFunction interface {
	MustGetFunction() golang.Function
}
