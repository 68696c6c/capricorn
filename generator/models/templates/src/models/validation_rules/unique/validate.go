package unique

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var ruleValidateBodyTemplate = `
	{{ .Field.Name.Unexported }}, ok := value.({{ .Field.Type }})
	if !ok {
		return errors.New("invalid {{ .Single.Space }} {{ .Field.Name.Space }}")
	}

	query := {{ .Receiver }}{ .DB }}.First(&{{ .Single.Exported }}{
		{{ .Field.Name.Exported }}: {{ .Field.Name.Unexported }},
	})
	if !query.RecordNotFound() {
		return errors.New("{{ .Single.Space }} {{ .Field.Name.Space }} already exists")
	}

	return nil
`

type validate struct {
	Receiver string
	DB       string
	Field    module.ResourceField
	Single   data.Name
}

func (m validate) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         m.GetName(),
		Imports:      m.GetImports(),
		Arguments:    m.GetArgs(),
		ReturnValues: m.GetReturns(),
		Receiver:     m.GetReceiver(),
		Body:         m.MustParse(),
	}
}

func (m validate) GetName() string {
	return "Validate"
}

func (m validate) GetImports() golang.Imports {
	return golang.Imports{
		Standard: []string{},
		App:      []string{},
		Vendor:   []string{data.ImportErrors, data.ImportGorm},
	}
}

func (m validate) GetReceiver() golang.Value {
	return golang.Value{
		Name: m.Receiver,
	}
}

func (m validate) GetArgs() []golang.Value {
	return []golang.Value{
		{
			Name: m.DB,
			Type: "*gorm.DB",
		},
	}
}

func (m validate) GetReturns() []golang.Value {
	return []golang.Value{
		{
			Type: "error",
		},
	}
}

func (m validate) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_rule_unique_validate_body", ruleValidateBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
