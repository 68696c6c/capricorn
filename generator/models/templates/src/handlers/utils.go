package handlers

import (
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
)

type MethodMeta struct {
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
