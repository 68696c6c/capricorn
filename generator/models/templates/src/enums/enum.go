package enums

import "github.com/68696c6c/capricorn/generator/models/templates/golang"

type Enum interface {
	MustGetFile() golang.File
}
