package utils

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
)

type ServiceMeta struct {
	Name         data.Name
	ReceiverName string
	FileName     string
	Resource     module.Resource
	PackageData  data.PackageData
}
