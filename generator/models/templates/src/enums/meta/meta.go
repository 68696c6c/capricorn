package meta

import "github.com/68696c6c/capricorn/generator/models/data"

type Meta struct {
	FileData         data.FileData
	PathData         data.PathData
	PackageData      data.PackageData
	TypeName         string
	TypeNameReadable string
	ReceiverName     string
	Values           []string
}
