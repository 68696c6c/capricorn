package handlers

import (
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
)

type Handler interface {
	GetRepoReference() string
	GetErrorsReference() string
	MustGetFunction() golang.Function
	GetImports() golang.Imports
	MustParse() string
}

type Meta struct {
	CreateRequestType string
	UpdateRequestType string
	ViewResponseType  string
	ListResponseType  string
	Resource          module.Resource
	Receiver          golang.Value
	RepoField         golang.Value
	ErrorsField       golang.Value
	ContextValue      golang.Value
}
