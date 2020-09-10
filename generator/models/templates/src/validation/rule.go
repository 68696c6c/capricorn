package validation

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var uniqueRuleBodyTemplate = `
	{{ .Field.Unexported }}, ok := value.({{ .Field.Type }})
	if !ok {
		return errors.New("invalid {{ .Single.Space }} {{ .Field.Space }}")
	}

	query := {{ .DB }}.First(&{{ .Single.Exported }}{
		{{ .Field.Exported }}: {{ .Field.Unexported }},
	})
	if !query.RecordNotFound() {
		return errors.New("{{ .Single.Space }} {{ .Field.Space }} already exists")
	}

	return nil
`

type UniqueRule struct {
	DB              string
	ConstructorName string
	Field           data.Name
	Single          data.Name
}

func MakeUniqueRule(dbArg string, singleName, fieldName data.Name) UniqueRule {
	return UniqueRule{
		DB:              dbArg,
		ConstructorName: makeRuleConstructorName(singleName, fieldName),
		Field:           fieldName,
		Single:          singleName,
	}
}

func (m UniqueRule) MustMakeFunction() golang.Function {
	return golang.Function{
		Name:    m.ConstructorName,
		Imports: golang.Imports{},
		Arguments: []golang.Value{
			{
				Name: m.DB,
				Type: "*gorm.DB",
			},
		},
		ReturnValues: []golang.Value{
			{
				Type: "error",
			},
		},
		Body: m.MustParse(),
	}
}

func (m UniqueRule) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_unique_rule_body", uniqueRuleBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
