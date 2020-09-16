package methods

import (
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
)

type Method interface {
	GetDbReference() string
	MustGetFunction() golang.Function
	GetImports() golang.Imports
	MustParse() string
}

type Meta struct {
	DBFieldName string
	Resource    module.Resource
	Receiver    golang.Value
	ModelType   string
}
