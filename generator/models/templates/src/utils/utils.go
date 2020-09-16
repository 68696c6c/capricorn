package utils

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
)

type ServiceMeta struct {
	Name         data.Name
	ReceiverName string
	FileName     string
	Resource     module.Resource
	PackageData  data.PackageData
	ModelType    string
}

type Service struct {
	FileData    data.FileData
	PathData    data.PathData
	PackageData data.PackageData
	Name        data.Name
	Resource    module.Resource
	Receiver    golang.Value

	Built      bool
	Imports    golang.Imports
	Interfaces []golang.Interface
	Structs    []golang.Struct
	Functions  []golang.Function
}

func NewService(meta ServiceMeta, receiverType string) Service {
	fileData, pathData := data.MakeGoFileData(meta.PackageData.GetImport(), meta.FileName)
	return Service{
		FileData:    fileData,
		PathData:    pathData,
		PackageData: meta.PackageData,
		Name:        meta.Name,
		Resource:    meta.Resource,
		Receiver: golang.Value{
			Name: meta.ReceiverName,
			Type: receiverType,
		},
	}
}
