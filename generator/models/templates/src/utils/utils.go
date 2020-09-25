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
	TypeData    data.TypeData
	Name        data.Name
	Resource    module.Resource

	Built      bool
	Imports    golang.Imports
	Interfaces []golang.Interface
	Structs    []golang.Struct
	Functions  []golang.Function
}

func NewService(meta ServiceMeta, receiverType string) Service {
	pkgData := meta.PackageData
	fileData, pathData := data.MakeGoFileData(pkgData.GetImport(), meta.FileName)
	return Service{
		FileData:    fileData,
		PathData:    pathData,
		PackageData: meta.PackageData,
		Name:        meta.Name,
		Resource:    meta.Resource,
		TypeData:    data.MakeTypeDataService(pkgData.Reference, meta.Name.Exported, receiverType, false),
		Receiver: golang.Value{
			Name: meta.ReceiverName,
			Type: receiverType,
		},
	}
}
