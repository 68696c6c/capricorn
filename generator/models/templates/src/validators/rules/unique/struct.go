package unique

import (
	"fmt"

	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/models/templates/src/validators/rules"
)

type Unique struct {
	ruleMeta rules.RuleMeta
}

func NewRule(dbArg, receiverName string, resourceSingleName data.Name, field module.ResourceField) rules.Rule {
	ruleName, constructorName := rules.MakeRuleName(resourceSingleName, field, "unique")
	return Unique{
		ruleMeta: rules.RuleMeta{
			RuleName:        ruleName,
			ConstructorName: constructorName,
			DBArgName:       dbArg,
			DBFieldName:     "db",
			Single:          resourceSingleName,
			Field:           field,
			Receiver: golang.Value{
				Name: receiverName,
				Type: "*" + ruleName,
			},
		},
	}
}

func (m Unique) GetUsage() string {
	return fmt.Sprintf("%s(%s)", m.ruleMeta.ConstructorName, m.ruleMeta.DBArgName)
}

func (m Unique) GetStructs() []golang.Struct {
	return []golang.Struct{
		{
			Name: m.ruleMeta.RuleName,
			Fields: []golang.Field{
				// @TODO how exactly is the message property used?
				{
					Name: "message",
					Type: "string",
				},
				{
					Name: m.ruleMeta.DBFieldName,
					Type: "*gorm.DB",
				},
			},
		},
	}
}

func (m Unique) MustGetFunctions() []golang.Function {
	c := newConstructor(m.ruleMeta)
	v := newValidate(m.ruleMeta)
	return []golang.Function{
		c.MustGetFunction(),
		v.MustGetFunction(),
	}
}
