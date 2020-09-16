package controllers

import "github.com/68696c6c/capricorn/generator/models/templates/golang"

func NewResourceResponse(responseType, modelType string) golang.Struct {
	return golang.Struct{
		Name: responseType,
		Fields: []golang.Field{
			{
				Name: modelType,
			},
		},
	}
}

func NewListResponse(responseType, modelType string) golang.Struct {
	return golang.Struct{
		Name: responseType,
		Fields: []golang.Field{
			{
				Name: "Data",
				Type: "[]*" + modelType,
				Tags: []golang.Tag{
					{
						Key:    "json",
						Values: []string{"data"},
					},
				},
			},
			{
				Type: "query.Pagination",
				Tags: []golang.Tag{
					{
						Key:    "json",
						Values: []string{"pagination"},
					},
				},
			},
		},
	}
}
