package validators

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
)

type ValidationMeta struct {
	DBFieldName   string
	ReceiverName  string
	ModelName     data.Name
	ModelReceiver golang.Value
}
