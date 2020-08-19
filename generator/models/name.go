package models

import "github.com/68696c6c/capricorn/generator/utils"

type Name struct {
	Snake      string
	Kebob      string
	Exported   string
	Unexported string
}

func MakeName(base string) Name {
	return Name{
		Snake:      utils.SeparatedToSnake(base),
		Kebob:      utils.SeparatedToKebob(base),
		Exported:   utils.SeparatedToExported(base),
		Unexported: utils.SeparatedToUnexported(base),
	}
}
