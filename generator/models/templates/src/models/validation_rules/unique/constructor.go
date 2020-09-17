package unique

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var ruleConstructorBodyTemplate = `
	return &{{ .RuleName }}{
		message: "{{ .Single.Space }} {{ .Field.Name.Space }} must be unique",
		db:      {{ .DB }},
	}
`

type constructor struct {
	RuleName string
	Name     string
	DB       string
	Field    module.ResourceField
	Single   data.Name
}

func (m constructor) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         m.GetName(),
		Imports:      m.GetImports(),
		Arguments:    m.GetArgs(),
		ReturnValues: m.GetReturns(),
		Body:         m.MustParse(),
	}
}

func (m constructor) GetName() string {
	return m.Name
}

func (m constructor) GetImports() golang.Imports {
	return golang.Imports{
		Standard: []string{},
		App:      []string{},
		Vendor:   []string{data.ImportGorm},
	}
}

func (m constructor) GetReceiver() golang.Value {
	return golang.Value{}
}

func (m constructor) GetArgs() []golang.Value {
	return []golang.Value{
		{
			Name: m.DB,
			Type: "*gorm.DB",
		},
	}
}

func (m constructor) GetReturns() []golang.Value {
	return []golang.Value{
		{
			Type: "error",
		},
	}
}

func (m constructor) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_rule_unique_constructor_body", ruleConstructorBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
