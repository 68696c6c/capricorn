package controllers

import "github.com/68696c6c/capricorn/generator/models/templates/golang"

func NewRequestStruct(requestType, modelType string) golang.Struct {
	return golang.Struct{
		Name: requestType,
		Fields: []golang.Field{
			{
				Name: modelType,
			},
		},
	}
}
