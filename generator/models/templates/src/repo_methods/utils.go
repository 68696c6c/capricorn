package repo_methods

import (
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
)

type MethodMeta struct {
	DBFieldName string
	Resource    module.Resource
	Receiver    golang.Value
}
