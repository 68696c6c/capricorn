package utils

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/pkg/errors"
)

const (
	ServiceTypeRepo    ServiceType = "repo"
	ServiceTypeService ServiceType = "service"
)

var validServiceTypes = map[string]ServiceType{"repo": ServiceTypeRepo, "service": ServiceTypeService}

type ServiceType string

func ServiceTypeFromString(input string) (ServiceType, error) {
	result, ok := validServiceTypes[input]
	if ok {
		return result, nil
	}
	return ServiceType(""), errors.Errorf("invalid ServiceType '%s'", input)
}

type ContainerFieldMeta struct {
	data.Name
	data.TypeData
	ServiceType
	golang.Field
	PackageImport string
	DomainKey     string
	Constructor   golang.Function
}

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
	Receiver    golang.Value

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
		TypeData:    data.MakeTypeDataService(pkgData.Reference, meta.Name.Exported, meta.ReceiverName, false),
		Receiver: golang.Value{
			Name: meta.ReceiverName,
			Type: receiverType,
		},
	}
}
