package unique

import (
	"fmt"

	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/models/templates/src/validation_rules/utils"
)

type Unique struct {
	RuleName        string
	ConstructorName string
	Receiver        string
	DB              string
	Field           module.ResourceField
	Single          data.Name
}

func NewRule(dbArg, receiverName string, resourceSingleName data.Name, field module.ResourceField) Unique {
	ruleName, constName := utils.MakeRuleName(resourceSingleName, field)
	return Unique{
		RuleName:        ruleName,
		ConstructorName: constName,
		DB:              dbArg,
		Receiver:        receiverName,
		Field:           field,
		Single:          resourceSingleName,
	}
}

func (m Unique) GetConstructorCall() string {
	return fmt.Sprintf("%s(%s)", m.ConstructorName, m.DB)
}

func (m Unique) GetStructs() []golang.Struct {
	return []golang.Struct{
		{
			Name: m.RuleName,
			Fields: []golang.Field{
				// @TODO how exactly is the message property used?
				{
					Name: "message",
					Type: "string",
				},
				{
					Name: m.DB,
					Type: "*gorm.DB",
				},
			},
		},
	}
}

func (m Unique) MustGetFunctions() []golang.Function {
	c := constructor{
		RuleName: m.RuleName,
		Name:     m.ConstructorName,
		DB:       m.DB,
		Field:    m.Field,
		Single:   m.Single,
	}

	v := validate{
		Receiver: m.Receiver,
		DB:       m.DB,
		Field:    m.Field,
		Single:   m.Single,
	}

	return []golang.Function{
		c.MustGetFunction(),
		v.MustGetFunction(),
	}
}
