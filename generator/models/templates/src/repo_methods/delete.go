package repo_methods

import (
	"fmt"

	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var deleteBodyTemplate = `
	errs :=  {{ .GetDbReference }}.Delete(m).GetErrors()
	if len(errs) > 0 {
		return goat.ErrorsToError(errs)
	}
	return nil
`

type Delete struct {
	dbFieldName string
	receiver    golang.Value
	Single      data.Name
}

func NewDelete(meta MethodMeta) Delete {
	return Delete{
		dbFieldName: meta.DBFieldName,
		receiver:    meta.Receiver,
		Single:      meta.Resource.Inflection.Single,
	}
}

func (m Delete) GetDbReference() string {
	return fmt.Sprintf("%s.%s", m.receiver.Name, m.dbFieldName)
}

func (m Delete) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         m.GetName(),
		Imports:      m.GetImports(),
		Receiver:     m.GetReceiver(),
		Arguments:    m.GetArgs(),
		ReturnValues: m.GetReturns(),
		Body:         m.MustParse(),
	}
}

func (m Delete) GetName() string {
	return "Delete"
}

func (m Delete) GetImports() golang.Imports {
	return golang.Imports{
		Standard: nil,
		App:      nil,
		Vendor:   []string{data.ImportGoat},
	}
}

func (m Delete) GetReceiver() golang.Value {
	return m.receiver
}

func (m Delete) GetArgs() []golang.Value {
	return []golang.Value{
		{
			Name: "m",
			Type: "*" + m.Single.Exported,
		},
	}
}

func (m Delete) GetReturns() []golang.Value {
	return []golang.Value{
		{
			Type: "error",
		},
	}
}

func (m Delete) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_repo_delete", deleteBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
